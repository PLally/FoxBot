package subtypes

import (
	"fmt"
	"github.com/plally/e621"
	"github.com/plally/subscription_api/subscription"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func RegisterE621() {
	handler := &E621Handler{
		Session: &e621.Session{
			BaseURL:   "https://e621.net",
			UserAgent: "FoxBotSubscriptions/0.1",
			Client:    &http.Client{},
			Username:  os.Getenv("E621_USERNAME"),
			ApiKey:    os.Getenv("E621_TOKEN"),
		},
	}
	handler.updatePostCache()
	go func() {
		time.Sleep(time.Minute * 15)
		handler.updatePostCache()

	}()
	subscription.SetSubTypeHandler("e621", handler)

}

type E621Handler struct {
	Session     *e621.Session
	LastUpdated int
	postCache   []*e621.Post
}

func (r *E621Handler) updatePostCache() {
	resp, _ := r.Session.GetPosts("", 1)
	lastId := resp.Posts[len(resp.Posts)-1].ID
	posts := resp.Posts
	for i:=1; i<100; i++ {
		resp, err := r.Session.GetPosts("id:<"+strconv.Itoa(lastId), 320)
		if err  != nil{
			log.Error(err)
			continue
		}

		posts = append(posts, resp.Posts...)
		lastId = resp.Posts[len(resp.Posts)-1].ID
		fmt.Println(len(posts), lastId)
		time.Sleep(time.Millisecond*500)
	}
	fmt.Println("UPDATED")
	r.postCache = posts
}

func (r *E621Handler) GetType() string { return "e621" }

func (r *E621Handler) GetNewItems(tags string) []subscription.SubscriptionItem {
	var items []subscription.SubscriptionItem
	parsed, _ := e621.ParseTags(tags, false)
	for _, post := range r.postCache {
		if !parsed.Matches(post.Tags) {
			continue
		}

		sub_item := subscription.SubscriptionItem{
			Title:       fmt.Sprintf( "E621 Post #%v", post.ID),
			Url:         r.Session.PostUrl(post),
			Description: fmt.Sprintf("Artists %v", strings.Join(post.Tags.Artist, ". ")), // todo use post.Description and truncate this
			Author:      strings.Join(post.Tags.Artist,", "),
			TimeID:      int64(post.ID),
			Image:       post.File.URL,
		}
		fmt.Println(sub_item.Url)
		items = append(items, sub_item)
	}
	return items
}

func (r *E621Handler) Validate(tags string) (string, error) {
	parsed, err := e621.ParseTags(tags, false)
	if err != nil {
		return "", err
	}
	return parsed.Normalized(), nil
}
