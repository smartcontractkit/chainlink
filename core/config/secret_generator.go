package config

import (
	"encoding/base64"
	"io/ioutil"
	"path/filepath"

	"github.com/gorilla/securecookie"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// SecretGenerator is the interface for objects that generate a secret
// used to sign or encrypt.
type SecretGenerator interface {
	Generate(string) ([]byte, error)
}

type FilePersistedSecretGenerator struct{}

func (f FilePersistedSecretGenerator) Generate(rootDir string) ([]byte, error) {
	sessionPath := filepath.Join(rootDir, "secret")
	if utils.FileExists(sessionPath) {
		data, err := ioutil.ReadFile(sessionPath)
		if err != nil {
			return data, err
		}
		return base64.StdEncoding.DecodeString(string(data))
	}
	key := securecookie.GenerateRandomKey(32)
	str := base64.StdEncoding.EncodeToString(key)
	err := utils.WriteFileWithMaxPerms(sessionPath, []byte(str), readWritePerms)
	return key, err
}
