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

// Register binds the redis client to ChatServiceServer
func Register(r *r.Redis) *ChatServiceServer {
	return &ChatServiceServer{
		r: r,
	}
}

// Connect adds a user into the redis servis
func (c *ChatServiceServer) Connect(ctx context.Context, req *pb.ConnectRequest) (*google_protobuf.Empty, error) {
	userKey := "online." + req.GetUser()

	err := c.r.Client.Publish("chat", req.GetUser()+" connected.").Err()
	_ = c.r.Client.SetNX(userKey, req.GetUser(), 0)
	if err != nil {
		panic(err)
	}
	// Add Redis key and make sure that
	fmt.Println(userKey + " connected!")
	return &google_protobuf.Empty{}, nil
}

// ListUsers lists all the current users in the room
func (c *ChatServiceServer) ListUsers(ctx context.Context, e *google_protobuf.Empty) (*pb.UserList, error) {
	return nil, nil
}

// ConsoleChat sends the messages and makes sure all subscribers receive it
func (c *ChatServiceServer) ConsoleChat(ctx context.Context, msg *pb.Message) (*google_protobuf.Empty, error) {
	return nil, nil
}

// Disconnect removes a user from the set of keys in redis
func (c *ChatServiceServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*google_protobuf.Empty, error) {
	return nil, nil
}
