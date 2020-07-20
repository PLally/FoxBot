package random

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
		ctx.Error(err)
		return
	}
	ctx.SendFile("fox.png", resp.Body)
}

func randomCat(ctx dgcommand.CommandContext) {
	resp, err := http.Get("http://aws.random.cat/meow")
	if err != nil {
		log.Error(err)
		ctx.Reply("something went wrong fetching the cat")
	}
	var dat map[string]string
	bytes, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(bytes, &dat)
	e := embed.NewEmbed()
	e.SetImageUrl(dat["file"])
	ctx.SendEmbed(e)

}

func randomUser(ctx dgcommand.CommandContext) {

	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		ctx.Error(err)
		return
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	user := guild.Members[random.Intn(guild.MemberCount)].User
	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, user.Username+"#"+user.Discriminator)
}