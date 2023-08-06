package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository/mariadb"
	"github.com/erdnaxeli/PlayBot/textbot"
	"github.com/erdnaxeli/PlayBot/types"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s CHANNEL PERSON MESSAGE", os.Args[0])
	}

	channel := os.Args[1]
	person := os.Args[2]
	msg := os.Args[3]

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

	bot := textbot.New(playbot.New(extractor, repository))
	result, cmd, err := bot.Execute(channel, person, msg)
	if err != nil {
		log.Fatal(err)
	}

	if result.ID != 0 {
		// A new record was saved, or a command returned a music record.
		printMusicRecord(
			result.ID,
			result.MusicRecord,
			result.Tags,
			result.IsNew,
		)
	} else if !cmd {
		// No music record was saved and no command was executed. We need to exit with
		// with an error so the perl code will try to parse the message.
		log.Fatal("the message cannot be interpreted")
	}
}

func printMusicRecord(
	recordId int64, record types.MusicRecord, tags []string, isNew bool,
) {
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
