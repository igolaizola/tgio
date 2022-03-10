package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/igolaizola/tgio"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	// Create signal based context
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
			cancel()
		}
		signal.Stop(c)
	}()

	// Launch command
	cmd := newForwardCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newForwardCommand() *ffcli.Command {
	fs := flag.NewFlagSet("tgio", flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	token := fs.String("token", "", "telegram bot token")
	chat := fs.Int("chat", 0, "telegram chat id")

	return &ffcli.Command{
		Name:       "tgio",
		ShortUsage: "tgio [flags] <key> <value data...>",
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithEnvVarPrefix("TGIO"),
		},
		ShortHelp: "run tgio forwarder",
		FlagSet:   fs,
		Exec: func(ctx context.Context, args []string) error {
			if *token == "" {
				return errors.New("missing token")
			}
			if *chat == 0 {
				return errors.New("missing chat id")
			}
			return tgio.Forward(ctx, os.Stdin, *token, *chat)
		},
	}
}
