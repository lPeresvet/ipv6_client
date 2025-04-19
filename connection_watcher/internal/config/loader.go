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

func (l *Loader) Load(configName string) (*config.WatcherConfig, error) {
	loadedConfig, err := l.configRepository.GetConfig(configName)
	if err != nil {
		return &config.WatcherConfig{}, err
	}

	if loadedConfig.Watcher.Reconnect.WaitingTimeout == 0 {
		loadedConfig.Watcher.Reconnect.WaitingTimeout = 5
	}

	if loadedConfig.Watcher.Reconnect.WatchingPeriod == 0 {
		loadedConfig.Watcher.Reconnect.WatchingPeriod = 5
	}

	return &loadedConfig.Watcher, nil
}
