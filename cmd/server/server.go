package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"

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

const STX = rune(2)

type textBot interface {
	Execute(string, string, string) (textbot.Result, bool, error)
}

type userNickAssociationRepository interface {
	GetUserFromNick(nick string) (string, error)
	GetUserFromCode(code string) (string, error)
	SaveAssociation(user string, nick string) error
}

type server struct {
	textBot    textBot
	repository userNickAssociationRepository

	ctcMutex     sync.RWMutex
	codesToCheck map[string]string
}

func NewServer(bot textBot, repository userNickAssociationRepository) *server {
	return &server{
		textBot:      bot,
		repository:   repository,
		codesToCheck: make(map[string]string),
	}
}

func (s *server) Execute(ctx context.Context, msg *pb.TextMessage) (*pb.Result, error) {
	log.Printf(
		"Parsing message by %s in %s: %s",
		msg.PersonName,
		msg.ChannelName,
		msg.Msg,
	)

	if msg.ChannelName == "PlayTest" && msg.Msg[0:2] == "PB" {
		// user authentication
		s.ctcMutex.Lock()
		defer s.ctcMutex.Unlock()

		s.codesToCheck[msg.PersonName] = msg.Msg

		return makeResult(
			&pb.IrcMessage{
				Msg: "Vérification en cours…",
				To:  msg.PersonName,
			},
			&pb.IrcMessage{
				Msg: fmt.Sprintf("info %s", msg.PersonName),
				To:  "NickServ",
			},
		), nil
	}
	if msg.PersonName == "NickServ" {
		if len(s.codesToCheck) == 0 {
			return &pb.Result{}, nil
		}

		log.Print("Received a message from NickServ.")
		s.ctcMutex.Lock()
		defer s.ctcMutex.Unlock()

		for nick, code := range s.codesToCheck {
			log.Printf("Trying to auth %s.", nick)

			if msg.Msg == fmt.Sprintf("Le pseudo %s%s%s n'est pas enregistré.", string(STX), nick, string(STX)) {
				log.Print("Unregistered nick")
				return makeResult(&pb.IrcMessage{
					Msg: "Il faut que ton pseudo soit enregistré auprès de NickServ pour pouvoir t'authentifier.",
					To:  nick,
				}), nil
			} else if msg.Msg != fmt.Sprintf("%s est actuellement connecté.", nick) {
				continue
			}

			log.Printf("Ok, authenticating nick %s.", nick)
			delete(s.codesToCheck, nick)

			user, err := s.repository.GetUserFromCode(code)
			if err != nil {
				return &pb.Result{}, err
			} else if user == "" {
				log.Printf("Unknown code.")
				return makeResult(&pb.IrcMessage{
					Msg: "Code inconnu. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
					To:  nick,
				}), nil
			}

			log.Printf("Code ok.")
			err = s.repository.SaveAssociation(user, nick)
			if err != nil {
				return &pb.Result{}, err
			}

			return makeResult(&pb.IrcMessage{
				Msg: "Association effectuée. Utilise la commande !fav pour enregistrer un lien dans tes favoris.",
				To:  nick,
			}), nil
		}

		return &pb.Result{}, nil
	}

	result, cmd, err := s.textBot.Execute(msg.ChannelName, msg.PersonName, msg.Msg)
	if err != nil {
		if errors.Is(err, playbot.NoRecordFoundError) {
			if result.Count > 0 {
				return makeResult(&pb.IrcMessage{
					Msg: "Tu tournes en rond, Jack !",
					To:  msg.ChannelName,
				}), nil
			} else {
				return makeResult(&pb.IrcMessage{
					Msg: "Je n'ai rien dans ce registre.",
					To:  msg.ChannelName,
				}), nil
			}
		} else if errors.Is(err, textbot.OffsetToBigError) {
			return makeResult(&pb.IrcMessage{
				Msg: "T'as compté tout ça sans te tromper, srsly ?",
				To:  msg.ChannelName,
			}), nil
		} else if errors.Is(err, playbot.InvalidOffsetError) {
			return makeResult(&pb.IrcMessage{
				Msg: "Offset invalide",
				To:  msg.ChannelName,
			}), nil
		} else if errors.Is(err, textbot.InvalidUsageError) {
			return makeResult(&pb.IrcMessage{
				Msg: insults[rand.Intn(len(insults))],
				To:  msg.ChannelName,
			}), nil
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
		return makeResult(&pb.IrcMessage{Msg: resultMsg, To: msg.ChannelName}), nil
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

		return makeResult(&pb.IrcMessage{
			Msg: resultMsg.String(),
			To:  msg.ChannelName,
		}), nil
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

func makeResult(messages ...*pb.IrcMessage) *pb.Result {
	return &pb.Result{
		Msg: messages,
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
	server := NewServer(bot, repository)
	handler := pb.NewPlaybotCliServer(server)

	log.Print("Starting the server")
	return http.ListenAndServe("localhost:1111", handler)
}
