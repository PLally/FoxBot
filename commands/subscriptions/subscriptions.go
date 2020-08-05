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

func (s subClient) subscribeCommand(ctx dgcommand.CommandContext) {
	subType := ctx.Args()[0]
	tags := ctx.Args()[1]

	sub, err := s.Subscribe("discord", ctx.Message.ChannelID, subType, tags)

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

func (s subClient) listSubscriptions(ctx dgcommand.CommandContext) {
	subs, err := s.FindChannelSubscriptions(ctx.Message.ChannelID)

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

	if len(subs) == 0 {
		msg = msg + "No subscriptions in this channel."
	}

	for _, sub := range subs {
		msg = msg + fmt.Sprintf("[%v]: %v - %v\n", sub.ID, sub.SubscriptionType.Type, sub.SubscriptionType.Tags)
	}

	msg += "```"

	ctx.Reply(msg)
}

func (s subClient) deleteSubscriptionID(ctx dgcommand.CommandContext) {
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

func (s subClient) deleteSubscription(ctx dgcommand.CommandContext) {

	subTypeType := ctx.Args()[0]
	subTypeTags := ctx.Args()[1]

	subs, err := s.FindChannelSubscriptions(ctx.Message.ChannelID)

	if err != nil {
		switch err.(type) {
		case subscription_client.SubError:
			ctx.Reply(err.Error())
		default:
			ctx.Error(err)
		}

		return
	}

	for _, sub := range subs {
		if sub.SubscriptionType.Type == subTypeType && sub.SubscriptionType.Tags == subTypeTags {
			_, err := s.DeleteSubscription(int(sub.ID))
			if err != nil {
				ctx.Error(err)
			}

			subPrintString := fmt.Sprintf("[%v]: %v - %v\n", sub.ID, sub.SubscriptionType.Type, sub.SubscriptionType.Tags)
			ctx.Reply("Deleted susbcription " + subPrintString)
			return
		}
	}

	ctx.Reply("Couldnt find that subscription")
}
