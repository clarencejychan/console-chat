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

// Connect adds a user into the redis service
func (c *ChatServiceServer) Connect(ctx context.Context, req *pb.ConnectRequest) (*google_protobuf.Empty, error) {
	var err error
	userKey := "online." + req.GetUser()

	// Check if user previously existed
	exists := c.r.KeyExists(userKey)
	if exists {
		err := errors.New("Someone is connected with your username, please choose another.")
		return nil, err
	}

	connectMsg := req.GetUser() + " connected."
	c.r.SetKey(userKey, req.GetUser())
	err = c.r.Publish(connectMsg)
	if err != nil {
		err := errors.New("Something went wrong when trying to publish your message.")
		return nil, err
	}

	fmt.Println(userKey + " connected!")
	return &google_protobuf.Empty{}, nil
}

// ListUsers lists all the current users in the room
func (c *ChatServiceServer) ListUsers(ctx context.Context, e *google_protobuf.Empty) (*pb.UserList, error) {
	return nil, nil
}

// ConsoleChat sends the messages and makes sure all subscribers receive it
func (c *ChatServiceServer) ConsoleChat(ctx context.Context, msg *pb.Message) (*google_protobuf.Empty, error) {
	var err error
	userKey := "online." + msg.GetUser()

	exists := c.r.KeyExists(userKey)
	if !exists {
		err := errors.New("Something went wrong with your key when sending a message")
		return nil, err
	}
	msgChat := "> " + msg.GetUser() + ": " + msg.GetSpeak()
	err = c.r.Publish(msgChat)

	if err != nil {
		err := errors.New("Something went wrong with your key when sending a message")
		return nil, err
	}

	return &google_protobuf.Empty{}, nil
}

// Disconnect removes a user from the set of keys in redis
func (c *ChatServiceServer) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*google_protobuf.Empty, error) {
	// Check if user already exists as a key, if it does delete it and unsubscribe it.
	userKey := "online." + req.GetUser()
	exists := c.r.KeyExists(userKey)
	if exists {
		exitMsg := "> " + req.GetUser() + " left the room"
		err := c.r.Publish(exitMsg)
		if err != nil {
			err := errors.New("Error when attempting to disconnect")
			return nil, err
		}

		err = c.r.DelKey(userKey)
		if err != nil {
			err := errors.New("Error when attempting to disconnect")
			return nil, err
		}
	}

	fmt.Println(userKey + " disconnected!")

	return &google_protobuf.Empty{}, nil
}
