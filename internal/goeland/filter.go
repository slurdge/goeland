package goeland

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/log"
)

type filter struct {
	help       string
	filterFunc func(source *Source)
}

var filters = map[string]filter{
	"all":     filter{"Default, include all entries", filterAll},
	"none":    filter{"Removes all entries", filterNone},
	"first":   filter{"Keep only the first entry", filterFirst},
	"last":    filter{"Keep only the last entry", filterLast},
	"reverse": filter{"Reverse the order of the entries", filterReverse},
	"today":   filter{"Keep only the entries for today", filterToday},
	"links":   filter{`Rewrite relative links src="// and href="// to have an https:// prefix`, filterRelativeLinks},
	"lebrief": filter{"Retrieves the full excerpts for Next INpact's Lebrief", filterLeBrief},
}

func GetFiltersHelp() string {
	lines := []string{}
	for name, value := range filters {
		line := fmt.Sprintf("\t- %s: %s", name, value.help)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func filterAll(source *Source) {

}

func filterNone(source *Source) {
	source.Entries = nil
}

func filterFirst(source *Source) {
	source.Entries = source.Entries[:1]
}

func filterLast(source *Source) {
	source.Entries = source.Entries[len(source.Entries)-1:]
}

func filterReverse(source *Source) {
	for i, j := 0, len(source.Entries)-1; i < j; i, j = i+1, j-1 {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	}
}

func filterToday(source *Source) {
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

func filterDigest(source *Source, level int, useFirstEntryTitle bool) {
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

func filterRelativeLinks(source *Source) {
	re := regexp.MustCompile(`(src|href)\s*=('|")\/\/`)
	for i, entry := range source.Entries {
		entry.Content = re.ReplaceAllString(entry.Content, "${1}=${2}https://")
		source.Entries[i] = entry
	}
}

func filterReplace(source *Source, key string, config config.Provider) {
	from := config.GetString(fmt.Sprintf("replace.%s.from", key))
	to := config.GetString(fmt.Sprintf("replace.%s.to", key))
	for i, entry := range source.Entries {
		entry.Content = strings.ReplaceAll(entry.Content, from, to)
		source.Entries[i] = entry
	}
}

func filterSource(config config.Provider, source *Source) {
	log.Infof("Retrieved %v feeds", len(source.Entries))
	filters := config.GetStringSlice(fmt.Sprintf("sources.%s.filters", source.Name))
	for _, filter := range filters {
		filterShort := filter
		filterParams := []string{}
		if strings.Contains(filter, "(") {
			filterShort = strings.Split(filter, "(")[0]
			filterParams = strings.Split(strings.ReplaceAll(strings.Split(filter, "(")[1], ")", ""), ",")
		}
		log.Infof("Executing %s filter with params: %v", filterShort, filterParams)
		switch filterShort {
		case "all":
			filterAll(source)
		case "today":
			filterToday(source)
		case "lebrief":
			filterLeBrief(source)
		case "digest":
			filterDigest(source, 1, false)
		case "digest2":
			filterDigest(source, 2, false)
		case "combine":
			filterDigest(source, 1, true)
		case "reverse":
			filterReverse(source)
		case "wikipedia":
			filterWikipedia(source)
		case "links":
			filterRelativeLinks(source)
		case "first":
			filterFirst(source)
		case "last":
			filterLast(source)
		case "replace":
			filterReplace(source, filterParams[0], config)
		case "language":
			filterLanguage(source, filterParams)
		default:
			log.Errorf("unknown filter: %s\n", filter)
		}
		log.Infof("After %s: %v feeds", filter, len(source.Entries))
		log.Debugf("After %s: %+v", filter, source.Entries)
	}
}
