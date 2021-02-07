package main

import (
	"flag"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/lordrusk/butchbot/boolbox"
)

type Bot struct {
	// context must not be embeded
	Ctx *bot.Context
}

var token = flag.String("t", "", "Set the bot token. Default uses 'BOT_TOKEN' env.")
var prefix = flag.String("p", "!", "Set the prefix.")
var botName = flag.String("n", "ButchBot", "Set the bot's name.")

var logger *log.Logger
var box *boolbox.Box // client for boolbox framework

func main() {
	flag.Parse()
	logger = log.New(os.Stdout, *botName+": ", 0)

	if *token == "" {
		toke := os.Getenv("BOT_TOKEN")
		if toke == "" {
			logger.Fatalln("No BOT_TOKEN: Set BOT_TOKEN or use '-t'")
		}

		token = &toke
	}

	commands := &Bot{}

	wait, err := bot.Start(*token, commands, func(ctx *bot.Context) error {
		ctx.HasPrefix = bot.NewPrefix(*prefix)
		ctx.EditableCommands = true // <- this is nice

		// get box
		var err error
		box, err = boolbox.NewBox(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Fatalf("Could not start bot: %s\n", err)
	}

	logger.Println("Bot Started")

	if err := wait(); err != nil {
		logger.Fatalf("Gateway fetal error: %s\n", err)
	}
}

const AUTHOR = "Prophet#5193"
