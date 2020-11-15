// this is where butch keeps his 4chan tools.
package boolbox

import (
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/discord"
	"github.com/lordrusk/godesu"
	"jaytaylor.com/html2text"
)

var (
	scheme          = "https"
	imgBaseURL      = "i.4cdn.org"
	defualt4chanURL = "boards.4chan.org"
	baseImageURL    = scheme + "://" + imgBaseURL
)

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
