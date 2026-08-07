package filters

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/i18n"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	i18n.Init("en-US")
	log.SetDefaultLogger(log.NewLogger(viper.New()))
	os.Exit(m.Run())
}

// newTestSource creates a source with fully populated entries: UID, title,
// HTML content, URL, distinct dates, image URL and a back-pointer to the source.
func newTestSource(numEntries int) *goeland.Source {
	source := &goeland.Source{
		Name:  "test",
		Title: "Test Source",
		URL:   "https://example.com/feed.xml",
	}
	base := time.Date(2026, time.January, 15, 12, 0, 0, 0, time.UTC)
	for i := 0; i < numEntries; i++ {
		source.Entries = append(source.Entries, goeland.Entry{
			UID:      fmt.Sprintf("uid-%d", i),
			Title:    fmt.Sprintf("Entry %d", i),
			Content:  fmt.Sprintf(`<p>Content of entry %d with a <a href="https://example.com/articles/%d">link</a>.</p>`, i, i),
			URL:      fmt.Sprintf("https://example.com/articles/%d", i),
			Date:     base.Add(time.Duration(i) * time.Hour),
			ImageURL: fmt.Sprintf("https://example.com/images/%d.png", i),
			Source:   source,
		})
	}
	return source
}

func entryTitles(source *goeland.Source) []string {
	titles := make([]string, 0, len(source.Entries))
	for _, entry := range source.Entries {
		titles = append(titles, entry.Title)
	}
	return titles
}

