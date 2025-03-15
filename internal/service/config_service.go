package service

import (
	"fmt"
	"implementation/internal/domain/config"
)

const xl2tpdDemonName = "xl2tpd.service"

type ConfigProvider interface {
	GetConfig(name string) (*config.Config, error)
}

type ConfigFiller interface {
	FillConfig(userConfig *config.Config) error
}

type ConfigService struct {
	configProvider ConfigProvider
	configFiller   ConfigFiller
	demonProvider  DemonProvider
}

func NewConfigService(provider ConfigProvider, filler ConfigFiller, demon DemonProvider) *ConfigService {
	return &ConfigService{
		configProvider: provider,
		configFiller:   filler,
		demonProvider:  demon,
	}
}

// GetConfig return parsed configuration config.Config
func (s *ConfigService) GetConfig(path string) (*config.Config, error) {
	return s.configProvider.GetConfig(path)
}

// InitConfig fill config files
func (s *ConfigService) InitConfig(path string) error {
	if err := s.demonProvider.StopDemon(xl2tpdDemonName); err != nil {
		return err
	}

	userConfig, err := s.configProvider.GetConfig(path)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := s.configFiller.FillConfig(userConfig); err != nil {
		return err
	}

	return s.demonProvider.StopDemon(xl2tpdDemonName)
}
