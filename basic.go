package main

import (
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

type arg struct {
	name       string
	isOptional bool
}

type cmdInfo struct {
	cmd   string
	args  []arg
	desc  string
	state int
}

type cmdGroup struct {
	name   string
	cmdArr []cmdInfo
}

var (
	/* Help message stuff */
	HelpDivider = "------------\n"
	helpColor   = "#fafafa"

	/* A commands state can be one of three values:
	 * 0 - Working order
	 * 1 - Work in progress
	 * 2 - Does not work
	 */

	/* command groups */
	basicCmdGroup = cmdGroup{
		name: "Basic",
		cmdArr: []cmdInfo{
			cmdInfo{
				cmd:  "help",
				desc: "Show's the help page",
			},
			cmdInfo{
				cmd:  "prefix",
				args: []arg{arg{name: "New Prefix", isOptional: false}},
				desc: "Set a new prefix",
			},
		},
	}

	rsvpCmdGroup = cmdGroup{
		name: "RSVP",
		cmdArr: []cmdInfo{
			cmdInfo{
				cmd:  "newBool",
				args: []arg{arg{name: "Bool Name", isOptional: false}, arg{name: "Bool Time", isOptional: false}, arg{name: "Bool Date", isOptional: false}, arg{name: "Bool Description (in quotes)", isOptional: false}},
				desc: "Create a new bool event",
			},
			cmdInfo{
				cmd:  "removeBool",
				args: []arg{arg{name: "Selected Bool", isOptional: false}},
				desc: "Remove a bool event",
			},
			cmdInfo{
				cmd:  "boolInfo",
				args: []arg{arg{name: "Selected Bool", isOptional: false}},
				desc: "Show info for a bool event.",
			},
			cmdInfo{
				cmd:  "rsvp",
				args: []arg{arg{name: "Selected Bool", isOptional: false}},
				desc: "rsvp for a bool event.",
			},
			cmdInfo{
				cmd:  "bools",
				desc: "Shows currently active bool events.",
			},
		},
	}

	profilingSystemCmdGroup = cmdGroup{
		name: "Profiling System",
		cmdArr: []cmdInfo{
			cmdInfo{
				cmd:  "profile",
				args: []arg{arg{name: "User Profile", isOptional: true}},
				desc: "Show a profile from the list of `!profiles`",
			},
			cmdInfo{
				cmd:  "profiles",
				desc: "Get a list of profiles",
			},
		},
	}

	chanCmdGroup = cmdGroup{
		name: "4chan",
		cmdArr: []cmdInfo{
			cmdInfo{
				cmd:  "post",
				desc: "Get a random 4chan post",
			},
			cmdInfo{
				cmd:  "board",
				args: []arg{arg{name: "A 4chan board", isOptional: false}},
				desc: "Get a random 4chan post from a specified board.",
			},
			cmdInfo{
				cmd:   "scope",
				args:  []arg{arg{name: "4chan Board", isOptional: false}, arg{name: "4chan Post No.", isOptional: false}},
				desc:  "Scope out a certain 4chan post.",
				state: 2,
			},
			cmdInfo{
				cmd:  "boards",
				desc: "Get a list of 4chan boards",
			},
		},
	}

	cmdGroupMap = map[string]cmdGroup{
		"basic":     basicCmdGroup,
		"rsvp":      rsvpCmdGroup,
		"profiling": profilingSystemCmdGroup,
		"4chan":     chanCmdGroup,
	}
)

/* Bot commands */
func (botStruct *Bot) Prefix(m *gateway.MessageCreateEvent, newPrefix string) (string, error) {
	Prefix = newPrefix
	botStruct.Ctx.HasPrefix = bot.NewPrefix(Prefix)

	return "`" + Prefix + "` is the new prefix!", nil
}

func (botStruct *Bot) Help(*gateway.MessageCreateEvent) (*discord.Embed, error) {
	/* generate the help command */
	var helpMsg strings.Builder

	helpMsg.WriteString(HelpDivider)
	helpMsg.WriteString("**Prefix:**  `")
	helpMsg.WriteString(Prefix)
	helpMsg.WriteString("`\n")
	helpMsg.WriteString(HelpDivider)
	helpMsg.WriteString("**Commands**\n")
	helpMsg.WriteString(HelpDivider)

	for _, cmdGroup := range cmdGroupMap {
		helpMsg.WriteString("***")
		helpMsg.WriteString(cmdGroup.name)
		helpMsg.WriteString(" Commands:***\n")
		for _, cmdInfo := range cmdGroup.cmdArr {
			if cmdInfo.state == 1 {
				helpMsg.WriteString("__[ Work In Progress ]__ ")
			} else if cmdInfo.state == 2 {
				helpMsg.WriteString("~~")
			}
			helpMsg.WriteString("**")
			helpMsg.WriteString(cmdInfo.cmd)
			helpMsg.WriteString("**")
			for i := 0; i < len(cmdInfo.args); i++ {
				helpMsg.WriteString(" [ ")
				if cmdInfo.args[i].isOptional == true {
					helpMsg.WriteString("*Optional* ")
				}
				helpMsg.WriteString(cmdInfo.args[i].name)
				helpMsg.WriteString(" ]")
			}
			helpMsg.WriteString(" -- *")
			helpMsg.WriteString(cmdInfo.desc)
			helpMsg.WriteString("*")
			if cmdInfo.state == 2 {
				helpMsg.WriteString("~~")
			}
			helpMsg.WriteString("\n")
		}
		helpMsg.WriteString(HelpDivider)
	}

	/* color */
	colorHex, err := strconv.ParseInt((helpColor)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* make the embed */
	embed := discord.Embed{
		Title:       BotName + " Help Page:",
		Description: helpMsg.String(),
		Color:       discord.Color(colorHex),
	}

	return &embed, nil
}
