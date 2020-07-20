package commands

import "github.com/plally/dgcommand"

func ping(ctx dgcommand.CommandContext) {
	ctx.Reply("Pong!")
}