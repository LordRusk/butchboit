package boolbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
)

// timeout for menus
var timeout = time.Second * 30

// errors
var timeoutErr = errors.New("Error: Timed out waiting for response!")
var msgErr = errors.New("Error: Could not send message!")

type Box struct {
	// context left unexported. Always use bot.Context
	ctx *bot.Context // context must not be embeded
	*discord.User
	BoomBoxes map[discord.GuildID]*BoomBox
}

func NewBox(ctx *bot.Context) (*Box, error) {
	if ctx == nil {
		return nil, errors.New("Error: Co client given!")
	}

	me, err := ctx.Me()
	if err != nil {
		return nil, fmt.Errorf("Error getting bot info: 'ctx.Me()': %s\n", err)
	}

	return &Box{ctx: ctx, User: me, BoomBoxes: make(map[discord.GuildID]*BoomBox)}, nil
}

// sends a message and waits for a respponse
func (box *Box) Ask(m *gateway.MessageCreateEvent, inquire string, timeoutMulti time.Duration) (string, error) {
	_, err := box.ctx.SendMessage(m.ChannelID, inquire, nil)
	if err != nil {
		return "", msgErr
	}

	if timeoutMulti <= 0 {
		timeoutMulti = 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*timeoutMulti)
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

// tracks all messages sent,
// and deletes once returned
// function is ran
// work in progress
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

// store a model in a json file
func StoreModel(path string, model interface{}) error {
	jsonBytes, err := json.MarshalIndent(model, "", "	")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, jsonBytes, 0666); err != nil {
		return err

	}

	return nil
}

// get stored model from json file.
func GetStoredModel(path string, model interface{}) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, model); err != nil {
		return err
	}

	return nil
}
