package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/clarencejychan/consolechat-grpc/console-chat"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
)

var (
	user       = ""
	redisAddr  = "localhost:6379"
	serverAddr = flag.String("server_addr", "localhost:3000", "The server address in the format of host:port")
)

// Intialize redis client
func initRedis(addr string) (*redis.Client, error) {
	// Set up Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Listening to redis on 6379")
	return redisClient, nil
}

// Receive messages in the pub/sub model
func listenForMessages(p *redis.PubSub, chat chan string) {
	for {
		m, err := p.ReceiveMessage()
		if err != nil {
			panic(err)
		} else {
			chat <- m.Payload
		}
	}
}

// Handle text input
func handleInput(scanner *bufio.Scanner, msg chan string) {
	// Read Input
	for scanner.Scan() {
		//fmt.Println("\033[8m") // Hide input
		m := scanner.Text()
		msg <- m
	}
}

// Initialize user entering the room
func initUser(ctx context.Context, scanner *bufio.Scanner, client pb.ChatServiceClient, conn *grpc.ClientConn) {
	var n string

	for {
		fmt.Println("> Enter a username")
		scanner.Scan()
		n = scanner.Text()
		user = n
		req := &pb.ConnectRequest{
			User: n,
		}
		_, err := client.Connect(ctx, req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		break
	}
}

func main() {
	flag.Parse()

	msg := make(chan string, 1)
	chat := make(chan string, 1)
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)

	// Init gRPC Client
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(*serverAddr, opt)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connecting to gRPC server on port 3000")

	// Init redis client
	r, err := initRedis(redisAddr)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	// Set up pub-sub
	pubSub := r.Subscribe("chat")
	go listenForMessages(pubSub, chat)

	fmt.Println("--------------------")

	ctx := context.Background()

	// Initialize ChatServiceClient Handler
	client := pb.NewChatServiceClient(conn)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// New Scanner
	scanner := bufio.NewScanner(os.Stdin)

	// Init User:
	initUser(ctx, scanner, client, conn)

	// Handle seperate thread in lightweight go routine
	go handleInput(scanner, msg)

	// Loop for Messaging
loop:
	for {
		select {
		case <-sigs:
			// Do things to exit client
			fmt.Println("Got shutdown, exiting")

			// Break out of the outer for statement and end the program
			break loop
		case s := <-msg:
			req := &pb.Message{
				User:  user,
				Speak: s,
			}
			_, err := client.ConsoleChat(ctx, req)

			if err != nil {
				panic(err)
			}
		case c := <-chat:
			fmt.Println(c)
		}
	}
}
