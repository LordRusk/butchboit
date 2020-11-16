package main

import (
	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
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

	Box.StoreModel(helpPath, help)

	return helpMsg, nil
}

func (b *Bot) Prefix(m *gateway.MessageCreateEvent, newPrefix string) (string, error) {
	Prefix = newPrefix
	b.Ctx.HasPrefix = bot.NewPrefix(Prefix)

	return "`" + Prefix + "` is the new prefix!", nil
}
