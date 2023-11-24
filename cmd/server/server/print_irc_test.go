package server_test

import (
	"testing"
	"time"

	"github.com/erdnaxeli/PlayBot/cmd/server/server"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
)

func TestIrcStaticPrinter(t *testing.T) {
	location := time.FixedZone("", -2*60*60)
	stats := playbot.MusicRecordStatistics{
		PostsCount:       10,
		PeopleCount:      3,
		ChannelsCount:    5,
		MaxPerson:        types.Person{Name: "George"},
		MaxPersonCount:   2,
		MaxChannelCount:  4,
		FirstPostPerson:  types.Person{Name: "Abitbol"},
		FirstPostChannel: types.Channel{Name: "#bigphatsubwoofer"},
		FirstPostDate:    time.Date(2023, 11, 19, 21, 22, 25, 12345, location),
		FavoritesCount:   15,
	}
	printer := server.IrcStatisticsPrinter{Location: time.FixedZone("", +2*60*60)}

	// test
	result := printer.Print(stats)

	// assertions
	assert.Equal(t,
		"Posté la 1re fois par Abitbol le 20-11-2023 01:22:25 sur #bigphatsubwoofer. Posté 10 fois par 3 personnes sur 5 channels. George l'a posté 2 fois. Sauvegardé dans ses favoris par 15 personnes.",
		result,
	)
}
