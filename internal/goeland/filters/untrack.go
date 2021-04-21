// Copyright (c) 2021 slurdge <slurdge@slurdge.org>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package filters

import (
	"regexp"

	"github.com/slurdge/goeland/internal/goeland"
)

func filterUntrack(source *goeland.Source, params *filterParams) {
	for index, entry := range source.Entries {
		reA := regexp.MustCompile(`<a\s+href="http://feeds\.feedburner\.com/.*?</a>`)
		reIMG := regexp.MustCompile(`<img src="http://feeds\.feedburner\.com/.*?/>`)
		entry.Content = reA.ReplaceAllString(entry.Content, "")
		entry.Content = reIMG.ReplaceAllString(entry.Content, "")
		source.Entries[index] = entry
	}
}
