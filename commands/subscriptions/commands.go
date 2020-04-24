package subscriptions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
	"github.com/plally/subscription_api/subscription"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)



func CommandGroup() *dgcommand.CommandGroup{
	var CommandGroup = dgcommand.Group()

	CommandGroup.Desc("Subsribe to updates from websites")

	db := setupDb()
	s := subClient{subscription_client.NewSubscriptionClient(db)}
	CommandGroup.Command("list", s.listSubscriptions).Use(middleware.Coooldown(5*time.Second, 3))

	CommandGroup.Command("delete <subtype> [tags...]", s.deleteSusbcription).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator), middleware.Coooldown(5*time.Second, 3))

	CommandGroup.Command("add <subtype> [tags...]", s.subscribeCommand).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator), middleware.Coooldown(7*time.Second, 3))

	go func() {
		for {
			subscription.CheckOutDatedSubscriptionTypes(db, 100)
			time.Sleep(time.Minute * 15)
		}
	}()

	return CommandGroup
}


func setupDb() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
	)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	db = db.LogMode(false)
	db.SetLogger(logrus.StandardLogger())
	return db
}