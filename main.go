package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	mastodon "github.com/mattn/go-mastodon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cron "github.com/robfig/cron"
)

type Configuration struct {
	ServerUrl   string
	AccessToken string

	MessageText       string
	MessageVisibility string
	Schedule          string
	OneShot           bool
	MetricsPort       string
}

var (
	imgCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "radpanda_current_image_count",
		Help: "The total number of images available to radpanda to post",
	})
)

func main() {
	config := Configuration{}
	flag.StringVar(&config.ServerUrl, "server", "", "(Required) The Mastodon Server")
	flag.StringVar(&config.AccessToken, "token", "", "(Required) The app access token")
	flag.StringVar(&config.MessageText, "message", "Red Pandas are rad! Have a panda! #RedPanda", "The message to send with each toot")
	flag.StringVar(&config.MessageVisibility, "visibility", "private", "The visibility of the toot (public, unlisted, private, direct)")
	flag.BoolVar(&config.OneShot, "one-shot", false, "Single shot message")
	flag.StringVar(&config.Schedule, "schedule", "@hourly", "A cron expression controlling when to send messages")
	flag.StringVar(&config.MetricsPort, "metrics-address", ":2112", "The address and port to listen on for prometheus metrics")
	flag.Parse()

	if config.AccessToken == "" || config.ServerUrl == "" {
		flag.PrintDefaults()
		return
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(config.MetricsPort, nil)
	}()

	files, err := ioutil.ReadDir("img")
	if err != nil {
		log.Panic(err)
	}

	imgCount.Set(float64(len(files)))

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

	imgCount.Set(float64(len(files)))

	fileIdx := rand.Intn(len(files))

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
