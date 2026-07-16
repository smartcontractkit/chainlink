package chainlink

import (
	"strings"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

func TestImportedAptosKeys_List(t *testing.T) {
	t.Parallel()

	secrets, err := parseImportedAptosSecrets(`
[Aptos]
[[Aptos.Keys]]
JSON = '{"id":"aptos-key-1"}'
ID = 4
Password = 'pw-1'

[[Aptos.Keys]]
JSON = '{"id":"aptos-key-2"}'
ID = 1
Password = 'pw-2'
`)
	require.NoError(t, err)

	cfg := &generalConfig{secrets: secrets}
	keys := cfg.ImportedAptosKeys().List()
	require.Len(t, keys, 2)

	expected4, err := chain_selectors.GetChainDetailsByChainIDAndFamily("4", chain_selectors.FamilyAptos)
	require.NoError(t, err)
	expected1, err := chain_selectors.GetChainDetailsByChainIDAndFamily("1", chain_selectors.FamilyAptos)
	require.NoError(t, err)

	require.JSONEq(t, `{"id":"aptos-key-1"}`, keys[0].JSON())
	require.Equal(t, "pw-1", keys[0].Password())
	require.Equal(t, expected4, keys[0].ChainDetails())

	require.JSONEq(t, `{"id":"aptos-key-2"}`, keys[1].JSON())
	require.Equal(t, "pw-2", keys[1].Password())
	require.Equal(t, expected1, keys[1].ChainDetails())
}

func TestImportedAptosKeys_ValidateRejectsUnknownChainID(t *testing.T) {
	t.Parallel()

	var secrets Secrets
	err := commonconfig.DecodeTOML(strings.NewReader(`
[Aptos]
[[Aptos.Keys]]
JSON = '{"id":"aptos-key-1"}'
ID = 999999
Password = 'pw-1'
	`), &secrets)
	require.NoError(t, err)
	require.ErrorContains(t, secrets.Aptos.ValidateConfig(), "invalid AptosKey")
}

func parseImportedAptosSecrets(secretsTOML string) (*Secrets, error) {
	var secrets Secrets
	if err := commonconfig.DecodeTOML(strings.NewReader(secretsTOML), &secrets); err != nil {
		return nil, err
	}

	return &secrets, nil
}
