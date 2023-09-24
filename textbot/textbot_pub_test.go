package textbot_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type TestPlayBot struct {
	mock.Mock
}

func (tp *TestPlayBot) GetTags(recordID int64) ([]string, error) {
	args := tp.Called(recordID)
	return args.Get(0).([]string), args.Error(1)
}

func (tp *TestPlayBot) GetLastID(channel types.Channel, offset int) (int64, error) {
	args := tp.Called(channel, offset)
	return args.Get(0).(int64), args.Error(1)
}

func (tp *TestPlayBot) GetMusicRecord(recordID int64) (types.MusicRecord, error) {
	args := tp.Called(recordID)
	return args.Get(0).(types.MusicRecord), args.Error(1)
}

func (tp *TestPlayBot) GetMusicRecordStatistics(recordID int64) (playbot.MusicRecordStatistics, error) {
	args := tp.Called(recordID)
	return args.Get(0).(playbot.MusicRecordStatistics), args.Error(1)
}

func (tp *TestPlayBot) ParseAndSaveMusicRecord(
	msg string, person types.Person, channel types.Channel,
) (int64, types.MusicRecord, bool, error) {
	args := tp.Called(msg, person, channel)
	return args.Get(0).(int64), args.Get(1).(types.MusicRecord), args.Bool(2), args.Error(3)
}

func (tp *TestPlayBot) SaveFav(user string, recordID int64) error {
	args := tp.Called(user, recordID)
	return args.Error(0)
}

func (tp *TestPlayBot) SaveMusicPost(
	recordID int64, channel types.Channel, person types.Person,
) error {
	args := tp.Called(recordID, channel, person)
	return args.Error(0)
}

func (tp *TestPlayBot) SaveTags(recordID int64, tags []string) error {
	args := tp.Called(recordID, tags)
	return args.Error(0)
}

func (tp *TestPlayBot) SearchMusicRecord(
	ctx context.Context, search playbot.Search,
) (int64, playbot.SearchResult, error) {
	args := tp.Called(ctx, search)
	return args.Get(0).(int64), args.Get(1).(playbot.SearchResult), args.Error(2)
}

func TestExecute_fav_id(t *testing.T) {
	// setup
	channel := "#channel"
	recordID := int64(42)
	msg := fmt.Sprintf("!fav %d", recordID)
	nick := "someNick"
	user := "someUser"

	tp := TestPlayBot{}
	tp.On(
		"ParseAndSaveMusicRecord",
		fmt.Sprint(recordID),
		types.Person{Name: nick},
		types.Channel{Name: channel},
	).Return(int64(0), types.MusicRecord{}, false, nil)
	tp.On("SaveFav", user, recordID).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel,
		nick,
		msg,
		user,
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}

func TestExecute_fav_offset(t *testing.T) {
	// setup
	channel := "#channel"
	offset := -2
	recordID := int64(42)
	msg := fmt.Sprintf("!fav %d", offset)
	nick := "someNick"
	user := "someUser"

	tp := TestPlayBot{}
	tp.On(
		"ParseAndSaveMusicRecord",
		fmt.Sprint(offset),
		types.Person{Name: nick},
		types.Channel{Name: channel},
	).Return(int64(0), types.MusicRecord{}, false, nil)
	tp.On("GetLastID", types.Channel{Name: channel}, offset).Return(recordID, nil)
	tp.On("SaveFav", user, recordID).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel,
		nick,
		msg,
		user,
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}

func TestExecute_fav_noIDnorOffset(t *testing.T) {
	// setup
	channel := "#channel"
	recordID := int64(42)
	msg := "!fav"
	nick := "someNick"
	user := "someUser"

	tp := TestPlayBot{}
	tp.On("GetLastID", types.Channel{Name: channel}, 0).Return(recordID, nil)
	tp.On("SaveFav", user, recordID).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel,
		nick,
		msg,
		user,
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}

func TestExecute_fav_musicRecord(t *testing.T) {
	// setup
	channel := "#channel"
	recordID := int64(42)
	newTags := []string{"tag1", "tag2"}
	tags := []string{"tag1", "tag2", "tag3"}
	// The given id "1" will be ignored.
	msg := "1 https://someURL #tag1 #tag2"
	nick := "someNick"
	user := "someUser"
	var musicRecord types.MusicRecord
	_ = gofakeit.Struct(&musicRecord)
	isNew := true

	tp := TestPlayBot{}
	tp.On(
		"ParseAndSaveMusicRecord",
		msg,
		types.Person{Name: nick},
		types.Channel{Name: channel},
	).Return(recordID, musicRecord, isNew, nil)
	tp.On("SaveTags", recordID, newTags).Return(nil)
	tp.On("GetTags", recordID).Return(tags, nil)
	tp.On("SaveFav", user, recordID).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel,
		nick,
		fmt.Sprintf("!fav %s", msg),
		user,
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(
		t,
		textbot.Result{
			ID:          recordID,
			MusicRecord: musicRecord,
			Tags:        tags,
			IsNew:       isNew,
		},
		result,
	)
}

