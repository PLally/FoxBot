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

func randomFox(ctx dgcommand.Context) {
	resp, err := http.Get("https://api.foxorsomething.net/fox/random.png")
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.SendFile("fox.png", resp.Body)
}

func randomCat(ctx dgcommand.Context) {
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

func randomUser(genericContext dgcommand.Context) {
	ctx := genericContext.(*dgcommand.DiscordContext)

	guild, err := ctx.S.State.Guild(ctx.M.GuildID)
	if err != nil {
		ctx.Error(err)
		return
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	user := guild.Members[random.Intn(guild.MemberCount)].User
	ctx.S.ChannelMessageSend(ctx.M.ChannelID, user.Username+"#"+user.Discriminator)
}