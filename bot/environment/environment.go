package environment

import (
	"github.com/spf13/viper"
	"log"
)

type Environment interface {
	ReadIn(file string) error
}

type Config interface {
	Load(map[string]any) error
}

func ReadConfigFromFileAndSystem(file string) map[string]any {
	vConfig := viper.New()
	vConfig.SetConfigFile(file)
	err := vConfig.ReadInConfig()
	if err != nil {
		log.Println("Failed to load", file, ", using only system environment variables instead")
	}
	vConfig.AutomaticEnv()
	return vConfig.AllSettings()
}
