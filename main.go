package main

import (
	"github.com/bwmarrin/discordgo"
	_ "github.com/plally/modular_bot/chat_tools"
	"github.com/plally/modular_bot/command"
	_ "github.com/plally/modular_bot/debug"
	_ "github.com/plally/modular_bot/nsfw"
	_ "github.com/plally/modular_bot/random"
	_ "github.com/plally/modular_bot/statistics"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
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
	bot.EnableModule("random")
	bot.EnableModule("nsfw")
	// bot.EnableModule("stats")

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
