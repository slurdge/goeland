package cmd

import (
	"bytes"
	"crypto/tls"
	"embed"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

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
	email "github.com/xhit/go-simple-mail/v2"
	"jaytaylor.com/html2text"
)

const logoAttachmentName = "logo.png"

//go:embed asset/email.default.html
var defaultEmailBytes []byte

//go:embed asset/default.css
var defaultCSS string

//go:embed asset/*
var asset embed.FS

//go:embed asset/goeland@250w.png
var logoBytes []byte

func createEmailTemplate(config config.Provider, pipe string) (*template.Template, error) {
	minifier := minify.New()
	minifier.Add("text/html", &mhtml.Minifier{
		KeepSpecialComments: true,
	})

	emailBytes := defaultEmailBytes

	var templateFilename string
	if len(pipe) > 0 {
		templateFilename = config.GetString(fmt.Sprintf("pipes.%s.template", pipe))
	}
	if len(templateFilename) == 0 {
		templateFilename = config.GetString("email.template")
	}
	if len(templateFilename) > 0 {
		var err error
		emailBytes, err = os.ReadFile(templateFilename)
		if err != nil {
			return nil, err
		}
	}

	minified, err := minifier.Bytes("text/html", emailBytes)
	if err != nil {
		return nil, err
	}

	tpl := template.Must(template.New("email").Parse(string(minified)))
	return tpl, nil
}

