package goeland

import "strings"

const headerToRemoveFR = "Article labellisé du jour "

// FilterWikipedia ...
func filterWikipedia(source *Source) {
	for i, entry := range source.Entries {
		entry.Content = strings.ReplaceAll(entry.Content, headerToRemoveFR, "")
		source.Entries[i] = entry
	}
}
