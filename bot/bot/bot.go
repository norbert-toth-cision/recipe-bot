package bot

import "recipebot/config"

type Bot interface {
	Configure(config.Config) error
	Start() error
	Stop() error
}
