package config

import (
	"github.com/spf13/viper"
	"log"
)

const File = ".env"

type Key string

const (
	BotToken = Key("BOT_TOKEN")
	Other    = Key("OTHER")
)

func init() {
	log.Println("Config read")
	viper.SetConfigFile(File)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not read config file: ", File)
	}
}

func GetConfig(key Key) (token string) {
	token = viper.GetString(string(key))
	if token == "" {
		log.Fatal("Environment variable ", key, " has no value")
	}
	return
}
