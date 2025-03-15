package linux

import (
	"bufio"
	"errors"
	"fmt"
	"implementation/internal/domain/config"
	"implementation/internal/domain/template"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

const (
	baseTemplatesPath = "config/templates"
	xl2tpPath         = "/etc/xl2tpd/xl2tpd.conf"
	optionsPath       = "/etc/ppp/options.l2tp"
	chapSecretsPath   = "/etc/ppp/chap-secrets"
)

type FileFiller func(config *config.Config) error

type ParsedSecrets map[string]map[string]string

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
		// todo could use different goroutines to fill in parallel
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
			contains, err := isContainsInFile(file, fmt.Sprintf("[lac %s]", user.Username))
			if err != nil {
				return fmt.Errorf("failed to check if contains %s: %w", user.Username, err)
			}

			if !contains {
				log.Printf("insert user [%v] data", user.Username)
				if err := filler.proceedTemplate(file, user, server.Address, "xl2tp_lac_template.ini"); err != nil {
					return fmt.Errorf("failed to insert lac section: %w", err)
				}
			}
		}
	}

	return nil
}

func (filler *ConfigFiller) proceedTemplate(file *os.File, userConfig config.UserConfig, serverAddress, templateName string) error {
	templateFile, err := os.Open(path.Join(filler.templatesPath, templateName))
	if err != nil {
		return fmt.Errorf("failed to open xl2tp_lac_template.ini: %w", err)
	}

	log.Printf("fill template file %s", templateName)

	scanner := bufio.NewScanner(templateFile)
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.Replace(line, template.UsernamePlaceholder, userConfig.Username, -1)
		line = strings.Replace(line, template.UsernamePlaceholder, userConfig.Password, -1)
		line = strings.Replace(line, template.ServerIpPlaceholder, serverAddress, -1) + "\n"

		if err := appendToFile(file, line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return templateFile.Close()
}

func (filler *ConfigFiller) fillOPTIONS_l2tp(userConfig *config.Config) error {
	for _, server := range userConfig.Servers {
		for _, user := range server.Users {
			if err := filler.fillOneOptionFile(user, server); err != nil {
				log.Printf("Failed to insert oprion filefor user <%s> %v:", user.Username, err)
				continue
			}
		}
	}
	return nil
}

func (filler *ConfigFiller) fillOneOptionFile(user config.UserConfig, server config.ServerConfig) error {
	filePath := buildUserOptionsPath(user.Username)

	if isFileExists(filePath) {
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to delete file <%s> %w", filePath, err)
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s file: %w", filePath, err)
	}

	defer file.Close()

	if err := filler.proceedTemplate(file, user, server.Address, "options_template.ini"); err != nil {
		return fmt.Errorf("failed to fill options template: %w", err)
	}

	return nil
}

func buildUserOptionsPath(user string) string {
	return fmt.Sprintf("%s.%s", optionsPath, user)
}

func (filler *ConfigFiller) fillCHAP_SECRETS(userConfig *config.Config) error {
	var err error
	data := make(ParsedSecrets)

	if isFileExists(chapSecretsPath) {
		data, err = filler.parseSecrets(chapSecretsPath)
		if err != nil {
			return fmt.Errorf("failed to parse chap secrets: %w", err)
		}

		if err := filler.deleteFile(chapSecretsPath); err != nil {
			return fmt.Errorf("failed to delete file <%s>: %w", chapSecretsPath, err)
		}
	}

	for _, server := range userConfig.Servers {
		for _, user := range server.Users {
			if _, ok := data[user.Username]; !ok {
				data[user.Username] = make(map[string]string)
			}
			data[user.Username]["*"] = user.Password
		}
	}

	return filler.writeConfigToSecrets(data, chapSecretsPath)
}

func (filler *ConfigFiller) writeConfigToSecrets(secrets ParsedSecrets, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %s file: %w", filePath, err)
	}
	defer file.Close()

	for username, passwords := range secrets {
		for server, password := range passwords {
			line := fmt.Sprintf("\"%s\" \"%s\" \"%s\" \"*\"\n", username, server, password)

			if _, err := file.WriteString(line); err != nil {
				return fmt.Errorf("failed to write to file <%s>: %w", filePath, err)
			}
		}
	}

	return nil
}

func (filler *ConfigFiller) deleteFile(secretsPath string) error {
	return os.Remove(secretsPath)
}

func (filler *ConfigFiller) parseSecrets(secretsPath string) (ParsedSecrets, error) {
	parsed := make(ParsedSecrets)

	file, err := os.OpenFile(secretsPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s file: %w", secretsPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		result := extractFromLine(line)

		if len(result) < 4 {
			return nil, fmt.Errorf("failed to parse line <%s>: invalid syntax", line)
		}

		if _, ok := parsed[result[0]]; !ok {
			parsed[result[0]] = make(map[string]string)
		}
		parsed[result[0]][result[1]] = result[2]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return parsed, nil
}

func extractFromLine(line string) []string {
	re := regexp.MustCompile(`"([^"]*)"|\*`)

	matches := re.FindAllStringSubmatch(line, -1)

	result := []string{}

	for _, match := range matches {
		if match[1] == "" {
			result = append(result, "*")
		} else {
			result = append(result, match[1])
		}
	}
	return result
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
