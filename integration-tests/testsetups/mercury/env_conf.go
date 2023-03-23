package mercury

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type MSInfoConf struct {
	RemoteUrl               string            `json:"remoteUrl"`
	LocalUrl                string            `json:"localUrl"`
	RemoteWsrpcUrl          string            `json:"remoteWsrpcUrl"`
	LocalWsrpcUrl           string            `json:"localWsrpcUrl"`
	UserId                  string            `json:"userId"`
	UserKey                 string            `json:"userKey"`
	UserEncryptedKey        string            `json:"userEncryptedKey"`
	RpcPubKey               ed25519.PublicKey `json:"rpcPubKey"`
	RpcNodesCsaPrivKeySeeds []string          `json:"rpcNodesCsaPrivKeys"`
}

type TestEnvConfig struct {
	Id            string                `json:"id"`
	K8Namespace   string                `json:"k8Namespace"`
	FeedId        string                `json:"feedId"`
	ChainId       int64                 `json:"chainId"`
	ContractsInfo map[string]string     `json:"contracts"`
	MSInfo        MSInfoConf            `json:"mercuryServer"`
	Actions       map[string]*envAction `json:"actions"`
}

func (c *TestEnvConfig) Save() (string, error) {
	// Create mercury env log dir if necessary
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	confDir := fmt.Sprintf("%s/logs", pwd)
	if _, err := os.Stat(confDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(confDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Save mercury env config to disk
	confPath := fmt.Sprintf("%s/%s-%s.json", confDir, c.Id, c.K8Namespace)
	f, _ := json.MarshalIndent(c, "", "   ")
	err = ioutil.WriteFile(confPath, f, 0644)

	return confPath, err
}
