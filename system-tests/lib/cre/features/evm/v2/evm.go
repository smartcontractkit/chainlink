package v2

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/durationpb"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"

	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	ks_contracts_op "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/operations/contracts"
	libc "github.com/smartcontractkit/chainlink/system-tests/lib/conversions"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
)

const (
	flag                = cre.EVMCapability
	registrationRefresh = 20 * time.Second
	registrationExpiry  = 60 * time.Second
	deltaStage          = 500*time.Millisecond + 1*time.Second // block time + 1 second delta
	requestTimeout      = 30 * time.Second
)

type EVM struct{}

func (o *EVM) Flag() cre.CapabilityFlag {
	return flag
}

func (o *EVM) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	topology *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	chainsWithForwarders := chainsWithForwarders(creEnv.Blockchains, cre.ConvertToNodeSetWithChainCapabilities(topology.NodeSets()))
	evmForwardersSelectors, exist := chainsWithForwarders[blockchain.FamilyEVM]

	if exist {
		selectorsToDeploy := make([]uint64, 0)
		for _, selector := range evmForwardersSelectors {
			// filter out EVM forwarder selectors that might have been already deployed by evm_v1 capability
			forwarderAddr := contracts.MightGetAddressFromDataStore(creEnv.CldfEnvironment.DataStore, selector, keystone_changeset.KeystoneForwarder.String(), creEnv.ContractVersions[keystone_changeset.KeystoneForwarder.String()], "")
			if forwarderAddr == nil {
				selectorsToDeploy = append(selectorsToDeploy, selector)
			}
		}

		if len(selectorsToDeploy) > 0 {
			deployErr := deployEVMForwarders(testLogger, creEnv.CldfEnvironment, selectorsToDeploy, creEnv.ContractVersions)
			if deployErr != nil {
				return nil, errors.Wrap(deployErr, "failed to deploy EVM Keystone forwarder")
			}
		}
	}

	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}

	capabilities := []keystone_changeset.DONCapabilityWithConfig{}

	for _, chainID := range enabledChainIDs {
		selector, selectorErr := chainselectors.SelectorFromChainId(chainID)
		if selectorErr != nil {
			return nil, errors.Wrapf(selectorErr, "failed to get selector from chainID: %d", chainID)
		}

		evmMethodConfigs, err := getEvmMethodConfigs(don.MustNodeSet())
		if err != nil {
			return nil, errors.Wrap(err, "there was an error getting EVM method configs")
		}

		capabilities = append(capabilities, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName: "evm" + ":ChainSelector:" + strconv.FormatUint(selector, 10),
				Version:      "1.0.0",
			},
			Config: &capabilitiespb.CapabilityConfig{
				MethodConfigs: evmMethodConfigs,
				LocalOnly:     don.HasOnlyLocalCapabilities(),
			},
			UseCapRegOCRConfig: true,
		})
	}

	capabilityToOCR3Config := make(map[string]*ocr3.OracleConfig, len(capabilities))
	for _, cap := range capabilities {
		capabilityToOCR3Config[cap.Capability.LabelledName] = contracts.DefaultChainCapabilityOCR3Config()
	}

	return &cre.PreEnvStartupOutput{
		DONCapabilityWithConfig: capabilities,
		CapabilityToOCR3Config:  capabilityToOCR3Config,
	}, nil
}

func (o *EVM) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	jobsErr := createJobs(
		ctx,
		don,
		dons,
		creEnv,
	)
	if jobsErr != nil {
		return jobsErr
	}

	// configure EVM forwarders for DONs that run consensus
	consensusDons := dons.DonsWithFlags(cre.ConsensusCapability)

	// Forwarders may be configured when multiple DONs share chains with EVM capability; duplicate configuration is harmless.
	chainsWithEVMCapability := chainsWithEVMCapability(creEnv.Blockchains, dons.DonsWithFlag(flag))
	if len(chainsWithEVMCapability) > 0 {
		evmChainsWithForwarders := make([]uint64, 0)
		for _, chainSelector := range chainsWithEVMCapability {
			evmChainsWithForwarders = append(evmChainsWithForwarders, uint64(chainSelector))
		}
		for _, don := range consensusDons {
			config, confErr := configureEVMForwarders(testLogger, creEnv.CldfEnvironment, evmChainsWithForwarders, don)
			if confErr != nil {
				return errors.Wrap(confErr, "failed to configure EVM forwarders")
			}
			testLogger.Info().Msgf("Configured EVM forwarders: %+v", config)
		}
	}

	return nil
}

