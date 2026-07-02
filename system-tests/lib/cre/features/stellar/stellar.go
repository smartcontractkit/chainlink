package stellar

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"dario.cat/mergo"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"google.golang.org/protobuf/types/known/durationpb"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	cre_jobs "github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	cre_jobs_ops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	job_types "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	stellchain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/stellar"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/features/jobhelpers"
)

// Stellar is the Local CRE feature for the Stellar chain capability (READ path).
//
// It registers the Stellar chain capability with its read methods and proposes
// the Stellar chain-capability worker jobs. It mirrors features/solana/v2 minus
// the forwarder (writes are Milestone B): no forwarder deploy, no WriteReport /
// LogTrigger method, no per-node transmitter. See STELLAR_LOCAL_CRE_PLAN.md A11.
type Stellar struct{}

var _ cre.Feature = (*Stellar)(nil)

const (
	flag              = cre.StellarCapability
	network           = "stellar"
	capabilityVersion = "1.0.0"
	requestTimeout    = 30 * time.Second
	// deltaStage is a write-transmission parameter; reads never use it, but the
	// capability config carries it. Soroban ledgers close ~5s.
	deltaStage = 5 * time.Second

	// placeholderForwarder is a syntactically-valid StrKey contract address used
	// only to satisfy the Stellar chain-capability config validation on the READ
	// path (config.go requires a valid C-address). Reads (GetLatestLedger /
	// ReadContract) never call the forwarder. Milestone B replaces this with the
	// address of the deployed CRE forwarder.
	placeholderForwarder = "CDLZFC3SYJYDZT7K67VZ75HPJVIEUVNIXF47ZG2FB2RMQQVU2HHGCYSC"
)

func (s *Stellar) Flag() cre.CapabilityFlag {
	return flag
}

func (s *Stellar) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	_ *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	stellarChain := extractStellarFromEnv(creEnv)
	if stellarChain == nil {
		// No Stellar blockchain in this environment; nothing to register.
		return &cre.PreEnvStartupOutput{}, nil
	}

	capabilities := registerStellarCapability(stellarChain.ChainSelector())

	capabilityToExtraSignerFamilies := make(map[string][]string, len(capabilities))
	ocrConfigs := map[string]*ocr3.OracleConfig{}
	for _, capability := range capabilities {
		// Chain read OCR & DON2DON use the EVM signing schema for all chains, so
		// EVM signers are required (the Soroban forwarder verifies secp256k1 too).
		capabilityToExtraSignerFamilies[capability.Capability.LabelledName] = []string{chainselectors.FamilyEVM}
		ocrConfigs[capability.Capability.LabelledName] = contracts.DefaultChainCapabilityOCR3Config()
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig:         capabilities,
		CapabilityToOCR3Config:          ocrConfigs,
		CapabilityToExtraSignerFamilies: capabilityToExtraSignerFamilies,
	}, nil
}

func (s *Stellar) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	if jobErr := createJobs(ctx, don, dons, creEnv); jobErr != nil {
		return errors.Wrap(jobErr, "failed to create jobs for stellar chain standard capability")
	}
	return nil
}

// registerStellarCapability registers the Stellar chain capability with its read
// methods. WriteReport / LogTrigger are Milestone B.
func registerStellarCapability(selector uint64) []keystone_changeset.DONCapabilityWithConfig {
	return []keystone_changeset.DONCapabilityWithConfig{{
		Capability: kcr.CapabilitiesRegistryCapability{
			LabelledName:   "stellar" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
			Version:        capabilityVersion,
			CapabilityType: 1,
		},
		Config: &capabilitiespb.CapabilityConfig{
			MethodConfigs: readMethodConfigs(),
		},
		UseCapRegOCRConfig: true,
	}}
}

// readMethodConfigs maps the Stellar capability's read RPCs (see the generated
// server switch in chainlink-common .../chain-capabilities/stellar) to a remote
// executable config. WriteReport is excluded (Milestone B).
func readMethodConfigs() map[string]*capabilitiespb.CapabilityMethodConfig {
	methodConfigs := make(map[string]*capabilitiespb.CapabilityMethodConfig)
	for _, action := range []string{"GetLatestLedger", "ReadContract"} {
		methodConfigs[action] = readActionConfig()
	}
	return methodConfigs
}

func readActionConfig() *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
			},
		},
	}
}

