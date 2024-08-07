package mariadb

import (
	"database/sql"
	"sort"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertEqualRecordRow(t *testing.T, tx *sql.Tx, record types.MusicRecord, rowId int64) {
	row := tx.QueryRow(
		`
			select
				type,
				url,
				sender_irc,
				sender,
				title,
				duration,
				file,
				broken,
				channel,
				playlist,
				external_id
			from playbot
			where
				id = ?
		`,
		rowId,
	)
	var type_, url, title string
	var senderIrc, sender, file, channel, externalId sql.Null[string]
	var duration, broken, playlist int
	err := row.Scan(&type_, &url, &senderIrc, &sender, &title, &duration, &file, &broken, &channel, &playlist, &externalId)
	require.Nil(t, err)
	assert.Equal(t, record.Source, type_)
	assert.Equal(t, record.Url, url)
	assert.False(t, senderIrc.Valid)
	assert.True(t, sender.Valid)
	assert.Equal(t, record.Band.Name, sender.V)
	assert.Equal(t, record.Name, title)
	assert.Equal(t, int(record.Duration.Seconds()), duration)
	assert.False(t, file.Valid)
	assert.Equal(t, 0, broken)
	assert.False(t, channel.Valid)
	assert.Equal(t, 0, playlist)
	assert.True(t, externalId.Valid)
	assert.Equal(t, record.RecordId, externalId.V)
}

func getTestRepository(t *testing.T) mariaDbRepository {
	r, err := New("test", "test", "localhost", "test")
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

func rollback(tx *sql.Tx) {
	_ = tx.Rollback()
}

func TestSearchResult(t *testing.T) {
	var record types.MusicRecord
	_ = gofakeit.Struct(&record)

	result := searchResult{
		id:          42,
		musicRecord: record,
	}

	assert.Equal(t, result.id, result.Id())
	assert.Equal(t, result.musicRecord, result.MusicRecord())
}

func TestInsertOrUpdateMusicRecord_Insert(t *testing.T) {
	recordDuration, _ := time.ParseDuration("1m35s")
	record := types.MusicRecord{
		Band:     types.Band{Name: "TestBand"},
		Duration: recordDuration,
		Name:     "testName",
		RecordId: "testRecordId",
		Source:   "testSource",
		Url:      "testUrl",
	}
	r := getTestRepository(t)
	defer r.db.Close()
	tx, _ := r.db.Begin()
	defer func() { _ = tx.Rollback() }()

	recordId, isNew, err := r.insertOrUpdateMusicRecord(tx, record)

	require.Nil(t, err)
	assertEqualRecordRow(t, tx, record, recordId)
	assert.True(t, isNew)
}

func TestInsertOrUpdateMusicRecord_Update(t *testing.T) {
	recordDuration, _ := time.ParseDuration("1m35s")
	record := types.MusicRecord{
		Band:     types.Band{Name: "TestBand"},
		Duration: recordDuration,
		Name:     "testName",
		RecordId: "testRecordId",
		Source:   "testSource",
		Url:      "testUrl",
	}
	r := getTestRepository(t)
	defer r.db.Close()
	tx, _ := r.db.Begin()
	defer rollback(tx)
	recordId, isNew, err := r.insertOrUpdateMusicRecord(tx, record)
	require.Nil(t, err)
	require.True(t, isNew)

	record.Band.Name = "NewBand"
	record.Duration += 1
	record.Name = "NewName"
	record.RecordId = "NewRecordId"
	record.Source = "NewSource"

	newRecordId, isNew, err := r.insertOrUpdateMusicRecord(tx, record)

	require.Nil(t, err)
	assert.Equal(t, recordId, newRecordId)
	assertEqualRecordRow(t, tx, record, recordId)
	assert.False(t, isNew)
}

func TestGetTags_noTags(t *testing.T) {
	// setup
	r := getTestRepository(t)
	defer r.db.Close()

	// test
	tags, err := r.GetTags(1987654334)
	require.Nil(t, err)

	// assertions
	assert.Equal(t, []string{}, tags)
}

func TestGetTags_tags(t *testing.T) {
	// setup
	r := getTestRepository(t)
	defer r.db.Close()

	// test data
	musicPost := getMusicPost()
	tags := []string{"some", "tags", "to", "test"}
	recordId, isNew, err := r.SaveMusicPost(musicPost)
	require.Nil(t, err)
	require.True(t, isNew)
	err = r.SaveTags(recordId, tags)
	require.Nil(t, err)

	// test
	foundTags, err := r.GetTags(recordId)
	require.Nil(t, err)

	// assertions
	sort.Strings(tags)
	sort.Strings(foundTags)
	assert.Equal(t, tags, foundTags)
}

func TestSaveChannelPost_ok(t *testing.T) {
	// setup
	r := getTestRepository(t)
	defer r.db.Close()
	tx, _ := r.db.Begin()
	defer rollback(tx)
	var record types.MusicRecord
	_ = gofakeit.Struct(&record)
	record.Duration, _ = time.ParseDuration("1m35s")
	recordId, isNew, err := r.insertOrUpdateMusicRecord(tx, record)
	require.Nil(t, err)
	require.True(t, isNew)

	// data to test
	var person types.Person
	var channel types.Channel
	_ = gofakeit.Struct(&person)
	_ = gofakeit.Struct(&channel)

	// test
	err = r.saveChannelPost(tx, recordId, person, channel)
	require.Nil(t, err)

	// assertions
	rows, err := tx.Query(
		`
			select
				sender_irc,
				chan
			from playbot_chan
			where
				content = ?
		`,
		recordId,
	)
	require.Nil(t, err)
	require.True(t, rows.Next())
	var senderIrc, channelName string
	err = rows.Scan(&senderIrc, &channelName)
	require.Nil(t, err)
	assert.Equal(t, person.Name, senderIrc)
	assert.Equal(t, channel.Name, channelName)
	// assert there are no more rows
	assert.False(t, rows.Next())
	assert.Nil(t, rows.Err())
}

func TestSaveChannelPost_RecordNotFound(t *testing.T) {
	// setup
	r := getTestRepository(t)
	defer r.db.Close()
	tx, _ := r.db.Begin()
	defer rollback(tx)
	var person types.Person
	var channel types.Channel
	_ = gofakeit.Struct(&person)
	_ = gofakeit.Struct(&channel)
	recordID := gofakeit.Int64()

	// test
	err := r.saveChannelPost(tx, recordID, person, channel)

	// assertions
	assert.NotNil(t, err)
	row := tx.QueryRow(
		`
			select
				sender_irc,
				chan
			from playbot_chan
			where
				content = ?
		`,
		recordID,
	)
	var senderIrc, channelName string
	err = row.Scan(&senderIrc, &channelName)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSaveMusicRecord_once(t *testing.T) {
	// setup
	r := getTestRepository(t)
	defer r.db.Close()
	post := getMusicPost()

	// test
	recordId, isNew, err := r.SaveMusicPost(post)
	require.Nil(t, err)

	// assertions
	assert.True(t, isNew)
	tx, _ := r.db.Begin()
	defer rollback(tx)
	assertEqualRecordRow(t, tx, post.MusicRecord, recordId)
	rows, err := tx.Query(
		`
			select
				sender_irc,
				chan
			from playbot_chan
			where
				content = ?
		`,
		recordId,
	)
	require.Nil(t, err)
	require.True(t, rows.Next())
	var senderIrc, channelName string
	err = rows.Scan(&senderIrc, &channelName)
	require.Nil(t, err)
	assert.Equal(t, post.Person.Name, senderIrc)
	assert.Equal(t, post.Channel.Name, channelName)
	// assert there are no more rows
	require.False(t, rows.Next())
	assert.Nil(t, rows.Err())
}

func TestSaveMusicRecord_twice(t *testing.T) {
	// setup

	r := getTestRepository(t)
	defer r.db.Close()

	post := getMusicPost()

	// first post
	recordId, isNew, err := r.SaveMusicPost(post)
	require.Nil(t, err)
	require.True(t, isNew)

	// second post
	secondPost := post
	_ = gofakeit.Struct(&secondPost.Person)
	_ = gofakeit.Struct(&secondPost.Channel)

	// test

	secondRecordId, isNew, err := r.SaveMusicPost(secondPost)
	require.Nil(t, err)

	// assertions

	assert.False(t, isNew)
	assert.Equal(t, recordId, secondRecordId)

	tx, _ := r.db.Begin()
	defer rollback(tx)
	assertEqualRecordRow(t, tx, post.MusicRecord, recordId)

	var senderIrc, channelName string
	rows, err := tx.Query(
		`
			select
				sender_irc,
				chan
			from playbot_chan
			where
				content = ?
			order by date
		`,
		recordId,
	)
	require.Nil(t, err)

	// first row
	require.True(t, rows.Next())
	err = rows.Scan(&senderIrc, &channelName)
	require.Nil(t, err)
	assert.Equal(t, post.Person.Name, senderIrc)
	assert.Equal(t, post.Channel.Name, channelName)

	// second row
	require.True(t, rows.Next())
	err = rows.Scan(&senderIrc, &channelName)
	require.Nil(t, err)
	assert.Equal(t, secondPost.Person.Name, senderIrc)
	assert.Equal(t, secondPost.Channel.Name, channelName)

	// assert there are no more rows
	require.False(t, rows.Next())
	assert.Nil(t, rows.Err())
}

func TestSaveTags(t *testing.T) {
	// setup

	r := getTestRepository(t)
	defer r.db.Close()

	post := getMusicPost()
	recordId, isNew, err := r.SaveMusicPost(post)
	require.Nil(t, err)
	require.True(t, isNew)

	var tags []string
	gofakeit.Slice(&tags)

	// test

	err = r.SaveTags(recordId, tags)
	require.Nil(t, err)

	// assertions

	tx, _ := r.db.Begin()
	defer rollback(tx)
	rows, err := tx.Query(
		`
			select tag
			from playbot_tags
			where id = ?
		`,
		recordId,
	)
	require.Nil(t, err)

	var savedTags []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		require.Nil(t, err)
		savedTags = append(savedTags, tag)
	}

	sort.Strings(tags)
	sort.Strings(savedTags)
	assert.Equal(t, tags, savedTags)
}
