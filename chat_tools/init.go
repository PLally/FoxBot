package chat_tools

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})

	Module := command.RegisterModule("chat_tools")
	Module.RegisterCommandFunc(">coinflip", coinFlip)
	Module.RegisterCommandFunc(">info", getDiscordObjectInfo)
	source := rand.NewSource(time.Now().UnixNano())
	random = rand.New(source)
}

func coinFlip(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	coinSides := []string{
		"heads",
		"tails",
	}
	side := coinSides[random.Intn(1)]
	return side
}
