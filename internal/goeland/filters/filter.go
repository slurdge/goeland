package filters

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	_ "github.com/slurdge/goeland/internal/goeland/i18n" //needed for now as we do not properly do i8n
	"github.com/slurdge/goeland/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	"all":        {"Default, include all entries", filterAll},
	"none":       {"Removes all entries", filterNone},
	"first":      {"Keep only the first (optional: N) entry. Use either 'first'  or 'first(3)'", filterFirst},
	"last":       {"Keep only the last  (optional: N) entry. Use either 'last'  or 'last(3)'", filterLast},
	"random":     {"Keep 1 or more random entries. Use either 'random' or 'random(5)'", filterRandom},
	"reverse":    {"Reverse the order of the entries", filterReverse},
	"today":      {"Keep only the entries for today", filterToday},
	"lasthours":  {"Keep only the entries that are from the X last hours (default 24)", filterLastHour},
	"digest":     {"Make a digest of all entries (optional heading level, default is " + fmt.Sprint(defaultHeaderLevel) + ")", filterDigest},
	"combine":    {"Combine all the entries into one source and use the first entry title as source title. Useful for merge sources", filterCombine},
	"links":      {`Rewrite relative links src="// and href="// to have an https:// prefix`, filterRelativeLinks},
	"embedimage": {`Embed a picture if the entry has an attachment with a type of picture (optional position: top|bottom|left|right, default is top)`, filterEmbedImage},
	"replace": {`Replace a string with another. Use with an argument like this: replace(myreplace) and define
		[replace.myreplace]
		from="A string"
		to="Another string"
	  in your config file.`, filterReplace},
	"includelink": {"Include the link of entries in the digest form", filterIncludeLink},
	"includesourcetitle": {"Include source titles of the entries in the digest form. Useful for merge sources", filterIncludeSourceTitle},
	"language":    {"Keep only the specified languages (best effort detection), use like this: language(en,de)", filterLanguage},
	"unseen":      {"Keep only unseen entry", filterUnSeen},
	"lebrief":     {"Deprecated. Use retrieve(div.content) instead. Retrieves the full excerpts for Next INpact's Lebrief", filterLeBrief},
	"retrieve":    {"Retrieves the full content from a goquery", filterRetrieveContent},
	"untrack":     {"Removes feedburner pixel tracking", filterUntrack},
	"reddit":      {"Better formatting for reddit rss", filterReddit},
	"sanitize":    {"Sanitize the content of entries (to be used if --unsafe-no-sanitize-filter was passed)", filterSanitize},
	"toc":         {"Create a table of content entry for all the entries (optional: title, will use the Title as a link)", filterToc},
	"limitwords":  {"Limit the number of words in the entry. Use limitwords(number).", filterLimitWords},
	"reskip":      {"Skip entries whose titles match a regular expression. Use reskip(regex).", filterRESkip},
}

// GetFiltersHelp returns a string that contains help for all filters
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
	number := extractNumber(source, params)
	source.Entries = source.Entries[:number]
}

func filterLast(source *goeland.Source, params *filterParams) {
	number := extractNumber(source, params)
	source.Entries = source.Entries[len(source.Entries)-number:]
}

func filterRandom(source *goeland.Source, params *filterParams) {
	number := extractNumber(source, params)
	rand.Shuffle(len(source.Entries), func(i, j int) {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	})
	source.Entries = source.Entries[:number]
}

func extractNumber(source *goeland.Source, params *filterParams) int {
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
	return number
}

func filterReverse(source *goeland.Source, params *filterParams) {
	for i, j := 0, len(source.Entries)-1; i < j; i, j = i+1, j-1 {
		source.Entries[i], source.Entries[j] = source.Entries[j], source.Entries[i]
	}
}

