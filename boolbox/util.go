// this is where butch keeps the rest of his tools.
package boolbox

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"time"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	// "github.com/diamondburned/arikawa/voice"
)

// Timeout for menu's
var timeout = time.Second * 30

// errors
var timeoutErr = errors.New("Error! Timed out wating for response!")
var msgErr = errors.New("Error! Could not send message!")

type Box struct {
	// context left unexported. Always use bot.Ctx.
	ctx *bot.Context
}

func NewBox(ctx *bot.Context) (*Box, error) {
	if ctx == nil {
		return nil, errors.New("Error! No client given!")
	}

	return &Box{ctx: ctx}, nil
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

// get stored model from json file.
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

// Ask is a easy function to get user input
// more than once in a function. Adds ability
// for easy scripting and wizards.
func (box *Box) Ask(m *gateway.MessageCreateEvent, inquire string) (string, error) {
	_, err := box.ctx.SendMessage(m.ChannelID, inquire, nil)
	if err != nil {
		return "", msgErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	v := box.ctx.WaitFor(ctx, func(v interface{}) bool {
		var pass bool
		for pass == false {
			mg, ok := v.(*gateway.MessageCreateEvent)
			if !ok {
				return false
			}

			if mg.ChannelID == m.ChannelID {
				return mg.Author.ID == m.Author.ID
			}
		}

		return false
	})

	if v == nil {
		return "", timeoutErr
	}

	resp := v.(*gateway.MessageCreateEvent)
	return resp.Content, nil
}

// creates handler that logs all message id's. Returns
// a function, when called, will delete all logged
// messages and cancel's handler.
func (box *Box) Track2Delete(targetID discord.ChannelID) func() {
	uc := make(chan struct{})
	mIDa := []discord.MessageID{}

	cancel := box.ctx.AddHandler(func(c *gateway.MessageCreateEvent) {
		if c.ChannelID == targetID {
			mIDa = append(mIDa, c.ID)
		}
	})

	go func() {
		_ = <-uc
		close(uc)

		cancel()
		if err := box.ctx.DeleteMessages(targetID, mIDa); err != nil {
			log.Println(err)
		}
	}()

	return func() { uc <- struct{}{} }
}
