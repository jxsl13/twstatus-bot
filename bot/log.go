package bot

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// started asynchronously
func (b *Bot) logWriter() {
	for {
		select {
		case logEntry := <-b.logChan:
			lvl := ""
			switch logEntry.Level {
			case slog.LevelError:
				lvl = " [ERROR]: "
			case slog.LevelWarn:
				lvl = " [WARN]: "
			case slog.LevelInfo:
				lvl = " [INFO]: "
			case slog.LevelDebug:
				lvl = " [DEBUG]: "
			}

			msg := fmt.Sprintf("%s%s",
				lvl,
				logEntry.Message,
			)
			log.Println(msg)
			_, err := b.state.SendMessageComplex(b.channelID, api.SendMessageData{
				Content: msg,
				Embeds:  logEntry.Embedding,
				Flags:   discord.SuppressEmbeds,
			})
			if err != nil {
				b.Errorf("failed to send log message: %v", err)
				continue
			}
		case <-b.ctx.Done():
			log.Println("closed async goroutine for log writer")
			return
		}
	}
}

type LogEntry struct {
	slog.Record
	Embedding []discord.Embed
}

func (b *Bot) Logf(level slog.Level, embed []discord.Embed, format string, args ...any) {
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

func (b *Bot) Errorf(format string, args ...any) {
	b.Logf(slog.LevelError, nil, format, args...)
}

func (b *Bot) Warnf(format string, args ...any) {
	b.Logf(slog.LevelWarn, nil, format, args...)
}

func (b *Bot) Infof(format string, args ...any) {
	b.Logf(slog.LevelInfo, nil, format, args...)
}

func (b *Bot) Debugf(format string, args ...any) {
	b.Logf(slog.LevelDebug, nil, format, args...)
}
