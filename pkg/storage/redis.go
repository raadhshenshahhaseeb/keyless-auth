package storage

// create redis connection
import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage interface {
	Save(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type Redis struct {
	Client *redis.Client
}

func NewRedisClient(ctx context.Context, config *redis.Options) (*Redis, error) {
	client := redis.NewClient(config)

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

func (r *Redis) Save(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return []byte(res), nil
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
