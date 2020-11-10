package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// apptsPath to where bools (appointments) are stored
	apptsPath = os.Getenv("HOME") + "/.local/share/bools.json"

	// get bools (appointments)
	Bools = boolbox.Appointments{}
	_     = Box.GetStoredModel(apptsPath, &Bools)

	// inquires
	dateInq    = "Date?"
	timeInq    = "Start of pickup time?"
	descInq    = "Description?"
	rsvpInq    = "What is your estimated pickup time?"
	dateInqDef = dateInq
	timeInqDef = timeInq
	descInqDef = descInq
	rsvpInqDef = rsvpInq

	// used for Editbool() menus
	rsvpNumOpts = "```\n[0] Edit Time\n[1] Delete rsvp\n[2] Exit```"
)

// great demonstrastion of the
// Ask() function.
func (botStruct *Bot) Newbool(m *gateway.MessageCreateEvent) (string, error) {
	var pass bool
	appointment := boolbox.Appointment{}

	uc := make(chan int)
	go Box.Track2Delete(m, uc)

	// get the name of the appointment
	resp, err := Box.Ask(m, "Name?")
	if err != nil {
		uc <- 0
		return "", err
	}

	if resp != "" {
		appointment.Name = resp
	}

	// get the date of the appointment
	for pass == false {
		resp, err := Box.Ask(m, dateInq)
		if err != nil {
			uc <- 0
			return "", err
		}

		if err := Box.CheckDate(resp); err == nil {
			appointment.Date = resp
			pass = true
		}

		dateInq = "Invalid date! Try 7/11, 23/12/2020, etc..."
	}

	// get the time of the appointment
	pass = false
	for pass == false {
		resp, err := Box.Ask(m, timeInq)
		if err != nil {
			uc <- 0
			return "", err
		}

		if err := Box.CheckTime(resp); err == nil {
			appointment.Time = resp
			pass = true
		}

		timeInq = "Invalid date! Try 7:30, 20:45, etc..."
	}

	// get the description of the appointment
	resp, err = Box.Ask(m, descInq)
	if err != nil {
		uc <- 0
		return "", err
	}

	if resp != "" {
		appointment.Desc = resp
	}

	// set inquires back to default
	dateInq = dateInqDef
	timeInq = timeInqDef
	descInq = descInqDef

	Bools.Appts = append(Bools.Appts, appointment)
	Box.StoreModel(apptsPath, Bools)
	uc <- 0

	return "New bool added! Check for a current list of bools with `" + Prefix + "bools`!", nil
}

func (botStruct *Bot) Removebool(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	var bwoolNum int
	var pass bool

	uc := make(chan int)
	go Box.Track2Delete(m, uc)

	for pass == false {
		resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like?\n```\n", "```", Bools.Appts))
		if err != nil {
			uc <- 0
			return "", err
		}

		intResp, _ := strconv.Atoi(resp)
		for num, _ := range Bools.Appts {
			if num == intResp {
				bwoolNum = num
				pass = true
			}
		}

		if pass == false {
			_, err := Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}
		}

	}

	resp, err := Box.Ask(m, "Do you really want to remove that bool? [y/N]")
	if err != nil {
		uc <- 0
		return "", nil
	}

	if resp == "y" || resp == "Y" {
		Bools.Appts = Box.RemoveAppointment(Bools.Appts, bwoolNum)
		if err := Box.StoreModel(apptsPath, Bools); err != nil {
			log.Fatalln(err)
		}

		uc <- 0
		return "Successfully removed bool!", nil
	}

	return "Bool not removed.", nil
}

