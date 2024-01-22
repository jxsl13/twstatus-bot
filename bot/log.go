package bot

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
)

// started asynchronously
func (b *Bot) logWriter() {
	for {
		select {
		case logEntry := <-b.l.Consume():
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
				b.l.Errorf("failed to send log message to %s: %v", b.channelID, err)
				continue
			}
		case <-b.ctx.Done():
			log.Println("closed async goroutine for log writer")
			return
		}
	}
}
