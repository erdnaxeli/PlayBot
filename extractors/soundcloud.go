package extractors

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/erdnaxeli/PlayBot/iso8601"
	"github.com/erdnaxeli/PlayBot/types"
	"golang.org/x/net/html"
)

type musicGroup struct {
	Name string `json:"name"`
}

type musicRecord struct {
	Name     string     `json:"name"`
	Duration string     `json:"duration"`
	Url      string     `json:"url"`
	ByArtist musicGroup `json:"byArtist"`
	Id       string     `json:"@id"`
}

type SoundCloudExtractor struct{}

func (*SoundCloudExtractor) Match(url string) (string, string) {
	re := regexp.MustCompile(`https?://(?:www\.)?soundcloud.com/([a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+)(?:\?.+)?`)
	groups := re.FindStringSubmatch(url)
	if groups == nil {
		return "", ""
	}

	return groups[0], groups[1]
}

func (e *SoundCloudExtractor) Extract(recordId string) (types.MusicRecord, error) {
	soundCloudRecord := e.getSoundcloudData("https://m.soundcloud.com/" + recordId)
	recordId, err := e.getRecordId(soundCloudRecord.Id)
	if err != nil {
		return types.MusicRecord{}, err
	}

	return types.MusicRecord{
		Band:     types.Band{Name: soundCloudRecord.ByArtist.Name},
		Duration: iso8601.ParseDuration(soundCloudRecord.Duration),
		Name:     soundCloudRecord.Name,
		RecordId: recordId,
		Url:      soundCloudRecord.Url,
	}, nil
}

func (*SoundCloudExtractor) getRecordId(recordId string) (string, error) {
	parts := strings.Split(recordId, ":")
	if len(parts) != 3 {
		return "", fmt.Errorf("unknown Soundcloud record id '%s'", recordId)
	}

	return parts[2], nil
}

func (e *SoundCloudExtractor) getSoundcloudData(url string) musicRecord {
	var resp *http.Response
	var err error

	for retries := 0; retries < 10; retries++ {
		resp, err = http.Get(url)

		if err != nil {
			log.Fatal("Error while Get", err)
		} else {
			if resp.StatusCode == http.StatusOK {
				break
			} else {
				log.Print("Received an HTTP error: ", resp.Status)
				time.Sleep(2 * time.Second)
			}
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Received an HTTP error: ", resp.Status)
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error while Parse: ", err)
	}

	ldjson := e.parse(document)
	if ldjson == "" {
		log.Fatal("No ld+json found")
	}

	var music musicRecord
	err = json.Unmarshal([]byte(ldjson), &music)
	if err != nil {
		log.Fatal("Error while Unmarshal: ", err)
	}

	return music
}

// Return the content of the first <script type="application/ld+json"> node found.
func (e *SoundCloudExtractor) parse(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "script" {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "application/ld+json" {
				return node.FirstChild.Data
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if content := e.parse(child); content != "" {
			return content
		}
	}

	return ""
}

func init() {
	extractors = append(extractors, &SoundCloudExtractor{})
}
