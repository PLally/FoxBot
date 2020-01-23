package commands

import (
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"github.com/plally/e621"
	"net/http"
	"os"
	"path"
	"strings"
)

var e6Session = e621.E621Session{
	BaseURL:   "https://e621.net",
	UserAgent: "FoxBot/0.1",
	Client:    &http.Client{},
	Username:  os.Getenv("E621_USERNAME"),
	ApiKey:    os.Getenv("E621_TOKEN"),
}

func e621Func(ctx dgcommand.CommandContext) {
	channel, err := ctx.S.Channel(ctx.M.ChannelID)
	if err != nil {
		ctx.Reply("Something went wrong")
		return
	}
	if !channel.NSFW {
		ctx.Reply("Command Only Available In NSFW channels")
		return
	}
	posts := e6Session.GetPosts(strings.Split("order:random "+ctx.Args[0], " "), 1)

	if len(posts) < 1 {
		ctx.Reply("No posts were found with those tags")
		return
	}
	post := posts[0]

	contentUrl := GetValidContentURL(post)
	description := strings.Builder{}
	for _, artist := range post.Artist {
		artistString := fmt.Sprintf("[%[1]v](https://e621.net/post?tags=%[1]v), ", artist)
		description.WriteString(artistString)
	}
	if contentUrl != post.FileURL {
		description.WriteString("\n*Click **E621 Post** to view content in its original form*")
	}

	e := embed.NewEmbed()
	e.SetTitle("E621 Post", post.PostURL())
	e.SetImageUrl(contentUrl)
	e.Description = description.String()
	ctx.S.ChannelMessageSendEmbed(ctx.M.ChannelID, e.MessageEmbed)
}

var E621Command = dgcommand.NewCommand("e621 [tags...]", e621Func)

func GetValidContentURL(p *e621.E621Post) string {
	urls := []string{
		p.FileURL,
		p.SampleURL,
		p.PreviewURL,
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
	return p.FileURL
}
