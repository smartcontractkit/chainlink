package operation

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/mcms"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	solanaUtils "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
	cldfsol "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	df_cache "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_cache"

	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

var Version1_0_0 = semver.MustParse("1.0.0")

var (
	InitCacheOp = operations.NewOperation(
		"init-cache-op",
		Version1_0_0,
		"Initialize DataFeeds Cache for Solana Chain",
		initCache,
	)
	DeployCacheOp = operations.NewOperation(
		"deploy-cache-op",
		Version1_0_0,
		"Deploys the DataFeeds Cache program for Solana Chain",
		commonOps.Deploy,
	)
	SetUpgradeAuthorityOp = operations.NewOperation(
		"set-upgrade-authority-op",
		Version1_0_0,
		"Sets Cache's upgrade authority for Solana Chain",
		setUpgradeAuthority,
	)
	ConfigureCacheDecimalReportOp = operations.NewOperation(
		"configure-cache-decimal-report-op",
		Version1_0_0,
		"Configure cache decimal report for Solana Chain",
		configureCacheDecimalReport,
	)
	InitCacheDecimalReportOp = operations.NewOperation(
		"init-cache-decimal-feed-op",
		Version1_0_0,
		"Initialize DataFeeds Cache Decimal Report for Solana Chain",
		initCacheDecimalReport,
	)
)

type (
	Deps struct {
		Env       cldf.Environment
		Chain     cldfsol.Chain
		Datastore datastore.DataStore
	}

	// For DataFeeds Cache initialization
	InitCacheInput struct {
		ProgramID  solana.PublicKey
		ChainSel   uint64
		FeedAdmins []solana.PublicKey // Feed admins to be added to the cache
	}

	InitCacheOutput struct {
		StatePubKey solana.PublicKey
	}

	SetUpgradeAuthorityInput struct {
		ChainSel            uint64
		ProgramID           string
		NewUpgradeAuthority string
		MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
	}

	SetUpgradeAuthorityOutput struct {
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}

	ConfigureCacheDecimalReportInput struct {
		ChainSel          uint64
		Descriptions      [][32]uint8
		DataIDs           [][16]uint8
		MCMS              *proposalutils.TimelockConfig // if set, assumes current owner is the timelock
		WorkflowMetadatas []df_cache.WorkflowMetadata
		FeedAdmin         solana.PublicKey
		State             solana.PublicKey
		Type              cldf.ContractType
		RemainingAccounts []solana.AccountMeta
	}

	ConfigureCacheOutput struct {
		Proposals []mcms.TimelockProposal // will be returned in case if timelock config is passed
	}

	InitCacheDecimalReportInput struct {
		ChainSel          uint64
		Version           string
		Qualifier         string
		MCMS              *proposalutils.TimelockConfig // if set, assumes current
		DataIDs           [][16]uint8
		FeedAdmin         solana.PublicKey
		State             solana.PublicKey
		Type              cldf.ContractType
		RemainingAccounts []solana.AccountMeta
	}
)

func confirmInstructionOrBuildProposal(
	deps Deps,
	chainSel uint64,
	instruction solana.Instruction,
	mcmsConfig *proposalutils.TimelockConfig,
	proposalDescription string,
) ([]mcms.TimelockProposal, error) {
	if mcmsConfig == nil {
		if err := deps.Chain.Confirm([]solana.Instruction{instruction}); err != nil {
			return nil, fmt.Errorf("failed to confirm instructions: %w", err)
		}
		return nil, nil
	}

	return buildMCMSProposal(deps, chainSel, instruction, mcmsConfig, proposalDescription)
}

func buildMCMSProposal(
	deps Deps,
	chainSel uint64,
	instruction solana.Instruction,
	mcmsConfig *proposalutils.TimelockConfig,
	description string,
) ([]mcms.TimelockProposal, error) {
	tx, err := helpers.BuildMCMSTxn(
		instruction,
		solana.BPFLoaderUpgradeableProgramID.String(),
		cldf.ContractType(solana.BPFLoaderUpgradeableProgramID.String()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	proposal, err := helpers.BuildProposalsForTxns(
		deps.Env,
		chainSel,
		description,
		mcmsConfig.MinDelay,
		[]mcmsTypes.Transaction{*tx},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}

	return []mcms.TimelockProposal{*proposal}, nil
}

func getCurrentAuthority(deps Deps, chainSel uint64, mcmsConfig *proposalutils.TimelockConfig) (solana.PublicKey, error) {
	if mcmsConfig == nil {
		return deps.Chain.DeployerKey.PublicKey(), nil
	}

	timelockSignerPDA, err := helpers.FetchTimelockSigner(
		deps.Datastore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel)),
	)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to get timelock signer: %w", err)
	}
	return timelockSignerPDA, nil
}

