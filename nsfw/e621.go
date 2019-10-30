package nsfw

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type JsonTime struct {
	JsonClass string
	S         int
	N         int
}

type E621Post struct {
	ID            int      `json:"id"`
	Tags          string   `json:"tags"`
	LockedTags    string   `json:"locked_tags"`
	Description   string   `json:"description"`
	CreatedAt     JsonTime `json:"created_at"`
	CreatorID     int      `json:"creator_id"`
	Author        string   `json:"author"`
	Change        int      `json:"change"`
	Source        string   `json:"source"`
	Score         int      `json:"score"`
	FavCount      int      `json:"fav_count"`
	Md5           string   `json:"md5"`
	FileSize      int      `json:"file_size"`
	FileURL       string   `json:"file_url"`
	FileExt       string   `json:"file_ext"`
	PreviewURL    string   `json:"preview_url"`
	PreviewWidth  int      `json:"preview_width"`
	PreviewHeight int      `json:"preview_height"`
	SampleURL     string   `json:"sample_url"`
	SampleWidth   int      `json:"sample_width"`
	SampleHeight  int      `json:"sample_height"`
	Rating        string   `json:"rating"`
	Status        string   `json:"status"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	HasComments   bool     `json:"has_comments"`
	HasNotes      bool     `json:"has_notes"`
	HasChildren   bool     `json:"has_children"`
	Children      string   `json:"children"`
	ParentID      int      `json:"parent_id"`
	Artist        []string `json:"artist"`
	Sources       []string `json:"sources"`
}

type E621Session struct {
	BaseURL   string
	Username  string
	UserAgent string
	ApiKey    string
	Client    *http.Client
}

func (e *E621Session) get(url string, params map[string]string) (*http.Response, error) {
	params["password_hash"] = e.ApiKey
	params["login"] = e.Username

	req, err := http.NewRequest("GET", e.BaseURL+url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Add("User-Agent", e.UserAgent)
	return e.Client.Do(req)
}

func (e *E621Session) getPosts(tags []string, limit int) (posts []E621Post) {
	resp, err := e.get("/post/index.json", map[string]string{
		"tags":  strings.Join(tags, " "),
		"limit": strconv.Itoa(limit),
	})

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return posts
	}
	json.Unmarshal(data, &posts)
	return posts

}
