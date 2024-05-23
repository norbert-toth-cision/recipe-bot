package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"recipebot/bot"
)

const (
	TokenVarName = "BOT_TOKEN"
	ConfigFile   = ".env"
)

func main() {
	token := obtainTokenOrFail()
	discord, err := discordgo.New("Bot " + token)
	onErrorFatal(err)

	discord.AddHandler(bot.OnNewMessage)
	err = discord.Open()
	onErrorFatal(err)

	defer func(discord *discordgo.Session) {
		err := discord.Close()
		onErrorFatal(err)
	}(discord)

	log.Println("Bot started")
	listenInterrupt()
}

func onErrorFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func obtainTokenOrFail() string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not read config file: ", ConfigFile)
	}

	token := viper.GetString(TokenVarName)
	if token == "" {
		log.Fatal("Environment variable ", TokenVarName, " has no value")
	}
	return token
}

func listenInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
