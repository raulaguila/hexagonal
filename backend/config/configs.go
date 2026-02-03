package config

import (
	"crypto/rand"
	"crypto/rsa"
	"embed"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/raulaguila/go-api/pkg/envx"
)

//go:embed locales/*
var Locales embed.FS

// Config holds all application configuration
type Environment struct {
	// System
	Timezone    *time.Location `env:"TZ" default:"America/Manaus"`
	ServiceName string         `env:"SYS_NAME" default:"API Backend"`
	Version     string         `env:"SYS_VERSION" default:"1.0.0"`
	Environment string         `env:"SYS_ENVIRONMENT" default:"production"`

	LogFormat string `env:"LOG_FORMAT" default:"json"`
	LogLevel  string `env:"LOG_LEVEL" default:"info"`

	Port          int  `env:"API_PORT" default:"9999"`
	EnableLogger  bool `env:"API_LOGGER" default:"1"`
	EnableSwagger bool `env:"API_SWAGGO" default:"1"`
	EnablePrefork bool `env:"API_PREFORK" default:"1"`

	// JWT
	AccessPrivateKey  *rsa.PrivateKey `env:"ACCESS_TOKEN" default:"new"`
	AccessExpiration  time.Duration   `env:"ACCESS_TOKEN_EXPIRE" default:"15m"`
	RefreshPrivateKey *rsa.PrivateKey `env:"RFRESH_TOKEN" default:"new"`
	RefreshExpiration time.Duration   `env:"RFRESH_TOKEN_EXPIRE" default:"60m"`

	// Database
	PGHost     string `env:"POSTGRES_HOST" default:"postgres"`
	PGPort     int    `env:"POSTGRES_PORT" default:"5438"`
	PGUser     string `env:"POSTGRES_USER" default:"root"`
	PGPassword string `env:"POSTGRES_PASS" default:"root"`
	PGBase     string `env:"POSTGRES_BASE" default:"api"`
	PGUrl      string `env:"POSTGRES_URL" default:"host=${POSTGRES_HOST} user=${POSTGRES_USER} password=${POSTGRES_PASS} dbname=${POSTGRES_BASE} port=${POSTGRES_PORT} sslmode=disable TimeZone=${TZ}"`

	// Redis
	RedisHost string        `env:"REDIS_HOST" default:"redis"`
	RedisPort int           `env:"REDIS_PORT" default:"6379"`
	RedisUser string        `env:"REDIS_USER" default:"default"`
	RedisPass string        `env:"REDIS_PASS" default:""`
	RedisDB   int           `env:"REDIS_DB" default:"0"`
	RedisTTL  time.Duration `env:"REDIS_TTL" default:"10m"`

	// MinIO
	MinioHost       string `env:"MINIO_HOST" default:"localhost"`
	MinioPort       int    `env:"MINIO_API_PORT" default:"9004"`
	MinioUrl        string `env:"MINIO_URL" default:"${MINIO_HOST}:${MINIO_API_PORT}"`
	MinioUser       string `env:"MINIO_USER" default:"minio"`
	MinioPassword   string `env:"MINIO_PASS" default:"miniopass"`
	MinioBucketName string `env:"MINIO_BUCKET" default:"api"`

	// OpenTelemetry
	OtelExporterOtlpEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT" default:"localhost:4317"`
}

func init() {
	// Register parser for *rsa.PrivateKey
	envx.RegisterParser(func(s string) (*rsa.PrivateKey, error) {
		if s == "new" {
			return rsa.GenerateKey(rand.Reader, 2048)
		}
		decoded, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return nil, err
		}
		t, err := jwt.ParseRSAPrivateKeyFromPEM(decoded)
		if err != nil {
			return nil, err
		}
		return t, nil
	})
}

// MustLoad loads environment or panics
func MustLoad() *Environment {
	cfg, err := load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Load loads configuration from environment
func load() (*Environment, error) {
	if err := envx.LoadDotEnvOverride(path.Join("config", ".env")); err != nil {
		_, b, _, _ := runtime.Caller(0)
		if err := envx.LoadDotEnvOverride(path.Join(path.Dir(b), "..", "config", ".env")); err != nil {
			fmt.Printf("Failed to load environment file: %v\n", err)
			os.Exit(1)
		}
	}

	env := &Environment{}
	if err := envx.Load(env); err != nil {
		fmt.Printf("Failed to parse environments: %v\n", err)
		os.Exit(1)
	}

	time.Local = env.Timezone

	return env, nil
}
