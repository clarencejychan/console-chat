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
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:3000", "The server address in the format of host:port")
)

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

	opt := grpc.WithInsecure()
	conn, err := grpc.Dial(*serverAddr, opt)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	ctx := context.Background()
	client := pb.NewChatServiceClient(conn)
	test := &pb.ConnectRequest{
		User: "test",
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
		}
	}
}
