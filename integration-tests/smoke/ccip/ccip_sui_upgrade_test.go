package ccip

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink-sui/contracts"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suiutil "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_cs "github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers/messagingtest"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	testsetups "github.com/smartcontractkit/chainlink/integration-tests/testsetups/ccip"
)

func Test_CCIP_Upgrade_Sui2EVM(t *testing.T) {
	ctx := testcontext.Get(t)

	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	sourceChain := suiChainSelectors[0]
	destChain := evmChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (SUI): ", destChain)

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	suiSenderAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	normalizedAddr, err := suiutil.ConvertStringToAddressBytes(suiSenderAddr)
	require.NoError(t, err)

	// SUI FeeToken
	// mint link token to use as feeToken
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.MintLinkToken{}, sui_cs.MintLinkTokenConfig{
			ChainSelector:  sourceChain,
			TokenPackageId: state.SuiChains[sourceChain].LinkTokenAddress,
			TreasuryCapId:  state.SuiChains[sourceChain].LinkTokenTreasuryCapId,
			Amount:         1000000000000, // 1000 Link with 1e9
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]
	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[linkops.MintLinkTokenOutput])
	require.True(t, ok)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(normalizedAddr[:], 32)
		out    messagingtest.TestCaseOutput
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // testRouter
		)
	)

	// upgrade contracts, upgrade onRamp to v2
	t.Log("Upgrading SUI contracts")
	upgradeCCIP(ctx, t, e, sourceChain, contracts.CCIP)
	upgradeSuiOnRamp(ctx, t, e, sourceChain, contracts.CCIPOnramp)

	t.Run("Sui OnRamp, CCIP FQ Upgraded: Message to EVM - Should Succeed", func(t *testing.T) {
		out = messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               state.Chains[destChain].Receiver.Address().Bytes(),
				ExtraArgs:              nil,
				Replayed:               true,
				FeeToken:               outputMap.Objects.MintedLinkTokenObjectId,
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})

	t.Logf("out: %v\n", out)
}

