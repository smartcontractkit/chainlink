package solana

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence/operation"

	df_cache "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_cache"
)

type Sender struct {
	ProgramID solana.PublicKey
	StateID   solana.PublicKey
}

type ConfigureCacheDecimalReportRequest struct {
	MCMS *proposalutils.TimelockConfig // if set, assumes current ownership is the timelock

	ChainSel  uint64
	Qualifier string
	Version   string

	SenderList []Sender

	AllowedWorkflowOwner [][20]uint8
	AllowedWorkflowName  [][10]uint8
	FeedAdmin            solana.PublicKey

	Descriptions [][32]uint8
	DataIDs      []string
}

var _ cldf.ChangeSetV2[*ConfigureCacheDecimalReportRequest] = ConfigureCacheDecimalReport{}

type ConfigureCacheDecimalReport struct{}

func (cs ConfigureCacheDecimalReport) VerifyPreconditions(env cldf.Environment, req *ConfigureCacheDecimalReportRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return err
	}
	// Check that AllowedSender, AllowedWorkflowOwner, and AllowedWorkflowName are all the same length
	// This is a requirement for the ConfigureCacheDecimalFeed operation
	if len(req.SenderList) != len(req.AllowedWorkflowOwner) || len(req.SenderList) != len(req.AllowedWorkflowName) {
		return errors.New("SenderList, AllowedWorkflowOwner, and AllowedWorkflowName must all have the same length")
	}

	// Check that Descriptions and DataIDs are all the same length
	if len(req.DataIDs) != len(req.Descriptions) {
		return errors.New("descriptions and DataIDs must all have the same length")
	}

	_, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return err
	}

	return nil
}

func (cs ConfigureCacheDecimalReport) Apply(env cldf.Environment, req *ConfigureCacheDecimalReportRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)

	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	cacheStateRef := datastore.NewAddressRefKey(req.ChainSel, CacheState, version, req.Qualifier)
	cacheRef := datastore.NewAddressRefKey(req.ChainSel, CacheContract, version, req.Qualifier)

	cacheState, err := env.DataStore.Addresses().Get(cacheStateRef)
	if err != nil {
		return out, fmt.Errorf("failed load cache state for chain sel %d", req.ChainSel)
	}

	cacheProgramID, err := env.DataStore.Addresses().Get(cacheRef)
	if err != nil {
		return out, fmt.Errorf("failed load cache program ID for chain sel %d", req.ChainSel)
	}

	var remainingAccounts []solana.AccountMeta
	dataIDs, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}

	// Create decimalReportAccounts by deriving PDAs for each DataID
	decimalReportAccounts, err := createRemainingAccounts(env.DataStore, "feed_config", req.ChainSel, req.Qualifier, req.Version, dataIDs)
	if err != nil {
		return out, fmt.Errorf("failed to create remaining accounts: %w", err)
	}

	cacheStateKey := solana.MustPublicKeyFromBase58(cacheState.Address)
	cacheProgramKey := solana.MustPublicKeyFromBase58(cacheProgramID.Address)
	remainingAccounts = append(remainingAccounts, decimalReportAccounts...)

	workflowMetadatas := make([]df_cache.WorkflowMetadata, len(req.SenderList))

	for idx, sender := range req.SenderList {
		allowedSender, err := deriveForwarderAuthority(sender.StateID, cacheProgramKey, sender.ProgramID)
		if err != nil {
			return out, fmt.Errorf("failed to derive forwarder authority: %w", err)
		}

		permissionAccounts, err := createPermissionFlagAccounts(cacheProgramKey, cacheStateKey, dataIDs,
			allowedSender, req.AllowedWorkflowOwner[idx], req.AllowedWorkflowName[idx])

		if err != nil {
			return out, fmt.Errorf("failed to create permission accounts: %w", err)
		}
		remainingAccounts = append(remainingAccounts, permissionAccounts...)

		workflowMetadatas[idx] = df_cache.WorkflowMetadata{
			AllowedSender:        allowedSender,
			AllowedWorkflowOwner: req.AllowedWorkflowOwner[idx],
			AllowedWorkflowName:  req.AllowedWorkflowName[idx],
		}
	}

	configureCacheDecimalReportInput := operation.ConfigureCacheDecimalReportInput{
		ChainSel:          req.ChainSel,
		MCMS:              req.MCMS,
		State:             solana.MustPublicKeyFromBase58(cacheState.Address),
		Type:              cldf.ContractType(CacheContract),
		WorkflowMetadatas: workflowMetadatas,
		FeedAdmin:         req.FeedAdmin,
		DataIDs:           dataIDs,
		Descriptions:      req.Descriptions,
		RemainingAccounts: remainingAccounts,
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}

	execSetAuthOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.ConfigureCacheDecimalReportOp, deps, configureCacheDecimalReportInput)
	if err != nil {
		return out, err
	}

	out.MCMSTimelockProposals = execSetAuthOut.Output.Proposals

	return out, nil
}

func deriveForwarderAuthority(forwarderState solana.PublicKey, receiverProgram solana.PublicKey, forwarderProgram solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte("forwarder"),
		forwarderState[:],
		receiverProgram[:],
	}
	ret, _, err := solana.FindProgramAddress(seeds, forwarderProgram)
	return ret, err
}
