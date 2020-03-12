package middleware

import (
	"github.com/plally/dgcommand"
	"github.com/sirupsen/logrus"
	"strings"
)

type MiddlewareFunc func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc
func RequireNSFW() MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		return func(ctx dgcommand.CommandContext) {
			channel, err := ctx.S.State.Channel(ctx.M.ChannelID)
			if err != nil || channel.NSFW {
				ctx.Reply("You must be in an nsfw channel to do this")
				return
			}
			h(ctx)
		}
	}
}
func LogWith(l *logrus.Logger) MiddlewareFunc{
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		return func(ctx dgcommand.CommandContext) {
			l.Infof("Handling args: %v", strings.Join(ctx.Args, ", "))
			h(ctx)
		}
	}

}

func RequirePermissions(perms ...int) MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		requiredPerms := perms
		return func(ctx dgcommand.CommandContext) {
			authorPerms, err := ctx.S.State.UserChannelPermissions(ctx.M.Author.ID, ctx.M.ChannelID)
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

func Wrap(f dgcommand.HandlerFunc, middleware ...MiddlewareFunc) dgcommand.HandlerFunc{
	fn := f
	for _, mid := range middleware {
		fn = mid(f)
	}

	return fn
}
