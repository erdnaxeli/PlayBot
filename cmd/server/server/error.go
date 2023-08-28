package server

import (
	"errors"
	"math/rand"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
)

func (s *server) handleError(msg *pb.TextMessage, result textbot.Result, err error) (*pb.Result, error) {
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
			Msg: "Offset invalide.",
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
					Msg: s.recordPrinter.Print(result),
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
