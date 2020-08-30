package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"github.com/plally/e621"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const (
	E621_NOSCORE                = "<:e6_noscore:739285090206613505>"
	E621_UP                     = "<:e6_upvote_new:739344985761251359>"
	E621_DOWN                   = "<:e6_downvote_new:739344985580765305>"
	E621_RATING_SAFE            = "<:e6_rating_s:739330560471728179>"
	E621_RATING_QUESTIONABLE    = "<:e6_rating_q:739330560664535040>"
	E621_RATING_EXPLICIT        = "<:e6_rating_e:739330560719060992>"
	E621_FAV_COUNT              = "<:e6_fav_count:739349545917481044>"
	E621_COMMENT_COUNT          = "<:e6_comment_count:739353756432597063>"
	SPACING_EMOTE               = "<:nothing:739347008979992636>"
	E621_MAX_DESCRIPTION_LENGTH = 300
)

var e621session = e621.NewSession("e621.net", "FoxBot/0.1")
var e926session = e621.NewSession("e926.net", "FoxBot/0.1")
func e926Func(ctx dgcommand.CommandContext) {
	e6command(e926session, ctx)
}

func e621Func(ctx dgcommand.CommandContext) {
	e6command(e621session, ctx)
}

func e6command(session *e621.Session, ctx dgcommand.CommandContext) {
	resp, err := session.GetPosts("order:random "+ctx.Args()[0], 1)
	if err != nil {
		ctx.Error(err)
		return
	}
	posts := resp.Posts
	if len(posts) < 1 {
		ctx.Reply("No posts were found with those tags")
		return
	}
	post := posts[0]

	contentUrl := GetValidContentURL(post)
	description := buildDescription(session.BaseURL, post)

	if path.Ext(contentUrl) == ".webm" {

		ctx.Reply(getEmbedLink(description, session, post))
		return
	}

	if contentUrl != post.File.URL {
		description = description + fmt.Sprintf("\n\n*Click **Post %v** to view content in its original form*", post.ID)
	}

	e := embed.NewEmbed()
	e.SetTitle(fmt.Sprintf("Post %v", post.ID), session.PostUrl(post))
	e.SetImageUrl(contentUrl)
	e.Description = description
	ctx.SendEmbed(e)
}

func buildDescription(baseUrl string, post *e621.Post) string {
	description := strings.Builder{}

	for i, artist := range post.Tags.Artist {
		if i != 0 {
			description.WriteString(", ")
		}
		artistString := fmt.Sprintf("[%[1]v](%v/post?tags=%[1]v)", artist, baseUrl)
		description.WriteString(artistString)
	}
	description.WriteByte('\n')
	var arrow string
	if post.Score.Total > 0 {
		arrow = E621_UP
	} else if post.Score.Total < 0 {
		arrow = E621_DOWN
	} else {
		arrow = E621_NOSCORE
	}

	rating := ""
	switch post.Rating {
	case "e":
		rating = E621_RATING_EXPLICIT
	case "q":
		rating = E621_RATING_QUESTIONABLE
	case "s":
		rating = E621_RATING_SAFE
	}

	infoLine := fmt.Sprintf("%v%v  %v%v  %v** %v**  %v\n\n",
		arrow, post.Score.Total, E621_FAV_COUNT, post.FavCount, E621_COMMENT_COUNT, post.CommentCount, rating)
	infoLine = strings.ReplaceAll(infoLine, "  ", SPACING_EMOTE)
	description.WriteString(infoLine)

	if len(post.Description) > E621_MAX_DESCRIPTION_LENGTH {
		post.Description = post.Description[:E621_MAX_DESCRIPTION_LENGTH] + "..."
	}

	description.WriteString(post.Description)
	return description.String()
}


type newEmbedPage struct {
	Meta map[string]string
	Title string
	Color string
	Redirect string
	Name string
}

func getEmbedLink(description string, session *e621.Session, post *e621.Post) string {
	data, _ := json.Marshal(newEmbedPage{
		Title: fmt.Sprintf("Post %v", post.ID),
		Redirect: session.PostUrl(post),
		Name: "e621/" + strconv.Itoa(post.ID),
		Meta: map[string]string{
			"og:description": description,
			"og:video": post.File.URL,
			"og:type": "video",
			"og:title": fmt.Sprintf("Post %v", post.ID),
			"theme-color": "#0000FF",
		},
	})

	bodyReader := bytes.NewReader(data)
	req, _ := http.NewRequest("POST", "https://embed.foxorsomething.net/newpage", bodyReader)
	req.Header.Set("Authorization", viper.GetString("embed_maker_auth"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return session.PostUrl(post)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("%v returned from embed maker", resp.Status)
		return session.PostUrl(post)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return "https://embed.foxorsomething.net/"+string(body)
}

func GetValidContentURL(p *e621.Post) string {
	urls := []string{
		p.File.URL,
		p.Sample.URL,
		p.Preview.URL,
	}

	validSuffixes := map[string]bool{
		".gif": true,
		".jpg": true,
		".png": true,
		".webm": true,
	}

	for _, url := range urls {
		fileExt := path.Ext(url)
		if validSuffixes[fileExt] {
			return url
		}
	}
	return p.File.URL
}
