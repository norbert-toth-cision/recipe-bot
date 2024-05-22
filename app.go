package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"my-first/hello/bot"
	"os"
	"os/signal"
)

const TokenVarName = "BOT_TOKEN"

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
	var token = os.Getenv(TokenVarName)
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
