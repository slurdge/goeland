# Changelog

## [0.21.0](https://github.com/slurdge/goeland/compare/v0.20.1...v0.21.0) (2025-12-18)


### Features

* setup go 1.25 ([ffbd81b](https://github.com/slurdge/goeland/commit/ffbd81b6ba7d07dfcdb197515070d3cb47970b68))
* update go to 1.25 ([7915b24](https://github.com/slurdge/goeland/commit/7915b248ed61a45a68ab1c727f8c3495922059a5))

## [0.20.1](https://github.com/slurdge/goeland/compare/v0.20.0...v0.20.1) (2025-06-02)


### Bug Fixes

* [#237](https://github.com/slurdge/goeland/issues/237) deprecated option in minifier ([996552d](https://github.com/slurdge/goeland/commit/996552d8eece0410686b9a633bddd01de8efe246))

## [0.20.0](https://github.com/slurdge/goeland/compare/v0.19.0...v0.20.0) (2025-06-01)


### Features

* allow to change the css (default, sakura, water) ([e5e5075](https://github.com/slurdge/goeland/commit/e5e507510e6328f03b23d59667623ee7757574c1))

## [0.19.0](https://github.com/slurdge/goeland/compare/v0.18.3...v0.19.0) (2025-04-03)


### Features

* reskip filter ([36f0e71](https://github.com/slurdge/goeland/commit/36f0e71c45440d3560e279a503cb644804c203aa))

## [0.18.3](https://github.com/slurdge/goeland/compare/v0.18.2...v0.18.3) (2024-01-06)


### Bug Fixes

* second try to fixup clientID missing ([5c778dc](https://github.com/slurdge/goeland/commit/5c778dc829778e9688de920ffa8c41e9339b08ed))

## [0.18.2](https://github.com/slurdge/goeland/compare/v0.18.1...v0.18.2) (2024-01-05)


### Bug Fixes

* imgur is not fecthing, probably missing client-id ([c7c1380](https://github.com/slurdge/goeland/commit/c7c13800b20726b343d0dcb3a19b5b79294b64b3))

## [0.18.1](https://github.com/slurdge/goeland/compare/v0.18.0...v0.18.1) (2023-12-18)


### Bug Fixes

* cannot build darwin docker images ([35c844c](https://github.com/slurdge/goeland/commit/35c844cf1444fe19ff13cf1963ac54d42e3dffe9))

## [0.18.0](https://github.com/slurdge/goeland/compare/v0.17.0...v0.18.0) (2023-12-18)


### Features

* add darwin/arm64 as a target ([c2ee684](https://github.com/slurdge/goeland/commit/c2ee6846125876d87f2280793416c373a35e8606))
* make includeLink works even with single articles ([7bf3d6d](https://github.com/slurdge/goeland/commit/7bf3d6df51ede9f5d4897bf21f8bcfa158788f06))

## [0.17.0](https://github.com/slurdge/goeland/compare/v0.16.0...v0.17.0) (2023-12-17)


### Features

* allow the user to set the user-agent ([1743092](https://github.com/slurdge/goeland/commit/17430921e2f632bf9dfddec07ece3a4cfd6a2b13))
* upgrade release-please to v4 ([81ed4c8](https://github.com/slurdge/goeland/commit/81ed4c8a6def28208a6cd0b1a361132781569489))

## [0.16.0](https://github.com/slurdge/goeland/compare/v0.15.0...v0.16.0) (2023-10-27)


### Features

* better links filter for relative URLs ([8bfa732](https://github.com/slurdge/goeland/commit/8bfa7321a21e62dfaea52458a1308575fd1686e7))


### Bug Fixes

* disable UPX as it seems to produce more problems than it solves ([6ce1de4](https://github.com/slurdge/goeland/commit/6ce1de478861170f23c551e124a66d216f6366ab))

## [0.15.0](https://github.com/slurdge/goeland/compare/v0.14.0...v0.15.0) (2023-03-26)


### Features

* Add fields to the Entry struct ([1b6106b](https://github.com/slurdge/goeland/commit/1b6106b93fd93af0e94b698adcb8fec58d7dd808))
* Add includesourcetitle filter ([6ba56f5](https://github.com/slurdge/goeland/commit/6ba56f5c5e479d123aa4c20ca62c42f634341d39))
* Set the source in the entry when fetching ([92d1f83](https://github.com/slurdge/goeland/commit/92d1f83be7fb162c7b30345a89912682d5e8f538))

## [0.14.0](https://github.com/slurdge/goeland/compare/v0.13.0...v0.14.0) (2023-03-16)


### Features

* add a second arg to image, 'link', which put a link to the entry associated with the image. ([eef64c6](https://github.com/slurdge/goeland/commit/eef64c6a142129cf7130c7d5825baf593156de74))
* get information from media:group if main info is empty. Should make Youtube work. Fix [#122](https://github.com/slurdge/goeland/issues/122) ([2c24d5c](https://github.com/slurdge/goeland/commit/2c24d5c404b5acd333e3cf11a604ce73803ca5e5))
* try to have a v in releases ([c147b9c](https://github.com/slurdge/goeland/commit/c147b9ccc0cf56b0a6d407e90b109129cc301ac4))


### Bug Fixes

* correctly parse change log as formatted by release-please. fixes [#126](https://github.com/slurdge/goeland/issues/126) ([0b6a9d7](https://github.com/slurdge/goeland/commit/0b6a9d7f9ee44f664d20381f0516aab012bd3ea2))
* python invocation for local build ([fbdb0b3](https://github.com/slurdge/goeland/commit/fbdb0b3d05363b2283f6d516f557cec107e22ede))

## [v0.13.0](https://github.com/slurdge/goeland/compare/v0.12.3...v0.13.0) (2023-01-18)


### Features

* allow insecure feed location ([71091b6](https://github.com/slurdge/goeland/commit/71091b66a9fab0ea6767c87f91845096d859f3ea))


### Bug Fixes

* prevent merges to have cycles and infinite recursion fixes [#100](https://github.com/slurdge/goeland/issues/100) ([8b10526](https://github.com/slurdge/goeland/commit/8b10526f95e652e46c16e8eaf53cf3181c475834))

## [v0.12.3](https://github.com/slurdge/goeland/compare/v0.12.2...v0.12.3) (2022-11-16)


### Bug Fixes

* another try for releases to trigger tag building ([3ff35c9](https://github.com/slurdge/goeland/commit/3ff35c92695d1fa8404ed74db37ec774cd0d2a9b))

## [v0.12.2](https://github.com/slurdge/goeland/compare/v0.12.1...v0.12.2) (2022-11-16)


### Bug Fixes

* try to fix workflow run with release please ([933c982](https://github.com/slurdge/goeland/commit/933c982edda74882a4c90f3234987e8b4eb16aff))

## [v0.12.1](https://github.com/slurdge/goeland/compare/v0.12.0...v0.12.1) (2022-11-15)


### Bug Fixes

* default to SSLTLS only on port 465 ([1bbe563](https://github.com/slurdge/goeland/commit/1bbe563c0492c298dc3ac054ffc0c5832608fc67))

## [v0.12.0](https://github.com/slurdge/goeland/compare/v0.11.0...v0.12.0) (2022-11-13)


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
