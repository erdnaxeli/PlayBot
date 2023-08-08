package mariadb_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/repository/mariadb"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type searchResult struct {
	id          int64
	musicRecord types.MusicRecord
}

func assertEqualSearchResults(
	t *testing.T, expected []searchResult, actual []playbot.SearchResult,
) {
	require.Len(t, actual, len(expected))
	sort.Slice(actual, func(i, j int) bool { return actual[i].Id() < actual[j].Id() })

	for idx := range expected {
		assert.Equal(t, expected[idx].id, actual[idx].Id())
		assert.Equal(t, expected[idx].musicRecord, actual[idx].MusicRecord())
	}
}

func getTestRepository(t *testing.T) playbot.Repository {
	r, err := mariadb.New("test:test@(localhost)/test")
	require.Nil(
		t,
		err,
		"A MariaDB server must be listening on localhost with a user 'test', a password 'test' and a database 'test' initialized with the test-db.sql file.",
	)
	return r
}

func getMusicPost() types.MusicPost {
	var post types.MusicPost
	_ = gofakeit.Struct(&post)
	post.MusicRecord.Duration, _ = time.ParseDuration("1m35s")

	return post
}

func TestSearchMusicRecord_all(t *testing.T) {
	// setup

	r := getTestRepository(t)
	channel1 := gofakeit.DomainName()
	channel2 := gofakeit.DomainName()

	// A post in channel1 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	post1 := getMusicPost()
	post1.Channel.Name = channel1
	post1.MusicRecord.Name = "1 - La Classe Américaine"
	post1.MusicRecord.Band.Name = "George Abitbol"
	post1RecordId, isNew, err := r.SaveMusicPost(post1)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post1RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - word "class" but not "bol"
	// - tags "movie" and "french"
	post2 := getMusicPost()
	post2.Channel.Name = channel1
	post2.MusicRecord.Name = "2 - La Classe Américaine"
	post2.MusicRecord.Band.Name = "George Brassens"
	post2RecordId, isNew, err := r.SaveMusicPost(post2)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post2RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	// - tag "movie" but not "french"
	post3 := getMusicPost()
	post3.Channel.Name = channel1
	post3.MusicRecord.Name = "3 - La Classe Américaine"
	post3.MusicRecord.Band.Name = "George Abitbol"
	post3RecordId, isNew, err := r.SaveMusicPost(post3)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post3RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	post4 := getMusicPost()
	post4.Channel.Name = channel1
	post4.MusicRecord.Name = "4 - George Abitbol, a memorial"
	post4.MusicRecord.Band.Name = "The American Class Fan Club"
	post4RecordId, isNew, err := r.SaveMusicPost(post4)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post4RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	post5 := getMusicPost()
	post5.Channel.Name = channel2
	post5.MusicRecord.Name = "5 - La Classe Américaine"
	post5.MusicRecord.Band.Name = "George Abitbol"
	post5RecordId, isNew, err := r.SaveMusicPost(post5)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post5RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - word "class" but not "bol"
	// - tags "movie" and "french"
	post6 := getMusicPost()
	post6.Channel.Name = channel2
	post6.MusicRecord.Name = "6 - La Classe Américaine"
	post6.MusicRecord.Band.Name = "George Brassens"
	post6RecordId, isNew, err := r.SaveMusicPost(post6)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post6RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	// - tag "movie" but not "french"
	post7 := getMusicPost()
	post7.Channel.Name = channel2
	post7.MusicRecord.Name = "7 - La Classe Américaine"
	post7.MusicRecord.Band.Name = "George Abitbol"
	post7RecordId, isNew, err := r.SaveMusicPost(post7)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post7RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	post8 := getMusicPost()
	post8.Channel.Name = channel2
	post8.MusicRecord.Name = "8 - George Abitbol, a memorial"
	post8.MusicRecord.Band.Name = "The American Class Fan Club"
	post8RecordId, isNew, err := r.SaveMusicPost(post8)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post8RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// test

	// A query to search music records in channel1 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	count, ch, err := r.SearchMusicRecord(
		context.Background(),
		types.Channel{Name: channel1},
		[]string{"class", "bol"},
		[]string{"movie", "french"},
	)
	require.Nil(t, err)

	// assertions

	assert.EqualValues(t, 2, count)
	// get all results
	var results []playbot.SearchResult
	for r := range ch {
		results = append(results, r)
	}
	assertEqualSearchResults(t, []searchResult{
		{id: post1RecordId, musicRecord: post1.MusicRecord},
		{id: post4RecordId, musicRecord: post4.MusicRecord},
	},
		results,
	)
}