func chainsWithEVMCapability(chains []blockchains.Blockchain, dons []*cre.Don) map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector {
	chainsWithEVMCapability := make(map[ks_contracts_op.EVMChainID]ks_contracts_op.Selector)
	for _, chain := range chains {
		for _, don := range dons {
			if flags.HasFlagForChain(don.Flags, cre.EVMCapability, chain.ChainID()) {
				if chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID())] != 0 {
					continue
				}
				chainsWithEVMCapability[ks_contracts_op.EVMChainID(chain.ChainID())] = ks_contracts_op.Selector(chain.ChainSelector())
			}
		}
	}

	return chainsWithEVMCapability
}

// The EVM binary runs as a capability runner rather than as a standard
// capability: it reaches the chain itself - its own RPC client, log poller and
// transaction manager over its own tables - and borrows from the crecore process
// beside it only what a node alone has. Those addresses are the ones the node
// launches crecore with, so they must match the [Capabilities.Proxy] block in the
// topology config (GRPCPort), which is what enables crecore in the first place.
const (
	// creCoreGRPCTarget is crecore's single gRPC address: the OCR proxy, the chain
	// keys it signs with, and the capabilities registry are all served on it.
	creCoreGRPCTarget = "localhost:50051"

	// runnerHTTPPortBase is where this capability's runners serve /metrics,
	// /debug/pprof, health and - the part the node needs - the settings reload
	// endpoint. One per chain, numbered from here, since a node running the EVM
	// capability on two chains runs two of these. Loopback inside the node's
	// container; it only has to avoid crecore's 50051/50052 and the ports the other
	// capability runners take (consensus 50053, cron 50054).
	runnerHTTPPortBase = 50060
)

// evmFlagPrefix is the namespace the EVM binary registers its own settings under,
// so a value from capability_defaults.toml becomes a flag.
const evmFlagPrefix = "--evm."

// evmValueFlags maps the keys accepted under [capability_configs.evm.values] to
// the binary's flag names. Unknown keys are rejected rather than dropped: a value
// set in the TOML and silently ignored looks like the binary disagreeing with its
// configuration.
//
// The keys that name something this feature works out for itself - the chain, the
// forwarder it deploys, the account a node signs with - are not here; they are
// passed below and would be two sources for one answer.
var evmValueFlags = map[string]string{
	"LogTriggerPollInterval":          "log-trigger-poll-interval",
	"LogTriggerSendChannelBufferSize": "log-trigger-send-channel-buffer-size",
	"LogTriggerLimitQueryLogSize":     "log-trigger-limit-query-log-size",
	"ReceiverGasMinimum":              "receiver-gas-minimum",
	"ForwarderLookbackBlocks":         "forwarder-lookback-blocks",
	"ObservationPollerWorkersCount":   "observation-poller-workers-count",
	"ObservationPollPeriod":           "observation-poll-period",
	"ChainHeightPollPeriod":           "chain-height-poll-period",
	"UnknownRequestsTTL":              "unknown-requests-ttl",
}

// runtimeValueFlags are the values this feature resolves rather than reads: they
// come from the environment it just built, so a TOML naming them is a
// configuration that can disagree with the deployment.
var runtimeValueFlags = map[string]bool{
	"ChainID": true, "NetworkFamily": true, "CreForwarderAddress": true, "NodeAddress": true, "DeltaStage": true,
}

