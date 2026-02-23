// Package ldjson a type to extract JSON-LD data from a webpage.
package ldjson

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/erdnaxeli/PlayBot/iso8601"
	"github.com/erdnaxeli/PlayBot/types"
	"golang.org/x/net/html"
)

// MusicAlbum contains data about the album where the music was publied.
type MusicAlbum struct {
	ByArtist MusicGroup `json:"byArtist"`
}

// MusicGroup contains data about the author of the music.
type MusicGroup struct {
	Name string `json:"name"`
}

// MusicRecording contains data about a music.
type MusicRecording struct {
	ByArtist         MusicGroup `json:"byArtist"`
	Duration         string     `json:"duration"`
	ID               string     `json:"@id"`
	InAlbum          MusicAlbum `json:"InAlbum"`
	MainEntityOfPage string     `json:"mainEntityOfPage"`
	Name             string     `json:"name"`
	URL              string     `json:"url"`
}

// Extractor is able to extracts music data from an URL exposing JSON-LD data.
type Extractor interface {
	Extract(string) (types.MusicRecord, error)
}

type extractor struct{}

// New returns a new instance of an Extractor object.
func New() Extractor {
	return extractor{}
}

func (e extractor) Extract(url string) (types.MusicRecord, error) {
	record, err := e.getMusicRecord(url)
	if err != nil {
		return types.MusicRecord{}, err
	}
	recordURL := record.URL
	if recordURL == "" {
		recordURL = record.MainEntityOfPage
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
		RecordID: record.ID,
		URL:      recordURL,
	}, nil
}

func (e extractor) getMusicRecord(url string) (MusicRecording, error) {
	var resp *http.Response
	var err error

	for range 10 {
		resp, err = http.Get(url)

		if err != nil {
			return MusicRecording{}, fmt.Errorf("error while Get: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		log.Print("Received an HTTP error: ", resp.Status)
		time.Sleep(2 * time.Second)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return MusicRecording{}, fmt.Errorf("%w: %d %s", ErrHTTPNotOk, resp.StatusCode, resp.Status)
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		return MusicRecording{}, fmt.Errorf("error while Parse: %w", err)
	}

	ldjson := e.parse(document)
	if ldjson == "" {
		return MusicRecording{}, ErrNoLDJSON
	}

	var music MusicRecording
	err = json.Unmarshal([]byte(ldjson), &music)
	if err != nil {
		return MusicRecording{}, fmt.Errorf("error while Unmarshal: %w", err)
	}

	return music, nil
}

// Return the content of the first <script type="application/ld+json"> node found.
func (e extractor) parse(node *html.Node) string {
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
