package filters

import (
	"strings"

	"github.com/abadojack/whatlanggo"
	"github.com/slurdge/goeland/internal/goeland"
	"jaytaylor.com/html2text"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func filterLanguage(source *goeland.Source, params *filterParams) {
	languages := params.args
	var current int
	for _, entry := range source.Entries {
		text, err := html2text.FromString(entry.Content, html2text.Options{OmitLinks: true})
		if err != nil {
			text = entry.Content
		}
		lang := whatlanggo.DetectLang(text)
		if !stringInSlice(lang.Iso6391(), languages) {
			continue
		}
		source.Entries[current] = entry
		current++
	}
	source.Entries = source.Entries[:current]
}
