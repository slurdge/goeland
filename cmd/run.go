package cmd

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"net/smtp"
	"net/textproto"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/jordan-wright/email"
	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/internal/goeland/fetch"
	"github.com/slurdge/goeland/internal/goeland/filters"
	"github.com/slurdge/goeland/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify/v2"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/vanng822/go-premailer/premailer"
	"jaytaylor.com/html2text"
)

const logoAttachmentName = "logo.png"

func createEmailTemplate(_ config.Provider) (*template.Template, error) {
	minifier := minify.New()
	minifier.Add("text/html", &mhtml.Minifier{
		KeepConditionalComments: true,
	})
	emailBytes, err := Asset("email.default.html")
	if err != nil {
		return nil, err
	}
	minified, err := minifier.Bytes("text/html", emailBytes)
	if err != nil {
		return nil, err
	}
	tpl := template.Must(template.New("email").Parse(string(minified)))
	return tpl, nil
}

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

func formatEmailSubject(source *goeland.Source, entry *goeland.Entry, templateString string) string {
	data := struct {
		EntryTitle  string
		SourceTitle string
		SourceName  string
		Today       time.Time
	}{EntryTitle: entry.Title,
		SourceTitle: source.Title,
		SourceName:  source.Name,
		Today:       time.Now()}
	var output bytes.Buffer
	if strings.TrimSpace(templateString) == "" {
		templateString = `{{.EntryTitle}}`
	}
	tpl := template.Must(template.New("email_title").Parse(templateString))
	tpl.Execute(&output, data)
	return output.String()
}
func formatHTMLEmail(entry *goeland.Entry, config config.Provider, tpl *template.Template) string {
	data := struct {
		EntryTitle    string
		EntryContent  string
		IncludeHeader bool
		IncludeTitle  bool
		IncludeFooter bool
		ContentID     string
	}{EntryTitle: html.EscapeString(entry.Title),
		EntryContent:  entry.Content,
		IncludeHeader: config.GetBool("email.include-header"),
		IncludeTitle:  config.GetBool("email.include-title"),
		IncludeFooter: config.GetBool("email.include-footer"),
		ContentID:     logoAttachmentName,
	}
	var output bytes.Buffer
	tpl.Execute(&output, data)

	prem, err := premailer.NewPremailerFromString(output.String(), premailer.NewOptions())
	if err != nil {
		log.Errorf("cannot instantiate premailer: %v", err)
		return output.String()
	}
	html, err := prem.Transform()
	if err != nil {
		log.Errorf("cannot inline css: %v", err)
		return output.String()
	}
	return html
}
func inlineImage(e *email.Email, r io.Reader, filename string, c string) (a *email.Attachment, err error) {
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, r); err != nil {
		return
	}
	at := &email.Attachment{
		Filename: filename,
		Header:   textproto.MIMEHeader{},
		Content:  buffer.Bytes(),
	}
	if c != "" {
		at.Header.Set("Content-Type", c)
	} else {
		at.Header.Set("Content-Type", "application/octet-stream")
	}
	at.Header.Set("Content-Disposition", fmt.Sprintf("inline;\r\n filename=\"%s\"", filename))
	at.Header.Set("Content-ID", fmt.Sprintf("<%s>", filename))
	at.Header.Set("Content-Transfer-Encoding", "base64")
	e.Attachments = append(e.Attachments, at)
	return at, nil
}

