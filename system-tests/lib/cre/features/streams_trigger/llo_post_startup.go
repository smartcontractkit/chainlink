package streams_trigger

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"text/template"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	ds "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	channel_config_store "github.com/smartcontractkit/data-streams-deploy/changeset/channel-config-store"
	configurator_v0_5_0 "github.com/smartcontractkit/data-streams-deploy/changeset/configurator/v0_5_0"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
)

// Static channel definition for LLO E2E: single definition with DON ID and optional S3.
// Channel defs are static data; only the ChannelConfigStore address comes from deployment.
const (
	staticChannelDefS3URL = ""
)

var staticChannelDefHash = [32]byte{}

// Stream job template: hardcoded values per stream ID to avoid bridge dependency.
// Stream 1=424242 (Format 5), 2=3000, 3=15, 4=111111 (Format 7 calculated base).
const streamJobSpecTmpl = `
type = "stream"
schemaVersion = 1
name = "{{.Name}}"
streamID = {{.StreamID}}
externalJobID = "{{.ExternalJobID}}"
observationSource = """
    result    [type=memo value="{{.HardcodedValue}}"];
    multiply  [type=multiply times=1 index=0];
    result    -> multiply;
"""
`

// LLO bootstrap job spec template.
const lloBootstrapJobSpecTmpl = `type = "bootstrap"
schemaVersion = 1
name = "llo-bootstrap"
externalJobID = "{{.ExternalJobID}}"
contractID = "{{.ConfiguratorAddr}}"
contractConfigTrackerPollInterval = "1s"
relay = "evm"

[relayConfig]
chainID = "{{.ChainID}}"
fromBlock = 0
lloDonID = {{.DonID}}
lloConfigMode = "bluegreen"
providerType = "llo"
`

// LLO worker job spec template with inline channel definitions (static JSON).
const lloWorkerJobSpecTmpl = `type = "offchainreporting2"
schemaVersion = 1
name = "llo-streams-don"
externalJobID = "{{.ExternalJobID}}"
forwardingAllowed = false
maxTaskDuration = "1s"
contractID = "{{.ConfiguratorAddr}}"
contractConfigTrackerPollInterval = "1s"
ocrKeyBundleID = "{{.OCRKeyBundleID}}"
p2pv2Bootstrappers = ["{{.BootstrapAddr}}"]
relay = "evm"
pluginType = "llo"
transmitterID = "{{.TransmitterID}}"

[pluginConfig]
channelDefinitions = """
{
  "1": {
    "reportFormat": 5,
    "streams": [{"streamId": 1, "aggregator": "median"}],
    "opts": {}
  },
  "2": {
    "reportFormat": 7,
    "streams": [
      {"streamId": 2, "aggregator": "median"},
      {"streamId": 3, "aggregator": "median"},
      {"streamId": 4, "aggregator": "median"},
      {"streamId": 5, "aggregator": "calculated"}
    ],
    "opts": {
      "feedId": "0x0001000000000000000000000000000000000000000000000000000000000001",
      "baseUSDFee": "0.1",
      "expirationWindow": 3600,
      "abi": [
        {"type": "int192", "expression": "Mul(s4, 5)", "expressionStreamId": 5}
      ]
    }
  }
}
"""
donID = {{.DonID}}

[[pluginConfig.transmitters]]
type = "cre"
[pluginConfig.transmitters.opts]
triggerCapabilityName = "streams-trigger"
triggerCapabilityVersion = "2.0.0"
triggerTickerMinResolutionMs = 1000
triggerSendChannelBufferSize = 1000
transmissionWindowMs = 100

[relayConfig]
chainID = "{{.ChainID}}"
lloDonID = {{.DonID}}
lloConfigMode = "bluegreen"
fromBlock = 1
`

