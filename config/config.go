package config

import (
	"encoding/json"
	"io"
	"os"
)

type IrcConfig struct {
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Nick             string   `json:"nick"`
	Channels         []string `json:"channels"`
	NickServPassword string   `json:"nickserv_password"`
}

type Config struct {
	YoutubeApiKey string    `json:"youtube_api_key"`
	DbName        string    `json:"bdd"`
	DbUser        string    `json:"user"`
	DbHost        string    `json:"host"`
	DbPassword    string    `json:"passwd"`
	Irc           IrcConfig `json:"irc"`
}

func ReadConfigFile(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, nil
	}

	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return Config{}, nil
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
