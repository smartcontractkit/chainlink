package solana

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	seq "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/solana/sequence/operation"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
)

const (
	CacheContract datastore.ContractType = "DataFeedsCache"
	CacheState    datastore.ContractType = "DataFeedsCacheState"
)

type DeployCacheRequest struct {
	ChainSel           uint64
	BuildConfig        *helpers.BuildSolanaConfig
	Qualifier          string
	LabelSet           datastore.LabelSet
	Version            string
	FeedAdmins         []solana.PublicKey // Feed admins to be added to the cache
	ForwarderProgramID solana.PublicKey   // ForwarderProgram that is allowed to write to this cache
}

var _ cldf.ChangeSetV2[*DeployCacheRequest] = DeployCache{}

type DeployCache struct{}

func (cs DeployCache) VerifyPreconditions(env cldf.Environment, req *DeployCacheRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return err
	}
	return nil
}

func (cs DeployCache) Apply(env cldf.Environment, req *DeployCacheRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	if req.BuildConfig != nil {
		// You may want to define a specific build params for the cache if needed
		err := helpers.BuildSolana(env, *req.BuildConfig, cacheBuildParams)
		if err != nil {
			return out, fmt.Errorf("failed build solana artifacts: %w", err)
		}
	}

	out.DataStore = datastore.NewMemoryDataStore()
	version := semver.MustParse(req.Version)
	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	deploySeqInput := seq.DeployCacheSeqInput{
		ChainSel:           req.ChainSel,
		ProgramName:        "data_feeds_cache",
		FeedAdmins:         req.FeedAdmins,
		ForwarderProgramID: req.ForwarderProgramID,
		ContractType:       CacheContract,
		Version:            version,
		Qualifier:          req.Qualifier,
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}

	deploySeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployCacheSeq, deps, deploySeqInput)
	if err != nil {
		return out, err
	}

	// Save programID
	err = out.DataStore.Addresses().Add(
		datastore.AddressRef{
			Address:       deploySeqReport.Output.ProgramID.String(),
			ChainSelector: req.ChainSel,
			Type:          CacheContract,
			Version:       version,
			Qualifier:     req.Qualifier,
			Labels:        req.LabelSet,
		},
	)
	if err != nil {
		return out, err
	}
	// Save StateID
	err = out.DataStore.Addresses().Add(
		datastore.AddressRef{
			Address:       deploySeqReport.Output.State.String(),
			ChainSelector: req.ChainSel,
			Type:          CacheState,
			Version:       version,
			Qualifier:     req.Qualifier,
			Labels:        req.LabelSet,
		},
	)
	if err != nil {
		return out, err
	}

	return out, nil
}

type SetCacheUpgradeAuthorityRequest struct {
	ChainSel            uint64
	NewUpgradeAuthority string // Use string for consistency with solana.PublicKey.String()
	Qualifier           string
	Version             string
	MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSetV2[*SetCacheUpgradeAuthorityRequest] = SetCacheUpgradeAuthority{}

type SetCacheUpgradeAuthority struct{}

func (cs SetCacheUpgradeAuthority) VerifyPreconditions(env cldf.Environment, req *SetCacheUpgradeAuthorityRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	cacheKey := datastore.NewAddressRefKey(req.ChainSel, CacheContract, version, req.Qualifier)
	_, err = env.DataStore.Addresses().Get(cacheKey)
	if err != nil {
		return fmt.Errorf("failed to load cache contract: %w", err)
	}

	if req.MCMS != nil {
		refs := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(req.ChainSel))
		_, err := helpers.FetchTimelockSigner(refs)
		if err != nil {
			return fmt.Errorf("failed fetch timelock signer: %w", err)
		}
	}

	return nil
}

func (cs SetCacheUpgradeAuthority) Apply(env cldf.Environment, req *SetCacheUpgradeAuthorityRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)
	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	cacheKey := datastore.NewAddressRefKey(req.ChainSel, CacheContract, version, req.Qualifier)
	addr, err := env.DataStore.Addresses().Get(cacheKey)
	if err != nil {
		return out, fmt.Errorf("failed to load cache contract: %w", err)
	}

	setAuthorityInput := operation.SetUpgradeAuthorityInput{
		ChainSel:            req.ChainSel,
		NewUpgradeAuthority: req.NewUpgradeAuthority,
		MCMS:                req.MCMS,
		ProgramID:           addr.Address,
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}

	execSetAuthOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.SetUpgradeAuthorityOp, deps, setAuthorityInput)
	if err != nil {
		return out, err
	}

	out.MCMSTimelockProposals = execSetAuthOut.Output.Proposals

	return out, nil
}

type InitCacheDecimalReportRequest struct {
	ChainSel  uint64
	Version   string
	Qualifier string
	MCMS      *proposalutils.TimelockConfig // if set, assumes current ownership
	DataIDs   []string
	FeedAdmin solana.PublicKey
}

var _ cldf.ChangeSetV2[*InitCacheDecimalReportRequest] = InitCacheDecimalReport{}

type InitCacheDecimalReport struct{}

func (cs InitCacheDecimalReport) VerifyPreconditions(env cldf.Environment, req *InitCacheDecimalReportRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	if _, err := semver.NewVersion(req.Version); err != nil {
		return err
	}

	_, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return err
	}

	cacheKey := datastore.NewAddressRefKey(req.ChainSel, CacheContract, semver.MustParse(req.Version), req.Qualifier)
	_, err = env.DataStore.Addresses().Get(cacheKey)
	if err != nil {
		return fmt.Errorf("failed to load cache contract: %w", err)
	}

	if len(req.DataIDs) == 0 {
		return errors.New("DataIDs cannot be empty")
	}

	if req.FeedAdmin.IsZero() {
		return errors.New("FeedAdmin cannot be zero")
	}

	return nil
}

