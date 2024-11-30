package changeset_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

func TestTransferLinkToken_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		config     changeset.TransferLinkTokenConfig
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "valid config",
			config: changeset.TransferLinkTokenConfig{
				Transfers: map[uint64]changeset.Transfer{
					chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector: {
						To:     common.HexToAddress("0x1"),
						Amount: big.NewInt(1),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing to address",
			config: changeset.TransferLinkTokenConfig{
				Transfers: map[uint64]changeset.Transfer{
					chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector: {
						Amount: big.NewInt(1),
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "to address must be set",
		},
		{
			name: "missing amount",
			config: changeset.TransferLinkTokenConfig{
				Transfers: map[uint64]changeset.Transfer{
					chainselectors.ETHEREUM_TESTNET_SEPOLIA.Selector: {
						To: common.HexToAddress("0x1"),
					},
				},
			},
			wantErr:    true,
			wantErrMsg: "amount must be set and positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.config.Validate()
			if test.wantErr {
				if test.wantErrMsg != "" {
					assert.Contains(t, err.Error(), test.wantErrMsg)
				}
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransferLinkToken(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelectorId := env.AllChainSelectors()[0]

	chain := env.Chains[chainSelectorId]
	deployer := chain.DeployerKey
	tokenContract, err := deployment.DeployContract(lggr, chain, env.ExistingAddresses,
		func(chain deployment.Chain) deployment.ContractDeploy[*link_token.LinkToken] {
			tokenAddress, tx, token, err2 := link_token.DeployLinkToken(
				deployer,
				chain.Client,
			)
			return deployment.ContractDeploy[*link_token.LinkToken]{
				tokenAddress, token, tx, deployment.NewTypeAndVersion(types.LinkToken, deployment.Version1_0_0), err2,
			}
		})
	require.NoError(t, err)

	tx, err := tokenContract.Contract.GrantMintRole(deployer, deployer.From)
	require.NoError(t, err)
	_, err = chain.Confirm(tx)

	tx, err = tokenContract.Contract.Mint(deployer, deployer.From, big.NewInt(100))
	require.NoError(t, err)
	_, err = chain.Confirm(tx)
	require.NoError(t, err)

	receiver := common.HexToAddress("0x1")
	_, err = changeset.TransferLinkToken(env,
		changeset.TransferLinkTokenConfig{
			Transfers: map[uint64]changeset.Transfer{
				chainSelectorId: {
					To:     receiver,
					Amount: big.NewInt(30),
				},
			},
		})
	require.NoError(t, err)

	balance, err := tokenContract.Contract.BalanceOf(nil, deployer.From)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(70), balance)

	balance, err = tokenContract.Contract.BalanceOf(nil, receiver)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(30), balance)
}