func assertTitles(t *testing.T, source *goeland.Source, want ...string) {
	t.Helper()
	got := entryTitles(source)
	if len(got) != len(want) {
		t.Fatalf("got %d entries %v, want %d %v", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d: got title %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGetFiltersHelp(t *testing.T) {
	help := GetFiltersHelp()
	for _, name := range []string{"all", "digest", "reskip", "embedimage"} {
		if !strings.Contains(help, "- "+name+":") {
			t.Errorf("help is missing filter %q", name)
		}
	}
}

func TestFilterAll(t *testing.T) {
	for _, numEntries := range []int{0, 3} {
		source := newTestSource(numEntries)
		filterAll(source, &filterParams{})
		if len(source.Entries) != numEntries {
			t.Errorf("with %d entries: got %d, want unchanged", numEntries, len(source.Entries))
		}
	}
}

func TestFilterNone(t *testing.T) {
	for _, numEntries := range []int{0, 3} {
		source := newTestSource(numEntries)
		filterNone(source, &filterParams{})
		if len(source.Entries) != 0 {
			t.Errorf("with %d entries: got %d, want 0", numEntries, len(source.Entries))
		}
	}
}

func TestFilterFirst(t *testing.T) {
	tests := []struct {
		name       string
		numEntries int
		args       []string
		want       []string
	}{
		{"default keeps one", 3, nil, []string{"Entry 0"}},
		{"explicit count", 4, []string{"3"}, []string{"Entry 0", "Entry 1", "Entry 2"}},
		{"count larger than source", 2, []string{"10"}, []string{"Entry 0", "Entry 1"}},
		{"garbage arg falls back to one", 3, []string{"abc"}, []string{"Entry 0"}},
		{"zero arg falls back to one", 3, []string{"0"}, []string{"Entry 0"}},
		{"empty source", 0, nil, nil},
		{"empty source with count", 0, []string{"3"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := newTestSource(tt.numEntries)
			filterFirst(source, &filterParams{args: tt.args})
			assertTitles(t, source, tt.want...)
		})
	}
}

func TestFilterLast(t *testing.T) {
	tests := []struct {
		name       string
		numEntries int
		args       []string
		want       []string
	}{
		{"default keeps one", 3, nil, []string{"Entry 2"}},
		{"explicit count", 4, []string{"2"}, []string{"Entry 2", "Entry 3"}},
		{"count larger than source", 2, []string{"10"}, []string{"Entry 0", "Entry 1"}},
		{"empty source", 0, nil, nil},
		{"empty source with count", 0, []string{"2"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := newTestSource(tt.numEntries)
			filterLast(source, &filterParams{args: tt.args})
			assertTitles(t, source, tt.want...)
		})
	}
}

func TestFilterRandom(t *testing.T) {
	tests := []struct {
		name       string
		numEntries int
		args       []string
		wantCount  int
	}{
		{"default keeps one", 5, nil, 1},
		{"explicit count", 5, []string{"3"}, 3},
		{"count equals source size", 4, []string{"4"}, 4},
		{"empty source", 0, nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := newTestSource(tt.numEntries)
			originalUIDs := make(map[string]bool)
			for _, entry := range source.Entries {
				originalUIDs[entry.UID] = true
			}
			filterRandom(source, &filterParams{args: tt.args})
			if len(source.Entries) != tt.wantCount {
				t.Fatalf("got %d entries, want %d", len(source.Entries), tt.wantCount)
			}
			seen := make(map[string]bool)
			for _, entry := range source.Entries {
				if !originalUIDs[entry.UID] {
					t.Errorf("entry %q was not in the original source", entry.UID)
				}
				if seen[entry.UID] {
					t.Errorf("entry %q appears twice", entry.UID)
				}
				seen[entry.UID] = true
			}
		})
	}
}

func TestFilterReverse(t *testing.T) {
	tests := []struct {
		name       string
		numEntries int
		want       []string
	}{
		{"even number of entries", 4, []string{"Entry 3", "Entry 2", "Entry 1", "Entry 0"}},
		{"odd number of entries", 3, []string{"Entry 2", "Entry 1", "Entry 0"}},
		{"single entry", 1, []string{"Entry 0"}},
		{"empty source", 0, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := newTestSource(tt.numEntries)
			filterReverse(source, &filterParams{})
			assertTitles(t, source, tt.want...)
		})
	}
}

func TestFilterToday(t *testing.T) {
	now := time.Now()
	source := newTestSource(4)
	source.Entries[0].Date = now
	source.Entries[1].Date = now.AddDate(0, 0, -1)
	source.Entries[2].Date = now
	source.Entries[3].Date = now.AddDate(-1, 0, 0)
	filterToday(source, &filterParams{})
	assertTitles(t, source, "Entry 0", "Entry 2")

	empty := newTestSource(0)
	filterToday(empty, &filterParams{})
	assertTitles(t, empty)
}

func TestFilterLastHour(t *testing.T) {
	now := time.Now()
	newSource := func() *goeland.Source {
		source := newTestSource(4)
		source.Entries[0].Date = now.Add(-1 * time.Hour)
		source.Entries[1].Date = now.Add(-30 * time.Hour)
		source.Entries[2].Date = now.Add(-5 * time.Minute)
		source.Entries[3].Date = now.Add(-72 * time.Hour)
		return source
	}

	t.Run("default 24 hours", func(t *testing.T) {
		source := newSource()
		filterLastHour(source, &filterParams{})
		assertTitles(t, source, "Entry 0", "Entry 2")
	})

	t.Run("custom range", func(t *testing.T) {
		source := newSource()
		filterLastHour(source, &filterParams{args: []string{"48"}})
		assertTitles(t, source, "Entry 0", "Entry 1", "Entry 2")
	})

	t.Run("empty source", func(t *testing.T) {
		source := newTestSource(0)
		filterLastHour(source, &filterParams{})
		assertTitles(t, source)
	})
}

func TestFilterDigest(t *testing.T) {
	t.Run("merges entries into one", func(t *testing.T) {
		source := newTestSource(3)
		filterDigest(source, &filterParams{})
		if len(source.Entries) != 1 {
			t.Fatalf("got %d entries, want 1", len(source.Entries))
		}
		digest := source.Entries[0]
		if digest.Title != "Digest for Test Source" {
			t.Errorf("got title %q, want %q", digest.Title, "Digest for Test Source")
		}
		if len(digest.UID) != 64 {
			t.Errorf("UID should be a sha256 hex string, got %q", digest.UID)
		}
		for i := 0; i < 3; i++ {
			heading := fmt.Sprintf("<h2>Entry %d</h2>", i)
			if !strings.Contains(digest.Content, heading) {
				t.Errorf("digest content is missing %q", heading)
			}
			body := fmt.Sprintf("Content of entry %d", i)
			if !strings.Contains(digest.Content, body) {
				t.Errorf("digest content is missing %q", body)
			}
		}
	})

	t.Run("custom heading level", func(t *testing.T) {
		source := newTestSource(2)
		filterDigest(source, &filterParams{args: []string{"3"}})
		if !strings.Contains(source.Entries[0].Content, "<h3>Entry 0</h3>") {
			t.Errorf("digest content should use h3 headings, got: %s", source.Entries[0].Content)
		}
	})

	t.Run("entries with include link become anchors", func(t *testing.T) {
		source := newTestSource(2)
		filterIncludeLink(source, &filterParams{})
		filterDigest(source, &filterParams{})
		want := `<h2><a href="https://example.com/articles/0">Entry 0</a></h2>`
		if !strings.Contains(source.Entries[0].Content, want) {
			t.Errorf("digest content is missing %q, got: %s", want, source.Entries[0].Content)
		}
	})

	t.Run("entries with include source title get a source heading", func(t *testing.T) {
		source := newTestSource(3)
		filterIncludeSourceTitle(source, &filterParams{})
		filterDigest(source, &filterParams{})
		want := `<h1><a href="https://example.com/feed.xml">Test Source</a></h1>`
		if got := strings.Count(source.Entries[0].Content, want); got != 1 {
			t.Errorf("source heading should appear exactly once, found %d in: %s", got, source.Entries[0].Content)
		}
	})

	t.Run("empty source stays empty", func(t *testing.T) {
		source := newTestSource(0)
		filterDigest(source, &filterParams{})
		assertTitles(t, source)
	})
}

func TestFilterCombine(t *testing.T) {
	t.Run("uses first entry title", func(t *testing.T) {
		source := newTestSource(3)
		filterCombine(source, &filterParams{})
		if len(source.Entries) != 1 {
			t.Fatalf("got %d entries, want 1", len(source.Entries))
		}
		if source.Entries[0].Title != "Entry 0" {
			t.Errorf("got title %q, want %q", source.Entries[0].Title, "Entry 0")
		}
	})

	t.Run("empty source stays empty", func(t *testing.T) {
		source := newTestSource(0)
		filterCombine(source, &filterParams{})
		assertTitles(t, source)
	})
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{"strips path query and fragment", "https://example.com/path/page?q=1#frag", "https://example.com", false},
		{"bare host", "https://example.com", "https://example.com", false},
		{"keeps port", "http://example.com:8080/feed", "http://example.com:8080", false},
		{"invalid url", "://missing-scheme", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBaseURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFilterRelativeLinks(t *testing.T) {
	source := newTestSource(2)
	source.Entries[0].Content = `<img src="//cdn.example.com/pic.png"><a href='//example.org/page'>x</a>`
	source.Entries[1].Content = `<a href="/local/page">y</a><img src="/img/logo.png">`
	filterRelativeLinks(source, &filterParams{})
	want0 := `<img src="https://cdn.example.com/pic.png"><a href='https://example.org/page'>x</a>`
	if source.Entries[0].Content != want0 {
		t.Errorf("protocol-relative links:\ngot  %s\nwant %s", source.Entries[0].Content, want0)
	}
	want1 := `<a href="https://example.com/local/page">y</a><img src="https://example.com/img/logo.png">`
	if source.Entries[1].Content != want1 {
		t.Errorf("root-relative links:\ngot  %s\nwant %s", source.Entries[1].Content, want1)
	}

	empty := newTestSource(0)
	filterRelativeLinks(empty, &filterParams{})
	assertTitles(t, empty)
}

func TestFilterReplace(t *testing.T) {
	config := viper.New()
	config.Set("replace.myreplace.from", "entry")
	config.Set("replace.myreplace.to", "article")

	source := newTestSource(2)
	filterReplace(source, &filterParams{args: []string{"myreplace"}, config: config})
	for i, entry := range source.Entries {
		if strings.Contains(entry.Content, "entry") {
			t.Errorf("entry %d still contains %q: %s", i, "entry", entry.Content)
		}
		if !strings.Contains(entry.Content, "article") {
			t.Errorf("entry %d is missing %q: %s", i, "article", entry.Content)
		}
	}

	empty := newTestSource(0)
	filterReplace(empty, &filterParams{args: []string{"myreplace"}, config: config})
	assertTitles(t, empty)
}

func TestFilterSanitize(t *testing.T) {
	source := newTestSource(1)
	source.Entries[0].Content = `<p>Hello</p><script>alert("evil")</script>`
	filterSanitize(source, &filterParams{})
	if strings.Contains(source.Entries[0].Content, "script") {
		t.Errorf("script tag survived sanitization: %s", source.Entries[0].Content)
	}
	if !strings.Contains(source.Entries[0].Content, "<p>Hello</p>") {
		t.Errorf("benign content was removed: %s", source.Entries[0].Content)
	}

	empty := newTestSource(0)
	filterSanitize(empty, &filterParams{})
	assertTitles(t, empty)
}

func TestFilterIncludeLinkAndSourceTitle(t *testing.T) {
	source := newTestSource(2)
	filterIncludeLink(source, &filterParams{})
	filterIncludeSourceTitle(source, &filterParams{})
	for i, entry := range source.Entries {
		if !entry.IncludeLink {
			t.Errorf("entry %d: IncludeLink not set", i)
		}
		if !entry.IncludeSourceTitle {
			t.Errorf("entry %d: IncludeSourceTitle not set", i)
		}
	}

	empty := newTestSource(0)
	filterIncludeLink(empty, &filterParams{})
	filterIncludeSourceTitle(empty, &filterParams{})
	assertTitles(t, empty)
}

func TestFilterEmbedImage(t *testing.T) {
	imageTag := `<img src="https://example.com/images/0.png"`

	t.Run("defaults to top", func(t *testing.T) {
		source := newTestSource(1)
		filterEmbedImage(source, &filterParams{})
		content := source.Entries[0].Content
		if !strings.HasPrefix(content, imageTag) {
			t.Errorf("image should be prepended, got: %s", content)
		}
		if !strings.Contains(content, `class="top"`) {
			t.Errorf("image should have class top, got: %s", content)
		}
	})

	t.Run("bottom position appends", func(t *testing.T) {
		source := newTestSource(1)
		filterEmbedImage(source, &filterParams{args: []string{"bottom"}})
		content := source.Entries[0].Content
		if !strings.HasSuffix(content, `class="bottom">`) {
			t.Errorf("image should be appended, got: %s", content)
		}
	})

	t.Run("left position adds clear break", func(t *testing.T) {
		source := newTestSource(1)
		filterEmbedImage(source, &filterParams{args: []string{"left"}})
		content := source.Entries[0].Content
		if !strings.HasPrefix(content, imageTag) || !strings.HasSuffix(content, `<br style="clear:both" />`) {
			t.Errorf("left position should prepend image and append a clearing break, got: %s", content)
		}
	})

	t.Run("link argument wraps image in anchor", func(t *testing.T) {
		source := newTestSource(1)
		filterEmbedImage(source, &filterParams{args: []string{"top", "link"}})
		want := `<a href="https://example.com/articles/0"><img src="https://example.com/images/0.png" class="top"></a>`
		if !strings.HasPrefix(source.Entries[0].Content, want) {
			t.Errorf("image should be wrapped in a link, got: %s", source.Entries[0].Content)
		}
	})

	t.Run("entry without image is untouched", func(t *testing.T) {
		source := newTestSource(1)
		source.Entries[0].ImageURL = "  "
		before := source.Entries[0].Content
		filterEmbedImage(source, &filterParams{})
		if source.Entries[0].Content != before {
			t.Errorf("content changed for an entry without image: %s", source.Entries[0].Content)
		}
	})

	t.Run("empty source", func(t *testing.T) {
		source := newTestSource(0)
		filterEmbedImage(source, &filterParams{})
		assertTitles(t, source)
	})
}

func TestFilterToc(t *testing.T) {
	t.Run("prepends a table of content", func(t *testing.T) {
		source := newTestSource(2)
		filterToc(source, &filterParams{})
		if len(source.Entries) != 3 {
			t.Fatalf("got %d entries, want 3", len(source.Entries))
		}
		toc := source.Entries[0]
		if toc.Title != "Table of Content for Test Source" {
			t.Errorf("got title %q", toc.Title)
		}
		if !strings.Contains(toc.Content, "<li>Entry 0</li>") || !strings.Contains(toc.Content, "<li>Entry 1</li>") {
			t.Errorf("toc content is missing entries: %s", toc.Content)
		}
		if len(toc.UID) != 64 {
			t.Errorf("UID should be a sha256 hex string, got %q", toc.UID)
		}
	})

	t.Run("title argument links the source", func(t *testing.T) {
		source := newTestSource(1)
		filterToc(source, &filterParams{args: []string{"title"}})
		want := `<a href="https://example.com/feed.xml">Test Source</a>`
		if source.Entries[0].Title != want {
			t.Errorf("got title %q, want %q", source.Entries[0].Title, want)
		}
	})

	t.Run("entries with include link become anchors", func(t *testing.T) {
		source := newTestSource(1)
		filterIncludeLink(source, &filterParams{})
		filterToc(source, &filterParams{})
		want := `<li><a href="https://example.com/articles/0">Entry 0</a></li>`
		if !strings.Contains(source.Entries[0].Content, want) {
			t.Errorf("toc content is missing %q: %s", want, source.Entries[0].Content)
		}
	})

	t.Run("empty source stays empty", func(t *testing.T) {
		source := newTestSource(0)
		filterToc(source, &filterParams{})
		assertTitles(t, source)
	})
}

func TestTruncateWordsToWholeSentence(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		length        int
		want          string
		wantTruncated bool
	}{
		{"cuts at sentence end", "One two three. Four five six. Seven.", 2, "One two three.", true},
		{"question mark ends sentence", "Is this it? More words follow here.", 2, "Is this it?", true},
		{"shorter than limit", "Just a few words.", 100, "Just a few words.", false},
		{"no sentence end after limit", "one two three four five", 2, "one two three four five", false},
		{"empty string", "", 5, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, truncated := truncateWordsToWholeSentence(tt.input, tt.length)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if truncated != tt.wantTruncated {
				t.Errorf("got truncated=%v, want %v", truncated, tt.wantTruncated)
			}
		})
	}
}

func TestFilterLimitWords(t *testing.T) {
	t.Run("limits content length", func(t *testing.T) {
		source := newTestSource(1)
		source.Entries[0].Content = "One two three. Four five six seven eight."
		filterLimitWords(source, &filterParams{args: []string{"2"}})
		if source.Entries[0].Content != "One two three." {
			t.Errorf("got %q", source.Entries[0].Content)
		}
	})

	t.Run("no argument leaves content unchanged", func(t *testing.T) {
		source := newTestSource(1)
		before := source.Entries[0].Content
		filterLimitWords(source, &filterParams{})
		if source.Entries[0].Content != before {
			t.Errorf("content changed without a word limit: %q", source.Entries[0].Content)
		}
	})

	t.Run("empty source", func(t *testing.T) {
		source := newTestSource(0)
		filterLimitWords(source, &filterParams{args: []string{"2"}})
		assertTitles(t, source)
	})
}

func TestFilterRESkip(t *testing.T) {
	tests := []struct {
		name       string
		numEntries int
		args       []string
		want       []string
	}{
		{"skips matching titles", 3, []string{"Entry [02]"}, []string{"Entry 1"}},
		{"no match keeps everything", 3, []string{"^Nothing$"}, []string{"Entry 0", "Entry 1", "Entry 2"}},
		{"invalid regex keeps everything", 2, []string{"("}, []string{"Entry 0", "Entry 1"}},
		{"missing argument keeps everything", 2, nil, []string{"Entry 0", "Entry 1"}},
		{"too many arguments keeps everything", 2, []string{"a", "b"}, []string{"Entry 0", "Entry 1"}},
		{"empty source", 0, []string{".*"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := newTestSource(tt.numEntries)
			filterRESkip(source, &filterParams{args: tt.args})
			assertTitles(t, source, tt.want...)
		})
	}
}

func TestFilterUntrack(t *testing.T) {
	source := newTestSource(1)
	source.Entries[0].Content = `<p>Hello</p><a href="http://feeds.feedburner.com/~ff/x">track</a><img src="http://feeds.feedburner.com/~r/x/><p>Bye</p>`
	filterUntrack(source, &filterParams{})
	if strings.Contains(source.Entries[0].Content, "feedburner") {
		t.Errorf("feedburner tracking survived: %s", source.Entries[0].Content)
	}
	if !strings.Contains(source.Entries[0].Content, "<p>Hello</p>") {
		t.Errorf("benign content was removed: %s", source.Entries[0].Content)
	}

	empty := newTestSource(0)
	filterUntrack(empty, &filterParams{})
	assertTitles(t, empty)
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name string
		a    string
		list []string
		want bool
	}{
		{"present", "en", []string{"en", "fr"}, true},
		{"case insensitive", "EN", []string{"en", "fr"}, true},
		{"absent", "de", []string{"en", "fr"}, false},
		{"empty list", "en", nil, false},
		{"empty string", "", []string{"en"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.a, tt.list); got != tt.want {
				t.Errorf("StringInSlice(%q, %v) = %v, want %v", tt.a, tt.list, got, tt.want)
			}
		})
	}
}

