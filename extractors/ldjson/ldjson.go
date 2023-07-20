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

type MusicGroup struct {
	Name string `json:"name"`
}

type MusicRecord struct {
	Name     string     `json:"name"`
	Duration string     `json:"duration"`
	Url      string     `json:"url"`
	ByArtist MusicGroup `json:"byArtist"`
	Id       string     `json:"@id"`
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

	return types.MusicRecord{
		Band:     types.Band{Name: record.ByArtist.Name},
		Duration: iso8601.ParseDuration(record.Duration),
		Name:     record.Name,
		RecordId: record.Id,
		Url:      record.Url,
	}, nil
}

func (e ldJsonExtractor) getMusicRecord(url string) MusicRecord {
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

	var music MusicRecord
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
