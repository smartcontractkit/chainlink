package cmd_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"chainlink/core/store/models/vrf_key"

	"github.com/stretchr/testify/require"
)

func TestImportVRFKeyFixture(t *testing.T) {
	rawPassword, err := ioutil.ReadFile("../../tools/clroot/password.txt")
	require.NoError(t, err)
	password := strings.TrimSpace(string(rawPassword))
	keyfile, err := ioutil.ReadFile("../../tools/clroot/vrfkey.json")
	require.NoError(t, err)
	encKey := vrf_key.EncryptedSecretKey{}
	require.NoError(t, json.Unmarshal(keyfile, &encKey))
	_, err = encKey.Decrypt(password)
	require.NoError(t, err)
}
