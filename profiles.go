package main

import (
	"errors"
	"strings"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	/* Tags */
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
	}

	/* Profiles */
	ProArr = []boolbox.Profile{
		boolbox.Profile{
			Name:     "Butchie",
			Nickname: "Butcheritto",
			Color:    "#ffffff",
			Triggers: []string{"Butch", "butch", "Butchie", "butchie", "Butcheritto", "butcheritto"},
			Tags:     []string{"OGB"},
			Info:     "**General Info:**\n- Butch is all knowing about bool and bool activities\n- He is the goodest boy that ever has lived or will live, he is ***Butchie***",
		},
		boolbox.Profile{
			Name:     "Parker Jacobsen",
			Nickname: "The President",
			Id:       654104550080249866,
			Color:    "#ffffff",
			Triggers: []string{"Parker", "parker", "Parker Jacobsen", "parker jacobsen", "The President"},
			Tags:     []string{"OGB", "GBV", "BDR", "COC"},
			Info:     "**General Info:**\n- Former Bool President before the (not so) great teardown of the bool hierarchy.\n- He started the bool in September 28th, 2019. He is one of only two people to be a bool member consistantly in the bool from the beggining\n- Has been known to drop buckets.\n- In 2015 Parker Jacobson won the aword for smallest penis in the world wide ClitDick Convention™.",
		},
		boolbox.Profile{
			Name:     "Russell Mears",
			Nickname: "Prophet",
			Id:       572291054947139596,
			Color:    "#bb00ff",
			Triggers: []string{"Russell", "russell", "Russ", "russ", "Prophet", "prophet", "McProphet", "mcprophet", "Russell Mears", "russell mears"},
			Tags:     []string{"GBV", "TCH", "COC", "IBDR", "BBC"},
			Info:     "**General Info:**\n- Russell joined the bool at the start of the school year of 2019.\n- Russell was the first of two people to originally drift in the Bool Machine.\n- He was one of three boolers to ever bool to the race track.\n- Russell was heavily addicted to benny along with Sub Bitch for a three month period at the end of school and early into the summer. Russell doesn't even remember a bool because of it.\n- Using Parker's observations, Russell was able to remove the GPS from the Bool Mobile™, greatly increasing what the bool can do and how often the bool can bool.\n- Russell is well known as the bool gopher boy.",
		},
		boolbox.Profile{
			Name:     "Erik Haack",
			Nickname: "TedNut",
			Id:       431696134130499584,
			Color:    "#4F42B5",
			Triggers: []string{"Erik", "erik", "Redneck", "redneck", "Tedneck", "tedneck", "Erik Haack", "erik haack"},
			Tags:     []string{"CKL", "IBDR", "TCH"},
			Info:     "**General Info:**\n- Erik joined the bool midway through the school year of 2019\n- His first bool was on may 29th, 2019\n- Erik started hanging out with the bool squad to simp for sara, but ended up booling, so he joined.\n- Does pop percs and xanny.\n- Good with any motorized vehicle.",
		},
		boolbox.Profile{
			Name:     "Sara Nicholson",
			Nickname: "Shlong Slinging Slasher",
			Id:       650488284991979521,
			Color:    "#e60000",
			Triggers: []string{"Sara", "sara", "Bub_Bitch", "Sub Bitch", "sub bitch", "Sara Nicholson", "sara nicholson", "slong slinging slasher"},
			Tags:     []string{"OGB", "IBDR"},
			Info:     "**General Info:**\n- Sara joined the bool in september 28th, 2019.\n- She is 1 of 4 original boolers.\n- She got the name Bool Mom because she always brings snacks. She is now the offical bool food lady.\n- Sara can, in fact, out smoke the mighty Weezus.\n- Has the ability to make people look.",
		},
		boolbox.Profile{
			Name:     "Micheal Pulver",
			Nickname: "Weezus",
			Id:       332733113769656340,
			Color:    "#285028",
			Triggers: []string{"Micheal", "micheal", "Weezus", "weezus", "Micheal Pulver", "micheal pulver"},
			Tags:     []string{"CAC", "OGB"},
			Info:     "**General Info:**\n- Weezus joined the bool September 28th of 2019\n- He is the only other person than Parker to be apart of the bool from the beginning.\n- If you couldn't tell by the name, Weezus is just that, weed jesus.\n- Weezus has infact, never smoked marijuana in his life!",
		},
		boolbox.Profile{
			Name:     "William Mcamis",
			Nickname: "pet retard",
			Id:       556675821012516886,
			Color:    "#2C2F33",
			Triggers: []string{"William", "william", "Thee Willy Dee", "thee willy dee", "Willy D", "willy d", "Pet Retard", "pet retard", "William Mcamis", "william mcamis"},
			Info:     "**General Info:**\n- William first joined the bool early school year of 2019\n- William is a cellist that burned his hand the day of the last concert performed before the great plague of corona\n- William can toilet teleport. From any toilet, he can instantly teleport to any other toilet in the world. For this, he has earned the role of *toilet teleporter*\n- William is the bool's lawyer.\n- William's only source of entertainment is hentai.\n- William has the ability to turn YOUR pickle dill.",
		},
		boolbox.Profile{
			Name:     "Ashley Nicholson",
			Nickname: "Smol Bool",
			Color:    "#ffffff",
			Id:       568913155451912220,
			Triggers: []string{"Ashley", "ashley", "Smol Bool", "smol bool", "Ashley Nicholson", "ashley nicholson"},
			Info:     "**General Info:**\n- Ashley joined the bool at the start of the great quarantine.\n- She got the name Smol Bool because she has the voice of a loli. (parker likes that)\n- Ashley liked the bool so much, after her first bool she booled the next three days in a row.",
		},
		boolbox.Profile{
			Name:     "Jack Nolan",
			Nickname: "Boomer",
			Color:    "#228B22",
			Triggers: []string{"Jack", "jack", "Jack Nolan", "jack nolan", "Boomer", "boomer"},
			Tags:     []string{"GBV", "MIB"},
			Info:     "**General Info:**\n- Jack joined the bool the bool XXXX-XX-XX.\n- Jack was one of 4 boolers to be in the Great Bool.\n- Jack has PTSD because of the great bool, and has since refused to narrate the story, which is dumb and stupid and gay.\n- Jack heavily simped for *Sub Bitch* during his time in the bool.\n- Jack left the bool because *Prophet* pranked Jack by telling him *Sub Bitch* had a dick.",
		},
		boolbox.Profile{
			Name:     "Jayden Barilo",
			Nickname: "Crit",
			Color:    "#AD03FC",
			Id:       542160535177658401,
			Triggers: []string{"Jayden", "jayden", "Crit", "crit", "CritMass", "critmass"},
			Info:     "**General Info:**\n- Jayden joined the bool September 13, 2020.\n- Jayden was in the bool mobile when the bool got its first ticket.\n- He can out drink anyone in the bool.\n- Parker and Jayden are in a long-term clandestine relationship.",
		},
	}
)

func (botStruct *Bot) Profile(m *gateway.MessageCreateEvent, input bot.RawArguments) (*discord.Embed, error) {
	if string(input) == "" {
		for _, profile := range ProArr {
			if profile.Id == m.Author.ID {
				embed, err := Box.GenProfileEmbed(profile, TagMap)
				if err != nil {
					return nil, err
				}

				return embed, nil
			}
		}
	} else {
		for _, profile := range ProArr {
			for _, trigger := range profile.Triggers {
				if string(input) == trigger {
					embed, err := Box.GenProfileEmbed(profile, TagMap)
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

func (botStruct *Bot) Profiles(*gateway.MessageCreateEvent) (*discord.Embed, error) {
	var profiles strings.Builder

	for _, profile := range ProArr {
		profiles.WriteString(profile.Name)
		profiles.WriteString(" (AKA: ")
		profiles.WriteString(profile.Nickname)
		profiles.WriteString(")\n")
	}

	embed := discord.Embed{
		Title:       "Bool Profiles",
		Description: profiles.String(),
	}

	return &embed, nil
}
