package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/voice/voicegateway"
)

// errors
var NoTreason = errors.New("Error! No treason is currently being commited in guild. Commit treason with " + Prefix + "treason")
var YesTreason = errors.New("Error! treason already commited in guild. Close current session with " + Prefix + "kill")

func (b *Bot) Treason(m *gateway.MessageCreateEvent) (string, error) {
	if Box.BoomBoxes[m.GuildID] != nil {
		return "", YesTreason
	}

	vst, err := b.Ctx.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		log.Println(err)
		return "", errors.New("Cannot join channel! " + m.Author.Username + " not in channel")
	}

	vs, err := Box.JoinChannel(m.GuildID, vst.ChannelID, false, true)
	if err != nil {
		log.Fatal(err)
		return "", errors.New("Cannot join channel!")
	}

	Box.BoomBoxes[m.GuildID] = Box.NewBoomBox(vs)

	// setup queue system
	go func() {
		for Box.BoomBoxes[m.GuildID] != nil {
			media := <-Box.BoomBoxes[m.GuildID].Player

			if len(Box.BoomBoxes[m.GuildID].Queue) != 0 {
				Box.BoomBoxes[m.GuildID].Queue = Box.BoomBoxes[m.GuildID].Queue[1:]
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cmd := exec.CommandContext(ctx,
				"ffmpeg",
				// Streaming is slow, so a single thread is all we need.
				"-hide_banner", "-threads", "1", "-loglevel", "error",
				"-i", "pipe:", "-filter:a", "volume=0.25", "-c:a", "libopus", "-b:a", "64k",
				"-f", "opus", "-",
			)

			oggWriter := Box.NewOggWriter(Box.BoomBoxes[m.GuildID])
			defer oggWriter.Close()
			Box.BoomBoxes[m.GuildID].Cancel = func() { cancel(); oggWriter.Close() }

			cmd.Stdin = media.Stream
			cmd.Stdout = oggWriter
			cmd.Stderr = os.Stderr

			_, err = b.Ctx.SendMessage(m.ChannelID, "Playing `"+media.Title+"`", nil)
			if err != nil {
				log.Println(err)
			}

			// start speaking
			if err := Box.BoomBoxes[m.GuildID].Speaking(voicegateway.Microphone); err != nil {
				log.Println("failed to send speaking:", err)
			}

			if err := cmd.Run(); err != nil {
				log.Println(err)
			}

			_, err = b.Ctx.SendMessage(m.ChannelID, "Finished playing `"+media.Title+"`", nil)
			if err != nil {
				log.Println(err)
			}

			if Box.BoomBoxes[m.GuildID] != nil {
				if len(Box.BoomBoxes[m.GuildID].Player) == 0 {
					_, err = b.Ctx.SendMessage(m.ChannelID, "Finished Queue", nil)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()

	return "Successfully commited treason in this channel. Use `" + Prefix + "Play [Search term || link]` to play a song", nil
}

func (b *Bot) Kill(m *gateway.MessageCreateEvent) (string, error) {
	if Box.BoomBoxes[m.GuildID] == nil {
		return "", NoTreason
	}

	if Box.BoomBoxes[m.GuildID].Cancel != nil {
		Box.BoomBoxes[m.GuildID].Cancel()
	}

	Box.RemoveSession(m.GuildID)
	Box.BoomBoxes[m.GuildID] = nil

	return "Successfully killed current treason session", nil
}

func (b *Bot) Play(m *gateway.MessageCreateEvent, input bot.RawArguments) error {
	if Box.BoomBoxes[m.GuildID] == nil {
		return NoTreason
	}

	if _, err := b.Ctx.VoiceState(m.GuildID, Box.ID); err != nil {
		return errors.New("Cannot play song! Not in channel")
	}

	var id string
	if Box.IsLink(string(input)) {
		id = string(input)
	} else {
		_, err := b.Ctx.SendMessage(m.ChannelID, "Searching `"+string(input)+"`", nil)
		if err != nil {
			return err
		}

		id, err = Box.GetVideoID(string(input))
		if err != nil {
			return err
		}
	}

	media, err := Box.GetVideo(id)
	if err != nil {
		log.Println(err)
	}

	Box.BoomBoxes[m.GuildID].Player <- media

	if len(Box.BoomBoxes[m.GuildID].Player) != 0 {
		Box.BoomBoxes[m.GuildID].Queue = append(Box.BoomBoxes[m.GuildID].Queue, media.Title)

		_, err := b.Ctx.SendMessage(m.ChannelID, "`"+media.Title+"` Added to queue", nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) Skip(m *gateway.MessageCreateEvent) error {
	if Box.BoomBoxes[m.GuildID] == nil {
		return NoTreason
	}

	if Box.BoomBoxes[m.GuildID].Cancel != nil {
		Box.BoomBoxes[m.GuildID].Cancel()
		Box.BoomBoxes[m.GuildID].Cancel = nil
	}

	return nil
}

func (b *Bot) Queue(m *gateway.MessageCreateEvent) (string, error) {
	if Box.BoomBoxes[m.GuildID] == nil {
		return "", NoTreason
	}

	if len(Box.BoomBoxes[m.GuildID].Queue) == 0 {
		return "", errors.New("No songs in queue")
	}

	var builder strings.Builder
	for num, title := range Box.BoomBoxes[m.GuildID].Queue {
		builder.WriteString(strconv.Itoa(num+1) + ": `" + title + "`\n")
	}

	return builder.String(), nil
}
