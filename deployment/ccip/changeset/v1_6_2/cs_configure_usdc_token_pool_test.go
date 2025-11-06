package v1_6_2_test

import (
	"maps"
	"math/big"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/v1_6_2"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_messenger"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/latest/mock_usdc_token_transmitter"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

func setupUSDCTokenPoolsEnvironmentForConfigure(t *testing.T, withPrereqs bool) *runtime.Runtime {
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulatedN(t, 2),
		environment.WithSolanaContainerN(t, 1, t.TempDir(), map[string]string{}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	selectors := rt.Environment().BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	if withPrereqs {
		var err error

		prereqCfg := make([]changeset.DeployPrerequisiteConfigPerChain, len(selectors))
		for i, selector := range selectors {
			prereqCfg[i] = changeset.DeployPrerequisiteConfigPerChain{
				ChainSelector: selector,
			}
		}

		err = rt.Exec(
			runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployPrerequisitesChangeset), changeset.DeployPrerequisiteConfig{
				Configs: prereqCfg,
			}),
		)
		require.NoError(t, err)
	}

	return rt
}

func setupUSDCTokenPoolsContractsForConfigure(
	t *testing.T,
	logger logger.Logger,
	chain cldf_evm.Chain,
	addressBook cldf.AddressBook,
) (
	*cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677],
	*cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger],
) {
	usdcToken, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677] {
			tokenAddress, tx, token, err := burn_mint_erc677.DeployBurnMintERC677(
				chain.DeployerKey,
				chain.Client,
				"USDC",
				"USDC",
				6,
				big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(1e18)),
			)
			return cldf.ContractDeploy[*burn_mint_erc677.BurnMintERC677]{
				Address:  tokenAddress,
				Contract: token,
				Tv:       cldf.NewTypeAndVersion(shared.BurnMintToken, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	transmitter, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter] {
			transmitterAddress, tx, transmitter, err := mock_usdc_token_transmitter.DeployMockE2EUSDCTransmitter(chain.DeployerKey, chain.Client, 0, 1, usdcToken.Address)
			return cldf.ContractDeploy[*mock_usdc_token_transmitter.MockE2EUSDCTransmitter]{
				Address:  transmitterAddress,
				Contract: transmitter,
				Tv:       cldf.NewTypeAndVersion(shared.USDCMockTransmitter, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	messenger, err := cldf.DeployContract(logger, chain, addressBook,
		func(chain cldf_evm.Chain) cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger] {
			messengerAddress, tx, messenger, err := mock_usdc_token_messenger.DeployMockE2EUSDCTokenMessenger(chain.DeployerKey, chain.Client, 0, transmitter.Address)
			return cldf.ContractDeploy[*mock_usdc_token_messenger.MockE2EUSDCTokenMessenger]{
				Address:  messengerAddress,
				Contract: messenger,
				Tv:       cldf.NewTypeAndVersion(shared.USDCTokenMessenger, deployment.Version1_6_2),
				Tx:       tx,
				Err:      err,
			}
		},
	)
	require.NoError(t, err)

	return usdcToken, messenger
}

func TestValidateConfigUSDCTokenPoolInput(t *testing.T) {
	t.Parallel()

	rt := setupUSDCTokenPoolsEnvironmentForConfigure(t, true)

	solChains := rt.Environment().BlockChains.SolanaChains()
	require.GreaterOrEqual(t, len(solChains), 1)
	solChain := slices.Collect(maps.Values(solChains))[0]

	evmChains := rt.Environment().BlockChains.EVMChains()
	require.GreaterOrEqual(t, len(evmChains), 1)
	evmChain := slices.Collect(maps.Values(evmChains))[0]

	addressBook := cldf.NewMemoryAddressBook()
	usdcToken, tokenMsngr := setupUSDCTokenPoolsContractsForConfigure(t, rt.Environment().Logger, evmChain, addressBook)

	err := rt.Exec(
		runtime.ChangesetTask(v1_6_2.DeployCCTPMessageTransmitterProxyNew, v1_6_2.DeployCCTPMessageTransmitterProxyContractConfig{
			USDCProxies: map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput{
				evmChain.Selector: {
					TokenMessenger: tokenMsngr.Address,
				},
			},
		}),
		runtime.ChangesetTask(v1_6_2.DeployUSDCTokenPoolNew, v1_6_2.DeployUSDCTokenPoolContractsConfig{
			USDCPools: map[uint64]v1_6_2.DeployUSDCTokenPoolInput{
				evmChain.Selector: {
					PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
					TokenMessenger:      tokenMsngr.Address,
					TokenAddress:        usdcToken.Address,
					PoolType:            shared.USDCTokenPool,
				},
			},
		}),
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(rt.Environment())
	require.NoError(t, err)

	minterPrivKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	callerPrivKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	dummyDomainID := uint32(0)
	tests := []struct {
		Msg    string
		Input  v1_6_2.ConfigUSDCTokenPoolInput
		ErrStr string
	}{
		{
			Msg: "Invalid chain selector",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					0: {
						MintRecipient:    "",
						AllowedCaller:    "",
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "invalid destination chain selector",
		},
		{
			Msg: "Solana mint recipient cannot be empty string",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					solChain.Selector: {
						MintRecipient:    "",
						AllowedCaller:    callerPrivKey.PublicKey().String(),
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "invalid mint recipient format",
		},
		{
			Msg: "Solana mint recipient cannot be zero address",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					solChain.Selector: {
						MintRecipient:    solana.PublicKey{}.String(),
						AllowedCaller:    callerPrivKey.PublicKey().String(),
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "mint recipient must be defined for Solana destination chain selector",
		},
		{
			Msg: "Solana allowed caller cannot be empty string",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					solChain.Selector: {
						MintRecipient:    minterPrivKey.PublicKey().String(),
						AllowedCaller:    "",
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "invalid allowed caller format",
		},
		{
			Msg: "Solana allowed caller cannot be zero address",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					solChain.Selector: {
						MintRecipient:    minterPrivKey.PublicKey().String(),
						AllowedCaller:    solana.PublicKey{}.String(),
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "allowed caller must be defined for Solana destination chain selector",
		},
		{
			Msg: "EVM allowed caller cannot be empty string",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					evmChain.Selector: {
						AllowedCaller:    "",
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "allowed caller must be defined for EVM destination chain selector",
		},
		{
			Msg: "EVM allowed caller cannot be zero address",
			Input: v1_6_2.ConfigUSDCTokenPoolInput{
				DestinationUpdates: map[uint64]v1_6_2.DomainUpdateInput{
					evmChain.Selector: {
						AllowedCaller:    utils.ZeroAddress.String(),
						DomainIdentifier: dummyDomainID,
						Enabled:          true,
					},
				},
			},
			ErrStr: "allowed caller must be defined for EVM destination chain selector",
		},
	}

	for _, test := range tests {
		t.Run(test.Msg, func(t *testing.T) {
			err := test.Input.Validate(t.Context(), evmChain, state.Chains[evmChain.Selector])
			require.ErrorContains(t, err, test.ErrStr)
		})
	}
}

func TestConfigureUSDCTokenPools(t *testing.T) {
	t.Parallel()

	rt := setupUSDCTokenPoolsEnvironmentForConfigure(t, true)

	solChains := slices.Collect(maps.Values(rt.Environment().BlockChains.SolanaChains()))
	require.GreaterOrEqual(t, len(solChains), 1)

	evmChains := slices.Collect(maps.Values(rt.Environment().BlockChains.EVMChains()))
	require.GreaterOrEqual(t, len(evmChains), 1)

	chainsLen := len(solChains) + len(evmChains)

	newUSDCMsgProxies := make(map[uint64]v1_6_2.DeployCCTPMessageTransmitterProxyInput, chainsLen)
	newUSDCTokenPools := make(map[uint64]v1_6_2.DeployUSDCTokenPoolInput, chainsLen)
	newUSDCConfigs := make(map[uint64]v1_6_2.ConfigUSDCTokenPoolInput, chainsLen)
	addrBook := cldf.NewMemoryAddressBook()
	dummySolDomainID := uint32(0)
	dummyEVMDomainID := uint32(1)
	for _, evmChain := range evmChains {
		usdcToken, tokenMessenger := setupUSDCTokenPoolsContractsForConfigure(t,
			rt.Environment().Logger,
			evmChain,
			addrBook,
		)

		newUSDCMsgProxies[evmChain.Selector] = v1_6_2.DeployCCTPMessageTransmitterProxyInput{
			TokenMessenger: tokenMessenger.Address,
		}

		newUSDCTokenPools[evmChain.Selector] = v1_6_2.DeployUSDCTokenPoolInput{
			PreviousPoolAddress: v1_6_2.USDCTokenPoolSentinelAddress,
			TokenMessenger:      tokenMessenger.Address,
			TokenAddress:        usdcToken.Address,
			PoolType:            shared.USDCTokenPool,
		}

		destUpdates := map[uint64]v1_6_2.DomainUpdateInput{}
		for _, solChain := range solChains {
			minterPrivKey, err := solana.NewRandomPrivateKey()
			require.NoError(t, err)

			callerPrivKey, err := solana.NewRandomPrivateKey()
			require.NoError(t, err)

			destUpdates[solChain.Selector] = v1_6_2.DomainUpdateInput{
				MintRecipient:    minterPrivKey.PublicKey().String(),
				AllowedCaller:    callerPrivKey.PublicKey().String(),
				DomainIdentifier: dummySolDomainID,
				Enabled:          true,
			}
		}

		for _, remoteEVMChain := range evmChains {
			if remoteEVMChain.Selector == evmChain.Selector {
				continue
			}

			// Add config for EVM to EVM domain update
			destUpdates[remoteEVMChain.Selector] = v1_6_2.DomainUpdateInput{
				AllowedCaller:    utils.RandomAddress().String(),
				DomainIdentifier: dummyEVMDomainID,
				Enabled:          true,
			}
		}

		newUSDCConfigs[evmChain.Selector] = v1_6_2.ConfigUSDCTokenPoolInput{
			DestinationUpdates: destUpdates,
		}
	}

	err := rt.Exec(
		runtime.ChangesetTask(v1_6_2.DeployCCTPMessageTransmitterProxyNew, v1_6_2.DeployCCTPMessageTransmitterProxyContractConfig{
			USDCProxies: newUSDCMsgProxies,
		}),
		runtime.ChangesetTask(v1_6_2.DeployUSDCTokenPoolNew, v1_6_2.DeployUSDCTokenPoolContractsConfig{
			USDCPools: newUSDCTokenPools,
		}),
		runtime.ChangesetTask(v1_6_2.ConfigUSDCTokenPoolChangeSet, v1_6_2.ConfigUSDCTokenPoolConfig{
			USDCPools: newUSDCConfigs,
		}),
	)
	require.NoError(t, err)

	state, err := stateview.LoadOnchainState(rt.Environment())
	require.NoError(t, err)
	for _, evmChain := range evmChains {
		pools := state.Chains[evmChain.Selector].USDCTokenPoolsV1_6
		require.Len(t, pools, 1)

		for _, solChain := range solChains {
			actualDomain, err := pools[deployment.Version1_6_2].GetDomain(nil, solChain.Selector)
			require.NoError(t, err)

			expectedDomain := newUSDCConfigs[evmChain.Selector].DestinationUpdates[solChain.Selector]

			allowedCallerAddr, err := solana.PublicKeyFromBase58(expectedDomain.AllowedCaller)
			require.NoError(t, err)
			mintRecipientAddr, err := solana.PublicKeyFromBase58(expectedDomain.MintRecipient)
			require.NoError(t, err)
			require.Equal(t, allowedCallerAddr.Bytes(), actualDomain.AllowedCaller[:])
			require.Equal(t, mintRecipientAddr.Bytes(), actualDomain.MintRecipient[:])
			require.Equal(t, expectedDomain.DomainIdentifier, actualDomain.DomainIdentifier)
			require.Equal(t, expectedDomain.Enabled, actualDomain.Enabled)
		}

		for _, remoteEVMChain := range evmChains {
			if remoteEVMChain.Selector == evmChain.Selector {
				continue
			}
			actualDomain, err := pools[deployment.Version1_6_2].GetDomain(nil, remoteEVMChain.Selector)
			require.NoError(t, err)

			expectedDomain := newUSDCConfigs[evmChain.Selector].DestinationUpdates[remoteEVMChain.Selector]

			allowedCallerAddr := common.LeftPadBytes(common.HexToAddress(expectedDomain.AllowedCaller).Bytes(), 32)
			require.Equal(t, allowedCallerAddr, actualDomain.AllowedCaller[:])
			require.Equal(t, expectedDomain.DomainIdentifier, actualDomain.DomainIdentifier)
			require.Equal(t, expectedDomain.Enabled, actualDomain.Enabled)
		}
	}
}
