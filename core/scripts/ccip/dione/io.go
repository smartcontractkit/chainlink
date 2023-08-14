package dione

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/smartcontractkit/chainlink/core/scripts/common"
)

type Folder string

const (
	JSON_FOLDER        Folder = "json"
	CHAIN_FOLDER       Folder = "chain"
	NODES_FOLDER       Folder = "nodes"
	CREDENTIALS_FOLDER Folder = "credentials"
)

func getFileLocation(env Environment, folder Folder) string {
	return fmt.Sprintf("%s/%s/%s.json", JSON_FOLDER, folder, env)
}

func ReadCredentials(env Environment) (DonCredentials, error) {
	path := getFileLocation(env, CREDENTIALS_FOLDER)
	jsonFile, err := os.Open(path)
	if err != nil {
		return DonCredentials{}, err
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return DonCredentials{}, err
	}

	var creds DonCredentials
	err = json.Unmarshal(byteValue, &creds)

	return creds, err
}

func MustReadNodeConfig(env Environment) NodesConfig {
	config, err := ReadNodeConfig(env)
	common.PanicErr(err)
	return config
}

func ReadNodeConfig(env Environment) (NodesConfig, error) {
	path := getFileLocation(env, NODES_FOLDER)
	jsonFile, err := os.Open(path)
	if err != nil {
		return NodesConfig{}, err
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return NodesConfig{}, err
	}

	var config NodesConfig
	err = json.Unmarshal(byteValue, &config)

	return config, err
}

func WriteJSON(path string, file []byte) error {
	return os.WriteFile(path, file, 0600)
}
