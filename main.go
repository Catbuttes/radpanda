package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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

func getConfig() Configuration {
	config := Configuration{}
	flag.StringVar(&config.ServerUrl, "server", "", "(Required) The Mastodon Server")
	flag.StringVar(&config.AccessToken, "token", "", "(Required) The app access token")
	flag.StringVar(&config.MessageText, "message", "", "The message to send with each toot")
	flag.StringVar(&config.MessageVisibility, "visibility", "unlisted", "The visibility of the toot (public, unlisted, private, direct)")
	flag.BoolVar(&config.OneShot, "one-shot", false, "Single shot message")
	flag.StringVar(&config.Schedule, "schedule", "@every 1h", "A cron expression controlling when to send messages")
	flag.StringVar(&config.MetricsPort, "metrics-address", "", "The address and port to listen on for prometheus metrics")
	flag.Parse()

	if config.AccessToken == "" {
		config.AccessToken = os.Getenv("RADPANDA_TOKEN")
	}

	if config.ServerUrl == "" {
		config.ServerUrl = os.Getenv("RADPANDA_SERVER")
	}

	if config.AccessToken == "" || config.ServerUrl == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if config.MessageText == "" {
		config.MessageText = os.Getenv("RADPANDA_TEXT")
	}

	if config.MessageVisibility == "unlisted" {
		config.MessageVisibility = os.Getenv("RADPANDA_VISIBILITY")
	}

	if config.Schedule == "@every 1h" {
		config.Schedule = os.Getenv("RADPANDA_SCHEDULE")
	}

	if config.MetricsPort == "" {
		config.MetricsPort = os.Getenv("RADPANDA_METRICS_PATH")
	}

	return config
}

func main() {

	config := getConfig()

	fmt.Printf("Running with config:\n-------\n%+v\n", config)

	if config.MetricsPort != "" {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			http.ListenAndServe(config.MetricsPort, nil)
		}()
	}

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
