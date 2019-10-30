package debug

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/modular_bot/command"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	Module := command.RegisterModule("debug")
	Module.RegisterCommandFunc(">ping", ping)
	Module.RegisterCommandFunc(">status", botStatus)
	Module.RegisterCommandFunc(">args", getArgs)
}

func getArgs(s *discordgo.Session, event *command.TextCommandEvent) string {
	return event.Args
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
