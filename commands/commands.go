package commands

import (
	"github.com/jinzhu/gorm"
	"github.com/plally/dgcommand"
	"github.com/sirupsen/logrus"
)

func RegisterCommands(rootHandler *dgcommand.CommandRoutingHandler, db *gorm.DB) {
	randomGroup := dgcommand.NewCommandHandler()
	l := logrus.StandardLogger()
	randomGroup.AddHandler(RandomFoxCommand.Name, withLogging(RandomFoxCommand, l))
	randomGroup.AddHandler(RandomCatCommand.Name, withLogging(RandomCatCommand, l))
	randomGroup.AddHandler(RandomUserCommand.Name, withLogging(RandomUserCommand, l))
	rootHandler.AddHandler("random", randomGroup)

	rootHandler.AddHandler(E621Command.Name, withLogging(E621Command, l))
	rootHandler.AddHandler(DiscordObjInfoCommand.Name, withLogging(DiscordObjInfoCommand, l))

	helpRootHandler = rootHandler
	rootHandler.AddHandler(HelpCommand.Name, withLogging(HelpCommand, l))

	subGroup := dgcommand.NewCommandHandler()
	RegisterSubCommands(subGroup, db)
	rootHandler.AddHandler("sub", subGroup)

}

func withLogging(command *dgcommand.Command,  l *logrus.Logger) *dgcommand.Command {
	originalCallback := command.Callback
	command.Callback = func(ctx dgcommand.CommandContext) {
		l.Infof("Invoked command %v", command.Name)
		originalCallback(ctx)
	}
	return command
}