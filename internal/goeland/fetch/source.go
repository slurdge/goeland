package fetch

import (
	"fmt"
	"strings"

	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/filters"
	"github.com/slurdge/goeland/log"
)

// FetchSource retrieves a source from either a feed, imgur or other sub-sources
func FetchSource(config config.Provider, sourceName string, parents []string) (*goeland.Source, error) {
	if !config.IsSet(fmt.Sprintf("sources.%s", sourceName)) {
		return nil, fmt.Errorf("cannot find source: %s", sourceName)
	}
	if filters.StringInSlice(strings.ToLower(sourceName), parents) {
		return nil, fmt.Errorf("source: %s. there seems to be a cycle in your configuration. Make sure no 'merge' source is referencing itself", sourceName)
	}
	sourceType := config.GetString(fmt.Sprintf("sources.%s.type", sourceName))
	log.Infof("Fetching source: %s of type %s", sourceName, sourceType)
	source := new(goeland.Source)
	var err error
	switch sourceType {
	case "feed":
		url := config.GetString(fmt.Sprintf("sources.%s.url", sourceName))
		allowInsecure := config.GetBool(fmt.Sprintf("sources.%s.allow-insecure", sourceName))
		err = fetchFeed(source, url, !strings.HasPrefix(url, "http"), allowInsecure)
		if err != nil {
			log.Errorf("Cannot retrieve feed: %s error: %v", url, err)
			return source, err
		}
	case "imgur":
		tag := config.GetString(fmt.Sprintf("sources.%s.tag", sourceName))
		sort := config.GetString(fmt.Sprintf("sources.%s.sort", sourceName))
		if !filters.StringInSlice(sort, []string{"top", "viral", "time"}) {
			sort = "top"
		}
		err = fetchImgurTag(source, tag, sort)
		if err != nil {
			log.Errorf("Cannot retrieve imgur tag: %s error: %v", tag, err)
			return source, err
		}
	case "merge":
		subSourceNames := config.GetStringSlice(fmt.Sprintf("sources.%s.sources", sourceName))
		for _, subSourceName := range subSourceNames {
			parents = append(parents, strings.ToLower(sourceName))
			subSource, err := FetchSource(config, subSourceName, parents)
			if err != nil {
				log.Errorf("cannot fetch source: %s (%v)", subSourceName, err)
				continue
			}
			source.Entries = append(source.Entries, subSource.Entries...)
		}
	default:
		return nil, fmt.Errorf("cannot understand source type: %s", sourceType)
	}
	source.Name = sourceName
	filters.FilterSource(source, config)
	log.Debugf("%v", source)
	return source, nil
}
