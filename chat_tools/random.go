package chat_tools

import (
	log "github.com/sirupsen/logrus"
	"github.com/plally/discord_modular_bot/command"

	"github.com/bwmarrin/discordgo"
	"strings"
	"net/http"
	"mime"
	"io/ioutil"
	"encoding/json"
	"fmt"
)


var randomHandlers = map[string]func(s *discordgo.Session, m *discordgo.MessageCreate) string{
	"fox":  randomFox,
	"cat":  randomCat,
	"user": randomUser,
}

func randomCommand(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	randomType := strings.TrimSpace(strings.Split(event.Message.Content, " ")[1])
	r, ok := randomHandlers[randomType]
	if !ok {
		return "Type not supported"
	}
	return r(s, event.Message)
}

func randomFox(s *discordgo.Session, m *discordgo.MessageCreate) string {
	resp, err := http.Get("https://api.foxorsomething.net/fox")
	if err != nil {
		log.Error(err)
		return "something went wrong"
	}
	contentType := resp.Header.Get("Content-Type")
	extensions, err := mime.ExtensionsByType(contentType)

	if err != nil {
		log.Error(err)
		return "something went wrong"
	}

	if extensions == nil || len(extensions) < 1 {
		return "something went wrong"
	}

	s.ChannelFileSend(m.ChannelID, "fox"+extensions[0], resp.Body)
	return ""
}

func randomCat(s *discordgo.Session, m *discordgo.MessageCreate) string {
	resp, err := http.Get("http://aws.random.cat/meow")
	if err != nil {
		log.Error(err)
		return "something went wrong"
	}
	var dat map[string]string
	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &dat)
	fmt.Println(string(bytes))
	embed := command.NewEmbed()
	embed.SetImageUrl(dat["file"])
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)

	return ""
}

func randomUser(s *discordgo.Session, m *discordgo.MessageCreate) string {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return "something went wrong"
	}
	user := guild.Members[random.Intn(guild.MemberCount)].User
	return user.Username + "#" + user.Discriminator
}
