package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	chat "github.com/clarencejychan/console-chat/server/cmd/chat"
	r "github.com/clarencejychan/console-chat/server/cmd/redis"
	pb "github.com/clarencejychan/consolechat-grpc/console-chat"
	"google.golang.org/grpc"
)

var (
	port      = flag.Int("port", 3000, "server port")
	redisPort = "localhost:6379"
)

func main() {

	flag.Parse()

	// Initialize Redis
	redis, err := r.InitRedis(redisPort)
	if err != nil {
		log.Fatalf("failed to connect to redis client on %s", redisPort)
	}

	// Subscribe to RedisPubSub
	redis.InitPubSubChannel()

	// Initialize gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	fmt.Printf("Starting console chat on port %d...\n", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// Register Service Methods
	pb.RegisterChatServiceServer(grpcServer, chat.Register(redis))

	// Serve gRPC server
	grpcServer.Serve(lis)
}
