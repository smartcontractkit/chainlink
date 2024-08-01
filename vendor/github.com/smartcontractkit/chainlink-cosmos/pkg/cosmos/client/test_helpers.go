package client

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/gjson"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/testutil"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Account struct {
	Name       string
	PrivateKey cryptotypes.PrivKey
	Address    sdk.AccAddress
}

// 0.001
var defaultCoin = sdk.NewDecWithPrec(1, 3)

// SetupLocalCosmosNode sets up a local terra node via wasmd, and returns pre-funded accounts, the test directory, and the url.
// Token name is for both staking and fee coin
func SetupLocalCosmosNode(t *testing.T, chainID string, token string) ([]Account, string, string) {
	minGasPrice := sdk.NewDecCoinFromDec(token, defaultCoin)
	testdir, err := os.MkdirTemp("", "integration-test")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(testdir))
	})
	t.Log(testdir)
	out, err := exec.Command("wasmd", "init", "integration-test", "-o", "--chain-id", chainID, "--home", testdir).Output()
	require.NoError(t, err, string(out))

	p := path.Join(testdir, "config", "app.toml")
	f, err := os.ReadFile(p)
	require.NoError(t, err)
	config, err := toml.Load(string(f))
	require.NoError(t, err)
	// Enable if desired to use lcd endpoints config.Set("api.enable", "true")
	config.Set("minimum-gas-prices", minGasPrice.String())
	require.NoError(t, os.WriteFile(p, []byte(config.String()), 0600))
	// TODO: could also speed up the block mining config

	p = path.Join(testdir, "config", "genesis.json")
	f, err = os.ReadFile(p)
	require.NoError(t, err)

	genesisData := string(f)
	// fix hardcoded staking/governance token, see
	// https://github.com/CosmWasm/wasmd/blob/develop/docker/setup_wasmd.sh
	// https://github.com/CosmWasm/wasmd/blob/develop/contrib/local/setup_wasmd.sh
	newStakingToken := fmt.Sprintf(`"%s"`, token)
	genesisData = strings.ReplaceAll(genesisData, "\"ustake\"", newStakingToken)
	genesisData = strings.ReplaceAll(genesisData, "\"stake\"", newStakingToken)
	require.NoError(t, os.WriteFile(p, []byte(genesisData), 0600))

	// Create 2 test accounts
	var accounts []Account
	for i := 0; i < 2; i++ {
		account := fmt.Sprintf("test%d", i)
		key, err2 := exec.Command("wasmd", "keys", "add", account, "--output", "json", "--keyring-backend", "test", "--home", testdir).CombinedOutput()
		require.NoError(t, err2, string(key))
		var k struct {
			Address  string `json:"address"`
			Mnemonic string `json:"mnemonic"`
		}
		require.NoError(t, json.Unmarshal(key, &k))
		expAcctAddr, err3 := sdk.AccAddressFromBech32(k.Address)
		require.NoError(t, err3)
		privateKey, address, err4 := testutil.CreateKeyFromMnemonic(k.Mnemonic)
		require.NoError(t, err4)
		require.Equal(t, expAcctAddr, address)
		// Give it 100000000ucosm
		out2, err2 := exec.Command("wasmd", "genesis", "add-genesis-account", k.Address, "100000000"+token, "--home", testdir).Output() //nolint:gosec
		require.NoError(t, err2, string(out2))
		accounts = append(accounts, Account{
			Name:       account,
			Address:    address,
			PrivateKey: privateKey,
		})
	}
	// Stake 10 tokens in first acct
	out, err = exec.Command("wasmd", "genesis", "gentx", accounts[0].Name, "10000000"+token, "--chain-id", chainID, "--keyring-backend", "test", "--home", testdir).CombinedOutput() //nolint:gosec
	require.NoError(t, err, string(out))
	out, err = exec.Command("wasmd", "genesis", "collect-gentxs", "--home", testdir).CombinedOutput()
	require.NoError(t, err, string(out))

	port := mustRandomPort()
	tendermintHost := fmt.Sprintf("127.0.0.1:%d", port)
	tendermintURL := "http://" + tendermintHost
	t.Log(tendermintURL)

	cmd := exec.Command("wasmd", "start", "--home", testdir,
		"--rpc.laddr", "tcp://"+tendermintHost,
		"--rpc.pprof_laddr", "127.0.0.1:0",
		"--grpc.address", "127.0.0.1:0",
		"--grpc-web.address", "127.0.0.1:0",
		"--p2p.laddr", "127.0.0.1:0")
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr
	require.NoError(t, cmd.Start())
	t.Cleanup(func() {
		assert.NoError(t, cmd.Process.Kill())
		if err2 := cmd.Wait(); assert.Error(t, err2) {
			if !assert.Contains(t, err2.Error(), "signal: killed", cmd.ProcessState.String()) {
				t.Log("wasmd stderr:", stdErr.String())
			}
		}
	})

	// Wait for api server to boot
	var ready bool
	for i := 0; i < 30; i++ {
		time.Sleep(time.Second)
		out, err = exec.Command("curl", tendermintURL+"/abci_info").Output() //nolint:gosec
		if err != nil {
			t.Logf("API server not ready yet (attempt %d): %v\n", i+1, err)
			continue
		}
		var a struct {
			Result struct {
				Response struct {
					LastBlockHeight string `json:"last_block_height"`
				} `json:"response"`
			} `json:"result"`
		}
		require.NoError(t, json.Unmarshal(out, &a), string(out))
		if a.Result.Response.LastBlockHeight == "" {
			t.Logf("API server not ready yet (attempt %d)\n", i+1)
			continue
		}
		ready = true
		break
	}
	require.True(t, ready)
	return accounts, testdir, tendermintURL
}