func createJobs(
	ctx context.Context,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
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

	bootstrap, isBootstrap := dons.Bootstrap()
	if !isBootstrap {
		return errors.New("could not find bootstrap node in topology, exactly one bootstrap node is required")
	}

	workerNodes, wErr := don.Workers()
	if wErr != nil {
		return errors.Wrap(wErr, "failed to find worker nodes")
	}

	enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}
	// Sorted so that a chain's HTTP port, which is its position in this list, is the
	// same on every run rather than following the order the capability was declared.
	slices.Sort(enabledChainIDs)

	registryChainID, rcErr := chainselectors.ChainIdFromSelector(creEnv.RegistryChainSelector)
	if rcErr != nil {
		return fmt.Errorf("failed to get chain ID from registry chain selector %d: %w", creEnv.RegistryChainSelector, rcErr)
	}

	specs := make(cre.DonJobs, 0, len(enabledChainIDs)*len(workerNodes))
	for i, chainID := range enabledChainIDs {
		chain, chainErr := evmChain(creEnv, chainID)
		if chainErr != nil {
			return chainErr
		}

		chainSelector, selErr := chainselectors.SelectorFromChainId(chainID)
		if selErr != nil {
			return errors.Wrapf(selErr, "failed to get chain selector from chainID %d", chainID)
		}

		capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
		if resolveErr != nil {
			return fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, resolveErr)
		}

		command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
		if cErr != nil {
			return errors.Wrap(cErr, "failed to get command for the EVM capability")
		}

		configFlags, configErr := buildConfigFlags(capabilityConfig)
		if configErr != nil {
			return configErr
		}

		forwarderKey := datastore.NewAddressRefKey(
			chainSelector,
			datastore.ContractType(keystone_changeset.KeystoneForwarder.String()),
			semver.MustParse("1.0.0"),
			"",
		)
		creForwarder, fErr := creEnv.CldfEnvironment.DataStore.Addresses().Get(forwarderKey)
		if fErr != nil {
			return errors.Wrap(fErr, "failed to get CRE Forwarder address")
		}

		for _, workerNode := range workerNodes {
			spec, specErr := capabilityRunnerJobSpec(runnerInputs{
				node:                workerNode,
				don:                 don,
				chain:               chain,
				chainID:             chainID,
				registryChainID:     registryChainID,
				httpPort:            runnerHTTPPortBase + i,
				command:             command,
				creForwarderAddress: creForwarder.Address,
				bootstrapPeer:       formatBootstrapPeer(bootstrap),
				configFlags:         configFlags,
			})
			if specErr != nil {
				return specErr
			}
			specs = append(specs, spec)
		}
	}
	if len(specs) == 0 {
		return fmt.Errorf("no worker nodes found in %s DON to run the EVM capability", don.GetName())
	}

	if err := jobs.Create(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to create EVM jobs: %w", err)
	}

	return nil
}

// runnerInputs is one job: one node, running this capability against one chain.
type runnerInputs struct {
	node            *cre.Node
	don             *cre.Don
	chain           blockchains.Blockchain
	chainID         uint64
	registryChainID uint64
	httpPort        int

	command             string
	creForwarderAddress string
	bootstrapPeer       string
	configFlags         []string
}

