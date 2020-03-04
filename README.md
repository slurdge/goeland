# goeland 📧

A RSS to email, ala rss2email written in Go.
Take one or more RSS feeds and transform them into a proper email format.

The available filters are as follow:

- all: Default, include all entries
- last: Keep only the last entry
- links: Rewrite relative links src="// and href="// to have an https:// prefix
- replace: Replace a string with another. Use with an argument like this: replace(myreplace) and define
        [replace.myreplace]
        from="A string"
        to="Another string"
    in your config file.
- none: Removes all entries
- first: Keep only the first entry
- reverse: Reverse the order of the entries
- today: Keep only the entries for today
- digest: Make a digest of all entries (optional heading level, default is 1)
- combine: Combine all the entries into one source and use the first entry title as source title. Useful for merge sources
- language: Keep only the specified languages (best effort detection), use like this: language(en,de)
- lebrief: Retrieves the full excerpts for Next INpact's Lebrief

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