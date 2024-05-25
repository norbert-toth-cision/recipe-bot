package environment

const (
	botToken = "BOT_TOKEN"
)

type RecipeBotConfig struct {
	BotToken string
}

func (rConf *RecipeBotConfig) Load(env map[string]any) error {
	var err error
	rConf.BotToken, err = GetRequiredString(env, botToken)
	return err
}
