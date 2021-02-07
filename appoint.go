package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// path to where appointments are stored
	apptsPath = os.Getenv("HOME") + "/.local/share/bools.json"

	// get appointments
	bools = boolbox.Appointments{}
	_     = boolbox.GetStoredModel(apptsPath, &bools)
)

// inquires
var dateInqDef = "Date? (mm,dd,(yy))"
var timeInqDef = "Start of pickup time?"
var descInqDef = "Description?"
var rsvpInqDef = "What is your estimated pickup time?"

// opts used for user-end menu's
var apptNumOpts = "```\n[0] Name\n[1] Date\n[2] Time\n[3] Description\n[4] Rsvp's\n[5] Exit```"
var rsvpNumOpts = "```\n[0] Edit Time\n[1] Delete rsvp\n[2] Exit```"

func (b *Bot) Newbool(m *gateway.MessageCreateEvent) (string, error) {
	dateInq := dateInqDef
	timeInq := timeInqDef
	descInq := descInqDef

	appointment := boolbox.Appointment{}

	resp, err := box.Ask(m, "Name?", 1)
	if err != nil {
		return "", err
	} else if resp == "" {
		return "", errors.New("Response cannot be ''!")
	}

	for { // get the date
		resp, err := box.Ask(m, dateInq, 1)
		if err != nil {
			return "", err
		}

		if strings.Contains(strings.ToLower(resp), "n/a") {
			appointment.Date.Ud = true
			break
		}

		date, err := boolbox.MakeDate(resp)
		if err != nil {
			dateInq = "Invalid date! Try 7/11, 23/12/2020, etc..."
			continue
		}

		appointment.Date = date
		break
	}

	for { // get the time
		resp, err := box.Ask(m, timeInq, 1)
		if err != nil {
			return "", err
		}

		if strings.Contains(strings.ToLower(resp), "n/a") {
			appointment.Time.Ud = true
			break
		}

		time, err := boolbox.MakeTime(resp)
		if err != nil {
			timeInq = "Invalid date! Try 7:30, 20:45, etc..."
			continue
		}

		appointment.Time = &time
		break
	}

	resp, err = box.Ask(m, descInq, 4) // extended timeout timer
	if err != nil {
		return "", err
	} else if resp == "" {
		return "", errors.New("Response cannot be ''!")
	}

	bools.Appts = append(bools.Appts, appointment)
	if err := boolbox.StoreModel(apptsPath, bools); err != nil {
		logger.Printf("Failed to store appointments: %s\n", err)
	}
	return "New bool added! Check for a current list of bools with `" + *prefix + "bools`!", nil
}

