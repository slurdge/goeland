package indigo

import (
	"fmt"
	"os"

	"github.com/mmcdole/gofeed"
)

func fetchFeed(feedLocation string, isFile bool) (*Source, error) {
	source := new(Source)
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	var err error
	if isFile {
		file, err := os.Open(feedLocation)
		if err != nil {
			return source, fmt.Errorf("cannot open: %s", feedLocation)
		}
		defer file.Close()
		feed, err = fp.Parse(file)
		if err != nil {
			return source, fmt.Errorf("cannot parse file: %s", feedLocation)
		}
	} else {
		feed, err = fp.ParseURL(feedLocation)
		if err != nil {
			return source, fmt.Errorf("cannot open or parse url: %s", feedLocation)
		}
	}
	for _, item := range feed.Items {
		entry := Entry{}
		entry.Title = item.Title
		entry.Content = item.Description
		entry.UID = item.GUID
		entry.Date = *item.PublishedParsed
		source.Entries = append(source.Entries, entry)
	}
	source.Title = feed.Title
	return source, nil
}
