package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	chat "github.com/clarencejychan/console-chat/server/cmd/chat"
	pb "github.com/clarencejychan/consolechat-grpc/console-chat"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 3000, "server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	fmt.Printf("Starting console chat on port %d...\n", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, chat.Register())
	grpcServer.Serve(lis)
}
