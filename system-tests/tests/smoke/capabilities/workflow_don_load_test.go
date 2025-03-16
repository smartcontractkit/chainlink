package capabilities

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/google/uuid"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	types2 "github.com/smartcontractkit/chainlink-common/pkg/types"
	v3 "github.com/smartcontractkit/chainlink-common/pkg/types/mercury/v3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	ns "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	"github.com/smartcontractkit/chainlink/deployment/environment/nodeclient"
	lidebug "github.com/smartcontractkit/chainlink/system-tests/lib/cre/debug"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/flags"
	keystonetypes "github.com/smartcontractkit/chainlink/system-tests/lib/cre/types"
	"github.com/smartcontractkit/chainlink/system-tests/lib/mock_capability"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/targets"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func TestKeystoneWithOCR3Workflow_TwoDons_MockCapabilities(t *testing.T) {
	testLogger := framework.L

	// Load and validate test configuration
	in, err := framework.Load[TestConfig](t)
	require.NoError(t, err, "couldn't load test config")
	validateEnvVars(t, in)
	require.Len(t, in.NodeSets, 2, "expected 2 node sets in the test config") // Nr of DONs + gateway

	mustSetCapabilitiesFn := func(input []*ns.Input) []*keystonetypes.CapabilitiesAwareNodeSet {
		return []*keystonetypes.CapabilitiesAwareNodeSet{
			{
				Input:              input[0],
				Capabilities:       []string{keystonetypes.OCR3Capability},
				DONTypes:           []string{keystonetypes.WorkflowDON},
				BootstrapNodeIndex: 0,
			},
			{
				Input:              input[1],
				Capabilities:       []string{keystonetypes.MockCapability},
				DONTypes:           []string{keystonetypes.CapabilitiesDON}, // <----- it's crucial to set the correct DON type
				BootstrapNodeIndex: 0,
			},
		}
	}

	feedsAddresses := make([][]string, in.WorkflowLoad.Jobs)
	for i := range in.WorkflowLoad.Jobs {
		feedsAddresses[i] = make([]string, 0)
		for range in.WorkflowLoad.Feeds {
			_, id := NewFeedID(t)
			feedsAddresses[i] = append(feedsAddresses[i], id)
		}
	}

	setupOutput := setupTestEnvironment(t, testLogger, in, nil, nil, mustSetCapabilitiesFn, feedsAddresses)

	ctx := tests.Context(t)
	// Log extra information that might help debugging
	t.Cleanup(func() {
		if t.Failed() {
			logTestInfo(testLogger, in.WorkflowConfig.FeedID, in.WorkflowConfig.WorkflowName, setupOutput.feedsConsumerAddress.Hex(), setupOutput.forwarderAddress.Hex())

			logDir := fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name())

			removeErr := os.RemoveAll(logDir)
			if removeErr != nil {
				testLogger.Error().Err(removeErr).Msg("failed to remove log directory")
				return
			}

			_, saveErr := framework.SaveContainerLogs(logDir)
			if saveErr != nil {
				testLogger.Error().Err(saveErr).Msg("failed to save container logs")
				return
			}

			debugDons := make([]*keystonetypes.DebugDon, 0, len(setupOutput.donTopology.DonsWithMetadata))
			for i, donWithMetadata := range setupOutput.donTopology.DonsWithMetadata {
				containerNames := make([]string, 0, len(donWithMetadata.NodesMetadata))
				for _, output := range setupOutput.nodeOutput[i].Output.CLNodes {
					containerNames = append(containerNames, output.Node.ContainerName)
				}
				debugDons = append(debugDons, &keystonetypes.DebugDon{
					NodesMetadata:  donWithMetadata.NodesMetadata,
					Flags:          donWithMetadata.Flags,
					ContainerNames: containerNames,
				})
			}

			debugInput := keystonetypes.DebugInput{
				DebugDons:        debugDons,
				BlockchainOutput: setupOutput.blockchainOutput,
			}
			lidebug.PrintTestDebug(t.Name(), testLogger, debugInput)
		}
	})

	//Get ocr key
	kb := make([]ocr2key.KeyBundle, 0)
	for _, don := range setupOutput.donTopology.DonsWithMetadata {
		if flags.HasFlag(don.Flags, keystonetypes.MockCapability) {
			for i, n := range don.DON.Nodes {
				if i == 0 {
					continue // Skip bootstrap nodeker
				}

				key, err := n.ExportOCR2Keys(n.Ocr2KeyBundleID)
				require.NoError(t, err)
				b, err2 := json.Marshal(key)
				require.NoError(t, err2)
				kk, err2 := ocr2key.FromEncryptedJSON(b, nodeclient.ChainlinkKeyPassword)
				require.NoError(t, err2)
				kb = append(kb, kk)
			}
		}
	}

	// Export key bundles so we can import them later in another test, used when crib cluster is already setup and we just want to connect to mocks for a different test
	require.NoError(t, saveKeyBundles(kb))
	require.NoError(t, saveFeedAddresses(feedsAddresses))

	testLogger.Info().Msg("Connecting to mock capabilities...")

	mocksClient := mock_capability.NewMockCapabilityController(testLogger)
	mockClientsAddress := make([]string, 0)
	if in.Infra.InfraType == "docker" {
		mockClientsAddress = []string{"127.0.0.1:13401", "127.0.0.1:13402", "127.0.0.1:13403", "127.0.0.1:13404"}
	} else {
		for i, _ := range setupOutput.nodeOutput[1].CLNodes {
			if i == 0 {
				continue
			}
			mockClientsAddress = append(mockClientsAddress, fmt.Sprintf("%s-%s-%d-mock.main.stage.cldev.sh:443", in.Infra.CRIB.Namespace, setupOutput.nodeOutput[1].NodeSetName, i-1))
		}
	}

	useInsecure := false
	if in.Infra.InfraType == "docker" {
		useInsecure = true
	}

	require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, useInsecure, true)) //Capability don ports

	testLogger.Info().Msg("Hooking into mock executable capabilities")

	receiveChannel := make(chan capabilities.CapabilityRequest, 1000)
	require.NoError(t, mocksClient.HookExecutables(ctx, receiveChannel))

	labels := map[string]string{
		"go_test_name": "test1",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 10,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 10*time.Minute),
			),
			Gun:                   NewStreamsGun(mocksClient, kb, feedsAddresses, "streams-trigger@1.0.0", receiveChannel),
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err)

	//go sendReports(t.Context(), t, mocksClient, testLogger, kb)
	//time.Sleep(time.Minute * 30)

}

