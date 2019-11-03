package statistics

import (
	"github.com/bwmarrin/discordgo"
	"github.com/plally/discord_modular_bot/command"
	influx "github.com/influxdata/influxdb1-client/v2"
	"fmt"
)
func stats(s *discordgo.Session, event *command.TextCommandEvent) (reply string) {
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