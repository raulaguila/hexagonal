package minio

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds MinIO configuration
type Config struct {
	Host       string
	Port       string
	User       string
	Password   string
	BucketName string
}

// NewConfigFromEnv creates a Config from environment variables
func NewConfigFromEnv() Config {
	return Config{
		Host:       os.Getenv("MINIO_HOST"),
		Port:       os.Getenv("MINIO_API_PORT"),
		User:       os.Getenv("MINIO_USER"),
		Password:   os.Getenv("MINIO_PASS"),
		BucketName: os.Getenv("MINIO_BUCKET_FILES"),
	}
}

// Connect establishes a connection to MinIO
func Connect(cfg Config) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to minio: %w", err)
	}

	// Initialize bucket
	if err := initBucket(client, cfg.BucketName); err != nil {
		return nil, err
	}

	return client, nil
}

// MustConnect establishes a connection or panics
func MustConnect(cfg Config) *minio.Client {
	client, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

// initBucket ensures the bucket exists and has versioning enabled
func initBucket(client *minio.Client, bucketName string) error {
	ctx := context.Background()

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	if err := client.EnableVersioning(ctx, bucketName); err != nil {
		return fmt.Errorf("failed to enable versioning: %w", err)
	}

	return nil
}
