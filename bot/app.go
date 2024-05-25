package main

import (
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"os/signal"
	"recipebot/bot"
)

const (
	ConfigFile = ".env"
)

func main() {
	vConfig := viper.New()
	vConfig.SetConfigFile(ConfigFile)
	err := vConfig.ReadInConfig()
	onErrorFatal(err)
	vConfig.AutomaticEnv()

	var recipeBot bot.Bot
	recipeBot = new(bot.RecipeBot)
	err = recipeBot.Configure(vConfig)
	onErrorFatal(err)
	err = recipeBot.Start()
	onErrorFatal(err)

	defer func() {
		err := recipeBot.Stop()
		onErrorFatal(err)
	}()

	log.Println("Bot started")
	listenMonitoring()
	listenInterrupt()
}

func listenMonitoring() {
	_, err := net.Listen("tcp", ":8300")
	onErrorFatal(err)
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