func createEmailPool(config config.Provider) (*email.SMTPClient, error) {
	host := config.GetString("email.host")
	port := config.GetInt("email.port")
	if port == 0 {
		port = 587
	}
	user := config.GetString("email.username")
	pass := config.GetString("email.password")
	passFile := config.GetString("email.password_file")
	if len(passFile) > 0 {
		if len(pass) > 0 {
			log.Warn("Both password and password_file are set. Using password_file.")
		}
		passFileContent, err := os.ReadFile(passFile)
		if err != nil {
			return nil, fmt.Errorf("error while reading password file: %v", err)
		}
		pass = string(passFileContent)
	}
	server := email.NewSMTPClient()
	authentications := map[string]email.AuthType{"none": email.AuthNone, "plain": email.AuthPlain, "login": email.AuthLogin, "crammd5": email.AuthCRAMMD5}
	authentication, found := authentications[config.GetString("email.authentication")]
	if !found {
		authentication = email.AuthPlain
	}
	server.Authentication = authentication
	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = pass
	encryptions := map[string]email.Encryption{"none": email.EncryptionNone, "tls": email.EncryptionSTARTTLS, "ssl": email.EncryptionSSLTLS}
	encryption, found := encryptions[config.GetString("email.encryption")]
	if !found {
		encryption = email.EncryptionSTARTTLS
		if port == 465 {
			encryption = email.EncryptionSSLTLS
		}
	}
	server.Encryption = encryption
	if config.GetBool("email.allow-insecure") {
		server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	server.KeepAlive = true
	emailTimeout := time.Duration(config.GetInt64("email.timeout-ms") * 1000 * 1000)
	server.ConnectTimeout = emailTimeout
	server.SendTimeout = emailTimeout
	smtpClient, err := server.Connect()

	if err != nil {
		return nil, err
	}
	return smtpClient, nil
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
func formatHTMLEmail(entry *goeland.Entry, config config.Provider, tpl *template.Template, destination string) string {
	footer := strings.TrimSpace(config.GetString("email.footer"))
	if footer == "" {
		footer = footers[rand.Intn(len(footers))]
	}
	cssContent := defaultCSS
	css := config.GetString("email.css")
	var cssMapping = map[string]string{
		"simple":      "simple.min.css",
		"sakura":      "sakura.css",
		"sakura-dark": "sakura-dark.css",
		"water":       "water.min.css",
	}
	filename, ok := cssMapping[strings.ToLower(css)]
	if ok {
		data, _ := asset.ReadFile(fmt.Sprintf("asset/%s", filename))
		cssContent = string(data)
	} else if css != "" && css != "default" {
		log.Warnf("Unknown css : %s, using default...", css)
	}
	data := struct {
		EntryTitle    string
		EntryContent  string
		EntryURL      string
		IncludeHeader bool
		IncludeTitle  bool
		IncludeFooter bool
		IncludeLink   bool
		EntryFooter   string
		ContentID     string
		CSS           string
	}{
		EntryTitle:    html.EscapeString(entry.Title),
		EntryContent:  entry.Content,
		EntryURL:      entry.URL,
		IncludeHeader: config.GetBool("email.include-header"),
		IncludeTitle:  config.GetBool("email.include-title"),
		IncludeFooter: config.GetBool("email.include-footer"),
		IncludeLink:   entry.IncludeLink,
		EntryFooter:   footer,
		ContentID:     "cid:" + logoAttachmentName,
		CSS:           cssContent,
	}
	if destination == "htmlfile" {
		data.ContentID = "data:image/png;base64," + base64.StdEncoding.EncodeToString(logoBytes)
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
func inlineImage(e *email.Email, r io.Reader, filename string, c string) (err error) {
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, r); err != nil {
		return err
	}
	at := &email.File{
		Name:   filename,
		Inline: true,
		Data:   buffer.Bytes(),
	}
	if c != "" {
		at.MimeType = c
	}
	e.Attach(at)
	if e.Error != nil {
		return e.Error
	}
	return nil
}

func stringInSlice(str string, slice []string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}

func run(cmd *cobra.Command, args []string) {
	log.Debugln("Running...")
	config := viper.GetViper()

	getSubString := func(root string, key string, tail string) string {
		return config.GetString(fmt.Sprintf("%s.%s.%s", root, key, tail))
	}

	dryRun := config.GetBool("dry-run")

	var pool *email.SMTPClient

	mainTemplate, err := createEmailTemplate(config, "")
	if err != nil {
		log.Fatalf("cannot create email template: %v", err)
	}

	logoFilename := config.GetString("email.logo")
	if logoFilename != "internal:goeland.png" {
		logoBytes, err = os.ReadFile(logoFilename)
		if err != nil {
			log.Fatalf("cannot read email logo file: %v", err)
		}
	}
	pipes := config.GetStringMapString("pipes")
	for pipe := range pipes {
		tpl := mainTemplate
		if len(getSubString("pipes", pipe, "template")) > 0 {
			tpl, err = createEmailTemplate(config, pipe)
			if err != nil {
				log.Fatalf("cannot create email template: %v", err)
			}
		}
		if len(args) != 0 && !stringInSlice(pipe, args) {
			continue
		}
		disabled := config.GetBool(fmt.Sprintf("pipes.%s.disabled", pipe))
		if disabled {
			log.Infof("Skipping disabled pipe: %s", pipe)
			continue
		}
		log.Infof("Executing pipe named: %s", pipe)
		sourceName := getSubString("pipes", pipe, "source")
		source, err := fetch.FetchSource(config, sourceName, []string{})
		if err != nil {
			log.Errorf("Error getting source: %s: %v", sourceName, err)
			continue
		}
		if dryRun {
			log.Infoln("Dry run has been specified, not outputting...")
			continue
		}
		destination := getSubString("pipes", pipe, "destination")
		switch destination {
		case "email":
			if pool == nil {
				pool, err = createEmailPool(config)
				if err != nil {
					log.Fatalf("cannot create email pool: %v", err)
				}
				if pool == nil {
					log.Errorf("cannot send email: no pool created")
				}
			}
			for _, entry := range source.Entries {
				message := email.NewMSG()
				message.SetFrom(getSubString("pipes", pipe, "email_from"))
				message.AddTo(config.GetStringSlice(fmt.Sprintf("pipes.%s.email_to", pipe))...)
				message.AddCc(config.GetStringSlice(fmt.Sprintf("pipes.%s.email_cc", pipe))...)
				message.AddBcc(config.GetStringSlice(fmt.Sprintf("pipes.%s.email_bcc", pipe))...)
				templateString := getSubString("pipes", pipe, "email_title")
				subject := formatEmailSubject(source, &entry, templateString)
				message.SetSubject(subject)
				entry.Title = subject
				if config.GetBool("email.include-header") {
					err := inlineImage(message, bytes.NewReader(logoBytes), logoAttachmentName, "image/png")
					if err != nil {
						log.Errorf("error attaching logo: %v", err)
					}
				}
				html := formatHTMLEmail(&entry, config, tpl, destination)
				text, err := html2text.FromString(entry.Content)
				if err != nil {
					text = "There was an error converting HTML content to text"
				}
				message.SetBody(email.TextPlain, text)
				message.AddAlternative(email.TextHTML, html)
				err = message.Send(pool)
				if err != nil {
					log.Errorf("error sending email: %v", err)
				}
			}
		case "htmlfile":
			for i, entry := range source.Entries {
				html := formatHTMLEmail(&entry, config, tpl, destination)
				var HTMLFile *os.File
				err := os.MkdirAll("data", os.ModePerm)
				if err != nil {
					fatalErr(fmt.Errorf("error while creating the data directory to put the html file in: %v", err))
				}
				if HTMLFile, err = os.Create(fmt.Sprintf(filepath.Join("data", "%s - %d.html"), pipe, i)); err != nil {
					fatalErr(fmt.Errorf("error while writing html file: %v", err))
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
	viper.BindPFlag("dry-run", runCmd.Flags().Lookup("dry-run"))
	runCmd.Flags().String("logo", "internal:goeland.png", "Override the logo file")
	viper.BindPFlag("email.logo", runCmd.Flags().Lookup("logo"))
	runCmd.Flags().String("footer", "", "Override the default footer")
	viper.BindPFlag("email.footer", runCmd.Flags().Lookup("footer"))
	runCmd.Flags().String("css", "default", "Override the embedded CSS (default, pico, sakura, sakura-dark)")
	viper.BindPFlag("email.css", runCmd.Flags().Lookup("css"))
	runCmd.Flags().Bool("unsafe-no-sanitize-filter", false, "Do not sanitize inputs. ⚠ Use at your own risk!")
	viper.BindPFlag("unsafe-no-sanitize-filter", runCmd.Flags().Lookup("unsafe-no-sanitize-filter"))
	bindFlags(runCmd, viper.GetViper())
	rootCmd.AddCommand(runCmd)
}
