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

// ElasticsearchSink sends log entries to Elasticsearch.
type ElasticsearchSink struct {
	mu       sync.Mutex
	client   *http.Client
	url      string
	index    string
	username string
	password string
}

// ElasticsearchConfig configures the Elasticsearch sink.
type ElasticsearchConfig struct {
	URL      string
	Index    string
	Username string
	Password string
	Timeout  time.Duration
}

// elasticDoc is the document format for Elasticsearch.
type elasticDoc struct {
	Timestamp string         `json:"@timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Caller    string         `json:"caller,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
}

// NewElasticsearch creates a new Elasticsearch sink.
func NewElasticsearch(cfg ElasticsearchConfig) *ElasticsearchSink {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &ElasticsearchSink{
		client: &http.Client{
			Timeout: timeout,
		},
		url:      cfg.URL,
		index:    cfg.Index,
		username: cfg.Username,
		password: cfg.Password,
	}
}

// Write sends the entry to Elasticsearch.
func (s *ElasticsearchSink) Write(entry *formatter.Entry) error {
	doc := elasticDoc{
		Timestamp: entry.Time,
		Level:     entry.Level,
		Message:   entry.Message,
		Caller:    entry.Caller,
		Fields:    entry.Fields,
	}

	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("elasticsearch: marshal error: %w", err)
	}

	// Build URL: POST /{index}/_doc
	url := fmt.Sprintf("%s/%s/_doc", s.url, s.index)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("elasticsearch: create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("elasticsearch: request error: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("elasticsearch: unexpected status: %d", resp.StatusCode)
	}

	return nil
}

// Close releases resources.
func (s *ElasticsearchSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.client.CloseIdleConnections()
	return nil
}
