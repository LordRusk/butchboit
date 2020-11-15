// where butchbot keeps his mutiny.
package boolbox

import (
	"errors"
	"io"
	"strings"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/state"
	"github.com/diamondburned/arikawa/voice"
	"github.com/jonas747/ogg"
	"github.com/kkdai/youtube"
)

var InvalidLink = errors.New("Error! Invalid link. Please try another.")

type BoomBox struct {
	V      *voice.Voice
	VSesh  *voice.Session
	Cancel func()
}

func NewBoomBox(s *state.State) *BoomBox {
	return &BoomBox{
		V: voice.NewVoice(s),
	}
}

// attached to Box because of needed map.
func (box *Box) CheckTreason(ID discord.GuildID) (boom *BoomBox, ok bool) {
	boom, ok = box.BoomBoxes[ID]

	return
}

// returns io.Reader of media, struct of media information,
// and an error.
func (boom *BoomBox) ParseLink(link string) (io.Reader, error) {
	plink := strings.Split(link, "/")
	if len(plink) != 4 {
		return nil, InvalidLink
	}

	plink = strings.Split(plink[3], "=")
	if len(plink) != 2 && len(plink) != 3 {
		return nil, InvalidLink
	}

	videoID := plink[1]
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	resp, err := client.GetStream(video, &video.Formats[0])
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// OggWriter is used to play sound through voice.
type OggWriter struct {
	pr    *io.PipeReader
	pw    *io.PipeWriter
	errCh chan error
}

func (box *Box) NewOggWriter(w io.Writer) *OggWriter {
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)

	go func() {
		oggDec := ogg.NewPacketDecoder(ogg.NewDecoder(pr))
		for {
			packet, _, err := oggDec.Decode()
			if err != nil {
				errCh <- err
				break
			}
			if _, err := w.Write(packet); err != nil {
				errCh <- err
				break
			}
		}
	}()

	return &OggWriter{
		pw:    pw,
		pr:    pr,
		errCh: errCh,
	}
}

func (w *OggWriter) Write(b []byte) (int, error) {
	select {
	case err := <-w.errCh:
		return 0, err
	default:
		return w.pw.Write(b)
	}
}

func (w *OggWriter) Close() error {
	w.pw.Close()
	return w.pr.Close()
}
