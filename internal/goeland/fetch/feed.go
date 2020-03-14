package fetch

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/slurdge/goeland/internal/goeland"
)

const minContentLen = 10

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
		feed, err = fp.ParseURL(feedLocation)
		if err != nil {
			return fmt.Errorf("cannot open or parse url: %s", feedLocation)
		}
	}
	for _, item := range feed.Items {

		entry := goeland.Entry{}
		entry.Title = item.Title
		entry.Content = item.Description
		if len(strings.TrimSpace(entry.Content)) < minContentLen {
			entry.Content = item.Content
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
