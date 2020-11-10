package boolbox

import (
	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
)

// A boolbox struct
type Box struct {
	// context must not be embeded
	Ctx *bot.Context
}

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

// holds all information about a
// bool profile
type Profile struct {
	Name     string         `json:"name,omitempty"`
	Nickname string         `json:"nickname,omitempty"`
	ID       discord.UserID `json:"id,string"`
	Color    string         `json:"color,omitempty"`
	Triggers []string       `json:"triggers,omitempty"`
	Tags     []string       `json:"tags,omitempty"`
	Info     string         `json:"info,omitempty"`
}

// Profile wrapper for json
type Profiles struct {
	Ps []Profile `json:"ps,omitempty"`
}

// rsvp struct used to keep track
// of pickup times and discord.user's
type Rsvp struct {
	User   discord.User `json:"user,omitempty"`
	PuTime string       `json:"puTime,omitempty"`
}

// appointment struct
type Appointment struct {
	Name string `json:"name,omitempty"`
	Date string `json:"date,omitempty"`
	Time string `json:"time,omitempty"`
	Desc string `json:"desc,omitempty"`
	Resv []Rsvp `json:"resv,omitempty"`
}

// appointment wrapper for json
type Appointments struct {
	Appts []Appointment `json:"appts,omitempty"`
}