// capabilityRunnerJobSpec builds the capabilityrunner job that launches the EVM
// binary on one node, for one chain.
//
// The node supervises the process over the empty LOOP and notifies it of CRE
// settings updates on --http.port; the database it keeps its chain state in comes
// from the node's own CL_DATABASE_URL, which the process inherits, in a schema of
// its own. Everything else it is told here, since a process that hosts no node
// cannot look any of it up: the chain to dial, the peer and keys crecore holds on
// its behalf, and which DON it is.
func capabilityRunnerJobSpec(in runnerInputs) (*jobv1.ProposeJobRequest, error) {
	if in.node.JobDistributorDetails == nil {
		return nil, fmt.Errorf("node %s has no job distributor details", in.node.Name)
	}
	peerID := strings.TrimPrefix(in.node.PeerID(), "p2p_")
	if peerID == "" {
		return nil, fmt.Errorf("node %s has no P2P peer ID", in.node.Name)
	}

	// The account this node transmits from on the registry chain, which is what the
	// OCR configuration lists it under, and the account it sends this chain's
	// transactions from - the same node, two chains, two keys.
	transmitAccount, err := transmitAccountFor(in.node, in.registryChainID)
	if err != nil {
		return nil, err
	}
	chainKey, ok := in.node.Keys.EVM[in.chainID]
	if !ok {
		return nil, fmt.Errorf("node %s has no EVM key for chain %d", in.node.Name, in.chainID)
	}

	rpc := in.chain.CtfOutput().Nodes[0]

	args := []string{
		// One instance, configured as it is here. The binary also has an "embed"
		// subcommand for running several in one process, which is for local runs
		// rather than for the node.
		"run",
		fmt.Sprintf("--http.port=%d", in.httpPort),
		"--ocr.proxy-address=" + creCoreGRPCTarget,
		"--ocr.peer-id=" + peerID,
		"--ocr.transmit-account=" + transmitAccount,
		"--ocr.bootstrappers=" + in.bootstrapPeer,
		"--keystore.proxy-address=" + creCoreGRPCTarget,
		"--capabilities.proxy-url=" + creCoreGRPCTarget,
		fmt.Sprintf("--capabilities.capability-don-id=%d", in.don.ID),
		// Its own schema in the node's database, created by the node's migrations
		// (0303_create_cre_standalone_schemas.sql): the tables are chainlink-evm's, and
		// the node's own copies of them are not this capability's to share. One schema
		// for every chain, not one per chain - these tables carry an evm_chain_id, which
		// is how the node keeps every chain in its own single evm schema.
		"--database.schema=evm_capability",
		fmt.Sprintf("--evm.chain-id=%d", in.chainID),
		"--evm.http-url=" + rpc.InternalHTTPUrl,
		"--evm.cre-forwarder-address=" + in.creForwarderAddress,
		"--evm.node-address=" + chainKey.PublicAddress.Hex(),
		"--evm.delta-stage=" + deltaStage.String(),
	}
	if rpc.InternalWSUrl != "" {
		args = append(args, "--evm.ws-url="+rpc.InternalWSUrl)
	} else {
		// Without a websocket the RPC pool has to poll for heads, and a pool doing
		// neither declares the chain unreachable.
		args = append(args, "--evm.new-heads-poll-interval=1s")
	}
	args = append(args, in.configFlags...)

	quoted := make([]string, 0, len(args))
	for _, a := range args {
		quoted = append(quoted, fmt.Sprintf("%q", a))
	}

	return &jobv1.ProposeJobRequest{
		NodeId: in.node.JobDistributorDetails.NodeID,
		Spec: fmt.Sprintf(`
type = "capabilityrunner"
schemaVersion = 1
externalJobID = "%s"
name = "evm-worker-%d"
command = "%s"
args = [%s]
`,
			uuid.NewString(),
			in.chainID,
			in.command,
			strings.Join(quoted, ", "),
		),
	}, nil
}

