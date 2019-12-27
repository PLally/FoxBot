package chat_tools

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/discord_modular_bot/command"
	"strconv"
	"time"
)

type Snowflake struct {
	ID                uint64
	Increment         uint64
	InternalProcessID uint64
	InternalWorkerID  uint64
	TimestampDiscord  uint64
	TimestampUnix     uint64
	Time              time.Time
}

func getDiscordObjectInfo(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	//TODO support channels, snowflakes, voice channels, emojis
	var user *discordgo.User
	if len(event.Message.Mentions) < 1 {
		return
	}
	user = event.Message.Mentions[0]

	embed := command.NewEmbed()
	embed.SetThumbnailUrl(user.AvatarURL("1024"))
	snowflake, _ := GetSnowflake(user.ID)
	snowflakeString := fmt.Sprintf(
		"```ID: %v\nIncrement: %d\nInternalProcessID: %d\nInternalWorkderID: %d\nTimestamp: %v```",
		snowflake.ID,
		snowflake.Increment,
		snowflake.InternalProcessID,
		snowflake.InternalWorkerID,
		snowflake.Time.UTC().String(),
	)
	embed.SetTitle(user.String(), "")
	embed.AddField("Snowflake Info", snowflakeString, true)
	s.ChannelMessageSendEmbed(event.Message.ChannelID, embed.MessageEmbed)

	return ""

}
func GetSnowflake(id string) (Snowflake, error) {
	snowflake, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return Snowflake{}, err
	}

	increment := snowflake & (2 ^ 12 - 1)
	internalProcessID := (snowflake & 0x1F000) >> 12
	internalWorkerID := (snowflake & 0x3E0000) >> 17
	timestampDiscord := snowflake >> 22
	timestampUnix := 1420070400000 + timestampDiscord
	creationTime := time.Unix(int64(timestampUnix/1000), 0)

	return Snowflake{
		snowflake,
		increment,
		internalProcessID,
		internalWorkerID,
		timestampDiscord,
		timestampUnix,
		creationTime,
	}, nil
}
