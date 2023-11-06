package server

import (
	"context"
	"log"

	pb "github.com/erdnaxeli/PlayBot/cmd/server/rpc"
	"github.com/erdnaxeli/PlayBot/playbot"
)

func (s *server) Execute(ctx context.Context, msg *pb.TextMessage) (*pb.Result, error) {
	log.Printf(
		"Parsing message by %s in %s: %s",
		msg.PersonName,
		msg.ChannelName,
		msg.Msg,
	)

	if msg.ChannelName == s.botNick && msg.Msg[0:2] == "PB" {
		return s.handleUserAuth(msg)
	} else if msg.PersonName == "NickServ" {
		return s.handleUserAuthCallback(msg)
	} else if msg.ChannelName == s.botNick {
		// This is a private discussion. We set the channel name as the sender
		// name.
		msg.ChannelName = msg.PersonName
	}

	user, err := s.repository.GetUserFromNick(msg.PersonName)
	if err != nil {
		return emptyResult, err
	}

	result, cmd, err := s.textBot.Execute(msg.ChannelName, msg.PersonName, msg.Msg, user)
	if err != nil {
		return s.handleError(msg, result, err)
	}

	if result.ID != 0 {
		// A new record was saved, or a command returned a music record.

		if !cmd {
			// It is a new record, we don't print the URL.
			result.Url = ""
		}

		resultMsg := s.recordPrinter.Print(result)
		return makeResult(&pb.IrcMessage{Msg: resultMsg, To: msg.ChannelName}), nil
	} else if (result.Statistics != playbot.MusicRecordStatistics{}) {
		// Statistics were returned.

		return makeResult(&pb.IrcMessage{
			Msg: s.statsPrinter.Print(result.Statistics),
			To:  msg.ChannelName,
		}), nil
	}

	return emptyResult, nil
}
