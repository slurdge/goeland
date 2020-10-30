package fetch

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/version"
)

const minContentLen = 10

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
		entry.Content = html.UnescapeString(item.Description)
		if len(strings.TrimSpace(entry.Content)) < minContentLen {
			entry.Content = html.UnescapeString(item.Content)
		}
		entry.UID = item.GUID
		if item.PublishedParsed != nil {
			entry.Date = *item.PublishedParsed
		} else {
			entry.Date = time.Now()
		}
		entry.URL = item.Link
		source.Entries = append(source.Entries, entry)
	}
	source.Title = feed.Title
	return nil
}
