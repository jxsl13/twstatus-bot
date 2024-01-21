package logging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/jxsl13/twstatus-bot/markdown"
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
}

func (b *Logger) Consume() <-chan LogEntry {
	return b.logChan
}

func (b *Logger) Logf(level slog.Level, embed []discord.Embed, format string, args ...any) {
	select {
	case b.logChan <- LogEntry{
		Record: slog.Record{
			Time:    time.Now(),
			Level:   level,
			Message: fmt.Sprintf(format, args...),
		},
		Embedding: embed,
	}:
	case <-b.ctx.Done():
		return
	}
}

func (b *Logger) Errorf(format string, args ...any) {
	b.Logf(slog.LevelError, nil, format, args...)
}

func (b *Logger) Warnf(format string, args ...any) {
	b.Logf(slog.LevelWarn, nil, format, args...)
}

func (b *Logger) Infof(format string, args ...any) {
	b.Logf(slog.LevelInfo, nil, format, args...)
}

func (b *Logger) Debugf(format string, args ...any) {
	b.Logf(slog.LevelDebug, nil, format, args...)
}

func (b *Logger) DebugAnyf(obj any, format string, args ...any) {
	b.Logf(slog.LevelDebug, []discord.Embed{
		{
			Title: "Debug",
			Type:  discord.NormalEmbed,
			Fields: []discord.EmbedField{
				{
					Value: markdown.CodeHighlight("go", fmt.Sprintf("%#v", obj)),
				},
			},
		},
	}, format, args...)
}
