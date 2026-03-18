package aptos

import (
	"bytes"
	"context"
	"encoding/hex"
	stderrors "errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"dario.cat/mergo"
	"github.com/Masterminds/semver/v3"
	aptossdk "github.com/aptos-labs/aptos-go-sdk"
	"github.com/pelletier/go-toml/v2"
	pkgerrors "github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/sethvargo/go-retry"
	"google.golang.org/protobuf/types/known/durationpb"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	aptosplatform "github.com/smartcontractkit/chainlink-aptos/bindings/platform"
	capabilitiespb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs"
	crejobops "github.com/smartcontractkit/chainlink/deployment/cre/jobs/operations"
	jobtypes "github.com/smartcontractkit/chainlink/deployment/cre/jobs/types"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
	aptoschangeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/aptos"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	crecontracts "github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
	credon "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don"
	crejobs "github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs/standardcapability"
	creblockchains "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains"
	aptoschain "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/blockchains/aptos"
	corechainlink "github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
)

const (
	flag                    = cre.WriteAptosCapability
	forwarderContractType   = "AptosForwarder"
	forwarderConfigVersion  = 1
	capabilityVersion       = "1.0.0"
	capabilityLabelPrefix   = "aptos:ChainSelector:"
	specConfigP2PMapKey     = "p2pToTransmitterMap"
	specConfigScheduleKey   = "transmissionSchedule"
	specConfigDeltaStageKey = "deltaStage"
	legacyTransmittersKey   = "aptosTransmitters"
	requestTimeoutKey       = "RequestTimeout"
	deltaStageKey           = "DeltaStage"
	transmissionScheduleKey = "TransmissionSchedule"
	forwarderQualifier      = ""
	zeroForwarderHex        = "0x0000000000000000000000000000000000000000000000000000000000000000"
	aptosConfigTemplate     = `{"chainId":"{{.ChainID}}","network":"aptos","creForwarderAddress":"{{.CREForwarderAddress}}"}`
	defaultWriteDeltaStage  = 500*time.Millisecond + 1*time.Second
	defaultRequestTimeout   = 30 * time.Second
)

var forwarderContractVersion = semver.MustParse("1.0.0")

type Aptos struct{}

type methodConfigSettings struct {
	RequestTimeout       time.Duration
	DeltaStage           time.Duration
	TransmissionSchedule capabilitiespb.TransmissionSchedule
}

func (a *Aptos) Flag() cre.CapabilityFlag {
	return flag
}

func CapabilityLabel(chainSelector uint64) string {
	return capabilityLabelPrefix + strconv.FormatUint(chainSelector, 10)
}

func ForwarderAddress(ds datastore.DataStore, chainSelector uint64) (string, bool) {
	key := datastore.NewAddressRefKey(
		chainSelector,
		datastore.ContractType(forwarderContractType),
		forwarderContractVersion,
		forwarderQualifier,
	)
	ref, err := ds.Addresses().Get(key)
	if err != nil {
		return "", false
	}
	return ref.Address, true
}

func MustForwarderAddress(ds datastore.DataStore, chainSelector uint64) string {
	addr, ok := ForwarderAddress(ds, chainSelector)
	if !ok {
		panic(fmt.Sprintf("missing Aptos forwarder address for chain selector %d", chainSelector))
	}
	return addr
}

func (a *Aptos) PreEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.DonMetadata,
	_ *cre.Topology,
	creEnv *cre.Environment,
) (*cre.PreEnvStartupOutput, error) {
	enabledChainIDs, err := don.MustNodeSet().GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}
	if len(enabledChainIDs) == 0 {
		return nil, nil
	}

	forwardersByChainID := make(map[uint64]string, len(enabledChainIDs))
	for _, chainID := range enabledChainIDs {
		aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if err != nil {
			return nil, err
		}

		forwarderAddress, err := ensureForwarder(ctx, testLogger, creEnv, aptosChain)
		if err != nil {
			return nil, err
		}
		forwardersByChainID[chainID] = forwarderAddress
	}

	if err := patchNodeTOML(don, forwardersByChainID); err != nil {
		return nil, err
	}

	return &cre.PreEnvStartupOutput{}, nil
}

