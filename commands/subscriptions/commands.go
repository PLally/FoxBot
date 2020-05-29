package subscriptions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
	"time"
)



func CommandGroup() *dgcommand.CommandGroup{
	var CommandGroup = dgcommand.Group()

	CommandGroup.Desc("Subsribe to updates from websites")
	s := subClient{subscription_client.NewSubscriptionClient("http://127.0.0.1:8000")}
	CommandGroup.Command("list", s.listSubscriptions).Use(middleware.Coooldown(5*time.Second, 3))

	/*
	CommandGroup.Command("delete <subtype> [tags...]", s.deleteSusbcription).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator), middleware.Coooldown(5*time.Second, 3))
*/

	CommandGroup.Command("deleteid <id>", s.deleteSubscriptionID)
	CommandGroup.Command("add <subtype> [tags...]", s.subscribeCommand).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator), middleware.Coooldown(7*time.Second, 3))

	return CommandGroup
}
