# Radpanda

Radpanda is a mastodon bot that will post an image from a local source on a configurable delay with a configurable message.

## Usage

```
Usage of radpanda
  -message string
    	The message to send with each toot (default "Red Pandas are rad! Have a panda!")
  -one-shot
    	Single shot message
  -schedule string
    	A cron expression controlling when to send messages (default "@hourly")
  -server string
    	The Mastodon Server
  -token string
    	the app access token
  -visibility string
    	The visibility of the toot (public, unlisted, private, direct) (default "private")
```