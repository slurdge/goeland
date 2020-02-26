package cmd

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
	"github.com/slurdge/indigo/config"
	"github.com/slurdge/indigo/internal/indigo"
	"github.com/slurdge/indigo/log"
	"github.com/spf13/cobra"
	"jaytaylor.com/html2text"
)

const tmpURL = "https://www.nextinpact.com/rss/lebrief.xml"

func createEmailPool(config config.Provider) (*email.Pool, error) {
	host := config.GetString("email.host")
	port := config.GetInt("email.port")
	user := config.GetString("email.username")
	pass := config.GetString("email.password")
	auth := smtp.PlainAuth("", user, pass, host)
	return email.NewPool(fmt.Sprintf("%s:%v", host, port), 8, auth)
}

func run(cmd *cobra.Command, args []string) {
	log.Debugln("Running...")
	config := config.Config()

	getSubString := func(root string, key string, tail string) string {
		return config.GetString(fmt.Sprintf("%s.%s.%s", root, key, tail))
	}

	dryRun := config.GetBool("dry-run")

	var pool *email.Pool
	var err error

	pipes := config.GetStringMapString("pipes")
	for pipe := range pipes {
		log.Infof("Executing pipe named: %v", pipe)
		sourceName := getSubString("pipes", pipe, "source")
		destinationName := getSubString("pipes", pipe, "destination")
		if !config.IsSet(fmt.Sprintf("sources.%s", sourceName)) {
			log.Errorf("cannot find source: %s", sourceName)
			continue
		}
		if !config.IsSet(fmt.Sprintf("destinations.%s", destinationName)) {
			log.Errorf("cannot find destination: %s", destinationName)
			continue
		}
		source, _ := indigo.GetSource(config, sourceName)
		log.Infof("Retrieved %v feeds", len(source.Entries))
		filters := strings.Split(getSubString("pipes", pipe, "filters"), ",")
		for i := range filters {
			filters[i] = strings.TrimSpace(filters[i])
		}
		for _, filter := range filters {
			switch filter {
			case "all":
				source.FilterAll()
			case "today":
				source.FilterToday()
			case "lebrief":
				source.FilterLeBrief()
			case "digest":
				source.FilterDigest()
			default:
				log.Errorf("unknown filter: %s\n", filter)
			}
			log.Infof("After %s: %v feeds", filter, len(source.Entries))
		}
		if dryRun {
			log.Infoln("Dry run has been specified, not outputting...")
			continue
		}
		if getSubString("destinations", destinationName, "type") == "email" {
			if pool == nil {
				pool, err = createEmailPool(config)
				if err != nil {
					log.Errorf("cannot create email pool: %v")
					continue
				}
			}
			for _, entry := range source.Entries {
				email := email.NewEmail()
				email.From = getSubString("destinations", destinationName, "email_from")
				email.To = []string{getSubString("destinations", destinationName, "email_to")}
				email.Subject = entry.Title
				html := entry.Content
				text, _ := html2text.FromString(html)
				email.Text = []byte(text)
				email.HTML = []byte(html)
				err := pool.Send(email, -1)
				if err != nil {
					fmt.Errorf("%v", err)
				}
			}
		} else {
			fmt.Printf("**%s**\n", source.Title)
			for _, entry := range source.Entries {
				text, _ := html2text.FromString(entry.Content, html2text.Options{})
				fmt.Printf("*%s*\n%s\n%s", entry.Title, entry.Date, text)
			}
		}
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
	config.Config().BindPFlag("dry-run", runCmd.Flags().Lookup("dry-run"))
	rootCmd.AddCommand(runCmd)
}
