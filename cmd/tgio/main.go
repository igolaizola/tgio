package main

import (
	"flag"
	"log"
	"os"

	"github.com/igolaizola/tgio"
)

func main() {
	token := flag.String("token", "", "telegram bot token")
	chat := flag.Int("chat", 0, "telegram chat id")
	flag.Parse()
	if *token == "" {
		log.Fatal("missing token")
	}
	if *chat == 0 {
		log.Fatal("missing chat id")
	}
	if err := tgio.Forward(os.Stdin, *token, *chat); err != nil {
		log.Fatal(err)
	}
}
