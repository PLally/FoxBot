package desttypes

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/dgcommand/embed"
	"github.com/plally/subscription_api/subscription"
	"github.com/sirupsen/logrus"
)

func RegisterDiscord(session *discordgo.Session) {
	subscription.SetDestinationHandler("discord", &DiscordDestinationHandler{session})
}

type DiscordDestinationHandler struct {
	session *discordgo.Session
}

func (d *DiscordDestinationHandler) GetType() string {
	return "discord"
}

func (d *DiscordDestinationHandler) Dispatch(id string, item subscription.SubscriptionItem) error {
	e := embed.NewEmbed()
	e.SetTitle(item.Title, item.Url)
	e.SetDescription(item.Description)
	e.SetImageUrl(item.Image)
	e.SetFooter(fmt.Sprintf("subscription: %v - %v", item.Type, item.Tags), "", "")
	_, err := d.session.ChannelMessageSendEmbed(id, e.MessageEmbed)
	if err != nil {
		logrus.Error(err)
	}
	return nil
}
