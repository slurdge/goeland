Changelog
=========

v0.6.0
------

- Add 'embedimage' filter.

v0.5.3
------

- Fix bug with NextInpact's lebrief filter.

v0.5.2
------

- Fix bug with Reddit RSS

v0.5.1
------

- Fix Email titles being escaped twice

v0.5.0
------

- True HTML messages with templates
- New filter: `lasthours` or `lasthours(X)` which allows to keep only entries that have less than X (default: 24) hours date.
- Add support for go generate.
- Assets are now put in binary form from `asset/` folder.

v0.4.4
-------

- Actual fix for v0.4.3

v0.4.3
------

- Fix the imgur clientID not being correctly provisionned

v0.4.2
------

- Fix the sub sources not properly being filtered

v0.4.1
------

- Fix the purge command

v0.4.0
------

This version adds the 'unseen' filters, which allow to filter entries that have not been already seen.
The key is the source name and entry UID. If you change your source name, it will invalidate your cache.

- New command: purge. goeland purge will remove all old (+15 days) entries of the database
- New config value: database. Use it to override the location of the database. Default is goeland.db.

v0.3.2
------

- arm64 releases are built

v0.3.1
------

- Bugfix for v0.3.0: ClientID was not included in build

v0.3.0
------

- More filters: random
- Imgur source (tag only)

v0.2.0
------

- Lists are now treated as such in the toml file
- More filters, do goeland help run to have a list
- Support for dates in digests
- Cleaner architecture for filters

v0.1.4
------

First public version, with basic functionality and filters