package fetch

import (
	"bytes"
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
	Hash        string              `json:"hash"`
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

type minifluxEntries struct {
	IDs    []int64 `json:"entry_ids"`
	Status string  `json:"status"`
}

func fetchMiniflux(source *goeland.Source, url string, apiToken string, allowInsecure bool, markAsRead bool, dryRun bool) error {

	if apiToken == "" {
		log.Warnln("Miniflux may fail: no api token given (set miniflux-api-token or the source's api-token)")
	}

	client := http.Client{Timeout: minifluxTimeout}
	if allowInsecure {
		log.Warningf("ignoring certificate security for url: %s\n", url)
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Auth-Token", apiToken)

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
	ids := make([]int64, len(miniflux.Entries))
	for index, minifluxEntry := range miniflux.Entries {
		entry := goeland.Entry{}
		entry.Source = source
		entry.Title = html.UnescapeString(policy.Sanitize(minifluxEntry.Title))
		entry.Content = minifluxEntry.Content
		if !viper.GetBool("unsafe-no-sanitize-filter") {
			entry.Content = policy.Sanitize(entry.Content)
		}
		// the hash survives feed re-creations and instance migrations, unlike the numeric id
		entry.UID = "miniflux-" + minifluxEntry.Hash
		if minifluxEntry.Hash == "" {
			entry.UID = fmt.Sprintf("miniflux-%d", minifluxEntry.ID)
		}
		entry.Date = minifluxEntry.PublishedAt
		entry.URL = minifluxEntry.URL
		for _, enclosure := range minifluxEntry.Enclosures {
			if strings.HasPrefix(enclosure.MimeType, "image") {
				entry.ImageURL = enclosure.URL
				break
			}
		}
		source.Entries = append(source.Entries, entry)
		ids[index] = minifluxEntry.ID
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

	if markAsRead && len(ids) > 0 {
		if dryRun {
			log.Infof("Dry run: not marking %d entries read for source: %s", len(ids), source.Name)
			return nil
		}
		log.Infof("Marking %d entries read for source: %s", len(ids), source.Name)
		// We build the PUT url from the request URL by splitting on the "v1" as we don't have much to work with
		base := strings.SplitN(url, "/v1/", 2)
		if len(base) < 2 {
			log.Infof("Cannot find the base URL for miniflux: %s, skipping mark-as-read", url)
			return nil
		}
		markURL := fmt.Sprintf("%s/v1/entries", base[0])
		payload := &minifluxEntries{IDs: ids, Status: "read"}
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(http.MethodPut, markURL, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("X-Auth-Token", apiToken)
		req.Header.Set("Content-Type", "application/json")
		res, err = client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusNoContent {
			return fmt.Errorf("received non-204 status %d marking entries read for source: %s", res.StatusCode, source.Name)
		}
	}

	return nil
}
