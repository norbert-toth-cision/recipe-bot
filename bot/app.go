package main

import (
	"github.com/spf13/viper"
	"log"
	"net"
	"os"
	"os/signal"
	"recipebot/bot"
	"recipebot/config"
)

const (
	ConfigFile = ".env"
)

func main() {
	vConfig := viper.New()
	vConfig.SetConfigFile(ConfigFile)
	err := vConfig.ReadInConfig()
	if err != nil {
		log.Println("Failed to load", ConfigFile, ", using only system environment variables instead")
	}
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
	listenMonitoring(vConfig)
	listenInterrupt()
}

func listenMonitoring(configs config.Config) {
	_, err := net.Listen("tcp", ":"+configs.GetString(config.MONITORING_PORT))
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
