package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
	"github.com/plally/subscription_api/subscription"
)

type subClient struct {
	subscription_client.SubscriptionClient
}

func (s subClient) subscribeCommand(ctx dgcommand.CommandContext) {
	subType := ctx.Args[0]
	tags    := ctx.Args[1]

	sub, err := s.Subscribe(subType, tags, ctx.M.ChannelID)
	if err != nil { ctx.Error(err); return }

	if sub.ID <= 0 {
		ctx.Reply("There was a problem subscribing to that")
	}

	ctx.Reply("Subscribed")
}

func (s subClient) listSubscriptions(ctx dgcommand.CommandContext) {
	subs, err := s.GetSubscriptions(ctx.M.ChannelID)
	if err != nil { ctx.Error(err); return }

	msg := "```"
	for _, sub := range subs {
		msg = msg + fmt.Sprintf("[%v]: %v - %v\n", sub.ID, sub.SubscriptionType.Type, sub.SubscriptionType.Tags)
	}

	msg += "```"

	ctx.Reply(msg)
}

func (s subClient) deleteSusbcription(ctx dgcommand.CommandContext) {
	err := s.DeleteSubscription(ctx.Args[0], ctx.Args[1], ctx.M.ChannelID)
	if err != nil { ctx.Error(err); return }

	ctx.Reply("Subscription deleted")
}

func RegisterSubCommands(r *dgcommand.CommandRoutingHandler, db *gorm.DB) {

	s := subClient{subscription_client.SubscriptionClient{DB:db}}

	deleteCallback := middleware.Wrap(
		s.deleteSusbcription,
		middleware.RequirePermissions(discordgo.PermissionAdministrator),
	)
	subscribeCallback := middleware.Wrap(
		s.subscribeCommand,
		middleware.RequirePermissions(discordgo.PermissionAdministrator),
	)
	r.AddHandler("delete", dgcommand.NewCommand("delete <subtype> [tags...]", deleteCallback))

	r.AddHandler("list", dgcommand.NewCommand("list", s.listSubscriptions))

	r.AddHandler("add", dgcommand.NewCommand("add <subtype> [tags...]", subscribeCallback))

	go subscription.CheckOutDatedSubscriptionTypes(db, 100)
}