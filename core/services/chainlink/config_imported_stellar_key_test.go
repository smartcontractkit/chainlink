package chainlink

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

func TestImportedStellarKeys_List(t *testing.T) {
	t.Parallel()

	chainIDs := firstNStellarChainIDs(t, 2)

	secrets, err := parseImportedStellarSecrets(fmt.Sprintf(`
[Stellar]
[[Stellar.Keys]]
JSON = '{"id":"stellar-key-1"}'
ID = %q
Password = 'pw-1'

[[Stellar.Keys]]
JSON = '{"id":"stellar-key-2"}'
ID = %q
Password = 'pw-2'
`, chainIDs[0], chainIDs[1]))
	require.NoError(t, err)

	cfg := &generalConfig{secrets: secrets}
	keys := cfg.ImportedStellarKeys().List()
	require.Len(t, keys, 2)

	expected0, err := chain_selectors.GetChainDetailsByChainIDAndFamily(chainIDs[0], chain_selectors.FamilyStellar)
	require.NoError(t, err)
	expected1, err := chain_selectors.GetChainDetailsByChainIDAndFamily(chainIDs[1], chain_selectors.FamilyStellar)
	require.NoError(t, err)

	require.JSONEq(t, `{"id":"stellar-key-1"}`, keys[0].JSON())
	require.Equal(t, "pw-1", keys[0].Password())
	require.Equal(t, expected0, keys[0].ChainDetails())

	require.JSONEq(t, `{"id":"stellar-key-2"}`, keys[1].JSON())
	require.Equal(t, "pw-2", keys[1].Password())
	require.Equal(t, expected1, keys[1].ChainDetails())
}

func TestImportedStellarKeys_ValidateRejectsUnknownChainID(t *testing.T) {
	t.Parallel()

	var secrets Secrets
	err := commonconfig.DecodeTOML(strings.NewReader(`
[Stellar]
[[Stellar.Keys]]
JSON = '{"id":"stellar-key-1"}'
ID = 'unknown-stellar-chain-id'
Password = 'pw-1'
	`), &secrets)
	require.NoError(t, err)
	require.ErrorContains(t, secrets.Stellar.ValidateConfig(), "invalid StellarKey")
}

func parseImportedStellarSecrets(secretsTOML string) (*Secrets, error) {
	var secrets Secrets
	if err := commonconfig.DecodeTOML(strings.NewReader(secretsTOML), &secrets); err != nil {
		return nil, err
	}

	return &secrets, nil
}

func firstNStellarChainIDs(t *testing.T, n int) []string {
	t.Helper()

	chainIDToSelector := chain_selectors.StellarChainIdToChainSelector()
	require.GreaterOrEqual(t, len(chainIDToSelector), n)

	chainIDs := make([]string, 0, len(chainIDToSelector))
	for chainID := range chainIDToSelector {
		chainIDs = append(chainIDs, chainID)
	}
	sort.Strings(chainIDs)

	return chainIDs[:n]
}
