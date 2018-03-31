package chat

import (
	"context"
	"fmt"

	r "github.com/clarencejychan/console-chat/server/cmd/redis"
	pb "github.com/clarencejychan/consolechat-grpc/console-chat"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
)

// ChatServedisiceServer implements the protobuf interface
type ChatServiceServer struct {
	r *r.Redis
}

func Register(r *r.Redis) *ChatServiceServer {
	return &ChatServiceServer{
		r: r,
	}
}

// Connect adds a user into the redis servis
func (c *ChatServiceServer) Connect(ctx context.Context, req *pb.ConnectRequest) (*google_protobuf.Empty, error) {
	userKey := "online." + req.GetUser()
	err := c.r.Client.Publish("chat", "it's working bitches").Err()
	if err != nil {
		panic(err)
	}
	// Add Redis key and make sure that
	fmt.Println(userKey + " connected!")
	return nil, nil
}

func (c *ChatServiceServer) ListUsers(context.Context, *google_protobuf.Empty) (*pb.UserList, error) {
	return nil, nil
}

func (c *ChatServiceServer) ConsoleChat(pb.ChatService_ConsoleChatServer) error {
	return nil
}
