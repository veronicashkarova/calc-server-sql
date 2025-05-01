package application

import (
	"github.com/veronicashkarova/agent/pkg/agent"
	"os"
	"strconv"
)

type Config struct {
	COMPUTING_POWER int
	IDLE_DELAY int
}

func ConfigFromEnv() *Config {

	config := new(Config)
	power, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		config.COMPUTING_POWER = 3
	} else {
		config.COMPUTING_POWER = power
	}

	idleDelay, err := strconv.Atoi(os.Getenv("IDLE_DELAY"))
	if err != nil {
		config.IDLE_DELAY = 5000
	} else {
		config.IDLE_DELAY = idleDelay
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	config := ConfigFromEnv()
	return &Application{
		config: config,
	}
}

func (a *Application) RunAgent() {
	agent.RunAgent(a.config.COMPUTING_POWER, a.config.IDLE_DELAY)
}
