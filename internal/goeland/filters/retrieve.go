package filters

import (
	"bytes"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/httpget"
	"github.com/spf13/viper"
)

var policy *bluemonday.Policy

// Deprecated: use filterRetrieveContent instead
func filterLeBrief(source *goeland.Source, params *filterParams) {
	params.args = []string{"div.content"}
	filterRetrieveContent(source, params)
}

func filterRetrieveContent(source *goeland.Source, params *filterParams) {
	args := params.args
	if len(args) < 1 {
		return
	}
	query := args[0]
	for index, entry := range source.Entries {
		link := entry.URL
		body, err := httpget.GetHTTPRessource(link)
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			continue
		}
		base, err := url.Parse(link)
		if err != nil {
			continue
		}
		fullContent := doc.Find(query)

		makeAttrFilter := func(attr string) func(_ int, selection *goquery.Selection) {
			return func(i int, selection *goquery.Selection) {
				src, exist := selection.Attr(attr)
				if !exist {
					return
				}
				relative, err := url.Parse(src)
				if err != nil {
					return
				}
				selection.SetAttr(attr, base.ResolveReference(relative).String())
			}
		}
		srcFilter := makeAttrFilter("src")
		hrefFilter := makeAttrFilter("href")
		fullContent.Find("img").Each(srcFilter)
		fullContent.Find("a").Each(hrefFilter)
		html, err := fullContent.Html()
		if err != nil {
			continue
		}
		if !viper.GetBool("unsafe-no-sanitize-filter") {
			entry.Content = policy.Sanitize(html)
		} else {
			entry.Content = html
		}
		source.Entries[index] = entry
	}
}

func init() {
	policy = bluemonday.UGCPolicy()
}
