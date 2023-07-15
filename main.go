package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erdnaxeli/PlayBot/extractors"
)

func main() {
	msg := os.Args[1]

	matchedUrl, music, err := extractors.Extract(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(matchedUrl)
	fmt.Println(music.RecordId)
	fmt.Println(music.Url)
	fmt.Println(music.Name)
	fmt.Println(music.Band.Name)
	fmt.Println(music.Duration.Seconds())
}
