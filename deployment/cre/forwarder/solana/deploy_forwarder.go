package solana

import (
	"fmt"
	"maps"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	seq "github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana/sequence/operation"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

const (
	ForwarderContract         datastore.ContractType = "SolanaForwarder"
	ForwarderState            datastore.ContractType = "SolanaForwarderState"
	DefaultForwarderQualifier                        = "ks_solana_forwarder"
)

var _ cldf.ChangeSetV2[*DeployForwarderRequest] = DeployForwarder{}

type DeployForwarder struct{}

func (cs DeployForwarder) VerifyPreconditions(env cldf.Environment, req *DeployForwarderRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}
	if _, err := semver.NewVersion(req.Version); err != nil {
		return err
	}

	return nil
}

type DeployForwarderRequest = struct {
	ChainSel    uint64
	BuildConfig *helpers.BuildSolanaConfig
	Qualifier   string
	LabelSet    datastore.LabelSet
	Version     string
}

func (cs DeployForwarder) Apply(env cldf.Environment, req *DeployForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	if req.BuildConfig != nil {
		err := helpers.BuildSolana(env, *req.BuildConfig, keystoneBuildParams)
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

	deploySeqInput := seq.DeployForwarderSeqInput{
		ChainSel:     req.ChainSel,
		ProgramName:  solutils.ProgKeystoneForwarder,
		Overallocate: true,
		ContractType: ForwarderContract,
		Qualifier:    req.Qualifier,
		Version:      version,
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}

	deploySeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.DeployForwarderSeq, deps, deploySeqInput)
	if err != nil {
		return out, err
	}

	// save programID
	err = out.DataStore.Addresses().Add(
		datastore.AddressRef{
			Address:       deploySeqReport.Output.ProgramID.String(),
			ChainSelector: req.ChainSel,
			Type:          ForwarderContract,
			Version:       version,
			Qualifier:     req.Qualifier,
			Labels:        req.LabelSet,
		},
	)

	if err != nil {
		return out, err
	}
	// save StateID
	err = out.DataStore.Addresses().Add(
		datastore.AddressRef{
			Address:       deploySeqReport.Output.State.String(),
			ChainSelector: req.ChainSel,
			Type:          ForwarderState,
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

// UpgradeForwarderRequest upgrades an already deployed forwarder program in place. The program ID
// does not change, so the forwarder state and every oracles config written against it stay valid.
type UpgradeForwarderRequest struct {
	ChainSel uint64
	// BuildConfig is optional. When building locally, the forwarder is built against the deployed
	// program ID unless LocalBuild.UpgradeKeys already pins one.
	BuildConfig *helpers.BuildSolanaConfig
	Qualifier   string
	// Version identifies the deployed forwarder to upgrade.
	Version string
	// NewVersion is optional. When set, the upgraded program and its state are additionally
	// recorded in the datastore under this version.
	NewVersion string
	LabelSet   datastore.LabelSet
	MCMS       *cldfproposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSetV2[*UpgradeForwarderRequest] = UpgradeForwarder{}

type UpgradeForwarder struct{}

func (cs UpgradeForwarder) VerifyPreconditions(env cldf.Environment, req *UpgradeForwarderRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}
	if req.NewVersion != "" {
		if _, err = semver.NewVersion(req.NewVersion); err != nil {
			return fmt.Errorf("invalid new version: %w", err)
		}
	}

	forwarderKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)
	if _, err = env.DataStore.Addresses().Get(forwarderKey); err != nil {
		return fmt.Errorf("failed to load forwarder: %w", err)
	}

	if req.MCMS != nil {
		refs := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(req.ChainSel))
		if _, err = helpers.FetchTimelockSigner(refs); err != nil {
			return fmt.Errorf("failed fetch timelock signer: %w", err)
		}
	}

	return nil
}

func (cs UpgradeForwarder) Apply(env cldf.Environment, req *UpgradeForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)

	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	forwarderKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)
	forwarderRef, err := env.DataStore.Addresses().Get(forwarderKey)
	if err != nil {
		return out, fmt.Errorf("failed to load forwarder: %w", err)
	}
	programID, err := solana.PublicKeyFromBase58(forwarderRef.Address)
	if err != nil {
		return out, fmt.Errorf("failed to parse forwarder program ID %q: %w", forwarderRef.Address, err)
	}

	if req.BuildConfig != nil {
		if err = buildForUpgrade(env, *req.BuildConfig, programID); err != nil {
			return out, err
		}
	}

	deps := operation.Deps{
		Datastore: env.DataStore,
		Env:       env,
		Chain:     ch,
	}

	upgradeSeqReport, err := operations.ExecuteSequence(env.OperationsBundle, seq.UpgradeForwarderSeq, deps, seq.UpgradeForwarderSeqInput{
		ChainSel:    req.ChainSel,
		ProgramName: solutils.ProgKeystoneForwarder,
		ProgramID:   programID,
		MCMS:        req.MCMS,
	})
	if err != nil {
		return out, err
	}

	out.MCMSTimelockProposals = upgradeSeqReport.Output.Proposals

	if req.NewVersion == "" {
		return out, nil
	}

	// The addresses are unchanged by an in place upgrade, only the recorded version moves. The refs
	// of the previous version are left in place, since the datastore of a changeset output is
	// merged into the environment rather than replacing it.
	newVersion := semver.MustParse(req.NewVersion)
	stateKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderState, version, req.Qualifier)
	stateRef, err := env.DataStore.Addresses().Get(stateKey)
	if err != nil {
		return out, fmt.Errorf("failed to load forwarder state: %w", err)
	}

	out.DataStore = datastore.NewMemoryDataStore()
	for _, ref := range []datastore.AddressRef{forwarderRef, stateRef} {
		err = out.DataStore.Addresses().Upsert(datastore.AddressRef{
			Address:       ref.Address,
			ChainSelector: req.ChainSel,
			Type:          ref.Type,
			Version:       newVersion,
			Qualifier:     req.Qualifier,
			Labels:        req.LabelSet,
		})
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

// buildForUpgrade builds the forwarder artifacts against the deployed program ID.
func buildForUpgrade(env cldf.Environment, buildConfig helpers.BuildSolanaConfig, programID solana.PublicKey) error {
	if err := helpers.BuildSolana(env, withPinnedUpgradeKey(buildConfig, programID), keystoneBuildParams); err != nil {
		return fmt.Errorf("failed build solana artifacts: %w", err)
	}

	return nil
}

// withPinnedUpgradeKey pins the program ID a local build declares to the deployed one. Without this
// the upgraded binary would reject every instruction addressed to it, since anchor checks the ID
// baked in at build time. A key the caller pinned itself wins, and the request of the caller is
// left untouched.
func withPinnedUpgradeKey(buildConfig helpers.BuildSolanaConfig, programID solana.PublicKey) helpers.BuildSolanaConfig {
	if !buildConfig.LocalBuild.BuildLocally {
		return buildConfig
	}

	upgradeKeys := maps.Clone(buildConfig.LocalBuild.UpgradeKeys)
	if upgradeKeys == nil {
		upgradeKeys = make(map[cldf.ContractType]string)
	}
	if _, ok := upgradeKeys[keystoneForwarder]; !ok {
		upgradeKeys[keystoneForwarder] = programID.String()
	}
	buildConfig.LocalBuild.UpgradeKeys = upgradeKeys

	return buildConfig
}

type SetForwarderUpgradeAuthorityRequest = struct {
	ChainSel            uint64
	NewUpgradeAuthority solana.PublicKey
	Qualifier           string
	Version             string
	MCMS                *cldfproposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
}

var _ cldf.ChangeSetV2[*SetForwarderUpgradeAuthorityRequest] = SetForwarderUpgradeAuthority{}

type SetForwarderUpgradeAuthority struct{}

func (cs SetForwarderUpgradeAuthority) VerifyPreconditions(env cldf.Environment, req *SetForwarderUpgradeAuthorityRequest) error {
	if _, ok := env.BlockChains.SolanaChains()[req.ChainSel]; !ok {
		return fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	forwarderKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)
	_, err = env.DataStore.Addresses().Get(forwarderKey)
	if err != nil {
		return fmt.Errorf("failed to load forwarder: %w", err)
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

func (cs SetForwarderUpgradeAuthority) Apply(env cldf.Environment, req *SetForwarderUpgradeAuthorityRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	version := semver.MustParse(req.Version)

	ch, ok := env.BlockChains.SolanaChains()[req.ChainSel]
	if !ok {
		return out, fmt.Errorf("solana chain not found for chain selector %d", req.ChainSel)
	}

	forwarderKey := datastore.NewAddressRefKey(req.ChainSel, ForwarderContract, version, req.Qualifier)
	addr, err := env.DataStore.Addresses().Get(forwarderKey)
	if err != nil {
		return out, fmt.Errorf("failed to load forwarder: %w", err)
	}

	setAuthorityInput := operation.SetUpgradeAuthorityInput{
		ChainSel:            req.ChainSel,
		NewUpgradeAuthority: req.NewUpgradeAuthority.String(),
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

type ConfigureForwarderRequest struct {
	DON forwarder.DonConfiguration

	MCMS *cldfproposalutils.TimelockConfig // if set, assumes current ownership is the timelock

	// Chains is optional. Defines chains for which request will be executed. If empty, runs for all available chains.
	Chains    map[uint64]struct{}
	Qualifier string
	Version   string
}

var _ cldf.ChangeSetV2[*ConfigureForwarderRequest] = ConfigureForwarders{}

type ConfigureForwarders struct{}

func (cs ConfigureForwarders) VerifyPreconditions(env cldf.Environment, req *ConfigureForwarderRequest) error {
	version, err := semver.NewVersion(req.Version)
	if err != nil {
		return err
	}

	return verifyForwarderChains(env, req.Chains, version, req.Qualifier, req.MCMS)
}

func (cs ConfigureForwarders) Apply(env cldf.Environment, req *ConfigureForwarderRequest) (cldf.ChangesetOutput, error) {
	var out cldf.ChangesetOutput

	mcmsBatches, err := configureForwarders(env, req)
	if err != nil {
		return out, fmt.Errorf("failed to configure forwarder: %w", err)
	}

	if req.MCMS == nil {
		return out, nil
	}

	out.MCMSTimelockProposals, err = buildTimelockProposals(env, mcmsBatches, *req.MCMS,
		"proposal to configure keystone forwarder contract")
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return out, nil
}

func configureForwarders(env cldf.Environment, req *ConfigureForwarderRequest) (map[uint64]mcmsTypes.BatchOperation, error) {
	version := semver.MustParse(req.Version)

	cfg, err := req.DON.ForwarderConfig(chain_selectors.FamilySolana, env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get forwarder config: %w", err)
	}
	signers := toSolSigners(cfg.Signers)

	batches := make(map[uint64]mcmsTypes.BatchOperation)
	for chain := range forwarderChains(env, req.Chains) {
		target, err := resolveForwarderConfigTarget(env, chain, version, req.Qualifier, cfg.DonID, cfg.ConfigVersion, req.MCMS)
		if err != nil {
			return nil, fmt.Errorf("chain selector %d: %w", chain.Selector, err)
		}

		deps := operation.Deps{
			Datastore: env.DataStore,
			Env:       env,
			Chain:     chain,
		}

		opOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.ConfigureForwarderOp, deps, operation.ConfigureForwarderInput{
			ForwarderConfigTarget: target,
			Signers:               signers,
			F:                     cfg.F,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}

		batches[chain.Selector] = opOut.Output.Batch
	}

	return batches, nil
}

func toSolSigners(ss []common.Address) [][20]uint8 {
	ret := make([][20]uint8, 0, len(ss))
	slices.SortFunc(ss, func(a, b common.Address) int {
		return slices.Compare(a.Bytes(), b.Bytes())
	})
	for _, s := range ss {
		ret = append(ret, s)
	}

	return ret
}