func (a *Aptos) PostDONStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	_ *cre.Dons,
	creEnv *cre.Environment,
) (*cre.PostDONStartupOutput, error) {
	enabledChainIDs, err := don.GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return nil, fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}
	if len(enabledChainIDs) == 0 {
		return nil, nil
	}

	workerNodeIDs, err := workerJDNodeIDs(don)
	if err != nil {
		return nil, err
	}
	nodes, err := deployment.NodeInfo(workerNodeIDs, creEnv.CldfEnvironment.Offchain)
	if err != nil {
		return nil, fmt.Errorf("failed to get Aptos worker node info for DON %q: %w", don.Name, err)
	}

	caps := make([]keystone_changeset.DONCapabilityWithConfig, 0, len(enabledChainIDs))
	ocrConfigs := make(map[string]*ocr3.OracleConfig, len(enabledChainIDs))
	capabilityToExtraSignerFamilies := make(map[string][]string, len(enabledChainIDs))
	for _, chainID := range enabledChainIDs {
		aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if err != nil {
			return nil, err
		}
		selector := aptosChain.ChainSelector()
		labelledName := CapabilityLabel(selector)

		p2pToTransmitterMap, err := p2pToTransmitterMapForChainSelector(nodes, selector)
		if err != nil {
			return nil, fmt.Errorf("failed to collect Aptos p2p transmitter map for capability %s: %w", labelledName, err)
		}
		capabilityConfig, err := cre.ResolveCapabilityConfig(don, flag, cre.ChainCapabilityScope(chainID))
		if err != nil {
			return nil, fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, err)
		}
		capConfig, err := BuildCapabilityConfig(capabilityConfig.Values, p2pToTransmitterMap, false)
		if err != nil {
			return nil, fmt.Errorf("failed to build Aptos capability config for capability %s: %w", labelledName, err)
		}

		caps = append(caps, keystone_changeset.DONCapabilityWithConfig{
			Capability: kcr.CapabilitiesRegistryCapability{
				LabelledName:   labelledName,
				Version:        capabilityVersion,
				CapabilityType: 1,
			},
			Config:             capConfig,
			UseCapRegOCRConfig: true,
		})
		ocrConfigs[labelledName] = crecontracts.DefaultChainCapabilityOCR3Config()
		capabilityToExtraSignerFamilies[labelledName] = []string{chainselectors.FamilyAptos}
	}

	return &cre.PostDONStartupOutput{
		DONCapabilityWithConfig:         caps,
		CapabilityToOCR3Config:          ocrConfigs,
		CapabilityToExtraSignerFamilies: capabilityToExtraSignerFamilies,
	}, nil
}

func (a *Aptos) PostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	specs := make(map[string][]string)

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

	enabledChainIDs, err := nodeSet.GetEnabledChainIDsForCapability(flag)
	if err != nil {
		return fmt.Errorf("could not find enabled chainIDs for '%s' in don '%s': %w", flag, don.Name, err)
	}
	if len(enabledChainIDs) == 0 {
		return nil
	}

	if configureErr := configureForwarders(ctx, testLogger, don, creEnv, enabledChainIDs); configureErr != nil {
		return configureErr
	}

	bootstrapNode, ok := dons.Bootstrap()
	if !ok {
		return pkgerrors.New("bootstrap node not found; required for Aptos OCR bootstrap peers")
	}
	bootstrapPeers := []string{
		fmt.Sprintf("%s@%s:%d", strings.TrimPrefix(bootstrapNode.Keys.PeerID(), "p2p_"), bootstrapNode.Host, cre.OCRPeeringPort),
	}

	capRegSemver, ok := creEnv.ContractVersions[keystone_changeset.CapabilitiesRegistry.String()]
	if !ok {
		return pkgerrors.New("CapabilitiesRegistry version not found in contract versions")
	}
	capRegVersion := capRegSemver.String()

	for _, chainID := range enabledChainIDs {
		aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if err != nil {
			return err
		}

		capabilityConfig, resolveErr := cre.ResolveCapabilityConfig(nodeSet, flag, cre.ChainCapabilityScope(chainID))
		if resolveErr != nil {
			return fmt.Errorf("could not resolve capability config for '%s' on chain %d: %w", flag, chainID, resolveErr)
		}
		command, cErr := standardcapability.GetCommand(capabilityConfig.BinaryName)
		if cErr != nil {
			return pkgerrors.Wrap(cErr, "failed to get command for Aptos capability")
		}

		tmpl, tmplErr := template.New("aptos-config").Parse(aptosConfigTemplate)
		if tmplErr != nil {
			return pkgerrors.Wrap(tmplErr, "failed to parse Aptos config template")
		}

		forwarderAddress := MustForwarderAddress(creEnv.CldfEnvironment.DataStore, aptosChain.ChainSelector())
		templateData := map[string]string{
			"ChainID":             strconv.FormatUint(chainID, 10),
			"CREForwarderAddress": forwarderAddress,
		}

		var configBuffer bytes.Buffer
		if execErr := tmpl.Execute(&configBuffer, templateData); execErr != nil {
			return pkgerrors.Wrap(execErr, "failed to execute Aptos config template")
		}

		configStr := configBuffer.String()
		if validationErr := credon.ValidateTemplateSubstitution(configStr, flag); validationErr != nil {
			return fmt.Errorf("aptos template validation failed: %w\nRendered: %s", validationErr, configStr)
		}

		workerInput := jobs.ProposeJobSpecInput{
			Domain:      offchain.ProductLabel,
			Environment: cre.EnvironmentName,
			DONName:     don.Name,
			JobName:     "write-aptos-worker-" + strconv.FormatUint(chainID, 10),
			ExtraLabels: map[string]string{cre.CapabilityLabelKey: flag},
			DONFilters: []offchain.TargetDONFilter{
				{Key: offchain.FilterKeyDONName, Value: don.Name},
			},
			Template: jobtypes.ReadContract,
			Inputs: jobtypes.JobSpecInput{
				"command":            command,
				"config":             configStr,
				"chainSelectorEVM":   creEnv.RegistryChainSelector,
				"chainSelectorAptos": aptosChain.ChainSelector(),
				"bootstrapPeers":     bootstrapPeers,
				"useCapRegOCRConfig": true,
				"capRegVersion":      capRegVersion,
			},
		}

		proposer := jobs.ProposeJobSpec{}
		if verifyErr := proposer.VerifyPreconditions(*creEnv.CldfEnvironment, workerInput); verifyErr != nil {
			return fmt.Errorf("precondition verification failed for Aptos worker job: %w", verifyErr)
		}
		workerReport, err := proposer.Apply(*creEnv.CldfEnvironment, workerInput)
		if err != nil {
			return fmt.Errorf("failed to propose Aptos worker job spec: %w", err)
		}

		for _, report := range workerReport.Reports {
			out, ok := report.Output.(crejobops.ProposeStandardCapabilityJobOutput)
			if !ok {
				return fmt.Errorf("unable to cast to ProposeStandardCapabilityJobOutput, actual type: %T", report.Output)
			}
			if mergeErr := mergo.Merge(&specs, out.Specs, mergo.WithAppendSlice); mergeErr != nil {
				return fmt.Errorf("failed to merge Aptos worker job specs: %w", mergeErr)
			}
		}
	}

	if len(specs) == 0 {
		return nil
	}
	if err := crejobs.Approve(ctx, creEnv.CldfEnvironment.Offchain, dons, specs); err != nil {
		return fmt.Errorf("failed to approve Aptos jobs: %w", err)
	}
	return nil
}

