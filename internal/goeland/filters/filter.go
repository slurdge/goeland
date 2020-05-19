package filters

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/log"
)

type filter struct {
	help       string
	filterFunc func(source *goeland.Source, params *filterParams)
}

type filterParams struct {
	args   []string
	config config.Provider
}

const defaultHeaderLevel = 2

var filters = map[string]filter{
	"all":     filter{"Default, include all entries", filterAll},
	"none":    filter{"Removes all entries", filterNone},
	"first":   filter{"Keep only the first entry", filterFirst},
	"last":    filter{"Keep only the last entry", filterLast},
	"random":  filter{"Keep 1 or more random entries. Use either 'random' or 'random(5)' for example.", filterRandom},
	"reverse": filter{"Reverse the order of the entries", filterReverse},
	"today":   filter{"Keep only the entries for today", filterToday},
	"digest":  filter{"Make a digest of all entries (optional heading level, default is " + string(defaultHeaderLevel) + ")", filterDigest},
	"combine": filter{"Combine all the entries into one source and use the first entry title as source title. Useful for merge sources", filterCombine},
	"links":   filter{`Rewrite relative links src="// and href="// to have an https:// prefix`, filterRelativeLinks},
	"replace": filter{`Replace a string with another. Use with an argument like this: replace(myreplace) and define
		[replace.myreplace]
		from="A string"
		to="Another string"
	  in your config file.`, filterReplace},
	"includelink": filter{"Include the link of entries in the digest form", filterIncludeLink},
	"language":    filter{"Keep only the specified languages (best effort detection), use like this: language(en,de)", filterLanguage},
	"unseen":      filter{"Keep only unseen entry", filterUnSeen},
	"lebrief":     filter{"Retrieves the full excerpts for Next INpact's Lebrief", filterLeBrief},
}

func GetFiltersHelp() string {
	lines := []string{}
	for name, value := range filters {
		line := fmt.Sprintf("\t- %s: %s", name, value.help)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
func filterAll(source *goeland.Source, params *filterParams) {

}

func filterNone(source *goeland.Source, params *filterParams) {
	source.Entries = nil
}

func filterFirst(source *goeland.Source, params *filterParams) {
	source.Entries = source.Entries[:1]
}

func filterLast(source *goeland.Source, params *filterParams) {
	source.Entries = source.Entries[len(source.Entries)-1:]
}

func filterRandom(source *goeland.Source, params *filterParams) {
	number := 1
	if len(params.args) > 0 {
		number, _ = strconv.Atoi(params.args[0])
	}
	if number <= 0 {
		number = 1
	}
	if number > len(source.Entries) {
		number = len(source.Entries)
	}
	rand.Shuffle(len(source.Entries), func(i, j int) {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	})
	source.Entries = source.Entries[:number]
}

func filterReverse(source *goeland.Source, params *filterParams) {
	for i, j := 0, len(source.Entries)-1; i < j; i, j = i+1, j-1 {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	}
}

func filterToday(source *goeland.Source, params *filterParams) {
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

func filterDigestGeneric(source *goeland.Source, level int, useFirstEntryTitle bool) {
	if len(source.Entries) < 1 {
		return
	}
	digest := goeland.Entry{}
	digest.Title = fmt.Sprintf("Digest for %s", source.Title)
	if useFirstEntryTitle && len(source.Entries) > 0 {
		digest.Title = source.Entries[0].Title
	}
	content := ""
	for _, entry := range source.Entries {
		if entry.IncludeLink {
			content += fmt.Sprintf(`<h%d><a href="%s">%s</a></h%d>`, level, entry.URL, entry.Title, level)
		} else {
			content += fmt.Sprintf("<h%d>%s</h%d>", level, entry.Title, level)
		}
		content += entry.Content
	}
	h := sha256.New()
	h.Write([]byte(content))
	digest.UID = fmt.Sprintf("%x", h.Sum(nil))
	digest.Date = time.Now()
	digest.Content = content
	source.Entries = []goeland.Entry{digest}
}

func filterDigest(source *goeland.Source, params *filterParams) {
	args := params.args
	level := defaultHeaderLevel
	if len(args) > 0 {
		level, _ = strconv.Atoi(args[0])
	}
	filterDigestGeneric(source, level, false)
}

func filterCombine(source *goeland.Source, params *filterParams) {
	args := params.args
	level := defaultHeaderLevel
	if len(args) > 0 {
		level, _ = strconv.Atoi(args[0])
	}
	filterDigestGeneric(source, level, true)
}

func filterRelativeLinks(source *goeland.Source, params *filterParams) {
	re := regexp.MustCompile(`(src|href)\s*=('|")\/\/`)
	for i, entry := range source.Entries {
		entry.Content = re.ReplaceAllString(entry.Content, "${1}=${2}https://")
		source.Entries[i] = entry
	}
}

func filterReplace(source *goeland.Source, params *filterParams) {
	key := params.args[0]
	config := params.config
	from := config.GetString(fmt.Sprintf("replace.%s.from", key))
	to := config.GetString(fmt.Sprintf("replace.%s.to", key))
	for i, entry := range source.Entries {
		entry.Content = strings.ReplaceAll(entry.Content, from, to)
		source.Entries[i] = entry
	}
}

func filterIncludeLink(source *goeland.Source, params *filterParams) {
	for i, _ := range source.Entries {
		source.Entries[i].IncludeLink = true
	}
}

// FilterSource filters a source according to the config
func FilterSource(source *goeland.Source, config config.Provider) {
	log.Infof("Retrieved %v feeds for source %v", len(source.Entries), source.Name)
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