func Test_CCIP_Upgrade_EVM2Sui(t *testing.T) {
	ctx := testcontext.Get(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// Deploy SUI Receiver
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]

	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(outputMap.PackageId, "0x")
	receiverByteDecoded, err := hex.DecodeString(id)
	require.NoError(t, err)

	// register the receiver
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)

	receiverByte := receiverByteDecoded

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	t.Log("Upgrading SUI contracts")
	upgradeCCIP(ctx, t, e, destChain, contracts.CCIP)
	upgradeSuiOffRamp(ctx, t, e, destChain, contracts.CCIPOfframp)

	// Block offramp v1
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.BlockVersion{}, sui_cs.BlockVersionConfig{
			SuiChainSelector: destChain,
			CCIPPackageId:    state.SuiChains[destChain].CCIPAddress,
			StateObjectId:    state.SuiChains[destChain].CCIPObjectRef,
			OwnerCapObjectId: state.SuiChains[destChain].CCIPOwnerCapObjectId,
			ModuleName:       "offramp",
			Version:          1,
		}),
	})
	require.NoError(t, err)

	// Block ccip v1 feequoter
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.BlockVersion{}, sui_cs.BlockVersionConfig{
			SuiChainSelector: destChain,
			CCIPPackageId:    state.SuiChains[destChain].CCIPAddress,
			StateObjectId:    state.SuiChains[destChain].CCIPObjectRef,
			OwnerCapObjectId: state.SuiChains[destChain].CCIPOwnerCapObjectId,
			ModuleName:       "fee_quoter",
			Version:          1,
		}),
	})
	require.NoError(t, err)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // test router
		)
	)

	t.Run("OffRamp, CCIP FQ upgraded and blocked v1: Message to Sui - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               receiverByte,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}

func Test_CCIP_Upgrade_NoBlock_EVM2Sui(t *testing.T) {
	ctx := testcontext.Get(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	// Deploy SUI Receiver
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]

	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(outputMap.PackageId, "0x")
	receiverByteDecoded, err := hex.DecodeString(id)
	require.NoError(t, err)

	// register the receiver
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)

	receiverByte := receiverByteDecoded

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	t.Log("Upgrading SUI contracts")
	upgradeCCIP(ctx, t, e, destChain, contracts.CCIP)
	upgradeSuiOffRamp(ctx, t, e, destChain, contracts.CCIPOfframp)

	state, err = stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // test router
		)
	)

	t.Run("OffRamp, CCIP FQ upgraded NoBlock: Message to Sui - Should Succeed", func(t *testing.T) {
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               receiverByte,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}

func Test_CCIP_Upgrade_CommonPkg_EVM2Sui(t *testing.T) {
	ctx := testcontext.Get(t)
	e, _, _ := testsetups.NewIntegrationEnvironment(
		t,
		testhelpers.WithNumOfChains(2),
		testhelpers.WithSuiChains(1),
	)

	evmChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyEVM))
	suiChainSelectors := e.Env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilySui))

	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	sourceChain := evmChainSelectors[0]
	destChain := suiChainSelectors[0]

	t.Log("Source chain (EVM): ", sourceChain, "Dest chain (Sui): ", destChain)

	err = testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &e, state, sourceChain, destChain, false)
	require.NoError(t, err)

	var (
		nonce  uint64
		sender = common.LeftPadBytes(e.Env.BlockChains.EVMChains()[sourceChain].DeployerKey.From.Bytes(), 32)
		setup  = messagingtest.NewTestSetupWithDeployedEnv(
			t,
			e,
			state,
			sourceChain,
			destChain,
			sender,
			false, // test router
		)
	)

	// Deploy SUI Receiver
	_, output, err := commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.DeployDummyReceiver{}, sui_cs.DeployDummyReceiverConfig{
			SuiChainSelector: destChain,
			McmsOwner:        "0x1",
		}),
	})
	require.NoError(t, err)

	rawOutput := output[0].Reports[0]

	outputMap, ok := rawOutput.Output.(sui_ops.OpTxResult[ccipops.DeployDummyReceiverObjects])
	require.True(t, ok)

	id := strings.TrimPrefix(outputMap.PackageId, "0x")
	receiverByteDecoded, err := hex.DecodeString(id)
	require.NoError(t, err)

	// register the receiver
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.RegisterDummyReceiver{}, sui_cs.RegisterDummyReceiverConfig{
			SuiChainSelector:       destChain,
			OwnerCapObjectId:       outputMap.Objects.OwnerCapObjectId,
			CCIPObjectRefObjectId:  state.SuiChains[destChain].CCIPObjectRef,
			DummyReceiverPackageId: outputMap.PackageId,
		}),
	})
	require.NoError(t, err)

	receiverByte := receiverByteDecoded

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		outputMap.Objects.CCIPReceiverStateObjectId,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	t.Log("Upgrading SUI contracts")
	upgradeCCIP(ctx, t, e, destChain, contracts.CCIP)

	// Block ccip v1 FQ
	_, _, err = commoncs.ApplyChangesets(t, e.Env, []commoncs.ConfiguredChangeSet{
		commoncs.Configure(sui_cs.BlockVersion{}, sui_cs.BlockVersionConfig{
			SuiChainSelector: destChain,
			CCIPPackageId:    state.SuiChains[destChain].CCIPAddress,
			StateObjectId:    state.SuiChains[destChain].CCIPObjectRef,
			OwnerCapObjectId: state.SuiChains[destChain].CCIPOwnerCapObjectId,
			ModuleName:       "fee_quoter",
			Version:          1,
		}),
	})
	require.NoError(t, err)

	t.Run("CCIP FQ upgraded blocked v1: Message to Sui - Should Succeed", func(t *testing.T) {
		// ccipChainState := state.SuiChains[destChain]
		message := []byte("Hello Sui, from EVM!")
		messagingtest.Run(t,
			messagingtest.TestCase{
				TestSetup:              setup,
				Nonce:                  &nonce,
				ValidationType:         messagingtest.ValidationTypeExec,
				Receiver:               receiverByte,
				MsgData:                message,
				ExtraArgs:              testhelpers.MakeSuiExtraArgs(1000000, true, receiverObjectIDs, [32]byte{}),
				ExpectedExecutionState: testhelpers.EXECUTION_STATE_SUCCESS,
			},
		)
	})
}

