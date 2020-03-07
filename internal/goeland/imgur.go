package goeland

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type imgurImage struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Animated    bool   `json:"animated"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type imgurItem struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	Link     string       `json:"link"`
	DateTime int64        `json:"datetime"`
	Images   []imgurImage `json:images`
}

type imgurData struct {
	Items []imgurItem `json:items`
}

type imgurRoot struct {
	Data imgurData `json:data`
}

var clientID = ""

func fetchImgurTag(tag string) (*Source, error) {
	source := new(Source)
	url := "https://api.imgur.com/3/gallery/t/" + tag + "/top/0/day"
	client := http.Client{Timeout: time.Second * 3}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return source, err
	}

	req.Header.Set("Authorization", "Client-ID "+clientID)

	res, err := client.Do(req)
	if err != nil {
		return source, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return source, err
	}
	imgurData := new(imgurRoot)
	if err := json.Unmarshal(data, imgurData); err != nil {
		return source, err
	}
	for _, item := range imgurData.Data.Items {
		if len(item.Images) < 1 {
			continue
		}
		image := item.Images[0]
		if image.Animated {
			continue
		}
		entry := Entry{}
		entry.Title = item.Title
		entry.Content = `<a href="` + item.Link + `"><img src="` + image.Link + `"></a><br>` + image.Description
		entry.UID = item.ID
		entry.Date = time.Unix(item.DateTime, 0)
		entry.URL = item.Link
		source.Entries = append(source.Entries, entry)
	}
	source.Title = "Imgur pictures for tag #" + tag
	return source, nil
}
