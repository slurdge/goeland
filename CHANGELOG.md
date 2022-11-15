# Changelog

## [0.12.1](https://github.com/slurdge/goeland/compare/v0.12.0...v0.12.1) (2022-11-15)


### Bug Fixes

* default to SSLTLS only on port 465 ([1bbe563](https://github.com/slurdge/goeland/commit/1bbe563c0492c298dc3ac054ffc0c5832608fc67))

## [0.12.0](https://github.com/slurdge/goeland/compare/v0.11.0...v0.12.0) (2022-11-13)


### Features

* add a way to have toc & link only digest with toc filter. fixes [#79](https://github.com/slurdge/goeland/issues/79) ([037a7c4](https://github.com/slurdge/goeland/commit/037a7c427a742566c155c67c72c5ba577fc46d59))
* add a way to limit the number of words in the entries. fixes [#74](https://github.com/slurdge/goeland/issues/74) ([59ad04f](https://github.com/slurdge/goeland/commit/59ad04fbb73fa1ee40f1c9caa4e51528e8537376))
* add basic i18n ([9e4d3e1](https://github.com/slurdge/goeland/commit/9e4d3e1b49823361e4322a6b6274f977f4b3de53))
* add password_file option ([e9b3322](https://github.com/slurdge/goeland/commit/e9b332274d533232b8acd65f05228505b97c2fa5))
* add Release please ([1078e19](https://github.com/slurdge/goeland/commit/1078e194faab8c968220c15317c8ed9103298d2a))
* allow defining pipe templates ([d5888f0](https://github.com/slurdge/goeland/commit/d5888f05bb37f93d8f67b6eb9746966db80cce92))
* get config from standard places ([9471b13](https://github.com/slurdge/goeland/commit/9471b137abf587b20b4d42f018965060a9b48933))
* logo base64 png output for debug html files ([9fb6010](https://github.com/slurdge/goeland/commit/9fb6010db88425b6e21a76f7522f2d66f765e9fd))
* scaffolding for supporting other css/theme, tested with sakura ([b2ce9fc](https://github.com/slurdge/goeland/commit/b2ce9fc028934b10ba142be7acec56a982b00e69))


### Bug Fixes

* default encryption set twice ([bf91552](https://github.com/slurdge/goeland/commit/bf915521a22fad436d5ce864b28ac139f6513fed))
* default to SSLTLS encryption when port is 465 ([fa44112](https://github.com/slurdge/goeland/commit/fa44112602c0962a2a5e65d05c60f43c97bc5b26))
* fix branch name for release ([ae48309](https://github.com/slurdge/goeland/commit/ae483092caad366c6c9532b80de1d506f7a24eec))
* put defaults values for ports related to encryption. never default to none ([068acd4](https://github.com/slurdge/goeland/commit/068acd41a3a3ef97aae332b27fbb17808aae18b6))
* wrong translations ([c305144](https://github.com/slurdge/goeland/commit/c305144257b98882baadff49caec667a60ae22eb))

Changelog
=========

v0.11.0
-------

- Feature: allow other email template (by @fabianofranz)
- Fixed environment variables not fetched (useful for Docker deployments)
- Feature: You can now cc and bcc recipients
- Feature: Sanitization can now be skipped, and added later through a filter.
- Bumped go to 1.18
- Updated dependencies

v0.10.3
--------

- Fixed a bug where empty UID feed wouldn't work with `unseen` filter.
- Updated dependencies

v0.10.2
-------

- Small release with Docker releases

v0.10.1
-------

- Fix by kylrth: fix pipes running only latest defined pipe in daemon mode.

v0.10.0
-------

- Add a `daemon` command for daemons. Also add a `--run-at-startup` flag for daemon
- Allow to override log level on the command line
- Allow to specify which pipes are running. Example: `goeland run <pipename>`
- Add a reddit filter that gets better pictures. Useful for picture intensive subreddits.
- Fix reddit rss download
- Create an email pool only if it's needed
- Various small fixes

v0.9.0
------

- Add the possibility to change authentication mechanism, thanks to dfosas

v0.8.1
------

- Fix GMail display by swapping HTML&Text
- Add footer override
- Add more nice footers

v0.8.0
------

- Add a filter to remove feedburner tracking
- Add some examples
- Replaced email library
- Allow to connect to TLS/SSL/Non secure servers
- Allow to override TLS/SSL security check
- Fixes email Titles containing escaped characters

v0.7.0
------

- Add a 'retrieve' filter. Retrieves the full articles from links. Use it like this: retrieve(#content). See goquery for queries.
- Bump go to 1.6.0
- Bump mod dependencies to latests versions
- Remove go generate and use the new go embed
- Add more footers
- Extract version from changelog
- First & Last filters now accept arguments

v0.6.1
------

- Fix #2: Not all embedded image are displayed by getting alternative sources for pictures.
- Fix a rendering bug with img classes.
- Sanitize HTML input with bluemonday.

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

- Fix the imgur clientID not being correctly provisioned

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
