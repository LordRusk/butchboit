package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"

	// "github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	// "github.com/diamondburned/arikawa/v2/voice"
	"github.com/diamondburned/arikawa/v2/voice/voicegateway"
	"github.com/lordrusk/butchbot/boolbox"
)

var (
	// errors
	NoTreason      = errors.New("Error! No treason is currently being commited in guild. Commit treason with " + Prefix + "treason")
	YesTreason     = errors.New("Error! treason already commited in guild. Close current session with " + Prefix + "kill")
	CannotAutoJoin = errors.New("Cannot auto join channel")
)

func (b *Bot) Treason(m *gateway.MessageCreateEvent) (string, error) {
	if _, ok := Box.CheckTreason(m.GuildID); ok {
		return "", YesTreason
	}

	boom := boolbox.NewBoomBox(b.Ctx.State)
	Box.BoomBoxes[m.GuildID] = boom

	return "Successfully commited treason in this channel. Use " + Prefix + "Play [link] to play a song", nil
}

func (b *Bot) Kill(m *gateway.MessageCreateEvent) (string, error) {
	boom, ok := Box.CheckTreason(m.GuildID)
	if !ok {
		return "", NoTreason
	}

	// Box.BoomBoxes[m.GuildID].V.RemoveSession(m.GuildID)
	boom.V.Close()
	Box.BoomBoxes[m.GuildID] = nil

	return "Successfully killed current session", nil
}

func (b *Bot) Join(m *gateway.MessageCreateEvent) (string, error) {
	boom, ok := Box.CheckTreason(m.GuildID)
	if !ok {
		return "", NoTreason
	}

	// fmt.Println's used for debugging
	fmt.Println("0")

	vst, err := b.Ctx.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return "", err
	}

	fmt.Println("1")

	if vst.ChannelID == 0 {
		return "", errors.New("Cannot auto join! User not in channel.")
	}

	fmt.Println("2")

	ch, err := b.Ctx.Channel(m.ChannelID)
	if err != nil {
		log.Fatalln("failed to get channel:", err)
		return "", errors.New("failed to get channel")
	}

	fmt.Println("3")

	voiceSession, err := boom.V.JoinChannel(ch.GuildID, ch.ID, false, true)
	if err != nil {
		log.Fatalln("failed to join channel:", err)
		return "", errors.New("failed to join channel")
	}

	fmt.Println("4")

	boom.VSesh = voiceSession

	return "Joined", nil
}

func (b *Bot) Play(m *gateway.MessageCreateEvent, link string) (string, error) {
	boom, ok := Box.CheckTreason(m.GuildID)
	if !ok {
		return "", NoTreason
	}

	media, err := boom.ParseLink(link)
	if err != nil {
		return "", err
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

	oggWriter := Box.NewOggWriter(boom.VSesh)
	defer oggWriter.Close()
	Box.BoomBoxes[m.GuildID].Cancel = cancel

	cmd.Stdin = media
	cmd.Stdout = oggWriter
	cmd.Stderr = os.Stderr

	done := make(chan error)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	// start speaking
	if err := boom.VSesh.Speaking(voicegateway.Microphone); err != nil {
		log.Fatalln("failed to send speaking:", err)
	}

	go func() { done <- cmd.Run() }()

	return "Playing", nil
}

func (b *Bot) Stop(m *gateway.MessageCreateEvent) (string, error) {
	boom, ok := Box.CheckTreason(m.GuildID)
	if !ok {
		return "", NoTreason
	}

	boom.Cancel()

	return "Stopped", nil
}
