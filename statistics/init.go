package statistics

import (
	"github.com/bwmarrin/discordgo"
	influx "github.com/influxdata/influxdb1-client/v2"
)

var client = struct {
	influx.Client
	messageConfig influx.BatchPointsConfig
	messageBatch  influx.BatchPoints
}{
	messageConfig: influx.BatchPointsConfig{
		Database: "discord_bot_stats",
	},
}

func initStats(session *discordgo.Session){
	session.AddHandler(onBotReady)
}
