package goeland

import (
	"time"
)

// Entry This represent an entry produced by a source
type Entry struct {
	UID         string
	Title       string
	Content     string
	URL         string
	Date        time.Time
	IncludeLink bool
}

// Source ...
type Source struct {
	Name    string
	Title   string
	Entries []Entry
}
