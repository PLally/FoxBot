package subscriptions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/help"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
	"github.com/spf13/viper"
	"time"
)



func CommandGroup() *dgcommand.CommandGroup{
	var CommandGroup = dgcommand.Group()

	CommandGroup.Default(dgcommand.HandlerFunc(help.DefaultHelpHandler))

	s := subClient{subscription_client.NewSubscriptionClient(viper.GetString("subapi_baseurl"), viper.GetString("subapi_token"))}

	CommandGroup.Desc("Subsribe to updates from websites")

	CommandGroup.Command("list", s.listSubscriptions).
		Use(middleware.Coooldown(5*time.Second, 3)).
		Desc("List subscriptions in a channel")

	CommandGroup.Command("deleteid <id>", s.deleteSubscriptionID).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator)).
		Desc("Delete a subscription using its id")

	CommandGroup.Command("delete <type> <tags>", s.deleteSubscription).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator)).
		Desc("Delete a subscription using its type and tags")

	CommandGroup.Command("add <subtype> [tags...]", s.subscribeCommand).
		Use(
			middleware.RequirePermissions(discordgo.PermissionAdministrator),
			middleware.Coooldown(7*time.Second, 3),
			).
		Desc("Create a subscription\n valid subtypes are e621 and rss \n")

	return CommandGroup
}
