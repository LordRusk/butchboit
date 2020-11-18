// where butchbot keeps his mutiny.
package boolbox

import (
	"errors"
	"io"
	"strings"

	"github.com/diamondburned/arikawa/voice"
	"github.com/jonas747/ogg"
)

var NonYoutubeLink = errors.New("Error! Not a youtube link!")

// This struct is an abstraction, making it easier to
// have multiple voices on different guilds.
type BoomBox struct {
	*voice.Session
	Cancel func()
}

func (box *Box) NewBoomBox(vs *voice.Session) *BoomBox {
	return &BoomBox{Session: vs}
}

func (box *Box) IsLink(input string) bool {
	plink := strings.Split(input, "/")

	if plink[0] == "https:" || plink[0] == "http:" {
		if plink[1] == "" {
			return true
		}
	}

	return false
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
