package redis

import (
	"errors"
	"fmt"

	r "github.com/go-redis/redis"
)

// Redis struct
type Redis struct {
	Client *r.Client
	PubSub *r.PubSub
}

// InitRedis initializes redis client.
func InitRedis(port string) (*Redis, error) {
	fmt.Println("Starting redis client listening on 6379")
	opts := &r.Options{
		Addr:     port,
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	client := r.NewClient(opts)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Redis{Client: client}, nil
}

// KeyExists checks if key exists in Redis
func (r *Redis) KeyExists(key string) bool {
	v := r.Client.Exists(key)
	if v.Val() == int64(1) {
		return true
	}
	return false
}

// SetKey sets the key value pair into redis
func (r *Redis) SetKey(key, value string) {
	r.Client.SetNX(key, value, 0)
}

// DelKey deletes the key from redis.
func (r *Redis) DelKey(key string) error {
	v := r.Client.Del(key)
	if v.Val() != int64(1) {
		err := errors.New("User already disconnected or key doesn't exist")
		return err
	}
	return nil
}

// Publish sends the value to the pub/sub channel
func (r *Redis) Publish(msg string) error {
	err := r.Client.Publish("chat", msg).Err()
	if err != nil {
		return err
	}
	return nil
}

// InitPubSubChannel initializes the Pub/Sub Channel
func (r *Redis) InitPubSubChannel() {
	r.PubSub = r.Client.Subscribe("chat")
}
