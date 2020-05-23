# goeland

![goeland](asset/goeland_small.png)

![version](https://img.shields.io/github/go-mod/go-version/slurdge/goeland)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/slurdge/goeland)
[![Build Status](https://travis-ci.com/slurdge/goeland.svg?branch=master)](https://travis-ci.com/slurdge/goeland)
![GitHub](https://img.shields.io/github/license/slurdge/goeland)
![Image license](https://img.shields.io/badge/Images-CC%20BY--SA%204.0-blueviolet)

A RSS to email, ala rss2email written in Go.
Take one or more RSS feeds and transform them into a proper email format.

The available filters are as follow:

- none: Removes all entries
- lasthours: Keep only the entries that are from the X last hours (default 24)
- unseen: Keep only unseen entry
- last: Keep only the last entry
- digest: Make a digest of all entries (optional heading level, default is )
- links: Rewrite relative links src="// and href="// to have an https:// prefix
- replace: Replace a string with another. Use with an argument like this: replace(myreplace) and define
```
[replace.myreplace]
        from="A string"
        to="Another string"
```
in your config file.
- includelink: Include the link of entries in the digest form
- lebrief: Retrieves the full excerpts for Next INpact's Lebrief
- first: Keep only the first entry
- random: Keep 1 or more random entries. Use either 'random' or 'random(5)' for example.
- today: Keep only the entries for today
- language: Keep only the specified languages (best effort detection), use like this: language(en,de)
- all: Default, include all entries
- reverse: Reverse the order of the entries
- combine: Combine all the entries into one source and use the first entry title as source title. Useful for merge sources

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

Linux:
```console
$ make
$ ./bin/goeland
```
All platforms:
```console
go run main.go
```

### Testing

No tests are written yet, contributions are welcomed!

```console
make test
```