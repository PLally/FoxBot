package commands

import (
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"strings"
)

var helpRootHandler *dgcommand.CommandRoutingHandler

func helpCommand(ctx dgcommand.CommandContext) {
	e := embed.NewEmbed()
	e.SetTitle("Help", "")

	args := strings.Split(ctx.Args[0], " ")

	e.AddField("Help", getCommandList(helpRootHandler, args), false)
	ctx.S.ChannelMessageSendEmbed(ctx.M.ChannelID, e.MessageEmbed)
}

var HelpCommand = dgcommand.NewCommand("help [command...]", helpCommand)

func getCommandList(h *dgcommand.CommandRoutingHandler, args []string) string {

	if len(args) > 0 && args[0] != "" {
		next := args[0]

		args = args[1:]
		nextHandler, ok := h.Commands()[next]
		if !ok {
			return "Couldn't Find a handler: " + next
		}

		switch v := nextHandler.(type) {
		case *dgcommand.CommandRoutingHandler:
			return getCommandList(v, args)
		case *dgcommand.Command:
			return v.String()
		}
	}

	var out string
	for name, c := range h.Commands() {

		switch v := c.(type) {
		case *dgcommand.CommandRoutingHandler:
			out += name + " <subcommand>"
		case *dgcommand.Command:
			out += v.String()
		}
		out += "\n"
	}
	return out
}
