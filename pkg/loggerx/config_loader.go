package loggerx

import (
	"github.com/raulaguila/go-api/pkg/loggerx/config"
	"github.com/raulaguila/go-api/pkg/loggerx/formatter"
	"github.com/raulaguila/go-api/pkg/loggerx/sink"
)

// NewFromConfig creates a new logger from a YAML configuration file.
func NewFromConfig(path string) (*Logger, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return NewFromConfigStruct(cfg)
}

// NewFromConfigStruct creates a new logger from a config struct.
func NewFromConfigStruct(cfg *config.Config) (*Logger, error) {
	level := ParseLevel(cfg.Logger.Level)
	logger := New(
		WithLevel(level),
		WithCaller(cfg.Logger.AddCaller),
		WithTimeFormat(cfg.Logger.TimeFormat),
	)

	var fmt formatter.Formatter
	if cfg.Logger.Format == "json" {
		fmt = formatter.NewJSON()
	} else {
		fmt = formatter.NewText()
	}

	if cfg.Logger.Stdout.Enabled {
		var textFmt *formatter.TextFormatter
		if cfg.Logger.Format == "text" {
			textFmt = formatter.NewText()
			textFmt.DisableColors = cfg.Logger.Stdout.DisableColors
			logger.AddSink(sink.NewStdout(sink.WithFormatter(textFmt)))
		} else {
			logger.AddSink(sink.NewStdout(sink.WithFormatter(fmt)))
		}
	}

	if cfg.Logger.File.Enabled {
		fileSink, err := sink.NewFile(cfg.Logger.File.Path, sink.WithFileFormatter(fmt))
		if err != nil {
			return nil, err
		}
		logger.AddSink(fileSink)
	}

	if cfg.Logger.Elasticsearch.Enabled {
		logger.AddSink(sink.NewElasticsearch(sink.ElasticsearchConfig{
			URL:      cfg.Logger.Elasticsearch.URL,
			Index:    cfg.Logger.Elasticsearch.Index,
			Username: cfg.Logger.Elasticsearch.Username,
			Password: cfg.Logger.Elasticsearch.Password,
		}))
	}

	if cfg.Logger.Loki.Enabled {
		logger.AddSink(sink.NewLoki(sink.LokiConfig{
			URL:    cfg.Logger.Loki.URL,
			Labels: cfg.Logger.Loki.Labels,
		}))
	}

	if cfg.Logger.Datadog.Enabled {
		logger.AddSink(sink.NewDatadog(sink.DatadogConfig{
			APIKey:  cfg.Logger.Datadog.APIKey,
			Site:    cfg.Logger.Datadog.Site,
			Service: cfg.Logger.Datadog.Service,
			Env:     cfg.Logger.Datadog.Env,
			Source:  cfg.Logger.Datadog.Source,
			Tags:    cfg.Logger.Datadog.Tags,
		}))
	}

	return logger, nil
}
