package chat

import (
	"context"
	"fmt"

	pb "github.com/clarencejychan/consolechat-grpc/console-chat"
	"github.com/go-redis/redis"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)

// ChatServiceServer implements the protobuf interface
type ChatServiceServer struct {
	r *redis.Client
}

func Register(r *redis.Client) *ChatServiceServer {
	return &ChatServiceServer{
		r: r,
	}
}

// Connect adds a user into the redis servis
func (c *ChatServiceServer) Connect(context.Context, *pb.ConnectRequest) (*google_protobuf.Empty, error) {
	fmt.Println("Connected!")
	return nil, nil
}

func (c *ChatServiceServer) ListUsers(context.Context, *google_protobuf.Empty) (*pb.UserList, error) {
	return nil, nil
}

func (c *ChatServiceServer) ConsoleChat(pb.ChatService_ConsoleChatServer) error {
	return nil
}
