package crecli

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	SharedSecretEnvVarSuffix = "_ENV_VAR_ALL"
)

type Secrets struct {
	SecretsNames map[string][]string `yaml:"secretsNames"`
}

func CreateSecretsFile(secretsNames map[string][]string) (*os.File, error) {
	secretsFile, cErr := os.CreateTemp("", "secrets.config.yaml")
	if cErr != nil {
		return nil, errors.Wrapf(cErr, "failed to create secrets file")
	}

	secrets := Secrets{
		SecretsNames: secretsNames,
	}

	secretsYAML, mErr := yaml.Marshal(secrets)
	if mErr != nil {
		return nil, errors.Wrapf(mErr, "failed to marshal secrets")
	}

	_, writeErr := secretsFile.Write(secretsYAML)
	if writeErr != nil {
		return nil, errors.Wrapf(writeErr, "failed to write secrets to file")
	}

	return secretsFile, nil
}
