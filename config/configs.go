package config

import (
	"crypto/rsa"
	"embed"
	"encoding/base64"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

//go:embed locales/*
var Locales embed.FS

// Config holds all application configuration
type Config struct {
	// System
	Timezone      string
	Version       string
	Port          string
	Environment   string
	LogLevel      string
	LogFormat     string
	EnablePrefork bool
	EnableLogger  bool
	EnableSwagger bool

	// JWT
	AccessPrivateKey  *rsa.PrivateKey
	AccessExpiration  time.Duration
	RefreshPrivateKey *rsa.PrivateKey
	RefreshExpiration time.Duration

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// MinIO
	MinioHost       string
	MinioPort       string
	MinioUser       string
	MinioPassword   string
	MinioBucketName string
}

// Load loads configuration from environment
func Load() (*Config, error) {
	// Try to load .env file
	if err := godotenv.Load(path.Join("config", ".env")); err != nil {
		_, b, _, _ := runtime.Caller(0)
		_ = godotenv.Load(path.Join(path.Dir(b), "..", "..", "config", ".env"))
	}

	// Set timezone
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "America/Manaus"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	time.Local = loc

	// Parse JWT keys
	accessPrivateKey, err := parseRSAPrivateKey(os.Getenv("ACCESS_TOKEN"))
	if err != nil {
		return nil, err
	}

	refreshPrivateKey, err := parseRSAPrivateKey(os.Getenv("RFRESH_TOKEN"))
	if err != nil {
		return nil, err
	}

	accessExpiration, _ := parseDuration(os.Getenv("ACCESS_TOKEN_EXPIRE"), 15)
	refreshExpiration, _ := parseDuration(os.Getenv("RFRESH_TOKEN_EXPIRE"), 60)

	return &Config{
		// System
		Timezone:      tz,
		Version:       getEnv("SYS_VERSION", "1.0.0"),
		Port:          getEnv("API_PORT", "9000"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		LogFormat:     getEnv("LOG_FORMAT", "json"),
		EnablePrefork: os.Getenv("API_ENABLE_PREFORK") == "1",
		EnableLogger:  os.Getenv("API_LOGGER") == "1",
		EnableSwagger: os.Getenv("API_SWAGGO") == "1",

		// JWT
		AccessPrivateKey:  accessPrivateKey,
		AccessExpiration:  accessExpiration,
		RefreshPrivateKey: refreshPrivateKey,
		RefreshExpiration: refreshExpiration,

		// Database
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", "root"),
		DBPassword: getEnv("POSTGRES_PASS", "root"),
		DBName:     getEnv("POSTGRES_BASE", "api"),

		// MinIO
		MinioHost:       os.Getenv("MINIO_HOST"),
		MinioPort:       getEnv("MINIO_API_PORT", "9004"),
		MinioUser:       getEnv("MINIO_USER", "minio"),
		MinioPassword:   getEnv("MINIO_PASS", "miniopass"),
		MinioBucketName: getEnv("MINIO_BUCKET_FILES", "api"),
	}, nil
}

// MustLoad loads configuration or panics
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// parseRSAPrivateKey parses a base64-encoded RSA private key
func parseRSAPrivateKey(encoded string) (*rsa.PrivateKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(decoded)
}

// parseDuration parses a duration string in minutes
func parseDuration(str string, defaultMinutes int) (time.Duration, error) {
	if str == "" {
		return time.Duration(defaultMinutes) * time.Minute, nil
	}
	var minutes int
	if _, err := os.Stat(str); err == nil {
		return time.Duration(defaultMinutes) * time.Minute, nil
	}
	for _, c := range str {
		if c >= '0' && c <= '9' {
			minutes = minutes*10 + int(c-'0')
		}
	}
	if minutes == 0 {
		minutes = defaultMinutes
	}
	return time.Duration(minutes) * time.Minute, nil
}

// getEnv returns environment variable value or default
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
