// this is where butch keeps his profiling system.
package boolbox

import (
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

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

// Turns a Profile into a discord.Embed.
func (box *Box) ProfileToEmbed(profile Profile, tagMap map[string]discord.EmbedField) (*discord.Embed, error) {
	/* title */
	var title strings.Builder
	title.WriteString("Bool profile for ")
	title.WriteString(profile.Name)
	title.WriteString(" (AKA: ")
	title.WriteString(profile.Nickname)
	title.WriteString(")")

	/* color */
	colorHex, err := strconv.ParseInt((profile.Color)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* tags */
	fields := []discord.EmbedField{}

	for _, tag := range profile.Tags {
		if _, ok := tagMap[tag]; ok {
			fields = append(fields, tagMap[tag])
		}
	}

	/* make the embed */
	embed := discord.Embed{
		Title:       title.String(),
		Description: profile.Info,
		Color:       discord.Color(colorHex),
		Fields:      fields,
	}

	return &embed, nil
}
