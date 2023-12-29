package bot

import (
	"log"
	"time"
)

func (b *Bot) async(duration time.Duration) {
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
				startTotal := time.Now()
				_, _, err := b.updateServers(b.ctx)
				if err != nil {
					log.Printf("failed to update servers: %v", err)
					return
				}
				_, err = b.updateDiscordMessages(b.ctx)
				if err != nil {
					log.Printf("failed to update discord messages: %v", err)
					return
				}
				log.Printf("Total time needed for updating server list and message list: %s", time.Since(startTotal))
			}()
		case <-b.ctx.Done():
			log.Println("closed async goroutine for server and message updates")
			return
		}

	}
}

// closeTimer should be used as a deferred function
// in order to cleanly shut down a timer
func closeTimer(timer *time.Timer, drained *bool) {
	if drained == nil {
		panic("drained bool pointer is nil")
	}
	if !timer.Stop() {
		if *drained {
			return
		}
		<-timer.C
		*drained = true
	}
}

// resetTimer sets drained to false after resetting the timer.
func resetTimer(timer *time.Timer, duration time.Duration, drained *bool) {
	if drained == nil {
		panic("drained bool pointer is nil")
	}
	if !timer.Stop() {
		if !*drained {
			<-timer.C
		}
	}
	timer.Reset(duration)
	*drained = false
}
