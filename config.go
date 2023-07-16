package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	YoutubeApiKey string `json:"youtube_api_key"`
}

func readConfigFile(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, nil
	}

	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, nil
	}

	var config Config
	json.Unmarshal(content, &config)

	return config, nil
}
