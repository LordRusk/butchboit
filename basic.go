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
	helpCmdInfo = cmdInfo{
		cmd:   "help",
		desc:  "Show's the help page",
		state: 1,
	}

	prefixCmdInfo = cmdInfo{
		cmd:  "prefix",
		args: []arg{arg{name: "New Prefix", isOptional: false}},
		desc: "Set a new prefix",
	}

	profileCmdInfo = cmdInfo{
		cmd:  "profile",
		args: []arg{arg{name: "User Profile", isOptional: true}},
		desc: "Show a profile from the list of `!profiles`",
	}

	profilesCmdInfo = cmdInfo{
		cmd:  "profiles",
		desc: "Get a list of profiles",
	}

	postCmdInfo = cmdInfo{
		cmd:  "post",
		desc: "Get a random 4chan post",
	}

	boardCmdInfo = cmdInfo{
		cmd:  "board",
		args: []arg{arg{name: "A 4chan board", isOptional: false}},
		desc: "Get a random 4chan post from a specified board.",
	}

	scopeCmdInfo = cmdInfo{
		cmd:   "scope",
		args:  []arg{arg{name: "4chan Board", isOptional: false}, arg{name: "4chan Post No.", isOptional: false}},
		desc:  "Scope out a certain 4chan post.",
		state: 2,
	}

	boardsCmdInfo = cmdInfo{
		cmd:  "boards",
		desc: "Get a list of 4chan boards",
	}

	/* command groups */
	basicCmdGroup = cmdGroup{
		name:   "Basic",
		cmdArr: []cmdInfo{helpCmdInfo, prefixCmdInfo},
	}

	profilingSystemCmdGroup = cmdGroup{
		name:   "Profiling System",
		cmdArr: []cmdInfo{profileCmdInfo, profilesCmdInfo},
	}

	chanCmdGroup = cmdGroup{
		name:   "4chan",
		cmdArr: []cmdInfo{postCmdInfo, boardCmdInfo, scopeCmdInfo, boardsCmdInfo},
	}

	cmdGroupMap = map[string]cmdGroup{
		"basic":     basicCmdGroup,
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
