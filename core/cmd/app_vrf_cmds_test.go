package cmd_test

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"

	"chainlink/core/store/models/vrfkey"

	"github.com/stretchr/testify/require"
)

// XXX: Remove before merging. See
// https://github.com/smartcontractkit/chainlink-private/pull/17#discussion_r366448801
func TestImportVRFKeyFixture(t *testing.T) {
	rawPassword, err := ioutil.ReadFile("../../tools/clroot/password.txt")
	require.NoError(t, err)
	password := strings.TrimSpace(string(rawPassword))
	keyfile, err := ioutil.ReadFile("../../tools/clroot/vrfkey.json")
	require.NoError(t, err)
	encKey := vrfkey.EncryptedSecretKey{}
	require.NoError(t, json.Unmarshal(keyfile, &encKey))
	_, err = encKey.Decrypt(password)
	require.NoError(t, err)
}
