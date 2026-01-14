package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// ElasticsearchConfig holds Elasticsearch connection settings
type ElasticsearchConfig struct {
	// URL of the Elasticsearch endpoint (e.g., "http://localhost:9200")
	URL string

	// Index name pattern (e.g., "logs-api-%s" where %s is date)
	IndexPattern string

	// Username for authentication (optional)
	Username string

	// Password for authentication (optional)
	Password string

	// BatchSize number of logs to batch before sending
	BatchSize int

	// FlushInterval interval to flush logs even if batch is not full
	FlushInterval time.Duration

	// Timeout for HTTP requests
	Timeout time.Duration

	// RetryCount number of retries on failure
	RetryCount int
}

// DefaultElasticsearchConfig returns default Elasticsearch configuration
func DefaultElasticsearchConfig() ElasticsearchConfig {
	return ElasticsearchConfig{
		URL:           "http://localhost:9200",
		IndexPattern:  "logs-api-%s",
		BatchSize:     100,
		FlushInterval: 5 * time.Second,
		Timeout:       10 * time.Second,
		RetryCount:    3,
	}
}

// ElasticsearchHandler implements slog.Handler for Elasticsearch
type ElasticsearchHandler struct {
	config    ElasticsearchConfig
	client    *http.Client
	buffer    []map[string]any
	mu        sync.Mutex
	level     slog.Level
	attrs     []slog.Attr
	groups    []string
	stopChan  chan struct{}
	flushChan chan struct{}
	enabled   bool
}

// NewElasticsearchHandler creates a new Elasticsearch handler
func NewElasticsearchHandler(cfg ElasticsearchConfig) *ElasticsearchHandler {
	h := &ElasticsearchHandler{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		buffer:    make([]map[string]any, 0, cfg.BatchSize),
		level:     LevelInfo,
		stopChan:  make(chan struct{}),
		flushChan: make(chan struct{}),
		enabled:   true,
	}

	// Start background flusher
	go h.backgroundFlusher()

	return h
}

// Enabled returns whether the handler is enabled for the given level
func (h *ElasticsearchHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.enabled && level >= h.level
}

// Handle processes a log record
func (h *ElasticsearchHandler) Handle(_ context.Context, r slog.Record) error {
	if !h.enabled {
		return nil
	}

	doc := make(map[string]any)

	// Add timestamp in Elasticsearch format
	doc["@timestamp"] = r.Time.Format(time.RFC3339Nano)
	doc["level"] = r.Level.String()
	doc["message"] = r.Message

	// Add groups as nested objects
	current := doc
	for _, g := range h.groups {
		nested := make(map[string]any)
		current[g] = nested
		current = nested
	}

	// Add pre-set attributes
	for _, attr := range h.attrs {
		h.addAttr(current, attr)
	}

	// Add record attributes
	r.Attrs(func(a slog.Attr) bool {
		h.addAttr(current, a)
		return true
	})

	// Add to buffer
	h.mu.Lock()
	h.buffer = append(h.buffer, doc)
	shouldFlush := len(h.buffer) >= h.config.BatchSize
	h.mu.Unlock()

	if shouldFlush {
		select {
		case h.flushChan <- struct{}{}:
		default:
		}
	}

	return nil
}

// addAttr adds an attribute to the document
func (h *ElasticsearchHandler) addAttr(doc map[string]any, attr slog.Attr) {
	if attr.Equal(slog.Attr{}) {
		return
	}

	key := attr.Key
	val := attr.Value

	switch val.Kind() {
	case slog.KindGroup:
		group := make(map[string]any)
		for _, a := range val.Group() {
			h.addAttr(group, a)
		}
		if len(group) > 0 {
			doc[key] = group
		}
	case slog.KindTime:
		doc[key] = val.Time().Format(time.RFC3339Nano)
	case slog.KindDuration:
		doc[key] = val.Duration().Nanoseconds()
	default:
		doc[key] = val.Any()
	}
}

// WithAttrs returns a new handler with additional attributes
func (h *ElasticsearchHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &ElasticsearchHandler{
		config:    h.config,
		client:    h.client,
		buffer:    h.buffer,
		level:     h.level,
		attrs:     append(h.attrs, attrs...),
		groups:    h.groups,
		stopChan:  h.stopChan,
		flushChan: h.flushChan,
		enabled:   h.enabled,
	}
	return newHandler
}

// WithGroup returns a new handler with a group
func (h *ElasticsearchHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newHandler := &ElasticsearchHandler{
		config:    h.config,
		client:    h.client,
		buffer:    h.buffer,
		level:     h.level,
		attrs:     h.attrs,
		groups:    append(h.groups, name),
		stopChan:  h.stopChan,
		flushChan: h.flushChan,
		enabled:   h.enabled,
	}
	return newHandler
}

// backgroundFlusher periodically flushes the buffer
func (h *ElasticsearchHandler) backgroundFlusher() {
	ticker := time.NewTicker(h.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.Flush()
		case <-h.flushChan:
			h.Flush()
		case <-h.stopChan:
			h.Flush()
			return
		}
	}
}

// Flush sends buffered logs to Elasticsearch
func (h *ElasticsearchHandler) Flush() error {
	h.mu.Lock()
	if len(h.buffer) == 0 {
		h.mu.Unlock()
		return nil
	}

	docs := make([]map[string]any, len(h.buffer))
	copy(docs, h.buffer)
	h.buffer = h.buffer[:0]
	h.mu.Unlock()

	return h.sendBulk(docs)
}

// sendBulk sends documents to Elasticsearch using bulk API
func (h *ElasticsearchHandler) sendBulk(docs []map[string]any) error {
	if len(docs) == 0 {
		return nil
	}

	indexName := fmt.Sprintf(h.config.IndexPattern, time.Now().Format("2006.01.02"))

	var buf bytes.Buffer
	for _, doc := range docs {
		// Action line
		action := map[string]any{
			"index": map[string]any{
				"_index": indexName,
			},
		}
		actionJSON, err := json.Marshal(action)
		if err != nil {
			continue
		}
		buf.Write(actionJSON)
		buf.WriteByte('\n')

		// Document line
		docJSON, err := json.Marshal(doc)
		if err != nil {
			continue
		}
		buf.Write(docJSON)
		buf.WriteByte('\n')
	}

	// Send to Elasticsearch
	url := fmt.Sprintf("%s/_bulk", h.config.URL)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-ndjson")

	if h.config.Username != "" {
		req.SetBasicAuth(h.config.Username, h.config.Password)
	}

	var lastErr error
	for i := 0; i <= h.config.RetryCount; i++ {
		resp, err := h.client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		lastErr = fmt.Errorf("elasticsearch returned status %d", resp.StatusCode)
	}

	return lastErr
}

// Stop stops the handler and flushes remaining logs
func (h *ElasticsearchHandler) Stop() error {
	h.enabled = false
	close(h.stopChan)
	return h.Flush()
}

// SetLevel sets the minimum log level
func (h *ElasticsearchHandler) SetLevel(level slog.Level) {
	h.level = level
}
