package sink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/raulaguila/go-api/pkg/loggerx/formatter"
)

// LokiSink sends log entries to Grafana Loki.
type LokiSink struct {
	mu     sync.Mutex
	client *http.Client
	url    string
	labels map[string]string
}

// LokiConfig configures the Loki sink.
type LokiConfig struct {
	URL     string
	Labels  map[string]string
	Timeout time.Duration
}

// lokiPushRequest is the Loki push API format.
type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// NewLoki creates a new Loki sink.
func NewLoki(cfg LokiConfig) *LokiSink {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	labels := cfg.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	return &LokiSink{
		client: &http.Client{
			Timeout: timeout,
		},
		url:    cfg.URL,
		labels: labels,
	}
}

// Write sends the entry to Loki.
func (s *LokiSink) Write(entry *formatter.Entry) error {
	// Build labels (include level)
	labels := make(map[string]string)
	for k, v := range s.labels {
		labels[k] = v
	}
	labels["level"] = entry.Level

	// Build log line
	logLine := entry.Message
	if len(entry.Fields) > 0 {
		fieldsJSON, _ := json.Marshal(entry.Fields)
		logLine = fmt.Sprintf("%s %s", entry.Message, string(fieldsJSON))
	}

	// Parse time and get nanoseconds
	t, err := time.Parse(time.RFC3339Nano, entry.Time)
	if err != nil {
		t = time.Now()
	}
	timestamp := fmt.Sprintf("%d", t.UnixNano())

	push := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: labels,
				Values: [][]string{
					{timestamp, logLine},
				},
			},
		},
	}

	data, err := json.Marshal(push)
	if err != nil {
		return fmt.Errorf("loki: marshal error: %w", err)
	}

	// Build URL: POST /loki/api/v1/push
	url := fmt.Sprintf("%s/loki/api/v1/push", s.url)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("loki: create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("loki: request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("loki: unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// Close releases resources.
func (s *LokiSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.client.CloseIdleConnections()
	return nil
}
