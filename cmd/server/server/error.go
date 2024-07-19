package server

import (
	"errors"
	"math/rand/v2"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
)

func (s *server) handleError(msg *pb.TextMessage, result textbot.Result, err error) (*pb.Result, error) {
	var messages []*pb.IrcMessage

	if errors.Is(err, playbot.NoRecordFoundError) {
		if result.Count > 0 {
			messages = append(messages, &pb.IrcMessage{
				Msg: "Tu tournes en rond, Jack !",
				To:  msg.ChannelName,
			})
		} else {
			messages = append(messages, &pb.IrcMessage{
				Msg: "Je n'ai rien dans ce registre.",
				To:  msg.ChannelName,
			})
		}
	} else if errors.Is(err, playbot.InvalidOffsetError) {
		messages = append(messages, &pb.IrcMessage{
			Msg: "Offset invalide.",
			To:  msg.ChannelName,
		})
	} else if errors.Is(err, textbot.OffsetToBigError) {
		messages = append(messages, &pb.IrcMessage{
			Msg: "T'as compté tout ça sans te tromper, srsly ?",
			To:  msg.ChannelName,
		})
	} else if errors.Is(err, textbot.InvalidUsageError) {
		messages = append(messages, &pb.IrcMessage{
			Msg: insults[rand.IntN(len(insults))],
			To:  msg.ChannelName,
		})
	} else if errors.Is(err, textbot.AuthenticationRequired) {
		messages = append(messages, &pb.IrcMessage{
			Msg: "Ce nick n'est associé à aucun login arise. Va sur http://nightiies.iiens.net/links/fav pour obtenir ton code personel.",
			To:  msg.PersonName,
		})
	}

	if len(messages) > 0 {
		// the error was handled
		err = nil
	}

	if result.ID != 0 {
		messages = append(messages, &pb.IrcMessage{
			Msg: s.recordPrinter.Print(result),
			To:  msg.ChannelName,
		})
	}
	return &pb.Result{Msg: messages}, err
}
