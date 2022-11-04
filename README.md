# goeland

![goeland](cmd/asset/goeland_small.png)

![GitHub release (latest by date)](https://img.shields.io/github/v/release/slurdge/goeland)
![version](https://img.shields.io/github/go-mod/go-version/slurdge/goeland)
[![Build Status](https://github.com/slurdge/goeland/actions/workflows/build.yml/badge.svg)](https://github.com/slurdge/goeland/actions/workflows/build.yml)
![GitHub](https://img.shields.io/github/license/slurdge/goeland)
![Image license](https://img.shields.io/badge/Images-CC%20BY--SA%204.0-blueviolet)
[![Docker images](https://github.com/slurdge/goeland/actions/workflows/docker.yml/badge.svg)](https://github.com/slurdge/goeland/actions/workflows/docker.yml)

A RSS to email, ala rss2email written in Go.

Support this project by giving it a ⭐️ and sharing it.

## About

Goeland excels at creating beautiful emails from RSS, tailored for daily or weekly digest.

It include a number of filters (see below) that can transform the RSS content along the way. It can also consume other sources, such as a Imgur tag.

Goeland transforms this...

```xml
<rss version="2.0">
<channel>
<title>Phoronix</title>
<link>https://www.phoronix.com/</link>
<description>
Linux Hardware Reviews, Benchmarks & Open-Source News
</description>
<language>en-us</language>
<item>
<title>
Google Announces KataOS As Security-Focused OS, Leveraging Rust & seL4 Microkernel
</title>
<link>https://www.phoronix.com/news/Google-KataOS</link>
<guid>https://www.phoronix.com/news/Google-KataOS</guid>
<description>
Google this week has announced the release of KataOS as their newest operating system effort focused on embedded devices running ambient machine learning workloads. KataOS is security-minded, exclusively uses the Rust programming language, and is built atop the seL4 microkernel as its foundation...
</description>
<pubDate>Sun, 16 Oct 2022 06:10:25 -0400</pubDate>
</item>
</rss>
```

into this

![email](documentation/screenshots/phoronix_mixed.png)

Goeland has a size-fits-all default template that works well with mobile, tablet, desktop and webmail clients.

Goeland can extract full text from most articles sources, enabling a ready to consume email.

## Status

Goeland is used in production with many email clients, and has sent over thousands of email. It is considered stable.

## Installation

Grab the latest binary release from the [release page](https://github.com/slurdge/goeland/releases/latest/).
Binaries are available for the following platforms:

* linux/386
* linux/amd64
* linux/arm
* linux/arm64
* darwin/amd64
* windows/amd64
* windows/386

Just put it in a folder where you have write permissions and run it first with :

```console
goeland run
```

If you are interested in another platform supported, please open a PR or submit a feature request.

## Usage

On first run, if it doesn't exist yet, goeland will create a `config.toml` with the default values. You need to adjust the `[email]` section with your SMTP server details.

### Sources

Afterwards, fill the `[sources]` and `[pipes]` sections.
Source are identified by their name after the `[source.]` field:

```toml
[sources.hackernews]
type = "feed"
url = "https://hnrss.org/newest"
filters = ["all", "today"]
```

You can then use `'hackernews'` in the following pipes.

The different source types are:

* `"feed"`: Regular RSS feed. Fill in the `url` field
* `"imgur"`: Return most recent results for a tag. Fill in the the `tag` field.
* `"merge"`: Will merge two or more sources together. Fill in the `sources` field with a list of sources: `sources = ["source1", "source2"]`. Especially useful to merge different sources on the same topic. Don't forget to `digest` or `combine` it later.

### Filtering

One powerful aspect of goeland is filtering. Instead of sending the content of the RSS directly to the email system, it can transform it in a number of ways in order to make it easier to read, process, etc.

Any number of filters can be defined, the order is important. For example, the following:

```toml
filters = ["unseen", "lebrief", "digest"]
```

Will first keep only previously `unseen` entries, then transform it nicer with `lebrief` filter, and, at last, will put them all together with `digest`. This will create only one email with a SourceTitle as the title of the RSS feed.

The available filters are as follow:

* none: Removes all entries
* all: Default, include all entries

* first: Keep only the first entry
* last: Keep only the last entry
* reverse: Reverse the order of the entries
* random: Keep 1 or more random entries. Use either 'random' or 'random(5)' for example.
* unseen: Keep only unseen entry. Entries that have been seen will be put in a `goeland.db` file. Use the `purge` command to remove seen entries.
* today: Keep only the entries for today
* lasthours: Keep only the entries that are from the X last hours (default 24)
* digest: Make a digest of all entries (optional heading level, default is 2)
* combine: Combine all the entries into one source and use the first entry title as source title. Useful for merge sources
* links: Rewrite relative links src="// and href="// to have an https:// prefix
* embedimage: Embed a picture if the entry has an attachment with a type of picture (optional position: top|bottom|left|right, default is top)
* replace: Replace a string with another. Use with an argument like this: replace(myreplace) and define

```toml
[replace.myreplace]
        from="A string"
        to="Another string"
```

in your config file.

* includelink: Include the link of entries in the digest form
* lebrief: Retrieves the full excerpts for Next INpact's Lebrief. Use only with a source from Next INpact.
* language: Keep only the specified languages (best effort detection), use like this: `language(en,de)`

### Pipes

After defining a number of sources, you can send them to a pipe. One source can be send to multiple pipes, but a pipe can only have one source. If you need to combine sources together, use the above special `merge` type to have this effect.

A pipe has the following structure:

```toml
[pipes.hackernews]
source = "hackernews"
destination = "email"
email_to = "example@example.com"
email_from = "HackerNews <goeland@example.com>"
email_title = "{{.EntryTitle}}"
```

You can use EntryTitle, SourceTitle and SourceName in the email template. SourceTitle is the title of the RSS stream.

For debug purposes, or in order to pipe in other systems, you can set destination to `terminal`.

### Email

In the email section you need to specify your outgoing mail server. From 0.8.0, you can specify both `encryption` and `allow-insecure` to connect to self hosted servers. You can also specify `authentication` to select the appropriate option for your server ( the options available are `"none"`, `"plain"`, `"login"` and `"crammd5"`; if unspecified it defaults to `"plain"`; see [`go-simple-mail`](https://pkg.go.dev/github.com/xhit/go-simple-mail/v2#AuthType)'s documentation for details).

```toml
[email]
host = "smtp.example.com"
port = 25
username = "default"
password = "p4ssw0rd"
encryption = "tls"
allow-insecure = false
authentication = "plain"
#Email customization
include-header = true
include-footer = true
include-toc = false
include-content = true
#footer = Your custom footer
#logo = internal:goeland.png
#template = /path/to/template.html
```

You can create your own template, see [relevant documentation](documentation/templates.md)

## Examples

This will bring you 6 puppies to your inbox.

```toml
loglevel = "info"
dry-run = false

[email]
host = "smtp.sendgrid.net"
port = 587
username = "apikey"
password = "<sendgridapikey>"

[sources]

[sources.insta]
url = "https://rssbridge.example.com/?action=display&bridge=Instagram&context=Hashtag&h=puppy&media_type=picture&direct_links=on&format=MRss"
type = "feed"
filters = ["random(3)"]

[sources.imgur]
type = "imgur"
tag = "puppy"
filters = ["random(3)"]

[sources.puppies]
type = "merge"
sources = ["insta", "imgur"]
filters = ["combine"]

[pipes]

[pipes.puppies]
source = "puppies"
destination = "email"
email_to = ["puppylover@example.com"]
email_from = "DailyPuppy <goeland@example.com>"
```

This will give you the latest article on a specific subreddit:

```toml
loglevel = "none"
dry-run = false
database = "goeland.db"

[email]
host = "example.com"
port = 25
username = "username"
password = "password"

[sources]

[sources.reddit]
url = "https://www.reddit.com/r/selfhosted/top.rss"
type = "feed"
filters = ["unseen", "includelink", "digest"]

[pipes.reddit]
source = "reddit"
destination = "email"
email_to = ["example@example.com"]
email_from = "Reddit <reddit@example.com>"
```

It is possible to send an email to multiple addresses, just put them in a list:

```toml
[pipes.reddit]
source = "reddit"
destination = "email"
email_to = ["bob@example.com", "alice@gmail.com", "charles@yahoo.com"]
email_from = "Reddit <reddit@example.com>"
```

See also the `examples/` folder.

## Contributing

Feel free to open bugs or PR for more sources, more filters suggestions.

If you encounter a problematic feed, please open a bug with the content of the feed attached.

## Future

Here is a list of things that could be nice

* image inliner
* embedded scripting language for filters&manipulation
* remove tags for instagram
* footer text
* use enclosure of the feed as header image
* <https://github.com/go-shiori/go-readability>
