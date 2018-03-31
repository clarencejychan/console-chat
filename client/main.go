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
	redisAddr  = "localhost:6379"
	serverAddr = flag.String("server_addr", "localhost:3000", "The server address in the format of host:port")
)

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
	} else {
		fmt.Println("Listening to redis on 6379")
	}
	return redisClient, nil
}

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

func handleInput(scanner *bufio.Scanner, msg chan string) {
	// Read Input
	for scanner.Scan() {
		m := scanner.Text()
		msg <- m
	}
}

func main() {
	// Connect gRPC
	flag.Parse()

	// Init redis client
	r, err := initRedis(redisAddr)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	// Set up pub-sub
	pubSub := r.Subscribe("chat")
	chat := make(chan string, 1)

	go listenForMessages(pubSub, chat)

	//pong, err := r.Ping().Result()
	//fmt.Println(pong, err)

	// Init gRPC Client
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(*serverAddr, opt)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	ctx := context.Background()
	client := pb.NewChatServiceClient(conn)

	// Test Input
	test := &pb.ConnectRequest{
		User: "Clarence",
	}
	client.Connect(ctx, test)

	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 1)

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	msg := make(chan string, 1)
	scanner := bufio.NewScanner(os.Stdin)

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
			fmt.Println("Echoing: ", s)
		case c := <-chat:
			fmt.Println(c)
		}
	}
}
