package indigo

import "strings"

const headerToRemoveFR = "Article labellisé du jour "

// FilterWikipedia ...
func (source *Source) FilterWikipedia() {
	for i, entry := range source.Entries {
		entry.Content = strings.ReplaceAll(entry.Content, headerToRemoveFR, "")
		source.Entries[i] = entry
	}
}
