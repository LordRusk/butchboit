// where butchbot keeps his mutiny.
package boolbox

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/diamondburned/arikawa/voice"
	"github.com/jonas747/ogg"
	"github.com/kkdai/youtube/v2"
)

// Maximum queue length
const MaxQueueLength = 512

var NonYoutubeLink = errors.New("Error! Not a youtube link!")

// base youtube search link
var (
	base        = "www.youtube.com"
	search      = "results?search_query="
	YtSearchURL = scheme + "://" + base + "/" + search
)

// Global youtube client
var uClient = youtube.Client{}

// struct to handle media
type Media struct {
	Stream io.Reader
	*youtube.Video
}

// This struct is an abstraction, making it easier to
// have multiple voices on different guilds.
type BoomBox struct {
	*voice.Session
	Player chan Media
	Cancel func()

	// only used for showing queue
	Queue []string
}

func (box *Box) NewBoomBox(vs *voice.Session) *BoomBox {
	return &BoomBox{Session: vs, Player: make(chan Media, MaxQueueLength)}
}

// a bufio split function that returns whats inbetween two
// two quoutes as a token
func scanQuotes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading qoutes.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != '"' {
			break
		}
	}

	// Scan until qoutes, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '"' {
			return i + width, data[start:i], nil
		}
	}

	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// Request more data.
	return start, nil, nil
}

// Get the video ID from the top search result from youtube.
func (box *Box) GetVideoID(sTerms string) (string, error) {
	sTerms = strings.ReplaceAll(sTerms, " ", "+")

	resp, err := http.Get(YtSearchURL + sTerms)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 1024)
	scanner.Buffer(buf, 512*1024)

	jsonbytes := []byte{}
	for scanner.Scan() {
		if scanner.Text() == "// scraper_data_begin" {
			scanner.Scan()
			jsonbytes = scanner.Bytes()
			break
		}
	}

	scanner = bufio.NewScanner(bytes.NewReader(jsonbytes))
	scanner.Split(scanQuotes)

	var passes int
	for scanner.Scan() {
		if scanner.Text() == "videoId" {
			passes++
			if passes == 2 {
				scanner.Scan()
				scanner.Scan()
				return scanner.Text(), nil
			}
		}
	}

	return "", errors.New("Error! Could not find video id")
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

func (box *Box) GetVideo(videoID string) (Media, error) {
	video, err := uClient.GetVideo(videoID)
	if err != nil {
		return Media{}, err
	}

	resp, err := uClient.GetStream(video, &video.Formats[0])
	if err != nil {
		return Media{}, err
	}

	return Media{Stream: resp.Body, Video: video}, nil
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
