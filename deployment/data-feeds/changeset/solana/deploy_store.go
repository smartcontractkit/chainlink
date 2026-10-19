package solana

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	df_store "github.com/smartcontractkit/chainlink-solana/contracts/generated/data_feeds_store"

	seq "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence/operation"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

type DeployStoreRequest struct {
	ChainSel    uint64
	BuildConfig *helpers.BuildSolanaConfig
	Qualifier   string
	LabelSet    datastore.LabelSet
	Version     string

	FeedAdmins         []solana.PublicKey // Feed admins of the store
	ForwarderProgramID solana.PublicKey   // Forwarder program allowed to report to the store
	// Ring size of new feeds, at most 313 so a feed account fits the 10 KiB account creation limit
	DefaultMaxStoredReports uint64
	// Max seconds a report's observation timestamp may lead the clock
	MaxObservationForwardSecs uint32
}

var _ cldf.ChangeSetV2[*DeployStoreRequest] = DeployStore{}

// DeployStore deploys the data feeds store program, unless it is already in the datastore, and
// initializes a store state. The deployer is the owner.
type DeployStore struct{}

func (cs DeployStore) VerifyPreconditions(env cldf.Environment, req *DeployStoreRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return err
	}
	if req.ForwarderProgramID.IsZero() {
		return errors.New("ForwarderProgramID cannot be zero")
	}
	if len(req.FeedAdmins) > 64 {
		return fmt.Errorf("at most 64 feed admins, got %d", len(req.FeedAdmins))
	}
	if _, err := sortedFeedAdmins(req.FeedAdmins); err != nil {
		return err
	}
	if req.DefaultMaxStoredReports < minMaxStoredReports || req.DefaultMaxStoredReports > maxDefaultStoredReports {
		return fmt.Errorf("DefaultMaxStoredReports must be in [%d, %d], got %d", minMaxStoredReports, maxDefaultStoredReports, req.DefaultMaxStoredReports)
	}
	if req.MaxObservationForwardSecs == 0 {
		return errors.New("MaxObservationForwardSecs cannot be zero")
	}
	return nil
}

func (cs DeployStore) Apply(env cldf.Environment, req *DeployStoreRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	if req.BuildConfig != nil {
		if err := helpers.BuildSolana(env, *req.BuildConfig, buildParams); err != nil {
			return out, fmt.Errorf("failed build solana artifacts: %w", err)
		}
	}

	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	version := semver.MustParse(req.Version)
	feedAdmins, err := sortedFeedAdmins(req.FeedAdmins)
	if err != nil {
		return out, err
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}
	report, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployStoreSeq, deps, seq.DeployStoreSeqInput{
		ChainSel:                  req.ChainSel,
		ProgramName:               string(dataFeedsStore),
		ForwarderProgramID:        req.ForwarderProgramID,
		FeedAdmins:                feedAdmins,
		DefaultMaxStoredReports:   req.DefaultMaxStoredReports,
		MaxObservationForwardSecs: req.MaxObservationForwardSecs,
		ContractType:              StoreContract,
		Version:                   version,
		Qualifier:                 req.Qualifier,
	})
	if err != nil {
		return out, err
	}

	ds := datastore.NewMemoryDataStore()
	if report.Output.Deployed {
		if err := ds.Addresses().Add(datastore.AddressRef{
			Address:       report.Output.ProgramID.String(),
			ChainSelector: req.ChainSel,
			Type:          StoreContract,
			Version:       version,
			Qualifier:     req.Qualifier,
			Labels:        req.LabelSet,
		}); err != nil {
			return out, err
		}
	}
	if err := ds.Addresses().Add(datastore.AddressRef{
		Address:       report.Output.State.String(),
		ChainSelector: req.ChainSel,
		Type:          StoreState,
		Version:       version,
		Qualifier:     req.Qualifier,
		Labels:        req.LabelSet,
	}); err != nil {
		return out, err
	}
	out.DataStore = ds

	return out, nil
}

// sortedFeedAdmins returns the feed admins strictly increasing, as initialize requires.
func sortedFeedAdmins(admins []solana.PublicKey) ([]solana.PublicKey, error) {
	sorted := slices.Clone(admins)
	slices.SortFunc(sorted, func(a, b solana.PublicKey) int { return bytes.Compare(a[:], b[:]) })
	for i, admin := range sorted {
		if admin.IsZero() {
			return nil, errors.New("feed admin cannot be zero")
		}
		if i > 0 && admin == sorted[i-1] {
			return nil, fmt.Errorf("duplicate feed admin %s", admin)
		}
	}
	return sorted, nil
}

