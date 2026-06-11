package drift

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type exception struct {
	File     string `yaml:"file"`
	Contains string `yaml:"contains"`
	Reason   string `yaml:"reason"`
}

func loadExceptions(root string) ([]exception, error) {
	path := filepath.Join(root, "tools", "version-exceptions.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []exception
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
