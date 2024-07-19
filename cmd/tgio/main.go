package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"

	"github.com/igolaizola/tgio"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// Build flags
var version = ""
var commit = ""
var date = ""

func main() {
	// Create signal based context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Launch command
	cmd := newCommand()
	if err := cmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func newCommand() *ffcli.Command {
	fs := flag.NewFlagSet("tgio", flag.ExitOnError)
	_ = fs.String("config", "", "config file (optional)")

	token := fs.String("token", "", "telegram bot token")
	chat := fs.Int("chat", 0, "telegram chat id")
	var includes, excludes []string
	fs.Var(fsStrings(&includes), "include", "include only messages that match this")
	fs.Var(fsStrings(&excludes), "exclude", "exclude messages that match this")

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
			return tgio.Forward(ctx, os.Stdin, *token, *chat, includes, excludes)
		},
		Subcommands: []*ffcli.Command{
			newVersionCommand(),
		},
	}
}

func newVersionCommand() *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "goobar version",
		ShortHelp:  "print version",
		Exec: func(ctx context.Context, args []string) error {
			v := version
			if v == "" {
				if buildInfo, ok := debug.ReadBuildInfo(); ok {
					v = buildInfo.Main.Version
				}
			}
			if v == "" {
				v = "dev"
			}
			versionFields := []string{v}
			if commit != "" {
				versionFields = append(versionFields, commit)
			}
			if date != "" {
				versionFields = append(versionFields, date)
			}
			fmt.Println(strings.Join(versionFields, " "))
			return nil
		},
	}
}

type stringsValue []string

func (f *stringsValue) String() string {
	return strings.Join(*f, ", ")
}

func (f *stringsValue) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func fsStrings(p *[]string) *stringsValue {
	return (*stringsValue)(p)
}
