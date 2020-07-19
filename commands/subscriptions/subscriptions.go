package subscriptions

import (
	"fmt"
	"github.com/plally/FoxBot/subscription_client"
	"github.com/plally/dgcommand"
	"strconv"
)

type subClient struct {
	*subscription_client.SubscriptionClient
}

func (s subClient) subscribeCommand(ctx dgcommand.Context) {
	subType := ctx.Args()[0]
	tags := ctx.Args()[1]

	sub, err := s.Subscribe("discord", ctx.Message().ChannelID, subType, tags)

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
	subs, err := s.FindChannelSubscriptions(ctx.Message().ChannelID)

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

func (s subClient) deleteSubscriptionID(ctx dgcommand.Context) {
	id, _ := strconv.Atoi(ctx.Args()[0])
	sub, err := s.DeleteSubscription(id)

	if err != nil {
		switch err.(type) {
		case subscription_client.SubError:
			ctx.Reply(err.Error())
		default:
			ctx.Error(err)
		}

		return
	}
	ctx.Reply(fmt.Sprintf("deleted subscription %v", sub.ID))

}
