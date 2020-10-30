package filters

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/slurdge/goeland/internal/goeland"
)

func filterLeBrief(source *goeland.Source, params *filterParams) {
	for index, entry := range source.Entries {
		link := entry.URL
		doc, err := goquery.NewDocument(link)
		if err != nil {
			continue
		}
		fullcontent, _ := doc.Find("div.content").Html()
		entry.Content = fullcontent
		source.Entries[index] = entry
	}
}
