package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/plally/FoxBot/commands"
	"github.com/plally/FoxBot/permissions"
	"github.com/plally/FoxBot/permissions/gormstore"
	"github.com/plally/dgcommand"
	dgcommandparsing "github.com/plally/dgcommand/parsing"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
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
	db := makeDB()

	// create and add command handlers
	rootHandler := commands.CommandGroup()

	prefixed := dgcommand.OnPrefix(viper.GetString("prefix"), rootHandler)

	var store permissions.Store = gormstore.New(db)
	registerCommandGroupPermissions("commands", store, rootHandler)

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered from fatal error: ", r)
			}
		}()

		ctx := dgcommand.CreatContext(s, m)
		ctx.WithValue("rootHandler", rootHandler)
		ctx.WithValue("permissionsStore", store)
		ctx.WithValue("permissionsSnowflake", permissions.GetPermissionsIdentifier(ctx.Message.GuildID, ctx.Message.Author.ID))
		ctx.WithValue("database", db)

		ctx.OnError = func(ctx dgcommand.CommandContext, err error) error {
			if errors.Is(err, dgcommandparsing.ErrMissingArg) {
				ctx.Reply(err.Error())
				return nil
			}
			return err
		}
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
		mw := io.MultiWriter(os.Stdout, file)
		logrus.SetOutput(mw)
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
	log.AddHook(WebhookHook(viper.GetString("logging_webhook")))
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateStatus(
		0,
		fmt.Sprintf("type `%vhelp` for a list of commands", viper.GetString("prefix")),
	)
}

func makeDB() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
	)

	log := logger.New(logrus.StandardLogger(), logger.Config{
		SlowThreshold: 0,
		Colorful:      false,
		LogLevel:      logger.Info,
	})

	db, _ := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: log,
	})

	return db
}
func makeSession() *discordgo.Session {

	session, err := discordgo.New("Bot "+viper.GetString("token"))
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	session.AddHandler(onReady)
	err = session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	return session
}