func runPostEnvStartup(
	ctx context.Context,
	testLogger zerolog.Logger,
	don *cre.Don,
	dons *cre.Dons,
	creEnv *cre.Environment,
) error {
	if creEnv.CldfEnvironment == nil || len(creEnv.Blockchains) == 0 {
		testLogger.Debug().Msg("streams-trigger PostEnvStartup: no CldfEnvironment or Blockchains, skipping LLO deployment")
		return nil
	}

	chainSelector := creEnv.Blockchains[0].ChainSelector()
	chainID := creEnv.Blockchains[0].ChainID()
	chainIDStr := fmt.Sprintf("%d", chainID)
	donID := uint32(don.ID)
	cldfEnv := creEnv.CldfEnvironment

	testLogger.Info().
		Uint64("chainSelector", chainSelector).
		Uint64("chainID", chainID).
		Uint32("donID", donID).
		Msg("Deploying LLO infrastructure via PostEnvStartup (changesets + static channel defs)")

	// 1. Deploy Configurator
	configuratorOutput, err := configurator_v0_5_0.DeployConfiguratorChangeset.Apply(*cldfEnv, configurator_v0_5_0.DeployConfiguratorConfig{
		ChainsToDeploy: []uint64{chainSelector},
	})
	if err != nil {
		return fmt.Errorf("deploy Configurator: %w", err)
	}
	mergeCldfDatastore(cldfEnv, configuratorOutput.DataStore, testLogger)

	configuratorAddr, err := extractContractAddress(cldfEnv.DataStore, chainSelector, "Configurator")
	if err != nil {
		return fmt.Errorf("get Configurator address: %w", err)
	}
	testLogger.Info().Str("configurator", configuratorAddr.Hex()).Msg("Configurator deployed")

	// 2. Deploy ChannelConfigStore
	ccsOutput, err := channel_config_store.DeployChannelConfigStoreChangeset.Apply(*cldfEnv, channel_config_store.DeployChannelConfigStoreConfig{
		ChainsToDeploy: []uint64{chainSelector},
	})
	if err != nil {
		return fmt.Errorf("deploy ChannelConfigStore: %w", err)
	}
	mergeCldfDatastore(cldfEnv, ccsOutput.DataStore, testLogger)

	ccsAddr, err := extractContractAddress(cldfEnv.DataStore, chainSelector, "ChannelConfigStore")
	if err != nil {
		return fmt.Errorf("get ChannelConfigStore address: %w", err)
	}
	testLogger.Info().Str("channelConfigStore", ccsAddr.Hex()).Msg("ChannelConfigStore deployed")

	// 3. Set channel definitions (static data)
	_, err = channel_config_store.SetChannelDefinitionChangeset.Apply(*cldfEnv, channel_config_store.SetChannelDefinitionsConfig{
		DefinitionsByChain: map[uint64][]channel_config_store.ChannelDefinition{
			chainSelector: {{
				ChannelConfigStore: ccsAddr,
				DonID:              donID,
				S3URL:              staticChannelDefS3URL,
				Hash:               staticChannelDefHash,
			}},
		},
	})
	if err != nil {
		return fmt.Errorf("set channel definitions: %w", err)
	}
	testLogger.Info().Msg("Channel definitions set (static)")

	// 4. Deploy stream jobs (text/template)
	workers, err := don.Workers()
	if err != nil {
		return fmt.Errorf("don workers: %w", err)
	}
	streamJobs := []struct{ streamID uint32; name string; hardcodedValue int64 }{
		{1, "stream-test-usd", 424242},
		{2, "stream-native-usd", 3000},
		{3, "stream-link-usd", 15},
		{4, "stream-data-usd", 111111},
	}
	streamTmpl := template.Must(template.New("stream").Parse(streamJobSpecTmpl))
	var jobSpecs cre.DonJobs
	for _, node := range workers {
		for _, sj := range streamJobs {
			var buf bytes.Buffer
			err := streamTmpl.Execute(&buf, map[string]interface{}{
				"Name":           fmt.Sprintf("%s-%s", sj.name, node.Name),
				"StreamID":       sj.streamID,
				"ExternalJobID":  uuid.New().String(),
				"HardcodedValue": sj.hardcodedValue,
			})
			if err != nil {
				return fmt.Errorf("stream job template: %w", err)
			}
			jobSpecs = append(jobSpecs, &jobv1.ProposeJobRequest{
				NodeId: node.JobDistributorDetails.NodeID,
				Spec:   buf.String(),
			})
		}
	}
	if err := jobs.Create(ctx, cldfEnv.Offchain, dons, jobSpecs); err != nil {
		return fmt.Errorf("create stream jobs: %w", err)
	}
	testLogger.Info().Int("count", len(jobSpecs)).Msg("Stream jobs deployed")

	// 5. Deploy LLO bootstrap + worker jobs
	bootstrap, ok := dons.Bootstrap()
	if !ok {
		return fmt.Errorf("bootstrap node not found")
	}
	_, ocrPeeringCfg, err := cre.PeeringCfgs(bootstrap)
	if err != nil {
		return fmt.Errorf("peering config: %w", err)
	}
	bootstrapAddr := fmt.Sprintf("%s@%s:%d", ocrPeeringCfg.OCRBootstraperPeerID, ocrPeeringCfg.OCRBootstraperHost, ocrPeeringCfg.Port)

	var lloJobSpecs cre.DonJobs
	bootstrapTmpl := template.Must(template.New("llo-bootstrap").Parse(lloBootstrapJobSpecTmpl))
	var bootstrapBuf bytes.Buffer
	if err := bootstrapTmpl.Execute(&bootstrapBuf, map[string]interface{}{
		"ExternalJobID":   uuid.New().String(),
		"ConfiguratorAddr": configuratorAddr.Hex(),
		"ChainID":         chainIDStr,
		"DonID":           donID,
	}); err != nil {
		return fmt.Errorf("llo bootstrap template: %w", err)
	}
	lloJobSpecs = append(lloJobSpecs, &jobv1.ProposeJobRequest{
		NodeId: bootstrap.JobDistributorDetails.NodeID,
		Spec:   bootstrapBuf.String(),
	})

	workerTmpl := template.Must(template.New("llo-worker").Parse(lloWorkerJobSpecTmpl))
	for _, node := range workers {
		ocrKeyBundleID := node.Keys.OCR2BundleIDs["evm"]
		csaPubKeyHex := strings.TrimPrefix(node.Keys.CSAKey.Key, "csa_")
		var buf bytes.Buffer
		if err := workerTmpl.Execute(&buf, map[string]interface{}{
			"ExternalJobID":    uuid.New().String(),
			"ConfiguratorAddr": configuratorAddr.Hex(),
			"OCRKeyBundleID":   ocrKeyBundleID,
			"BootstrapAddr":    bootstrapAddr,
			"TransmitterID":    csaPubKeyHex,
			"DonID":            donID,
			"ChainID":          chainIDStr,
		}); err != nil {
			return fmt.Errorf("llo worker template: %w", err)
		}
		lloJobSpecs = append(lloJobSpecs, &jobv1.ProposeJobRequest{
			NodeId: node.JobDistributorDetails.NodeID,
			Spec:   buf.String(),
		})
	}
	if err := jobs.Create(ctx, cldfEnv.Offchain, dons, lloJobSpecs); err != nil {
		return fmt.Errorf("create LLO jobs: %w", err)
	}
	testLogger.Info().Msg("LLO jobs (bootstrap + workers) deployed")

	// 6. Wait for LogPoller to register
	time.Sleep(3 * time.Second)

	// 7. Build oracle identities and set OCR config
	oracles, err := buildOracles(testLogger, workers)
	if err != nil {
		return fmt.Errorf("build oracles: %w", err)
	}
	n := len(oracles)
	f := (n - 1) / 3
	if f < 1 {
		f = 1
	}
	var donConfigIDBytes [32]byte
	big.NewInt(int64(donID)).FillBytes(donConfigIDBytes[:])
	donConfigID := "0x" + hex.EncodeToString(donConfigIDBytes[:])
	configParams := configurator_v0_5_0.NewConfiguratorConfig(configurator_v0_5_0.ConfiguratorSetParamsOptions{
		DONConfigID:         &donConfigID,
		ConfiguratorAddress: &configuratorAddr,
		OCROptions: &configurator_v0_5_0.OCR3DataStreamsOptions{
			Oracles: oracles,
			F:       &f,
			S:       []int{n},
			OnchainConfigOptions: &datastreamsllo.OnchainConfig{
				Version:                 1,
				PredecessorConfigDigest: nil,
			},
			OffchainConfigOptions: datastreamsllo.OffchainConfig{
				ProtocolVersion:                     1,
				DefaultMinReportIntervalNanoseconds: uint64(1 * time.Second),
				EnableObservationCompression:        true,
			},
		},
	})
	_, err = configurator_v0_5_0.SetProductionConfigChangeset.Apply(*cldfEnv, configurator_v0_5_0.SetProductionConfig{
		ConfigurationsByChain: map[uint64][]configurator_v0_5_0.ConfiguratorSetParams{
			chainSelector: {*configParams},
		},
	})
	if err != nil {
		return fmt.Errorf("set production config: %w", err)
	}
	testLogger.Info().Msg("OCR production config set")

	return nil
}

