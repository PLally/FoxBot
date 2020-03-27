package commands

import (
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func randomFox(ctx dgcommand.CommandContext) {
	resp, err := http.Get("https://api.foxorsomething.net/fox/random.png")
	if err != nil {
		log.Error(err)
		ctx.Reply("something went wrong fetching the fox")
		return
	}

	ctx.S.ChannelFileSend(ctx.M.ChannelID, "fox.png", resp.Body)
}

var RandomFoxCommand = dgcommand.NewCommand("fox", randomFox)

func randomCat(ctx dgcommand.CommandContext) {
	resp, err := http.Get("http://aws.random.cat/meow")
	if err != nil {
		log.Error(err)
		ctx.S.ChannelMessageSend(ctx.M.ChannelID, "something went wrong fetching the cat")
	}
	var dat map[string]string
	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &dat)
	e := embed.NewEmbed()
	e.SetImageUrl(dat["file"])
	ctx.S.ChannelMessageSendEmbed(ctx.M.ChannelID, e.MessageEmbed)

}

var RandomCatCommand = dgcommand.NewCommand("cat", randomCat)

func randomUser(ctx dgcommand.CommandContext) {
	guild, err := ctx.S.State.Guild(ctx.M.GuildID)
	if err != nil {
		ctx.S.ChannelMessageSend(ctx.M.ChannelID, "something went wrong")
	}
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	user := guild.Members[random.Intn(guild.MemberCount)].User
	ctx.S.ChannelMessageSend(ctx.M.ChannelID, user.Username+"#"+user.Discriminator)
}

var RandomUserCommand = dgcommand.NewCommand("user", randomUser)
