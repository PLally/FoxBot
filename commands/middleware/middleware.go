package middleware

import (
	"github.com/plally/dgcommand"
	log "github.com/sirupsen/logrus"
)

func RequireNSFW() dgcommand.MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		return func(ctx dgcommand.CommandContext) {
			channel, err := ctx.Session.State.Channel(ctx.Message.ChannelID)
			if err != nil || !channel.NSFW {
				ctx.Reply("You must be in an nsfw channel to do this")
				return
			}
			h(ctx)
		}
	}
}

func RequirePermissions(perms ...int) dgcommand.MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		requiredPerms := perms
		return func(ctx dgcommand.CommandContext) {
			authorPerms, err := ctx.Session.UserChannelPermissions(ctx.Message.Author.ID, ctx.Message.ChannelID)
			if err != nil {
				log.Error(err)
				ctx.Reply("You dont have the required permissions to do this")
				return
			}
			for _, perm := range requiredPerms {
				if !(authorPerms & perm == perm) {
					ctx.Reply("You dont have the required permissions to do this")
					return
				}
			}
			h(ctx)
		}
	}
}
