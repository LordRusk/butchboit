package main

import (
	"log"
	"os"

	"github.com/diamondburned/arikawa/bot"
)

type Bot struct {
	/* Context must not be embedded */
	Ctx *bot.Context
}

var (
	token   = os.Getenv("BOT_TOKEN")
	Prefix  = "!"
	BotName = "ButchBot"
)

func main() {
	if token == "" {
		log.Fatalln("No $BOT_TOKEN")
	}

	commands := &Bot{}

	wait, err := bot.Start(token, commands, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix(Prefix)
		ctx.EditableCommands = true

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(BotName, "has started")

	if err := wait(); err != nil {
		log.Fatalln("Gateway fetal error:", err)
	}
}
