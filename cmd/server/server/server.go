package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
)

var insults = [...]string{
	"Ahahahah ! 23 à 0 !",
	"C'est la piquette, Jack !",
	"Tu sais pas jouer, Jack !",
	"T'es mauvais, Jack !",
}

const STX = rune(2)

var emptyResult = &pb.Result{}

type textBot interface {
	Execute(
		channelName string, personName string, msg string, user string,
	) (textbot.Result, bool, error)
}

type UserNickAssociationRepository interface {
	GetUserFromNick(nick string) (string, error)
	GetUserFromCode(code string) (string, error)
	SaveAssociation(user string, nick string) error
}

type server struct {
	textBot    textBot
	repository UserNickAssociationRepository

	ctcMutex     sync.RWMutex
	codesToCheck map[string]string
}

func NewServer(bot textBot, repository UserNickAssociationRepository) *server {
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
			return emptyResult, nil
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
				return emptyResult, err
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
				return emptyResult, err
			}

			return makeResult(&pb.IrcMessage{
				Msg: "Association effectuée. Utilise la commande !fav pour enregistrer un lien dans tes favoris.",
				To:  nick,
			}), nil
		}

		return emptyResult, nil
	}

	user, err := s.repository.GetUserFromNick(msg.PersonName)
	if err != nil {
		return emptyResult, nil
	}

	result, cmd, err := s.textBot.Execute(msg.ChannelName, msg.PersonName, msg.Msg, user)
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

		} else if errors.Is(err, playbot.InvalidOffsetError) {
			return makeResult(&pb.IrcMessage{
				Msg: "Offset invalide",
				To:  msg.ChannelName,
			}), nil
		} else if errors.Is(err, textbot.OffsetToBigError) {
			return makeResult(&pb.IrcMessage{
				Msg: "T'as compté tout ça sans te tromper, srsly ?",
				To:  msg.ChannelName,
			}), nil
		} else if errors.Is(err, textbot.InvalidUsageError) {
			return makeResult(&pb.IrcMessage{
				Msg: insults[rand.Intn(len(insults))],
				To:  msg.ChannelName,
			}), nil
		} else if errors.Is(err, textbot.AuthenticationRequired) {
			authMsg := &pb.IrcMessage{
				Msg: "Ce nick n'est associé à aucun login arise. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
				To:  msg.PersonName,
			}

			if result.ID != 0 {
				return makeResult(
					&pb.IrcMessage{
						Msg: printMusicRecord(result),
						To:  msg.ChannelName,
					},
					authMsg,
				), nil
			} else {
				return makeResult(authMsg), nil
			}
		}

		return emptyResult, err
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
	}

	return emptyResult, nil
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
