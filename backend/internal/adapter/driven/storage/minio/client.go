package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds MinIO configuration
type Config struct {
	Url        string
	User       string
	Password   string
	BucketName string
}

// MustConnect establishes a connection or panics
func MustConnect(cfg *Config) *minio.Client {
	client, err := connect(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

// Connect establishes a connection to MinIO
func connect(cfg *Config) (*minio.Client, error) {
	client, err := minio.New(cfg.Url, &minio.Options{
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