// workerConfig is the JSON config consumed by the Stellar chain-capability worker
// LOOPP (capabilities/chain_capabilities/stellar/config). For reads only network,
// chainId and a valid creForwarderAddress are meaningful; the rest carry defaults.
type workerConfig struct {
	Network             string `json:"network"`
	ChainID             string `json:"chainId"`
	CREForwarderAddress string `json:"creForwarderAddress"`
	DeltaStage          int64  `json:"deltaStage"`
	IsLocal             bool   `json:"isLocal"`
}

func createJobs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	stellarChain := extractStellarFromEnv(creEnv)
	if stellarChain == nil {
		return errors.New("no stellar blockchain found in environment")
	}

	var nodeSet cre.NodeSetWithCapabilityConfigs
	for _, ns := range dons.AsNodeSetWithChainCapabilities() {
		if ns.GetName() == don.Name {
			nodeSet = ns
			break
		}
	}
	if nodeSet == nil {
		return fmt.Errorf("could not find node set for Don named '%s'", don.Name)
	}

	config, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.CapabilityScope{})
	if resolveErr != nil {
		return errors.Wrap(resolveErr, "unable to find stellar capability config")
	}

	command, cErr := standardcapability.GetCommand(config.BinaryName)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to get command for stellar capability")
	}

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}
	bootstrapPeers := []string{fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"), bootstrap.Host, cre.OCRPeeringPort)}

	capRegVersion, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
	if !ok {
		return errors.New("CapabilitiesRegistry version not found in contract versions")
	}

	// Read config is uniform across worker nodes (no per-node transmitter), so a
	// single proposal fans out to every node in the DON.
	cfgBytes, mErr := json.Marshal(workerConfig{
		Network:             network,
		ChainID:             stellarChain.StellarChainID(),
		CREForwarderAddress: placeholderForwarder,
		DeltaStage:          int64(deltaStage),
		IsLocal:             true,
	})
	if mErr != nil {
		return errors.Wrap(mErr, "failed to marshal stellar worker config")
	}

	workerInput := cre_jobs.ProposeJobSpecInput{
		Domain:      offchain.ProductLabel,
		Environment: cre.EnvironmentName,
		DONName:     don.Name,
		// JD enforces a 63-char cap on job names (K8s label rule); the 64-char hex
		// Stellar chain id overflows it, so key the job on the short numeric selector.
		JobName:     "stellar-worker-" + strconv.FormatUint(stellarChain.ChainSelector(), 10),
		ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
		DONFilters: []offchain.TargetDONFilter{
			{Key: offchain.FilterKeyDONName, Value: don.Name},
		},
		Template: job_types.Stellar,
		Inputs: job_types.JobSpecInput{
			"command":            command,
			"config":             string(cfgBytes),
			"chainSelectorEVM":   creEnv.RegistryChainSelector,
			"bootstrapPeers":     bootstrapPeers,
			"useCapRegOCRConfig": true,
			"capRegVersion":      capRegVersion.String(),
		},
	}

	if verErr := (cre_jobs.ProposeJobSpec{}).VerifyPreconditions(*creEnv.CldfEnvironment, workerInput); verErr != nil {
		return fmt.Errorf("precondition verification failed for Stellar worker job: %w", verErr)
	}

	workerReport, workerErr := (cre_jobs.ProposeJobSpec{}).Apply(*creEnv.CldfEnvironment, workerInput)
	if workerErr != nil {
		return fmt.Errorf("failed to propose Stellar worker job spec: %w", workerErr)
	}

	specs := make(map[string][]string)
	for _, r := range workerReport.Reports {
		out, castOK := r.Output.(cre_jobs_ops.ProposeStandardCapabilityJobOutput)
		if !castOK {
			return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", r.Output)
		}
		if mergeErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice); mergeErr != nil {
			return fmt.Errorf("failed to merge worker job specs: %w", mergeErr)
		}
	}

	mergedSpecs, mergeErr := jobhelpers.MergeSpecsByIndex([]map[string][]string{specs})
	if mergeErr != nil {
		return mergeErr
	}

	if approveErr := jobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, mergedSpecs); approveErr != nil {
		return fmt.Errorf("failed to approve Stellar jobs: %w", approveErr)
	}

	return nil
}

func extractStellarFromEnv(creEnv *cre.Environment) *stellchain.Blockchain {
	for _, bcOut := range creEnv.Blockchains {
		if bcOut.IsFamily(chainselectors.FamilyStellar) {
			if sc, ok := bcOut.(*stellchain.Blockchain); ok {
				return sc
			}
		}
	}
	return nil
}
