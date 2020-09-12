package main

import (
	"time"
	"errors"
	"strings"
	"strconv"
	"math/rand"

	"jaytaylor.com/html2text"
	"github.com/lordrusk/godesu"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

type IntRange struct {
	min, max int
}

var (
	/* godesu stuff */
	gochan = godesu.New()
	_, boards = gochan.GetBoards()
	r4color = "#006500"
	boardsColor = "#12d7a9"
	scheme = "https"
	imgBaseURL = "i.4cdn.org"
	defualt4chanURL = "boards.4chan.org"
	baseImageURL = scheme+"://"+imgBaseURL
)

/* misc functions */
func (ir *IntRange) NextRandom(r* rand.Rand) int {
    return r.Intn(ir.max - ir.min +1) + ir.min
}

func post2Embed(thread godesu.Thread, postNum int) (*discord.Embed, error) {
	posts := thread.Posts
	post := posts[postNum]

	var title strings.Builder

	title.WriteString("Board: ")
	title.WriteString(thread.Board)
	title.WriteString("\n")

	if posts[0].No != post.No {
		title.WriteString("Thread No. `")
		title.WriteString(strconv.FormatInt(int64(posts[0].No), 10))
		title.WriteString("`\n")
		if posts[0].Sub != "" {
			origPostSub, err := html2text.FromString(posts[0].Sub, html2text.Options{PrettyTables: true})
			if err != nil {
				return nil, err
			}

			title.WriteString("Thread: ")
			title.WriteString(origPostSub)
			title.WriteString("\n")
		}
	}

	title.WriteString(HelpDivider)

	title.WriteString("Post No. `")
	title.WriteString(strconv.FormatInt(int64(post.No), 10))
	title.WriteString("`\n")

	name, err := html2text.FromString(post.Name, html2text.Options{PrettyTables: true})
	if err != nil {
		return nil, err
	}

	title.WriteString("Name: ")
	title.WriteString(name)
	title.WriteString("\n")
	if post.Sub != "" {
		sub, err := html2text.FromString(post.Sub, html2text.Options{PrettyTables: true})
		if err != nil {
			return nil, err
		}

		title.WriteString("Subject: ")
		title.WriteString(sub)
		title.WriteString("\n")
	}

	/* get the description */
	description, err := html2text.FromString(post.Com, html2text.Options{PrettyTables: true})
	if err != nil {
		return nil, err
	}

	/* color */
	colorHex, err := strconv.ParseInt((r4color)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* get the image */
	image := discord.EmbedImage {
		URL:	string(baseImageURL+"/"+thread.Board+"/"+strconv.FormatInt(post.Tim, 10)+post.Ext),
	}

	/* get the thread URL */
	fields := []discord.EmbedField {
		discord.EmbedField {
			Name:	"Thread URL",
			Value:	string(scheme+"://"+defualt4chanURL+"/"+thread.Board+"/thread/"+strconv.FormatInt(int64(posts[0].No), 10)),
			Inline:	true,
		},
	}

	/* make it into a discord.Embed */
	embed := discord.Embed {
		Title:		title.String(),
		Description:	description,
		Color:		discord.Color(colorHex),
		Image:		&image,
		Fields:		fields,

	}

	return &embed, nil
}

/* Bot commands */
func (botStruct *Bot) Post(*gateway.MessageCreateEvent) (*discord.Embed, error) {
	/* backend stuff */
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	/* get a random board */
	boardMap := make(map[int]string)
	for num, board := range boards.All {
		boardMap[num] = board.Board
	}

	irb := IntRange{0, len(boardMap)-1}
	boardName := boardMap[irb.NextRandom(r)]

	board := gochan.Board(boardName)

	/* get a random thread */
	err, catalog := board.GetCatalog()
	if err != nil {
		return nil, err
	}
	threadMap := make(map[int]int)
	for _, page := range catalog.Pages {
		for num, thread := range page.Threads {
			threadMap[num] = thread.No
		}
	}

	irt := IntRange{0, len(threadMap)-1}
	err, thread := board.GetThread(threadMap[irt.NextRandom(r)])
	if err != nil {
		return nil, err
	}

	/* get a random post */
	posts := thread.Posts

	irp := IntRange{0, len(posts)-1}
	postNum := irp.NextRandom(r)

	/*** BUILD THE EMBED ***/
	embed, err := post2Embed(thread, postNum)
	if err != nil {
		return nil, err
	}

	return embed, nil
}

func (botStruct *Bot) Board(m *gateway.MessageCreateEvent, boardName string) (*discord.Embed, error) {
	if boardName == "" {
		return nil, errors.New("Error! No Boards Specified. Please use `!Boards` for a list.")
	}

	/* backend stuff */
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var boardExists bool
	for _, board := range boards.All {
		if boardName == board.Board {
			boardExists = true
		}
	}

	if boardExists != true {
		return nil, errors.New("Error! Board does not exist. Get a list with `!boards`!")
	}

	board := gochan.Board(boardName)

	err, catalog := board.GetCatalog()
	if err != nil {
		return nil, err
	}
	threadMap := make(map[int]int)
	for _, page := range catalog.Pages {
		for num, thread := range page.Threads {
			threadMap[num] = thread.No
		}
	}

	irt := IntRange{0, len(threadMap)-1}
	err, thread := board.GetThread(threadMap[irt.NextRandom(r)])
	if err != nil {
		return nil, err
	}

	/* get a random post */
	posts := thread.Posts

	irp := IntRange{0, len(posts)-1}
	postNum := irp.NextRandom(r)

	/*** BUILD THE EMBED ***/
	embed, err := post2Embed(thread, postNum)
	if err != nil {
		return nil, err
	}

	return embed, nil
}

func (botStruct *Bot) Scope(m *gateway.MessageCreateEvent, boardName string, postNo int) (*discord.Embed, error) {
	/* godesu stuff */
	board := gochan.Board(boardName)

	err, thread := board.GetThread(postNo)
	if err != nil {
		return nil, err
	}

	var postNum int
	for num, post := range thread.Posts {
		if post.No == postNo {
			postNum = num
		}
	}

	/*** BUILD THE EMBED ***/
	embed, err := post2Embed(thread, postNum)
	if err != nil {
		return nil, err
	}

	return embed, nil
}

func (botStruct *Bot) Boards(*gateway.MessageCreateEvent) (*discord.Embed, error) {
	/* godesu stuff */
	_, boards := gochan.GetBoards()

	var description strings.Builder
	for _, board := range boards.All {
		description.WriteString("**")
		description.WriteString(board.Board)
		description.WriteString("**")
		description.WriteString(" ")
	}

	/* color */
	colorHex, err := strconv.ParseInt((boardsColor)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	embed := discord.Embed {
		Title:		"Possible Boards",
		Description:	description.String(),
		Color:		discord.Color(colorHex),
	}

	return &embed, nil
}
