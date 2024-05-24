package queue

import (
	"log"
	"recipebot/config"
)

func PushIntoQueue(message interface{}) {

	log.Println("pushing message ", message, " to ", config.GetConfig(config.Other))
}