// BuildCapabilityConfig builds the Aptos capability config using the same config
// split as runtime capabilities: method execution policy in MethodConfigs and
// Aptos-specific runtime inputs in SpecConfig. OCR contract config remains
// separate and is only attached through UseCapRegOCRConfig.
func BuildCapabilityConfig(values map[string]any, p2pToTransmitterMap map[string]string, localOnly bool) (*capabilitiespb.CapabilityConfig, error) {
	methodSettings, err := resolveMethodConfigSettings(values)
	if err != nil {
		return nil, err
	}

	capConfig := &capabilitiespb.CapabilityConfig{
		MethodConfigs: methodConfigs(methodSettings),
		LocalOnly:     localOnly,
	}
	if err := setRuntimeSpecConfig(capConfig, methodSettings, p2pToTransmitterMap); err != nil {
		return nil, err
	}
	return capConfig, nil
}

func methodConfigs(settings methodConfigSettings) map[string]*capabilitiespb.CapabilityMethodConfig {
	return map[string]*capabilitiespb.CapabilityMethodConfig{
		"View": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					TransmissionSchedule:      capabilitiespb.TransmissionSchedule_AllAtOnce,
					RequestTimeout:            durationpb.New(settings.RequestTimeout),
					ServerMaxParallelRequests: 10,
					RequestHasherType:         capabilitiespb.RequestHasherType_Simple,
				},
			},
		},
		"WriteReport": {
			RemoteConfig: &capabilitiespb.CapabilityMethodConfig_RemoteExecutableConfig{
				RemoteExecutableConfig: &capabilitiespb.RemoteExecutableConfig{
					TransmissionSchedule:      settings.TransmissionSchedule,
					DeltaStage:                durationpb.New(settings.DeltaStage),
					RequestTimeout:            durationpb.New(settings.RequestTimeout),
					ServerMaxParallelRequests: 10,
					RequestHasherType:         capabilitiespb.RequestHasherType_WriteReportExcludeSignatures,
				},
			},
		},
	}
}

func resolveMethodConfigSettings(values map[string]any) (methodConfigSettings, error) {
	settings := methodConfigSettings{
		RequestTimeout:       defaultRequestTimeout,
		DeltaStage:           defaultWriteDeltaStage,
		TransmissionSchedule: capabilitiespb.TransmissionSchedule_AllAtOnce,
	}

	if values == nil {
		return settings, nil
	}

	requestTimeout, ok, err := durationValue(values, requestTimeoutKey)
	if err != nil {
		return methodConfigSettings{}, err
	}
	if ok {
		settings.RequestTimeout = requestTimeout
	}

	deltaStage, ok, err := durationValue(values, deltaStageKey)
	if err != nil {
		return methodConfigSettings{}, err
	}
	if ok {
		settings.DeltaStage = deltaStage
	}

	transmissionSchedule, ok, err := transmissionScheduleValue(values, transmissionScheduleKey)
	if err != nil {
		return methodConfigSettings{}, err
	}
	if ok {
		settings.TransmissionSchedule = transmissionSchedule
	}

	return settings, nil
}

