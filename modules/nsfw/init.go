package nsfw

import (
	"fmt"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"strings"
)

var E6Session E621Session = E621Session{
	BaseURL:   "https://e621.net",
	UserAgent: "FoxBot/0.1",
	Client:    &http.Client{},
	Username:  os.Getenv("E621_USERNAME"),
	ApiKey:    os.Getenv("E621_TOKEN"),
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	Module := command.RegisterModule("nsfw")
	Module.RegisterCommandFunc(">e621", e621Command)
}

func e621Command(ctx *command.CommandContext) (reply string) {
	channel, err := ctx.Bot.Channel(ctx.Message.ChannelID)
	if err != nil {
		return ""
	}
	if !channel.NSFW {
		return "Command Only Available In NSFW channels"
	}
	posts := E6Session.GetPosts(strings.Split("order:random "+strings.Join(ctx.Args[1:]," "), " "), 1)

	if len(posts) < 1 {
		return "No posts were found with those tags "
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

	e := command.NewEmbed()
	e.SetTitle("E621 Post", post.PostURL())
	e.SetImageUrl(contentUrl)
	e.Description = description.String()
	ctx.Bot.ChannelMessageSendEmbed(ctx.Message.ChannelID, e.MessageEmbed)
	return ""
}

func GetValidContentURL(p *E621Post) string {
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
