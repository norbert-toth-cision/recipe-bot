package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"recipebot/bot"
	"recipebot/config"
)

func main() {
	discord, err := discordgo.New("Bot " + config.GetConfig(config.BotToken))
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

func listenInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
