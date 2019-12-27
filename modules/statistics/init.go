package statistics

import (
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/plally/discord_modular_bot/command"
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
var module *command.Module
func init() {
	module = command.RegisterModule("stats")
	module.RegisterCommandFunc(">user_stats", userStatsCommand)
	module.RegisterCommandFunc(">stats", guildStatsCommand)
	module.OnEnable = OnEnable
}

func OnEnable(b *command.Bot) {
	b.AddHandler(onBotReady)
}
