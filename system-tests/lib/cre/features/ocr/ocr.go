package ocr

import (
	"context"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/offchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"

	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

const flag = cre.ConsensusCapability

type OCR struct{}

func (o *OCR) Flag() cre.CapabilityFlag {
	return flag
}

func (o *OCR) PreDONStartup(
	testLogger zerolog.Logger,
	registryChainSelector uint64,
	cldfEnv *cldf.Environment,
	provider infra.Provider,
	topology *cre.Topology,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	capabilityConfigs cre.CapabilityConfigs,
	contractVersions map[string]string,
) error {
	// nothing to do
	return nil
}

const (
	OCR3ContractQualifier    = "capability_ocr3"
	DONTimeContractQualifier = "capability_dontime"
)

func (o *OCR) PostDONStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	nodeSetOutput []*cre.WrappedNodeOutput,
	blockchainOutputs []*cre.WrappedBlockchainOutput,
	contractVersions map[string]string,
) (*cre.PostDONStartupOutput, error) {
	memoryDatastore := datastore.NewMemoryDataStore()

	// load all existing addresses into memory datastore
	mergeErr := memoryDatastore.Merge(creEnv.CldfEnvironment.DataStore)
	if mergeErr != nil {
		return nil, fmt.Errorf("failed to merge existing datastore into memory datastore: %w", mergeErr)
	}

	// deploy OCR3 contract
	_, ocrErr := deployOCR3Contract(OCR3ContractQualifier, creEnv.DonTopology.HomeChainSelector, creEnv.CldfEnvironment, memoryDatastore)
	if ocrErr != nil {
		return nil, fmt.Errorf("failed to deploy OCR3 contract %w", ocrErr)
	}

	ocr3Addr := MustGetAddressFromMemoryDataStore(memoryDatastore, creEnv.DonTopology.HomeChainSelector, keystone_changeset.OCR3Capability.String(), contractVersions[keystone_changeset.OCR3Capability.String()], OCR3ContractQualifier)
	testLogger.Info().Msgf("Deployed OCR3 %s contract on chain %d at %s", contractVersions[keystone_changeset.OCR3Capability.String()], creEnv.DonTopology.HomeChainSelector, ocr3Addr)

	// deploy don time contract
	_, timeErr := deployOCR3Contract(DONTimeContractQualifier, creEnv.DonTopology.HomeChainSelector, creEnv.CldfEnvironment, memoryDatastore) // Switch to dedicated config type once available
	if timeErr != nil {
		return nil, fmt.Errorf("failed to deploy DONTime contract %w", timeErr)
	}
	donTimeAddr := MustGetAddressFromMemoryDataStore(memoryDatastore, creEnv.DonTopology.HomeChainSelector, keystone_changeset.OCR3Capability.String(), contractVersions[keystone_changeset.OCR3Capability.String()], DONTimeContractQualifier)
	testLogger.Info().Msgf("Deployed OCR3 %s (DON Time) contract on chain %d at %s", contractVersions[keystone_changeset.OCR3Capability.String()], creEnv.DonTopology.HomeChainSelector, donTimeAddr)

	// update the CRE environment datastore to include the newly deployed contracts
	creEnv.CldfEnvironment.DataStore = memoryDatastore.Seal()

	// create OCR3 and don time jobs
	jobErr := createJobs(ctx, creEnv)
	if jobErr != nil {
		return nil, fmt.Errorf("failed to create OCR3 jobs: %w", jobErr)
	}

	testLogger.Info().Msg("OCR3 jobs created")

	// wait for LP to be started (otherwise it won't pick up contract's configuration events)
	if err := waitForLogPollerToBeHealthy(creEnv.DonTopology.Dons.List(), nodeSetOutput); err != nil {
		return nil, errors.Wrap(err, "failed while waiting for Log Poller to become healthy")
	}

	// configure OCR3 contract
	dons, donsErr := toDons(creEnv)
	if donsErr != nil {
		return nil, fmt.Errorf("failed to convert to dons: %w", donsErr)
	}

	consensusV1DON, donErr := dons.shouldBeOneDon(cre.ConsensusCapability)
	if donErr != nil {
		return nil, fmt.Errorf("failed to get consensus v1 DON: %w", donErr)
	}

	ocr3Config, ocr3confErr := defaultOCR3Config(creEnv.DonTopology)
	if ocr3confErr != nil {
		return nil, fmt.Errorf("failed to get default OCR3 config: %w", ocr3confErr)
	}

	_, ocr3Err := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: &ocr3Addr,
			ChainSelector:   creEnv.DonTopology.HomeChainSelector,
			DON:             consensusV1DON.keystoneDonConfig(),
			Config:          consensusV1DON.resolveOcr3Config(*ocr3Config),
			DryRun:          false,
		},
	)
	if ocr3Err != nil {
		return nil, errors.Wrap(ocr3Err, "failed to configure OCR3 contract")
	}

	// configure DON Time contract
	// don time happens to be the same as consensus v1 DON, but it doesn't have to be
	_, donTimeErr := operations.ExecuteOperation(
		creEnv.CldfEnvironment.OperationsBundle,
		ks_contracts_op.ConfigureOCR3Op,
		ks_contracts_op.ConfigureOCR3OpDeps{
			Env: creEnv.CldfEnvironment,
		},
		ks_contracts_op.ConfigureOCR3OpInput{
			ContractAddress: &donTimeAddr,
			ChainSelector:   creEnv.DonTopology.HomeChainSelector,
			DON:             consensusV1DON.keystoneDonConfig(),
			Config:          consensusV1DON.resolveOcr3Config(*ocr3Config),
			DryRun:          false,
		},
	)
	if donTimeErr != nil {
		return nil, errors.Wrap(donTimeErr, "failed to configure DON Time contract")
	}

	// return capabilities registry configuration data
	capabilities := make(map[int][]keystone_changeset.DONCapabilityWithConfig)
	for donIdx, don := range creEnv.DonTopology.Dons.List() {
		if capabilities[donIdx] == nil {
			capabilities[donIdx] = []keystone_changeset.DONCapabilityWithConfig{}
		}
		if don.HasFlag(flag) {
			capabilities[donIdx] = append(capabilities[donIdx], keystone_changeset.DONCapabilityWithConfig{
				Capability: kcr.CapabilitiesRegistryCapability{
					LabelledName:   "offchain_reporting",
					Version:        "1.0.0",
					CapabilityType: 2, // CONSENSUS
					ResponseType:   0, // REPORT
				},
				Config: &capabilitiespb.CapabilityConfig{},
			})
		}
	}

	return &cre.PostDONStartupOutput{
		DONCapabilityWithConfigs: capabilities,
	}, nil
}

