package goeland

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/log"
)

type filter struct {
	help       string
	filterFunc func(source *Source, params *filterParams)
}

type filterParams struct {
	args   []string
	config config.Provider
}

var filters = map[string]filter{
	"all":     filter{"Default, include all entries", filterAll},
	"none":    filter{"Removes all entries", filterNone},
	"first":   filter{"Keep only the first entry", filterFirst},
	"last":    filter{"Keep only the last entry", filterLast},
	"reverse": filter{"Reverse the order of the entries", filterReverse},
	"today":   filter{"Keep only the entries for today", filterToday},
	"digest":  filter{"Make a digest of all entries (optional heading level, default is 1)", filterDigest},
	"combine": filter{"Combine all the entries into one source and use the first entry title as source title. Useful for merge sources", filterCombine},
	"links":   filter{`Rewrite relative links src="// and href="// to have an https:// prefix`, filterRelativeLinks},
	"replace": filter{`Replace a string with another. Use with an argument like this: replace(myreplace) and define
[replace.myreplace]
from="A string"
to="Another string"`, filterReplace},
	"language": filter{"Keep only the specified languages (best effort detection), use like this: language(en,de)", filterLanguage},
	"lebrief":  filter{"Retrieves the full excerpts for Next INpact's Lebrief", filterLeBrief},
}

func GetFiltersHelp() string {
	lines := []string{}
	for name, value := range filters {
		line := fmt.Sprintf("\t- %s: %s", name, value.help)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func filterAll(source *Source, params *filterParams) {

}

func filterNone(source *Source, params *filterParams) {
	source.Entries = nil
}

func filterFirst(source *Source, params *filterParams) {
	source.Entries = source.Entries[:1]
}

func filterLast(source *Source, params *filterParams) {
	source.Entries = source.Entries[len(source.Entries)-1:]
}

func filterReverse(source *Source, params *filterParams) {
	for i, j := 0, len(source.Entries)-1; i < j; i, j = i+1, j-1 {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	}
}

func filterToday(source *Source, params *filterParams) {
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

func filterDigestGeneric(source *Source, level int, useFirstEntryTitle bool) {
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

func filterDigest(source *Source, params *filterParams) {
	args := params.args
	level := 1
	if len(args) > 0 {
		level, _ = strconv.Atoi(args[0])
	}
	filterDigestGeneric(source, level, false)
}

func filterCombine(source *Source, params *filterParams) {
	filterDigestGeneric(source, 1, true)
}

func filterRelativeLinks(source *Source, params *filterParams) {
	re := regexp.MustCompile(`(src|href)\s*=('|")\/\/`)
	for i, entry := range source.Entries {
		entry.Content = re.ReplaceAllString(entry.Content, "${1}=${2}https://")
		source.Entries[i] = entry
	}
}

func filterReplace(source *Source, params *filterParams) {
	key := params.args[0]
	config := params.config
	from := config.GetString(fmt.Sprintf("replace.%s.from", key))
	to := config.GetString(fmt.Sprintf("replace.%s.to", key))
	for i, entry := range source.Entries {
		entry.Content = strings.ReplaceAll(entry.Content, from, to)
		source.Entries[i] = entry
	}
}

func filterSource(source *Source, config config.Provider) {
	log.Infof("Retrieved %v feeds", len(source.Entries))
	filterNames := config.GetStringSlice(fmt.Sprintf("sources.%s.filters", source.Name))
	for _, filterName := range filterNames {
		filterShort := filterName
		args := []string{}
		if strings.Contains(filterName, "(") {
			filterShort = strings.Split(filterName, "(")[0]
			args = strings.Split(strings.ReplaceAll(strings.Split(filterName, "(")[1], ")", ""), ",")
			for i, arg := range args {
				args[i] = strings.TrimSpace(arg)
			}
		}
		if filter, found := filters[filterShort]; found {
			log.Debugf("Executing %s filter with args: %v", filterShort, args)
			params := filterParams{args: args, config: config}
			filter.filterFunc(source, &params)
		} else {
			log.Errorf("unknown filter: %s\n", filterName)
		}
		log.Infof("After %s: %v feeds", filterName, len(source.Entries))
		log.Debugf("After %s: %+v", filterName, source.Entries)
	}
}
