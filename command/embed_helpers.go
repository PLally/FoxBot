package command

import "github.com/bwmarrin/discordgo"

type Embed struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

func (e *Embed) AddField(name, value string, inline bool) {
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

func (e *Embed) SetTitle(title string) {
	e.Title = title
}

func (e *Embed) SetImageUrl(url string) {
	e.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
}

func (e *Embed) SetThumbnailUrl(url string) {
	e.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
}
