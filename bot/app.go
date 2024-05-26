package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"recipebot/bot"
	"recipebot/cloud"
	"recipebot/environment"
	"recipebot/monitoring"
	"recipebot/urlProcessor"
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

	recipeBot, err := createBot(appConfig)
	if err != nil {
		onErrorFatal(err)
	}

	if err := recipeBot.Start(); err != nil {
		graceFullyStopBotAndFatal(recipeBot, err)
	}

	defer graceFullyStopBot(recipeBot)

	log.Println("Bot started")

	if err := createMonitor(appConfig.MonitoringConfig).Monitor(); err != nil {
		graceFullyStopBotAndFatal(recipeBot, err)
	}
	listenInterrupt()
}

func createBot(env *environment.VarOrFileEnvironment) (bot.Bot, error) {
	recipeBot := new(bot.RecipeBot)
	recipeBot.Configure(env.BotConfig)

	dropbox := cloud.NewDropbox(env.DropboxConfig)
	proc, err := urlProcessor.NewTiktokProc(env.TiktokProcConfig, dropbox)
	if err != nil {
		return nil, err
	}
	recipeBot.AddProcessor(proc)
	return recipeBot, nil

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

func graceFullyStopBotAndFatal(bot bot.Bot, original error) {
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
