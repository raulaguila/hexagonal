package formatter

import (
	"bytes"
	"fmt"
)

// TextFormatter formats log entries as human-readable text.
type TextFormatter struct {
	// DisableColors disables ANSI color output.
	DisableColors bool

	// DisableTimestamp hides the timestamp.
	DisableTimestamp bool

	// DisableCaller hides the caller information.
	DisableCaller bool

	// TimestampFormat is the format for timestamps.
	TimestampFormat string
}

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

// Format formats the entry as text.
func (f *TextFormatter) Format(entry *Entry) ([]byte, error) {
	var buf bytes.Buffer

	// Timestamp
	if !f.DisableTimestamp && entry.Time != "" {
		buf.WriteString(f.colorize(colorGray, entry.Time))
		buf.WriteString(" ")
	}

	// Level
	levelColor := f.getLevelColor(entry.Level)
	buf.WriteString(f.colorize(levelColor, fmt.Sprintf("%-5s", entry.Level)))
	buf.WriteString(" ")

	// Caller
	if !f.DisableCaller && entry.Caller != "" {
		buf.WriteString(f.colorize(colorCyan, entry.Caller))
		buf.WriteString(" ")
	}

	// Message
	buf.WriteString(entry.Message)

	// Fields
	if len(entry.Fields) > 0 {
		buf.WriteString(" ")
		first := true
		for k, v := range entry.Fields {
			if !first {
				buf.WriteString(" ")
			}
			buf.WriteString(f.colorize(colorBlue, k))
			buf.WriteString("=")
			buf.WriteString(fmt.Sprintf("%v", v))
			first = false
		}
	}

	buf.WriteString("\n")
	return buf.Bytes(), nil
}

func (f *TextFormatter) colorize(color, text string) string {
	if f.DisableColors {
		return text
	}
	return color + text + colorReset
}

func (f *TextFormatter) getLevelColor(level string) string {
	switch level {
	case "DEBUG":
		return colorGray
	case "INFO":
		return colorGreen
	case "WARN":
		return colorYellow
	case "ERROR":
		return colorRed
	case "FATAL":
		return colorPurple
	default:
		return colorReset
	}
}

// NewText creates a new text formatter.
func NewText() *TextFormatter {
	return &TextFormatter{}
}

// NewTextNoColor creates a new text formatter without colors.
func NewTextNoColor() *TextFormatter {
	return &TextFormatter{DisableColors: true}
}
