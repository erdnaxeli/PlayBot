package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	pb "github.com/erdnaxeli/PlayBot/cmd/cli/rpc"
	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository/mariadb"
	"github.com/erdnaxeli/PlayBot/textbot"
)

type textBot interface {
	Execute(string, string, string) (textbot.Result, bool, error)
}

type server struct {
	textBot textBot
}

func (s server) Execute(ctx context.Context, msg *pb.TextMessage) (*pb.Result, error) {
	log.Printf(
		"Parsing message by %s in %s: %s",
		msg.PersonName,
		msg.ChannelName,
		msg.Msg,
	)

	result, cmd, err := s.textBot.Execute(msg.ChannelName, msg.PersonName, msg.Msg)
	if err != nil {
		if errors.Is(err, textbot.NoRecordFound{}) {
			return &pb.Result{Msg: "Je n'ai rien dans ce registre."}, nil
		}

		return &pb.Result{}, err
	}

	if result.ID != 0 {
		// A new record was saved, or a command returned a music record.

		if !cmd {
			// It is a new record, we don't print the URL.
			result.Url = ""
		}

		resultMsg := printMusicRecord(result)
		return &pb.Result{Msg: resultMsg}, nil
	} else if !cmd {
		// No record was saved nor command executed.
		log.Print("unknown command or record")
		return &pb.Result{}, errors.New("unknown command or record")
	}

	return &pb.Result{}, nil
}

func startServer() error {
	config, err := config.ReadConfigFile("playbot.conf")
	if err != nil {
		return err
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
		return err
	}

	bot := textbot.New(playbot.New(extractor, repository))
	server := server{textBot: bot}
	handler := pb.NewPlaybotCliServer(server)

	log.Print("Starting the server")
	return http.ListenAndServe("localhost:1111", handler)
}
