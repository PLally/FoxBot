package nsfw

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"os"
)

var E6Session E621Session = E621Session{
	BaseURL:   "https://e621.net/",
	UserAgent: "FoxBot/0.1",
	Client:    &http.Client{},
	Username: os.Getenv("E621_USERNAME"),
	ApiKey: os.Getenv("E621_TOKEN"),
}

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	Module := command.RegisterModule("nsfw")
	Module.RegisterCommandFunc(">e621", e621Command)
}

func e621Command(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	channel, err := s.Channel(event.Message.ChannelID)
	if err != nil {
		return ""
	}
	if !channel.NSFW {
		return "Command Only Available In NSFW channels"
	}
	posts := E6Session.getPosts(strings.Split(event.Args, " "), 1)

	if len(posts) < 1 {
		return "No posts were found with those tags "
	}
	post := posts[0]
	e := command.NewEmbed()
	e.SetTitle("E621")
	e.SetImageUrl(post.FileURL)
	e.MessageEmbed.Description = strings.Join(post.Artist, ", ")

	s.ChannelMessageSendEmbed(event.Message.ChannelID, e.MessageEmbed)
	return ""
}
