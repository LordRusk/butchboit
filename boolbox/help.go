// this is where butch keeps his help system.
package boolbox

import (
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/discord"
)

var (
	HelpDivider = "------------\n" // used elsewhere for generation
	HelpColor   = "#fafafa"
)

// singular arguement for command
type Arg struct {
	Name       string `json:"name,omitempty"`
	IsOptional bool   `json:"isoptional,omitempty"`
}

// singular command
// A commands state can be one of three values:
// 0 - Working order
// 1 - Work in progress
// 2 - Does not work
type Cmd struct {
	Cmd   string `json:"cmd,omitempty"`
	Args  []Arg  `json:"args,omitempty"`
	Desc  string `json:"desc,omitempty"`
	State int    `json:"State,omitempty"`
}

// group of commands used to organize
// commands
type CmdGroup struct {
	Name   string `json:"name,omitempty"`
	CmdArr []Cmd  `json:"cmdarr,omitempty"`
}

// CmdGroup wrapper for json
type CmdGroups struct {
	Cm map[string]CmdGroup `json:"name,omitempty"`
}

// generate the help message
func (box *Box) GenHelpMsg(prefix string, botName string, cmdGroupMap map[string]CmdGroup) (*discord.Embed, error) {
	/* generate the help command */
	var helpMsg strings.Builder

	helpMsg.WriteString(HelpDivider)
	helpMsg.WriteString("**Prefix:**  `" + prefix + "`\n" + HelpDivider + "**Commands**\n")
	helpMsg.WriteString(HelpDivider)

	for _, cmdGroup := range cmdGroupMap {
		helpMsg.WriteString("***" + cmdGroup.Name + " Commands:***\n")
		for _, cmdInfo := range cmdGroup.CmdArr {
			if cmdInfo.State == 1 {
				helpMsg.WriteString("__[ Work In Progress ]__ ")
			} else if cmdInfo.State == 2 {
				helpMsg.WriteString("~~")
			}
			helpMsg.WriteString("**" + cmdInfo.Cmd + "**")
			for i := 0; i < len(cmdInfo.Args); i++ {
				helpMsg.WriteString(" [ ")
				if cmdInfo.Args[i].IsOptional == true {
					helpMsg.WriteString("*Optional* ")
				}
				helpMsg.WriteString(cmdInfo.Args[i].Name + " ]")
			}
			helpMsg.WriteString(" -- *" + cmdInfo.Desc + "*")
			if cmdInfo.State == 2 {
				helpMsg.WriteString("~~")
			}
			helpMsg.WriteString("\n")
		}
		helpMsg.WriteString(HelpDivider)
	}

	/* color */
	colorHex, err := strconv.ParseInt((HelpColor)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* make the embed */
	embed := discord.Embed{
		Title:       botName + " Help Page:",
		Description: helpMsg.String(),
		Color:       discord.Color(colorHex),
	}

	return &embed, nil
}
