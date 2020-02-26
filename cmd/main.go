package main

import (
	"github.com/bwmarrin/discordgo
	"github.com/plally/FoxBot/commands"
	"github.com/plally/dgcommand"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// setup discord bot
	TOKEN := os.Getenv("DISCORD_BOT_TOKEN")
	session, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	// create and add command handlers

	rootHandler := dgcommand.NewCommandHandler()

	getPrefix := func(dgcommand.CommandContext) string { return ">" }
	commands.RegisterCommands(rootHandler)
	prefixedRootHandler := dgcommand.WithPrefix(rootHandler, getPrefix)

	session.AddHandler(dgcommand.DiscordHandle(prefixedRootHandler))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}
