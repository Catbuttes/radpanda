package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
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
	flag.StringVar(&config.MessageText, "message", "Red Pandas are rad! Have a panda!", "The message to send with each toot")
	flag.BoolVar(&config.OneShot, "one-shot", false, "Single shot message")
	flag.IntVar(&config.RepeatEveryMinutes, "repeat-duration", 60, "How many minutes per cycle")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	client := mastodon.NewClient(&mastodon.Config{
		Server:      config.ServerUrl,
		AccessToken: config.AccessToken,
	})

	ctx := context.Background()

	for {

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

		if config.OneShot {
			return
		}

		time.Sleep(time.Hour)
	}

}