func TestFilterLanguage(t *testing.T) {
	newSource := func() *goeland.Source {
		source := newTestSource(2)
		source.Entries[0].Title = "English entry"
		source.Entries[0].Content = "<p>This is a wonderful story about the sea and the many birds that fly over it every single morning.</p>"
		source.Entries[1].Title = "French entry"
		source.Entries[1].Content = "<p>Ceci est une merveilleuse histoire qui parle de la mer et des nombreux oiseaux qui la survolent chaque matin.</p>"
		return source
	}

	t.Run("keeps only requested language", func(t *testing.T) {
		source := newSource()
		filterLanguage(source, &filterParams{args: []string{"en"}})
		assertTitles(t, source, "English entry")
	})

	t.Run("keeps several languages", func(t *testing.T) {
		source := newSource()
		filterLanguage(source, &filterParams{args: []string{"en", "fr"}})
		assertTitles(t, source, "English entry", "French entry")
	})

	t.Run("no language drops everything", func(t *testing.T) {
		source := newSource()
		filterLanguage(source, &filterParams{})
		assertTitles(t, source)
	})

	t.Run("empty source", func(t *testing.T) {
		source := newTestSource(0)
		filterLanguage(source, &filterParams{args: []string{"en"}})
		assertTitles(t, source)
	})
}