func upgradeSuiOnRamp(ctx context.Context, t *testing.T, e testhelpers.DeployedEnv, sourceChain uint64, version contracts.Package) {
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	signerAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	// compile packages
	compiledPackage, err := suiBind.CompilePackage(version, map[string]string{
		"ccip":        state.SuiChains[sourceChain].CCIPAddress,
		"ccip_onramp": "0x0", // old onRamp address
		"mcms":        state.SuiChains[sourceChain].MCMSPackageID,
		"mcms_owner":  "0x1",

		"latest_ccip_pkg":     state.SuiChains[sourceChain].CCIPMockV2PackageId,
		"original_onramp_pkg": state.SuiChains[sourceChain].OnRampAddress,
		"upgrade_cap":         state.SuiChains[sourceChain].OnRampUpgradeCapId,
		"signer":              signerAddr,
	}, true, "")
	require.NoError(t, err)

	// decode modules from base64 -> [][]byte
	moduleBytes := make([][]byte, len(compiledPackage.Modules))
	for i, moduleBase64 := range compiledPackage.Modules {
		decoded, err := base64.StdEncoding.DecodeString(moduleBase64)
		require.NoError(t, err)

		moduleBytes[i] = decoded
	}

	depAddresses := make([]models.SuiAddress, len(compiledPackage.Dependencies))
	for i, dep := range compiledPackage.Dependencies {
		depAddresses[i] = models.SuiAddress(dep)
	}

	policy := byte(0)

	// upgrade the onRamp
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},

		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		state.SuiChains[sourceChain].OnRampAddress,
		state.SuiChains[sourceChain].OnRampUpgradeCapId,
		moduleBytes,
		depAddresses,
		policy,
		compiledPackage.Digest,
	)
	require.NoError(t, err)

	newOnRampPkgID, err := suiBind.FindPackageIdFromPublishTx(*resp)
	require.NoError(t, err)

	// Add new PackageID
	onRamp, err := module_onramp.NewOnramp(state.SuiChains[sourceChain].OnRampAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	// add new pkgId to state
	_, err = onRamp.AddPackageId(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].OnRampStateObjectId}, suiBind.Object{Id: state.SuiChains[sourceChain].OnRampOwnerCapObjectId}, newOnRampPkgID)
	require.NoError(t, err)

	newOnRamp, err := module_onramp.NewOnramp(newOnRampPkgID, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	typeAndVersion, err := newOnRamp.DevInspect().TypeAndVersion(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	})
	require.NoError(t, err)

	require.Equal(t, "OnRamp 1.6.1", typeAndVersion)

	// save the new pkgId to addressbook
	typeAndVersionOnRampMockV2 := cldf.NewTypeAndVersion(deployment.SuiOnRampMockV2, deployment.Version1_0_0)
	//nolint:staticcheck // using ExistingAddresses temporarily until Sui migration to datastore is complete
	err = e.Env.ExistingAddresses.Save(sourceChain, newOnRampPkgID, typeAndVersionOnRampMockV2)
	require.NoError(t, err)

	t.Log("Upgraded SUI onRamp")
}

func upgradeSuiOffRamp(ctx context.Context, t *testing.T, e testhelpers.DeployedEnv, sourceChain uint64, version contracts.Package) {
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	signerAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	// compile packages
	compiledPackage, err := suiBind.CompilePackage(version, map[string]string{
		"ccip":         state.SuiChains[sourceChain].CCIPAddress,
		"ccip_offramp": "0x0",
		"mcms":         state.SuiChains[sourceChain].MCMSPackageID,
		"mcms_owner":   "0x1",

		"latest_ccip_pkg":      state.SuiChains[sourceChain].CCIPMockV2PackageId,
		"original_offramp_pkg": state.SuiChains[sourceChain].OffRampAddress,
		"upgrade_cap":          state.SuiChains[sourceChain].OffRampUpgradeCapId,
		"signer":               signerAddr,
	}, true, "")
	require.NoError(t, err)

	// decode modules from base64 -> [][]byte
	moduleBytes := make([][]byte, len(compiledPackage.Modules))
	for i, moduleBase64 := range compiledPackage.Modules {
		decoded, err := base64.StdEncoding.DecodeString(moduleBase64)
		require.NoError(t, err)

		moduleBytes[i] = decoded
	}

	depAddresses := make([]models.SuiAddress, len(compiledPackage.Dependencies))
	for i, dep := range compiledPackage.Dependencies {
		depAddresses[i] = models.SuiAddress(dep)
	}

	policy := byte(0)

	// upgrade the offramp
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},

		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		state.SuiChains[sourceChain].OffRampAddress,
		state.SuiChains[sourceChain].OffRampUpgradeCapId,
		moduleBytes,
		depAddresses,
		policy,
		compiledPackage.Digest,
	)
	require.NoError(t, err)

	newOffRampPkgID, err := suiBind.FindPackageIdFromPublishTx(*resp)
	require.NoError(t, err)

	// Add new PackageID
	offRamp, err := module_offramp.NewOfframp(state.SuiChains[sourceChain].OffRampAddress, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	// add new pkgId to state
	_, err = offRamp.AddPackageId(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].OffRampStateObjectId}, suiBind.Object{Id: state.SuiChains[sourceChain].OffRampOwnerCapId}, newOffRampPkgID)
	require.NoError(t, err)

	newOffRamp, err := module_offramp.NewOfframp(newOffRampPkgID, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	typeAndVersion, err := newOffRamp.DevInspect().TypeAndVersion(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	})
	require.NoError(t, err)

	require.Equal(t, "OffRamp 1.6.1", typeAndVersion)

	// save the new pkgId to addressbook
	typeAndVersionOffRampMockV2 := cldf.NewTypeAndVersion(deployment.SuiOffRampMockV2, deployment.Version1_0_0)
	//nolint:staticcheck // using ExistingAddresses temporarily until Sui migration to datastore is complete
	err = e.Env.ExistingAddresses.Save(sourceChain, newOffRampPkgID, typeAndVersionOffRampMockV2)
	require.NoError(t, err)

	t.Log("Upgraded SUI offRamp")
}

