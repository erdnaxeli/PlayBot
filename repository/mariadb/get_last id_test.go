package mariadb_test

import (
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/erdnaxeli/PlayBot/playbot"
	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLastID_ok(t *testing.T) {
	r := getTestRepository(t)
	var channel types.Channel
	_ = gofakeit.Struct(&channel)

	var ids []int64
	for i := 0; i < 100; i++ {
		post := getMusicPost()
		post.Channel = channel

		id, isNew, err := r.SaveMusicPost(post)
		require.True(t, isNew)
		require.Nil(t, err)

		ids = append(ids, id)
	}

	tests := []struct {
		offset int
		id     int64
	}{
		{0, ids[99]},
		{1, ids[98]},
		{3, ids[96]},
		{6, ids[93]},
		{10, ids[89]},
	}
	for _, test := range tests {
		t.Run(
			fmt.Sprint(test.offset),
			func(t *testing.T) {
				result, err := r.GetLastID(channel, test.offset)
				require.Nil(t, err)

				assert.Equal(t, test.id, result)
			},
		)
	}
}

func TestGetLastID_noRecordFound(t *testing.T) {
	r := getTestRepository(t)
	var channel types.Channel
	_ = gofakeit.Struct(&channel)

	result, err := r.GetLastID(channel, 0)

	assert.Equal(t, result, int64(0))
	assert.ErrorIs(t, err, playbot.NoRecordFoundError)
}
