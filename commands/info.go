package commands

import (
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"github.com/plally/dgcommand/snowflake"
	"regexp"
	"strings"
)

var (
	mentionPattern   = regexp.MustCompile("^ *<@![0-9]+>")
	snowflakePattern = regexp.MustCompile("^ *[0-9]+")
)

func objInfoFunc(ctx dgcommand.Context) {
	switch ctx := ctx.(type) {
	case *dgcommand.DiscordContext:
		obj := ctx.Args()[0]
		if mentionPattern.MatchString(obj) {
			mentionInfo(ctx, obj)
		} else if snowflakePattern.MatchString(obj) {
			snowflakeInfo(ctx, obj)
		} else {
			ctx.Reply("Couldn't give you any info about that object")
		}
	default:
		ctx.Error("Command info only accepts a discord context")
	}
}

func mentionInfo(ctx *dgcommand.DiscordContext, obj string) {
	objID := strings.ReplaceAll(obj, "!", "")
	objID = strings.ReplaceAll(objID, "<", "")
	objID = strings.ReplaceAll(objID, ">", "")
	objID = strings.ReplaceAll(objID, "@", "")

	for _, user := range ctx.M.Mentions {
		if user.ID == objID {
			e := embed.NewEmbed()
			e.SetThumbnailUrl(user.AvatarURL("1024"))
			snow, err := snowflake.NewSnowflake(user.ID)
			if err != nil {
				break
			}
			e.SetTitle(user.String(), "")
			e.AddField("Snowflake Info", getSnowflakeString(snow), true)
			ctx.S.ChannelMessageSendEmbed(ctx.M.ChannelID, e.MessageEmbed)
			return
		}
	}
	ctx.Reply("Couldn't give you any info about that mention")
}

func snowflakeInfo(ctx *dgcommand.DiscordContext, obj string) {
	user, err := ctx.S.User(obj)

	e := embed.NewEmbed()
	e.SetTitle("Snowflake info", "")
	if user != nil && err == nil {
		e.SetThumbnailUrl(user.AvatarURL("1024"))
		e.SetTitle(user.String(), "")
	}
	snow, err := snowflake.NewSnowflake(obj)
	if err != nil {
		ctx.Reply("Error processing snowflake")
		return
	}
	e.AddField("Snowflake Info", getSnowflakeString(snow), true)
	ctx.SendEmbed(e)
}

func getSnowflakeString(snow snowflake.Snowflake) string {
	return fmt.Sprintf(
		"```ID: %v\nIncrement: %d\nInternalProcessID: %d\nInternalWorkderID: %d\nTimestamp: %v```",
		snow.ID,
		snow.Increment,
		snow.InternalProcessID,
		snow.InternalWorkerID,
		snow.Time.UTC().String(),
	)
}