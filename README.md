# Webex RSS Bot

A simple RSS reader that publishes its messages to Webex as a POC.

> **Note:** It will send *all* items it finds in the feed first time around as a demonstration, so be careful if you
> have a large feed.

## Usage

1. Ensure you have a webex bot token as `WEBEX_TEAMS_ACCESS_TOKEN` in the environment. This will also be automatically
loaded from a `.env` file if present.
1. Run `go run ./ --email <email address>`

You can use `webex-rss-bot -h` to get help:

```bash
$ webex-rss-bot -h

Usage: webex-rss-bot [--interval INTERVAL] [--feed FEED] --email EMAIL

Options:
  --interval INTERVAL, -i INTERVAL
                         interval to check feed [default: 10s]
  --feed FEED, -f FEED   url for rss feed [default: https://status.webex.com/history.rss]
  --email EMAIL, -e EMAIL
                         email address to send notifications (required)
  --help, -h             display this help and exit
```

## TODO

Some things we could do to make this more useful.

* [ ] Refactor to enable tests
* [ ] Add tests
* [ ] Don't send all items on first pass
* [ ] Accept multiple feeds, perhaps from a config file
* [ ] Add persistent storage to keep a list of seen items between restarts
* [ ] Allow use of room id to send notifications
