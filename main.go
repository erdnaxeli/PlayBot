package main

import (
	"fmt"
	"log"
	"os"

	"github.com/erdnaxeli/PlayBot/extractors"
)

func main() {
	url := os.Args[1]

	music, err := extractors.Extract(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(music.RecordId)
	fmt.Println(music.Url)
	fmt.Println(music.Name)
	fmt.Println(music.Band.Name)
	fmt.Println(music.Duration.Seconds())
}
