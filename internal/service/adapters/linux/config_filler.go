package linux

import (
	"bufio"
	"errors"
	"fmt"
	"implementation/internal/domain/config"
	"log"
	"os"
	"path"
	"strings"
)

const (
	baseTemplatesPath = "config/templates"
	xl2tpPath         = "/etc/xl2tpd/xl2tpd.conf"
	optionsPath       = "/etc/ppp/options.l2tp"
	chapSecretsPath   = "/etc/ppp/chap-secrets"
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
		"xl2tpd.conf":  filler.fillXL2TP_conf,
		"options.l2tp": filler.fillOPTIONS_l2tp,
		"chap-secrets": filler.fillCHAP_SECRETS,
	}

	for fileName, fileFiller := range fileFillers {
		if err := fileFiller(userConfig); err != nil {
			return fmt.Errorf("failed to fill %s: %w", fileName, err)
		}
	}

	return nil
}

func (filler *ConfigFiller) fillXL2TP_conf(userConfig *config.Config) error {
	var (
		file *os.File
		err  error
	)

	if !isFileExists(xl2tpPath) {
		log.Printf("xl2tp config file %s does not exist", xl2tpPath)

		if err := copyFile(path.Join(filler.templatesPath, "xl2tp_global_template.ini"), xl2tpPath); err != nil {
			return fmt.Errorf("failed to copy xl2tp_global_template.ini: %w", err)
		}
	}
	log.Printf("filling file %s", xl2tpPath)

	file, err = os.OpenFile(xl2tpPath, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open xl2tp.config file: %v", err)
	}

	defer file.Close()

	for _, server := range userConfig.Servers {
		for _, user := range server.Users {
			log.Printf("filling user [%v]", user.Username)

			contains, err := isContainsInFile(file, fmt.Sprintf("[lac %s]", user.Username))
			if err != nil {
				return fmt.Errorf("failed to check if contains %s: %w", user.Username, err)
			}

			if !contains {
				if err := filler.insertLNSSection(file, user.Username, server.Address); err != nil {
					return fmt.Errorf("failed to insert lns section: %w", err)
				}
			}
		}
	}

	return nil
}

func (filler *ConfigFiller) insertLNSSection(file *os.File, userName, serverAddress string) error {
	templateFile, err := os.Open(path.Join(filler.templatesPath, "xl2tp_lns_template.ini"))
	if err != nil {
		return fmt.Errorf("failed to open xl2tp_lns_template.ini: %w", err)
	}

	scanner := bufio.NewScanner(templateFile)
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.Replace(line, "%username%", userName, 1)
		line = strings.Replace(line, "%server_ip%", serverAddress, 1) + "\n"

		if err := appendToFile(file, line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return templateFile.Close()
}

func (filler *ConfigFiller) fillOPTIONS_l2tp(config *config.Config) error {
	return nil
}

func (filler *ConfigFiller) fillCHAP_SECRETS(config *config.Config) error {
	return nil
}

func isFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func appendToFile(file *os.File, data string) error {
	_, err := file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func isContainsInFile(file *os.File, substring string) (bool, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, substring) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func copyFile(src, destination string) (err error) {
	data, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("failed to read src: %v", err)
	}

	if err := os.WriteFile(destination, data, 0644); err != nil {
		log.Fatalf("failed to write to target: %v", err)
	}

	return nil
}
