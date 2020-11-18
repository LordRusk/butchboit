package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/voice/voicegateway"
)

var (
	// errors
	NoTreason      = errors.New("Error! No treason is currently being commited in guild. Commit treason with " + Prefix + "treason")
	YesTreason     = errors.New("Error! treason already commited in guild. Close current session with " + Prefix + "kill")
	CannotAutoJoin = errors.New("Cannot auto join channel")
)

func (b *Bot) Treason(m *gateway.MessageCreateEvent) (string, error) {
	if Box.BoomBoxes[m.GuildID] != nil {
		return "", YesTreason
	}

	vst, err := b.Ctx.VoiceState(m.GuildID, m.Author.ID)
	if err != nil {
		return "", errors.New("Cannot join channel! " + m.Author.Username + " not in channel")
	}

	vs, err := Box.JoinChannel(m.GuildID, vst.ChannelID, false, true)
	if err != nil {
		log.Fatal(err)
		return "", errors.New("Cannot join channel!")
	}

	Box.BoomBoxes[m.GuildID] = Box.NewBoomBox(vs)
	return "Successfully commited treason in this channel. Use `" + Prefix + "Play [link]` to play a song", nil
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

	return "Successfully killed current session", nil
}

func (b *Bot) Play(m *gateway.MessageCreateEvent, input bot.RawArguments) (string, error) {
	if Box.BoomBoxes[m.GuildID] == nil {
		return "", NoTreason
	}

	if _, err := b.Ctx.VoiceState(m.GuildID, Box.ID); err != nil {
		return "", errors.New("Cannot play song! Not in channel")
	}

	if Box.BoomBoxes[m.GuildID].Cancel != nil {
		Box.BoomBoxes[m.GuildID].Cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cmd *exec.Cmd
	if isLink := Box.IsLink(string(input)); isLink {
		cmd = exec.CommandContext(ctx,
			"youtube-dl",
			"-f", "250", "-o", "-",
			"--audio-quality", "0",
			"--external-downloader", "ffmpeg",
			"--external-downloader-args", `
        			-hide_banner -threads 1 -loglevel error
        			-c:a copy -f opus -
    			`,
			string(input),
		)

	} else {
		cmd = exec.CommandContext(ctx,
			"youtube-dl",
			"-f", "250", "-o", "-",
			"--default-search", "ytsearch1:",
			"--audio-quality", "0",
			"--external-downloader", "ffmpeg",
			"--external-downloader-args", `
        			-hide_banner -threads 1 -loglevel error
        			-c:a copy -f opus -
    			`,
			string(input),
		)
	}

	// media, info, err := Box.GetVideo(id)
	// if err != nil {
	// 	return "", err
	// }

	oggWriter := Box.NewOggWriter(Box.BoomBoxes[m.GuildID])
	defer oggWriter.Close()
	Box.BoomBoxes[m.GuildID].Cancel = func() { cancel(); oggWriter.Close() }

	cmd.Stdout = oggWriter
	cmd.Stderr = os.Stderr

	done := make(chan error)
	go func() { done <- cmd.Run() }()

	_, err := b.Ctx.SendMessage(m.ChannelID, "Playing", nil)
	if err != nil {
		return "", err
	}

	// start speaking
	if err := Box.BoomBoxes[m.GuildID].Speaking(voicegateway.Microphone); err != nil {
		log.Fatalln("failed to send speaking:", err)
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	// Block until either SIGINT is received OR ffmpeg is done.
	select {
	case <-sig:
	case err = <-done:
	}

	if err != nil {
		log.Println("youtube-dl or ffmpeg failed, exiting.")
	}

	if err := Box.BoomBoxes[m.GuildID].StopSpeaking(); err != nil {
		log.Println("failed to stop speaking:", err)
	}

	return "Finished playing `", nil
}

func (b *Bot) Stop(m *gateway.MessageCreateEvent) error {
	if Box.BoomBoxes[m.GuildID] == nil {
		return NoTreason
	}

	if Box.BoomBoxes[m.GuildID].Cancel != nil {
		if err := Box.BoomBoxes[m.GuildID].StopSpeaking(); err != nil {
			log.Fatalln("failed to send stop speaking:", err)
		}

		Box.BoomBoxes[m.GuildID].Cancel()
		Box.BoomBoxes[m.GuildID].Cancel = nil
	}

	return nil
}
