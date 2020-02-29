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
	Name    string
	Title   string
	Entries []Entry
}

// GetSource ...
func GetSource(config config.Provider, sourceName string) (*Source, error) {
	if !config.IsSet(fmt.Sprintf("sources.%s", sourceName)) {
		return nil, fmt.Errorf("cannot find source: %s", sourceName)
	}
	url := config.GetString(fmt.Sprintf("sources.%s.url", sourceName))
	sourceType := config.GetString(fmt.Sprintf("sources.%s.type", sourceName))
	log.Debugf("Fetching source: %s of type %s", sourceName, sourceType)
	var source *Source
	var err error
	switch sourceType {
	case "feed":
		source, err = fetchFeed(url, !strings.HasPrefix(url, "http"))
		if err != nil {
			log.Errorf("Cannot retrieve feed: %s error: %v", sourceName, err)
			return source, err
		}
	case "merge":
		subSourceNames := SplitAndTrimString(config.GetString(fmt.Sprintf("sources.%s.sources", sourceName)))
		for _, subSourceName := range subSourceNames {
			currentSource, err := GetSource(config, subSourceName)
			if source == nil {
				source = currentSource
			} else {
				source.Entries = append(source.Entries, currentSource.Entries...)
			}
			if err != nil {
				log.Errorf("cannot fetch source: %s (%v)", subSourceName, err)
				continue
			}
		}
	default:
		return nil, fmt.Errorf("cannot understand source type: %s", sourceType)
	}
	if source != nil {
		source.Name = sourceName
		filterSource(config, source)
	}
	log.Debugf("%v", source)
	return source, nil
}
