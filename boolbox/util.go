package boolbox

// This file is where all utilities and
// helper functions of boolbox call home.

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

var (
	// Timeout for menu's
	timeout = time.Second * 30

	// errors
	timeoutErr = errors.New("Error! Timed out wating for response!")
	msgErr     = errors.New("Error! Could not send message!")
)

// initialize a new Box
func NewBox(ctx *bot.Context) (*Box, error) {
	if ctx == nil {
		return nil, errors.New("Error! No client given!")
	}

	return &Box{Ctx: ctx}, nil
}

// store a model in a json file
func (box *Box) StoreModel(path string, model interface{}) error {
	jsonBytes, err := json.MarshalIndent(model, "", "	")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, jsonBytes, 0777); err != nil {
		return err

	}

	return nil
}

// get stored appointments from json file.
// returns blank Appointments for simplicities
// sake.
func (box *Box) GetStoredModel(path string, model interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, model); err != nil {
		return err
	}

	return nil
}

// remove an item from an array of interface{}
func (box *Box) RemoveAppointment(s []Appointment, i int) []Appointment {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// remove an item from an array of interface{}
func (box *Box) RemoveRsvp(s []Rsvp, i int) []Rsvp {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// helpful function to check if
// an inputted date is valid
func (box *Box) CheckDate(input string) error {
	pDate := strings.Split(input, "/")
	if len(pDate) == 2 {
		_, firstErr := strconv.Atoi(pDate[0])
		_, secondErr := strconv.Atoi(pDate[1])
		if firstErr == nil || secondErr == nil {
			return nil
		}
	}
	if len(pDate) == 3 {
		_, thirdErr := strconv.Atoi(pDate[2])
		if thirdErr != nil {
			return nil
		}
	}

	return errors.New("Invalid date")
}

// helpful function to check if
// an inputted time is valid
func (box *Box) CheckTime(input string) error {
	pTime := strings.Split(input, ":")
	if len(pTime) == 2 {
		_, firstErr := strconv.Atoi(pTime[0])
		_, secondErr := strconv.Atoi(pTime[1])
		if firstErr == nil || secondErr == nil {
			return nil
		}
	}

	return errors.New("Invalid time")
}

// gives an array with ints referring
// to, in order, the names of fields
// from Appointment. Helpful for scripting.
func (box *Box) GetApptSects() []string {
	s := make([]string, 5)

	s[0] = "Name"
	s[1] = "Date"
	s[2] = "Time"
	s[3] = "Decs"
	s[4] = "[]Rsvp"

	return s
}

// Ask is a easy function to get user input
// more than once in a function. Adds ability
// for easy scripting and wizards. Returns the
// discord.MessageID's of all messages sent.
func (box *Box) Ask(m *gateway.MessageCreateEvent, inquire string) (string, error) {
	_, err := box.Ctx.SendMessage(m.ChannelID, inquire, nil)
	if err != nil {
		return "", msgErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	v := box.Ctx.WaitFor(ctx, func(v interface{}) bool {
		mg, ok := v.(*gateway.MessageCreateEvent)
		if !ok {
			return false
		}

		return mg.Author.ID == m.Author.ID
	})

	if v == nil {
		return "", timeoutErr
	}

	resp := v.(*gateway.MessageCreateEvent)
	return resp.Content, nil
}

// extremely helpful function that, tracks
// all new messages sent after initilization,
// until it recieves a signal from uc, when it
// does, it deletes all messages tracked. Useful
// for interactive bot scripts to keep channels
// from looking ugly with a bunch of leftover menus.
func (box *Box) Track2Delete(m *gateway.MessageCreateEvent, uc chan int) error {
	var pass bool
	ic := make(chan discord.MessageID)
	mIDa := []discord.MessageID{}

	go func() {
		for pass == false {
			for mID := range ic {
				mIDa = append(mIDa, mID)
			}
		}

		close(ic)
	}()

	go func() {
		_ = <-uc

		pass = true
		close(uc)

		box.Ctx.DeleteMessages(m.ChannelID, mIDa)
	}()

	box.Ctx.AddHandler(func(m *gateway.MessageCreateEvent) {
		ic <- m.ID
	})

	select {}
}
