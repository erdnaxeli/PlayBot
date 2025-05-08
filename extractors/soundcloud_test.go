package extractors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSoundcloudMatch(t *testing.T) {
	tests := []struct {
		in         string
		matchedURL string
		recordID   string
	}{
		{"", "", ""},
		{
			"hello world",
			"",
			"",
		},
		{
			"https://soundcloud.com/hate_music/frederic-hate-podcast-332",
			"https://soundcloud.com/hate_music/frederic-hate-podcast-332",
			"hate_music/frederic-hate-podcast-332",
		},
		{
			"https://soundcloud.com/hate_music/frederic-hate-podcast-332?set=toto",
			"https://soundcloud.com/hate_music/frederic-hate-podcast-332?set=toto",
			"hate_music/frederic-hate-podcast-332",
		},
		{
			"!https://soundcloud.com/hate_music/frederic-hate-podcast-332",
			"",
			"",
		},
		{
			"hello world https://soundcloud.com/hate_music/frederic-hate-podcast-332 #techno #mix",
			"https://soundcloud.com/hate_music/frederic-hate-podcast-332",
			"hate_music/frederic-hate-podcast-332",
		},
		{
			"hello world https://soundcloud.com/hate_music/frederic-hate-podcast-332 #techno #mix",
			"https://soundcloud.com/hate_music/frederic-hate-podcast-332",
			"hate_music/frederic-hate-podcast-332",
		},
		{
			"hello world !https://soundcloud.com/hate_music/frederic-hate-podcast-332",
			"",
			"",
		},
	}

	for _, test := range tests {
		t.Run(
			test.in,
			func(t *testing.T) {
				matchedURL, recordID := SoundCloudExtractor{}.Match(test.in)
				assert.Equal(t, test.matchedURL, matchedURL)
				assert.Equal(t, test.recordID, recordID)
			},
		)
	}
}