func transmissionScheduleValue(values map[string]any, key string) (capabilitiespb.TransmissionSchedule, bool, error) {
	raw, ok := values[key]
	if !ok {
		return 0, false, nil
	}

	schedule, ok := raw.(string)
	if !ok {
		return 0, false, fmt.Errorf("%s must be a string, got %T", key, raw)
	}

	switch strings.TrimSpace(schedule) {
	case "allAtOnce":
		return capabilitiespb.TransmissionSchedule_AllAtOnce, true, nil
	case "oneAtATime":
		return capabilitiespb.TransmissionSchedule_OneAtATime, true, nil
	default:
		return 0, false, fmt.Errorf("%s must be allAtOnce or oneAtATime, got %q", key, schedule)
	}
}

func durationValue(values map[string]any, key string) (time.Duration, bool, error) {
	raw, ok := values[key]
	if !ok {
		return 0, false, nil
	}

	switch v := raw.(type) {
	case string:
		parsed, err := time.ParseDuration(strings.TrimSpace(v))
		if err != nil {
			return 0, false, fmt.Errorf("%s must be a valid duration string: %w", key, err)
		}
		return parsed, true, nil
	case time.Duration:
		return v, true, nil
	default:
		return 0, false, fmt.Errorf("%s must be a duration string, got %T", key, raw)
	}
}

func patchNodeTOML(don *cre.DonMetadata, forwardersByChainID map[uint64]string) error {
	for nodeIndex := range don.MustNodeSet().NodeSpecs {
		currentConfig := don.MustNodeSet().NodeSpecs[nodeIndex].Node.TestConfigOverrides
		if strings.TrimSpace(currentConfig) == "" {
			return fmt.Errorf("missing node config for node index %d in DON %q", nodeIndex, don.Name)
		}

		var typedConfig corechainlink.Config
		if err := toml.Unmarshal([]byte(currentConfig), &typedConfig); err != nil {
			return fmt.Errorf("failed to unmarshal config for node index %d: %w", nodeIndex, err)
		}

		for chainID, forwarderAddress := range forwardersByChainID {
			if err := setForwarderAddress(&typedConfig, strconv.FormatUint(chainID, 10), forwarderAddress); err != nil {
				return fmt.Errorf("failed to patch Aptos forwarder address for node index %d: %w", nodeIndex, err)
			}
		}

		stringifiedConfig, err := toml.Marshal(typedConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal patched config for node index %d: %w", nodeIndex, err)
		}
		don.MustNodeSet().NodeSpecs[nodeIndex].Node.TestConfigOverrides = string(stringifiedConfig)
	}

	return nil
}

func setForwarderAddress(cfg *corechainlink.Config, chainID, forwarderAddress string) error {
	for i := range cfg.Aptos {
		raw := map[string]any(cfg.Aptos[i])
		if fmt.Sprint(raw["ChainID"]) != chainID {
			continue
		}

		workflow := make(map[string]any)
		switch existing := raw["Workflow"].(type) {
		case map[string]any:
			for k, v := range existing {
				workflow[k] = v
			}
		case corechainlink.RawConfig:
			for k, v := range existing {
				workflow[k] = v
			}
		case nil:
		default:
			return fmt.Errorf("unexpected Aptos workflow config type %T", existing)
		}
		workflow["ForwarderAddress"] = forwarderAddress
		raw["Workflow"] = workflow
		cfg.Aptos[i] = corechainlink.RawConfig(raw)
		return nil
	}

	return fmt.Errorf("Aptos chain %s not found in node config", chainID)
}

