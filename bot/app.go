package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"recipebot/bot"
	"recipebot/environment"
	"recipebot/monitoring"
	"recipebot/queue"
	"syscall"
)

const (
	ConfigFile = ".env"
)

func main() {
	appConfig := new(environment.VarOrFileEnvironment)
	if err := appConfig.ReadIn(ConfigFile); err != nil {
		onErrorFatal(err)
	}

	recipeBot := createBot(appConfig.BotConfig)
	if queue, err := createQueue(appConfig.RabbitMQConfig); err != nil {
		onErrorFatal(err)
	} else {
		recipeBot.WithQueue(queue)
	}

	if err := recipeBot.Start(); err != nil {
		graceFullyStopBotAndExit(recipeBot, err)
	}

	defer graceFullyStopBot(recipeBot)

	log.Println("Bot started")

	if err := createMonitor(appConfig.MonitoringConfig).Monitor(); err != nil {
		graceFullyStopBotAndExit(recipeBot, err)
	}
	listenInterrupt()
}

func createBot(botConfig *environment.RecipeBotConfig) bot.Bot {
	recipeBot := new(bot.RecipeBot)
	recipeBot.Configure(botConfig)
	return recipeBot
}

func createQueue(config *environment.RmqConfig) (queue.Queue, error) {
	rmq := new(queue.RMQueue)
	if err := rmq.Configure(config); err != nil {
		closeErr := rmq.Close()
		return nil, errors.Join(err, closeErr)
	}
	return rmq, nil
}

func createMonitor(config *environment.SimpleActuatorConfig) monitoring.Monitor {
	monitor := new(monitoring.SimpleMonitor)
	monitor.Configure(config)
	return monitor
}

func graceFullyStopBot(bot bot.Bot) {
	err := bot.Stop()
	onErrorFatal(err)
}

func graceFullyStopBotAndExit(bot bot.Bot, original error) {
	err := bot.Stop()
	onErrorFatal(errors.Join(original, err))
}

func onErrorFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func listenInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-c
}
