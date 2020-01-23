package statistics

import (
	"github.com/bwmarrin/discordgo"
	influx "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// statistics collection

func onBotReady(s *discordgo.Session, ready *discordgo.Ready) {
	c, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     os.Getenv("INFLUXDB_ADDRESS"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
		Username: os.Getenv("INFLUXDB_USERNAME"),
	})

	go func() {
		client.Client = c
		client.messageBatch, err = influx.NewBatchPoints(client.messageConfig)

		if err != nil {
			log.Error("InfluxDB client failed: " + err.Error())
			return
		}
		s.AddHandler(onMessageCreate)
		defer client.Close()
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
		"guild_id":   m.GuildID,
		"channel_id": m.ChannelID,
	}
	fields := map[string]interface{}{
		"quantity": 1,
		"user_id":  m.Author.ID,
	}
	point, err := influx.NewPoint("discord_messages", tags, fields, time.Now().UTC())
	if err != nil {
		log.Error(err)
		return
	}
	client.messageBatch.AddPoint(point)
}
