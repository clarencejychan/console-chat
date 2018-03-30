package redis

import (
	"fmt"

	r "github.com/go-redis/redis"
)

// InitRedis initializes redis client.
func InitRedis(port string) (*r.Client, error) {
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
	return client, nil
}
