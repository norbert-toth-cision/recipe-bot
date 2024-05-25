package environment

const (
	user        = "RABBITMQ_USER"
	password    = "RABBITMQ_PASSWORD"
	host        = "RABBITMQ_HOST"
	port        = "RABBITMQ_PORT"
	virtHost    = "RABBITMQ_VIRT_HOST"
	outputQueue = "RABBITMQ_OUTPUT_QUEUE_NAME"
)

type RmqConfig struct {
	User        string
	Password    string
	Host        string
	Port        int
	VirtualHost string
	OutputQueue string
}

func (rmqConf *RmqConfig) Load(env map[string]any) error {
	var err error
	if rmqConf.User, err = GetRequiredString(env, user); err != nil {
		return err
	}
	if rmqConf.Password, err = GetRequiredString(env, password); err != nil {
		return err
	}
	if rmqConf.Host, err = GetRequiredString(env, host); err != nil {
		return err
	}
	if rmqConf.Port, err = GetRequiredInt(env, port); err != nil {
		return err
	}
	if rmqConf.VirtualHost, err = GetRequiredString(env, virtHost); err != nil {
		return err
	}
	if rmqConf.OutputQueue, err = GetRequiredString(env, outputQueue); err != nil {
		return err
	}
	return nil
}