// DeployTestContract deploys a test contract.
func DeployTestContract(t *testing.T, tendermintURL, chainID string, token string, deployAccount, ownerAccount Account, tc *Client, testdir, wasmTestContractPath string) sdk.AccAddress {
	minGasPrice := sdk.NewDecCoinFromDec(token, defaultCoin)
	//nolint:gosec
	submitResp, err2 := exec.Command("wasmd", "tx", "wasm", "store", wasmTestContractPath, "--node", tendermintURL,
		"--from", deployAccount.Name, "--gas", "auto", "--fees", "100000"+token, "--gas-adjustment", "1.3", "--chain-id", chainID, "--home", testdir, "--keyring-backend", "test", "--keyring-dir", testdir, "--yes", "--output", "json").Output()
	require.NoError(t, err2, string(submitResp))

	// wait for tx to be committed
	txHash := gjson.Get(string(submitResp), "txhash")
	require.True(t, txHash.Exists())
	storeTx, success := AwaitTxCommitted(t, tc, txHash.String())
	require.True(t, success)

	// get code id from tx receipt
	storeCodeLog := storeTx.TxResponse.Logs[len(storeTx.TxResponse.Logs)-1]
	codeID, err := strconv.ParseUint(storeCodeLog.GetEvents()[1].Attributes[1].Value, 10, 64)
	require.NoError(t, err, "failed to parse code id from tx receipt")

	accountNumber, sequenceNumber, err := tc.Account(ownerAccount.Address)
	require.NoError(t, err)
	deployTx, err3 := tc.SignAndBroadcast([]sdk.Msg{
		&wasmtypes.MsgInstantiateContract{
			Sender: ownerAccount.Address.String(),
			Admin:  "",
			CodeID: codeID,
			Label:  "testcontract",
			Msg:    []byte(`{"count":0}`),
			Funds:  sdk.Coins{},
		},
	}, accountNumber, sequenceNumber, minGasPrice, ownerAccount.PrivateKey, txtypes.BroadcastMode_BROADCAST_MODE_SYNC)
	require.NoError(t, err3)

	// wait for tx to be committed
	deployTxReceipt, success := AwaitTxCommitted(t, tc, deployTx.TxResponse.TxHash)
	require.True(t, success)

	return GetContractAddr(t, deployTxReceipt.GetTxResponse())
}

func GetContractAddr(t *testing.T, deployTxReceipt *sdk.TxResponse) sdk.AccAddress {
	var contractAddr string
	for _, etype := range deployTxReceipt.Events {
		if etype.Type == "wasm" {
			for _, attr := range etype.Attributes {
				if attr.Key == "_contract_address" {
					contractAddr = attr.Value
				}
			}
		}
	}
	require.NotEqual(t, "", contractAddr)
	contract, err := sdk.AccAddressFromBech32(contractAddr)
	require.NoError(t, err)
	return contract
}

func mustRandomPort() int {
	r, err := rand.Int(rand.Reader, big.NewInt(65535-1023))
	if err != nil {
		panic(fmt.Errorf("unexpected error generating random port: %w", err))
	}
	return int(r.Int64() + 1024)
}

// AwaitTxCommitted waits for a transaction to be committed on chain and returns the tx receipt
func AwaitTxCommitted(t *testing.T, tc *Client, txHash string) (response *txtypes.GetTxResponse, success bool) {
	for i := 0; i < 10; i++ { // max poll attempts to wait for tx commitment
		txReceipt, err := tc.Tx(txHash)
		if err == nil {
			return txReceipt, true
		}
		time.Sleep(time.Second * 1) // TODO: configure dynamically based on block times
	}
	return nil, false
}
