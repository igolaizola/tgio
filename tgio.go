package tgio

import (
	"errors"
	"fmt"
	"io"
	"log"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Forward forwards reader data to a telegram chat by sending
// messages as a given bot.
func Forward(reader io.Reader, token string, chat int) error {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("couldn't create bot api: %w", err)
	}
	data := make([]byte, 1024)
	for {
		n, err := reader.Read(data)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("tgio: couldn't read: %w", err)
		}
		if n == 0 {
			continue
		}
		msg := tgbot.NewMessage(int64(chat), string(data[:n]))
		if _, err := bot.Send(msg); err != nil {
			log.Println(fmt.Sprintf("tgio: %v", err))
		}
	}
}