func filterLastHour(source *goeland.Source, params *filterParams) {
	hours := 24
	if len(params.args) > 0 {
		hours, _ = strconv.Atoi(params.args[0])
	}
	var current int
	for _, entry := range source.Entries {
		if entry.Date.Before(time.Now().Add(time.Hour * time.Duration(-hours))) {
			continue
		}
		source.Entries[current] = entry
		current++
	}
	source.Entries = source.Entries[:current]
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
	i8n := message.NewPrinter(language.BritishEnglish)
	digest.Title = i8n.Sprintf("Digest for %s", source.Title)
	if useFirstEntryTitle && len(source.Entries) > 0 {
		digest.Title = source.Entries[0].Title
	}
	var previousSource *goeland.Source
	content := ""
	for _, entry := range source.Entries {
		if entry.IncludeSourceTitle && previousSource != entry.Source  {
			content += fmt.Sprintf(`<h%d><a href="%s">%s</a></h%d>`, level - 1, entry.Source.URL, entry.Source.Title, level - 1)
			previousSource = entry.Source
		}
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

func getBaseURL(URL string) (string, error) {
	url, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	url.Path = ""
	url.RawQuery = ""
	url.Fragment = ""
	return url.String(), nil
}

func filterRelativeLinks(source *goeland.Source, params *filterParams) {
	re := regexp.MustCompile(`(src|href)\s*=\s*('|")\/\/`)
	for i, entry := range source.Entries {
		entry.Content = re.ReplaceAllString(entry.Content, "${1}=${2}https://")
		source.Entries[i] = entry
	}
	baseURL, err := getBaseURL(source.URL)
	if err != nil {
		return
	}
	for !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	re = regexp.MustCompile(`(src|href)\s*=\s*('|")\/`)
	for i, entry := range source.Entries {
		entry.Content = re.ReplaceAllString(entry.Content, "${1}=${2}"+baseURL)
		source.Entries[i] = entry
	}
}

func filterReplace(source *goeland.Source, params *filterParams) {
	key := ""
	if len(params.args) > 0 {
		key = params.args[0]
	}
	config := params.config
	from := config.GetString(fmt.Sprintf("replace.%s.from", key))
	to := config.GetString(fmt.Sprintf("replace.%s.to", key))
	for i, entry := range source.Entries {
		entry.Content = strings.ReplaceAll(entry.Content, from, to)
		source.Entries[i] = entry
	}
}

func filterSanitize(source *goeland.Source, params *filterParams) {
	for i, entry := range source.Entries {
		entry.Content = policy.Sanitize(entry.Content)
		source.Entries[i] = entry
	}
}

func filterIncludeLink(source *goeland.Source, params *filterParams) {
	for i := range source.Entries {
		source.Entries[i].IncludeLink = true
	}
}

func filterIncludeSourceTitle(source *goeland.Source, params *filterParams) {
	for i := range source.Entries {
		source.Entries[i].IncludeSourceTitle = true
	}
}

func filterEmbedImage(source *goeland.Source, params *filterParams) {
	args := params.args
	positions := []string{"top", "bottom", "left", "right"}
	position := 0
	if len(args) > 0 {
		for i, v := range positions {
			if v == args[0] {
				position = i
			}
		}
	}
	includeEntryLink := false
	if len(args) > 1 {
		if strings.ToLower(args[1]) == "link" {
			includeEntryLink = true
		}
	}
	for i, entry := range source.Entries {
		if strings.TrimSpace(entry.ImageURL) == "" {
			continue
		}
		imageLink := fmt.Sprintf(`<img src="%s" class="%s">`, entry.ImageURL, positions[position])
		if includeEntryLink {
			imageLink = fmt.Sprintf(`<a href="%s">%s</a>`, entry.URL, imageLink)
		}
		switch position {
		case 0:
			entry.Content = imageLink + entry.Content
		case 1:
			entry.Content = entry.Content + imageLink
		case 2:
			entry.Content = imageLink + entry.Content + `<br style="clear:both" />`
		case 3:
			entry.Content = imageLink + entry.Content + `<br style="clear:both" />`
		}

		source.Entries[i] = entry
	}
}

func filterToc(source *goeland.Source, params *filterParams) {
	if len(source.Entries) == 0 {
		return
	}
	toc := goeland.Entry{}
	args := params.args
	if len(args) > 0 && strings.ToLower(args[0]) == "title" {
		toc.Title = fmt.Sprintf(`<a href="%s">%s</a>`, source.URL, source.Title)
	} else {
		toc.Title = fmt.Sprintf("Table of Content for %s", source.Title)
	}
	content := "<ul>"
	for _, entry := range source.Entries {
		if entry.IncludeLink {
			content += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, entry.URL, entry.Title)
		} else {
			content += fmt.Sprintf("<li>%s</li>", entry.Title)
		}
	}
	content += "</ul>"
	h := sha256.New()
	h.Write([]byte(content))
	toc.UID = fmt.Sprintf("%x", h.Sum(nil))
	toc.Date = time.Now()
	toc.Content = content
	source.Entries = append([]goeland.Entry{toc}, source.Entries...)
}

func isEndOfSentence(r rune) bool {
	return r == '.' || r == '?' || r == '!' || r == '"' || r == '\n'
}

// taken from https://github.com/gohugoio/hugo/blob/3854a6fa6c323d1c09aa71a0626c9eef62709294/helpers/content.go#L219
func truncateWordsToWholeSentence(s string, summaryLength int) (string, bool) {
	var (
		wordCount     = 0
		lastWordIndex = -1
	)

	for i, r := range s {
		if unicode.IsSpace(r) {
			wordCount++
			lastWordIndex = i

			if wordCount >= summaryLength {
				break
			}

		}
	}

	if lastWordIndex == -1 {
		return s, false
	}

	endIndex := -1

	for j, r := range s[lastWordIndex:] {
		if isEndOfSentence(r) {
			endIndex = j + lastWordIndex + utf8.RuneLen(r)
			break
		}
	}

	if endIndex == -1 {
		return s, false
	}

	return strings.TrimSpace(s[:endIndex]), endIndex < len(s)
}

func filterLimitWords(source *goeland.Source, params *filterParams) {
	number := -1
	if len(params.args) > 0 {
		number, _ = strconv.Atoi(params.args[0])
	}
	for i, entry := range source.Entries {
		if number > 0 {
			//var isTruncated bool
			//add a more link here ?
			entry.Content, _ = truncateWordsToWholeSentence(entry.Content, number)
		}
		source.Entries[i] = entry
	}
}

func filterRESkip(source *goeland.Source, params *filterParams) {
	args := params.args
	if len(args) != 1 {
		log.Errorf("reskip takes exactly one parameter; got %d", len(args))
		return
	}
	re, err := regexp.Compile(args[0])
	if err != nil {
		log.Errorf("error compiling reskip regex: %v", err)
		return
	}

	var current int
	for _, entry := range source.Entries {
		if re.Match([]byte(entry.Title)) {
			log.Debugf("skipping entry with title '%s' due to regex match '%s'", entry.Title, args[0])
			continue
		}
		source.Entries[current] = entry
		current++
	}
	source.Entries = source.Entries[:current]
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
