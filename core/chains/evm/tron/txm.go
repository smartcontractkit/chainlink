package tron

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	tronkeystore "github.com/smartcontractkit/chainlink-tron/relayer/keystore"
	tronclient "github.com/smartcontractkit/chainlink-tron/relayer/sdk"
	trontxm "github.com/smartcontractkit/chainlink-tron/relayer/txm"
)

func ConstructTronTxm(logger logger.Logger, cfg *config.ChainScoped, nodes []*toml.Node, keystore keys.Store) (*trontxm.TronTxm, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("Tron chain requires at least one node")
	}

	fullNodeURL := nodes[0].HTTPURLExtraWrite.URL()
	tronClient, err := tronclient.CreateFullNodeClient(fullNodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create tron client: %w", err)
	}

	return trontxm.New(logger, tronkeystore.NewLoopKeystoreAdapter(keystore), tronClient, trontxm.TronTxmConfig{
		// From testing, this multiplier ensures all exec messages are fully executed.
		// Energy estimation doesn't seem to account for more complex smart contract execution.
		// Given that Tron has static gas prices, we don't expect this to be a problem as this multiplier is sufficiently high.
		EnergyMultiplier: 3,
		// Maximum number of transactions to buffer in the broadcast channel.
		BroadcastChanSize: 100,
		// Number of seconds to wait between polling the blockchain for transaction confirmation.
		ConfirmPollSecs: 5,
	}), nil
}