func upgradeCCIP(ctx context.Context, t *testing.T, e testhelpers.DeployedEnv, sourceChain uint64, version contracts.Package) string {
	state, err := stateview.LoadOnchainState(e.Env)
	require.NoError(t, err)

	signerAddr, err := e.Env.BlockChains.SuiChains()[sourceChain].Signer.GetAddress()
	require.NoError(t, err)

	t.Log("UPGRADECAP, SIGNER: ", state.SuiChains[sourceChain].CCIPUpgradeCapObjectId, signerAddr)
	// compile packages
	compiledPackage, err := suiBind.CompilePackage(version, map[string]string{
		"ccip":       "0x0",
		"mcms":       state.SuiChains[sourceChain].MCMSPackageID,
		"mcms_owner": signerAddr,

		"original_ccip_pkg": state.SuiChains[sourceChain].CCIPAddress,
		"upgrade_cap":       state.SuiChains[sourceChain].CCIPUpgradeCapObjectId,
		"signer":            signerAddr,
	}, true, "")
	require.NoError(t, err)

	// decode modules from base64 -> [][]byte
	moduleBytes := make([][]byte, len(compiledPackage.Modules))
	for i, moduleBase64 := range compiledPackage.Modules {
		decoded, err := base64.StdEncoding.DecodeString(moduleBase64)
		require.NoError(t, err)

		moduleBytes[i] = decoded
	}

	depAddresses := make([]models.SuiAddress, len(compiledPackage.Dependencies))
	for i, dep := range compiledPackage.Dependencies {
		depAddresses[i] = models.SuiAddress(dep)
	}

	policy := byte(0)

	// upgrade the ccipPkg
	b := uint64(700_000_000)
	resp, err := testhelpers.UpgradeContractDirect(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	},

		e.Env.BlockChains.SuiChains()[sourceChain].Client,
		state.SuiChains[sourceChain].CCIPAddress,
		state.SuiChains[sourceChain].CCIPUpgradeCapObjectId,
		moduleBytes,
		depAddresses,
		policy,
		compiledPackage.Digest,
	)
	require.NoError(t, err)

	newCCIPPkgID, err := suiBind.FindPackageIdFromPublishTx(*resp)
	require.NoError(t, err)

	// Add new PackageID
	ccipStateObject, err := module_state_object.NewStateObject(newCCIPPkgID, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	// add new pkgId to state
	_, err = ccipStateObject.AddPackageId(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	}, suiBind.Object{Id: state.SuiChains[sourceChain].CCIPObjectRef}, suiBind.Object{Id: state.SuiChains[sourceChain].CCIPOwnerCapObjectId}, newCCIPPkgID)
	require.NoError(t, err)

	newFQ, err := module_fee_quoter.NewFeeQuoter(newCCIPPkgID, e.Env.BlockChains.SuiChains()[sourceChain].Client)
	require.NoError(t, err)

	typeAndVersion, err := newFQ.DevInspect().TypeAndVersion(ctx, &suiBind.CallOpts{
		Signer:           e.Env.BlockChains.SuiChains()[sourceChain].Signer,
		WaitForExecution: true,
		GasBudget:        &b,
	})
	require.NoError(t, err)

	require.Equal(t, "FeeQuoter 1.6.1", typeAndVersion)

	// save the new pkgId to addressbook
	typeAndVersionCCIPMockV2 := cldf.NewTypeAndVersion(deployment.SuiCCIPMockV2, deployment.Version1_0_0)
	//nolint:staticcheck // using ExistingAddresses temporarily until Sui migration to datastore is complete
	err = e.Env.ExistingAddresses.Save(sourceChain, newCCIPPkgID, typeAndVersionCCIPMockV2)
	require.NoError(t, err)

	t.Log("Upgraded SUI CCIP")

	return newCCIPPkgID
}
