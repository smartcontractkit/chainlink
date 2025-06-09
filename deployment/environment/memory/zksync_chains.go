package memory

import (
	"context"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"

	cs "github.com/smartcontractkit/chain-selectors"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
)

func GenerateChainsZk(t *testing.T, numChains int) map[uint64]cldf_evm.Chain {
	chains := make(map[uint64]cldf_evm.Chain)

	for i := 0; i < numChains; i++ {
		chainID := cs.TEST_90000051.EvmChainID + uint64(i) //nolint:gosec // it shouldn't overflow

		output, err := blockchain.NewBlockchainNetwork(&blockchain.Input{
			Type:    "anvil-zksync",
			ChainID: strconv.FormatUint(chainID, 10),
			Port:    strconv.FormatInt(int64(freeport.GetN(t, 1)[0]), 10),
		})
		require.NoError(t, err)

		testcontainers.CleanupContainer(t, output.Container)

		sel, err := cs.SelectorFromChainId(chainID)
		require.NoError(t, err)

		client, err := ethclient.Dial(output.Nodes[0].ExternalHTTPUrl)
		require.NoError(t, err)

		gasPrice, err := client.SuggestGasPrice(t.Context())
		require.NoError(t, err)

		require.Greater(t, len(blockchain.AnvilZKSyncRichAccountPks), 1)

		keyedTransactors := make([]*bind.TransactOpts, 0)
		for _, pk := range blockchain.AnvilZKSyncRichAccountPks {
			privateKey, err := crypto.HexToECDSA(pk)
			require.NoError(t, err)
			transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, new(big.Int).SetUint64(chainID))
			transactor.GasPrice = gasPrice
			require.NoError(t, err)
			keyedTransactors = append(keyedTransactors, transactor)
		}

		require.Equal(t, len(blockchain.AnvilZKSyncRichAccountPks), len(keyedTransactors))

		clientZk := clients.NewClient(client.Client())
		deployerZk, err := accounts.NewWallet(common.Hex2Bytes(blockchain.AnvilZKSyncRichAccountPks[0]), clientZk, nil)
		require.NoError(t, err)

		chain := cldf_evm.Chain{
			Selector:    sel,
			Client:      client,
			DeployerKey: keyedTransactors[0], // to use to interact with contracts
			Users:       keyedTransactors[1:],
			Confirm: func(tx *types.Transaction) (uint64, error) {
				ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
				defer cancel()
				receipt, err := bind.WaitMined(ctx, client, tx)
				if err != nil {
					return 0, err
				}
				return receipt.Status, nil
			},
			IsZkSyncVM:          true,
			ClientZkSyncVM:      clientZk,
			DeployerKeyZkSyncVM: deployerZk,
		}

		chains[sel] = chain
	}

	return chains
}
