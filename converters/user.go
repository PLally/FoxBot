package converters

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/dgcommand"
	"regexp"
	"strings"
)

var mentionPattern = regexp.MustCompile("(:?^<@!?)?([0-9]+)(:?>)?")

func ParseMember(ctx dgcommand.CommandContext, arg string) (*discordgo.Member) {
	arg = strings.Trim(arg, " ")

	match := mentionPattern.FindStringSubmatch(arg)

	if len(match) >= 2 {
		userID := match[2]
		member, _ := ctx.Session.GuildMember(ctx.Message.GuildID, userID)
		return member
	}

	guild, err := ctx.Session.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	for _, member := range guild.Members {
		if member.User.Username + member.User.Discriminator == arg {
			return member
		}
	}
	return nil
}
