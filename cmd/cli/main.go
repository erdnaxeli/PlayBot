package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository"
	"github.com/erdnaxeli/PlayBot/types"
)

var bot *playbot.Playbot

func init() {
	config, err := config.ReadConfigFile("playbot.conf")
	if err != nil {
		log.Fatal(err)
	}

	ldjsonExtractor := ldjson.NewLdJsonExtractor()
	extractor := extractors.New(
		extractors.NewBandcampExtractor(ldjsonExtractor),
		extractors.NewSoundCloudExtractor(ldjsonExtractor),
		&extractors.YoutubeExtractor{
			ApiKey: config.YoutubeApiKey,
		},
	)

	repository, err := repository.NewSqlRepository(
		fmt.Sprintf(
			"%s:%s@(%s)/%s",
			config.DbUser,
			config.DbPassword,
			config.DbHost,
			config.DbName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	bot = playbot.New(extractor, repository)
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s CHANNEL PERSON MESSAGE", os.Args[0])
	}

	channel := types.Channel{Name: os.Args[1]}
	person := types.Person{Name: os.Args[2]}
	msg := os.Args[3]

	recordId := saveMusicRecord(msg, person, channel)
	saveTags(msg, recordId)
}

func saveMusicRecord(msg string, person types.Person, channel types.Channel) int64 {
	recordId, err := bot.SaveMusicRecord(msg, person, channel)
	if err != nil {
		log.Fatal(err)
	}

	return recordId
}

func saveTags(msg string, recordId int64) {
	re := regexp.MustCompile(`\s+`)
	var tags []string
	for _, word := range re.Split(msg, -1) {
		if word[0] == '#' {
			tags = append(tags, word[0:])
		}
	}

	err := bot.SaveTags(recordId, tags)
	if err != nil {
		log.Fatal(err)
	}
}
