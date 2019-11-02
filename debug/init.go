package debug

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
	"fmt"
)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	Module := command.RegisterModule("debug")
	Module.RegisterCommandFunc(">ping", ping)
	Module.RegisterCommandFunc(">status", botStatus)
	Module.RegisterCommandFunc(">help", helpCommand)
}

func ping(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	return "pong... "
}

func botStatus(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	embed := command.NewEmbed()
	embed.Color = 0x0000ff
	embed.AddField("Active Goroutines", strconv.Itoa(runtime.NumGoroutine()), true)
	s.ChannelMessageSendEmbed(event.Message.ChannelID, embed.MessageEmbed)
	return
}

func helpCommand(s *discordgo.Session, event *command.TextCommandEvent) (reply string)  {
	e := command.NewEmbed()
	e.SetTitle("Command List")
	for _, module := range event.Bot.EnabledModules {
		var b strings.Builder
		for _, cmd := range module.Commands {
			b.WriteString(cmd.Name+"\n")
		}
		fieldValue := b.String()
		if len(fieldValue) < 1 {
			continue
		}
		e.AddField(module.Name, fieldValue, false)
	}
	_, err := s.ChannelMessageSendEmbed(event.Message.ChannelID, e.MessageEmbed)
	fmt.Println(err)
	return "TEST"
}