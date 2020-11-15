package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// path to where bools (appointments) are stored
	apptsPath = os.Getenv("HOME") + "/.local/share/bools.json"

	// get bools (appointments)
	Bools = boolbox.Appointments{}
	_     = Box.GetStoredModel(apptsPath, &Bools)

	// inquires
	dateInqDef = "Date?"
	timeInqDef = "Start of pickup time?"
	descInqDef = "Description?"
	rsvpInqDef = "What is your estimated pickup time?"
)

// opts used for user-end menu's
var rsvpNumOpts = "```\n[0] Edit Time\n[1] Delete rsvp\n[2] Exit```"
var apptNumOpts = "```\n[0] Name\n[1] Date\n[2] Time\n[3] Description\n[4] Rsvp's```"

// great demonstrastion of the
// Ask() function.
func (botStruct *Bot) Newbool(m *gateway.MessageCreateEvent) (string, error) {
	var pass bool

	dateInq := dateInqDef
	timeInq := timeInqDef
	descInq := descInqDef

	appointment := boolbox.Appointment{}

	// get the name of the appointment
	resp, err := Box.Ask(m, "Name?")
	if err != nil {
		return "", err
	}

	if resp != "" {
		appointment.Name = resp
	}

	// get the date of the appointment
	for pass == false {
		resp, err := Box.Ask(m, dateInq)
		if err != nil {
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
		return "", err
	}

	if resp != "" {
		appointment.Desc = resp
	}

	Bools.Appts = append(Bools.Appts, appointment)
	Box.StoreModel(apptsPath, Bools)

	return "New bool added! Check for a current list of bools with `" + Prefix + "bools`!", nil
}

func (botStruct *Bot) Removebool(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	var pass bool
	var bwoolNum int
	var builder strings.Builder

	builder.Write([]byte("Which bool would you like?\n```\n"))
	for num, appointment := range Bools.Appts {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(appointment.Name))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	for pass == false {
		resp, err := Box.Ask(m, builder.String())
		if err != nil {
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
		return "", nil
	}

	if resp == "y" || resp == "Y" {
		Bools.Appts = Box.RemoveAppointment(Bools.Appts, bwoolNum)
		if err := Box.StoreModel(apptsPath, Bools); err != nil {
			log.Fatalln(err)
		}

		return "Successfully removed bool!", nil
	}

	return "Bool not removed.", nil
}

// another demonstration on
// the usefullness of Ask()
func (botStruct *Bot) Rsvp(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	var bwoolNum int
	var pass bool
	var builder strings.Builder
	rsvpInq := rsvpInqDef

	builder.Write([]byte("Which bool would you like?\n```\n"))
	for num, appointment := range Bools.Appts {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(appointment.Name))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	for pass == false {
		resp, err := Box.Ask(m, builder.String())
		if err != nil {
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

			return "Successfully un-RSVP'd!", nil
		}
	}

	rsvp := boolbox.Rsvp{User: m.Author}

	pass = false
	for pass == false {
		resp, err := Box.Ask(m, rsvpInq)
		if err != nil {
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

	return "Successfully RSVP'd!", nil
}

func (botStruct *Bot) Editbool(m *gateway.MessageCreateEvent) (string, error) {
	var bwoolNum int
	var rsvpNum int
	var sectNum int
	var rsvpMenuNum int
	var pass bool
	var builder strings.Builder

	builder.Write([]byte("Which bool would you like?\n```\n"))
	for num, appointment := range Bools.Appts {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(appointment.Name))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	for pass == false {
		resp, err := Box.Ask(m, builder.String())
		if err != nil {
			return "", err
		}

		intResp, err := strconv.Atoi(resp)
		if err == nil {
			for num, _ := range Bools.Appts {
				if num == intResp {
					bwoolNum = num
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

	pass = false
	for pass == false {
		resp, err := Box.Ask(m, "Which part of the bool would you like to edit?\n"+apptNumOpts)
		if err != nil {
			return "", nil
		}

		sectNum, _ = strconv.Atoi(resp)
		if sectNum > 0 && sectNum < 4 {
			pass = true
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
			return "", err
		}

		Bools.Appts[bwoolNum].Name = resp

		return "Successfully changed bool name!", nil
	} else if sectNum == 1 {
		resp, err := Box.Ask(m, "What would you like to change the date to?")
		if err != nil {
			return "", err
		}

		if err := Box.CheckDate(resp); err != nil {
			return "", err
		}

		Bools.Appts[bwoolNum].Date = resp

		return "Successfully changed bool date!", nil
	} else if sectNum == 2 {
		resp, err := Box.Ask(m, "What would you like to change the time to?")
		if err != nil {
			return "", err
		}

		if err := Box.CheckTime(resp); err != nil {
			return "", err
		}

		Bools.Appts[bwoolNum].Time = resp

		return "Successfully changed bool time!", nil
	} else if sectNum == 3 {
		resp, err := Box.Ask(m, "What would you like to change the description to?")
		if err != nil {
			return "", err
		}

		Bools.Appts[bwoolNum].Desc = resp

		return "Successfully changed bool description!", nil
	}

	if len(Bools.Appts[bwoolNum].Resv) < 1 {
		return "", errors.New("Nobody has rsvp'd, so there are no rsvp's to edit.")
	}

	builder.Reset()
	builder.Write([]byte("```\n"))
	for num, rsvp := range Bools.Appts[bwoolNum].Resv {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(rsvp.User.Username))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	pass = false
	for pass == false {
		resp, err := Box.Ask(m, builder.String())
		if err != nil {
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

					return "Successfully changed rsvp time!", nil
				}
			}
		} else if rsvpMenuNum == 1 {
			resp, err := Box.Ask(m, "Are you sure you want to delete this rsvp? [y/N]")
			if err != nil {
				return "", err
			}

			if resp == "y" || resp == "Y" {
				Bools.Appts[bwoolNum].Resv = Box.RemoveRsvp(Bools.Appts[bwoolNum].Resv, rsvpNum)
				if err := Box.StoreModel(apptsPath, Bools); err != nil {
					log.Fatalln(err)
				}

				return "Successfully deleted rsvp!", nil
			}

			return "Rsvp not deleted.", nil
		}
	}

	return "", errors.New("There should be no way you get this error...so good job!")
}

func (botStruct *Bot) Pickedup(m *gateway.MessageCreateEvent) (string, error) {
	var bwoolNum int
	var pass bool
	var builder strings.Builder

	builder.Write([]byte("Which bool would you like?\n```\n"))
	for num, appointment := range Bools.Appts {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(appointment.Name))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	for pass == false {
		resp, err := Box.Ask(m, builder.String())
		if err != nil {
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

	pass = false
	var rsvpNum int
	for num, rsvp := range Bools.Appts[bwoolNum].Resv {
		if rsvp.User == m.Author {
			rsvpNum = num
			pass = true
		}
	}

	if pass == false {
		return "", errors.New("Error! You have not rsvp'd!")
	}

	Bools.Appts[bwoolNum].Resv[rsvpNum].PuTime = "*Picked up*"
	return "Marked as picked up.", nil
}

func (botStruct *Bot) Bool(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	var builder strings.Builder

	builder.Write([]byte("```\n"))
	for num, appointment := range Bools.Appts {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(appointment.Name))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte("```"))

	resp, err := Box.Ask(m, builder.String())
	if err != nil {
		return nil, err
	}

	fields := []discord.EmbedField{}

	intResp, err := strconv.Atoi(resp)
	if err == nil {
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

				return &embed, nil
			}
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
