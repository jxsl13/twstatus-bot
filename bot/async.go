package bot

import (
	"log"
	"time"

	"github.com/jxsl13/twstatus-bot/model"
)

// start this asynchronously once
func (b *Bot) serverUpdater(duration time.Duration) {
	var (
		timer   = time.NewTimer(0)
		drained = false
	)
	defer closeTimer(timer, &drained)
	for {
		select {
		case <-timer.C:
			drained = true
			// do something
			resetTimer(timer, duration, &drained)
			func() {
				_, _, err := b.updateServers()
				if err != nil {
					log.Printf("failed to update servers: %v", err)
					return
				}

				// publish changed servers
				err = b.changedServers()
				if err != nil {
					log.Printf("failed to get changed server messages from db: %v", err)
					return
				}
			}()
		case <-b.ctx.Done():
			log.Println("closed async goroutine for server and message updates")
			return
		}
	}
}

func (b *Bot) messageUpdater(id int) {
	log.Printf("goroutine %d starting async goroutine for message updates", id)

loop:
	for {
		select {
		case <-b.ctx.Done():
			break loop
		case server, ok := <-b.c:
			if !ok {
				break loop
			}
			err := b.updateDiscordMessage(server)
			if err != nil {
				log.Printf("goroutine %0d: failed to update discord message %v: %v", id, server.Target, err)
			}

		}
	}

	log.Printf("goroutine %d: closed async goroutine for message updates", id)
}

func (b *Bot) cacheCleanup() {
	var (
		cleanupInterval = 20 * b.pollingInterval
		timer           = time.NewTimer(cleanupInterval)
		drained         = false
	)
	defer closeTimer(timer, &drained)
	for {
		select {
		case <-timer.C:
			drained = true
			// do something
			resetTimer(timer, cleanupInterval, &drained)

			size := b.conflictMap.Size()
			if size == 0 {
				// nothing to do
				continue
			}

			now := time.Now()
			log.Printf("cache contains %d entries before cleanup at %s", size, now)
			b.conflictMap.Range(func(key model.MessageTarget, value Backoff) bool {
				// remove expired keys
				if now.After(value.Until) {
					b.conflictMap.Delete(key)
				}
				return true
			})
			log.Printf("cache contains %d entries after cleanup at %s", b.conflictMap.Size(), now)

		case <-b.ctx.Done():
			log.Println("closed async goroutine for cache cleanup")
			return
		}
	}
}
