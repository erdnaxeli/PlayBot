package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

// Return the content of the first <script type="application/ld+json"> node found.
func parse(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "script" {
		for _, attr := range node.Attr {
			if attr.Key == "type" && attr.Val == "application/ld+json" {
				return node.FirstChild.Data
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if content := parse(child); content != "" {
			return content
		}
	}

	return ""
}

func getSoundcloudData(url string) musicRecord {
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

	ldjson := parse(document)
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

func parseDuration(duration string) int {
	re := regexp.MustCompile(`PT(?:(?P<hours>\d\d?)H)?(?:(?P<minutes>\d\d?)M)(?:(?P<secondes>\d\d?)S)`)

	groups := re.FindStringSubmatch(duration)

	var hours, minutes, seconds int
	if groups[1] != "" {
		hours, _ = strconv.Atoi(groups[1])
	} else {
		hours = 0
	}

	if groups[2] != "" {
		minutes, _ = strconv.Atoi(groups[2])
	} else {
		minutes = 0
	}

	if groups[3] != "" {
		seconds, _ = strconv.Atoi(groups[3])
	} else {
		seconds = 0
	}

	return hours*3600 + minutes*60 + seconds
}

func main() {
	soundcloudId := os.Args[1]

	music := getSoundcloudData("https://m.soundcloud.com/" + soundcloudId)

	id := strings.Split(music.Id, ":")

	fmt.Println(id[2])
	fmt.Println(music.Url)
	fmt.Println(music.Name)
	fmt.Println(music.ByArtist.Name)
	fmt.Println(parseDuration(music.Duration))
}
