package main

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// path to CmdGroups for help menus
	helpPath = "./help.json"

	// get boolbox.CmdGroups for help message generation
	help = boolbox.CmdGroups{}
	_    = boolbox.GetStoredModel(helpPath, &help)
)

func (b *Bot) Help(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	helpMsg, err := boolbox.GenHelpMsg(*prefix, *botName, help.Cm)
	if err != nil {
		return nil, err
	}

	return helpMsg, nil
}

func (b *Bot) Prefix(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	*prefix = string(input)
	b.Ctx.HasPrefix = bot.NewPrefix(*prefix)

	return fmt.Sprintf("`%s` is the new prefix!", *prefix), nil
}
