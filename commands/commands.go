package commands

import (
	"github.com/jinzhu/gorm"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/dgcommand"
)

func RegisterCommands(root *dgcommand.CommandGroup, db *gorm.DB) {
	randomGroup := dgcommand.Group()

	root.AddHandler("random", randomGroup)
	randomGroup.Command("fox", randomFox).
		Desc("a random picture of a fox")
	randomGroup.Command("cat", randomCat).
		Desc("a random picture of a cat")
	randomGroup.Command("user", randomUser).
		Desc("a random user in this guild")

	root.Command("e621 [tags...]", e621Func).Use(middleware.RequireNSFW()).
		Desc("A random picture from e621")
	root.Command("info <object>", objInfoFunc).
		Desc("Gets info about the given discord object")
	root.Command("help [command...]", helpGroup(*root).helpCommand).
		Desc("shows the help message")

	subGroup := dgcommand.Group()
	subGroup.Desc("Subscribe to updates from websites")

	RegisterSubCommands(subGroup, db)
	root.AddHandler("sub", subGroup)
}