func ensureForwarder(
	ctx context.Context,
	testLogger zerolog.Logger,
	creEnv *cre.Environment,
	chain *aptoschain.Blockchain,
) (string, error) {
	if addr, ok := ForwarderAddress(creEnv.CldfEnvironment.DataStore, chain.ChainSelector()); ok {
		return addr, nil
	}
	if !creEnv.Provider.IsDocker() {
		return "", fmt.Errorf("missing Aptos forwarder address for chain selector %d", chain.ChainSelector())
	}

	nodeURL, err := chain.NodeURL()
	if err != nil {
		return "", fmt.Errorf("invalid Aptos node URL for chain selector %d: %w", chain.ChainSelector(), err)
	}
	client, err := chain.NodeClient()
	if err != nil {
		return "", fmt.Errorf("failed to create Aptos client for chain selector %d (%s): %w", chain.ChainSelector(), nodeURL, err)
	}
	deployerAccount, err := chain.LocalDeployerAccount()
	if err != nil {
		return "", fmt.Errorf("failed to create Aptos deployer signer: %w", err)
	}
	deploymentChain, err := chain.LocalDeploymentChain()
	if err != nil {
		return "", fmt.Errorf("failed to build Aptos deployment chain for chain selector %d: %w", chain.ChainSelector(), err)
	}

	owner := deployerAccount.AccountAddress()
	containerName := ""
	if output := chain.CtfOutput(); output != nil {
		containerName = output.ContainerName
	}
	if ensureErr := ensureAccountVisible(ctx, testLogger, client, nodeURL, owner, chain.ChainSelector(), containerName); ensureErr != nil {
		testLogger.Warn().
			Uint64("chainSelector", chain.ChainSelector()).
			Str("nodeURL", nodeURL).
			Err(ensureErr).
			Msg("Aptos deployer account not confirmed visible yet; proceeding with deploy retries")
	}

	var deployedAddress string
	var pendingTxHash string
	var lastDeployErr error
	if retryErr := retry.Do(ctx, retry.WithMaxDuration(3*time.Minute, retry.NewFibonacci(500*time.Millisecond)), func(ctx context.Context) error {
		deploymentResp, deployErr := aptoschangeset.DeployPlatform(deploymentChain, owner, nil)
		if deployErr != nil {
			lastDeployErr = deployErr
			if fundErr := chain.Fund(ctx, owner.StringLong(), 1_000_000_000_000); fundErr != nil {
				testLogger.Warn().
					Uint64("chainSelector", chain.ChainSelector()).
					Err(fundErr).
					Msg("failed to re-fund Aptos deployer account during deploy retry")
			}
			return retry.RetryableError(fmt.Errorf("deploy-to-object failed: %w", deployErr))
		}
		if deploymentResp == nil {
			lastDeployErr = pkgerrors.New("nil deployment response")
			return retry.RetryableError(pkgerrors.New("DeployPlatform returned nil response"))
		}
		deployedAddress = deploymentResp.Address.StringLong()
		pendingTxHash = deploymentResp.Tx
		return nil
	}); retryErr != nil {
		if lastDeployErr != nil {
			return "", fmt.Errorf("failed to deploy Aptos platform forwarder for chain selector %d after retries: %w", chain.ChainSelector(), stderrors.Join(lastDeployErr, retryErr))
		}
		return "", fmt.Errorf("failed to deploy Aptos platform forwarder for chain selector %d after retries: %w", chain.ChainSelector(), retryErr)
	}

	addr, err := normalizeForwarderAddress(deployedAddress)
	if err != nil {
		return "", fmt.Errorf("invalid Aptos forwarder address parsed from deployment output for chain selector %d: %w", chain.ChainSelector(), err)
	}

	if err := addForwarderToDataStore(creEnv, chain.ChainSelector(), addr); err != nil {
		return "", err
	}

	testLogger.Info().
		Uint64("chainSelector", chain.ChainSelector()).
		Str("nodeURL", nodeURL).
		Str("txHash", pendingTxHash).
		Str("forwarderAddress", addr).
		Msg("Aptos platform forwarder deployed")

	return addr, nil
}

func addForwarderToDataStore(creEnv *cre.Environment, chainSelector uint64, address string) error {
	memoryDatastore, err := crecontracts.NewDataStoreFromExisting(creEnv.CldfEnvironment.DataStore)
	if err != nil {
		return fmt.Errorf("failed to create memory datastore: %w", err)
	}

	err = memoryDatastore.AddressRefStore.Add(datastore.AddressRef{
		Address:       address,
		ChainSelector: chainSelector,
		Type:          datastore.ContractType(forwarderContractType),
		Version:       forwarderContractVersion,
		Qualifier:     forwarderQualifier,
	})
	if err != nil && !stderrors.Is(err, datastore.ErrAddressRefExists) {
		return fmt.Errorf("failed to add Aptos forwarder address to datastore: %w", err)
	}

	creEnv.CldfEnvironment.DataStore = memoryDatastore.Seal()
	return nil
}

