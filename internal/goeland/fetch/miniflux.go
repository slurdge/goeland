package fetch

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/i18n"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/viper"
)

// search queries can take a while on large instances
const minifluxTimeout = time.Second * 10

type minifluxResponse struct {
	Total   int             `json:"total"`
	Entries []minifluxEntry `json:"entries"`
}

type minifluxEntry struct {
	ID          int64               `json:"id"`
	Title       string              `json:"title"`
	URL         string              `json:"url"`
	Content     string              `json:"content"`
	PublishedAt time.Time           `json:"published_at"`
	Enclosures  []minifluxEnclosure `json:"enclosures"`
	Feed        minifluxFeed        `json:"feed"`
}

type minifluxFeed struct {
	Title    string           `json:"title"`
	SiteURL  string           `json:"site_url"`
	Category minifluxCategory `json:"category"`
}

type minifluxCategory struct {
	Title string `json:"title"`
}

type minifluxEnclosure struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
}

func fetchMiniflux(source *goeland.Source, url string, allowInsecure bool) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	config := viper.GetViper()
	apiToken := config.GetString("miniflux-api-token")
	if apiToken == "" {
		log.Warnln("Miniflux may fail: miniflux-api-token is empty")
	}
	req.Header.Set("X-Auth-Token", apiToken)

	client := http.Client{Timeout: minifluxTimeout}
	if allowInsecure {
		log.Warningf("ignoring certificate security for url: %s\n", url)
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("received error code %d", res.StatusCode)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	miniflux := new(minifluxResponse)
	if err := json.Unmarshal(data, miniflux); err != nil {
		return err
	}
	if miniflux.Total > len(miniflux.Entries) {
		log.Warningf("miniflux returned %d out of %d matching entries, increase the 'limit' query parameter of url: %s", len(miniflux.Entries), miniflux.Total, url)
	}
	for _, minifluxEntry := range miniflux.Entries {
		entry := goeland.Entry{}
		entry.Source = source
		entry.Title = html.UnescapeString(policy.Sanitize(minifluxEntry.Title))
		entry.Content = minifluxEntry.Content
		if !viper.GetBool("unsafe-no-sanitize-filter") {
			entry.Content = policy.Sanitize(entry.Content)
		}
		entry.UID = fmt.Sprintf("miniflux-%d", minifluxEntry.ID)
		entry.Date = minifluxEntry.PublishedAt
		entry.URL = minifluxEntry.URL
		for _, enclosure := range minifluxEntry.Enclosures {
			if strings.HasPrefix(enclosure.MimeType, "image") {
				entry.ImageURL = enclosure.URL
				break
			}
		}
		source.Entries = append(source.Entries, entry)
	}

	// single feed or category urls carry a better title than the generic one
	if len(miniflux.Entries) > 0 {
		first := miniflux.Entries[0]
		if strings.Contains(req.URL.Path, "/v1/feeds/") {
			source.Title = first.Feed.Title
			source.URL = first.Feed.SiteURL
		} else if strings.Contains(req.URL.Path, "/v1/categories/") {
			source.Title = first.Feed.Category.Title
		}
	}
	if source.Title == "" {
		source.Title = i18n.T("Miniflux entries")
	}
	if source.URL == "" {
		source.URL = url
	}

	return nil
}
