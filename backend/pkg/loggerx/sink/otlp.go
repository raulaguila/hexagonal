package sink

import (
	"context"
	"fmt"
	"time"

	"github.com/raulaguila/go-api/pkg/loggerx/formatter"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

// OTLPSink sends log entries to the OpenTelemetry Logger Provider.
type OTLPSink struct {
	logger log.Logger
}

// NewOTLP creates a new OTLP sink.
// It uses the global logger provider by default.
func NewOTLP() *OTLPSink {
	return &OTLPSink{
		logger: global.Logger("github.com/raulaguila/go-api/pkg/loggerx"),
	}
}

// Write sends the entry to OTEL Logger.
func (s *OTLPSink) Write(entry *formatter.Entry) error {
	var severity log.Severity
	switch entry.Level {
	case "DEBUG":
		severity = log.SeverityDebug
	case "INFO":
		severity = log.SeverityInfo
	case "WARN":
		severity = log.SeverityWarn
	case "ERROR":
		severity = log.SeverityError
	case "FATAL":
		severity = log.SeverityFatal
	default:
		severity = log.SeverityInfo
	}

	// Convert fields to KeyValues
	var attributes []log.KeyValue
	for k, v := range entry.Fields {
		attributes = append(attributes, log.String(k, formatValue(v)))
	}
	if entry.Caller != "" {
		attributes = append(attributes, log.String("caller", entry.Caller))
	}

	record := log.Record{}
	record.SetTimestamp(time.Now())
	record.SetSeverity(severity)
	record.SetSeverityText(entry.Level)
	record.SetBody(log.StringValue(entry.Message))
	record.AddAttributes(attributes...)

	// If context has trace info, OTEL SDK automatically extracts it if using the right Context.
	// However, 'formatter.Entry' has a Context field. We should use it.
	ctx := entry.Context
	if ctx == nil {
		ctx = context.Background()
	}

	s.logger.Emit(ctx, record)
	return nil
}

// Close does nothing as the Provider is managed globally.
func (s *OTLPSink) Close() error {
	return nil
}

func formatValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
