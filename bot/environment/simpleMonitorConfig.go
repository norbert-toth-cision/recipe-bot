package environment

const (
	monitoringPort = "MONITORING_PORT"
)

type SimpleActuatorConfig struct {
	MonitoringPort int
}

func (mConf *SimpleActuatorConfig) Load(env map[string]any) error {
	var err error
	mConf.MonitoringPort, err = GetRequiredInt(env, monitoringPort)
	return err
}
