package middleware

import (
	"github.com/plally/dgcommand"

)

func RequireNSFW() dgcommand.MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		return func(ctx dgcommand.Context) {
			switch ctx := (ctx).(type) {
			case *dgcommand.DiscordContext:
				channel, err := ctx.S.State.Channel(ctx.M.ChannelID)
				if err != nil || !channel.NSFW {
					ctx.Reply("You must be in an nsfw channel to do this")
					return
				}
				h(ctx)
			}

			ctx.Reply("you cant do this")
		}
	}
}


func RequirePermissions(perms ...int) dgcommand.MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		requiredPerms := perms
		return func(ctx dgcommand.Context) {
			authorPerms, err := ctx.(*dgcommand.DiscordContext).S.State.UserChannelPermissions(ctx.Message().ID, ctx.Message().ChannelID)
			if err != nil {
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