func createJobs(
	ctx context.Context,
	creEnv *cre.Environment,
) error {
	ocr3Key := datastore.NewAddressRefKey(
		creEnv.DonTopology.HomeChainSelector,
		datastore.ContractType(keystone_changeset.OCR3Capability.String()),
		semver.MustParse("1.0.0"),
		contracts.OCR3ContractQualifier,
	)
	ocr3CapabilityAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(ocr3Key)
	if err != nil {
		return errors.Wrap(err, "failed to get OCR3 capability address")
	}

	donTimeKey := datastore.NewAddressRefKey(
		creEnv.DonTopology.HomeChainSelector,
		datastore.ContractType(keystone_changeset.OCR3Capability.String()),
		semver.MustParse("1.0.0"),
		contracts.DONTimeContractQualifier,
	)
	donTimeAddress, err := creEnv.CldfEnvironment.DataStore.Addresses().Get(donTimeKey)
	if err != nil {
		return errors.Wrap(err, "failed to get DON Time address")
	}

	jobSpecs := []*job.ProposeJobRequest{}

	chainID, chErr := chainselectors.ChainIdFromSelector(creEnv.DonTopology.HomeChainSelector)
	if chErr != nil {
		return errors.Wrapf(chErr, "failed to get chain ID from chain selector %d", creEnv.DonTopology.HomeChainSelector)
	}

	for _, don := range creEnv.DonTopology.Dons.List() {
		if !don.HasFlag(flag) {
			continue
		}

		workerNodes, wErr := don.WorkerNodes()
		if wErr != nil {
			return errors.Wrap(wErr, "failed to find worker nodes")
		}

		bootstrapNode, err := creEnv.DonTopology.BootstrapNode()
		if err != nil {
			return errors.Wrap(err, "failed to get bootstrap node from DON metadata")
		}

		_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrapNode)
		if err != nil {
			return errors.Wrap(err, "failed to get peering configs")
		}

		jobSpecs = append(jobSpecs, jobs.BootstrapOCR3(bootstrapNode.JobDistributorDetails.NodeID, "ocr3-capability", ocr3CapabilityAddress.Address, chainID))

		for _, workerNode := range workerNodes {
			evmKey, ok := workerNode.Keys.EVM[chainID]
			if !ok {
				return fmt.Errorf("failed to get EVM key (chainID %d, node index %d)", chainID, workerNode.Index)
			}

			// we need the OCR2 key bundle for the EVM chain, because OCR jobs currently run only on EVM chains
			evmOCR2KeyBundle, ok := workerNode.Keys.OCR2BundleIDs[chainselectors.FamilyEVM]
			if !ok {
				return fmt.Errorf("node %s does not have OCR2 key bundle for evm", workerNode.Name)
			}

			// we pass here bundles for all chains to enable multi-chain signing
			jobSpecs = append(jobSpecs, jobs.WorkerOCR3(workerNode.JobDistributorDetails.NodeID, ocr3CapabilityAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, workerNode.Keys.OCR2BundleIDs, ocrPeeringCfg, chainID))
			jobSpecs = append(jobSpecs, jobs.DonTimeJob(workerNode.JobDistributorDetails.NodeID, donTimeAddress.Address, evmKey.PublicAddress.Hex(), evmOCR2KeyBundle, ocrPeeringCfg, chainID))
		}
	}

	return jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, creEnv.DonTopology, jobSpecs)
}

