package debug

import (
	"github.com/plally/discord_modular_bot/command"
	"runtime"
	"strconv"
)

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