func configureForwarders(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	creEnv *cre.Environment,
	chainIDs []uint64,
) error {
	workers, err := don.Workers()
	if err != nil {
		return fmt.Errorf("failed to get worker nodes for DON %q: %w", don.Name, err)
	}
	f := (len(workers) - 1) / 3
	if f <= 0 {
		return fmt.Errorf("invalid Aptos DON %q fault tolerance F=%d (workers=%d)", don.Name, f, len(workers))
	}
	if f > 255 {
		return fmt.Errorf("aptos DON %q fault tolerance F=%d exceeds u8", don.Name, f)
	}

	donIDUint32, err := aptosDonIDUint32(don.ID)
	if err != nil {
		return fmt.Errorf("invalid DON id for Aptos forwarder config: %w", err)
	}

	oracles, err := donOraclePublicKeys(ctx, don)
	if err != nil {
		return err
	}

	for _, chainID := range chainIDs {
		aptosChain, err := findAptosChainByChainID(creEnv.Blockchains, chainID)
		if err != nil {
			return err
		}

		nodeURL, err := aptosChain.NodeURL()
		if err != nil {
			return fmt.Errorf("invalid Aptos node URL for chain selector %d: %w", aptosChain.ChainSelector(), err)
		}
		client, err := aptosChain.NodeClient()
		if err != nil {
			return fmt.Errorf("failed to create Aptos client for chain selector %d (%s): %w", aptosChain.ChainSelector(), nodeURL, err)
		}
		deployerAccount, err := aptosChain.LocalDeployerAccount()
		if err != nil {
			return fmt.Errorf("failed to create Aptos deployer signer for forwarder config: %w", err)
		}
		deployerAddress := deployerAccount.AccountAddress()

		containerName := ""
		if output := aptosChain.CtfOutput(); output != nil {
			containerName = output.ContainerName
		}

		if err := ensureAccountVisible(ctx, testLogger, client, nodeURL, deployerAddress, aptosChain.ChainSelector(), containerName); err != nil {
			testLogger.Warn().
				Uint64("chainSelector", aptosChain.ChainSelector()).
				Str("nodeURL", nodeURL).
				Err(err).
				Msg("Aptos deployer account not confirmed visible yet; proceeding with forwarder set_config retries")
		}

		forwarderHex := MustForwarderAddress(creEnv.CldfEnvironment.DataStore, aptosChain.ChainSelector())
		var forwarderAddr aptossdk.AccountAddress
		if err := forwarderAddr.ParseStringRelaxed(forwarderHex); err != nil {
			return fmt.Errorf("invalid Aptos forwarder address for chain selector %d: %w", aptosChain.ChainSelector(), err)
		}
		forwarderContract := aptosplatform.Bind(forwarderAddr, client).Forwarder()

		var pendingTxHash string
		var lastSetConfigErr error
		if err := retry.Do(ctx, retry.WithMaxDuration(2*time.Minute, retry.NewFibonacci(500*time.Millisecond)), func(ctx context.Context) error {
			pendingTx, err := forwarderContract.SetConfig(&bind.TransactOpts{Signer: deployerAccount}, donIDUint32, forwarderConfigVersion, byte(f), oracles)
			if err != nil {
				lastSetConfigErr = err
				if fundErr := aptosChain.Fund(ctx, deployerAddress.StringLong(), 1_000_000_000_000); fundErr != nil {
					testLogger.Warn().
						Uint64("chainSelector", aptosChain.ChainSelector()).
						Err(fundErr).
						Msg("failed to fund Aptos deployer account during set_config retry")
				}
				return retry.RetryableError(fmt.Errorf("set_config transaction submit failed: %w", err))
			}
			pendingTxHash = pendingTx.Hash
			receipt, err := client.WaitForTransaction(pendingTxHash)
			if err != nil {
				lastSetConfigErr = err
				return retry.RetryableError(fmt.Errorf("waiting for set_config transaction failed: %w", err))
			}
			if !receipt.Success {
				lastSetConfigErr = fmt.Errorf("vm status: %s", receipt.VmStatus)
				return retry.RetryableError(fmt.Errorf("set_config transaction failed: %s", receipt.VmStatus))
			}
			return nil
		}); err != nil {
			if lastSetConfigErr != nil {
				return fmt.Errorf("failed to configure Aptos forwarder %s for DON %q on chain selector %d: %w", forwarderHex, don.Name, aptosChain.ChainSelector(), stderrors.Join(lastSetConfigErr, err))
			}
			return fmt.Errorf("failed to configure Aptos forwarder %s for DON %q on chain selector %d: %w", forwarderHex, don.Name, aptosChain.ChainSelector(), err)
		}

		testLogger.Info().
			Str("donName", don.Name).
			Uint64("donID", don.ID).
			Uint64("chainSelector", aptosChain.ChainSelector()).
			Str("txHash", pendingTxHash).
			Str("forwarderAddress", forwarderHex).
			Msg("configured Aptos forwarder set_config")
	}

	return nil
}

func donOraclePublicKeys(ctx context.Context, don *cre.Don) ([][]byte, error) {
	workers, err := don.Workers()
	if err != nil {
		return nil, fmt.Errorf("failed to list worker nodes for DON %q: %w", don.Name, err)
	}

	oracles := make([][]byte, 0, len(workers))
	for _, worker := range workers {
		ocr2ID := ""
		if worker.Keys != nil && worker.Keys.OCR2BundleIDs != nil {
			ocr2ID = worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos]
		}
		if ocr2ID == "" {
			fetchedID, err := worker.Clients.GQLClient.FetchOCR2KeyBundleID(ctx, strings.ToUpper(chainselectors.FamilyAptos))
			if err != nil {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q and fallback fetch failed: %w", worker.Name, don.Name, err)
			}
			if fetchedID == "" {
				return nil, fmt.Errorf("missing Aptos OCR2 bundle id for worker %q in DON %q", worker.Name, don.Name)
			}
			ocr2ID = fetchedID
			if worker.Keys != nil {
				if worker.Keys.OCR2BundleIDs == nil {
					worker.Keys.OCR2BundleIDs = make(map[string]string)
				}
				worker.Keys.OCR2BundleIDs[chainselectors.FamilyAptos] = ocr2ID
			}
		}

		exported, err := worker.ExportOCR2Keys(ocr2ID)
		if err != nil {
			return nil, fmt.Errorf("failed to export Aptos OCR2 key for worker %q (bundle %s): %w", worker.Name, ocr2ID, err)
		}
		pubkey, err := parseOCR2OnchainPublicKey(exported.OnchainPublicKey)
		if err != nil {
			return nil, fmt.Errorf("invalid Aptos OCR2 onchain public key for worker %q: %w", worker.Name, err)
		}
		oracles = append(oracles, pubkey)
	}

	return oracles, nil
}

