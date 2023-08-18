package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/config"
	"github.com/erdnaxeli/PlayBot/extractors"
	"github.com/erdnaxeli/PlayBot/extractors/ldjson"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository/mariadb"
	"github.com/erdnaxeli/PlayBot/textbot"
)

var insults = [...]string{
	"Ahahahah ! 23 à 0 !",
	"C'est la piquette, Jack !",
	"Tu sais pas jouer, Jack !",
	"T'es mauvais, Jack !",
}

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
		if errors.Is(err, playbot.NoRecordFoundError) {
			if result.Count > 0 {
				return &pb.Result{
					Msg: "Tu tournes en rond, Jack !",
					To:  msg.ChannelName,
				}, nil
			} else {
				return &pb.Result{
					Msg: "Je n'ai rien dans ce registre.",
					To:  msg.ChannelName,
				}, nil
			}
		} else if errors.Is(err, textbot.OffsetToBigError) {
			return &pb.Result{
				Msg: "T'as compté tout ça sans te tromper, srsly ?",
				To:  msg.ChannelName,
			}, nil
		} else if errors.Is(err, playbot.InvalidOffsetError) {
			return &pb.Result{
				Msg: "Offset invalide",
				To:  msg.ChannelName,
			}, nil
		} else if errors.Is(err, textbot.InvalidUsageError) {
			return &pb.Result{
				Msg: insults[rand.Intn(len(insults))],
				To:  msg.ChannelName,
			}, nil
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
		return &pb.Result{Msg: resultMsg, To: msg.ChannelName}, nil
	} else if (result.Statistics != playbot.MusicRecordStatistics{}) {
		// Statistics were returned.
		var resultMsg strings.Builder
		fmt.Fprintf(
			&resultMsg,
			"Posté la 1re fois par %s le %s sur %s.",
			result.Statistics.FirstPostPerson.Name,
			result.Statistics.FirstPostDate.Format("01-02-2006 15:04:05"),
			result.Statistics.FirstPostChannel.Name,
		)

		fmt.Fprintf(
			&resultMsg,
			" Posté %d fois par %d personne",
			result.Statistics.PostsCount,
			result.Statistics.PeopleCount,
		)
		plural(result.Statistics.PeopleCount, &resultMsg)

		fmt.Fprintf(
			&resultMsg, " sur %d channel", result.Statistics.ChannelsCount,
		)
		plural(result.Statistics.ChannelsCount, &resultMsg)
		fmt.Fprint(&resultMsg, ".")

		fmt.Fprintf(
			&resultMsg,
			" %s l'a posté %d fois.",
			result.Statistics.MaxPerson.Name,
			result.Statistics.MaxPersonCount,
		)

		fmt.Fprintf(
			&resultMsg,
			" Sauvegardé dans ses favoris par %d personne",
			result.Statistics.FavoritesCount,
		)
		plural(result.Statistics.FavoritesCount, &resultMsg)
		fmt.Fprint(&resultMsg, ".")

		return &pb.Result{
			Msg: resultMsg.String(),
			To:  msg.ChannelName,
		}, nil
	} else if !cmd {
		// No record was saved nor command executed.
		log.Print("unknown command or record")
		return &pb.Result{}, errors.New("unknown command or record")
	}

	return &pb.Result{}, nil
}

func plural(count int, builder *strings.Builder) {
	if count == 0 || count > 1 {
		fmt.Fprint(builder, "s")
	}
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
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbName,
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