func TestSearchMusicRecord_noChannel(t *testing.T) {
	// setup

	r := getTestRepository(t)
	channel1 := gofakeit.DomainName()
	channel2 := gofakeit.DomainName()
	tag1 := gofakeit.Noun()
	tag2 := gofakeit.Noun()

	// A post in channel1 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	post1 := getMusicPost()
	post1.Channel.Name = channel1
	post1.MusicRecord.Name = "1 - La Classe Américaine"
	post1.MusicRecord.Band.Name = "George Abitbol"
	post1RecordId, isNew, err := r.SaveMusicPost(post1)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post1RecordId, []string{tag1, tag2, "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - word "class" but not "bol"
	// - tags tag1 and tag2
	post2 := getMusicPost()
	post2.Channel.Name = channel1
	post2.MusicRecord.Name = "2 - La Classe Américaine"
	post2.MusicRecord.Band.Name = "George Brassens"
	post2RecordId, isNew, err := r.SaveMusicPost(post2)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post2RecordId, []string{tag1, tag2, "wrong"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	// - tag tag1 but not tag2
	post3 := getMusicPost()
	post3.Channel.Name = channel1
	post3.MusicRecord.Name = "3 - La Classe Américaine"
	post3.MusicRecord.Band.Name = "George Abitbol"
	post3RecordId, isNew, err := r.SaveMusicPost(post3)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post3RecordId, []string{tag1, "english", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	// - tags tag1 and tag2
	post4 := getMusicPost()
	post4.Channel.Name = channel1
	post4.MusicRecord.Name = "4 - George Abitbol, a memorial"
	post4.MusicRecord.Band.Name = "The American Class Fan Club"
	post4RecordId, isNew, err := r.SaveMusicPost(post4)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post4RecordId, []string{tag1, tag2, "documentary"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	// - tags tag1 and tag2
	post5 := getMusicPost()
	post5.Channel.Name = channel2
	post5.MusicRecord.Name = "5 - La Classe Américaine"
	post5.MusicRecord.Band.Name = "George Abitbol"
	post5RecordId, isNew, err := r.SaveMusicPost(post5)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post5RecordId, []string{tag1, tag2, "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - word "class" but not "bol"
	// - tags tag1 and tag2
	post6 := getMusicPost()
	post6.Channel.Name = channel2
	post6.MusicRecord.Name = "6 - La Classe Américaine"
	post6.MusicRecord.Band.Name = "George Brassens"
	post6RecordId, isNew, err := r.SaveMusicPost(post6)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post6RecordId, []string{tag1, tag2, "wrong"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	// - tag tag1 but not tag2
	post7 := getMusicPost()
	post7.Channel.Name = channel2
	post7.MusicRecord.Name = "7 - La Classe Américaine"
	post7.MusicRecord.Band.Name = "George Abitbol"
	post7RecordId, isNew, err := r.SaveMusicPost(post7)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post7RecordId, []string{tag1, "english", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	// - tags tag1 and tag2
	post8 := getMusicPost()
	post8.Channel.Name = channel2
	post8.MusicRecord.Name = "8 - George Abitbol, a memorial"
	post8.MusicRecord.Band.Name = "The American Class Fan Club"
	post8RecordId, isNew, err := r.SaveMusicPost(post8)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post8RecordId, []string{tag1, tag2, "documentary"})
	require.Nil(t, err)

	// test

	// A query to search music records in channel1 matching:
	// - words "class" and "bol"
	// - tags tag1 and tag2
	count, ch, err := r.SearchMusicRecord(
		context.Background(),
		types.Channel{},
		[]string{"class", "bol"},
		[]string{tag1, tag2},
	)
	require.Nil(t, err)

	// assertions

	assert.EqualValues(t, 4, count)
	// get all results
	var results []playbot.SearchResult
	for r := range ch {
		results = append(results, r)
	}
	assertEqualSearchResults(
		t,
		[]searchResult{
			{id: post1RecordId, musicRecord: post1.MusicRecord},
			{id: post4RecordId, musicRecord: post4.MusicRecord},
			{id: post5RecordId, musicRecord: post5.MusicRecord},
			{id: post8RecordId, musicRecord: post8.MusicRecord},
		},
		results,
	)
}

func TestSearchMusicRecord_noWords(t *testing.T) {
	// setup

	r := getTestRepository(t)
	channel1 := gofakeit.DomainName()
	channel2 := gofakeit.DomainName()

	// A post in channel1 matching:
	// - tags "movie" and "french"
	post1 := getMusicPost()
	post1.Channel.Name = channel1
	post1.MusicRecord.Name = "1 - La Classe Américaine"
	post1.MusicRecord.Band.Name = "George Abitbol"
	post1RecordId, isNew, err := r.SaveMusicPost(post1)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post1RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - tags "movie" and "french"
	post2 := getMusicPost()
	post2.Channel.Name = channel1
	post2.MusicRecord.Name = "2 - La Classe Américaine"
	post2.MusicRecord.Band.Name = "George Brassens"
	post2RecordId, isNew, err := r.SaveMusicPost(post2)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post2RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - tag "movie" but not "french"
	post3 := getMusicPost()
	post3.Channel.Name = channel1
	post3.MusicRecord.Name = "3 - La Classe Américaine"
	post3.MusicRecord.Band.Name = "George Abitbol"
	post3RecordId, isNew, err := r.SaveMusicPost(post3)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post3RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - tags "movie" and "french"
	post4 := getMusicPost()
	post4.Channel.Name = channel1
	post4.MusicRecord.Name = "4 - George Abitbol, a memorial"
	post4.MusicRecord.Band.Name = "The American Class Fan Club"
	post4RecordId, isNew, err := r.SaveMusicPost(post4)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post4RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - tags "movie" and "french"
	post5 := getMusicPost()
	post5.Channel.Name = channel2
	post5.MusicRecord.Name = "5 - La Classe Américaine"
	post5.MusicRecord.Band.Name = "George Abitbol"
	post5RecordId, isNew, err := r.SaveMusicPost(post5)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post5RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - tags "movie" and "french"
	post6 := getMusicPost()
	post6.Channel.Name = channel2
	post6.MusicRecord.Name = "6 - La Classe Américaine"
	post6.MusicRecord.Band.Name = "George Brassens"
	post6RecordId, isNew, err := r.SaveMusicPost(post6)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post6RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - tag "movie" but not "french"
	post7 := getMusicPost()
	post7.Channel.Name = channel2
	post7.MusicRecord.Name = "7 - La Classe Américaine"
	post7.MusicRecord.Band.Name = "George Abitbol"
	post7RecordId, isNew, err := r.SaveMusicPost(post7)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post7RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - tags "movie" and "french"
	post8 := getMusicPost()
	post8.Channel.Name = channel2
	post8.MusicRecord.Name = "8 - George Abitbol, a memorial"
	post8.MusicRecord.Band.Name = "The American Class Fan Club"
	post8RecordId, isNew, err := r.SaveMusicPost(post8)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post8RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// test

	// A query to search music records in channel1 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	count, ch, err := r.SearchMusicRecord(
		context.Background(),
		types.Channel{Name: channel1},
		[]string{},
		[]string{"movie", "french"},
	)
	require.Nil(t, err)

	// assertions

	assert.EqualValues(t, 3, count)
	// get all results
	var results []playbot.SearchResult
	for r := range ch {
		results = append(results, r)
	}
	assertEqualSearchResults(
		t,
		[]searchResult{
			{id: post1RecordId, musicRecord: post1.MusicRecord},
			{id: post2RecordId, musicRecord: post2.MusicRecord},
			{id: post4RecordId, musicRecord: post4.MusicRecord},
		},
		results,
	)
}

func TestSearchMusicRecord_noTags(t *testing.T) {
	// setup

	r := getTestRepository(t)
	channel1 := gofakeit.DomainName()
	channel2 := gofakeit.DomainName()

	// A post in channel1 matching:
	// - words "class" and "bol"
	post1 := getMusicPost()
	post1.Channel.Name = channel1
	post1.MusicRecord.Name = "1 - La Classe Américaine"
	post1.MusicRecord.Band.Name = "George Abitbol"
	post1RecordId, isNew, err := r.SaveMusicPost(post1)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post1RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - word "class" but not "bol"
	post2 := getMusicPost()
	post2.Channel.Name = channel1
	post2.MusicRecord.Name = "2 - La Classe Américaine"
	post2.MusicRecord.Band.Name = "George Brassens"
	post2RecordId, isNew, err := r.SaveMusicPost(post2)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post2RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	post3 := getMusicPost()
	post3.Channel.Name = channel1
	post3.MusicRecord.Name = "3 - La Classe Américaine"
	post3.MusicRecord.Band.Name = "George Abitbol"
	post3RecordId, isNew, err := r.SaveMusicPost(post3)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post3RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	post4 := getMusicPost()
	post4.Channel.Name = channel1
	post4.MusicRecord.Name = "4 - George Abitbol, a memorial"
	post4.MusicRecord.Band.Name = "The American Class Fan Club"
	post4RecordId, isNew, err := r.SaveMusicPost(post4)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post4RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	post5 := getMusicPost()
	post5.Channel.Name = channel2
	post5.MusicRecord.Name = "5 - La Classe Américaine"
	post5.MusicRecord.Band.Name = "George Abitbol"
	post5RecordId, isNew, err := r.SaveMusicPost(post5)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post5RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - word "class" but not "bol"
	post6 := getMusicPost()
	post6.Channel.Name = channel2
	post6.MusicRecord.Name = "6 - La Classe Américaine"
	post6.MusicRecord.Band.Name = "George Brassens"
	post6RecordId, isNew, err := r.SaveMusicPost(post6)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post6RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	post7 := getMusicPost()
	post7.Channel.Name = channel2
	post7.MusicRecord.Name = "7 - La Classe Américaine"
	post7.MusicRecord.Band.Name = "George Abitbol"
	post7RecordId, isNew, err := r.SaveMusicPost(post7)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post7RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	post8 := getMusicPost()
	post8.Channel.Name = channel2
	post8.MusicRecord.Name = "8 - George Abitbol, a memorial"
	post8.MusicRecord.Band.Name = "The American Class Fan Club"
	post8RecordId, isNew, err := r.SaveMusicPost(post8)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post8RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// test

	// A query to search music records in channel1 matching:
	// - words "class" and "bol"
	// - tags "movie" and "french"
	count, ch, err := r.SearchMusicRecord(
		context.Background(),
		types.Channel{Name: channel1},
		[]string{"class", "bol"},
		[]string{},
	)
	require.Nil(t, err)

	// assertions

	assert.EqualValues(t, 3, count)
	// get all results
	var results []playbot.SearchResult
	for r := range ch {
		results = append(results, r)
	}
	assertEqualSearchResults(
		t,
		[]searchResult{
			{id: post1RecordId, musicRecord: post1.MusicRecord},
			{id: post3RecordId, musicRecord: post3.MusicRecord},
			{id: post4RecordId, musicRecord: post4.MusicRecord},
		},
		results,
	)
}

func TestSearchMusicRecord_noWords_noTags(t *testing.T) {
	// setup

	r := getTestRepository(t)
	channel1 := gofakeit.DomainName()
	channel2 := gofakeit.DomainName()

	// A post in channel1 matching:
	// - words "class" and "bol"
	post1 := getMusicPost()
	post1.Channel.Name = channel1
	post1.MusicRecord.Name = "1 - La Classe Américaine"
	post1.MusicRecord.Band.Name = "George Abitbol"
	post1RecordId, isNew, err := r.SaveMusicPost(post1)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post1RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - word "class" but not "bol"
	post2 := getMusicPost()
	post2.Channel.Name = channel1
	post2.MusicRecord.Name = "2 - La Classe Américaine"
	post2.MusicRecord.Band.Name = "George Brassens"
	post2RecordId, isNew, err := r.SaveMusicPost(post2)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post2RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	post3 := getMusicPost()
	post3.Channel.Name = channel1
	post3.MusicRecord.Name = "3 - La Classe Américaine"
	post3.MusicRecord.Band.Name = "George Abitbol"
	post3RecordId, isNew, err := r.SaveMusicPost(post3)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post3RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel1 matching:
	// - words "class" and "bol"
	post4 := getMusicPost()
	post4.Channel.Name = channel1
	post4.MusicRecord.Name = "4 - George Abitbol, a memorial"
	post4.MusicRecord.Band.Name = "The American Class Fan Club"
	post4RecordId, isNew, err := r.SaveMusicPost(post4)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post4RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	post5 := getMusicPost()
	post5.Channel.Name = channel2
	post5.MusicRecord.Name = "5 - La Classe Américaine"
	post5.MusicRecord.Band.Name = "George Abitbol"
	post5RecordId, isNew, err := r.SaveMusicPost(post5)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post5RecordId, []string{"movie", "french", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - word "class" but not "bol"
	post6 := getMusicPost()
	post6.Channel.Name = channel2
	post6.MusicRecord.Name = "6 - La Classe Américaine"
	post6.MusicRecord.Band.Name = "George Brassens"
	post6RecordId, isNew, err := r.SaveMusicPost(post6)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post6RecordId, []string{"movie", "french", "wrong"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	post7 := getMusicPost()
	post7.Channel.Name = channel2
	post7.MusicRecord.Name = "7 - La Classe Américaine"
	post7.MusicRecord.Band.Name = "George Abitbol"
	post7RecordId, isNew, err := r.SaveMusicPost(post7)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post7RecordId, []string{"movie", "english", "classic"})
	require.Nil(t, err)

	// A post in channel2 matching:
	// - words "class" and "bol"
	post8 := getMusicPost()
	post8.Channel.Name = channel2
	post8.MusicRecord.Name = "8 - George Abitbol, a memorial"
	post8.MusicRecord.Band.Name = "The American Class Fan Club"
	post8RecordId, isNew, err := r.SaveMusicPost(post8)
	require.True(t, isNew)
	require.Nil(t, err)
	err = r.SaveTags(post8RecordId, []string{"movie", "french", "documentary"})
	require.Nil(t, err)

	// test

	// A query to search music records in any channel
	ctx, cancel := context.WithCancel(context.Background())
	_, _, err = r.SearchMusicRecord(
		ctx,
		types.Channel{},
		[]string{},
		[]string{},
	)
	cancel()

	// assertions

	assert.Nil(t, err)
}

func TestSearchMusicRecord_noResult(t *testing.T) {
	// setup
	r := getTestRepository(t)

	// test
	count, ch, err := r.SearchMusicRecord(
		context.Background(),
		types.Channel{Name: gofakeit.DomainName()},
		[]string{},
		[]string{},
	)
	require.Nil(t, err)

	// assertions

	assert.EqualValues(t, 0, count)
	// get all results
	var results []playbot.SearchResult
	for r := range ch {
		results = append(results, r)
	}

	assert.Empty(t, results)
}

func TestSearchMusicRecord_contextDone(t *testing.T) {
	// setup

	r := getTestRepository(t)
	channel := gofakeit.DomainName()

	// create 10 posts in channel
	for i := 0; i < 10; i++ {
		post := getMusicPost()
		post.Channel.Name = channel
		_, isNew, err := r.SaveMusicPost(post)
		require.True(t, isNew)
		require.Nil(t, err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// test

	count, ch, err := r.SearchMusicRecord(
		ctx,
		types.Channel{Name: channel},
		[]string{},
		[]string{},
	)
	require.Nil(t, err)

	// assertions

	assert.EqualValues(t, 10, count)
	// get a result
	_, ok := <-ch
	assert.True(t, ok)
	// cancel the context
	cancel()
	// get another result
	_, ok = <-ch
	assert.False(t, ok)
}
