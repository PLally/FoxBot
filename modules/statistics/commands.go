package statistics

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/plally/discord_modular_bot/command"
	log "github.com/sirupsen/logrus"
)
func userStatsCommand(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	for _, user := range event.Message.Mentions {
		return userStats(s, user, event.Message.GuildID)
	}
	return
}

func userStats(s *discordgo.Session, user *discordgo.User, guildID string) (reply string) {
	queryString := `select sum("quantity") from discord_messages where "guild_id" = '%s' AND "user_id" = '%s'`
	queryString = fmt.Sprintf(queryString,guildID, user.ID)
	q := influx.NewQuery(queryString, "discord_bot_stats", "")
	response, err := client.Query(q); if err != nil || response.Error() != nil {

		return response.Error().Error()
	}
	fmt.Println(queryString)

	result := response.Results[0]
	if len(result.Series) < 1 { return "No Stats"}
	values := result.Series[0].Values

	return fmt.Sprintf("```%v messages```", values[0][1])
}

func guildStatsCommand(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
	guildID := event.Message.GuildID
	queryString := `select sum(quantity) from discord_messages where guild_id='%s' group by channel_id`
	queryString = fmt.Sprintf(queryString, guildID)
	q := influx.NewQuery(queryString, "discord_bot_stats", "")
	response, err := client.Query(q); if err != nil || response.Error() != nil {
		return response.Error().Error()
	}
	fmt.Println(queryString)
	result := response.Results[0]
	if len(result.Series) < 1 { return "No Stats"}
	channelActivityText := ""
	for _, series := range result.Series {
		channel, err := s.State.Channel(series.Tags["channel_id"])
		if err != nil {
			log.Error(err)
			continue
		}
		channelActivityText += fmt.Sprintf("%s: %s\n", channel.Name, series.Values[0][1])
	}
	embed := command.NewEmbed()
	embed.SetTitle("Guild Statistics", "")
	embed.AddField("Channel Activity", channelActivityText, false)
	s.ChannelMessageSendEmbed(event.Message.ChannelID, embed.MessageEmbed)
	return ""
}
