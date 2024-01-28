package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type Logger struct {
	ctx     context.Context
	logChan chan LogEntry
}

func NewLogger(ctx context.Context) *Logger {
	return &Logger{
		ctx:     ctx,
		logChan: make(chan LogEntry, 1024),
	}
}

type LogEntry struct {
	slog.Record
	Embedding []discord.Embed
	Files     []sendpart.File
}

func (b *Logger) Consume() <-chan LogEntry {
	return b.logChan
}

func (b *Logger) Logf(level slog.Level, format string, args ...any) {
	select {
	case b.logChan <- LogEntry{
		Record: slog.Record{
			Time:    time.Now(),
			Level:   level,
			Message: fmt.Sprintf(format, args...),
		},
	}:
	case <-b.ctx.Done():
		return
	}
}

func (b *Logger) Errorf(format string, args ...any) {
	b.Logf(slog.LevelError, format, args...)
}

func (b *Logger) Warnf(format string, args ...any) {
	b.Logf(slog.LevelWarn, format, args...)
}

func (b *Logger) Infof(format string, args ...any) {
	b.Logf(slog.LevelInfo, format, args...)
}

func (b *Logger) Debugf(format string, args ...any) {
	b.Logf(slog.LevelDebug, format, args...)
}

func (b *Logger) DebugAnyf(obj any, format string, args ...any) {
	buf := &bytes.Buffer{}
	buf.Grow(1014 * 1024)

	switch o := obj.(type) {
	case []byte:
		_, _ = buf.Write(o)
	case json.RawMessage:
		_, _ = buf.Write(o)
	default:
		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")
		_ = enc.Encode(obj)
	}

	select {
	case b.logChan <- LogEntry{
		Record: slog.Record{
			Time:    time.Now(),
			Level:   slog.LevelDebug,
			Message: fmt.Sprintf(format, args...),
		},
		Files: []sendpart.File{
			{
				Name:   "debug.json",
				Reader: buf,
			},
		}}:
	case <-b.ctx.Done():
		return
	}
}
