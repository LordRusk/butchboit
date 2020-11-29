package main

import (
	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// path to CmdGroups for help menus
	helpPath = "./help.json"

	// get boolbox.CmdGroups for help messase generation
	help = boolbox.CmdGroups{}
	_    = Box.GetStoredModel(helpPath, &help)
)

func (b *Bot) Help(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	helpMsg, err := Box.GenHelpMsg(Prefix, BotName, help.Cm)
	if err != nil {
		return nil, err
	}

	return helpMsg, nil
}

func (b *Bot) Prefix(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	Prefix = string(input)
	b.Ctx.HasPrefix = bot.NewPrefix(Prefix)

	return "`" + Prefix + "` is the new prefix!", nil
}
