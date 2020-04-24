package commands

import "github.com/plally/dgcommand"

func ping(ctx dgcommand.Context) {
	ctx.Reply("Pong!")
}