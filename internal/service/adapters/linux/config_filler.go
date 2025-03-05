package linux

import (
	"errors"
	"fmt"
	"github.com/katalix/go-l2tp/config"
)

const (
	baseTemplatesPath = "config/templates"
)

type FileFiller func(config *config.Config) error

type ConfigFiller struct {
	templatesPath string
}

func NewConfigFiller(templatesPath string) *ConfigFiller {
	if templatesPath == "" {
		templatesPath = baseTemplatesPath
	}

	return &ConfigFiller{
		templatesPath: templatesPath,
	}
}

func (filler *ConfigFiller) FillConfig(userConfig *config.Config) error {
	if userConfig == nil {
		return errors.New("userConfig is nil")
	}

	fileFillers := map[string]FileFiller{
		"xl2tpd.conf":  FillXL2TP_conf,
		"options.l2tp": FillOPTIONS_l2tp,
		"chap-secrets": FillCHAP_SECRETS,
	}

	for fileName, fileFiller := range fileFillers {
		if err := fileFiller(userConfig); err != nil {
			return fmt.Errorf("failed to fill %s: %w", fileName, err)
		}
	}

	return nil
}

func FillXL2TP_conf(config *config.Config) error {
	return nil
}

func FillOPTIONS_l2tp(config *config.Config) error {
	return nil
}

func FillCHAP_SECRETS(config *config.Config) error {
	return nil
}