// buildConfigFlags turns [capability_configs.evm.values] into --evm.* flags. An
// empty set leaves every field at the binary's own default.
func buildConfigFlags(capConfig cre.CapabilityConfig) ([]string, error) {
	keys := make([]string, 0, len(capConfig.Values))
	for k := range capConfig.Values {
		keys = append(keys, k)
	}
	// Sorted so the args - and so the job spec a node is proposed - are the same on
	// every run rather than following map iteration order.
	sort.Strings(keys)

	configFlags := make([]string, 0, len(keys))
	for _, k := range keys {
		if runtimeValueFlags[k] {
			// Resolved from the environment below; a configured copy could only disagree
			// with it, so it is ignored rather than argued with.
			continue
		}
		name, ok := evmValueFlags[k]
		if !ok {
			return nil, fmt.Errorf("unknown EVM capability config value %q; the EVM binary accepts %s",
				k, strings.Join(sortedKeys(evmValueFlags), ", "))
		}
		configFlags = append(configFlags, fmt.Sprintf("%s%s=%v", evmFlagPrefix, name, capConfig.Values[k]))
	}
	return configFlags, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// evmChain is the blockchain this capability reads, as the environment built it.
func evmChain(creEnv *cre.Environment, chainID uint64) (blockchains.Blockchain, error) {
	for _, bc := range creEnv.Blockchains {
		if bc.IsFamily(blockchain.FamilyEVM) && bc.ChainID() == chainID {
			return bc, nil
		}
	}
	return nil, fmt.Errorf("no EVM chain %d in this environment", chainID)
}

// transmitAccountFor returns the account the node is registered to transmit from:
// its EVM address on the registry chain, as lowercase hex with no 0x prefix.
//
// The presentation is the point, not just the address. libocr compares this to the
// config as an account string, and the config's copy is produced by whoever decodes
// it: for a capability read out of the CapabilitiesRegistry that is
// capabilitiespb.OCR3ConfigFromProto, which renders each on-chain transmitter as
// hex.EncodeToString of its 20 bytes. So this renders it the same way rather than
// reading it back from JD, which reports the EIP-55 form and would never match.
func transmitAccountFor(node *cre.Node, registryChainID uint64) (string, error) {
	if node.Keys == nil {
		return "", fmt.Errorf("node %s has no keys", node.Name)
	}
	key, ok := node.Keys.EVM[registryChainID]
	if !ok || key == nil {
		return "", fmt.Errorf("node %s has no EVM key for the registry chain %d, so it has no account to transmit from", node.Name, registryChainID)
	}
	return hex.EncodeToString(key.PublicAddress.Bytes()), nil
}

func formatBootstrapPeer(bootstrap *cre.Node) string {
	return fmt.Sprintf("%s@%s:%d",
		strings.TrimPrefix(bootstrap.Keys.PeerID(), "p2p_"),
		bootstrap.Host,
		cre.OCRPeeringPort)
}

// getEvmMethodConfigs returns the method configs for all EVM methods we want to support, if any method is missing it
// will not be reached by the node when running evm capability in remote don
func getEvmMethodConfigs(nodeSet *cre.NodeSet) (map[string]*capabilitiespb.CapabilityMethodConfig, error) {
	evmMethodConfigs := map[string]*capabilitiespb.CapabilityMethodConfig{}

	// the read actions should be all defined in the proto that are neither a LogTrigger type, not a WriteReport type
	// see the RPC methods to map here: https://github.com/smartcontractkit/chainlink-protos/blob/main/cre/capabilities/blockchain/evm/v1alpha/client.proto
	readActions := []string{
		"CallContract",
		"FilterLogs",
		"BalanceAt",
		"EstimateGas",
		"GetTransactionByHash",
		"GetTransactionReceipt",
		"HeaderByNumber",
	}
	for _, action := range readActions {
		evmMethodConfigs[action] = readActionConfig()
	}

	triggerConfig, err := logTriggerConfig(nodeSet)
	if err != nil {
		return nil, errors.Wrap(err, "failed get config for LogTrigger")
	}

	evmMethodConfigs["LogTrigger"] = triggerConfig
	evmMethodConfigs["WriteReport"] = writeReportActionConfig()
	return evmMethodConfigs, nil
}

func logTriggerConfig(nodeSet *cre.NodeSet) (*capabilitiespb.CapabilityMethodConfig, error) {
	faultyNodes, faultyErr := nodeSet.MaxFaultyNodes()
	if faultyErr != nil {
		return nil, errors.Wrap(faultyErr, "failed to get faulty nodes")
	}

	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteTriggerConfig{
			RemoteTriggerConfig: &capabilitiespb.RemoteTriggerConfig{
				RegistrationRefresh:     durationpb.New(registrationRefresh),
				RegistrationExpiry:      durationpb.New(registrationExpiry),
				MinResponsesToAggregate: faultyNodes + 1,
				MessageExpiry:           durationpb.New(2 * registrationExpiry),
				MaxBatchSize:            25,
				BatchCollectionPeriod:   durationpb.New(200 * time.Millisecond),
			},
		},
	}, nil
}

func writeReportActionConfig() *capabilitiespb.CapabilityMethodConfig {
	return &capabilitiespb.CapabilityMethodConfig{
		RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
			RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
				RequestTimeout:            durationpb.New(requestTimeout),
				ServerMaxParallelRequests: 10,
				RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
			},
		},
	}
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

func deployEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, contractVersions map[cre.ContractType]*semver.Version) error {
	memoryDatastore, mErr := contracts.NewDataStoreFromExisting(cldfEnv.DataStore)
	if mErr != nil {
		return fmt.Errorf("failed to create memory datastore: %w", mErr)
	}

	evmForwardersReport, deployErr := operations.ExecuteSequence(
		cldfEnv.OperationsBundle,
		forwarder.DeploySequence,
		forwarder.DeploySequenceDeps{
			Env: cldfEnv,
		},
		forwarder.DeploySequenceInput{
			Targets: chainSelectors,
		},
	)
	if deployErr != nil {
		return errors.Wrap(deployErr, "failed to deploy evm forwarder")
	}

	if err := memoryDatastore.Merge(evmForwardersReport.Output.Datastore); err != nil {
		return errors.Wrap(err, "failed to merge datastore with Keystone contracts addresses")
	}

	for _, selector := range chainSelectors {
		forwarderAddr := contracts.MustGetAddressFromMemoryDataStore(memoryDatastore, selector, keystone_changeset.KeystoneForwarder.String(), contractVersions[keystone_changeset.KeystoneForwarder.String()], "")
		testLogger.Info().Msgf("Deployed EVM Forwarder %s contract on chain %d at %s", contractVersions[keystone_changeset.KeystoneForwarder.String()], selector, forwarderAddr)
	}

	cldfEnv.DataStore = memoryDatastore.Seal()

	return nil
}

