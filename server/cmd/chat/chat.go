package chat

import (
	"context"
	"fmt"

	pb "github.com/clarencejychan/consolechat-grpc/console-chat"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)

// ChatServiceServer implements the protobuf interface
type ChatServiceServer struct{}

func Register() *ChatServiceServer {
	return &ChatServiceServer{}
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
