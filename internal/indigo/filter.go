package indigo

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"time"

	"github.com/slurdge/indigo/config"
	"github.com/slurdge/indigo/log"
)

func (source *Source) filterAll() {

}

func (source *Source) filterNone() {
	source.Entries = nil
}

func (source *Source) filterFirst() {
	source.Entries = source.Entries[:1]
}

func (source *Source) filterReverse() {
	for i, j := 0, len(source.Entries)-1; i < j; i, j = i+1, j-1 {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	}
}

func (source *Source) filterToday() {
	var current int
	for _, entry := range source.Entries {
		if entry.Date.Day() != time.Now().Day() {
			continue
		}
		source.Entries[current] = entry
		current++
	}
	source.Entries = source.Entries[:current]
}

func (source *Source) filterDigest(level int, useFirstEntryTitle bool) {
	if len(source.Entries) <= 1 {
		return
	}
	digest := Entry{}
	digest.Title = fmt.Sprintf("Digest for %s", source.Title)
	if useFirstEntryTitle && len(source.Entries) > 0 {
		digest.Title = source.Entries[0].Title
	}
	content := ""
	for _, entry := range source.Entries {
		content += fmt.Sprintf("<h%d>%s</h%d>", level, entry.Title, level)
		content += entry.Content
	}
	h := sha256.New()
	h.Write([]byte(content))
	digest.UID = fmt.Sprintf("%x", h.Sum(nil))
	digest.Date = time.Now()
	digest.Content = content
	source.Entries = []Entry{digest}
}

func (source *Source) filterRelativeLinks() {
	re := regexp.MustCompile(`(src|href)\s*=('|")\/\/`)
	for i, entry := range source.Entries {
		entry.Content = re.ReplaceAllString(entry.Content, "${1}=${2}https://")
		source.Entries[i] = entry
	}
}

func filterSource(config config.Provider, source *Source) {
	log.Infof("Retrieved %v feeds", len(source.Entries))
	filters := SplitAndTrimString(config.GetString(fmt.Sprintf("sources.%s.filters", source.Name)))
	for _, filter := range filters {
		switch filter {
		case "all":
			source.filterAll()
		case "today":
			source.filterToday()
		case "lebrief":
			source.filterLeBrief()
		case "digest":
			source.filterDigest(1, false)
		case "digest2":
			source.filterDigest(2, false)
		case "combine":
			source.filterDigest(1, true)
		case "reverse":
			source.filterReverse()
		case "wikipedia":
			source.filterWikipedia()
		case "links":
			source.filterRelativeLinks()
		case "first":
			source.filterFirst()
		default:
			log.Errorf("unknown filter: %s\n", filter)
		}
		log.Infof("After %s: %v feeds", filter, len(source.Entries))
		log.Debugf("After %s: %+v", filter, source.Entries)
	}
}
