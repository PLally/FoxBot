package debug

import (
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"

)

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	Module := command.RegisterModule("debug")
	Module.RegisterCommandFunc("ping", ping).
		SetUsage("").
		SetDescription("Pong")
	Module.RegisterCommandFunc("status", botStatus).
		SetDescription("Returns some information about the bots status")
	Module.RegisterCommandFunc("help", helpCommand).
		SetUsage("[Command]").
		SetUsage("Returns a list of commands, or usage info for a specific command")

}