func workerJDNodeIDs(don *cre.Don) ([]string, error) {
	workers, err := don.Workers()
	if err != nil {
		return nil, err
	}

	nodeIDs := make([]string, 0, len(workers))
	for _, worker := range workers {
		if worker.JobDistributorDetails.NodeID == "" {
			return nil, fmt.Errorf("missing Job Distributor node id for worker %q in DON %q", worker.Name, don.Name)
		}
		nodeIDs = append(nodeIDs, worker.JobDistributorDetails.NodeID)
	}

	return nodeIDs, nil
}

func p2pToTransmitterMapForChainSelector(nodes deployment.Nodes, chainSelector uint64) (map[string]string, error) {
	if len(nodes) == 0 {
		return nil, pkgerrors.New("no DON worker nodes provided")
	}

	p2pToTransmitterMap := make(map[string]string)
	for _, node := range nodes {
		ocrCfg, ok := node.OCRConfigForChainSelector(chainSelector)
		if !ok {
			continue
		}

		transmitter, err := normalizeTransmitter(string(ocrCfg.TransmitAccount))
		if err != nil {
			return nil, fmt.Errorf("invalid Aptos transmitter for node %s: %w", node.Name, err)
		}

		peerKey := hex.EncodeToString(node.PeerID[:])
		p2pToTransmitterMap[peerKey] = transmitter
	}

	if len(p2pToTransmitterMap) == 0 {
		return nil, fmt.Errorf("no Aptos OCR config/transmitters found for chain selector %d", chainSelector)
	}

	return p2pToTransmitterMap, nil
}

func setRuntimeSpecConfig(capConfig *capabilitiespb.CapabilityConfig, settings methodConfigSettings, p2pToTransmitterMap map[string]string) error {
	if capConfig == nil {
		return pkgerrors.New("capability config is nil")
	}

	specConfig, err := values.FromMapValueProto(capConfig.SpecConfig)
	if err != nil {
		return fmt.Errorf("failed to decode existing spec config: %w", err)
	}
	if specConfig == nil {
		specConfig = values.EmptyMap()
	}

	delete(specConfig.Underlying, legacyTransmittersKey)

	scheduleValue, err := values.Wrap(remoteTransmissionScheduleString(settings.TransmissionSchedule))
	if err != nil {
		return fmt.Errorf("failed to wrap transmission schedule: %w", err)
	}
	specConfig.Underlying[specConfigScheduleKey] = scheduleValue

	deltaStageValue, err := values.Wrap(settings.DeltaStage)
	if err != nil {
		return fmt.Errorf("failed to wrap delta stage: %w", err)
	}
	specConfig.Underlying[specConfigDeltaStageKey] = deltaStageValue

	if len(p2pToTransmitterMap) > 0 {
		mapValue, err := values.Wrap(p2pToTransmitterMap)
		if err != nil {
			return fmt.Errorf("failed to wrap p2p transmitter map: %w", err)
		}
		specConfig.Underlying[specConfigP2PMapKey] = mapValue
	}

	capConfig.SpecConfig = values.ProtoMap(specConfig)
	return nil
}

func remoteTransmissionScheduleString(schedule capabilitiespb.TransmissionSchedule) string {
	switch schedule {
	case capabilitiespb.TransmissionSchedule_OneAtATime:
		return "oneAtATime"
	default:
		return "allAtOnce"
	}
}

func normalizeTransmitter(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", pkgerrors.New("empty Aptos transmitter")
	}

	var addr aptossdk.AccountAddress
	if err := addr.ParseStringRelaxed(s); err != nil {
		return "", err
	}
	return addr.StringLong(), nil
}

func normalizeForwarderAddress(raw string) (string, error) {
	var addr aptossdk.AccountAddress
	if err := addr.ParseStringRelaxed(strings.TrimSpace(raw)); err != nil {
		return "", err
	}
	return addr.StringLong(), nil
}

func findAptosChainByChainID(chains []creblockchains.Blockchain, chainID uint64) (*aptoschain.Blockchain, error) {
	for _, bc := range chains {
		if bc.IsFamily(chainselectors.FamilyAptos) && bc.ChainID() == chainID {
			aptosBlockchain, ok := bc.(*aptoschain.Blockchain)
			if !ok {
				return nil, fmt.Errorf("Aptos blockchain for chain id %d has unexpected type %T", chainID, bc)
			}
			return aptosBlockchain, nil
		}
	}
	return nil, fmt.Errorf("Aptos blockchain for chain id %d not found", chainID)
}

func normalizeNodeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("aptos node URL %q must include scheme and host", rawURL)
	}
	path := strings.TrimRight(u.Path, "/")
	if path == "" || path != "/v1" {
		u.Path = "/v1"
	}
	return u.String(), nil
}

func NormalizeNodeURL(rawURL string) (string, error) {
	return normalizeNodeURL(rawURL)
}

