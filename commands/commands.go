package commands

import (
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/commands/random"
	"github.com/plally/FoxBot/commands/subscriptions"
	"github.com/plally/dgcommand"
)




func CommandGroup() *dgcommand.CommandGroup {
	var CommandGroup = dgcommand.Group()
	CommandGroup.AddHandler("random", random.CommandGroup())

	CommandGroup.Command("e621 [tags...]", e621Func).Use(middleware.RequireNSFW()).
		Desc("A random picture from e621")
	CommandGroup.Command("info <object>", objInfoFunc).
		Desc("Gets info about the given discord object")
	CommandGroup.Command("help [command...]", helpGroup(*CommandGroup).helpCommand).
		Desc("shows the help message")

	CommandGroup.AddHandler("sub", subscriptions.CommandGroup())

	return CommandGroup
}
