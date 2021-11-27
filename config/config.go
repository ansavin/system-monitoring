package config

import (
	"fmt"
	"os"

	"go.uber.org/config"
)

// Cfg represents our service config
type Cfg struct {
	DisableCPUStats bool `yaml:"DisableCPUStats"`
	DisableDevStats bool `yaml:"DisableDevStats"`
	DisableFsStats  bool `yaml:"DisableFsStats"`
}

const configPath = "./config.yml"

// NewConfig get data from config
func NewConfig() (Cfg, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return Cfg{}, fmt.Errorf("cannot get config: %s", err.Error())
	}

	provider, err := config.NewYAML(config.Source(file))
	if err != nil {
		panic(err)
		return Cfg{}, fmt.Errorf("cannot get config: %s", err.Error())
	}

	var c Cfg
	if err := provider.Get("").Populate(&c); err != nil {
		return Cfg{}, fmt.Errorf("cannot populate config: %s", err.Error())
	}
	return c, nil
}
