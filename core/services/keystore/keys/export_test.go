package keys_export_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocrkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/solkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/terrakey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
)

type KeyType interface {
	ethkey.KeyV2 | csakey.KeyV2 | p2pkey.KeyV2 | vrfkey.KeyV2 | ocrkey.KeyV2 | solkey.Key | terrakey.Key
	ToEncryptedJSON(password string, scryptParams utils.ScryptParams) (export []byte, err error)
	String() string
}

func TestKeyExportImport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		testCase func(*testing.T)
	}{
		{"ethkey", runTestCase[ethkey.KeyV2]},
		{"csakey", runTestCase[csakey.KeyV2]},
		{"p2pkey", runTestCase[p2pkey.KeyV2]},
		{"vrfkey", runTestCase[vrfkey.KeyV2]},
		{"ocrkey", runTestCase[ocrkey.KeyV2]},
		{"solkey", runTestCase[solkey.Key]},
		{"terrakey", runTestCase[solkey.Key]},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			test.testCase(t)
		})
	}
}

func runTestCase[K KeyType](t *testing.T) {
	key, err := createKey[K]()
	require.NoError(t, err)

	json, err := key.ToEncryptedJSON("password", utils.DefaultScryptParams)
	require.NoError(t, err)
	require.NotEmpty(t, json)

	imported, err := decrypt[K](t, json, "password")
	require.NoError(t, err)

	require.Equal(t, key.String(), imported.String())
}

func createKey[K KeyType]() (ret K, err error) {
	switch r := any(&ret).(type) {
	case *ethkey.KeyV2:
		*r, err = ethkey.NewV2()
	case *csakey.KeyV2:
		*r, err = csakey.NewV2()
	case *p2pkey.KeyV2:
		*r, err = p2pkey.NewV2()
	case *vrfkey.KeyV2:
		*r, err = vrfkey.NewV2()
	case *ocrkey.KeyV2:
		*r, err = ocrkey.NewV2()
	case *solkey.Key:
		*r, err = solkey.New()
	case *terrakey.Key:
		*r = terrakey.New()
	}
	return
}

func decrypt[K KeyType](t *testing.T, keyJSON []byte, password string) (ret K, err error) {
	switch r := any(&ret).(type) {
	case *ethkey.KeyV2:
		t.SkipNow()
	case *csakey.KeyV2:
		*r, err = csakey.FromEncryptedJSON(keyJSON, password)
	case *p2pkey.KeyV2:
		*r, err = p2pkey.FromEncryptedJSON(keyJSON, password)
	case *vrfkey.KeyV2:
		*r, err = vrfkey.FromEncryptedJSON(keyJSON, password)
	case *ocrkey.KeyV2:
		*r, err = ocrkey.FromEncryptedJSON(keyJSON, password)
	case *solkey.Key:
		*r, err = solkey.FromEncryptedJSON(keyJSON, password)
	case *terrakey.Key:
		*r, err = terrakey.FromEncryptedJSON(keyJSON, password)
	}
	return
}