func TestExecucet_fav_noUser(t *testing.T) {
	// setup
	channel := "#channel"
	msg := "!fav"
	nick := "someNick"

	tp := TestPlayBot{}

	tb := textbot.New(&tp)
	tp.On("GetLastID", types.Channel{Name: channel}, 0).Return(int64(42), nil)

	// test
	_, _, err := tb.Execute(
		channel,
		nick,
		msg,
		"",
	)

	// assertions
	assert.ErrorIs(t, err, textbot.AuthenticationRequired)
	tp.AssertExpectations(t)
}

func TestExecute_fav_noUser_musicRecord(t *testing.T) {
	// setup
	channel := "#channel"
	recordID := int64(42)
	newTags := []string{"tag1", "tag2"}
	tags := []string{"tag1", "tag2", "tag3"}
	// The given id "1" will be ignored.
	msg := "1 https://someURL #tag1 #tag2"
	nick := "someNick"
	var musicRecord types.MusicRecord
	_ = gofakeit.Struct(&musicRecord)
	isNew := true

	tp := TestPlayBot{}
	tp.On(
		"ParseAndSaveMusicRecord",
		msg,
		types.Person{Name: nick},
		types.Channel{Name: channel},
	).Return(recordID, musicRecord, isNew, nil)
	tp.On("SaveTags", recordID, newTags).Return(nil)
	tp.On("GetTags", recordID).Return(tags, nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel,
		nick,
		fmt.Sprintf("!fav %s", msg),
		"",
	)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(
		t,
		textbot.Result{
			ID:          recordID,
			MusicRecord: musicRecord,
			Tags:        tags,
			IsNew:       isNew,
		},
		result,
	)
	assert.ErrorIs(t, err, textbot.AuthenticationRequired)
}

func TestExecute_saveTagsCmd_tagsWithoutHash(t *testing.T) {
	// setup
	recordID := int64(42)
	tags := []string{"some", "tags", "to", "test"}

	tp := TestPlayBot{}
	tp.On("SaveTags", recordID, tags).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		"#channel", "someUser",
		fmt.Sprintf("!tag %d %s", recordID, strings.Join(tags, " ")),
		"",
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}

func TestExecute_saveTagsCmd_id(t *testing.T) {
	// setup
	recordID := int64(42)
	tags := []string{"#some", "#tags", "#to", "#test"}

	tp := TestPlayBot{}
	tp.On("SaveTags", recordID, []string{"some", "tags", "to", "test"}).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		"#channel", "someUser",
		fmt.Sprintf("!tag %d %s", recordID, strings.Join(tags, " ")),
		"",
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}

func TestExecute_saveTagsCmd_offset(t *testing.T) {
	// setup
	channel := types.Channel{Name: "#channel"}
	offset := -2
	recordID := int64(42)
	tags := []string{"#some", "#tags", "#to", "#test"}

	tp := TestPlayBot{}
	tp.On("GetLastID", channel, offset).Return(recordID, nil)
	tp.On("SaveTags", recordID, []string{"some", "tags", "to", "test"}).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel.Name, "someUser",
		fmt.Sprintf("!tag %d %s", offset, strings.Join(tags, " ")),
		"",
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}

func TestExecute_saveTagsCmd_noIDnorOffset(t *testing.T) {
	// setup
	channel := types.Channel{Name: "#channel"}
	recordID := int64(42)
	tags := []string{"#some", "#tags", "#to", "#test"}

	tp := TestPlayBot{}
	tp.On("GetLastID", channel, 0).Return(recordID, nil)
	tp.On("SaveTags", recordID, []string{"some", "tags", "to", "test"}).Return(nil)

	tb := textbot.New(&tp)

	// test
	result, ok, err := tb.Execute(
		channel.Name, "someUser",
		fmt.Sprintf("!tag %s", strings.Join(tags, " ")),
		"",
	)
	require.Nil(t, err)

	// assertions
	tp.AssertExpectations(t)
	assert.True(t, ok)
	assert.Equal(t, textbot.Result{}, result)
}
