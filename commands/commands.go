package commands

import "github.com/plally/dgcommand"

func RegisterCommands(rootHandler *dgcommand.CommandRoutingHandler) {
	randomGroup := dgcommand.NewCommandHandler()
	randomGroup.AddHandler(RandomFoxCommand.Name, RandomFoxCommand)
	randomGroup.AddHandler(RandomCatCommand.Name, RandomCatCommand)
	randomGroup.AddHandler(RandomUserCommand.Name, RandomUserCommand)
	rootHandler.AddHandler("random", randomGroup)

	rootHandler.AddHandler(E621Command.Name, E621Command)
	rootHandler.AddHandler(DiscordObjInfoCommand.Name, DiscordObjInfoCommand)

	helpRootHandler = rootHandler
	rootHandler.AddHandler(HelpCommand.Name, HelpCommand)
}