func TestFilterReddit(t *testing.T) {
	t.Run("sanitizes regular post content", func(t *testing.T) {
		source := newTestSource(1)
		source.Entries[0].URL = "https://www.reddit.com/r/golang/comments/abc123/some_post/"
		source.Entries[0].Content = `<p>Hello</p><script>alert("evil")</script>`
		filterReddit(source, &filterParams{})
		if strings.Contains(source.Entries[0].Content, "script") {
			t.Errorf("script tag survived sanitization: %s", source.Entries[0].Content)
		}
		if !strings.Contains(source.Entries[0].Content, "<p>Hello</p>") {
			t.Errorf("benign content was removed: %s", source.Entries[0].Content)
		}
	})

	t.Run("entry without post id is left alone", func(t *testing.T) {
		source := newTestSource(1)
		source.Entries[0].URL = "https://www.reddit.com/r/golang/"
		before := source.Entries[0].Content
		filterReddit(source, &filterParams{})
		if source.Entries[0].Content != before {
			t.Errorf("content changed for entry without post ID: %s", source.Entries[0].Content)
		}
	})
}

func TestPurgeUnseen(t *testing.T) {
	config := viper.New()
	config.Set("database", filepath.Join(t.TempDir(), "goeland-test.db"))

	source := newTestSource(2)
	filterUnSeen(source, &filterParams{config: config})
	assertTitles(t, source, "Entry 0", "Entry 1")

	// A large retention window purges nothing: the entries stay seen.
	if err := PurgeUnseen(config, "test", 15); err != nil {
		t.Fatalf("PurgeUnseen: %v", err)
	}
	source = newTestSource(2)
	filterUnSeen(source, &filterParams{config: config})
	assertTitles(t, source)

	// A zero-day window purges all seen records: the entries are unseen again.
	if err := PurgeUnseen(config, "test", 0); err != nil {
		t.Fatalf("PurgeUnseen: %v", err)
	}
	source = newTestSource(2)
	filterUnSeen(source, &filterParams{config: config})
	assertTitles(t, source, "Entry 0", "Entry 1")
}

