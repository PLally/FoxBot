package subscriptions

import (
	"fmt"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
)

type subClient struct {
	*subscription_client.SubscriptionClient
}

func (s subClient) subscribeCommand(ctx dgcommand.Context) {
	subType := ctx.Args()[0]
	tags := ctx.Args()[1]

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

	ctx.Reply(fmt.Sprintf("Created subscription %v: %v", sub.SubscriptionType.Type, sub.SubscriptionType.Tags))
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