package main

import (
	"log"
	"os"

	"github.com/diamondburned/arikawa/bot"
	"github.com/lordrusk/butchbot/boolbox"
)

type Bot struct {
	// context must not be embedded
	Ctx *bot.Context
}

var (
	token   = os.Getenv("BOT_TOKEN")
	Prefix  = "!"
	BotName = "ButchBot"

	Box *boolbox.Box
)

func main() {
	if token == "" {
		log.Fatalln("No $BOTTOKEN")
	}

	commands := &Bot{}

	wait, err := bot.Start(token, commands, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix(Prefix)
		ctx.EditableCommands = true

		// get boolbox
		var err error
		Box, err = boolbox.NewBox(commands.Ctx)
		if err != nil {
			log.Fatalln(err)
		}

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
