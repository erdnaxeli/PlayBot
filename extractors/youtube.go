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

type YoutubeExtractor struct {
	ApiKey string
}

func (*YoutubeExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(`(?:^|[^!])https?://(?:(?:www|music).youtube.com/watch\?[a-zA-Z0-9_=&-]*v=|youtu.be/)([a-zA-Z0-9_-]+)`)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	return groups[0], groups[1]
}

func (e *YoutubeExtractor) Extract(recordId string) (types.MusicRecord, error) {
	youtube, err := youtube.NewService(context.Background(), option.WithAPIKey(e.ApiKey))
	if err != nil {
		return types.MusicRecord{}, err
	}

	call := youtube.Videos.List([]string{"snippet", "contentDetails"}).Id(recordId)
	response, err := call.Do()
	if err != nil {
		return types.MusicRecord{}, err
	}

	if len(response.Items) == 0 {
		return types.MusicRecord{}, fmt.Errorf("unknown video %s", recordId)
	}

	video := response.Items[0]
	return types.MusicRecord{
		Band:     types.Band{Name: video.Snippet.ChannelTitle},
		Duration: iso8601.ParseDuration(video.ContentDetails.Duration),
		Name:     video.Snippet.Title,
		RecordId: recordId,
		Source:   "youtube",
		Url:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", recordId),
	}, nil
}
