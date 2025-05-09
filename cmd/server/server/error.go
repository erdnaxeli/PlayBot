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

	if errors.Is(err, playbot.ErrNoRecordFound) {
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
	} else if errors.Is(err, playbot.ErrInvalidOffset) {
		messages = append(messages, &pb.IrcMessage{
			Msg: "Offset invalide.",
			To:  msg.ChannelName,
		})
	} else if errors.Is(err, textbot.ErrOffsetToBig) {
		messages = append(messages, &pb.IrcMessage{
			Msg: "T'as compté tout ça sans te tromper, srsly ?",
			To:  msg.ChannelName,
		})
	} else if errors.Is(err, textbot.ErrInvalidUsage) {
		messages = append(messages, &pb.IrcMessage{
			Msg: insults[rand.IntN(len(insults))],
			To:  msg.ChannelName,
		})
	} else if errors.Is(err, textbot.ErrAuthenticationRequired) {
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
			Msg: s.recordPrinter.PrintResult(result),
			To:  msg.ChannelName,
		})
	}
	return &pb.Result{Msg: messages}, err
}
