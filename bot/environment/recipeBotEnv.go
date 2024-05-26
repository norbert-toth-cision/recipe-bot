package environment

import (
	"errors"
	"github.com/spf13/cast"
	"strings"
)

type VarOrFileEnvironment struct {
	MonitoringConfig *SimpleActuatorConfig
	BotConfig        *RecipeBotConfig
	TiktokProcConfig *TiktokProcConfig
	DropboxConfig    *DropboxConfig
}

func (e *VarOrFileEnvironment) ReadIn(file string) error {
	env := ReadConfigFromFileAndSystem(file)

	e.MonitoringConfig = new(SimpleActuatorConfig)
	if err := e.MonitoringConfig.Load(env); err != nil {
		return err
	}

	e.BotConfig = new(RecipeBotConfig)
	if err := e.BotConfig.Load(env); err != nil {
		return err
	}

	e.TiktokProcConfig = new(TiktokProcConfig)
	if err := e.TiktokProcConfig.Load(env); err != nil {
		return err
	}
	e.DropboxConfig = new(DropboxConfig)
	if err := e.DropboxConfig.Load(env); err != nil {
		return err
	}
	return nil
}

func GetRequiredString(env map[string]any, key string) (string, error) {
	v, ok := env[strings.ToLower(key)]
	if !ok {
		return "", errors.New("Required setting " + key + " is missing")
	}
	return cast.ToString(v), nil
}

func GetRequiredInt(env map[string]any, key string) (int, error) {
	v, ok := env[strings.ToLower(key)]
	if !ok {
		return 0, errors.New("Required setting " + key + " is missing")
	}
	return cast.ToInt(v), nil
}
