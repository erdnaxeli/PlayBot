package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository/mariadb"
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

	repository, err := mariadb.New(
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

	recordId, musicRecord, isNew := saveMusicRecord(msg, person, channel)
	saveTags(msg, recordId)
	tags := getTags(recordId)

	log.Println("Record saved", recordId, musicRecord)
	printMusicRecord(recordId, musicRecord, tags, isNew)
}

func saveMusicRecord(msg string, person types.Person, channel types.Channel) (int64, types.MusicRecord, bool) {
	recordId, musicRecord, isNew, err := bot.SaveMusicRecord(msg, person, channel)
	if err != nil {
		log.Fatal(err)
	}

	return recordId, musicRecord, isNew
}

func saveTags(msg string, recordId int64) {
	re := regexp.MustCompile(`\s+`)
	var tags []string
	for _, word := range re.Split(msg, -1) {
		if word[0] == '#' {
			tags = append(tags, word[1:])
		}
	}

	err := bot.SaveTags(recordId, tags)
	if err != nil {
		log.Fatal(err)
	}
}

func getTags(recordId int64) []string {
	tags, err := bot.GetTags(recordId)
	if err != nil {
		log.Fatal(err)
	}

	return tags
}

func printMusicRecord(recordId int64, record types.MusicRecord, tags []string, isNew bool) {
	if isNew {
		fmt.Print("+", recordId, "\n")
	} else {
		fmt.Println(recordId)
	}
	fmt.Println(record.RecordId)
	fmt.Println("")
	fmt.Println("🎉", record.Name)
	fmt.Println(record.Band.Name)
	fmt.Println(record.Duration.Seconds())

	for idx := range tags {
		tags[idx] = "#" + tags[idx]
	}
	fmt.Println(strings.Join(tags, " "))
}