// another demonstration on
// the usefullness of Ask()
func (botStruct *Bot) Rsvp(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	var bwoolNum int
	var pass bool

	uc := make(chan int)
	go Box.Track2Delete(m, uc)

	for pass == false {
		resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like to rsvp for?\n```\n", "```", Bools.Appts))
		if err != nil {
			uc <- 0
			return "", err
		}

		intResp, _ := strconv.Atoi(resp)
		for num, _ := range Bools.Appts {
			if num == intResp {
				bwoolNum = num
				pass = true
			}
		}

		if pass == false {
			_, err := Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}
		}
	}

	for Num, rsvp := range Bools.Appts[bwoolNum].Resv {
		if rsvp.User == m.Author {
			Bools.Appts[bwoolNum].Resv = Box.RemoveRsvp(Bools.Appts[bwoolNum].Resv, Num)
			uc <- 0

			return "Successfully un-RSVP'd!", nil
		}
	}

	rsvp := boolbox.Rsvp{User: m.Author}

	pass = false
	for pass == false {
		resp, err := Box.Ask(m, rsvpInq)
		if err != nil {
			uc <- 0
			return "", err
		}

		if err := Box.CheckTime(resp); err == nil {
			rsvp.PuTime = resp
			pass = true

		}

		rsvpInq = "Invalid time! Try 7:30, 20:45, etc..."
	}

	rsvpInq = rsvpInqDef

	Bools.Appts[bwoolNum].Resv = append(Bools.Appts[bwoolNum].Resv, rsvp)
	Box.StoreModel(apptsPath, Bools)
	uc <- 0

	return "Successfully RSVP'd!", nil
}

func (botStruct *Bot) Editbool(m *gateway.MessageCreateEvent) (string, error) {
	var bwoolNum int
	var rsvpNum int
	var sectNum int
	var rsvpMenuNum int
	var pass bool
	var builder strings.Builder

	uc := make(chan int)
	go Box.Track2Delete(m, uc)

	for pass == false {
		resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like to edit?\n```\n", "```", Bools.Appts))
		if err != nil {
			uc <- 0
			return "", err
		}

		intResp, _ := strconv.Atoi(resp)
		for num, _ := range Bools.Appts {
			if num == intResp {
				bwoolNum = num
				pass = true
			}
		}

		if pass == false {
			_, err := Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}
		}
	}

	boolFields := Box.GetApptSects()

	builder.Write([]byte("```\n"))
	for i := 0; i < len(boolFields); i++ {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(i)))
		builder.Write([]byte("] "))
		builder.Write([]byte(boolFields[i]))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	pass = false
	for pass == false {
		resp, err := Box.Ask(m, "Which part of the bool would you like to edit?\n"+builder.String())
		if err != nil {
			uc <- 0
			return "", nil
		}

		sectNum, _ = strconv.Atoi(resp)
		for key, _ := range boolFields {
			if key == sectNum {
				pass = true
			}
		}

		if pass == false {
			_, err := Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}
		}
	}

	if sectNum == 0 {
		resp, err := Box.Ask(m, "What would you like to change the name to?")
		if err != nil {
			uc <- 0
			return "", err
		}

		Bools.Appts[bwoolNum].Name = resp
		uc <- 0

		return "Successfully changed bool name!", nil
	} else if sectNum == 1 {
		resp, err := Box.Ask(m, "What would you like to change the date to?")
		if err != nil {
			uc <- 0
			return "", err
		}

		if err := Box.CheckDate(resp); err != nil {
			return "", err
		}

		Bools.Appts[bwoolNum].Date = resp
		uc <- 0

		return "Successfully changed bool date!", nil
	} else if sectNum == 2 {
		resp, err := Box.Ask(m, "What would you like to change the time to?")
		if err != nil {
			uc <- 0
			return "", err
		}

		if err := Box.CheckTime(resp); err != nil {
			return "", err
		}

		Bools.Appts[bwoolNum].Time = resp
		uc <- 0

		return "Successfully changed bool time!", nil
	} else if sectNum == 3 {
		resp, err := Box.Ask(m, "What would you like to change the description to?")
		if err != nil {
			uc <- 0
			return "", err
		}

		Bools.Appts[bwoolNum].Desc = resp
		uc <- 0

		return "Successfully changed bool description!", nil
	}

	if len(Bools.Appts[bwoolNum].Resv) < 1 {
		uc <- 0
		return "", errors.New("Nobody has rsvp'd, so there are no rsvp's to edit.")
	}

	pass = false
	for pass == false {
		resp, err := Box.Ask(m, Box.NumRsvpList("Which rsvp would you like to change??\n```\n", "```", Bools.Appts[bwoolNum].Resv))
		if err != nil {
			uc <- 0
			return "", err
		}

		intResp, _ := strconv.Atoi(resp)
		for num, _ := range Bools.Appts[bwoolNum].Resv {
			if num == intResp {
				rsvpNum = num
				pass = true
			}
		}

		if pass == false {
			_, err := Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}

		}

		pass = false
		for pass == false {
			resp, err := Box.Ask(m, "What would you like to do to the rsvp?\n"+rsvpNumOpts)
			if err != nil {
				uc <- 0
				return "", err
			}

			intResp, err := strconv.Atoi(resp)
			if err == nil {
				for i := 0; i < 3; i++ {
					if i == intResp {
						rsvpMenuNum = i
						pass = true
					}
				}
			}

			if pass == false {
				_, err := Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
				if err != nil {
					return "", err
				}
			}
		}

		var passed bool
		if rsvpMenuNum == 0 {
			for passed == false {
				resp, err := Box.Ask(m, "What would you like the new pickup time be?")
				if err != nil {
					uc <- 0
					return "", err
				}

				if err := Box.CheckTime(resp); err != nil {
					_, err := Box.Ctx.SendMessage(m.ChannelID, "Invalid date! Try 7:30, 20:45, etc...", nil)
					if err != nil {
						return "", err
					}
				} else {
					passed = true
				}

				if passed == true {
					Bools.Appts[bwoolNum].Resv[rsvpNum].PuTime = resp
					uc <- 0

					return "Successfully changed rsvp time!", nil
				}
			}
		} else if rsvpMenuNum == 1 {
			resp, err := Box.Ask(m, "Are you sure you want to delete this rsvp? [y/N]")
			if err != nil {
				uc <- 0
				return "", err
			}

			if resp == "y" || resp == "Y" {
				Bools.Appts[bwoolNum].Resv = Box.RemoveRsvp(Bools.Appts[bwoolNum].Resv, rsvpNum)
				uc <- 0

				if err := Box.StoreModel(apptsPath, Bools); err != nil {
					log.Fatalln(err)
				}

				return "Successfully deleted rsvp!", nil
			}

			uc <- 0
			return "Rsvp not deleted.", nil
		}
	}

	return "", errors.New("There should be no way you get this error...so good job!")
}

