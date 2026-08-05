package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadCacheClient(t *testing.T) {
	env, _, _ := newTestEnv(t)
	seedCacheRef(t, &env, testContractID, "test-cache", "1.0.0")

	client, ref, err := LoadCacheClient(env, testChainSel, "test-cache", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, testContractID, ref.Address)
	require.Equal(t, testContractID, client.ContractID())

	// missing ref must fail
	_, _, err = LoadCacheClient(env, testChainSel, "does-not-exist", "1.0.0")
	require.Error(t, err)

	// unknown chain must fail
	_, _, err = LoadCacheClient(env, 999999, "test-cache", "1.0.0")
	require.Error(t, err)
}

func TestLoadProxyClient(t *testing.T) {
	env, _, _ := newTestEnv(t)
	seedProxyRef(t, &env, testProxyAddress, "test-proxy", "1.0.0")

	client, ref, err := LoadProxyClient(env, testChainSel, "test-proxy", "1.0.0")
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Equal(t, testProxyAddress, ref.Address)
	require.Equal(t, testProxyAddress, client.ContractID())

	// missing ref must fail
	_, _, err = LoadProxyClient(env, testChainSel, "does-not-exist", "1.0.0")
	require.Error(t, err)

	// unknown chain must fail
	_, _, err = LoadProxyClient(env, 999999, "test-proxy", "1.0.0")
	require.Error(t, err)
}
