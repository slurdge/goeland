package fetch

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"html"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/httpget"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/viper"
)

const minContentLen = 10

var linkRegExp *regexp.Regexp
var newlineRegExp *regexp.Regexp

var policy *bluemonday.Policy

func fetchFeed(source *goeland.Source, feedLocation string, isFile bool, allowInsecure bool) error {
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	if isFile {
		file, err := os.Open(feedLocation)
		if err != nil {
			return fmt.Errorf("cannot open: %s", feedLocation)
		}
		defer file.Close()
		feed, err = fp.Parse(file)
		if err != nil {
			return fmt.Errorf("cannot parse file: %s", feedLocation)
		}
	} else {
		var body []byte
		var err error
		if !allowInsecure {
			body, err = httpget.GetHTTPRessource(feedLocation)
		} else {
			log.Warningf("ignoring certificate security for url: %s\n", feedLocation)
			body, err = httpget.GetHTTPRessourceInsecure(feedLocation)
		}

		if err != nil {
			return fmt.Errorf("cannot open or parse url: %s (%v)", feedLocation, err)
		}
		feed, err = fp.ParseString(string(body))
		if err != nil {
			return fmt.Errorf("cannot open or parse url: %s (%v)", feedLocation, err)
		}
	}
	for _, item := range feed.Items {

		entry := goeland.Entry{}
		//order is important for title as Sanitize will escape
		entry.Title = item.Title
		entry.Title = policy.Sanitize(entry.Title)
		entry.Title = html.UnescapeString(entry.Title)
		contentLength := len(strings.TrimSpace(item.Content))
		descriptionLength := len(strings.TrimSpace(item.Description))
		entry.Content = html.UnescapeString(item.Description)
		// 'smart' check to get the most interesting content
		if descriptionLength < minContentLen || contentLength > descriptionLength+minContentLen {
			entry.Content = html.UnescapeString(item.Content)
		}
		if !viper.GetBool("unsafe-no-sanitize-filter") {
			entry.Content = policy.Sanitize(entry.Content)
		}
		entry.UID = item.GUID
		if item.PublishedParsed != nil {
			entry.Date = *item.PublishedParsed
		} else {
			entry.Date = time.Now()
		}
		entry.URL = item.Link
		if len(item.Enclosures) > 0 && strings.HasPrefix(item.Enclosures[0].Type, "image") {
			entry.ImageURL = item.Enclosures[0].URL
		} else if item.Image != nil && item.Image.URL != "" {
			entry.ImageURL = item.Image.URL
		} else {
			for _, extension := range item.Extensions["media"]["content"] {
				if strings.ToLower(extension.Name) == "content" {
					if strings.HasPrefix(extension.Attrs["type"], "image") && extension.Attrs["url"] != "" {
						entry.ImageURL = extension.Attrs["url"]
					}
				}
			}
		}
		if len(entry.Content) == 0 && len(item.Extensions["media"]["group"]) > 0 {
			for _, extensions := range item.Extensions["media"]["group"][0].Children {
				if len(extensions) > 0 {
					extension := extensions[0]
					if strings.ToLower(extension.Name) == "description" {
						if strings.ToLower(extension.Attrs["type"]) == "html" {
							entry.Content = extension.Value
						} else {
							content := extension.Value
							content = linkRegExp.ReplaceAllString(content, `<a href="$1">$1</a>`)
							content = newlineRegExp.ReplaceAllString(content, "<br>")
							entry.Content = content
						}

					}
					if strings.ToLower(extension.Name) == "thumbnail" {
						entry.ImageURL = extension.Attrs["url"]
					}
				}
			}
		}
		if strings.TrimSpace(entry.UID) == "" {
			hash := fnv.New64a()
			hash.Write([]byte(entry.URL))
			entry.UID = hex.EncodeToString(hash.Sum([]byte{}))
		}
		source.Entries = append(source.Entries, entry)
	}
	source.Title = feed.Title
	source.URL = feed.Link
	return nil
}

func init() {
	policy = bluemonday.UGCPolicy()
	linkRegExp = regexp.MustCompile(`((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[.\!\/\\w]*))?)`)
	newlineRegExp = regexp.MustCompile(`\r?\n`)
}
