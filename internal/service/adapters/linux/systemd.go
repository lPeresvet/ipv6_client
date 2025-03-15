package linux

import (
	"fmt"
	"implementation/internal/domain/connections"
	"implementation/internal/parsers"
	"os/exec"
)

const stateFieldName = "ActiveState"

type SystemdProvider struct{}

func NewSystemdProvider() *SystemdProvider {
	return &SystemdProvider{}
}

func (s *SystemdProvider) StartDemon(demonName string) error {
	if err := exec.Command("systemctl", "start", demonName).Run(); err != nil {
		return fmt.Errorf("failed to start %s demon: %w", demonName, err)
	}

	return nil
}

func (s *SystemdProvider) StopDemon(demonName string) error {
	if err := exec.Command("systemctl", "stop", demonName).Run(); err != nil {
		return fmt.Errorf("failed to stop %s demon: %w", demonName, err)
	}

	return nil
}

func (s *SystemdProvider) DemonStatus(demonName string) (*connections.DemonInfo, error) {
	output, err := exec.Command("systemctl", "show", "--no-pager", demonName).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s status: %w", demonName, err)
	}

	info, err := parsers.ParseIni(string(output))
	if err != nil {
		return nil, err
	}
	status := connections.DemonInactive

	if val := info[stateFieldName]; val == string(connections.DemonActive) {
		status = connections.DemonActive
	}

	return &connections.DemonInfo{Status: status}, nil
}
