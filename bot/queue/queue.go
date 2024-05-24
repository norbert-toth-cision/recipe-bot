package queue

import "recipebot/config"

type Queue interface {
	Configure(config.Config)
	SendMessage(interface{})
	Close() error
}
