package random

import (
	"github.com/plally/FoxBot/help"
	"github.com/plally/dgcommand"
)



func CommandGroup() *dgcommand.CommandGroup {
	var CommandGroup = dgcommand.Group()

	CommandGroup.Default(dgcommand.HandlerFunc(help.DefaultHelpHandler))

	CommandGroup.Desc("get a random item")

	CommandGroup.Command("fox", randomFox).
		Desc("a random picture of a fox")
	CommandGroup.Command("cat", randomCat).
		Desc("a random picture of a cat")
	CommandGroup.Command("user", randomUser).
		Desc("a random user in this guild")

	return CommandGroup
}