func configureEVMForwarders(testLogger zerolog.Logger, cldfEnv *cldf.Environment, chainSelectors []uint64, ocr3DON *cre.Don) (*forwarder.Config, error) {
	forwarderCfg := forwarder.DonConfiguration{
		Name:    ocr3DON.Name,
		ID:      libc.MustSafeUint32FromUint64(ocr3DON.ID),
		F:       ocr3DON.F,
		Version: 1, // TODO this should be dynamic, but we don't have cap reg configured at this point, can we get that version from forwarder contract?
		NodeIDs: ocr3DON.KeystoneDONConfig().NodeIDs,
	}

	if len(chainSelectors) == 0 {
		for _, chain := range cldfEnv.BlockChains.EVMChains() {
			chainSelectors = append(chainSelectors, chain.Selector)
		}
	}

	chainsByQualifier := make(map[string]map[uint64]struct{})
	for _, selector := range chainSelectors {
		refs := cldfEnv.DataStore.Addresses().Filter(
			datastore.AddressRefByChainSelector(selector),
			datastore.AddressRefByType(datastore.ContractType(keystone_changeset.KeystoneForwarder.String())),
		)
		if len(refs) == 0 {
			return nil, fmt.Errorf("failed to resolve deployed forwarder for chain selector %d", selector)
		}

		for _, ref := range refs {
			if chainsByQualifier[ref.Qualifier] == nil {
				chainsByQualifier[ref.Qualifier] = make(map[uint64]struct{})
			}
			chainsByQualifier[ref.Qualifier][selector] = struct{}{}
		}
	}

	qualifiers := make([]string, 0, len(chainsByQualifier))
	for qualifier := range chainsByQualifier {
		qualifiers = append(qualifiers, qualifier)
	}
	sort.Strings(qualifiers)

	var configuredConfig forwarder.Config
	for _, qualifier := range qualifiers {
		fout, err := operations.ExecuteSequence(
			cldfEnv.OperationsBundle,
			forwarder.ConfigureSeq,
			forwarder.ConfigureSeqDeps{
				Env: cldfEnv,
			},
			forwarder.ConfigureSeqInput{
				DON:       forwarderCfg,
				Qualifier: qualifier,
				Chains:    chainsByQualifier[qualifier],
			},
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to configure forwarders with qualifier %q", qualifier)
		}
		configuredConfig = fout.Output.Config
	}

	return &configuredConfig, nil
}

func chainsWithForwarders(blockchains []blockchains.Blockchain, nodeSets []cre.NodeSetWithCapabilityConfigs) map[string][]uint64 {
	chainsWithForwarders := make(map[string][]uint64)

	for _, bcOut := range blockchains {
		for _, nodeSet := range nodeSets {
			if chainSelectors, familyExists := chainsWithForwarders[bcOut.ChainFamily()]; familyExists {
				if slices.Contains(chainSelectors, bcOut.ChainSelector()) {
					continue
				}
			}

			if !bcOut.IsFamily(chainselectors.FamilyEVM) && !bcOut.IsFamily(chainselectors.FamilyTron) {
				continue
			}

			if flags.RequiresForwarderContract(nodeSet.GetCapabilityFlags(), bcOut.ChainID()) {
				if _, exists := chainsWithForwarders[bcOut.ChainFamily()]; !exists {
					chainsWithForwarders[bcOut.ChainFamily()] = []uint64{}
				}
				chainsWithForwarders[bcOut.ChainFamily()] = append(chainsWithForwarders[bcOut.ChainFamily()], bcOut.ChainSelector())
			}
		}
	}

	return chainsWithForwarders
}
