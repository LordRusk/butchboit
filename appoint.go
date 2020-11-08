package main

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// path to where bools (appointments) are stored
	path = os.Getenv("HOME") + "/.local/share/bools.json"

	// get bools (appointments)
	Bools = Box.GetStoredAppointments(path)

	// inquires
	dateInq = "Date?"
	timeInq = "Start of pickup time?"
	descInq = "Description?"
	rsvpInq = "What is your estimated pickup time?"
)

// great demonstrastion of the
// Ask() function.
func (botStruct *Bot) Newbool(m *gateway.MessageCreateEvent) (string, error) {
	var pass bool
	appointment := boolbox.Appointment{}

	// get the name of the appointment
	resp, err := Box.Ask(m, "Name?")
	if err != nil {
		return "", err
	}

	if resp != "" {
		appointment.Name = resp
		pass = true
	}

	// get the date of the appointment
	pass = false
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

	// set inquires back to default
	dateInq = "Date?"
	timeInq = "Time?"
	descInq = "Description?"

	Bools.Appts = append(Bools.Appts, appointment)

	Box.StoreAppointments(path, Bools)
	return "New bool added! Check for a current list of bools with `" + Prefix + "bools`!", nil
}

func (botStruct *Bot) Removebool(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	var bwoolNum int
	var pass bool

	for pass == false {
		resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like?\n```\n", "```", Bools.Appts))
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
			_, err = Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}
		}

	}

	resp, err := Box.Ask(m, "Do you really want to remove that bool? [y/N]")
	if err != nil {
		return "", nil
	}

	if resp != "y" && resp != "Y" {
		return "", errors.New("Bool not removed")
	}

	Bools.Appts = Box.RemoveAppointment(Bools.Appts, bwoolNum)
	Box.StoreAppointments(path, Bools)

	return "Successfully removed bool!", nil
}

// another demonstration on
// the usefullness of Ask()
func (botStruct *Bot) Rsvp(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	var bwoolNum int
	var pass bool

	for pass == false {
		resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like to rsvp for?\n```\n", "```", Bools.Appts))
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
			_, err = Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
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

		// check that time is correctly formatted
		pBoolTime := strings.Split(resp, ":")
		if len(pBoolTime) == 2 {
			_, firstErr := strconv.Atoi(pBoolTime[0])
			_, secondErr := strconv.Atoi(pBoolTime[1])
			if firstErr == nil || secondErr == nil {
				pass = true
				rsvp.PuTime = resp
			}
		}

		timeInq = "Invalid date! Try 7:30, 20:45, etc..."
	}

	rsvpInq = "What is your estimated pickup time?"

	Bools.Appts[bwoolNum].Resv = append(Bools.Appts[bwoolNum].Resv, rsvp)
	Box.StoreAppointments(path, Bools)

	return "Successfully RSVP'd!", nil
}

func (botStruct *Bot) Editbool(m *gateway.MessageCreateEvent) (string, error) {
	var bwoolNum int
	var sectNum int
	var pass bool
	var builder strings.Builder

	for pass == false {
		resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like to edit?\n```\n", "```", Bools.Appts))
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
			_, err = Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
			if err != nil {
				return "", err
			}
		}
	}

	boolFields := Box.GetApptSects()

	builder.Write([]byte("```\n"))
	for i := 0; i < len(boolFields)-1; i++ {
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
			return "", nil
		}

		sectNum, _ = strconv.Atoi(resp)
		for key, _ := range boolFields {
			if key == sectNum {
				pass = true
			}
		}

		if pass == false {
			_, err = Box.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil)
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

	return "", errors.New("There should be no way you get this error...so good job!")
}

func (botStruct *Bot) Bool(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	resp, err := Box.Ask(m, Box.NumApptList("Which bool would you like?\n```\n", "```", Bools.Appts))
	if err != nil {
		return nil, err
	}

	fields := []discord.EmbedField{}

	for num, bwool := range Bools.Appts {
		intResp, _ := strconv.Atoi(resp)
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

	return nil, errors.New("Bool does not exist, get a list with `" + Prefix + "bools`.")
}

func (botStruct *Bot) Bools(_ *gateway.MessageCreateEvent) (*discord.Embed, error) {
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