func mergeCldfDatastore(cldfEnv *cldf.Environment, from interface{ Seal() ds.DataStore }, testLogger zerolog.Logger) {
	if from == nil {
		return
	}
	mergedDS := ds.NewMemoryDataStore()
	if err := mergedDS.Merge(from.Seal()); err != nil {
		testLogger.Warn().Err(err).Msg("merge changeset datastore (first)")
		return
	}
	if cldfEnv.DataStore != nil {
		if err := mergedDS.Merge(cldfEnv.DataStore); err != nil {
			testLogger.Warn().Err(err).Msg("merge existing cldf datastore")
			return
		}
	}
	cldfEnv.DataStore = mergedDS.Seal()
}

func extractContractAddress(dataStore ds.DataStore, chainSelector uint64, contractType string) (common.Address, error) {
	if dataStore == nil {
		return common.Address{}, fmt.Errorf("datastore is nil")
	}
	addresses, err := dataStore.Addresses().Fetch()
	if err != nil {
		return common.Address{}, fmt.Errorf("fetch addresses: %w", err)
	}
	for _, addr := range addresses {
		if addr.ChainSelector == chainSelector && string(addr.Type) == contractType {
			return common.HexToAddress(addr.Address), nil
		}
	}
	return common.Address{}, fmt.Errorf("contract %s not found for chain %d", contractType, chainSelector)
}

