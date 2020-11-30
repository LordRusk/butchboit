package main

import (
	"errors"
	"os"
	"strings"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// path to profile json file
	proPath = os.Getenv("HOME") + "/.local/share/profiles.json"

	// get profiles
	Profiles = boolbox.Profiles{}
	_        = boolbox.GetStoredModel(proPath, &Profiles)
)

var (
	// Tags
	ButchBotCreator = discord.EmbedField{
		Name:   "Butch Bot Creator",
		Value:  "Created Butch Bot",
		Inline: true,
	}

	OgBooler = discord.EmbedField{
		Name:   "OG Booler",
		Value:  "Booled in the first bool",
		Inline: true,
	}

	GreatBoolVeteran = discord.EmbedField{
		Name:   "Great Bool Veteran",
		Value:  "Fought in the great bool",
		Inline: true,
	}

	BoolDriver = discord.EmbedField{
		Name:   "Bool Driver",
		Value:  "Can legally drive the bool",
		Inline: true,
	}

	IllegalBoolDriver = discord.EmbedField{
		Name:   "Illegal Bool Driver",
		Value:  "Has driven the bool illegally",
		Inline: true,
	}

	BoolTech = discord.EmbedField{
		Name:   "Bool Tech",
		Value:  "Good with a form of technology",
		Inline: true,
	}

	MissingInBool = discord.EmbedField{
		Name:   "Missing In Bool",
		Value:  "No longer apart of the bool",
		Inline: true,
	}

	ConeKiller = discord.EmbedField{
		Name:   "Cone Killer",
		Value:  "XXXX-XXXX-XXXXX-XXX",
		Inline: true,
	}

	ConeAcomplice = discord.EmbedField{
		Name:   "Cone Accomplice",
		Value:  "XXXX-XXXX-XXXXX-XXX",
		Inline: true,
	}

	ConeOnConscious = discord.EmbedField{
		Name:   "Cone On The Conscious",
		Value:  "XXXX-XXXX-XXXXX-XXX",
		Inline: true,
	}

	DeerSummoner = discord.EmbedField{
		Name:   "Will Summon Deers",
		Value:  "Your car will not last ;)",
		Inline: true,
	}

	TagMap = map[string]discord.EmbedField{
		"BBC":  ButchBotCreator,
		"OGB":  OgBooler,
		"GBV":  GreatBoolVeteran,
		"BDR":  BoolDriver,
		"IBDR": IllegalBoolDriver,
		"TCH":  BoolTech,
		"MIB":  MissingInBool,
		"CKL":  ConeKiller,
		"CAC":  ConeAcomplice,
		"COC":  ConeOnConscious,
		"DSR":  DeerSummoner,
	}
)

func (b *Bot) Profile(m *gateway.MessageCreateEvent, input bot.RawArguments) (*discord.Embed, error) {
	if string(input) == "" {
		for _, profile := range Profiles.Ps {
			if profile.ID == m.Author.ID {
				embed, err := boolbox.ProfileToEmbed(profile, TagMap)
				if err != nil {
					return nil, err
				}

				return embed, nil
			}
		}
	} else {
		for _, profile := range Profiles.Ps {
			for _, trigger := range profile.Triggers {
				if string(input) == trigger {
					embed, err := boolbox.ProfileToEmbed(profile, TagMap)
					if err != nil {
						return nil, err
					}

					return embed, nil
				}
			}
		}
	}

	return nil, errors.New("Error! Profile not found. Please run `!profiles` for a list")
}

func (b *Bot) Profiles(*gateway.MessageCreateEvent) (*discord.Embed, error) {
	var desc strings.Builder

	for _, profile := range Profiles.Ps {
		desc.WriteString(profile.Name)
		desc.WriteString(" (AKA: ")
		desc.WriteString(profile.Nickname)
		desc.WriteString(")\n")
	}

	embed := discord.Embed{
		Title:       "Bool Profiles",
		Description: desc.String(),
	}

	return &embed, nil
}
