package solana

import (
	"encoding/binary"
	"errors"
	"fmt"
	"slices"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	mcmsSolana "github.com/smartcontractkit/mcms/sdk/solana"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	commonstate "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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

type SetForwarderUpgradeAuthorityRequest = struct {
	ChainSel            uint64
	NewUpgradeAuthority solana.PublicKey
	Qualifier           string
	Version             string
	MCMS                *proposalutils.TimelockConfig // if set, assumes current upgrade authority is the timelock
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

	MCMS *proposalutils.TimelockConfig // if set, assumes current ownership is the timelock

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

	if req.Chains != nil {
		for sel := range req.Chains {
			if _, ok := env.BlockChains.SolanaChains()[sel]; !ok {
				return fmt.Errorf("solana chain not found for chain selector %d", sel)
			}
			forwarderKey := datastore.NewAddressRefKey(sel, ForwarderContract, version, req.Qualifier)
			_, err := env.DataStore.Addresses().Get(forwarderKey)

			if err != nil {
				return fmt.Errorf("failed get fowarder for chain selector %d: %w", sel, err)
			}
			if req.MCMS != nil {
				_, err = commonstate.MaybeLoadMCMSWithTimelockChainStateSolanaV2(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(sel)))
				if err != nil {
					return fmt.Errorf("failed to load MCMS for chain selector %d: %w", sel, err)
				}
			}
		}
	}

	return nil
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
	env.Logger.Info("req delay", req.MCMS.MinDelay)

	var proposals []mcms.TimelockProposal
	for chainSel, batch := range mcmsBatches {
		// get timelocks, proposers, inspectors per chain
		solChain := env.BlockChains.SolanaChains()[chainSel]

		addresses := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSel))
		mcmState, _ := commonstate.MaybeLoadMCMSWithTimelockChainStateSolanaV2(addresses)
		if mcmState.TimelockProgram.IsZero() {
			return cldf.ChangesetOutput{}, errors.New("timelock is not found")
		}

		timelocks := map[uint64]string{}
		proposers := map[uint64]string{}
		inspectors := map[uint64]sdk.Inspector{}
		timelocks[solChain.Selector] = mcmsSolana.ContractAddress(
			mcmState.TimelockProgram,
			mcmsSolana.PDASeed(mcmState.TimelockSeed),
		)

		proposers[solChain.Selector] = mcmsSolana.ContractAddress(mcmState.McmProgram, mcmsSolana.PDASeed(mcmState.ProposerMcmSeed))
		inspectors[solChain.Selector] = mcmsSolana.NewInspector(solChain.Client)
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocks,
			proposers,
			inspectors,
			[]mcmsTypes.BatchOperation{batch},
			"proposal to transfer ownership of keystone forwarder contract to timelock",
			*req.MCMS)

		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		proposals = append(proposals, *proposal)
	}
	out.MCMSTimelockProposals = proposals

	return out, nil
}

func configureForwarders(env cldf.Environment, req *ConfigureForwarderRequest) (map[uint64]mcmsTypes.BatchOperation, error) {
	ops := make(map[uint64]mcmsTypes.BatchOperation)
	version := semver.MustParse(req.Version)

	cfg, err := req.DON.ForwarderConfig(chain_selectors.FamilySolana, env.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get forwarder config: %w", err)
	}

	for _, chain := range env.BlockChains.SolanaChains() {
		if _, shouldInclude := req.Chains[chain.Selector]; len(req.Chains) > 0 && !shouldInclude {
			continue
		}
		forwarderStateRef := datastore.NewAddressRefKey(chain.Selector, ForwarderState, version, req.Qualifier)
		forwarderRef := datastore.NewAddressRefKey(chain.Selector, ForwarderContract, version, req.Qualifier)
		forwarderState, err := env.DataStore.Addresses().Get(forwarderStateRef)
		if err != nil {
			return nil, fmt.Errorf("failed load forwarder state for chain sel %d", chain.Selector)
		}
		forwarderProgramID, err := env.DataStore.Addresses().Get(forwarderRef)
		if err != nil {
			return nil, fmt.Errorf("failed load forwarder for chain sel %d", chain.Selector)
		}
		configPDA := getConfigPDA(solana.MustPublicKeyFromBase58(forwarderState.Address),
			cfg.DonID, cfg.ConfigVersion, solana.MustPublicKeyFromBase58(forwarderProgramID.Address))

		owner := chain.DeployerKey.PublicKey()
		if req.MCMS != nil {
			// get timelock from datastore
			timelockPDA, err := helpers.FetchTimelockSigner(env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chain.Selector)))
			if err != nil {
				return nil, err
			}
			owner = timelockPDA
		}

		deps := operation.Deps{
			Datastore: env.DataStore,
			Env:       env,
			Chain:     chain,
		}
		signers := toSolSigners(cfg.Signers)
		opOut, err := operations.ExecuteOperation(env.OperationsBundle, operation.ConfigureForwarderOp, deps, operation.ConfigureForwarderInput{
			ProgramID:      solana.MustPublicKeyFromBase58(forwarderProgramID.Address),
			MCMS:           req.MCMS,
			Owner:          owner.String(),
			Signers:        signers,
			DonID:          cfg.DonID,
			ConfigVersion:  cfg.ConfigVersion,
			F:              cfg.F,
			ForwarderState: solana.MustPublicKeyFromBase58(forwarderState.Address),
			ConfigPDA:      configPDA.String(),
			Type:           cldf.ContractType(ForwarderContract),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to configure forwarder for chain selector %d: %w", chain.Selector, err)
		}

		ops[chain.Selector] = opOut.Output.Batch
	}

	return ops, nil
}

func getConfigPDA(statePubkey solana.PublicKey, donID uint32, configVersion uint32, programID solana.PublicKey) solana.PublicKey {
	configID := getConfigID(donID, configVersion)
	reqIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(reqIDBytes, configID)

	seeds := [][]byte{
		[]byte("config"),
		statePubkey.Bytes(),
		reqIDBytes,
	}

	addr, _, _ := solana.FindProgramAddress(seeds, programID)
	return addr
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

func getConfigID(donID uint32, configVersion uint32) uint64 {
	return (uint64(donID) << 32) | uint64(configVersion)
}
