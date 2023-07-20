package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
)

func main() {
	msg := os.Args[1]

	config, err := readConfigFile("playbot.conf")
	if err != nil {
		log.Fatal(err)
	}
	extractor := extractors.New(
		extractors.NewSoundCloudExtractor(
			ldjson.NewLdJsonExtractor(),
		),
		&extractors.YoutubeExtractor{
			ApiKey: config.YoutubeApiKey,
		},
	)

	matchedUrl, music, err := extractor.Extract(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matchedUrl)
	fmt.Println(music.RecordId)
	fmt.Println(music.Url)
	fmt.Println(music.Name)
	fmt.Println(music.Band.Name)
	fmt.Println(music.Duration.Seconds())
}
