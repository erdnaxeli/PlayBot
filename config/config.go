package config

import (
	"io/fs"

	"github.com/spf13/viper"
)

type IrcConfig struct {
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Nick             string   `json:"nick"`
	Channels         []string `json:"channels"`
	NickServPassword string   `json:"nickserv_password"`
}

type Config struct {
	YoutubeApiKey string
	DbName        string
	DbUser        string
	DbHost        string
	DbPassword    string
	Irc           IrcConfig
	Timezone      string
	ServerAddress string
}

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
		YoutubeApiKey: viper.GetString("youtube_api_key"),
		DbName:        viper.GetString("bdd"),
		DbHost:        viper.GetString("host"),
		DbUser:        viper.GetString("user"),
		DbPassword:    viper.GetString("passwd"),
		Irc: IrcConfig{
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
