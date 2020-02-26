package indigo

import (
	"fmt"
	"strings"
	"time"

	"github.com/slurdge/indigo/config"
	"github.com/slurdge/indigo/log"
)

// Entry This represent an entry produced by a source
type Entry struct {
	UID     string
	Title   string
	Content string
	Date    time.Time
}

// Source ...
type Source struct {
	Title   string
	Entries []Entry
}

// GetSource ...
func GetSource(config config.Provider, sourceName string) (*Source, error) {
	url := config.GetString(fmt.Sprintf("sources.%s.url", sourceName))

	//todo: get type of feed
	content, err := fetchFeed(url, !strings.HasPrefix(url, "http"))
	if err != nil {
		log.Errorf("Cannot retrieve feed: %s error: %v", sourceName, err)
		return content, err
	}
	return content, nil
}
