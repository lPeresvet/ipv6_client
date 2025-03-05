package service

import (
	"fmt"
	"implementation/internal/domain/config"
)

type ConfigProvider interface {
	ReadConfig() (*config.Config, error)
}

type ConfigFiller interface {
	FillConfig(userConfig *config.Config) error
}

type ConfigService struct {
	configProvider ConfigProvider
	configFiller   ConfigFiller
}

func NewConfigService(provider ConfigProvider, filler ConfigFiller) ConfigService {
	return ConfigService{
		configProvider: provider,
		configFiller:   filler,
	}
}

// GetConfig return parsed configuration config.Config
func (s *ConfigService) GetConfig() (*config.Config, error) {
	return s.configProvider.ReadConfig()
}

// InitConfig fill config files
func (s *ConfigService) InitConfig() error {
	userConfig, err := s.configProvider.ReadConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	return s.configFiller.FillConfig(userConfig)
}
