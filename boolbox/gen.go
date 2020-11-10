package boolbox

// This file is where all of boolbox's
// generation functions like GenHelpMessage()
// call home.

import (
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/discord"
	"github.com/lordrusk/godesu"
	"jaytaylor.com/html2text"
)

var (
	// defaults
	HelpDivider = "------------\n"
	HelpColor   = "#fafafa"

	// four
	scheme          = "https"
	imgBaseURL      = "i.4cdn.org"
	defualt4chanURL = "boards.4chan.org"
	baseImageURL    = scheme + "://" + imgBaseURL
)

// turned []Appointment into an
// string formatted to be a numbered
// list. may be helpful in scripting.
// writes prefix and suffix as well.
func (box *Box) NumApptList(prefix string, suffix string, appointments []Appointment) string {
	var builder strings.Builder

	builder.Write([]byte(prefix))
	for num, appointment := range appointments {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(appointment.Name))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte(suffix))

	return builder.String()
}

// same function as appove but for
// []Rsvp. Also helpful for scripting.
func (box *Box) NumRsvpList(prefix string, suffix string, resv []Rsvp) string {
	var builder strings.Builder

	builder.Write([]byte(prefix))
	for num, rsvp := range resv {
		builder.Write([]byte("["))
		builder.Write([]byte(strconv.Itoa(num)))
		builder.Write([]byte("] "))
		builder.Write([]byte(rsvp.User.Username))
		builder.Write([]byte("\n"))
	}
	builder.Write([]byte(suffix))

	return builder.String()
}

// generate the help message
func (box *Box) GenHelpMsg(prefix string, botName string, cmdGroupMap map[string]CmdGroup) (*discord.Embed, error) {
	/* generate the help command */
	var helpMsg strings.Builder

	helpMsg.WriteString(HelpDivider)
	helpMsg.WriteString("**Prefix:**  `")
	helpMsg.WriteString(prefix)
	helpMsg.WriteString("`\n")
	helpMsg.WriteString(HelpDivider)
	helpMsg.WriteString("**Commands**\n")
	helpMsg.WriteString(HelpDivider)

	for _, cmdGroup := range cmdGroupMap {
		helpMsg.WriteString("***")
		helpMsg.WriteString(cmdGroup.Name)
		helpMsg.WriteString(" Commands:***\n")
		for _, cmdInfo := range cmdGroup.CmdArr {
			if cmdInfo.State == 1 {
				helpMsg.WriteString("__[ Work In Progress ]__ ")
			} else if cmdInfo.State == 2 {
				helpMsg.WriteString("~~")
			}
			helpMsg.WriteString("**")
			helpMsg.WriteString(cmdInfo.Cmd)
			helpMsg.WriteString("**")
			for i := 0; i < len(cmdInfo.Args); i++ {
				helpMsg.WriteString(" [ ")
				if cmdInfo.Args[i].IsOptional == true {
					helpMsg.WriteString("*Optional* ")
				}
				helpMsg.WriteString(cmdInfo.Args[i].Name)
				helpMsg.WriteString(" ]")
			}
			helpMsg.WriteString(" -- *")
			helpMsg.WriteString(cmdInfo.Desc)
			helpMsg.WriteString("*")
			if cmdInfo.State == 2 {
				helpMsg.WriteString("~~")
			}
			helpMsg.WriteString("\n")
		}
		helpMsg.WriteString(HelpDivider)
	}

	/* color */
	colorHex, err := strconv.ParseInt((HelpColor)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* make the embed */
	embed := discord.Embed{
		Title:       botName + " Help Page:",
		Description: helpMsg.String(),
		Color:       discord.Color(colorHex),
	}

	return &embed, nil
}

// build appointment description
func (box *Box) BuildApptDesc(appointment Appointment) string {
	var desc strings.Builder

	desc.WriteString("**Time: ")
	desc.WriteString(appointment.Time)
	desc.WriteString("\nDate: ")
	desc.WriteString(appointment.Date)
	desc.WriteString("**\n")
	desc.WriteString(HelpDivider)
	desc.WriteString(appointment.Desc)
	desc.WriteString("\n")

	return desc.String()
}

// Turns a boolbox.Profile into
// a discord.Embed.
func (box *Box) GenProfileEmbed(profile Profile, tagMap map[string]discord.EmbedField) (*discord.Embed, error) {
	/* title */
	var title strings.Builder
	title.WriteString("Bool profile for ")
	title.WriteString(profile.Name)
	title.WriteString(" (AKA: ")
	title.WriteString(profile.Nickname)
	title.WriteString(")")

	/* color */
	colorHex, err := strconv.ParseInt((profile.Color)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* tags */
	fields := []discord.EmbedField{}

	for _, tag := range profile.Tags {
		if _, ok := tagMap[tag]; ok {
			fields = append(fields, tagMap[tag])
		}
	}

	/* make the embed */
	embed := discord.Embed{
		Title:       title.String(),
		Description: profile.Info,
		Color:       discord.Color(colorHex),
		Fields:      fields,
	}

	return &embed, nil
}

// turn a four post into a
// quite nice looking discord
// embed.
func (box *Box) FourToEmbed(color string, thread godesu.Thread, postNum int) (*discord.Embed, error) {
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
	colorHex, err := strconv.ParseInt((color)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	/* get the image */
	image := discord.EmbedImage{
		URL: string(baseImageURL + "/" + thread.Board + "/" + strconv.FormatInt(post.Tim, 10) + post.Ext),
	}

	/* get the thread URL */
	fields := []discord.EmbedField{
		discord.EmbedField{
			Name:   "Thread URL",
			Value:  string(scheme + "://" + defualt4chanURL + "/" + thread.Board + "/thread/" + strconv.FormatInt(int64(posts[0].No), 10)),
			Inline: true,
		},
	}

	/* make it into a discord.Embed */
	embed := discord.Embed{
		Title:       title.String(),
		Description: description,
		Color:       discord.Color(colorHex),
		Image:       &image,
		Fields:      fields,
	}

	return &embed, nil
}
