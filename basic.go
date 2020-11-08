package main

import (
	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// A commands state can be one of three values:
	// 0 - Working order
	// 1 - Work in progress
	// 2 - Does not work

	/* command groups */
	basicCmdGroup = boolbox.CmdGroup{
		Name: "Basic",
		CmdArr: []boolbox.Cmd{
			boolbox.Cmd{
				Cmd:  "help",
				Desc: "Show's the help page",
			},
			boolbox.Cmd{
				Cmd:  "prefix",
				Args: []boolbox.Arg{boolbox.Arg{Name: "New Prefix", IsOptional: false}},
				Desc: "Set a new prefix",
			},
		},
	}

	appointCmdGroup = boolbox.CmdGroup{
		Name: "Appoint",
		CmdArr: []boolbox.Cmd{
			boolbox.Cmd{
				Cmd:  "newbool",
				Args: []boolbox.Arg{boolbox.Arg{Name: "Script", IsOptional: false}},
				Desc: "Create a new bool event",
			},
			boolbox.Cmd{
				Cmd:  "removebool",
				Args: []boolbox.Arg{boolbox.Arg{Name: "Script", IsOptional: false}},
				Desc: "Remove a bool event",
			},
			boolbox.Cmd{
				Cmd:  "editbool",
				Args: []boolbox.Arg{boolbox.Arg{Name: "Script", IsOptional: false}},
				Desc: "Edit a bool",
			},
			boolbox.Cmd{
				Cmd:  "bool",
				Args: []boolbox.Arg{boolbox.Arg{Name: "Script", IsOptional: false}},
				Desc: "Show info for a bool event.",
			},
			boolbox.Cmd{
				Cmd:  "rsvp",
				Args: []boolbox.Arg{boolbox.Arg{Name: "Script", IsOptional: false}},
				Desc: "rsvp for a bool event.",
			},
			boolbox.Cmd{
				Cmd:  "bools",
				Desc: "Shows currently active bool events.",
			},
		},
	}

	profileSystemCmdGroup = boolbox.CmdGroup{
		Name: "Profiles",
		CmdArr: []boolbox.Cmd{
			boolbox.Cmd{
				Cmd:  "profile",
				Args: []boolbox.Arg{boolbox.Arg{Name: "User Profile", IsOptional: true}},
				Desc: "Show a profile from the list of `!profiles`",
			},
			boolbox.Cmd{
				Cmd:  "profiles",
				Desc: "Get a list of profiles",
			},
		},
	}

	fourCmdGroup = boolbox.CmdGroup{
		Name: "Four",
		CmdArr: []boolbox.Cmd{
			boolbox.Cmd{
				Cmd:  "post",
				Desc: "Get a random 4chan post",
			},
			boolbox.Cmd{
				Cmd:  "board",
				Args: []boolbox.Arg{boolbox.Arg{Name: "A 4chan board", IsOptional: false}},
				Desc: "Get a random 4chan post from a specified board.",
			},
			boolbox.Cmd{
				Cmd:   "scope",
				Args:  []boolbox.Arg{boolbox.Arg{Name: "4chan Board", IsOptional: false}, boolbox.Arg{Name: "4chan Post No.", IsOptional: false}},
				Desc:  "Scope out a certain 4chan post.",
				State: 2,
			},
			boolbox.Cmd{
				Cmd:  "boards",
				Desc: "Get a list of 4chan boards",
			},
		},
	}

	cmdGroupMap = map[string]boolbox.CmdGroup{
		"basic":    basicCmdGroup,
		"appoint":  appointCmdGroup,
		"profiles": profileSystemCmdGroup,
		"four":     fourCmdGroup,
	}
)

func (botStruct *Bot) Help(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	helpMsg, err := Box.GenHelpMsg(Prefix, BotName, cmdGroupMap)
	if err != nil {
		return nil, err
	}

	return helpMsg, nil
}

func (botStruct *Bot) Prefix(m *gateway.MessageCreateEvent, newPrefix string) (string, error) {
	Prefix = newPrefix
	botStruct.Ctx.HasPrefix = bot.NewPrefix(Prefix)

	return "`" + Prefix + "` is the new prefix!", nil
}
