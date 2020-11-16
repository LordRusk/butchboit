package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"

	// "github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	// "github.com/diamondburned/arikawa/state"
	// "github.com/diamondburned/arikawa/voice"
	"github.com/diamondburned/arikawa/voice/voicegateway"
	"github.com/lordrusk/butchbot/boolbox"
)

// Global boombox
var BoomBox *boolbox.BoomBox

var (
	// checks
	Active  bool
	InChan  bool
	Playing bool
)

var (
	// errors
	NoTreason      = errors.New("Error! No treason is currently being commited in guild. Commit treason with " + Prefix + "treason")
	YesTreason     = errors.New("Error! treason already commited in guild. Close current session with " + Prefix + "kill")
	CannotAutoJoin = errors.New("Cannot auto join channel")
)

func (b *Bot) Treason(m *gateway.MessageCreateEvent) (string, error) {
	if Active {
		return "", YesTreason
	}

	BoomBox = Box.NewBoomBox(b.Ctx)
	Active = true

	return "Successfully commited treason in this channel. Use " + Prefix + "Play [link] to play a song", nil
}

func (b *Bot) Kill(m *gateway.MessageCreateEvent) (string, error) {
	if !Active {
		return "", NoTreason
	}

	if err := BoomBox.Close(); err != nil {
		return "", errors.New("Failed to close voice!")
	}

	BoomBox = nil
	Active = false

	return "Successfully killed current session", nil
}

func (b *Bot) Join(m *gateway.MessageCreateEvent) error {
	if !Active {
		return NoTreason
	}

	vst, err := b.Ctx.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return errors.New("Cannot join channel! " + m.Author.Username + " not in channel")
	}

	vs, err := BoomBox.JoinChannel(m.GuildID, vst.ChannelID, false, true)
	if err != nil {
		log.Fatal(err)
		return errors.New("Cannot join channel!")
	}

	BoomBox.VS = vs
	InChan = true

	return nil
}

func (b *Bot) Play(m *gateway.MessageCreateEvent, link string) error {
	if !Active {
		return NoTreason
	}

	if BoomBox.Cancel != nil {
		BoomBox.Cancel()
	}

	media, err := BoomBox.ParseLink(link)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"ffmpeg",
		// Streaming is slow, so a single thread is all we need.
		"-hide_banner", "-threads", "1", "-loglevel", "error",
		"-i", "pipe:", "-c:a", "libopus", "-b:a", "64k",
		"-f", "opus", "-",
	)

	oggWriter := Box.NewOggWriter(BoomBox.VS)
	defer oggWriter.Close()
	BoomBox.Cancel = func() { cancel(); oggWriter.Close() }

	cmd.Stdin = media
	cmd.Stdout = oggWriter
	cmd.Stderr = os.Stderr

	done := make(chan error)
	go func() { done <- cmd.Run() }()

	// start speaking
	if err := BoomBox.VS.Speaking(voicegateway.Microphone); err != nil {
		log.Fatalln("failed to send speaking:", err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	// Block until either SIGINT is received OR ffmpeg is done.
	select {
	case <-sig:
	case err = <-done:
	}

	return nil
}

func (b *Bot) Stop(m *gateway.MessageCreateEvent) error {
	if !Active {
		return NoTreason
	}

	if BoomBox.Cancel != nil {
		if err := BoomBox.VS.Speaking(voicegateway.SpeakingFlag(0)); err != nil {
			log.Fatalln("failed to send stop speaking:", err)
		}

		BoomBox.Cancel()
	}

	return nil
}
