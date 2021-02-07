package boolbox

import (
	"strconv"

	"github.com/diamondburned/arikawa/v2/discord"
)

// holds a profile
type Profile struct {
	Name     string         `json:"name,omitempty"`
	Nickname string         `json:"nickname,omitempty"`
	ID       discord.UserID `json:"id,string"`
	Color    string         `json:"color,omitempty"`
	Triggers []string       `json:"triggers,omitempty"`
	Info     string         `json:"info,omitempty"`

	// array of keys which will be tried against
	// a map[string]discord.Embed to find tags.
	Tags []string `json:"tags,omitempty"`
}

// Profile wrapper for json
type Profiles struct {
	Ps []Profile `json:"ps,omitempty"`
}

func ProfileToEmbed(profile Profile, tagMap map[string]discord.EmbedField) (*discord.Embed, error) {
	colorHex, err := strconv.ParseInt((profile.Color)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	// tags
	fields := []discord.EmbedField{}
	for _, tag := range profile.Tags {
		if _, ok := tagMap[tag]; ok {
			fields = append(fields, tagMap[tag])
		}
	}

	embed := discord.Embed{
		Title:       "Bool profile for " + profile.Name + " (AKA: " + profile.Nickname + ")",
		Description: profile.Info,
		Color:       discord.Color(colorHex),
		Fields:      fields,
	}

	return &embed, nil
}
