package monitoring

import (
	"github.com/spf13/cast"
	"net"
	"recipebot/environment"
)

type SimpleMonitor struct {
	config *environment.SimpleActuatorConfig
}

func (m *SimpleMonitor) Configure(config *environment.SimpleActuatorConfig) {
	m.config = config
}

func (m *SimpleMonitor) Monitor() error {
	_, err := net.Listen("tcp", ":"+cast.ToString(m.config.MonitoringPort))
	return err
}
