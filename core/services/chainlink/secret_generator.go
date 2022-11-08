package chainlink

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gorilla/securecookie"

	"github.com/smartcontractkit/chainlink/core/utils"
)

// this permission grants read / write access to file owners only
const readWritePerms = os.FileMode(0600)

// SecretGenerator is the interface for objects that generate a secret
// used to sign or encrypt.
type SecretGenerator interface {
	Generate(string) ([]byte, error)
}

type FilePersistedSecretGenerator struct{}

func (f FilePersistedSecretGenerator) Generate(rootDir string) ([]byte, error) {
	fmt.Println("STARTING GENERATION")
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Printf("CWD '%s'\n", wd)
	sessionPath := filepath.Join(rootDir, "secret")
	fmt.Printf("GENERATION PATH '%s' Root '%s'\n", sessionPath, rootDir)
	if exists, err := utils.FileExists(sessionPath); err != nil {
		return nil, err
	} else if exists {
		data, err := os.ReadFile(sessionPath)
		if err != nil {
			return data, err
		}
		return base64.StdEncoding.DecodeString(string(data))
	}
	fmt.Printf("File Doesn't Exist, Writing One '%s'\n", filepath.Join(wd, sessionPath))
	key := securecookie.GenerateRandomKey(32)
	str := base64.StdEncoding.EncodeToString(key)
	err = utils.WriteFileWithMaxPerms(sessionPath, []byte(str), readWritePerms)
	return key, err
}
