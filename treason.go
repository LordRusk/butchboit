package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v2/bot"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/voice"
	"github.com/diamondburned/arikawa/v2/voice/voicegateway"
	"github.com/lordrusk/butchbot/boolbox"
)

// errors
var NoTreason = fmt.Errorf("Error! No treason is currently being commited in this guild. Commit treason with `%streason`", *prefix)
var YesTreason = fmt.Errorf("Error! Treason already commited in this guild. Close current session with `%skill`", *prefix)

func (b *Bot) Treason(m *gateway.MessageCreateEvent) (string, error) {
	if box.BoomBoxes[m.GuildID] != nil {
		return "", YesTreason
	}

	v, err := voice.NewSession(b.Ctx.State)
	if err != nil {
		return "", fmt.Errorf("Could not create voice session: %s\n", err)
	}

	vst, err := b.Ctx.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		logger.Printf("Failed to get voice state of %s: %s\n", m.Author.Username, err)
		return "", fmt.Errorf("Cannot join channel! %s not in channel", m.Author.Username)
	}

	box.BoomBoxes[m.GuildID] = box.NewBoomBox(v)

	if err := box.BoomBoxes[m.GuildID].JoinChannel(m.GuildID, vst.ChannelID, false, true); err != nil {
		logger.Printf("Failed to join channel: %s\n", err)
		return "", errors.New("Cannot join channel!")
	}

	// setup queue system
	go func() {
		for box.BoomBoxes[m.GuildID] != nil {
			media := <-box.BoomBoxes[m.GuildID].Player

			box.BoomBoxes[m.GuildID].Playing = &media

			if len(box.BoomBoxes[m.GuildID].Queue) != 0 {
				box.BoomBoxes[m.GuildID].Queue = box.BoomBoxes[m.GuildID].Queue[1:]
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cmd := exec.CommandContext(ctx,
				"ffmpeg",
				// Streaming is slow, so a single thread is all we need.
				"-hide_banner", "-threads", "1", "-loglevel", "error", "-ss",
				strconv.Itoa(media.StartAt), "-i", "pipe:", "-filter:a", "volume=0.25",
				"-c:a", "libopus", "-b:a", "64k", "-f", "opus", "-",
			)

			oggWriter := boolbox.NewOggWriter(box.BoomBoxes[m.GuildID])
			defer oggWriter.Close()
			box.BoomBoxes[m.GuildID].Cancel = func() { cancel(); oggWriter.Close() }

			cmd.Stdin = media.Stream
			cmd.Stdout = oggWriter
			cmd.Stderr = os.Stderr

			if _, err = b.Ctx.SendMessage(m.ChannelID, fmt.Sprintf("Playing `%s`", media.Title), nil); err != nil {
				logger.Printf("Failed to send message: %s\n", err)
			}

			// start speaking
			if err := box.BoomBoxes[m.GuildID].Speaking(voicegateway.Microphone); err != nil {
				logger.Printf("treason: failed to start speaking: %s\n", err)
			}

			if err := cmd.Run(); err != nil {
				logger.Printf("treason: Failed to run cmd: %s\n", err)
			}

			if _, err = b.Ctx.SendMessage(m.ChannelID, fmt.Sprintf("Finished playing `%s`", media.Title), nil); err != nil {
				logger.Printf("Failed to send message: %s\n")
			}

			if box.BoomBoxes[m.GuildID] != nil {
				if len(box.BoomBoxes[m.GuildID].Player) == 0 {
					_, err = b.Ctx.SendMessage(m.ChannelID, "Finished Queue", nil)
					if err != nil {
						log.Println(err)
					}
				}
			}

			if box.BoomBoxes[m.GuildID] != nil {
				box.BoomBoxes[m.GuildID].Playing = nil
			}
		}
	}()

	return fmt.Sprintf("Successfully commited treason on this channel. Use `%splay [Searcg term || link]` to play a song", *prefix), nil
}

func (b *Bot) Kill(m *gateway.MessageCreateEvent) (string, error) {
	if box.BoomBoxes[m.GuildID] == nil {
		return "", NoTreason
	}

	if box.BoomBoxes[m.GuildID].Cancel != nil {
		box.BoomBoxes[m.GuildID].Cancel()
	}

	box.BoomBoxes[m.GuildID] = nil
	return "Successfully killed current treason session", nil
}

func (b *Bot) Play(m *gateway.MessageCreateEvent, input bot.RawArguments) error {
	if box.BoomBoxes[m.GuildID] == nil {
		return NoTreason
	}

	if _, err := b.Ctx.VoiceState(m.GuildID, box.ID); err != nil {
		return errors.New("Cannot play song! Not in channel")
	}

	var id string
	if boolbox.IsLink(string(input)) {
		id = string(input)
	} else {
		_, err := b.Ctx.SendMessage(m.ChannelID, fmt.Sprintf("Searching `%s`", input), nil)
		if err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}

		id, err = boolbox.GetVidID(string(input))
		if err != nil {
			return err
		}
	}

	media, err := boolbox.GetVideo(id)
	if err != nil {
		log.Println(err)
	}

	box.BoomBoxes[m.GuildID].Player <- media

	if len(box.BoomBoxes[m.GuildID].Player) != 0 {
		box.BoomBoxes[m.GuildID].Queue = append(box.BoomBoxes[m.GuildID].Queue, media.Title)

		if _, err := b.Ctx.SendMessage(m.ChannelID, fmt.Sprintf("`%s` Added to queue", media.Title), nil); err != nil {
			logger.Printf("Failed to send message: %s\n", err)
		}
	}

	return nil
}

func (b *Bot) Skip(m *gateway.MessageCreateEvent) error {
	if box.BoomBoxes[m.GuildID] == nil {
		return NoTreason
	}

	if box.BoomBoxes[m.GuildID].Cancel != nil {
		box.BoomBoxes[m.GuildID].Cancel()
		box.BoomBoxes[m.GuildID].Cancel = nil
	}

	return nil
}

func (b *Bot) Queue(m *gateway.MessageCreateEvent) (string, error) {
	if box.BoomBoxes[m.GuildID] == nil {
		return "", NoTreason
	}

	if len(box.BoomBoxes[m.GuildID].Queue) == 0 {
		return "", errors.New("No songs in queue")
	}

	var builder strings.Builder
	for num, title := range box.BoomBoxes[m.GuildID].Queue {
		builder.WriteString(fmt.Sprintf("%s: `%s`\n", strconv.Itoa(num+1), title))
	}

	return builder.String(), nil
}
