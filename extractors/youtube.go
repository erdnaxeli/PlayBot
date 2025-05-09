package extractors

import (
	"context"
	"fmt"
	"regexp"

	"github.com/erdnaxeli/PlayBot/iso8601"
	"github.com/erdnaxeli/PlayBot/types"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YoutubeUnknownVideoError is the error when a video is unknown.
type YoutubeUnknownVideoError struct {
	// RecordID is the ID of the video.
	RecordID string
}

func (err YoutubeUnknownVideoError) Error() string {
	return fmt.Sprintf("unknow video %s", err.RecordID)
}

// YoutubeExtractor implements the MatcherExtractor interface for Youtube URLs.
type YoutubeExtractor struct {
	// APIKey is the Youtube API key used to make API calls.
	APIKey string
}

// Match returns the URL matched and the record ID.
func (*YoutubeExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(`(?:^|[^!])https?://(?:(?:www|music)\.youtube\.com/watch\?[a-zA-Z0-9_=&-]*v=|youtu.be/)([a-zA-Z0-9_-]+)`)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	return groups[0], groups[1]
}

// Extract returns a Youtube video data from the video ID.
func (e *YoutubeExtractor) Extract(recordID string) (types.MusicRecord, error) {
	youtube, err := youtube.NewService(context.Background(), option.WithAPIKey(e.APIKey))
	if err != nil {
		return types.MusicRecord{}, err
	}

	call := youtube.Videos.List([]string{"snippet", "contentDetails"}).Id(recordID)
	response, err := call.Do()
	if err != nil {
		return types.MusicRecord{}, err
	}

	if len(response.Items) == 0 {
		return types.MusicRecord{}, YoutubeUnknownVideoError{recordID}
	}

	video := response.Items[0]
	duration, err := iso8601.ParseDuration(video.ContentDetails.Duration)
	if err != nil {
		return types.MusicRecord{}, err
	}

	return types.MusicRecord{
		Band:     types.Band{Name: video.Snippet.ChannelTitle},
		Duration: duration,
		Name:     video.Snippet.Title,
		RecordID: recordID,
		Source:   "youtube",
		URL:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", recordID),
	}, nil
}
