package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	mastodon "github.com/mattn/go-mastodon"
	cron "github.com/robfig/cron"
)

type Configuration struct {
	ServerUrl   string
	AccessToken string

	MessageText       string
	MessageVisibility string
	Schedule          string
	OneShot           bool
}

func main() {
	config := Configuration{}
	flag.StringVar(&config.ServerUrl, "server", "", "(Required) The Mastodon Server")
	flag.StringVar(&config.AccessToken, "token", "", "(Required) The app access token")
	flag.StringVar(&config.MessageText, "message", "Red Pandas are rad! Have a panda! #RedPanda", "The message to send with each toot")
	flag.StringVar(&config.MessageVisibility, "visibility", "private", "The visibility of the toot (public, unlisted, private, direct)")
	flag.BoolVar(&config.OneShot, "one-shot", false, "Single shot message")
	flag.StringVar(&config.Schedule, "schedule", "@hourly", "A cron expression controlling when to send messages")
	flag.Parse()

	if config.AccessToken == "" || config.ServerUrl == "" {
		flag.PrintDefaults()
		return
	}

	rand.Seed(time.Now().UnixNano())

	if config.OneShot {
		process(config)
		return
	} else {
		scheduler := cron.New()
		err := scheduler.AddFunc(config.Schedule, func() { process(config) })
		if err != nil {
			log.Panic(err)
		}

		scheduler.Run()
	}
}

func process(config Configuration) {

	client := mastodon.NewClient(&mastodon.Config{
		Server:      config.ServerUrl,
		AccessToken: config.AccessToken,
	})

	ctx := context.Background()

	files, err := ioutil.ReadDir("img")
	if err != nil {
		log.Panic(err)
	}

	fileIdx := rand.Intn(len(files)+1) + 1

	attachment, err := client.UploadMedia(ctx, "img/"+files[fileIdx].Name())
	if err != nil {
		log.Printf("Error Uploading Media, %s\n", files[fileIdx].Name())
		log.Fatal(err)
	}

	toot := mastodon.Toot{
		Status:     config.MessageText,
		Visibility: "private",
		MediaIDs: []mastodon.ID{
			attachment.ID,
		},
	}

	_, err = client.PostStatus(ctx, &toot)
	if err != nil {
		log.Println("Error Tooting")
		log.Fatal(err)
	}

}
