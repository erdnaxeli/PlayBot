package ldjson

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/erdnaxeli/PlayBot/iso8601"
	"github.com/erdnaxeli/PlayBot/types"
	"golang.org/x/net/html"
)

type MusicAlbum struct {
	ByArtist MusicGroup `json:"byArtist"`
}

type MusicGroup struct {
	Name string `json:"name"`
}

type MusicRecording struct {
	ByArtist         MusicGroup `json:"byArtist"`
	Duration         string     `json:"duration"`
	Id               string     `json:"@id"`
	InAlbum          MusicAlbum `json:"InAlbum"`
	MainEntityOfPage string     `json:"mainEntityOfPage"`
	Name             string     `json:"name"`
	Url              string     `json:"url"`
}

type LdJsonExtractor interface {
	Extract(string) (types.MusicRecord, error)
}

type ldJsonExtractor struct{}

func NewLdJsonExtractor() LdJsonExtractor {
	return ldJsonExtractor{}
}

func (e ldJsonExtractor) Extract(url string) (types.MusicRecord, error) {
	record := e.getMusicRecord(url)
	recordUrl := record.Url
	if recordUrl == "" {
		recordUrl = record.MainEntityOfPage
	}

	band := record.InAlbum.ByArtist.Name
	if band == "" {
		band = record.ByArtist.Name
	}

	duration, err := iso8601.ParseDuration(record.Duration)
	if err != nil {
		return types.MusicRecord{}, err
	}

	return types.MusicRecord{
		Band:     types.Band{Name: band},
		Duration: duration,
		Name:     record.Name,
		RecordID: record.Id,
		URL:      recordUrl,
	}, nil
}

func (e ldJsonExtractor) getMusicRecord(url string) MusicRecording {
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
	defer func() { _ = resp.Body.Close() }()

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

	var music MusicRecording
	err = json.Unmarshal([]byte(ldjson), &music)
	if err != nil {
		log.Fatal("Error while Unmarshal: ", err)
	}

	return music
}

// Return the content of the first <script type="application/ld+json"> node found.
func (e ldJsonExtractor) parse(node *html.Node) string {
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