func run(cmd *cobra.Command, args []string) {
	log.Debugln("Running...")
	config := viper.GetViper()

	emailTimeoutInNs := time.Duration(config.GetInt64("email.timeout-ms") * 1000 * 1000)

	getSubString := func(root string, key string, tail string) string {
		return config.GetString(fmt.Sprintf("%s.%s.%s", root, key, tail))
	}

	dryRun := config.GetBool("dry-run")

	pool, err := createEmailPool(config)
	if err != nil {
		log.Errorf("cannot create email pool: %v", err)
	}
	tpl, err := createEmailTemplate(config)
	if err != nil {
		log.Errorf("cannot create email template: %v", err)
	}
	var logoBytes []byte
	logoFilename := config.GetString("email.logo")
	if logoFilename == "internal:goeland.png" {
		logoBytes, err = Asset("goeland@250w.png")
		if err != nil {
			log.Errorf("cannot create email logo: %v", err)
		}
	} else {
		logoBytes, err = ioutil.ReadFile(logoFilename)
		if err != nil {
			log.Errorf("cannot read email logo file: %v", err)
		}
	}
	pipes := config.GetStringMapString("pipes")
	for pipe := range pipes {
		disabled := config.GetBool(fmt.Sprintf("pipes.%s.disabled", pipe))
		if disabled {
			log.Infof("Skipping disabled pipe: %s", pipe)
			continue
		}
		log.Infof("Executing pipe named: %s", pipe)
		sourceName := getSubString("pipes", pipe, "source")
		source, err := fetch.FetchSource(config, sourceName)
		if err != nil {
			log.Errorf("Error getting source: %s", sourceName)
		}
		if dryRun {
			log.Infoln("Dry run has been specified, not outputting...")
			continue
		}
		destination := getSubString("pipes", pipe, "destination")
		switch destination {
		case "email":
			if pool == nil {
				log.Errorf("cannot send email: no pool created")
			}
			for _, entry := range source.Entries {
				email := email.NewEmail()
				email.From = getSubString("pipes", pipe, "email_from")
				email.To = config.GetStringSlice(fmt.Sprintf("pipes.%s.email_to", pipe))
				templateString := getSubString("pipes", pipe, "email_title")
				email.Subject = formatEmailSubject(source, &entry, templateString)
				entry.Title = email.Subject
				text, err := html2text.FromString(entry.Content)
				if err == nil {
					email.Text = []byte(text)
				} else {
					email.Text = []byte("There was an error converting HTML content to text")
				}
				if config.GetBool("email.include-header") {
					//at, err := email.Attach(bytes.NewReader(logoBytes), logoAttachmentName, "image/png")
					at, err := inlineImage(email, bytes.NewReader(logoBytes), logoAttachmentName, "image/png")
					if err != nil {
						log.Errorf("error attaching logo: %v", err)
					}
					at.HTMLRelated = true
				}
				html := formatHTMLEmail(&entry, config, tpl)
				email.HTML = []byte(html)
				err = pool.Send(email, emailTimeoutInNs)
				if err != nil {
					log.Errorf("error sending email: %v", err)
				}
			}
		case "htmlfile":
			for i, entry := range source.Entries {
				html := formatHTMLEmail(&entry, config, tpl)
				var HTMLFile *os.File
				if HTMLFile, err = os.Create(fmt.Sprintf("%s - %d.html", pipe, i)); err != nil {
					fatalErr(fmt.Errorf("cannot open config.toml for writing"))
				}
				HTMLFile.Write([]byte(html))
			}
		case "console":
		case "terminal":
			fmt.Printf("**%s**\n", source.Title)
			for _, entry := range source.Entries {
				text, _ := html2text.FromString(entry.Content, html2text.Options{})
				fmt.Printf("*%s*\n%s\n%s\n", entry.Title, entry.Date, text)
			}
		case "null", "none":
		default:
			log.Infof("unknown destination type: %s", destination)
		}
	}
	if config.GetBool("auto-purge") {
		purge(nil, nil)
	}
}

// versionCmd represents the version command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Fetch the RSS and emails it",
	Long: strings.Join([]string{
		`Take one or more RSS feeds and transform them into a proper email format.
		
The available filters are as follow:`,
		filters.GetFiltersHelp(),
	}, "\n"),
	Run: run,
}

func init() {
	runCmd.Flags().Bool("dry-run", false, "Do not output anything, just fetch and filter the content")
	viper.GetViper().BindPFlag("dry-run", runCmd.Flags().Lookup("dry-run"))
	runCmd.Flags().String("logo", "internal:goeland.png", "Override the logo file")
	viper.GetViper().BindPFlag("email.logo", runCmd.Flags().Lookup("logo"))
	rootCmd.AddCommand(runCmd)
}
