package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/plally/FoxBot/commands"
	"github.com/plally/FoxBot/subscription_client/desttypes"
	"github.com/plally/FoxBot/subscription_client/subtypes"
	"github.com/plally/dgcommand"
	"github.com/plally/subscription_api/database"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// setup discord bot
	viper.SetEnvPrefix("FOX_BOT")
	viper.SetConfigName("foxbot_config")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/foxbot/")
	viper.AddConfigPath("$HOME/.foxbot")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil { log.Fatal(err) }
	TOKEN := viper.GetString("TOKEN")

	session, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	subtypes.RegisterE621()
	desttypes.RegisterDiscord(session)
	//database setup
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
	)
	db, err :=gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	db = db.LogMode(false)
	db.SetLogger(logrus.StandardLogger())
	database.Migrate(db)
	// create and add command handlers
	rootHandler := dgcommand.NewCommandHandler()

	getPrefix := func(dgcommand.CommandContext) string { return ">>" }
	commands.RegisterCommands(rootHandler, db)

	prefixedRootHandler := dgcommand.WithPrefix(rootHandler, getPrefix)

	session.AddHandler(dgcommand.DiscordHandle(prefixedRootHandler))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}
