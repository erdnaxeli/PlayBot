package textbot_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

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
