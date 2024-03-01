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

// Start of text char
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

type MusicRecordPrinter interface {
	Print(record textbot.Result) string
}

type StatisticsPrinter interface {
	Print(stats playbot.MusicRecordStatistics) string
}

type server struct {
	botNick       string
	textBot       textBot
	repository    UserNickAssociationRepository
	recordPrinter MusicRecordPrinter
	statsPrinter  StatisticsPrinter

	ctcMutex     sync.RWMutex
	codesToCheck map[string]string
}

func NewServer(
	nick string,
	bot textBot,
	repository UserNickAssociationRepository,
	recordPrinter MusicRecordPrinter,
	statsPrinter StatisticsPrinter,
) *server {
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