func TestReconnectMock(t *testing.T) {
	testLogger := framework.L
	ctx := tests.Context(t)

	kb, err := loadKeyBundlesFromCache()
	require.NoError(t, err)

	feedAddresses, err := loadFeedAddressesFromCache()
	require.NoError(t, err)
	testLogger.Info().Msg("Connecting to mock capabilities...")
	var mocksClient *mock_capability.MockCapabilityController

	mocksClient, err = mock_capability.NewMockCapabilityControllerFromCache(testLogger, false)
	if err != nil {
		mocksClient = mock_capability.NewMockCapabilityController(testLogger)
		mockClientsAddress := []string{"crib-local-capabilities-0-mock.main.stage.cldev.sh", "crib-local-capabilities-1-mock.main.stage.cldev.sh", "crib-local-capabilities-2-mock.main.stage.cldev.sh", "crib-local-capabilities-3-mock.main.stage.cldev.sh"}

		require.NoError(t, mocksClient.ConnectAll(mockClientsAddress, false, true)) //Capability don ports
	}

	testLogger.Info().Msg("Hooking into mock executable capabilities")

	receiveChannel := make(chan capabilities.CapabilityRequest, 1000)
	require.NoError(t, mocksClient.HookExecutables(ctx, receiveChannel))

	labels := map[string]string{
		"go_test_name": "Workflow DON Load Test",
		"branch":       "profile-check",
		"commit":       "profile-check",
	}

	_, err = wasp.NewProfile().
		Add(wasp.NewGenerator(&wasp.Config{
			CallTimeout: time.Minute * 5,
			LoadType:    wasp.RPS,
			Schedule: wasp.Combine(
				wasp.Plain(4, 10*time.Minute),
			),
			Gun:                   NewStreamsGun(mocksClient, kb, feedAddresses, "streams-trigger@1.0.0", receiveChannel),
			Labels:                labels,
			LokiConfig:            wasp.NewEnvLokiConfig(),
			RateLimitUnitDuration: time.Minute,
		})).
		Run(true)
	require.NoError(t, err)
}

var _ wasp.Gun = (*StreamsGun)(nil)

type StreamsGun struct {
	capProxy    *mock_capability.MockCapabilityController
	keyBundles  []ocr2key.KeyBundle
	feeds       [][]string
	triggerID   string
	waitChans   map[int64]chan interface{}
	recieveChan <-chan capabilities.CapabilityRequest
	mu          sync.Mutex
}

func NewStreamsGun(capProxy *mock_capability.MockCapabilityController, keyBundles []ocr2key.KeyBundle, feeds [][]string, triggerID string, ch <-chan capabilities.CapabilityRequest) *StreamsGun {
	sg := &StreamsGun{
		capProxy:    capProxy,
		keyBundles:  keyBundles,
		feeds:       feeds,
		triggerID:   triggerID,
		recieveChan: ch,
	}
	go sg.waitHOOKloop()
	return sg
}

