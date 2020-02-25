package cmd

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/slurdge/indigo/config"
	"github.com/slurdge/indigo/log"
	"github.com/spf13/cobra"
)

const tmpURL = "https://www.nextinpact.com/rss/lebrief.xml"

// Entry This represent an entry produced by a source
type Entry struct {
	UID     string
	Title   string
	Content string
	Date    time.Time
}

// Source ...
type Source struct {
	Title   string
	Entries []Entry
}

func fetchFeed(feedLocation string, isFile bool) (Source, error) {
	var source Source
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	var err error
	if isFile {
		//todo : handle error
		file, err := os.Open(feedLocation)
		if err != nil {
			return source, fmt.Errorf("cannot open: %s", feedLocation)
		}
		defer file.Close()
		feed, err = fp.Parse(file)
		if err != nil {
			return source, fmt.Errorf("cannot parse file: %s", feedLocation)
		}
	} else {
		feed, err = fp.ParseURL(feedLocation)
		if err != nil {
			return source, fmt.Errorf("cannot open or parse url: %s", feedLocation)
		}
	}
	for _, item := range feed.Items {
		entry := Entry{}
		entry.Title = item.Title
		entry.Content = item.Description
		fmt.Println(entry.Title, item.Description)
		entry.UID = item.GUID
		entry.Date = *item.PublishedParsed
		source.Entries = append(source.Entries, entry)
	}
	source.Title = feed.Title
	return source, nil
}

func filterAll(source Source) Source {
	return source
}

func filterNone(source Source) Source {
	source.Entries = nil
	return source
}

func filterToday(source Source) Source {
	var current int
	for _, entry := range source.Entries {
		if entry.Date.Day() != time.Now().Day() {
			continue
		}
		source.Entries[current] = entry
		current++
	}
	source.Entries = source.Entries[:current]
	return source
}

func filterDigest(source Source) Source {
	digest := Entry{}
	digest.Title = fmt.Sprintf("Digest for %s", source.Title)
	content := ""
	for _, entry := range source.Entries {
		content += fmt.Sprintf("<h1>%s</h1><br>", entry.Title)
		content += entry.Content
	}
	h := sha256.New()
	h.Write([]byte(content))
	digest.UID = fmt.Sprintf("%x", h.Sum(nil))
	digest.Date = time.Now()
	digest.Content = content
	source.Entries = []Entry{digest}
	return source
}

func run(cmd *cobra.Command, args []string) {
	log.Debugln("Running...")
	defaultConfig := config.Config()

	// todo: remove this line
	fmt.Println(defaultConfig.GetBool("main.dry-run"))

	feeds := defaultConfig.GetStringMapString("sources")
	for feedname := range feeds {
		filename := defaultConfig.GetString(fmt.Sprintf("sources.%s.url", feedname))
		content, err := fetchFeed(filename, true)
		if err != nil {
			log.Errorf("Cannot retrieve feed: %s error: %v", feedname, err)
		}
		log.Infof("Retrieved %v feeds\n", len(content.Entries))
		content = filterToday(content)
		content = filterDigest(content)
		log.Infof("Today we have %v feeds\n", len(content.Entries))
		log.Infof("%v", content)
	}

}

// versionCmd represents the version command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch the RSS and emails it",
	Run:   run,
}

func init() {
	runCmd.Flags().Bool("dry-run", false, "Do only a dry-run")
	config.Config().BindPFlag("main.dry-run", runCmd.Flags().Lookup("dry-run"))
	rootCmd.AddCommand(runCmd)
}