type SetStoreUpgradeAuthorityRequest struct {
	ChainSel            uint64
	NewUpgradeAuthority string
	Qualifier           string
	Version             string
	MCMS                *cldfproposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSetV2[*SetStoreUpgradeAuthorityRequest] = SetStoreUpgradeAuthority{}

type SetStoreUpgradeAuthority struct{}

func (cs SetStoreUpgradeAuthority) VerifyPreconditions(env cldf.Environment, req *SetStoreUpgradeAuthorityRequest) error {
	if _, err := solana.PublicKeyFromBase58(req.NewUpgradeAuthority); err != nil {
		return fmt.Errorf("invalid NewUpgradeAuthority: %w", err)
	}
	return verifyStoreRequest(env, storeRef{req.ChainSel, req.Version, req.Qualifier}, req.MCMS)
}

func (cs SetStoreUpgradeAuthority) Apply(env cldf.Environment, req *SetStoreUpgradeAuthorityRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	store, err := loadStore(env.DataStore, storeRef{req.ChainSel, req.Version, req.Qualifier})
	if err != nil {
		return out, err
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     env.BlockChains.SolanaChains()[req.ChainSel],
	}
	report, err := operations.ExecuteOperation(env.OperationsBundle, operation.SetUpgradeAuthorityOp, deps, operation.SetUpgradeAuthorityInput{
		ChainSel:            req.ChainSel,
		NewUpgradeAuthority: req.NewUpgradeAuthority,
		MCMS:                req.MCMS,
		ProgramID:           store.ProgramID.String(),
	})
	if err != nil {
		return out, err
	}
	out.MCMSTimelockProposals = report.Output.Proposals
	return out, nil
}

type TransferOwnershipStoreRequest struct {
	ChainSel  uint64
	Version   string
	Qualifier string
	MCMSCfg   cldfproposalutils.TimelockConfig
}

var _ cldf.ChangeSetV2[*TransferOwnershipStoreRequest] = TransferOwnershipStore{}

// TransferOwnershipStore proposes the timelock as the store owner, signed by the deployer as the
// current owner, and returns the timelock proposal that accepts it.
type TransferOwnershipStore struct{}

func (cs TransferOwnershipStore) VerifyPreconditions(env cldf.Environment, req *TransferOwnershipStoreRequest) error {
	return verifyStoreRequest(env, storeRef{req.ChainSel, req.Version, req.Qualifier}, &req.MCMSCfg)
}

func (cs TransferOwnershipStore) Apply(env cldf.Environment, req *TransferOwnershipStoreRequest) (cldf.ChangesetOutput, error) {
	store, err := loadStore(env.DataStore, storeRef{req.ChainSel, req.Version, req.Qualifier})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	timelockSigner, err := authority(env, req.ChainSel, &req.MCMSCfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	build := storeInstruction(store.ProgramID, nil)

	transfer, err := build(df_store.NewTransferOwnershipInstruction(timelockSigner,
		env.BlockChains.SolanaChains()[req.ChainSel].DeployerKey.PublicKey(), store.State))
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	if _, err = sendInstructions(env, req.ChainSel, nil, "transfer data feeds store ownership", []solana.Instruction{transfer}); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	accept, err := build(df_store.NewAcceptOwnershipInstruction(timelockSigner, store.State))
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return sendInstructions(env, req.ChainSel, &req.MCMSCfg, "accept data feeds store ownership", []solana.Instruction{accept})
}

type SetStoreFeedAdminsRequest struct {
	ChainSel  uint64
	Version   string
	Qualifier string
	MCMS      *cldfproposalutils.TimelockConfig // if set, assumes the timelock is the store owner

	Add    []solana.PublicKey
	Remove []solana.PublicKey
}

var _ cldf.ChangeSetV2[*SetStoreFeedAdminsRequest] = SetStoreFeedAdmins{}

// SetStoreFeedAdmins grants or revokes feed admin permission. Signed by the store owner.
type SetStoreFeedAdmins struct{}

func (cs SetStoreFeedAdmins) VerifyPreconditions(env cldf.Environment, req *SetStoreFeedAdminsRequest) error {
	if len(req.Add)+len(req.Remove) == 0 {
		return errors.New("no feed admins to add or remove")
	}
	seen := map[solana.PublicKey]bool{}
	for _, admin := range slices.Concat(req.Add, req.Remove) {
		if admin.IsZero() {
			return errors.New("feed admin cannot be zero")
		}
		if seen[admin] {
			return fmt.Errorf("feed admin %s listed twice", admin)
		}
		seen[admin] = true
	}
	return verifyStoreRequest(env, storeRef{req.ChainSel, req.Version, req.Qualifier}, req.MCMS)
}

func (cs SetStoreFeedAdmins) Apply(env cldf.Environment, req *SetStoreFeedAdminsRequest) (cldf.ChangesetOutput, error) {
	store, err := loadStore(env.DataStore, storeRef{req.ChainSel, req.Version, req.Qualifier})
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	owner, err := authority(env, req.ChainSel, req.MCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	build := storeInstruction(store.ProgramID, nil)

	var ixs []solana.Instruction
	for _, admin := range req.Add {
		ix, err := build(df_store.NewSetFeedAdminInstruction(admin, true, owner, store.State))
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		ixs = append(ixs, ix)
	}
	for _, admin := range req.Remove {
		ix, err := build(df_store.NewSetFeedAdminInstruction(admin, false, owner, store.State))
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		ixs = append(ixs, ix)
	}
	return sendInstructions(env, req.ChainSel, req.MCMS, "set data feeds store feed admins", ixs)
}
