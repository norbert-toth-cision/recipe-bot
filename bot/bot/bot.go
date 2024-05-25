package bot

import "recipebot/queue"

type Bot interface {
	WithQueue(queue queue.Queue)
	Start() error
	Stop() error
}
