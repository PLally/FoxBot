package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/plally/FoxBot/commands"
	"github.com/plally/dgcommand"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)


// TODO verbose descriptions with examples
// TODO only allow n subscriptions per guild

func main() {
	setupConfig()

	session := makeSession()
	createLogger(session)

	// create and add command handlers
	rootHandler := commands.CommandGroup()

	prefixed := dgcommand.OnPrefix(viper.GetString("prefix"), rootHandler)

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered from fatal error: ", r)
			}
		}()

		ctx := dgcommand.CreatContext(s, m)
		ctx.WithValue("rootHandler", rootHandler)
		prefixed.Handle(ctx)
	})

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
	if err != nil {
		log.Fatal(err)
	}
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(
		0,
		fmt.Sprintf("type `%vhelp` for a list of commands", viper.GetString("prefix")),
	)
}
func makeSession() *discordgo.Session {

	TOKEN := viper.GetString("TOKEN")

	session, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal(err.Error())
	}
	session.AddHandler(onReady)
	err = session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	return session
}