func TestFilterUnSeen(t *testing.T) {
	config := viper.New()
	config.Set("database", filepath.Join(t.TempDir(), "goeland-test.db"))

	source := newTestSource(3)
	filterUnSeen(source, &filterParams{config: config})
	assertTitles(t, source, "Entry 0", "Entry 1", "Entry 2")

	// The same entries plus a new one: only the new one is unseen.
	source = newTestSource(4)
	filterUnSeen(source, &filterParams{config: config})
	assertTitles(t, source, "Entry 3")

	empty := newTestSource(0)
	filterUnSeen(empty, &filterParams{config: config})
	assertTitles(t, empty)

	// An unopenable database (a directory) keeps the entries untouched.
	badConfig := viper.New()
	badConfig.Set("database", t.TempDir())
	source = newTestSource(2)
	filterUnSeen(source, &filterParams{config: badConfig})
	assertTitles(t, source, "Entry 0", "Entry 1")
}

func TestFilterSource(t *testing.T) {
	newConfig := func(filters ...string) *viper.Viper {
		config := viper.New()
		config.Set("sources.test.filters", filters)
		return config
	}

	t.Run("applies filters in order", func(t *testing.T) {
		source := newTestSource(5)
		FilterSource(source, newConfig("first(4)", "last(2)"))
		assertTitles(t, source, "Entry 2", "Entry 3")
	})

	t.Run("trims spaces in arguments", func(t *testing.T) {
		source := newTestSource(5)
		FilterSource(source, newConfig("first( 2 )"))
		assertTitles(t, source, "Entry 0", "Entry 1")
	})

	t.Run("multiple arguments", func(t *testing.T) {
		source := newTestSource(1)
		FilterSource(source, newConfig("embedimage(top, link)"))
		if !strings.HasPrefix(source.Entries[0].Content, `<a href="https://example.com/articles/0">`) {
			t.Errorf("embedimage(top, link) not applied: %s", source.Entries[0].Content)
		}
	})

	t.Run("unknown filter is ignored", func(t *testing.T) {
		source := newTestSource(3)
		FilterSource(source, newConfig("doesnotexist", "first(2)"))
		assertTitles(t, source, "Entry 0", "Entry 1")
	})

	t.Run("no filters configured", func(t *testing.T) {
		source := newTestSource(3)
		FilterSource(source, viper.New())
		assertTitles(t, source, "Entry 0", "Entry 1", "Entry 2")
	})

	t.Run("empty source", func(t *testing.T) {
		source := newTestSource(0)
		FilterSource(source, newConfig("first", "digest", "toc"))
		assertTitles(t, source)
	})
}

// TestAllFiltersHandleEmptySource runs every registered filter against an
// empty source to make sure none of them panics or invents entries.
// Filters that produce entries from existing ones (digest, toc, ...) must
// leave an empty source empty.
func TestAllFiltersHandleEmptySource(t *testing.T) {
	for name, f := range filters {
		t.Run(name, func(t *testing.T) {
			config := viper.New()
			config.Set("database", filepath.Join(t.TempDir(), "goeland-test.db"))
			source := newTestSource(0)
			f.filterFunc(source, &filterParams{args: nil, config: config})
			if len(source.Entries) != 0 {
				t.Errorf("filter %q produced %d entries from an empty source", name, len(source.Entries))
			}
		})
	}
}
