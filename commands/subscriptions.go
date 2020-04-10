package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"github.com/plally/FoxBot/commands/middleware"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
	"github.com/plally/subscription_api/subscription"
	"time"
)

type subClient struct {
	subscription_client.SubscriptionClient
}

func (s subClient) subscribeCommand(ctx dgcommand.Context) {
	subType := ctx.Args()[0]
	tags    := ctx.Args()[1]

	sub, err := s.Subscribe(subType, tags, ctx.Message().ChannelID)

	if err != nil {
		switch err.(type) {
		case subscription_client.SubError:
			ctx.Reply(err.Error())
		default:
			ctx.Error(err)
		}

		return
	}

	if sub.ID <= 0 {
		ctx.Reply("There was a problem subscribing to that")
	}

	ctx.Reply("Subscribed")
}

func (s subClient) listSubscriptions(ctx dgcommand.Context) {
	subs, err := s.GetSubscriptions(ctx.Message().ChannelID)

	if err != nil {
		switch err.(type) {
		case subscription_client.SubError:
			ctx.Reply(err.Error())
		default:
			ctx.Error(err)
		}

		return
	}

	msg := "```"
	for _, sub := range subs {
		msg = msg + fmt.Sprintf("[%v]: %v - %v\n", sub.ID, sub.SubscriptionType.Type, sub.SubscriptionType.Tags)
	}

	msg += "```"

	ctx.Reply(msg)
}

func (s subClient) deleteSusbcription(ctx dgcommand.Context) {
	message := ctx.Message()
	err := s.DeleteSubscription(ctx.Args()[0], ctx.Args()[1], message.ChannelID)

	if err != nil {
		switch err.(type) {
		case subscription_client.SubError:
			ctx.Reply(err.Error())
		default:
			ctx.Error(err)
		}

		return
	}
	ctx.Reply("Subscription deleted")
}

func RegisterSubCommands(r *dgcommand.CommandGroup, db *gorm.DB) {
	s := subClient{subscription_client.SubscriptionClient{DB:db}}
	r.AddHandler("list", dgcommand.NewCommand("list", s.listSubscriptions))

	r.Command("delete <subtype> [tags...]", s.deleteSusbcription).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator))

	r.Command("add <subtype> [tags...]", s.subscribeCommand).
		Use(middleware.RequirePermissions(discordgo.PermissionAdministrator))

	go func() {
		for {
			subscription.CheckOutDatedSubscriptionTypes(db, 100)
			time.Sleep(time.Minute * 15)
		}
	}()
}