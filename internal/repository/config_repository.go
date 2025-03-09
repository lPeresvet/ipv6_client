package repository

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
	"implementation/internal/domain/config"
)

type FileRepository struct {
	path string
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		path: path,
	}
}

func (r *FileRepository) GetConfig(name string) (*config.Config, error) {
	file, err := os.Open(path.Join(r.path, name))
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from config file: %w", err)
	}

	connectionConfig := config.Config{}

	if err := yaml.Unmarshal(data, &connectionConfig); err != nil {
		log.Fatalf("error: %v", err)
	}

	return &connectionConfig, nil
}
