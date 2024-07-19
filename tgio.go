package tgio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Forward forwards reader data to a telegram chat by sending
// messages as a given bot.
func Forward(ctx context.Context, reader io.Reader, token string, chat int, includes, excludes []string) error {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("couldn't create bot api: %w", err)
	}
	data := make([]byte, 1024)
	var n int
	errC := make(chan error)
	for {
		// Read message
		go func() {
			n, err = reader.Read(data)
			errC <- err
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-errC:
		}
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("tgio: couldn't read: %w", err)
		}
		if n == 0 {
			continue
		}

		// Skip message if it doesn't match any include
		var skip bool
		if len(includes) > 0 {
			for _, include := range includes {
				if strings.Contains(string(data[:n]), include) {
					continue
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Skip message if it matches any exclude
		if len(excludes) > 0 {
			for _, exclude := range excludes {
				if !strings.Contains(string(data[:n]), exclude) {
					continue
				}
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Send message
		msg := tgbot.NewMessage(int64(chat), string(data[:n]))
		if _, err := bot.Send(msg); err != nil {
			log.Printf("tgio: %v\n", err)
		}
	}
}
