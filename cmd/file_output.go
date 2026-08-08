package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/inbucket/html2text"
	"github.com/slurdge/goeland/config"
	"github.com/slurdge/goeland/internal/goeland"
	"github.com/slurdge/goeland/log"
)

const defaultHTMLFilePath = "data"
const defaultHTMLFileName = "{{.Pipe}} - {{.EntryNumber}}.html"

func getHTMLFileOption(config config.Provider, pipe string, key string, fallback string) string {
	value := config.GetString(fmt.Sprintf("pipes.%s.htmlfile_%s", pipe, key))
	if value == "" {
		value = config.GetString(fmt.Sprintf("htmlfile.%s", key))
	}
	if value == "" {
		value = fallback
	}
	return value
}

func formatHTMLFileName(templateString string, pipe string, source *goeland.Source, index int) (string, error) {
	entry := &source.Entries[index]
	data := struct {
		Pipe        string
		SourceName  string
		SourceTitle string
		EntryNumber int
		EntryTitle  string
		EntryUID    string
		EntryDate   time.Time
		Today       time.Time
	}{
		Pipe:        pipe,
		SourceName:  source.Name,
		SourceTitle: source.Title,
		EntryNumber: index,
		EntryTitle:  entry.Title,
		EntryUID:    entry.UID,
		EntryDate:   entry.Date,
		Today:       time.Now(),
	}
	tpl, err := template.New("htmlfile_name").Parse(templateString)
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err := tpl.Execute(&output, data); err != nil {
		return "", err
	}
	return output.String(), nil
}

func sanitizeFileName(name string) string {
	var builder strings.Builder
	for _, r := range name {
		switch {
		case r < 0x20 || r == 0x7f:
		case strings.ContainsRune(`/\:*?"<>|`, r):
			builder.WriteRune('-')
		default:
			builder.WriteRune(r)
		}
	}
	name = strings.TrimRight(builder.String(), " .")
	const maxNameLen = 128
	if utf8.RuneCountInString(name) > maxNameLen {
		ext := filepath.Ext(name)
		base := []rune(strings.TrimSuffix(name, ext))
		keep := maxNameLen - utf8.RuneCountInString(ext)
		if keep < 1 {
			keep = 1
		}
		if len(base) > keep {
			base = base[:keep]
		}
		name = strings.TrimRight(string(base), " .") + ext
	}
	return name
}

func writeHTMLFiles(config config.Provider, pipe string, source *goeland.Source, tpl *template.Template) {
	dir := getHTMLFileOption(config, pipe, "path", defaultHTMLFilePath)
	nameTemplate := getHTMLFileOption(config, pipe, "filename", defaultHTMLFileName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Errorf("cannot create the html output directory '%s': %v", dir, err)
		return
	}
	written := 0
	for i, entry := range source.Entries {
		html := formatHTMLEmail(&entry, config, tpl, "htmlfile")
		fileName, err := formatHTMLFileName(nameTemplate, pipe, source, i)
		if err != nil {
			log.Errorf("cannot execute html filename template '%s': %v, using default", nameTemplate, err)
			fileName, err = formatHTMLFileName(defaultHTMLFileName, pipe, source, i)
			if err != nil {
				log.Fatalf("cannot execute default html filename template :%v", err)
			}
		}
		fileName = sanitizeFileName(fileName)
		if fileName == "" {
			log.Warningf("empty filename for source: %s entry: %s", source.Name, entry.Title)
			fileName = fmt.Sprintf("entry-%d.html", i)
		}
		filePath := filepath.Join(dir, fileName)
		if err := os.WriteFile(filePath, []byte(html), 0644); err != nil {
			log.Errorf("cannot write html file '%s': %v", filePath, err)
			continue
		}
		log.Debugf("Wrote %s", filePath)
		written++
	}
	log.Infof("Wrote %d/%d entries to %s for pipe %s", written, len(source.Entries), dir, pipe)
}

func printToConsole(source *goeland.Source) {
	fmt.Printf("**%s**\n", source.Title)
	for _, entry := range source.Entries {
		text, _ := html2text.FromString(entry.Content, html2text.Options{})
		fmt.Printf("*%s*\n%s\n%s\n", entry.Title, entry.Date, text)
	}
}
