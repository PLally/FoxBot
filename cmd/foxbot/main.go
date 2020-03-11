package main

import (
	"flag"
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
var (
	shouldMigrate = flag.Bool("migrate", false, "perform database migration")
)
func main() {
	setupConfig()

	session := makeSession()
	db := setupDb()
	
	subtypes.RegisterE621()
	subtypes.RegisterRSS()
	desttypes.RegisterDiscord(session)
	//database setup

	// create and add command handlers
	rootHandler := dgcommand.NewCommandHandler()

	getPrefix := func(dgcommand.CommandContext) string { return viper.GetString("prefix") }
	commands.RegisterCommands(rootHandler, db)

	prefixedRootHandler := dgcommand.WithPrefix(rootHandler, getPrefix)

	session.AddHandler(dgcommand.DiscordHandle(prefixedRootHandler))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}

func setupConfig() {
	viper.SetEnvPrefix("FOX_BOT")
	viper.SetConfigName("foxbot_config")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/foxbot/")
	viper.AddConfigPath("$HOME/.foxbot")
	viper.AddConfigPath(".")
	viper.SetDefault("prefix", ">")
	err := viper.ReadInConfig()
	if err != nil { log.Fatal(err) }
}

func makeSession() *discordgo.Session {

	TOKEN := viper.GetString("TOKEN")

	session, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	return session
}

func setupDb() *gorm.DB {
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
	if *shouldMigrate {
		database.Migrate(db)
		os.Exit(0)
	}

	return db
}