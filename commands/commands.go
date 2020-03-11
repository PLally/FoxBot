package commands

import (
	"github.com/jinzhu/gorm"
	"github.com/plally/dgcommand"
)

func RegisterCommands(rootHandler *dgcommand.CommandRoutingHandler, db *gorm.DB) {
	randomGroup := dgcommand.NewCommandHandler()
	randomGroup.AddHandler(RandomFoxCommand.Name, RandomFoxCommand)
	randomGroup.AddHandler(RandomCatCommand.Name, RandomCatCommand)
	randomGroup.AddHandler(RandomUserCommand.Name, RandomUserCommand)
	rootHandler.AddHandler("random", randomGroup)

	rootHandler.AddHandler(E621Command.Name, E621Command)
	rootHandler.AddHandler(DiscordObjInfoCommand.Name, DiscordObjInfoCommand)

	helpRootHandler = rootHandler
	rootHandler.AddHandler(HelpCommand.Name, HelpCommand)

	subGroup := dgcommand.NewCommandHandler()
	RegisterSubCommands(subGroup, db)
	rootHandler.AddHandler("sub", subGroup)
}