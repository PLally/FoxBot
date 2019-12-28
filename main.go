package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/plally/discord_modular_bot/modules/chat_tools"
	_ "github.com/plally/discord_modular_bot/modules/debug"
	_ "github.com/plally/discord_modular_bot/modules/nsfw"
	_ "github.com/plally/discord_modular_bot/modules/statistics"
)

var TOKEN = os.Getenv("DISCORD_BOT_TOKEN")

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(os.Stdout)

	session, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Error(err)
		return
	}
	bot := command.NewBot(session)
	bot.EnableModule("debug")
	bot.EnableModule("chat_tools")
	bot.EnableModule("nsfw")
	bot.EnableModule("stats")

	err = session.Open()
	if err != nil {
		log.Error("error opening connection,", err)
		return
	}
	session.UpdateListeningStatus("everything")

	log.Info("Bot Started")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

}
