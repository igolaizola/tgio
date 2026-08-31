package tgio

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TestSendMessageToTopic(t *testing.T) {
	var requests []url.Values
	bot := testBot(func(r *http.Request) (*http.Response, error) {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		requests = append(requests, r.PostForm)
		return telegramResponse(), nil
	})

	if err := sendMessage(bot, -100123, 42, "hello"); err != nil {
		t.Fatalf("sendMessage() error = %v", err)
	}

	if len(requests) != 1 {
		t.Fatalf("got %d requests, want 1", len(requests))
	}
	if got := requests[0].Get("chat_id"); got != "-100123" {
		t.Errorf("chat_id = %q, want %q", got, "-100123")
	}
	if got := requests[0].Get("message_thread_id"); got != "42" {
		t.Errorf("message_thread_id = %q, want %q", got, "42")
	}
	if got := requests[0].Get("text"); got != "hello" {
		t.Errorf("text = %q, want %q", got, "hello")
	}
}

func TestSendMessageWithoutTopic(t *testing.T) {
	var requests []url.Values
	bot := testBot(func(r *http.Request) (*http.Response, error) {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		requests = append(requests, r.PostForm)
		return telegramResponse(), nil
	})

	if err := sendMessage(bot, -100123, 0, "hello"); err != nil {
		t.Fatalf("sendMessage() error = %v", err)
	}

	if len(requests) != 1 {
		t.Fatalf("got %d requests, want 1", len(requests))
	}
	if got := requests[0].Get("message_thread_id"); got != "" {
		t.Errorf("message_thread_id = %q, want it omitted", got)
	}
}

func TestForwardToTopicRejectsNegativeTopic(t *testing.T) {
	err := ForwardToTopic(context.Background(), strings.NewReader("hello"), "token", -100123, -1, nil, nil)
	if err == nil {
		t.Fatal("ForwardToTopic() error = nil, want an error")
	}
	if got, want := err.Error(), "tgio: topic id must be non-negative"; got != want {
		t.Errorf("ForwardToTopic() error = %q, want %q", got, want)
	}
}

func TestShouldSkip(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		includes []string
		excludes []string
		wantSkip bool
	}{
		{name: "no filters", message: "hello"},
		{name: "matches any include", message: "hello world", includes: []string{"missing", "world"}},
		{name: "matches no includes", message: "hello world", includes: []string{"missing", "absent"}, wantSkip: true},
		{name: "matches exclude", message: "hello world", excludes: []string{"world"}, wantSkip: true},
		{name: "include and exclude", message: "hello world", includes: []string{"hello"}, excludes: []string{"world"}, wantSkip: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkip(tt.message, tt.includes, tt.excludes); got != tt.wantSkip {
				t.Errorf("shouldSkip() = %v, want %v", got, tt.wantSkip)
			}
		})
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func testBot(roundTrip roundTripperFunc) *tgbot.BotAPI {
	return &tgbot.BotAPI{
		Token:  "token",
		Client: &http.Client{Transport: roundTrip},
	}
}

func telegramResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"ok":true,"result":{}}`)),
		Header:     make(http.Header),
	}
}
