// Package config provides functions to read the app configuration.
package config

import (
	"io/fs"

	"github.com/spf13/viper"
)

// IRCConfig represent the specific configuration for the IRC client of the app.
type IRCConfig struct {
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Nick             string   `json:"nick"`
	Channels         []string `json:"channels"`
	NickServPassword string   `json:"nickserv_password"`
}

// Config represents the configuration to run the playbot app.
type Config struct {
	YoutubeAPIKey string
	DbName        string
	DbUser        string
	DbHost        string
	DbPassword    string
	IRC           IRCConfig
	Timezone      string
	ServerAddress string
}

// ReadConfig read the configuratio from a config file or env variables.
func ReadConfig() (Config, error) {
	viper.SetDefault("server.address", "localhost:1111")

	viper.SetConfigFile("playbot.conf")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		// we ignore the error if it is about the file being not found
		if _, ok := err.(*fs.PathError); !ok {
			return Config{}, err
		}
	}

	_ = viper.BindEnv("youtube_api_key", "YOUTUBE_API_KEY")
	_ = viper.BindEnv("host", "DB_HOST")
	_ = viper.BindEnv("bdd", "DB_NAME")
	_ = viper.BindEnv("user", "DB_USER")
	_ = viper.BindEnv("passwd", "DB_PWD")
	_ = viper.BindEnv("irc.host", "IRC_HOST")
	_ = viper.BindEnv("irc.port", "IRC_PORT")
	_ = viper.BindEnv("irc.nick", "IRC_NICK")
	_ = viper.BindEnv("irc.channels", "IRC_CHANNELS")
	_ = viper.BindEnv("irc.nickserv_password", "IRC_NICKSERV_PWD")
	_ = viper.BindEnv("server.address", "SERVER_ADDR")
	_ = viper.BindEnv("timezone", "TIMEZONE")

	return Config{
		YoutubeAPIKey: viper.GetString("youtube_api_key"),
		DbName:        viper.GetString("bdd"),
		DbHost:        viper.GetString("host"),
		DbUser:        viper.GetString("user"),
		DbPassword:    viper.GetString("passwd"),
		IRC: IRCConfig{
			Host:             viper.GetString("irc.host"),
			Port:             viper.GetInt("irc.port"),
			Nick:             viper.GetString("irc.nick"),
			Channels:         viper.GetStringSlice("irc.channels"),
			NickServPassword: viper.GetString("irc.nickserv_password"),
		},
		ServerAddress: viper.GetString("server.address"),
		Timezone:      viper.GetString("timezone"),
	}, nil
}