// TODO make our DON use clclient from the CTF instead, since then we could remove nodeSetOutput as DON has Nodes that would have this client initialised
func waitForLogPollerToBeHealthy(dons []*cre.DON, nodeSetOutput []*cre.WrappedNodeOutput) error {
	for idx, nodeSetOut := range nodeSetOutput {
		if !dons[idx].HasFlag(cre.ConsensusCapability) || !dons[idx].HasFlag(cre.VaultCapability) {
			continue
		}
		nsClients, cErr := clclient.New(nodeSetOut.CLNodes)
		if cErr != nil {
			return errors.Wrap(cErr, "failed to create node set clients")
		}
		eg := &errgroup.Group{}
		for _, c := range nsClients {
			eg.Go(func() error {
				return c.WaitHealthy(".*ConfigWatcher", "passing", 100)
			})
		}
		if waitErr := eg.Wait(); waitErr != nil {
			return errors.Wrap(waitErr, "failed to wait for ConfigWatcher health check")
		}
	}

	return nil
}

/*
	CODE BELOW WAS COPIED FROM VAROIUS PLACES IN THE SYSTEM TESTS AND (SOMETIMES) MODIFIED
	TO SHOWCASE THE CONCEPT OF FEATURES.

	IN THE FUTURE, WE SHOULD REFACTOR THE CODE TO AVOID DUPLICATIONS.
*/

