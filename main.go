package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"time"

	mastodon "github.com/mattn/go-mastodon"
)

type Configuration struct {
	ServerUrl   string
	AccessToken string

	MessageText        string
	RepeatEveryMinutes int
	OneShot            bool
}

func main() {
	config := Configuration{}
	flag.StringVar(&config.ServerUrl, "server", "", "The Mastodon Server")
	flag.StringVar(&config.AccessToken, "token", "", "the app access token")
	flag.StringVar(&config.MessageText, "message", "Have a panda!", "The message to send with each toot")
	flag.BoolVar(&config.OneShot, "one-shot", false, "Single shot message")
	flag.IntVar(&config.RepeatEveryMinutes, "repeat-duration", 60, "How many minutes per cycle")
	flag.Parse()

	client := mastodon.NewClient(&mastodon.Config{
		Server:      config.ServerUrl,
		AccessToken: config.AccessToken,
	})

	ctx := context.Background()

	files, err := ioutil.ReadDir("img")
	if err != nil {
		log.Panic(err)
	}

	attachment, err := client.UploadMedia(ctx, "img/"+files[10].Name())
	if err != nil {
		log.Printf("Error Uploading Media, %s\n", files[10].Name())
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

	for {

		if config.OneShot {
			return
		}

		time.Sleep(time.Hour)
	}

}
