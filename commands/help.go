package commands

import (
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"strings"
)

type helpGroup dgcommand.CommandGroup


func (g helpGroup) helpCommand(ctx dgcommand.Context) {
	ctx = ctx.(*dgcommand.DiscordContext)
	e := embed.NewEmbed()
	e.SetTitle("Help", "")

	args := strings.Split(ctx.Args()[0], " ")
	group := dgcommand.CommandGroup(g)
	e.AddField("Help", getCommandList(&group, args), false)

	ctx.SendEmbed(e)
}

func getCommandList(h *dgcommand.CommandGroup, args []string) string {

	if len(args) > 0 && args[0] != "" {
		next := args[0]

		args = args[1:]
		nextHandler, ok := h.Commands[next]
		if !ok {
			return "Couldn't Find a handler: " + next
		}

		switch v := nextHandler.(type) {
		case *dgcommand.CommandGroup:
			return getCommandList(v, args)
		case *dgcommand.Command:
			return fmt.Sprintf("`%v` : %v", v.String(), v.Description)
		}
	}

	var out string
	for name, c := range h.Commands {

		switch v := c.(type) {
		case *dgcommand.CommandGroup:
			out += "` "+name + " <subcommand>"+" `"
		case *dgcommand.Command:
			out += "` "+v.String()+" `"
		}
		out += "\n"
	}
	return out
}
