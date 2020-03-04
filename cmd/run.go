package cmd

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jaytaylor.com/html2text"
)

const tmpURL = "https://www.nextinpact.com/rss/lebrief.xml"

func createEmailPool(config config.Provider) (*email.Pool, error) {
	host := config.GetString("email.host")
	port := config.GetInt("email.port")
	if port == 0 {
		port = 587
	}
	fullhost := fmt.Sprintf("%s:%v", host, port)
	user := config.GetString("email.username")
	pass := config.GetString("email.password")
	auth := smtp.PlainAuth("", user, pass, host)
	return email.NewPool(fullhost, 8, auth)
}

func run(cmd *cobra.Command, args []string) {
	log.Debugln("Running...")
	config := viper.GetViper()

	emailTimeoutInMs := time.Duration(config.GetInt64("email_timeout_ms"))

	getSubString := func(root string, key string, tail string) string {
		return config.GetString(fmt.Sprintf("%s.%s.%s", root, key, tail))
	}

	dryRun := config.GetBool("dry-run")

	var pool *email.Pool
	pipes := config.GetStringMapString("pipes")
	for pipe := range pipes {
		disabled := config.GetBool(fmt.Sprintf("pipes.%s.disabled", pipe))
		if disabled {
			log.Infof("Skipping disabled pipe: %s", pipe)
			continue
		}
		log.Infof("Executing pipe named: %s", pipe)
		sourceName := getSubString("pipes", pipe, "source")
		source, err := goeland.GetSource(config, sourceName)
		if err != nil {
			log.Errorf("Error getting source: %s", sourceName)
		}
		if dryRun {
			log.Infoln("Dry run has been specified, not outputting...")
			continue
		}
		if getSubString("pipes", pipe, "destination") == "email" {
			if pool == nil {
				pool, err = createEmailPool(config)
				if err != nil {
					log.Errorf("cannot create email pool: %v", err)
					continue
				}
			}
			for _, entry := range source.Entries {
				email := email.NewEmail()
				email.From = getSubString("pipes", pipe, "email_from")
				email.To = config.GetStringSlice(fmt.Sprintf("pipes.%s.email_to", pipe))
				data := struct {
					EntryTitle  string
					SourceTitle string
					SourceName  string
				}{EntryTitle: entry.Title, SourceTitle: source.Title, SourceName: source.Name}
				var output bytes.Buffer
				templateString := getSubString("pipes", pipe, "email_title")
				if strings.TrimSpace(templateString) == "" {
					templateString = `{{.EntryTitle}}`
				}
				tpl := template.Must(template.New("title").Parse(templateString))
				tpl.Execute(&output, data)
				email.Subject = output.String()
				html := entry.Content
				text, err := html2text.FromString(html)
				if err == nil {
					email.Text = []byte(text)
				} else {
					email.Text = []byte("There was an error converting HTML content to text")
				}
				email.HTML = []byte(html)
				err = pool.Send(email, emailTimeoutInMs*1000*1000)
				if err != nil {
					log.Errorf("error sending email: %v", err)
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
	Long: strings.Join([]string{
		`Take one or more RSS feeds and transform them into a proper email format.
		
The available filters are as follow:`,
		goeland.GetFiltersHelp(),
	}, "\n"),
	Run: run,
}

func init() {
	runCmd.Flags().Bool("dry-run", false, "Do not output anything, just fetch and filter the content")
	viper.GetViper().BindPFlag("dry-run", runCmd.Flags().Lookup("dry-run"))
	rootCmd.AddCommand(runCmd)
}
