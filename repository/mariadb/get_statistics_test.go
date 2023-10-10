package mariadb_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getMusicRecord() types.MusicRecord {
	var record types.MusicRecord
	_ = gofakeit.Struct(&record)
	record.Duration = time.Second * time.Duration(rand.Int63n(300))

	return record
}

func TestGetStatistics_noRecord(t *testing.T) {
	r := getTestRepository(t)

	stats, err := r.GetMusicRecordStatistics(-1)
	require.Nil(t, err)

	assert.Nil(t, err)
	assert.Equal(t, playbot.MusicRecordStatistics{}, stats)
}

func TestGetStatistics_ok(t *testing.T) {
	r := getTestRepository(t)

	// two records
	record1 := getMusicRecord()
	record2 := getMusicRecord()

	// three people
	var person1 types.Person
	var person2 types.Person
	var person3 types.Person
	_ = gofakeit.Struct(&person1)
	_ = gofakeit.Struct(&person2)
	_ = gofakeit.Struct(&person3)

	// three channels
	var channel1 types.Channel
	var channel2 types.Channel
	var channel3 types.Channel
	_ = gofakeit.Struct(&channel1)
	_ = gofakeit.Struct(&channel2)
	_ = gofakeit.Struct(&channel3)

	// posts

	// record1, person1, channel1
	recordID, isNew, err := r.SaveMusicPost(types.MusicPost{
		MusicRecord: record1,
		Person:      person1,
		Channel:     channel1,
	})
	require.Nil(t, err)
	require.True(t, isNew)
	// record1, person2, channel2
	_, isNew, err = r.SaveMusicPost(types.MusicPost{
		MusicRecord: record1,
		Person:      person2,
		Channel:     channel2,
	})
	require.Nil(t, err)
	require.False(t, isNew)
	// record1, person2, channel3
	_, isNew, err = r.SaveMusicPost(types.MusicPost{
		MusicRecord: record1,
		Person:      person2,
		Channel:     channel3,
	})
	require.Nil(t, err)
	require.False(t, isNew)
	// record2, person3, channel3
	_, isNew, err = r.SaveMusicPost(types.MusicPost{
		MusicRecord: record2,
		Person:      person3,
		Channel:     channel3,
	})
	require.Nil(t, err)
	require.True(t, isNew)
	// record2, person1, channel1
	_, isNew, err = r.SaveMusicPost(types.MusicPost{
		MusicRecord: record2,
		Person:      person1,
		Channel:     channel1,
	})
	require.Nil(t, err)
	require.False(t, isNew)
	// record1, person3, channel2
	_, isNew, err = r.SaveMusicPost(types.MusicPost{
		MusicRecord: record1,
		Person:      person3,
		Channel:     channel2,
	})
	require.Nil(t, err)
	require.False(t, isNew)

	// test
	stats, err := r.GetMusicRecordStatistics(recordID)
	require.Nil(t, err)

	// assertions
	assert.Equal(
		t,
		playbot.MusicRecordStatistics{
			PostsCount:       4,
			PeopleCount:      3,
			ChannelsCount:    3,
			MaxPerson:        person2,
			MaxPersonCount:   2,
			MaxChannel:       channel2,
			MaxChannelCount:  2,
			FirstPostPerson:  person1,
			FirstPostChannel: channel1,
			FirstPostDate:    stats.FirstPostDate, // no assertion for this field
			FavoritesCount:   0,
		},
		stats,
	)
}
