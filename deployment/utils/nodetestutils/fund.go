package nodetestutils

import (
	"testing"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"

	"github.com/smartcontractkit/chainlink/deployment/internal/aptostestutils"

	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf_ton "github.com/smartcontractkit/chainlink-deployments-framework/chain/ton"
)

var (
	// tonTransmitterFundAmount is the amount of TON to fund the transmitter address.
	tonTransmitterFundAmount = tlb.MustFromTON("200")
)

// fundNodesTon funds the given nodes with the given amount of TON.
func fundNodesTon(t *testing.T, tonChain cldf_ton.Chain, nodes []*Node) {
	messages := make([]*wallet.Message, 0, len(nodes))
	for _, node := range nodes {
		tonkeys, err := node.App.GetKeyStore().TON().GetAll()
		require.NoError(t, err)
		require.Len(t, tonkeys, 1)
		transmitter := tonkeys[0].PubkeyToAddress()
		msg, err := tonChain.Wallet.BuildTransfer(transmitter, tonTransmitterFundAmount, false, "")
		require.NoError(t, err)
		messages = append(messages, msg)
	}
	_, _, err := tonChain.Wallet.SendManyWaitTransaction(t.Context(), messages)
	require.NoError(t, err)
}

// fundNodesAptos funds the given nodes with the given amount of APT.
func fundNodesAptos(t *testing.T, aptosChain cldf_aptos.Chain, nodes []*Node) {
	for _, node := range nodes {
		aptoskeys, err := node.App.GetKeyStore().Aptos().GetAll()
		require.NoError(t, err)
		require.Len(t, aptoskeys, 1)
		transmitter := aptoskeys[0]
		transmitterAccountAddress := aptos.AccountAddress{}
		require.NoError(t, transmitterAccountAddress.ParseStringRelaxed(transmitter.Account()))
		aptostestutils.FundAccount(t, aptosChain.DeployerSigner, transmitterAccountAddress, 100*1e8, aptosChain.Client)
	}
}

// fundNodesSol funds the given nodes with the given amount of SOL.
func fundNodesSol(t *testing.T, solChain cldf_solana.Chain, nodes []*Node) {
	for _, node := range nodes {
		solkeys, err := node.App.GetKeyStore().Solana().GetAll()
		require.NoError(t, err)
		require.Len(t, solkeys, 1)
		transmitter := solkeys[0]
		_, err = solChain.Client.RequestAirdrop(t.Context(), transmitter.PublicKey(), 1000*solana.LAMPORTS_PER_SOL, solRpc.CommitmentConfirmed)
		require.NoError(t, err)
		// we don't wait for confirmation so we don't block the tests, it'll take a while before nodes start transmitting
	}
}
