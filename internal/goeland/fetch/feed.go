package fetch

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"html"
	"os"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/httpget"
	"github.com/spf13/viper"
)

const minContentLen = 10

var policy *bluemonday.Policy

func fetchFeed(source *goeland.Source, feedLocation string, isFile bool) error {
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	if isFile {
		file, err := os.Open(feedLocation)
		if err != nil {
			return fmt.Errorf("cannot open: %s", feedLocation)
		}
		defer file.Close()
		feed, err = fp.Parse(file)
		if err != nil {
			return fmt.Errorf("cannot parse file: %s", feedLocation)
		}
	} else {
		body, err := httpget.GetHTTPRessource(feedLocation)
		if err != nil {
			return fmt.Errorf("cannot open or parse url: %s (%v)", feedLocation, err)
		}
		feed, err = fp.ParseString(string(body))
		if err != nil {
			return fmt.Errorf("cannot open or parse url: %s (%v)", feedLocation, err)
		}
	}
	for _, item := range feed.Items {

		entry := goeland.Entry{}
		//order is important for title as Sanitize will escape
		entry.Title = item.Title
		entry.Title = policy.Sanitize(entry.Title)
		entry.Title = html.UnescapeString(entry.Title)
		contentLength := len(strings.TrimSpace(item.Content))
		descriptionLength := len(strings.TrimSpace(item.Description))
		entry.Content = html.UnescapeString(item.Description)
		// 'smart' check to get the most interesting content
		if descriptionLength < minContentLen || contentLength > descriptionLength+minContentLen {
			entry.Content = html.UnescapeString(item.Content)
		}
		if !viper.GetBool("unsafe-no-sanitize-filter") {
			entry.Content = policy.Sanitize(entry.Content)
		}
		entry.UID = item.GUID
		if item.PublishedParsed != nil {
			entry.Date = *item.PublishedParsed
		} else {
			entry.Date = time.Now()
		}
		entry.URL = item.Link
		if len(item.Enclosures) > 0 && strings.HasPrefix(item.Enclosures[0].Type, "image") {
			entry.ImageURL = item.Enclosures[0].URL
		} else if item.Image != nil && item.Image.URL != "" {
			entry.ImageURL = item.Image.URL
		} else {
			for _, extension := range item.Extensions["media"]["content"] {
				if strings.ToLower(extension.Name) == "content" {
					if strings.HasPrefix(extension.Attrs["type"], "image") && extension.Attrs["url"] != "" {
						entry.ImageURL = extension.Attrs["url"]
					}
				}
			}
		}
		if strings.TrimSpace(entry.UID) == "" {
			hash := fnv.New64a()
			hash.Write([]byte(entry.URL))
			entry.UID = hex.EncodeToString(hash.Sum([]byte{}))
		}
		source.Entries = append(source.Entries, entry)
	}
	source.Title = feed.Title
	source.URL = feed.Link
	if feed.Image != nil && feed.Image.URL != "" {
		source.ImageURL = feed.Image.URL
	}
	return nil
}

func init() {
	policy = bluemonday.UGCPolicy()
}
