package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/viper"
)

type MinifluxApiResponse struct {
	Total   int             `json:"total"`
	Entries []MinifluxEntry `json:"entries"`
}

type MinifluxEntry struct {
	ID          int                 `json:"id"`
	UserID      int                 `json:"user_id"`
	FeedID      int                 `json:"feed_id"`
	Title       string              `json:"title"`
	URL         string              `json:"url"`
	CommentsURL string              `json:"comments_url"`
	Author      string              `json:"author"`
	Content     string              `json:"content"`
	Hash        string              `json:"hash"`
	PublishedAt time.Time           `json:"published_at"`
	CreatedAt   time.Time           `json:"created_at"`
	Status      string              `json:"status"`
	ShareCode   string              `json:"share_code"`
	Starred     bool                `json:"starred"`
	ReadingTime int                 `json:"reading_time"`
	Enclosures  []MiniFluxEnclosure `json:"enclosures"`
	Feed        MiniFluxFeed        `json:"feed"`
}

type MiniFluxFeed struct {
	ID                  int              `json:"id"`
	UserID              int              `json:"user_id"`
	Title               string           `json:"title"`
	SiteURL             string           `json:"site_url"`
	FeedURL             string           `json:"feed_url"`
	CheckedAt           time.Time        `json:"checked_at"`
	EtagHeader          string           `json:"etag_header"`
	LastModifiedHeader  string           `json:"last_modified_header"`
	ParsingErrorMessage string           `json:"parsing_error_message"`
	ParsingErrorCount   int              `json:"parsing_error_count"`
	ScraperRules        string           `json:"scraper_rules"`
	RewriteRules        string           `json:"rewrite_rules"`
	Crawler             bool             `json:"crawler"`
	BlocklistRules      string           `json:"blocklist_rules"`
	KeeplistRules       string           `json:"keeplist_rules"`
	UserAgent           string           `json:"user_agent"`
	Username            string           `json:"username"`
	Password            string           `json:"password"`
	Disabled            bool             `json:"disabled"`
	IgnoreHTTPCache     bool             `json:"ignore_http_cache"`
	FetchViaProxy       bool             `json:"fetch_via_proxy"`
	Category            MiniFluxCategory `json:"category"`
	Icon                MiniFluxIcon     `json:"icon"`
}

type MiniFluxEnclosure struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	EntryID          int    `json:"entry_id"`
	URL              string `json:"url"`
	MimeType         string `json:"mime_type"`
	Size             int64  `json:"size"`
	MediaProgression int    `json:"media_progression"`
}

type MiniFluxCategory struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
}

type MiniFluxIcon struct {
	FeedID int `json:"feed_id"`
	IconID int `json:"icon_id"`
}

func fecthMiniFlux(source *goeland.Source, url string, allowInsecure bool) error {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := http.Client{Timeout: time.Second * 3}

	config := viper.GetViper()
	minifluxApiToken := config.GetString("miniflux-api-token")

	if minifluxApiToken == "" {
		log.Warnln("Miniflux may fail: miniflux-api-token is empty")
	}

	req.Header.Set("X-Auth-Token", minifluxApiToken)

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
	miniflux := new(MinifluxApiResponse)
	if err := json.Unmarshal(data, miniflux); err != nil {
		return err
	}
	for _, minifluxEntry := range miniflux.Entries {
		log.Debugf("%v", minifluxEntry)
		entry := goeland.Entry{}
		entry.Source = source
		entry.Title = minifluxEntry.Title
		entry.Content = minifluxEntry.Content
		entry.UID = fmt.Sprintf("miniflux-%i", minifluxEntry.ID)
		entry.Date = minifluxEntry.PublishedAt
		entry.URL = minifluxEntry.URL
		for _, enclosure := range minifluxEntry.Enclosures {
			if strings.HasPrefix(enclosure.MimeType, "image") {
				entry.ImageURL = enclosure.URL
			}
		}
		source.Entries = append(source.Entries, entry)
	}

	return nil
}
