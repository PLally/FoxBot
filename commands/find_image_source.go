package commands

import (
	"encoding/json"
	"fmt"
	"github.com/plally/dgcommand"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func findSource(imgUrl string) []string{
	api_key := viper.GetString("sauce_api_key")
	url := fmt.Sprintf("https://saucenao.com/search.php?db=999&output_type=2&numres=4&url=%v&api_key=%v", imgUrl, api_key)
	resp, err := http.Get(url)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	var obj SauceResponse
	json.Unmarshal(data, &obj)

	urls := make([]string, 1)
	for _, result := range obj.Results {
		similarity, err := strconv.ParseFloat(result.Header.Similarity, 64)
		if err != nil {
			continue
		}

		if similarity > 80 && len(result.Data.ExtUrls) > 0 {
			urls = append(urls, result.Data.ExtUrls[0])
		}
	}
	if len(obj.Results) < 1 {
		return []string{}
	}
	return urls
}

func findSourceCommand(ctx dgcommand.CommandContext) {
	reference := ctx.Message.MessageReference
	msg, err := ctx.Session.ChannelMessage(reference.ChannelID, reference.MessageID)
	if err != nil {
		ctx.Error(err)
	}

	urls := findSource(msg.Content)
	if len(urls) < 1 {
		ctx.Reply("Not Found")
		return
	}
	ctx.Reply(strings.Join(urls, "\n"))
}

type SauceResponse struct {
	Header struct {
		UserID           int    `json:"user_id"`
		AccountType      int    `json:"account_type"`
		ShortLimit       string `json:"short_limit"`
		LongLimit        string `json:"long_limit"`
		LongRemaining    int    `json:"long_remaining"`
		ShortRemaining   int    `json:"short_remaining"`
		Status           int    `json:"status"`
		ResultsRequested int    `json:"results_requested"`
		Index map[string]struct{
			Status   int `json:"status"`
			ParentID int `json:"parent_id"`
			ID       int `json:"id"`
			Results  int `json:"results"`
		} `json:"index"`
		SearchDepth       string  `json:"search_depth"`
		MinimumSimilarity float64 `json:"minimum_similarity"`
		QueryImageDisplay string  `json:"query_image_display"`
		QueryImage        string  `json:"query_image"`
		ResultsReturned   int     `json:"results_returned"`
	} `json:"header"`
	Results []struct {
		Header struct {
			Similarity string `json:"similarity"`
			Thumbnail  string `json:"thumbnail"`
			IndexID    int    `json:"index_id"`
			IndexName  string `json:"index_name"`
			Dupes      int    `json:"dupes"`
		} `json:"header"`
		Data struct {
			ExtUrls    []string `json:"ext_urls"`
			Title      string   `json:"title"`
		} `json:"data,omitempty"`
	} `json:"results"`
}