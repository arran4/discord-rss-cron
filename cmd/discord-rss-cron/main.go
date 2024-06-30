package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adrg/xdg"
	"github.com/mmcdole/gofeed"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Last struct {
	Check   *time.Time
	HookUrl string
	FeedUrl string
}

func main() {
	cfnp := fmt.Sprintf("discord-rss-webhook/%s", "channel-1")
	cfn, err := xdg.ConfigFile(cfnp)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	log.Printf("Config: %s", cfn)
	cfnb, err := os.ReadFile(cfn)
	last := &Last{}
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
		default:
			log.Panicf("Error: %s", err)
		}
	} else if err := json.Unmarshal(cfnb, &last); err != nil {
		log.Panicf("Error: %s", err)
		return
	}
	if last.FeedUrl == "" || last.HookUrl == "" {
		log.Printf("FeedUrl, or HookUrl in %s empty / not present. Please edit %s and add the values", cfn, cfn)
		cfnb, err = json.Marshal(last)
		if err != nil {
			log.Panicf("Error: %s", err)
		}
		err = os.WriteFile(cfn, cfnb, 0644)
		if err != nil {
			log.Panicf("Error: %s", err)
		}
		os.Exit(-1)
		return
	}
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(last.FeedUrl)
	var largestUnseen *gofeed.Item
	for _, item := range feed.Items {
		if item.PublishedParsed == nil {
			continue
		}
		if last.Check != nil && !last.Check.Before(*item.PublishedParsed) {
			continue
		}
		if largestUnseen != nil && largestUnseen.PublishedParsed.After(*item.PublishedParsed) {
			continue
		}
		largestUnseen = item
	}
	if largestUnseen == nil {
		log.Printf("Done nothing found")
		return
	}
	last.Check = largestUnseen.PublishedParsed
	log.Printf("Sending %s", largestUnseen.Link)
	omb, err := json.Marshal(map[string]any{
		"content":  largestUnseen.Link,
		"username": "NPR Music",
		//"avatarURL": "https://i.imgur.com/AfFp7pu.png",
	})
	if err != nil {
		log.Panicf("Error: %s", err)
	}

	r, err := http.NewRequest("POST", last.HookUrl, bytes.NewReader(omb))
	if err != nil {
		log.Panicf("Error: %s", err)
		return
	}

	r.Header.Add("Content-Type", "application/json")

	do, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Panicf("Error: %s", err)
		return
	}
	respB, err := io.ReadAll(do.Body)
	if err != nil {
		log.Panicf("Error: %s", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicf("Error: %s", err)
		}
	}(do.Body)

	log.Printf("%s", string(respB))

	cfnb, err = json.Marshal(last)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	err = os.WriteFile(cfn, cfnb, 0644)
	if err != nil {
		log.Panicf("Error: %s", err)
	}
	log.Printf("Done")
}
