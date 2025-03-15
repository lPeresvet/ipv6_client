package linux

import (
	"fmt"
	"implementation/internal/domain/connections"
	"log"
	"os/exec"
)

type SystemdProvider struct{}

func NewSystemdProvider() *SystemdProvider {
	return &SystemdProvider{}
}

func (s *SystemdProvider) StartDemon(demonName string) error {
	if err := exec.Command("systemctl", "start", demonName).Run(); err != nil {
		return fmt.Errorf("failed to start demon: %w", err)
	}

	return nil
}

func (s *SystemdProvider) StopDemon(demonName string) error {
	if err := exec.Command("systemctl", "stop", demonName).Run(); err != nil {
		return fmt.Errorf("failed to stop demon: %w", err)
	}

	return nil
}

func (s *SystemdProvider) DemonStatus(demonName string) (*connections.DemonInfo, error) {
	output, err := exec.Command("systemctl", "show", "--no-pager", demonName).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s status: %w", demonName, err)
	}

	log.Printf("<%s>", string(output))

	return &connections.DemonInfo{}, nil
}