func (botStruct *Bot) Bool(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	uc := make(chan int)
	go Box.Track2Delete(m, uc)

	resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like?\n```\n", "```", Bools.Appts))
	if err != nil {
		uc <- 0
		return nil, err
	}

	fields := []discord.EmbedField{}

	intResp, _ := strconv.Atoi(resp)
	for num, bwool := range Bools.Appts {
		if num == intResp {
			if len(bwool.Resv) > 0 {
				for _, rsvp := range bwool.Resv {
					field := discord.EmbedField{
						Name:   rsvp.User.Username,
						Value:  "Pickup time: " + rsvp.PuTime,
						Inline: true,
					}

					fields = append(fields, field)
				}
			}

			embed := discord.Embed{
				Title:       bwool.Name,
				Description: Box.BuildApptDesc(bwool),
				Fields:      fields,
			}

			uc <- 0
			return &embed, nil
		}
	}

	return nil, errors.New("Bool does not exist, get a list with `" + Prefix + "bools`.")
}

func (botStruct *Bot) Bools(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	if len(Bools.Appts) == 0 {
		return nil, errors.New("No bools currently active. Use `" + Prefix + "newbool` to add a new scheduled bool")
	}

	fields := []discord.EmbedField{}

	for _, Bool := range Bools.Appts {
		field := discord.EmbedField{
			Name:   "`" + Bool.Name + "`",
			Value:  Bool.Desc,
			Inline: false,
		}
		fields = append(fields, field)
	}

	embed := discord.Embed{
		Title:  "Current Bools",
		Fields: fields,
	}

	return &embed, nil
}
