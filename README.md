[![made-with-golang](https://img.shields.io/badge/Made%20with-Golang-blue.svg?style=flat-square)](https://golang.org/)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%203-blue.svg?style=flat-square)](https://github.com/l3uddz/nabarr/blob/master/LICENSE.md)
[![Discord](https://img.shields.io/discord/381077432285003776.svg?colorB=177DC1&label=Discord&style=flat-square)](https://discord.io/cloudbox)
[![Donate](https://img.shields.io/badge/Donate-gray.svg?style=flat-square)](#donate)

# Nabarr

Nabarr monitors Newznab/Torznab RSS feeds to find new media to add to Sonarr and or Radarr.

## Table of contents

- [Installing nabarr](#installing-nabarr)
- [Introduction](#introduction)
    - [Media](#media)
    - [PVR](#pvr)
    - [RSS](#rss)
    - [Full config file](#full-config-file)
- [Other installation options](#other-installation-options)
    - [Docker](#docker)
- [Donate](#donate)

## Installing nabarr

Nabarr offers [pre-compiled binaries](https://github.com/l3uddz/nabarr/releases/latest) for Linux, MacOS and Windows for each official release. In addition, there is also a [Docker image](#docker)!

Alternatively, you can build the Nabarr binary yourself.
To build nabarr on your system, make sure:

1. Your machine runs Linux, macOS or Windows
2. You have [Go](https://golang.org/doc/install) installed (1.14 or later preferred)
3. Clone this repository and cd into it from the terminal
4. Run `make build` from the terminal

You should now have a binary with the name `nabarr` in the appropriate dist sub-directory of the project.

If you need to debug certain Nabarr behaviour, either add the `-v` flag for debug mode or the `-vv` flag for trace mode to get even more details about internal behaviour.

## Introduction

Nabarr configuration is split into three distinct modules:

- Media
- PVR
- RSS

### Media

The media configuration section has only one requirement, a trakt `client_id` must be present as this will be used to fetch metadata for any shows/movies that appear in your RSS feeds.

```yaml
media:
  trakt:
    client_id: trakt-client-id
  omdb:
    api_key: omdb-api-key
```

An omdb `api_key` can also be provided which will be used to supplement trakt data with additional information such as IMDb rating, Metascore rating and Rotten Tomatoes rating.

### PVR

The pvrs configuration section is where you will specify the PVR's that Nabarr will work with.

```yaml
pvrs:
  - name: sonarr
    type: sonarr
    url: https://sonarr.domain.com
    api_key: sonarr-api-key
    quality_profile: WEBDL-1080p
    root_folder: /mnt/unionfs/Media/TV
    filters:
      ignores:
        - 'not (FeedTitle matches "(?i)S\\d\\d?E?\\d?\\d?")'
        - 'FeedTitle matches "(?i)\\d\\d\\d\\d\\s?[\\s\\.\\-]\\d\\d?\\s?[\\s\\.\\-]\\d\\d?"'
        - 'len(Languages) != 1 || "en" not in Languages'
        - 'Runtime < 10 || Runtime > 70'
        - 'Network == ""'
        - 'any (["Hallmark Movies"], {Network contains #})'
        - 'not (any(Country, {# in ["us", "gb", "au", "ca"]}))'
        - 'Year < 2000'
        - 'Year < 2021 && Omdb.ImdbRating < 7.5'
        - 'AiredEpisodes > 200'
        - 'Year > (Now().Year() + 1)'
        - 'any (["WWE", "AEW", "WWF", "NXT", "Live:", "Concert", "Musical", " Edition", "Wrestling"], {Title contains #})'
        - 'len(Genres) == 0'
        - 'any (Genres, {# in ["anime", "talk-show", "news"]})'
        - 'Network in ["Twitch", "Xbox Video", "YouTube"]'
        - 'any (["harry", "potter", "horrid", "henry", "minions", "WWE", "WWF"], {Summary contains #})'
        - 'Title matches "(?i)ru ?wwe.+events.+"'
        - 'Title contains "My 600"'
        - 'TvdbId in ["248783"]'
```

### RSS

The rss configuration section is where you will specify the RSS feeds that Nabarr will work with.

```yaml
rss:
  feeds:
    - name: series premiere
      url: https://rss.indexer.me/rss-search.php?catid=19,20&user=your-username&api=your-api-key&search=S01E01&langs=11&nuke=1&pw=2&nodupe=1&limit=200
      cron: '*/10 * * * *'
      pvrs:
        - sonarr
```

In order for Nabarr to be-able to process items in these feeds, a tvdb and/or imdb id must be present in the feed items.

If there is a tvdb id present, it is assumed that the feed item relates to a TV Series and thus, the item will propagate to any Sonarr PVR specified.

If there is a imdb id present, it is assumed that the feed item relates to a Movie and thus, the item will propagate to any Radarr PVR specified.

### Full config file

With the examples given in the [media](#media), [pvr](#pvr) and [rss](#rss) sections, here is what your full config file *could* look like:

```yaml
media:
  trakt:
    client_id: trakt-client-id
  omdb:
    api_key: omdb-api-key
pvrs:
  - name: sonarr
    type: sonarr
    url: https://sonarr.domain.com
    api_key: sonarr-api-key
    quality_profile: WEBDL-1080p
    root_folder: /mnt/unionfs/Media/TV
    filters:
      ignores:
        - 'not (FeedTitle matches "(?i)S\\d\\d?E?\\d?\\d?")'
        - 'FeedTitle matches "(?i)\\d\\d\\d\\d\\s?[\\s\\.\\-]\\d\\d?\\s?[\\s\\.\\-]\\d\\d?"'
        - 'len(Languages) != 1 || "en" not in Languages'
        - 'Runtime < 10 || Runtime > 70'
        - 'Network == ""'
        - 'any (["Hallmark Movies"], {Network contains #})'
        - 'not (any(Country, {# in ["us", "gb", "au", "ca"]}))'
        - 'Year < 2000'
        - 'Year < 2021 && Omdb.ImdbRating < 7.5'
        - 'AiredEpisodes > 200'
        - 'Year > (Now().Year() + 1)'
        - 'any (["WWE", "AEW", "WWF", "NXT", "Live:", "Concert", "Musical", " Edition", "Wrestling"], {Title contains #})'
        - 'len(Genres) == 0'
        - 'any (Genres, {# in ["anime", "talk-show", "news"]})'
        - 'Network in ["Twitch", "Xbox Video", "YouTube"]'
        - 'any (["harry", "potter", "horrid", "henry", "minions", "WWE", "WWF"], {Summary contains #})'
        - 'Title matches "(?i)ru ?wwe.+events.+"'
        - 'Title contains "My 600"'
        - 'TvdbId in ["248783"]'
  - name: radarr
    type: radarr
    url: https://radarr.domain.com
    api_key: radarr-api-key
    quality_profile: Remux
    root_folder: /mnt/unionfs/Media/Movies
    filters:
      ignores:
        - 'len(Languages) != 1 || "en" not in Languages'
        - 'Runtime < 60'
        - 'len(Genres) == 0'
        - '("music" in Genres || "documentary" in Genres)'
        - 'Year > (Now().Year() + 1)'
        - 'Year < 1980'
        - 'Year < 2021 && (Omdb.Metascore < 55 || Omdb.RottenTomatoes < 55)'
        - 'Title startsWith "Untitled"'
        - 'any (["WWE", "AEW", "WWF", "NXT", "Live:", "Concert", "Musical", " Edition", "Paglaki Ko", "Wrestling ", "UFC on"], {Title contains #})'
        - 'any (["harry", "potter", "horrid", "henry", "minions", "WWE", "WWF"], {Summary contains #})'
        - 'Title matches "^UFC.?\\d.+\\:"'
        - 'ImdbId in ["tt0765458", "tt0892255"]'
        - 'TmdbId in ["11910", "8881"]'
rss:
  feeds:
    - name: series premiere
      url: https://rss.indexer.me/rss-search.php?catid=19,20&user=your-username&api=your-api-key&search=S01E01&langs=11&nuke=1&pw=2&nodupe=1&limit=200
      cron: '*/10 * * * *'
      pvrs:
        - sonarr
```

## Other installation options

### Docker

Nabarr's Docker image provides various versions that are available via tags. The `latest` tag usually provides the latest stable version. Others are considered under development and caution must be exercised when using them.

| Tag | Description |
| :----: | --- |
| latest | Latest stable version from a tagged GitHub release |
| master | Most recent GitHub master commit |

#### Usage

```bash
docker run \
  --name=nabarr \
  -e "PUID=1000" \
  -e "PGID=1001" \
  -v "/opt/nabarr:/config" \
  --restart=unless-stopped \
  -d cloudb0x/nabarr:latest
```

#### Parameters

Nabarr's Docker image supports the following parameters.

| Parameter | Function |
| :----: | --- |
| `-e PUID=1000` | The UserID to run the Nabarr binary as |
| `-e PGID=1000` | The GroupID to run the Nabarr binary as |
| `-e APP_VERBOSITY=0` | The Nabarr logging verbosity level to use. (0 = info, 1 = debug, 2 = trace) |
| `-v /config` | Nabarr's config |

#### Cloudbox

The following Docker setup should work for many Cloudbox users.

**WARNING: You still need to configure the `config.yml` file!**

```bash
docker run \
  --name=nabarr \
  -e "PUID=1000" \
  -e "PGID=1001" \
  -v "/opt/nabarr:/config" \
  --label="com.github.cloudbox.cloudbox_managed=true" \
  --network=cloudbox \
  --network-alias=nabarr  \
  --restart=unless-stopped \
  -d cloudb0x/nabarr:latest
```

## Donate

If you find this project helpful, feel free to make a small donation:

- [Monzo](https://monzo.me/today): Credit Cards, Apple Pay, Google Pay

- [Paypal: l3uddz@gmail.com](https://www.paypal.me/l3uddz)

- [GitHub Sponsor](https://github.com/sponsors/l3uddz): GitHub matches contributions for first 12 months.

- BTC: 3CiHME1HZQsNNcDL6BArG7PbZLa8zUUgjL