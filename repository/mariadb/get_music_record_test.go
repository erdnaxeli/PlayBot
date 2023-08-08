package mariadb_test

import (
	"testing"

	"github.com/erdnaxeli/PlayBot/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMusicRecord_noResult(t *testing.T) {
	r := getTestRepository(t)

	result, err := r.GetMusicRecord(9223372036854775807)
	require.Nil(t, err)

	assert.Equal(t, types.MusicRecord{}, result)
}

func TestGetMusicRecord_ok(t *testing.T) {
	var records []types.MusicRecord
	var recordIDs []int64
	r := getTestRepository(t)
	for i := 0; i < 10; i++ {
		post := getMusicPost()
		recordID, isNew, err := r.SaveMusicPost(post)
		require.Nil(t, err)
		require.True(t, isNew)

		records = append(records, post.MusicRecord)
		recordIDs = append(recordIDs, recordID)
	}

	result, err := r.GetMusicRecord(recordIDs[4])
	require.Nil(t, err)

	assert.Equal(t, records[4], result)
}