func (s *StreamsGun) Call(l *wasp.Generator) *wasp.Response {
	timesteamp := time.Now().Unix()
	err := s.prepareWaitHOOK(timesteamp)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	reports := make([]datastreams.FeedReport, 0)
	for i := range s.feeds {
		for _, feed := range s.feeds[i] {
			r, err := createFeedReport(big.NewInt(int64(rand.IntN(100))), timesteamp, feed, s.keyBundles)
			if err != nil {
				return &wasp.Response{Error: err.Error()}
			}
			reports = append(reports, *r)
		}
	}

	event, err := values.WrapMap(datastreams.StreamsTriggerEvent{
		Payload:   reports,
		Metadata:  datastreams.Metadata{},
		Timestamp: timesteamp,
	})
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	framework.L.Info().Msg(fmt.Sprintf("Sending trigger response %v+", event))
	eventBytes, err := mock_capability.MapToBytes(event)

	err = s.capProxy.SendTrigger(context.TODO(), s.triggerID, uuid.NewString(), eventBytes)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}

	//Wait for hook back with the same eventID
	err = s.waitForHOOK(timesteamp)
	if err != nil {
		return &wasp.Response{Error: err.Error()}
	}
	return &wasp.Response{}
}

func (s *StreamsGun) waitHOOKloop() error {
	for {
		select {
		case m, ok := <-s.recieveChan:
			if !ok {

				return fmt.Errorf("channel closed")
			}

			inputs, err := decodeTargetInput(m.Inputs)
			if err != nil {
				fmt.Println("error decoding inputs")
			}

			//To get the timestamp we look at the last 64 chars of the hex encoded report
			hexReport := hex.EncodeToString(inputs.Inputs.SignedReport.Report)
			timestampInHex := hexReport[len(hexReport)-64:]
			timestamp, err := strconv.ParseInt(timestampInHex, 16, 64)
			if err != nil {
				fmt.Println("error parsing timestamp")
			}

			s.mu.Lock()
			//Check if exist
			if ch, exist := s.waitChans[int64(timestamp)]; exist {
				s.mu.Unlock()
				ch <- m //This is blocking
			} else {
				s.mu.Unlock()
			}
		}
	}

	return nil
}

func (s *StreamsGun) prepareWaitHOOK(reportTimestamp int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.waitChans == nil {
		s.waitChans = make(map[int64]chan interface{})
	}
	if _, exists := s.waitChans[reportTimestamp]; exists {
		return fmt.Errorf("cannot prepare for HOOK, timestamp  %d already exits", reportTimestamp)
	}
	s.waitChans[reportTimestamp] = make(chan interface{})
	return nil
}

func (s *StreamsGun) waitForHOOK(timestamp int64) error {
	s.mu.Lock()

	if ch, exists := s.waitChans[timestamp]; !exists {
		s.mu.Unlock()
		return fmt.Errorf("cannot wait for HOOK, timestamp  %q does not exist", timestamp)
	} else {
		s.mu.Unlock()
		<-ch
		delete(s.waitChans, timestamp)
		return nil
	}

}

func createFeedReport(price *big.Int, observationTimestamp int64,
	feedIDString string,
	keyBundles []ocr2key.KeyBundle) (*datastreams.FeedReport, error) {
	reportCtx := ocrTypes.ReportContext{}
	rawCtx := RawReportContext(reportCtx)

	bytes, err := hex.DecodeString(feedIDString[2:])
	if err != nil {
		return nil, err
	}
	var feedIDBytes [32]byte
	copy(feedIDBytes[:], bytes)

	fullReport, err := newReport(feedIDBytes, price, observationTimestamp)
	if err != nil {
		return nil, err
	}

	report := &datastreams.FeedReport{
		FeedID:               feedIDString,
		FullReport:           fullReport,
		BenchmarkPrice:       price.Bytes(),
		ObservationTimestamp: observationTimestamp,
		Signatures:           [][]byte{},
		ReportContext:        rawCtx,
	}

	for _, key := range keyBundles {
		sig, err := key.Sign(reportCtx, report.FullReport)
		if err != nil {
			return nil, err
		}
		report.Signatures = append(report.Signatures, sig)
	}

	return report, nil
}
func RawReportContext(reportCtx ocrTypes.ReportContext) []byte {
	rc := evmutil.RawReportContext(reportCtx)
	flat := []byte{}
	for _, r := range rc {
		flat = append(flat, r[:]...)
	}
	return flat
}

