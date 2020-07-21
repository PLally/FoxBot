package commands

import (
	"github.com/plally/FoxBot/help"
	"github.com/plally/dgcommand"
	"strings"
)

func helpCommand(ctx dgcommand.CommandContext) {

	args := strings.Split(ctx.Args()[0], " ")

	commandGroup, ok := ctx.Value("rootHandler").(*dgcommand.CommandGroup)
	if !ok { return }

	ctx.Reply(help.GetHelp(commandGroup, args))
}
