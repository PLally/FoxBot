package debug

import (
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	Module := command.RegisterModule("debug")
	Module.RegisterCommandFunc(">ping", ping)
	Module.RegisterCommandFunc(">status", botStatus)
	Module.RegisterCommandFunc(">help", helpCommand)
}

func ping(ctx *command.CommandContext) (reply string) {
	return "pong... "
}

func botStatus(ctx *command.CommandContext) (reply string) {
	embed := command.NewEmbed()
	embed.Color = 0x0000ff
	embed.AddField("Active Goroutines", strconv.Itoa(runtime.NumGoroutine()), true)
	ctx.Bot.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed.MessageEmbed)
	return
}

func helpCommand(ctx *command.CommandContext) (reply string)  {
	e := command.NewEmbed()
	e.SetTitle("Command List", "")
	for _, module := range ctx.Bot.EnabledModules {
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
	_, err := ctx.Bot.ChannelMessageSendEmbed(ctx.Message.ChannelID, e.MessageEmbed)
	if err != nil {
		log.Error(err)
	}
	return
}