func (cs InitCacheDecimalReport) Apply(env cldf.Environment, req *InitCacheDecimalReportRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)
	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	cacheStateRef := datastore.NewAddressRefKey(req.ChainSel, CacheState, version, req.Qualifier)
	cacheState, err := env.DataStore.Addresses().Get(cacheStateRef)
	if err != nil {
		return out, fmt.Errorf("failed load cache state for chain sel %d", req.ChainSel)
	}

	dataIDs, err := dataIDsToBytes(req.DataIDs)
	if err != nil {
		return out, err
	}

	// Create remaining accounts by deriving PDAs for each DataID
	decimalReportAccounts, err := createRemainingAccounts(env.DataStore, "decimal_report", req.ChainSel, req.Qualifier, req.Version, dataIDs)
	if err != nil {
		return out, fmt.Errorf("failed to create remaining accounts: %w", err)
	}

	initInput := operation.InitCacheDecimalReportInput{
		ChainSel:          req.ChainSel,
		MCMS:              req.MCMS,
		State:             solana.MustPublicKeyFromBase58(cacheState.Address),
		Type:              cldf.ContractType(CacheContract),
		DataIDs:           dataIDs,
		FeedAdmin:         req.FeedAdmin,
		RemainingAccounts: decimalReportAccounts,
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}

	execInitOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.InitCacheDecimalReportOp, deps, initInput)
	if err != nil {
		return out, err
	}

	out.MCMSTimelockProposals = execInitOut.Output.Proposals

	return out, nil
}

// createRemainingAccounts creates the remaining accounts needed for InitCacheDecimalFeed
// by deriving the decimal report PDAs for each DataID
func createRemainingAccounts(ds datastore.DataStore, seed string, chainSel uint64, qualifier, version string, dataIDs [][16]uint8) ([]solana.AccountMeta, error) {
	// Get the deployed cache state and program ID from the datastore
	parsedVersion := semver.MustParse(version)
	cacheStateRef := datastore.NewAddressRefKey(chainSel, CacheState, parsedVersion, qualifier)
	cacheRef := datastore.NewAddressRefKey(chainSel, CacheContract, parsedVersion, qualifier)

	cacheState, err := ds.Addresses().Get(cacheStateRef)
	if err != nil {
		return nil, fmt.Errorf("failed load cache state for chain sel %d", chainSel)
	}
	cacheProgramID, err := ds.Addresses().Get(cacheRef)
	if err != nil {
		return nil, fmt.Errorf("failed load cache program ID for chain sel %d", chainSel)
	}

	cacheStateKey := solana.MustPublicKeyFromBase58(cacheState.Address)
	cacheProgramKey := solana.MustPublicKeyFromBase58(cacheProgramID.Address)

	remainingAccounts := make([]solana.AccountMeta, len(dataIDs))
	for i, dataID := range dataIDs {
		// Derive decimal report PDA for each data ID
		seeds := [][]byte{
			[]byte(seed),
			cacheStateKey.Bytes(),
			dataID[:],
		}
		reportPDA, _, err := solana.FindProgramAddress(seeds, cacheProgramKey)
		if err != nil {
			return nil, fmt.Errorf("failed to derive decimal report PDA for data ID %x: %w", dataID, err)
		}
		// Use the Meta helper for consistency with generated bindings
		remainingAccounts[i] = *solana.Meta(reportPDA).WRITE()
	}

	return remainingAccounts, nil
}

func createPermissionFlagAccounts(programID, state solana.PublicKey, dataIDs [][16]byte, sender solana.PublicKey, owner [20]byte, name [10]byte) ([]solana.AccountMeta, error) {
	var ret []solana.AccountMeta

	for _, dataID := range dataIDs {
		repHash := createReportHash(dataID, sender, owner, name)
		flagPDA, _, err := solana.FindProgramAddress([][]byte{
			[]byte("permission_flag"),
			state.Bytes(),
			repHash[:],
		}, programID)
		if err != nil {
			return nil, fmt.Errorf("failed to derive permission_flag PDA: %w", err)
		}

		ret = append(ret, *solana.Meta(flagPDA).WRITE())
	}
	return ret, nil
}

// createReportHash matches Rust's hash::hash(&[data_id, sender, owner, name].concat()).to_bytes()
// https://github.com/smartcontractkit/chainlink-solana/blob/develop/contracts/programs/data-feeds-cache/src/lib.rs#L1089
func createReportHash(dataID [16]byte, sender solana.PublicKey, owner [20]byte, name [10]byte) [32]byte {
	buf := bytes.Join([][]byte{
		dataID[:],
		sender.Bytes(),
		owner[:],
		name[:],
	}, nil)
	return sha256.Sum256(buf)
}

func dataIDsToBytes(dataIDs []string) ([][16]byte, error) {
	var out [][16]byte
	for _, dataID := range dataIDs {
		id, err := dataIDtoBytes(dataID)
		if err != nil {
			return nil, err
		}

		out = append(out, id)
	}

	return out, nil
}

func dataIDtoBytes(dataID string) ([16]byte, error) {
	var out [16]byte
	bigID, ok := new(big.Int).SetString(dataID, 0)
	if !ok {
		return out, fmt.Errorf("invalid data_id: %v", dataID)
	}
	if bigID.BitLen() > 128 {
		return out, fmt.Errorf("data_id is too long: %d", bigID.BitLen())
	}

	copy(out[:], bigID.Bytes())
	return out, nil
}
