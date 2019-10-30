package statistics

import (
	"github.com/bwmarrin/discordgo"
	influx "github.com/influxdata/influxdb1-client/v2"
	"github.com/plally/modular_bot/command"
	log "github.com/sirupsen/logrus"
	"time"
)

var client = struct {
	influx.Client
	messageConfig influx.BatchPointsConfig
	messageBatch  influx.BatchPoints
}{
	messageConfig: influx.BatchPointsConfig{
		Database: "bot_test_db",
	},
}

func init() {
	command.RegisterModule("stats").OnEnable = OnEnable
}

func OnEnable(b *command.Bot) {
	b.AddHandler(onBotReady)
}
func onBotReady(s *discordgo.Session, ready *discordgo.Ready) {
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: "http://127.0.0.1:8086",
	})

	go func() {
		client.Client = c
		client.messageBatch, err = influx.NewBatchPoints(client.messageConfig)

		if err != nil {
			log.Error(err)
			return
		}
		s.AddHandler(onMessageCreate)
		defer client.Close()
		if err != nil {
			log.Error("InfluxDB client failed: " + err.Error())
			return
		}
		for {
			err := client.Client.Write(client.messageBatch)
			if err != nil {
				log.Error(err.Error())
			}
			client.messageBatch, err = influx.NewBatchPoints(client.messageConfig)
			if err != nil {
				log.Error(err)
				continue
			}

			time.Sleep(time.Second * 5)
		}
	}()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	tags := map[string]string{
		"user_id":    m.Author.ID,
		"guild_id":   m.GuildID,
		"channel_id": m.ChannelID,
	}
	fields := map[string]interface{}{
		"quantity": 1,
	}
	point, err := influx.NewPoint("discord_messages", tags, fields, time.Now().UTC())
	if err != nil {
		log.Error(err)
		return
	}
	client.messageBatch.AddPoint(point)
}
