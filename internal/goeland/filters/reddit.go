package filters

import (
	"html"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/httpget"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/viper"
)

func filterReddit(source *goeland.Source, params *filterParams) {
	policy := bluemonday.NewPolicy()
	policy.AllowImages()
	policy.AllowStandardURLs()
	policy.AllowAttrs("href").OnElements("a")
	policy.AllowElements("p")

	re := regexp.MustCompile(`\/comments\/([a-z0-9]+)\/`)
	for i, entry := range source.Entries {
		postId := re.FindStringSubmatch(entry.URL)[1]
		if !viper.GetBool("unsafe-no-sanitize-filter") {
			entry.Content = policy.Sanitize(entry.Content)
		}
		if strings.Contains(entry.Content, "b.thumbs.redditmedia.com") ||
			strings.Contains(entry.Content, "external-preview.redd.it") {
			//we consider this is only a picture post
			image_url, err := getBetterPreview(postId)
			if err != nil {
				log.Warningf("Cannot get better preview picture %v, ignoring", err)
				continue
			}
			ore := regexp.MustCompile(`<img\s+src="[^"]+"[^>]*>`)
			entry.Content = ore.ReplaceAllString(entry.Content, `<img src="`+image_url+`">`)
		}
		source.Entries[i] = entry
	}
}

func getBetterPreview(postId string) (string, error) {
	jsonURL := "https://api.reddit.com/api/info/?id=t3_" + postId

	body, err := httpget.GetHTTPRessource(jsonURL)
	if err != nil {
		return "", err
	}

	media_id, err := jsonparser.GetString(body, "data", "children", "[0]", "data", "gallery_data", "items", "[0]", "media_id")
	if err == nil {
		image_url, err := jsonparser.GetString(body, "data", "children", "[0]", "data", "media_metadata", media_id, "p", "[3]", "u")
		if err == nil {
			image_url = html.UnescapeString(image_url)
			return image_url, nil
		}
	}
	preview, err := jsonparser.GetString(body, "data", "children", "[0]", "data", "preview", "images", "[0]", "source", "url")
	if err != nil {
		return "", err
	}
	return html.UnescapeString(preview), nil
}
