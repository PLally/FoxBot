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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)
var (
	shouldMigrate = flag.Bool("migrate", false, "perform database migration")
)
func main() {
	setupConfig()


	session := makeSession()
	createLogger(session)
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
	session.UpdateStatus(
		0,
		fmt.Sprintf("type `%vhelp` for a list of commands", viper.GetString("prefix")),
	)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()
}

func createLogger(session *discordgo.Session) {
	file, err := os.OpenFile(viper.GetString("logfile"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	if viper.GetString("environment") == "dev" {
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetOutput(file)
	}
	session.LogLevel = discordgo.LogDebug
	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		msg = fmt.Sprintf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)
		switch msgL {
		case discordgo.LogError:
			log.Error(msg)
		case discordgo.LogWarning:
			log.Warn(msg)
		case discordgo.LogDebug:
			log.Debug(msg)
		default:
			log.Info(msg)
		}
	}
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
	viper.SetDefault("logfile", "foxbot.log")
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
		log.Info("running database.Migrate")
		database.Migrate(db)
		os.Exit(0)
	}

	return db
}