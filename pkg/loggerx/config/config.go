package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the complete logger configuration.
type Config struct {
	Logger LoggerConfig `json:"logger"`
}

// LoggerConfig contains all logger settings.
type LoggerConfig struct {
	// Level is the minimum log level (debug, info, warn, error, fatal)
	Level string `json:"level"`

	// Format is the output format (json, text)
	Format string `json:"format"`

	// TimeFormat is the time format for logs (default: RFC3339Nano)
	TimeFormat string `json:"time_format"`

	// AddCaller enables caller information
	AddCaller bool `json:"add_caller"`

	// Stdout configuration
	Stdout StdoutConfig `json:"stdout"`

	// File configuration
	File FileConfig `json:"file"`

	// Elasticsearch configuration
	Elasticsearch ElasticsearchConfig `json:"elasticsearch"`

	// Loki configuration
	Loki LokiConfig `json:"loki"`

	// Datadog configuration
	Datadog DatadogConfig `json:"datadog"`
}

// StdoutConfig configures stdout output.
type StdoutConfig struct {
	Enabled       bool `json:"enabled"`
	DisableColors bool `json:"disable_colors"`
}

// FileConfig configures file output.
type FileConfig struct {
	Enabled    bool   `json:"enabled"`
	Path       string `json:"path"`
	MaxSizeMB  int    `json:"max_size_mb"`
	MaxBackups int    `json:"max_backups"`
}

// ElasticsearchConfig configures Elasticsearch output.
type ElasticsearchConfig struct {
	Enabled  bool   `json:"enabled"`
	URL      string `json:"url"`
	Index    string `json:"index"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LokiConfig configures Grafana Loki output.
type LokiConfig struct {
	Enabled bool              `json:"enabled"`
	URL     string            `json:"url"`
	Labels  map[string]string `json:"labels"`
}

// DatadogConfig configures Datadog output.
type DatadogConfig struct {
	Enabled bool     `json:"enabled"`
	APIKey  string   `json:"api_key"`
	Site    string   `json:"site"`
	Service string   `json:"service"`
	Env     string   `json:"env"`
	Source  string   `json:"source"`
	Tags    []string `json:"tags"`
}

// Load reads and parses a configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Set defaults
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.Format == "" {
		cfg.Logger.Format = "text"
	}

	return &cfg, nil
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	// Validate level
	level := c.Logger.Level
	if level != "debug" && level != "info" && level != "warn" && level != "error" && level != "fatal" {
		return fmt.Errorf("invalid log level: %s", level)
	}

	// Validate format
	format := c.Logger.Format
	if format != "json" && format != "text" {
		return fmt.Errorf("invalid format: %s", format)
	}

	// Validate file path if enabled
	if c.Logger.File.Enabled && c.Logger.File.Path == "" {
		return fmt.Errorf("file path is required when file output is enabled")
	}

	// Validate Elasticsearch
	if c.Logger.Elasticsearch.Enabled {
		if c.Logger.Elasticsearch.URL == "" {
			return fmt.Errorf("elasticsearch URL is required")
		}
		if c.Logger.Elasticsearch.Index == "" {
			return fmt.Errorf("elasticsearch index is required")
		}
	}

	// Validate Loki
	if c.Logger.Loki.Enabled && c.Logger.Loki.URL == "" {
		return fmt.Errorf("loki URL is required")
	}

	// Validate Datadog
	if c.Logger.Datadog.Enabled {
		if c.Logger.Datadog.APIKey == "" {
			return fmt.Errorf("datadog API key is required")
		}
		if c.Logger.Datadog.Service == "" {
			return fmt.Errorf("datadog service name is required")
		}
	}

	return nil
}
