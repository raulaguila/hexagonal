package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/raulaguila/go-api/pkg/loggerx/formatter"
)

// DatadogSink sends log entries to Datadog.
type DatadogSink struct {
	mu      sync.Mutex
	client  *http.Client
	apiKey  string
	site    string
	service string
	env     string
	source  string
	tags    []string
}

// DatadogConfig configures the Datadog sink.
type DatadogConfig struct {
	APIKey  string
	Site    string // e.g., "datadoghq.com", "datadoghq.eu"
	Service string
	Env     string
	Source  string
	Tags    []string
	Timeout time.Duration
}

// datadogLog is the Datadog log format.
type datadogLog struct {
	DDSource string         `json:"ddsource"`
	DDTags   string         `json:"ddtags"`
	Hostname string         `json:"hostname"`
	Message  string         `json:"message"`
	Service  string         `json:"service"`
	Status   string         `json:"status"`
	Time     string         `json:"time"`
	Caller   string         `json:"caller,omitempty"`
	Fields   map[string]any `json:"fields,omitempty"`
}

// NewDatadog creates a new Datadog sink.
func NewDatadog(cfg DatadogConfig) *DatadogSink {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	site := cfg.Site
	if site == "" {
		site = "datadoghq.com"
	}

	source := cfg.Source
	if source == "" {
		source = "go"
	}

	return &DatadogSink{
		client: &http.Client{
			Timeout: timeout,
		},
		apiKey:  cfg.APIKey,
		site:    site,
		service: cfg.Service,
		env:     cfg.Env,
		source:  source,
		tags:    cfg.Tags,
	}
}

// Write sends the entry to Datadog.
func (s *DatadogSink) Write(entry *formatter.Entry) error {
	// Map log level to Datadog status
	status := s.mapLevel(entry.Level)

	// Build tags
	var tags strings.Builder
	fmt.Fprintf(&tags, "env:%s", s.env)
	for _, tag := range s.tags {
		tags.WriteString("," + tag)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	log := datadogLog{
		DDSource: s.source,
		DDTags:   tags.String(),
		Hostname: hostname,
		Message:  entry.Message,
		Service:  s.service,
		Status:   status,
		Time:     entry.Time,
		Caller:   entry.Caller,
		Fields:   entry.Fields,
	}

	data, err := json.Marshal([]datadogLog{log})
	if err != nil {
		return fmt.Errorf("datadog: marshal error: %w", err)
	}

	// Build URL
	url := fmt.Sprintf("https://http-intake.logs.%s/api/v2/logs", s.site)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("datadog: create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", s.apiKey)

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// Close releases resources.
func (s *DatadogSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.client.CloseIdleConnections()
	return nil
}

func (s *DatadogSink) mapLevel(level string) string {
	switch level {
	case "DEBUG":
		return "debug"
	case "INFO":
		return "info"
	case "WARN":
		return "warn"
	case "ERROR":
		return "error"
	case "FATAL":
		return "critical"
	default:
		return "info"
	}
}
