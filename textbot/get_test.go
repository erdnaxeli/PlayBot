package textbot_test

import (
	"context"
	"testing"

	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/textbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type PlaybotMock struct {
	mock.Mock
	// Stores the calls to SearchMusicRecord. We ignore the first ctx argument.
	searchMusicRecordCalls []playbot.Search
}

func (m *PlaybotMock) GetTags(recordID int64) ([]string, error) {
	args := m.Called(recordID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *PlaybotMock) GetLastID(channel types.Channel, limit int) (int64, error) {
	args := m.Called(channel, limit)
	return args.Get(0).(int64), args.Error(1)
}

func (m *PlaybotMock) GetMusicRecord(recordID int64) (types.MusicRecord, error) {
	args := m.Called(recordID)
	return args.Get(0).(types.MusicRecord), args.Error(1)
}

func (m *PlaybotMock) GetMusicRecordStatistics(recordID int64) (playbot.MusicRecordStatistics, error) {
	args := m.Called(recordID)
	return args.Get(0).(playbot.MusicRecordStatistics), args.Error(1)
}

func (m *PlaybotMock) ParseAndSaveMusicRecord(url string, person types.Person, channel types.Channel) (int64, types.MusicRecord, bool, error) {
	args := m.Called(url, person, channel)
	return args.Get(0).(int64), args.Get(1).(types.MusicRecord), args.Bool(2), args.Error(3)
}

func (m *PlaybotMock) SaveFav(person string, recordID int64) error {
	args := m.Called(person, recordID)
	return args.Error(0)
}

func (m *PlaybotMock) SaveMusicPost(recordID int64, channel types.Channel, person types.Person) error {
	args := m.Called(recordID, channel, person)
	return args.Error(0)
}

func (m *PlaybotMock) SaveTags(recordID int64, tags []string) error {
	args := m.Called(recordID, tags)
	return args.Error(0)
}

func (m *PlaybotMock) SearchMusicRecord(ctx context.Context, search playbot.Search) (int64, playbot.SearchResult, error) {
	// search contains a context object, so we can check if it is the same as the one we expect
	m.searchMusicRecordCalls = append(m.searchMusicRecordCalls, search)
	return 42, SearchResult{
		id:          42,
		musicRecord: types.MusicRecord{Name: "Some music record"},
	}, nil
}

type SearchResult struct {
	id          int64
	musicRecord types.MusicRecord
}

func (s SearchResult) Id() int64 {
	return s.id
}

func (s SearchResult) MusicRecord() types.MusicRecord {
	return s.musicRecord
}

func TestGet(t *testing.T) {
	tests := []struct {
		msg          string
		all          bool
		id           int64
		tags         []string
		excludedTags []string
		words        []string
		searchById   bool
	}{
		{
			msg:          "!get",
			id:           0,
			all:          false,
			tags:         nil,
			excludedTags: nil,
			words:        nil,
			searchById:   false,
		},
		{
			msg:          "!get -a",
			id:           0,
			all:          true,
			tags:         nil,
			excludedTags: nil,
			words:        nil,
			searchById:   false,
		},
		{
			msg:          "!get -a some thing #else",
			id:           0,
			all:          true,
			tags:         []string{"else"},
			excludedTags: nil,
			words:        []string{"some", "thing"},
			searchById:   false,
		},
		{
			msg:          "!get some thing #else ##excluded",
			id:           0,
			all:          false,
			tags:         []string{"else"},
			excludedTags: []string{"excluded"},
			words:        []string{"some", "thing"},
			searchById:   false,
		},
		{
			msg:          "!get 42",
			id:           42,
			all:          false,
			tags:         nil,
			excludedTags: nil,
			words:        nil,
			searchById:   true,
		},
		{
			msg:          "!get 42 some thing #else ##excluded",
			id:           42,
			all:          false,
			tags:         nil,
			excludedTags: nil,
			words:        nil,
			searchById:   false,
		},
	}

	for _, test := range tests {
		t.Run(
			test.msg,
			func(t *testing.T) {
				// Given
				var recordId int64 = 42
				playbotMock := &PlaybotMock{}
				playbotMock.On("GetTags", recordId).Return([]string{"tag1", "tag2"}, nil)
				playbotMock.On("SaveMusicPost", recordId, types.Channel{Name: "channel"}, types.Person{Name: "PlayBot"}).Return(nil)

				if test.searchById {
					playbotMock.On("GetMusicRecord", test.id).Return(
						types.MusicRecord{Name: "Some music record", RecordId: "42"}, nil,
					)
				}

				textBot := textbot.New(playbotMock)

				// When
				_, isCmd, err := textBot.Execute("channel", "george", test.msg, "")
				require.Nil(t, err)

				// Then
				assert.True(t, isCmd)
				playbotMock.AssertExpectations(t)

				if test.id == 0 {
					assert.Len(t, playbotMock.searchMusicRecordCalls, 1)
					search := playbotMock.searchMusicRecordCalls[0]
					assert.Equal(t, types.Channel{Name: "channel"}, search.Channel)
					assert.Equal(t, test.all, search.GlobalSearch)
					assert.Equal(t, test.words, search.Words)
					assert.Equal(t, test.tags, search.Tags)
					assert.Equal(t, test.excludedTags, search.ExcludedTags)
				}
			},
		)
	}
}
