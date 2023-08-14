package redis

import (
	"TESTGO/pkg/database"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

var SeeksterTokenNamespace = "seekster:access_token"

type RealRedisClient struct {
	client *redis.Client
}

func NewRedisClient() database.RedisClientInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	fmt.Println("InitializeRedis")
	return &RealRedisClient{client: client}
}

/*
	func InitializeRedis() {
		fmt.Println("InitializeRedis")
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	}
*/
func GetSeeksterTokenNamespace() string {
	return SeeksterTokenNamespace
}

func (r *RealRedisClient) GetSeeksterToken(ctx context.Context, ssoid string) (string, error) {
	seeksterToken, err := r.client.Get(ctx, SeeksterTokenNamespace+":"+ssoid).Result()
	if err != nil {
		return "", err
	}
	return seeksterToken, nil
}

func (r *RealRedisClient) SetSeeksterToken(ctx context.Context, ssoid string, seeksterToken string) error {
	err := r.client.Set(ctx, SeeksterTokenNamespace+":"+ssoid, seeksterToken, 24*7*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RealRedisClient) CloseRedis() {
	_ = r.client.Close()
}
