package server_test

import (
	"testing"
	"time"

	"github.com/erdnaxeli/PlayBot/cmd/server/server"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
)

func TestMatrixMusicRecordPrinter_new(t *testing.T) {
	result := textbot.Result{
		ID: 1243,
		MusicRecord: types.MusicRecord{
			Band: types.Band{
				Name: "Some Band",
			},
			Duration: 100 * time.Second,
			Name:     "Some track",
			RecordId: "some-id",
			Source:   "some-source",
			Url:      "some-url",
		},
		Tags:       []string{"some", "tags"},
		IsNew:      true,
		Count:      42,
		Statistics: playbot.MusicRecordStatistics{},
	}
	printer := server.MatrixMusicRecordPrinter{}

	// test
	msg := printer.Print(result)

	// assertions
	assert.Equal(
		t,
		`<font color="#FFFF00">[+1243]</font> <font color="#009300">Some track | Some Band</font> <font color="#0000FC">1m40s</font> =&gt; some-url <font color="#FC7F00">#some #tags</font> <font color="#7F7F7F">[42 résultats]</font>`,
		msg,
	)
}

func TestMatrixMusicRecordPrinter_old(t *testing.T) {
	result := textbot.Result{
		ID: 1243,
		MusicRecord: types.MusicRecord{
			Band: types.Band{
				Name: "Some Band",
			},
			Duration: 100 * time.Second,
			Name:     "Some track",
			RecordId: "some-id",
			Source:   "some-source",
			Url:      "some-url",
		},
		Tags:       []string{"some", "tags"},
		IsNew:      false,
		Count:      42,
		Statistics: playbot.MusicRecordStatistics{},
	}
	printer := server.MatrixMusicRecordPrinter{}

	// test
	msg := printer.Print(result)

	// assertions
	assert.Equal(
		t,
		`<font color="#FFFF00">[1243]</font> <font color="#009300">Some track | Some Band</font> <font color="#0000FC">1m40s</font> =&gt; some-url <font color="#FC7F00">#some #tags</font> <font color="#7F7F7F">[42 résultats]</font>`,
		msg,
	)
}

func TestMatrixMusicRecordPrinter_noUrl(t *testing.T) {
	result := textbot.Result{
		ID: 1243,
		MusicRecord: types.MusicRecord{
			Band: types.Band{
				Name: "Some Band",
			},
			Duration: 100 * time.Second,
			Name:     "Some track",
			RecordId: "some-id",
			Source:   "some-source",
			Url:      "",
		},
		Tags:       []string{"some", "tags"},
		IsNew:      true,
		Count:      42,
		Statistics: playbot.MusicRecordStatistics{},
	}
	printer := server.MatrixMusicRecordPrinter{}

	// test
	msg := printer.Print(result)

	// assertions
	assert.Equal(
		t,
		`<font color="#FFFF00">[+1243]</font> <font color="#009300">Some track | Some Band</font> <font color="#0000FC">1m40s</font> <font color="#FC7F00">#some #tags</font> <font color="#7F7F7F">[42 résultats]</font>`,
		msg,
	)
}

func TestMatrixStatisticsPrinter_Print(t *testing.T) {
	/*
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
		printer := server.MatrixMusicRecordPrinter{Location: time.FixedZone("", +2*60*60)}

		// test
		result := printer.Print(stats)

		// assertions
		assert.Equal(t,
			"Posté la 1re fois par Abitbol le 20-11-2023 01:22:25 sur #bigphatsubwoofer. Posté 10 fois par 3 personnes sur 5 channels. George l'a posté 2 fois. Sauvegardé dans ses favoris par 15 personnes.",
			result,
		)
	*/
}