// for now copy from system-tests/lib/cre/contracts/keystone.go
func deployOCR3Contract(qualifier string, selector uint64, env *cldf.Environment, ds datastore.MutableDataStore) (*ks_contracts_op.DeployOCR3ContractSequenceOutput, error) {
	ocr3DeployReport, err := operations.ExecuteSequence(
		env.OperationsBundle,
		ks_contracts_op.DeployOCR3ContractsSequence,
		ks_contracts_op.DeployOCR3ContractSequenceDeps{
			Env: env,
		},
		ks_contracts_op.DeployOCR3ContractSequenceInput{
			ChainSelector: selector,
			Qualifier:     qualifier,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy OCR3 contract '%s' on chain %d: %w", qualifier, selector, err)
	}
	// TODO: CRE-742 remove address book
	if err = env.ExistingAddresses.Merge(ocr3DeployReport.Output.AddressBook); err != nil { //nolint:staticcheck // won't migrate now
		return nil, fmt.Errorf("failed to merge address book with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}
	if err = ds.Merge(ocr3DeployReport.Output.Datastore); err != nil {
		return nil, fmt.Errorf("failed to merge datastore with OCR3 contract address for '%s' on chain %d: %w", qualifier, selector, err)
	}
	return &ocr3DeployReport.Output, nil
}

// for now copy from system-tests/lib/cre/contracts/keystone.go
func MustGetAddressFromMemoryDataStore(dataStore *datastore.MemoryDataStore, chainSel uint64, contractType string, version string, qualifier string) common.Address {
	key := datastore.NewAddressRefKey(
		chainSel,
		datastore.ContractType(contractType),
		semver.MustParse(version),
		qualifier,
	)
	addrRef, err := dataStore.Addresses().Get(key)
	if err != nil {
		panic(fmt.Sprintf("Failed to get %s %s (qualifier=%s) address for chain %d: %s", contractType, version, qualifier, chainSel, err.Error()))
	}
	return common.HexToAddress(addrRef.Address)
}

// copied from system-tests/lib/cre/contracts/contracts.go and modified
func toDons(
	creEnv *cre.Environment,
) (*dons, error) {
	dons := &dons{
		c:        make(map[string]donConfig),
		offChain: creEnv.CldfEnvironment.Offchain,
	}

	for _, don := range creEnv.DonTopology.Dons.List() {
		// if it's only a gateway DON, we don't want to register it with the Capabilities Registry
		// since it doesn't have any capabilities
		if flags.HasOnlyOneFlag(don.Flags, cre.GatewayDON) {
			continue
		}

		var capabilities []keystone_changeset.DONCapabilityWithConfig

		// // check what capabilities each DON has and register them with Capabilities Registry contract
		// for _, configFn := range input.CapabilityRegistryConfigFns {
		// 	if configFn == nil {
		// 		continue
		// 	}

		// 	enabledCapabilities, err2 := configFn(don.Flags, input.NodeSets[donIdx])
		// 	if err2 != nil {
		// 		return nil, errors.Wrap(err2, "failed to get capabilities from config function")
		// 	}

		// 	capabilities = append(capabilities, enabledCapabilities...)
		// }

		workerNodes, wErr := don.WorkerNodes()
		if wErr != nil {
			return nil, errors.Wrap(wErr, "failed to find worker nodes")
		}

		donPeerIDs := make([]string, len(workerNodes))
		for i, node := range workerNodes {
			// we need to use p2pID here with the "p2p_" prefix
			donPeerIDs[i] = node.Keys.P2PKey.PeerID.String()
		}

		forwarderF := (len(workerNodes) - 1) / 3
		if forwarderF == 0 {
			if flags.HasFlag(don.Flags, cre.ConsensusCapability) || flags.HasFlag(don.Flags, cre.ConsensusCapabilityV2) {
				return nil, fmt.Errorf("incorrect number of worker nodes: %d. Resulting F must conform to formula: mod((N-1)/3) > 0", len(workerNodes))
			}
			// for other capabilities, we can use 1 as F
			forwarderF = 1
		}

		// we only need to assign P2P IDs to NOPs, since `ConfigureInitialContractsChangeset` method
		// will take care of creating DON to Nodes mapping
		nop := keystone_changeset.NOP{
			Name:  fmt.Sprintf("NOP for %s DON", don.Name),
			Nodes: donPeerIDs,
		}
		donName := don.Name + "-don"
		c := keystone_changeset.DonCapabilities{
			Name:         donName,
			F:            libc.MustSafeUint8(forwarderF),
			Nops:         []keystone_changeset.NOP{nop},
			Capabilities: capabilities,
		}

		dons.c[donName] = donConfig{
			id:              uint32(don.ID), //nolint:gosec // G115
			DonCapabilities: c,
			flags:           don.Flags,
		}
	}

	return dons, nil
}

// copied from system-tests/lib/cre/contracts/contracts.go
type donConfig struct {
	id uint32 // the DON id as registered in the capabilities registry
	keystone_changeset.DonCapabilities
	flags []cre.CapabilityFlag
}

func (d *donConfig) resolveOcr3Config(c keystone_changeset.OracleConfig) *keystone_changeset.OracleConfig {
	c.TransmissionSchedule = []int{d.N()}
	return &c
}

func (d *donConfig) keystoneDonConfig() ks_contracts_op.ConfigureKeystoneDON {
	don := ks_contracts_op.ConfigureKeystoneDON{
		Name: d.Name,
	}
	for _, nop := range d.Nops {
		don.NodeIDs = append(don.NodeIDs, nop.Nodes...)
	}
	return don
}

type dons struct {
	c        map[string]donConfig
	offChain offchain.Client
}

func (d *dons) donsOrderedByID() []donConfig {
	out := make([]donConfig, 0, len(d.c))
	for _, don := range d.c {
		out = append(out, don)
	}

	// Use sort library to sort by ID
	sort.Slice(out, func(i, j int) bool {
		return out[i].id < out[j].id
	})

	return out
}

func (d *dons) ListByFlag(flag cre.CapabilityFlag) ([]donConfig, error) {
	out := make([]donConfig, 0)
	for _, don := range d.donsOrderedByID() {
		if flags.HasFlag(don.flags, flag) {
			out = append(out, don)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("don with flag %s not found", flag)
	}
	return out, nil
}

func (d *dons) shouldBeOneDon(flag cre.CapabilityFlag) (donConfig, error) {
	dons, err := d.ListByFlag(flag)
	if err != nil {
		return donConfig{}, err
	}
	if len(dons) != 1 {
		return donConfig{}, fmt.Errorf("expected exactly one DON with flag %s, found %d", flag, len(dons))
	}
	return dons[0], nil
}

// copied from system-tests/lib/cre/contracts/contracts.go
func defaultOCR3Config(donTopology *cre.DonTopology) (*keystone_changeset.OracleConfig, error) {
	var transmissionSchedule []int

	for _, don := range donTopology.Dons.List() {
		if don.HasFlag(cre.ConsensusCapability) || don.HasFlag(cre.ConsensusCapabilityV2) {
			workerNodes, wErr := don.WorkerNodes()
			if wErr != nil {
				return nil, errors.Wrap(wErr, "failed to find worker nodes")
			}

			// this schedule makes sure that all worker nodes are transmitting OCR3 reports
			transmissionSchedule = []int{len(workerNodes)}
			break
		}
	}

	if len(transmissionSchedule) == 0 {
		return nil, errors.New("no OCR3-capable DON found in the topology")
	}

	// values supplied by Alexandr Yepishev as the expected values for OCR3 config
	oracleConfig := &keystone_changeset.OracleConfig{
		DeltaProgressMillis:               5000,
		DeltaResendMillis:                 5000,
		DeltaInitialMillis:                5000,
		DeltaRoundMillis:                  2000,
		DeltaGraceMillis:                  500,
		DeltaCertifiedCommitRequestMillis: 1000,
		DeltaStageMillis:                  30000,
		MaxRoundsPerEpoch:                 10,
		TransmissionSchedule:              transmissionSchedule,
		MaxDurationQueryMillis:            1000,
		MaxDurationObservationMillis:      1000,
		MaxDurationShouldAcceptMillis:     1000,
		MaxDurationShouldTransmitMillis:   1000,
		MaxFaultyOracles:                  1,
		ConsensusCapOffchainConfig: &ocr3.ConsensusCapOffchainConfig{
			MaxQueryLengthBytes:       1000000,
			MaxObservationLengthBytes: 1000000,
			MaxOutcomeLengthBytes:     1000000,
			MaxReportLengthBytes:      1000000,
			MaxBatchSize:              1000,
		},
		UniqueReports: true,
	}

	return oracleConfig, nil
}
