package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DB       int
	TTL      time.Duration
}

// Service provides Redis operations
type Service struct {
	client *redis.Client
	TTL    time.Duration
}

// MustNew creates a new Redis service or panics
func MustNew(cfg Config) *Service {
	svc, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return svc
}

// New creates a new Redis service
func New(cfg Config) (*Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Service{client: client, TTL: cfg.TTL}, nil
}

// GetClient returns the underlying redis client
func (s *Service) GetClient() *redis.Client {
	return s.client
}

// Close closes the connection
func (s *Service) Close() error {
	return s.client.Close()
}

// Get retrieves a value from Redis and unmarshals it into dest
func (s *Service) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

// Set sets a value in Redis (marshaled to JSON) with specific expiration
func (s *Service) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, bytes, expiration).Err()
}

// Del deletes values from Redis
func (s *Service) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return s.client.Del(ctx, keys...).Err()
}
