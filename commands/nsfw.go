package commands

import (
	"fmt"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"github.com/plally/e621"
	"path"
	"strings"
)

var e6Session = e621.NewSession("e621.net", "FoxBot/0.1",)
func e621Func(ctx dgcommand.CommandContext) {
	resp, err := e6Session.GetPosts("order:random "+ctx.Args[0], 1)
	if err != nil {
		ctx.Error(err)
		return
	}
	posts := resp.Posts
	if len(posts) < 1 {
		ctx.Reply("No posts were found with those tags")
		return
	}
	post := posts[0]

	contentUrl := GetValidContentURL(post)
	description := strings.Builder{}
	for _, artist := range post.Tags.Artist {
		artistString := fmt.Sprintf("[%[1]v](https://e621.net/post?tags=%[1]v), ", artist)
		description.WriteString(artistString)
	}
	if contentUrl != post.File.URL{
		description.WriteString("\n*Click **E621 Post** to view content in its original form*")
	}

	e := embed.NewEmbed()
	e.SetTitle("E621 Post", e6Session.PostUrl(post))
	e.SetImageUrl(contentUrl)
	e.Description = description.String()
	ctx.S.ChannelMessageSendEmbed(ctx.M.ChannelID, e.MessageEmbed)
}

var E621Command =  dgcommand.NewCommand("e621 [tags...]", middleware.Wrap(e621Func, middleware.RequireNSFW()))

func GetValidContentURL(p *e621.Post) string {
	urls := []string{
		p.File.URL,
		p.Sample.URL,
		p.Preview.URL,
	}

	validSuffixes := map[string]bool{
		".gif": true,
		".jpg": true,
		".png": true,
	}

	for _, url := range urls {
		fileExt := path.Ext(url)
		if validSuffixes[fileExt] {
			return url
		}
	}
	return p.File.URL
}
