package config

import (
	"implementation/client_src/pkg/config"
)

type ConfigRepository interface {
	GetConfig(name string) (*config.Config, error)
}

type Loader struct {
	configRepository ConfigRepository
}

func NewLoader(repository ConfigRepository) *Loader {
	return &Loader{
		configRepository: repository,
	}
}

func (l *Loader) Load(configName string) (*config.Config, error) {
	loadedConfig, err := l.configRepository.GetConfig(configName)
	if err != nil {
		return &config.Config{}, err
	}

	return loadedConfig, nil
}
