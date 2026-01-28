package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Service provides Redis operations
type Service struct {
	client *redis.Client
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
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Service{client: client}, nil
}

// GetClient returns the underlying redis client
func (s *Service) GetClient() *redis.Client {
	return s.client
}

// Close closes the connection
func (s *Service) Close() error {
	return s.client.Close()
}
