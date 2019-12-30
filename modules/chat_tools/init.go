package chat_tools

import (
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	Module := command.RegisterModule("chat_tools")

	Module.RegisterCommandFunc("coinflip", coinFlip).
		SetUsage("").
		SetDescription("flips a coin, returns either heads or tails")

	Module.RegisterCommandFunc("info", getDiscordObjectInfo).
		SetUsage("[user...]").
		SetDescription("Gets info about users")

	Module.RegisterCommandFunc("random", randomCommand).
		SetUsage("[set]").
		SetDescription("Gets a random item from a set.\nvalid values for set are:\n\tfox\n\tuser")
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

func coinFlip(ctx *command.CommandContext) (reply string) {
	coinSides := []string{
		"heads",
		"tails",
	}
	side := coinSides[random.Intn(1)]
	return side
}
