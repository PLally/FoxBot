package debug

import (
	"fmt"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"strings"
)

func helpCommand(ctx *command.CommandContext) (reply string)  {
	fmt.Println(ctx.Args)
	if len(ctx.Args) > 1 {
		cmd := ctx.Bot.GetCommand(ctx.Args[1])
		if cmd == nil {
			return fmt.Sprintf("%v is not a valid command", ctx.Args[1])
		}
		if len(cmd.Usage) > 0 {
			return "```Usage: "+cmd.Name+" "+cmd.Usage+"\n"+cmd.Description+"```"
		} else {
			return "```"+cmd.Name+"```"
		}
	}
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