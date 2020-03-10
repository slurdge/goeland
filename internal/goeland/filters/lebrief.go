package filters

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/slurdge/goeland/internal/goeland"
)

func filterLeBrief(source *goeland.Source, params *filterParams) {
	for index, entry := range source.Entries {
		re := regexp.MustCompile(`<a href="([^"]+)">Lire la suite</a>`)
		link := re.FindStringSubmatch(entry.Content)[1]
		doc, err := goquery.NewDocument(link)
		if err != nil {
			continue
		}
		fullcontent, _ := doc.Find(".brief-inner-content").Html()
		entry.Content = fullcontent
		source.Entries[index] = entry
	}
}
