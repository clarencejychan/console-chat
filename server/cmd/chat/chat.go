package chat

import (
	"context"
	"errors"
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

	// Check if user previously existed
	e := c.r.Client.Exists(userKey)
	if e.Val() == int64(1) {
		err := errors.New("Someone is connected with your username, please choose another.")
		return nil, err
	}

	_ = c.r.Client.SetNX(userKey, req.GetUser(), 0)
	err := c.r.Client.Publish("chat", req.GetUser()+" connected.").Err()
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
	user := msg.GetUser()
	speak := msg.GetSpeak()
	userKey := "online." + user

	e := c.r.Client.Exists(userKey)
	if e.Val() == int64(1) {
		err := c.r.Client.Publish("chat", "> "+user+": "+speak).Err()
		if err != nil {
			panic(err)
		}
	}

	return &google_protobuf.Empty{}, nil
}

// Disconnect removes a user from the set of keys in redis
func (c *ChatServiceServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*google_protobuf.Empty, error) {
	return nil, nil
}
