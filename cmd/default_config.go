package cmd

// DefaultConfig contains the configuration that will be output if none is present
const DefaultConfig string = `loglevel = "none"
dry-run = false

[email]
host = "smtp.exmaple.com"
port = 587
username = "user"
password = "pass"

[sources]

[sources.hackernews]
url = "https://hnrss.org/newest"
type = "feed"
# See doc for available filters
filters = ""

[pipes]

[pipes.hackernews]
#Either put disabled = true or prefix pipes with disabled like this: disabled.pipes.hackernews
disabled = true
source = "hackernews"
destination = "email"
email_to = "example@example.com"
email_from = "HackerNews <goeland@example.com>"
#Default: you can use EntryTitle, SourceTitle and SourceName in the template
#email_title = "{{.EntryTitle}}"
`
