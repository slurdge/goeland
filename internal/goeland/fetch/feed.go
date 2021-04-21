package fetch

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/version"
)

const minContentLen = 10

var policy *bluemonday.Policy

//from https://github.com/mmcdole/gofeed/issues/74#
type userAgentTransport struct {
	http.RoundTripper
}

func (c *userAgentTransport) roundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "multiple:goeland:"+version.Version+" (commit id:"+version.GitCommit+") (by /u/goelandrss)")
	return c.RoundTripper.RoundTrip(r)
}

func fetchFeed(source *goeland.Source, feedLocation string, isFile bool) error {
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	var err error
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
		fp.Client = &http.Client{
			Transport: &userAgentTransport{http.DefaultTransport},
		}
		feed, err = fp.ParseURL(feedLocation)
		if err != nil {
			return fmt.Errorf("cannot open or parse url: %s", feedLocation)
		}
	}
	for _, item := range feed.Items {

		entry := goeland.Entry{}
		entry.Title = html.UnescapeString(item.Title)
		contentLength := len(strings.TrimSpace(item.Content))
		descriptionLength := len(strings.TrimSpace(item.Description))
		entry.Content = html.UnescapeString(item.Description)
		// 'smart' check to get the most interresting content
		if descriptionLength < minContentLen || contentLength > descriptionLength+minContentLen {
			entry.Content = html.UnescapeString(item.Content)
		}
		entry.Title = policy.Sanitize(entry.Title)
		entry.Content = policy.Sanitize(entry.Content)
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
		source.Entries = append(source.Entries, entry)
	}
	source.Title = feed.Title
	return nil
}

func init() {
	policy = bluemonday.UGCPolicy()
}
