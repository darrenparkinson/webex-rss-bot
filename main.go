package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/darrenparkinson/webex-rss-bot/internal/cache"

	"github.com/alexflint/go-arg"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
)

var args struct {
	Interval time.Duration `arg:"-i" help:"interval to check feed" default:"10s"`
	Feed     string        `arg:"-f" help:"url for rss feed" default:"https://status.webex.com/history.rss"`
	Email    string        `arg:"-e,required" help:"email address to send notifications (required)"`
}

func main() {
	godotenv.Load()
	arg.MustParse(&args)

	// set up webex connectivity
	wbx := webexteams.NewClient()

	// keep track in memory of the items we've sent
	sentItems := cache.New[string, bool]()

	// function to check if an item has been sent
	isItemSent := func(itemID string) bool {
		_, found := sentItems.Get(itemID)
		return found
	}

	// function to mark an item as sent
	markItemSent := func(itemID string) {
		sentItems.Set(itemID, true)
	}

	// start a ticker
	ticker := time.NewTicker(args.Interval)
	defer ticker.Stop()

	// listen for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			log.Println("checking feed")
			fp := gofeed.NewParser()
			feed, _ := fp.ParseURL(args.Feed)
			for _, item := range feed.Items {
				if isItemSent(item.GUID) {
					continue
				}
				log.Println("sending", item.Published, item.Title, item.Link)
				markItemSent(item.GUID)
				var markdown strings.Builder
				markdown.WriteString(fmt.Sprintf("**%s**  \n", feed.Title))
				markdown.WriteString(fmt.Sprintf("[%s](%s)  \n", item.Title, item.Link))
				markdown.WriteString("Details:  \n")
				markdown.WriteString(item.Description)
				mcr := &webexteams.MessageCreateRequest{
					ToPersonEmail: args.Email,
					Markdown:      markdown.String(),
				}
				wbx.Messages.CreateMessage(mcr)
			}

		case s := <-quit:
			log.Println("caught interrupt:", s.String())
			os.Exit(0)
		}
	}
}