func newReport(feedID [32]byte, price *big.Int, timestamp int64) ([]byte, error) {
	lggr, _ := logger.NewLogger()
	v3Codec := reportcodec.NewReportCodec(feedID, lggr)
	raw, err := v3Codec.BuildReport(context.TODO(), v3.ReportFields{
		BenchmarkPrice: price,
		Timestamp:      uint32(timestamp),
		Bid:            big.NewInt(0),
		Ask:            big.NewInt(0),
		LinkFee:        big.NewInt(0),
		NativeFee:      big.NewInt(0),
	})
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func decodeTargetInput(inputs *values.Map) (targets.Request, error) {
	var r targets.Request
	const signedReportField = "signed_report"
	signedReport, ok := inputs.Underlying[signedReportField]
	if !ok {
		return r, fmt.Errorf("missing required field %s", signedReportField)
	}

	if err := signedReport.UnwrapTo(&r.Inputs.SignedReport); err != nil {
		return r, err
	}

	return r, nil
}

func saveKeyBundles(keyBundles []ocr2key.KeyBundle) error {
	cacheDir := "cache/keys"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	for i, kb := range keyBundles {
		bytes, err := kb.Marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal key bundle %d: %w", i, err)
		}

		filename := fmt.Sprintf("%s/key_bundle_%d.json", cacheDir, i)
		if err := os.WriteFile(filename, bytes, 0644); err != nil {
			return fmt.Errorf("failed to write key bundle %d to file: %w", i, err)
		}
	}
	return nil
}

func loadKeyBundlesFromCache() ([]ocr2key.KeyBundle, error) {
	cacheDir := "cache/keys"
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var keyBundles []ocr2key.KeyBundle
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "key_bundle_") {
			bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", cacheDir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read key bundle file %s: %w", file.Name(), err)
			}

			kb, err := ocr2key.New(chaintype.EVM)
			if err != nil {
				return nil, fmt.Errorf("cannot create new key bundle from %s: %w", file.Name(), err)
			}
			if err := kb.Unmarshal(bytes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal key bundle from %s: %w", file.Name(), err)
			}
			keyBundles = append(keyBundles, kb)
		}
	}

	if len(keyBundles) == 0 {
		return nil, fmt.Errorf("no key bundles found in cache directory")
	}
	return keyBundles, nil
}

func Decode(report []byte) (uint32, error) {
	//(bytes32 FeedID, uint224 Price, uint32 Timestamp)[] Reports
	framework.L.Info().Msgf("Attemptiong to decode %s", hex.EncodeToString(report))
	config := map[string]any{
		"abi": "(bytes32,uint224,uint32)[]",
	}
	wrapped, err := values.NewMap(config)
	if err != nil {
		return 0, err
	}
	c, err := NewEVMDecoder(wrapped)
	if err != nil {
		return 0, err
	}

	var r []struct{}

	err = c.Decode(context.TODO(), report, &r, "user")
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func NewEVMDecoder(config *values.Map) (types2.Decoder, error) {
	// parse the "inner" encoder config - user-defined fields
	wrappedSelector, err := config.Underlying["abi"].Unwrap()
	if err != nil {
		return nil, err
	}
	selectorStr, ok := wrappedSelector.(string)
	if !ok {
		return nil, fmt.Errorf("expected %s to be a string", "abi")
	}
	selector, err := abi.ParseSelector("inner(" + selectorStr + ")")
	if err != nil {
		return nil, err
	}
	jsonSelector, err := json.Marshal(selector.Inputs)
	if err != nil {
		return nil, err
	}

	chainCodecConfig := types.ChainCodecConfig{
		TypeABI: string(jsonSelector),
	}

	codecConfig := types.CodecConfig{Configs: map[string]types.ChainCodecConfig{
		"user": chainCodecConfig,
	}}

	c, err := codec.NewCodec(codecConfig)
	if err != nil {
		return nil, err
	}

	return c, nil
}
func NewFeedID(t *testing.T) ([32]byte, string) {
	buf := [32]byte{}
	_, err := crand.Read(buf[:])
	require.NoError(t, err)
	return buf, "0x" + hex.EncodeToString(buf[:])
}

func saveFeedAddresses(feedsAddresses [][]string) error {
	cacheDir := "cache/feeds"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	filename := fmt.Sprintf("%s/feed_addresses.json", cacheDir)
	bytes, err := json.Marshal(feedsAddresses)
	if err != nil {
		return fmt.Errorf("failed to marshal feed addresses: %w", err)
	}

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write feed addresses to file: %w", err)
	}

	return nil
}

func loadFeedAddressesFromCache() ([][]string, error) {
	cacheDir := "cache/feeds"
	filename := fmt.Sprintf("%s/feed_addresses.json", cacheDir)

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read feed addresses file: %w", err)
	}

	var feedsAddresses [][]string
	if err := json.Unmarshal(bytes, &feedsAddresses); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed addresses: %w", err)
	}

	return feedsAddresses, nil
}
