package main

import (
	"log"
	"net/http"
	"time"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/cmd/server/server"
	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository/mariadb"
	"github.com/erdnaxeli/PlayBot/textbot"
)

func startServer() error {
	config, err := config.ReadConfig()
	if err != nil {
		return err
	}

	ldjsonExtractor := ldjson.New()
	extractor := extractors.New(
		extractors.NewBandcampExtractor(ldjsonExtractor),
		extractors.NewSoundCloudExtractor(ldjsonExtractor),
		&extractors.YoutubeExtractor{
			APIKey: config.YoutubeAPIKey,
		},
	)

	repository, err := mariadb.New(
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
	)
	if err != nil {
		return err
	}

	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		return err
	}

	bot := textbot.New(playbot.New(extractor, repository))
	server := server.NewServer(
		config.IRC.Nick,
		bot,
		repository,
		server.IrcMusicRecordPrinter{},
		server.IrcStatisticsPrinter{Location: location},
	)
	handler := pb.NewPlaybotCliServer(server)

	log.Print("Starting the server")
	return http.ListenAndServe(config.ServerAddress, handler)
}