func (b *Bot) Removebool(m *gateway.MessageCreateEvent) (string, error) {
	var builder strings.Builder

	builder.WriteString("Which bool would you like?\n```\n")
	for num, appointment := range bools.Appts {
		builder.WriteString(fmt.Sprintf("[%d] %s\n", strconv.Itoa(num), appointment.Name))
	}
	builder.WriteString("```")

	var bwoolNum int
	for {
		resp, err := box.Ask(m, builder.String(), 1)
		if err != nil {
			return "", err
		}

		iResp, err := strconv.Atoi(resp) // purposefully not check error
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		for num, _ := range bools.Appts {
			if num == iResp {
				bwoolNum = num
				break
			}
		}

		if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range! Try again...\n", nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	resp, err := box.Ask(m, "Do you really want to remove that bool? [y/N]", 1)
	if err != nil {
		return "", nil
	}

	if resp == "y" || resp == "Y" {
		bools.Appts = boolbox.RemoveAppointment(bools.Appts, bwoolNum)
		if err := boolbox.StoreModel(apptsPath, bools); err != nil {
			logger.Printf("Failed to store appointments: %s\n", err)
		}

		return "Successfully removed bool!", nil
	}

	return "Bool not removed", nil
}

func (b *Bot) Rsvp(m *gateway.MessageCreateEvent) (string, error) {
	var builder strings.Builder
	rsvpInq := rsvpInqDef

	builder.WriteString("Which bool would you like to rsvp for?\n```\n")
	for num, appointment := range bools.Appts {
		builder.WriteString(fmt.Sprintf("[%d] %s\n", strconv.Itoa(num), appointment.Name))
	}
	builder.WriteString("```")

	var bwoolNum int
	for {
		resp, err := box.Ask(m, builder.String(), 1)
		if err != nil {
			return "", err
		}

		iResp, err := strconv.Atoi(resp) // intentionally not check error
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		for num, _ := range bools.Appts {
			if num == iResp {
				bwoolNum = num
				break
			}
		}

		if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	for num, rsvp := range bools.Appts[bwoolNum].Resv {
		if rsvp.User == m.Author {
			bools.Appts[bwoolNum].Resv = boolbox.RemoveRsvp(bools.Appts[bwoolNum].Resv, num)
			return "Successfully un-RSVP'd", nil
		}
	}

	rsvp := boolbox.Rsvp{User: m.Author}

	for {
		resp, err := box.Ask(m, rsvpInq, 1)
		if err != nil {
			return "", err
		}

		time, err := boolbox.MakeTime(resp)
		if err != nil {
			rsvpInq = "Invalid time! Try 7:30, 20:45, etc..."
			continue
		}
		rsvp.Time = time

		bools.Appts[bwoolNum].Resv = append(bools.Appts[bwoolNum].Resv, rsvp)
		if err := boolbox.StoreModel(apptsPath, bools); err != nil {
			logger.Printf("Failed to store appointments: %s\n", err)
		}

		return "Successfully RSVP'd!", nil

	}
}

func (b *Bot) Editbool(m *gateway.MessageCreateEvent) (string, error) {
	var builder strings.Builder
	builder.WriteString("Which bool would you like?\n```\n")
	for num, appointment := range bools.Appts {
		builder.WriteString(fmt.Sprintf("[%d] %s\n", strconv.Itoa(num), appointment.Name))
	}
	builder.WriteString("```")

	var bwoolNum int
	for {
		resp, err := box.Ask(m, builder.String(), 1)
		if err != nil {
			return "", err
		}

		iResp, err := strconv.Atoi(resp) // purposefully not check error
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		for num, _ := range bools.Appts {
			if num == iResp {
				bwoolNum = num
				break
			}
		}

		if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	var sNum int
	for {
		resp, err := box.Ask(m, fmt.Sprintf("Which part of the bool would you like to edit?\n%s\n", apptNumOpts), 1)
		if err != nil {
			return "", err
		}

		sNum, err = strconv.Atoi(resp)
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, fmt.Sprintf("Ivalid option: %s\n", err), nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		if sNum < 0 && sNum > 5 {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		break
	}

	if sNum == 5 {
		return "Exited", nil
	} else if sNum == 0 {
		resp, err := box.Ask(m, "What would you like to change the name to?", 1)
		if err != nil {
			return "", err
		}

		bools.Appts[bwoolNum].Name = resp
		return "Successfully changed bool name!", nil
	} else if sNum == 1 {
		for {
			resp, err := box.Ask(m, "What would you like to change the date to?", 1)
			if err != nil {
				return "", err
			}

			if strings.Contains(strings.ToLower(resp), "n/a") {
				bools.Appts[bwoolNum].Date.Ud = true
				return "Successfully changed bool date!", nil
			}

			date, err := boolbox.MakeDate(resp)
			if err != nil {
				if _, err := b.Ctx.SendMessage(m.ChannelID, "Invalid date!", nil); err != nil {
					logger.Printf("Failed to send message: %s\n", err)
				}

				continue
			}

			bools.Appts[bwoolNum].Date = date
			return "Successfully changed bool date!", nil
		}
	} else if sNum == 2 {
		for {
			resp, err := box.Ask(m, "What would you like to change the time to?", 1)
			if err != nil {
				return "", err
			}

			if strings.Contains(strings.ToLower(resp), "n/a") {
				bools.Appts[bwoolNum].Time.Ud = true
				return "Successfully changed bool time!", nil
			}

			time, err := boolbox.MakeTime(resp)
			if err != nil {
				if _, err := b.Ctx.SendMessage(m.ChannelID, "Invalid time! Try 7:30, 20:45, etc...", nil); err != nil {
					logger.Printf("Failed to send message: %s\n", err)
				}

				continue
			}

			bools.Appts[bwoolNum].Time = &time
			return "Successfully changed bool time!", nil
		}
	} else if sNum == 3 {
		resp, err := box.Ask(m, "What would you like to change the description to?", 4)
		if err != nil {
			return "", err
		}

		bools.Appts[bwoolNum].Desc = resp
		return "Successfully changed bool description!", nil
	}

	if len(bools.Appts[bwoolNum].Resv) < 1 {
		return "", errors.New("No rsvp's to edit!")
	}

	builder.Reset()
	builder.WriteString("```\n")
	for num, rsvp := range bools.Appts[bwoolNum].Resv {
		builder.WriteString(fmt.Sprintf("[%d] %s\n", strconv.Itoa(num), rsvp.User.Username))
	}
	builder.WriteString("```")

	var rsvpNum int
	for {
		resp, err := box.Ask(m, builder.String(), 1)
		if err != nil {
			return "", err
		}

		iResp, err := strconv.Atoi(resp)
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		for num, _ := range bools.Appts[bwoolNum].Resv {
			if num == iResp {
				rsvpNum = num
				break
			}
		}

		if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	for {
		resp, err := box.Ask(m, fmt.Sprintf("What would you like to do to the rsvp?\n%s", rsvpNumOpts), 1)
		if err != nil {
			return "", err
		}

		sNum, err = strconv.Atoi(resp)
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		if sNum >= 0 && sNum <= 2 {
			break
		}

		if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	if sNum == 2 {
		return "Exited", nil
	} else if sNum == 0 {
		for {
			resp, err := box.Ask(m, "What would you like the new pickup time to be?", 1)
			if err != nil {
				return "", err
			}

			time, err := boolbox.MakeTime(resp)
			if err != nil {
				if _, err := b.Ctx.SendMessage(m.ChannelID, "Invalid time! Try 7:30, 20:45, etc...", nil); err != nil {
					logger.Printf("Failed to send message: %s\n", err)
				}

				continue
			}

			bools.Appts[bwoolNum].Resv[rsvpNum].Time = time
			return "Successfully changed rsvp time!", nil
		}
	} else if sNum == 1 {
		resp, err := box.Ask(m, "Are you sure you want to delete this rsvp? [y/N]", 1)
		if err != nil {
			return "", err
		}

		if resp == "y" || resp == "Y" {
			bools.Appts[bwoolNum].Resv = boolbox.RemoveRsvp(bools.Appts[bwoolNum].Resv, rsvpNum)
			if err := boolbox.StoreModel(apptsPath, bools); err != nil {
				logger.Printf("Failed to store appointments: %s\n", err)
			}

			return "Successfully deleted rsvp!", nil
		}

		return "Rsvp not deleted", nil
	}

	return "", fmt.Errorf("There should be no way to get this error...so good job! Let %s know!", AUTHOR)
}

func (b *Bot) Pickedup(m *gateway.MessageCreateEvent) (string, error) {
	var builder strings.Builder
	builder.WriteString("Wich bool would you like?\n```\n")
	for num, appointment := range bools.Appts {
		builder.WriteString(fmt.Sprintf("[%d] %s\n", strconv.Itoa(num), appointment.Name))
	}
	builder.WriteString("```")

	var bwoolNum int
	for {
		resp, err := box.Ask(m, builder.String(), 1)
		if err != nil {
			return "", err
		}

		iResp, err := strconv.Atoi(resp)
		if err != nil {
			if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			continue
		}

		for num, _ := range bools.Appts {
			if num == iResp {
				bwoolNum = num
				break
			}
		}

		if _, err := b.Ctx.SendMessage(m.ChannelID, "Choice out of range!, try again...\n", nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	var pass bool
	var rsvpNum int
	for num, rsvp := range bools.Appts[bwoolNum].Resv {
		if rsvp.User == m.Author {
			rsvpNum = num
			pass = !pass
		}
	}

	bools.Appts[bwoolNum].Resv[rsvpNum].Pickedup = true
	return "marked as picked up", nil
}

func (b *Bot) Bool(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	if len(bools.Appts) == 0 {
		return nil, fmt.Errorf("No bools currently active. use `%snewbool` to ass a new scheduled bool", *prefix)
	}

	var builder strings.Builder
	builder.WriteString("```\n")
	for num, appointment := range bools.Appts {
		builder.WriteString(fmt.Sprintf("[%d] %s\n", strconv.Itoa(num), appointment.Name))
	}
	builder.WriteString("```")

	resp, err := box.Ask(m, builder.String(), 1)
	if err != nil {
		return nil, err
	}

	fields := []discord.EmbedField{}

	iResp, err := strconv.Atoi(resp)
	if err == nil {
		for num, bwool := range bools.Appts {
			if num == iResp {
				if len(bwool.Resv) > 0 {
					for _, rsvp := range bwool.Resv {
						field := discord.EmbedField{
							Name:   rsvp.User.Username,
							Inline: true,
						}

						if rsvp.Pickedup == true {
							field.Value = "*Picked Up*"
						} else {
							field.Value = fmt.Sprintf("Pickup time: %s", boolbox.BuildTime(&boolbox.Time{Time: rsvp.Time.Time}))
						}

						fields = append(fields, field)
					}

					return &discord.Embed{
						Title:       bwool.Name,
						Description: boolbox.BuildApptDesc(bwool),
						Fields:      fields,
					}, nil
				}
			}
		}
	}

	return nil, errors.New("Choice out of range")
}

func (b *Bot) Bools(m *gateway.MessageCreateEvent) (*discord.Embed, error) {
	if len(bools.Appts) == 0 {
		return nil, fmt.Errorf("No bools currently active. use `%snewbool` to ass a new scheduled bool", *prefix)
	}

	fields := []discord.EmbedField{}

	for _, bwool := range bools.Appts {
		field := discord.EmbedField{
			Name:   fmt.Sprintf("`%s`", bwool.Name),
			Value:  bwool.Desc,
			Inline: false,
		}
		fields = append(fields, field)
	}

	return &discord.Embed{
		Title:  "Current Bools",
		Fields: fields,
	}, nil
}
