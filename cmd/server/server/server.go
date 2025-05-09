// Package server provides an object implementing the pb.PlaybotCli interface.
package server

import (
	"fmt"
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

// STX represents the start of text character.
const STX = rune(2)

var emptyResult = &pb.Result{}

// TextBot is the interface that wraps the main Execute method.
type TextBot interface {
	// Execute parses a message from a channel, store any record found in the message, and execute any command found.
	//
	// The format of the command depends on the implementation.
	Execute(
		channelName string, personName string, msg string, user string,
	) (textbot.Result, bool, error)
}

// UserNickAssociationRepository is the interface that wraps methods to authenticate an user.
type UserNickAssociationRepository interface {
	GetUserFromNick(nick string) (string, error)
	GetUserFromCode(code string) (string, error)
	SaveAssociation(user string, nick string) error
}

// MusicRecordPrinter is the interface that wraps the PrintResult method.
type MusicRecordPrinter interface {
	PrintResult(record textbot.Result) string
}

// StatisticsPrinter is the interface that wraps the PrintStatistics methods.
type StatisticsPrinter interface {
	PrintStatistics(stats playbot.MusicRecordStatistics) string
}

type server struct {
	botNick       string
	textBot       TextBot
	repository    UserNickAssociationRepository
	recordPrinter MusicRecordPrinter
	statsPrinter  StatisticsPrinter

	ctcMutex     sync.RWMutex
	codesToCheck map[string]string
}

// NewServer returns a new instance of a Server.
func NewServer(
	nick string,
	bot TextBot,
	repository UserNickAssociationRepository,
	recordPrinter MusicRecordPrinter,
	statsPrinter StatisticsPrinter,
) pb.PlaybotCli {
	return &server{
		botNick:       nick,
		textBot:       bot,
		repository:    repository,
		recordPrinter: recordPrinter,
		statsPrinter:  statsPrinter,
		codesToCheck:  make(map[string]string),
	}
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
