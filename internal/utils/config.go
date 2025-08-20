package utils

import (
	"os"

	"gopkg.in/yaml.v3"
	"github.com/michaeldvinci/syllabus/internal/models"
)

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*models.Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg models.Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}