package commands

import (
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/plally/dgcommand/embed"
	"github.com/plally/e621"
	"path"
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

var e6Session = e621.NewSession("e621.net", "FoxBot/0.1")

func e621Func(ctx dgcommand.CommandContext) {
	resp, err := e6Session.GetPosts("order:random "+ctx.Args()[0], 1)
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
	description := strings.Builder{}

	for i, artist := range post.Tags.Artist {
		if i != 0 {
			description.WriteString(", ")
		}
		artistString := fmt.Sprintf("[%[1]v](https://e621.net/post?tags=%[1]v)", artist)
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
	if contentUrl != post.File.URL {
		description.WriteString(fmt.Sprintf("\n\n*Click **E621 Post %v** to view content in its original form*", post.ID))
	}

	e := embed.NewEmbed()
	e.SetTitle(fmt.Sprintf("E621 Post %v", post.ID), e6Session.PostUrl(post))
	e.SetImageUrl(contentUrl)
	e.Description = description.String()
	ctx.SendEmbed(e)
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
	}

	for _, url := range urls {
		fileExt := path.Ext(url)
		if validSuffixes[fileExt] {
			return url
		}
	}
	return p.File.URL
}