func initCache(b operations.Bundle, deps Deps, in InitCacheInput) (InitCacheOutput, error) {
	var out InitCacheOutput

	if df_cache.ProgramID.IsZero() {
		df_cache.SetProgramID(in.ProgramID)
	}

	stateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		return out, fmt.Errorf("failed to create random keys: %w", err)
	}

	instruction, err := df_cache.NewInitializeInstruction(
		in.FeedAdmins,
		deps.Chain.DeployerKey.PublicKey(),
		stateKey.PublicKey(),
		solana.SystemProgramID,
	).ValidateAndBuild()
	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	if err = deps.Chain.Confirm([]solana.Instruction{instruction}, solanaUtils.AddSigners(stateKey)); err != nil {
		return out, errors.New("failed to confirm")
	}

	out.StatePubKey = stateKey.PublicKey()
	return out, nil
}

func setUpgradeAuthority(b operations.Bundle, deps Deps, in SetUpgradeAuthorityInput) (SetUpgradeAuthorityOutput, error) {
	var out SetUpgradeAuthorityOutput

	programID, err := solana.PublicKeyFromBase58(in.ProgramID)
	if err != nil {
		return out, fmt.Errorf("failed parse programID: %w", err)
	}

	newAuthority, err := solana.PublicKeyFromBase58(in.NewUpgradeAuthority)
	if err != nil {
		return out, fmt.Errorf("failed parse upgrade authority: %w", err)
	}

	currentAuthority, err := getCurrentAuthority(deps, in.ChainSel, in.MCMS)
	if err != nil {
		return out, err
	}

	instruction := helpers.SetUpgradeAuthority(&deps.Env, programID, currentAuthority, newAuthority, false)

	proposals, err := confirmInstructionOrBuildProposal(
		deps,
		in.ChainSel,
		instruction,
		in.MCMS,
		"proposal to SetUpgradeAuthority in Solana",
	)
	if err != nil {
		return out, err
	}

	if proposals != nil {
		out.Proposals = proposals
	}

	return out, nil
}

func initCacheDecimalReport(b operations.Bundle, deps Deps, in InitCacheDecimalReportInput) (ConfigureCacheOutput, error) {
	var out ConfigureCacheOutput

	instruction := df_cache.NewInitDecimalReportsInstruction(
		in.DataIDs,
		in.FeedAdmin,
		in.State,
		solana.SystemProgramID,
	)

	for _, acc := range in.RemainingAccounts {
		instruction.AccountMetaSlice = append(instruction.AccountMetaSlice, &acc)
	}

	tx, err := instruction.ValidateAndBuild()

	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	proposals, err := confirmInstructionOrBuildProposal(
		deps,
		in.ChainSel,
		tx,
		in.MCMS,
		"proposal to InitDecimalReports in Solana",
	)
	if err != nil {
		return out, err
	}

	if proposals != nil {
		out.Proposals = proposals
	}

	return out, nil
}

func configureCacheDecimalReport(b operations.Bundle, deps Deps, in ConfigureCacheDecimalReportInput) (ConfigureCacheOutput, error) {
	var out ConfigureCacheOutput

	instruction := df_cache.NewSetDecimalFeedConfigsInstruction(
		in.DataIDs,
		in.Descriptions,
		in.WorkflowMetadatas,
		in.FeedAdmin,
		in.State,
		solana.SystemProgramID,
	)

	for _, acc := range in.RemainingAccounts {
		instruction.AccountMetaSlice = append(instruction.AccountMetaSlice, &acc)
	}

	tx, err := instruction.ValidateAndBuild()

	if err != nil {
		return out, fmt.Errorf("failed to build and validate initialize instruction %w", err)
	}

	proposals, err := confirmInstructionOrBuildProposal(
		deps,
		in.ChainSel,
		tx,
		in.MCMS,
		"proposal to SetDecimalFeedConfigs in Solana",
	)
	if err != nil {
		return out, err
	}

	if proposals != nil {
		out.Proposals = proposals
	}

	return out, nil
}
