package indigo

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (source *Source) FilterAll() {

}

func (source *Source) FilterNone() {
	source.Entries = nil
}

func (source *Source) FilterToday() {
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

func (source *Source) FilterDigest() {
	if len(source.Entries) <= 1 {
		return
	}
	digest := Entry{}
	digest.Title = fmt.Sprintf("Digest for %s", source.Title)
	content := ""
	for _, entry := range source.Entries {
		content += fmt.Sprintf("<h1>%s</h1>", entry.Title)
		content += entry.Content
	}
	h := sha256.New()
	h.Write([]byte(content))
	digest.UID = fmt.Sprintf("%x", h.Sum(nil))
	digest.Date = time.Now()
	digest.Content = content
	source.Entries = []Entry{digest}
}

func (source *Source) FilterLeBrief() {
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
