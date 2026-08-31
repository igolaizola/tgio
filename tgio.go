package tgio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Forward forwards reader data to a telegram chat by sending
// messages as a given bot.
func Forward(ctx context.Context, reader io.Reader, token string, chat int, includes, excludes []string) error {
	return forward(ctx, reader, token, chat, 0, includes, excludes)
}

// ForwardToTopic forwards reader data to a Telegram group topic by sending
// messages as a given bot. A topic ID of zero sends to the chat without a
// message thread, preserving the behavior of Forward.
func ForwardToTopic(ctx context.Context, reader io.Reader, token string, chat, topic int, includes, excludes []string) error {
	if topic < 0 {
		return errors.New("tgio: topic id must be non-negative")
	}

	return forward(ctx, reader, token, chat, topic, includes, excludes)
}

func forward(ctx context.Context, reader io.Reader, token string, chat, topic int, includes, excludes []string) error {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("couldn't create bot api: %w", err)
	}
	data := make([]byte, 1024)
	type readResult struct {
		n   int
		err error
	}
	errC := make(chan readResult, 1)
	for {
		// Read message
		go func() {
			n, err := reader.Read(data)
			errC <- readResult{n: n, err: err}
		}()
		var result readResult
		select {
		case <-ctx.Done():
			return ctx.Err()
		case result = <-errC:
		}

		n, err := result.n, result.err
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("tgio: couldn't read: %w", err)
		}
		if n == 0 {
			continue
		}

		message := string(data[:n])
		if shouldSkip(message, includes, excludes) {
			continue
		}

		// Send message
		if err := sendMessage(bot, chat, topic, message); err != nil {
			log.Printf("tgio: %v\n", err)
		}
	}
}

func shouldSkip(message string, includes, excludes []string) bool {
	if len(includes) > 0 {
		matched := false
		for _, include := range includes {
			if strings.Contains(message, include) {
				matched = true
				break
			}
		}
		if !matched {
			return true
		}
	}

	for _, exclude := range excludes {
		if strings.Contains(message, exclude) {
			return true
		}
	}

	return false
}

func sendMessage(bot *tgbot.BotAPI, chat, topic int, text string) error {
	if topic == 0 {
		msg := tgbot.NewMessage(int64(chat), text)
		_, err := bot.Send(msg)
		return err
	}

	params := url.Values{}
	params.Set("chat_id", strconv.FormatInt(int64(chat), 10))
	params.Set("message_thread_id", strconv.Itoa(topic))
	params.Set("text", text)
	_, err := bot.MakeRequest("sendMessage", params)
	return err
}
