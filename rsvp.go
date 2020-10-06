package main

import (
	"errors"
	"strconv"
	"strings"

	//"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

type Bool struct {
	Name  string
	Date  string
	Time  string
	Desc  string
	RSVPS []discord.User
}

var (
	Bools = map[string]Bool{}

	/* errors */
	nameErr = errors.New("Error! Invalid name.")
	timeErr = errors.New("Error! Invalid time.")
	dateErr = errors.New("Error! Invalid date.")
)

func (botStruct *Bot) NewBool(m *gateway.MessageCreateEvent, boolName string, boolTime string, boolDate string, boolDescription string) (string, error) {
	newBool := Bool{}

	/* check the name */
	if boolName == "" {
		return "", nameErr
	}

	/* check the time */
	pBoolTime := strings.Split(boolTime, ":")
	if len(pBoolTime) != 2 {
		return "", timeErr
	}
	_, firstErr := strconv.Atoi(pBoolTime[0])
	_, secondErr := strconv.Atoi(pBoolTime[1])
	if firstErr != nil || secondErr != nil {
		return "", timeErr
	}

	/* check the date */
	pBoolDate := strings.Split(boolDate, "/")
	if len(pBoolDate) != 2 && len(pBoolDate) != 3 {
		return "", dateErr
	}
	_, firstErr = strconv.Atoi(pBoolDate[0])
	_, secondErr = strconv.Atoi(pBoolDate[1])
	if firstErr != nil || secondErr != nil {
		return "", dateErr
	}
	if len(pBoolDate) == 3 {
		_, thirdErr := strconv.Atoi(pBoolTime[1])
		if thirdErr != nil {
			return "", dateErr
		}
	}

	newBool = Bool{
		Name:  boolName,
		Date:  boolDate,
		Time:  boolTime,
		Desc:  boolDescription,
		RSVPS: []discord.User{},
	}

	Bools[newBool.Name] = newBool

	return "New bool successfully added!", nil
}

func (botStruct *Bot) Bools(*gateway.MessageCreateEvent) (*discord.Embed, error) {
	if len(Bools) == 0 {
		return nil, errors.New("No bools currently active. Use `" + Prefix + "newBool` to add a new scheduled bool")
	}

	fields := []discord.EmbedField{}

	for _, Bool := range Bools {
		field := discord.EmbedField{Name: Bool.Name, Value: Bool.Desc, Inline: false}
		fields = append(fields, field)
	}

	embed := discord.Embed{
		Title:  "Current Bools",
		Fields: fields,
	}

	return &embed, nil
}

func (botStruct *Bot) BoolInfo(m *gateway.MessageCreateEvent, selectedBool string) (*discord.Embed, error) {
	if _, ok := Bools[selectedBool]; !ok {
		return nil, errors.New("Bool does not exist, get a list with `" + Prefix + "bools`.")
	}

	var desc strings.Builder

	desc.WriteString("**Time: ")
	desc.WriteString(Bools[selectedBool].Time)
	desc.WriteString("\nDate: ")
	desc.WriteString(Bools[selectedBool].Date)
	desc.WriteString("**\n")
	desc.WriteString(HelpDivider)
	desc.WriteString(Bools[selectedBool].Desc)
	desc.WriteString("\n")

	if len(Bools[selectedBool].RSVPS) > 0 {
		desc.WriteString(HelpDivider)
		desc.WriteString("**Boolers that have RSVP'd**\n")
		for _, booler := range Bools[selectedBool].RSVPS {
			desc.WriteString(booler.Username)
			desc.WriteString("\n")
		}
	}

	embed := discord.Embed{
		Title:       Bools[selectedBool].Name,
		Description: desc.String(),
	}

	return &embed, nil
}

func (botStruct *Bot) RemoveBool(m *gateway.MessageCreateEvent, selectedBool string) (string, error) {
	if _, ok := Bools[selectedBool]; ok {
		delete(Bools, selectedBool)

		return "Bool successfully removed.", nil
	} else {
		return "", errors.New("Bool does not exist, get a list with `" + Prefix + "bools`.")

	}
}

func (botStruct *Bot) Rsvp(m *gateway.MessageCreateEvent, selectedBool string) (string, error) {
	if Bool, ok := Bools[selectedBool]; ok {
		Bool.RSVPS = append(Bool.RSVPS, m.Author)
		Bools[selectedBool] = Bool

		return "Successfully RSVP'd.", nil
	} else {
		return "", errors.New("Bool does not exist, get a list with `" + Prefix + "bools`.")

	}
}