func faucetURLFromNodeURL(nodeURL string) (string, error) {
	parsed, err := url.Parse(nodeURL)
	if err != nil {
		return "", err
	}
	host := parsed.Hostname()
	if host == "" {
		return "", fmt.Errorf("aptos node URL %q has empty host", nodeURL)
	}
	parsed.Host = fmt.Sprintf("%s:%s", host, blockchain.DefaultAptosFaucetPort)
	parsed.Path = ""
	parsed.RawPath = ""
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func FaucetURLFromNodeURL(nodeURL string) (string, error) {
	return faucetURLFromNodeURL(nodeURL)
}

func aptosChainIDUint8(chainID uint64) (uint8, error) {
	if chainID > 255 {
		return 0, fmt.Errorf("aptos chain id %d does not fit in uint8", chainID)
	}
	return uint8(chainID), nil
}

func ChainIDUint8(chainID uint64) (uint8, error) {
	return aptosChainIDUint8(chainID)
}

func aptosDonIDUint32(donID uint64) (uint32, error) {
	if donID > uint64(^uint32(0)) {
		return 0, fmt.Errorf("don id %d exceeds u32", donID)
	}
	return uint32(donID), nil
}

func ensureAccountVisible(
	ctx context.Context,
	testLogger zerolog.Logger,
	client *aptossdk.NodeClient,
	nodeURL string,
	address aptossdk.AccountAddress,
	chainSelector uint64,
	containerName string,
) error {
	if err := FundAccountBestEffort(ctx, testLogger, client, nodeURL, containerName, address, 0, 100_000_000); err == nil {
		return nil
	}

	testLogger.Warn().
		Uint64("chainSelector", chainSelector).
		Str("nodeURL", nodeURL).
		Str("account", address.StringLong()).
		Msg("Aptos account not confirmed visible after funding attempts")
	return fmt.Errorf("aptos account %s not visible yet on %s", address.StringLong(), nodeURL)
}

func waitForAccountVisible(ctx context.Context, client *aptossdk.NodeClient, address aptossdk.AccountAddress) error {
	var lastErr error
	err := retry.Do(ctx, retry.WithMaxDuration(20*time.Second, retry.NewFibonacci(250*time.Millisecond)), func(context.Context) error {
		if _, err := client.Account(address); err != nil {
			lastErr = err
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		if lastErr != nil {
			return lastErr
		}
		return err
	}
	return nil
}

func WaitForAccountVisible(ctx context.Context, client *aptossdk.NodeClient, address aptossdk.AccountAddress, timeout time.Duration) error {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return waitForAccountVisible(waitCtx, client, address)
}

func fundAccountInContainer(ctx context.Context, containerName, address string, amount uint64) error {
	dc, err := framework.NewDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client: %w", err)
	}
	cmd := []string{
		"aptos", "account", "fund-with-faucet",
		"--account", address,
		"--amount", strconv.FormatUint(amount, 10),
	}
	if _, err = dc.ExecContainerWithContext(ctx, containerName, cmd); err != nil {
		return fmt.Errorf("failed to execute aptos faucet funding command in container %s: %w", containerName, err)
	}
	return nil
}

func FundAccountInContainer(ctx context.Context, containerName, address string, amount uint64) error {
	return fundAccountInContainer(ctx, containerName, address, amount)
}

func FundAccountBestEffort(
	ctx context.Context,
	testLogger zerolog.Logger,
	client *aptossdk.NodeClient,
	nodeURL string,
	containerName string,
	address aptossdk.AccountAddress,
	minBalance uint64,
	fundingAmount uint64,
) error {
	if _, err := client.Account(address); err == nil {
		if minBalance == 0 {
			return nil
		}
		if balance, balErr := client.AccountAPTBalance(address); balErr == nil && balance >= minBalance {
			return nil
		}
	}

	if faucetURL, err := faucetURLFromNodeURL(nodeURL); err == nil {
		if faucetClient, err := aptossdk.NewFaucetClient(client, faucetURL); err == nil {
			if fundErr := faucetClient.Fund(address, fundingAmount); fundErr != nil {
				testLogger.Warn().
					Err(fundErr).
					Str("faucet_url", faucetURL).
					Str("account", address.StringLong()).
					Msg("Aptos host faucet funding failed")
			} else if waitErr := WaitForAccountVisible(ctx, client, address, 8*time.Second); waitErr == nil {
				return nil
			}
		}
	}

	if strings.TrimSpace(containerName) == "" {
		return fmt.Errorf("aptos account %s not visible after host funding attempts", address.StringLong())
	}
	if err := fundAccountInContainer(ctx, containerName, address.StringLong(), fundingAmount); err != nil {
		return err
	}
	return WaitForAccountVisible(ctx, client, address, 8*time.Second)
}

func WaitForTransactionSuccess(client *aptossdk.NodeClient, txHash, label string) error {
	tx, err := client.WaitForTransaction(txHash)
	if err != nil {
		return fmt.Errorf("failed waiting for Aptos tx %s: %w", label, err)
	}
	if !tx.Success {
		return fmt.Errorf("aptos tx failed: %s vm_status=%s", label, tx.VmStatus)
	}
	return nil
}

func parseOCR2OnchainPublicKey(hexValue string) ([]byte, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(hexValue), "ocr2on_aptos_")
	decoded, err := hex.DecodeString(trimmed)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

var (
	_ cre.Feature = (*Aptos)(nil)
)