func buildOracles(testLogger zerolog.Logger, workers []*cre.Node) ([]confighelper.OracleIdentityExtra, error) {
	oracles := make([]confighelper.OracleIdentityExtra, len(workers))
	for i, node := range workers {
		if node.Keys.CSAKey == nil {
			return nil, fmt.Errorf("node %s has no CSA key", node.Name)
		}
		csaKeyHex := strings.TrimPrefix(node.Keys.CSAKey.Key, "csa_")
		csaBytes, err := hex.DecodeString(csaKeyHex)
		if err != nil {
			return nil, fmt.Errorf("decode CSA key node %s: %w", node.Name, err)
		}
		peerID := strings.TrimPrefix(node.Keys.P2PKey.PeerID.String(), "p2p_")

		ocr2Keys, err := node.Clients.RestClient.MustReadOCR2Keys()
		if err != nil {
			return nil, fmt.Errorf("read OCR2 keys node %s: %w", node.Name, err)
		}
		var offchainPK types.OffchainPublicKey
		var configPK types.ConfigEncryptionPublicKey
		var onchainPK []byte
		found := false
		for _, keyData := range ocr2Keys.Data {
			if keyData.Attributes.ChainType == "evm" {
				offchainHex := strings.TrimPrefix(keyData.Attributes.OffChainPublicKey, "ocr2off_evm_")
				configHex := strings.TrimPrefix(keyData.Attributes.ConfigPublicKey, "ocr2cfg_evm_")
				onchainHex := strings.TrimPrefix(keyData.Attributes.OnChainPublicKey, "ocr2on_evm_")
				ob, err := hex.DecodeString(offchainHex)
				if err != nil {
					return nil, errors.Wrapf(err, "offchain key node %s", node.Name)
				}
				copy(offchainPK[:], ob)
				cb, err := hex.DecodeString(configHex)
				if err != nil {
					return nil, errors.Wrapf(err, "config key node %s", node.Name)
				}
				copy(configPK[:], cb)
				onchainPK, err = hex.DecodeString(onchainHex)
				if err != nil {
					return nil, errors.Wrapf(err, "onchain key node %s", node.Name)
				}
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("no EVM OCR2 key for node %s", node.Name)
		}
		csaFormatted := "0x" + hex.EncodeToString(csaBytes)
		oracles[i] = confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPK,
				OffchainPublicKey: offchainPK,
				PeerID:            peerID,
				TransmitAccount:   types.Account(csaFormatted),
			},
			ConfigEncryptionPublicKey: configPK,
		}
	}
	return oracles, nil
}
