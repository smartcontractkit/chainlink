// This file contains data generators and utilities to simplify tests.
// The data generated here shouldn't be used to run OCR instances
package monitoring

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/linkedin/goavro"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/smartcontractkit/chainlink-relay/pkg/monitoring/pb"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"
)

// Sources

func NewFakeRDDSource(minFeeds, maxFeeds uint8) Source {
	return &fakeRddSource{minFeeds, maxFeeds}
}

type fakeRddSource struct {
	minFeeds, maxFeeds uint8
}

func (f *fakeRddSource) Name() string {
	return "fake-rdd"
}

func (f *fakeRddSource) Fetch(_ context.Context) (interface{}, error) {
	numFeeds := int(f.minFeeds) + rand.Intn(int(f.maxFeeds-f.minFeeds))
	feeds := make([]FeedConfig, numFeeds)
	for i := 0; i < numFeeds; i++ {
		feeds[i] = generateFeedConfig()
	}
	return feeds, nil
}

func (f *fakeRandomDataSourceFactory) Run(ctx context.Context, log Logger) {
	update, err := generateEnvelope()
	if err != nil {
		log.Errorw("failed to generate fake read from chain", "error", err)
	}
	for {
		select {
		case f.updates <- update:
			log.Infow("generate envelope")
			update, err = generateEnvelope()
			if err != nil {
				log.Errorw("failed to generate fake read from chain", "error", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

var _ SourceFactory = (*fakeRandomDataSourceFactory)(nil)

type fakeRandomDataSourceFactory struct {
	updates chan Envelope
	ctx     context.Context
}

func (f *fakeRandomDataSourceFactory) NewSource(chainConfig ChainConfig, feedConfig FeedConfig) (Source, error) {
	return &fakeSource{f}, nil
}

type fakeSource struct {
	factory *fakeRandomDataSourceFactory
}

func (f *fakeSource) Fetch(ctx context.Context) (interface{}, error) {
	var update Envelope
	select {
	case update = <-f.factory.updates:
		return update, nil
	case <-f.factory.ctx.Done():
		return nil, fmt.Errorf("source closed")
	}
}

type fakeSourceWithWait struct {
	waitOnRead time.Duration
}

func (f *fakeSourceWithWait) Fetch(ctx context.Context) (interface{}, error) {
	select {
	case <-time.After(f.waitOnRead):
		return 1, nil
	case <-ctx.Done():
		return 0, nil
	}
}

type fakeSourceFactoryWithError struct {
	updates     chan interface{}
	errors      chan error
	returnError bool
}

func (f *fakeSourceFactoryWithError) NewSource(_ ChainConfig, _ FeedConfig) (Source, error) {
	if f.returnError {
		return nil, fmt.Errorf("fake source factory error")
	}
	return &fakeSourceWithError{
		f.updates,
		f.errors,
	}, nil
}

type fakeSourceWithError struct {
	updates chan interface{}
	errors  chan error
}

func (f *fakeSourceWithError) Fetch(ctx context.Context) (interface{}, error) {
	select {
	case update := <-f.updates:
		return update, nil
	case err := <-f.errors:
		return nil, err
	case <-ctx.Done():
		return nil, nil
	}
}

// Exporters

type fakeExporterFactory struct {
	data        chan interface{}
	returnError bool
}

func (f *fakeExporterFactory) NewExporter(chainConfig ChainConfig, feedConfig FeedConfig) (Exporter, error) {
	if f.returnError {
		return nil, fmt.Errorf("fake exporter factory error")
	}
	return &fakeExporter{
		f.data,
	}, nil
}

type fakeExporter struct {
	data chan interface{}
}

func (f *fakeExporter) Export(ctx context.Context, data interface{}) {
	select {
	case f.data <- data:
	case <-ctx.Done():
	}
}

func (f *fakeExporter) Cleanup(_ context.Context) {
}

// Generators

func generate32ByteArr() [32]byte {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		panic("unable to generate [32]byte from rand")
	}
	var out [32]byte
	copy(out[:], buf[:32])
	return out
}

type fakeFeedConfig struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Path           string `json:"path,omitempty"`
	Symbol         string `json:"symbol,omitempty"`
	HeartbeatSec   int64  `json:"heartbeat,omitempty"`
	ContractType   string `json:"contract_type,omitempty"`
	ContractStatus string `json:"status,omitempty"`
	// This functions as a feed identifier.
	ContractAddress []byte `json:"contract_address,omitempty"`
	Multiply        uint64 `json:"multiply,omitempty"`
}

func (f fakeFeedConfig) GetID() string             { return f.ID }
func (f fakeFeedConfig) GetName() string           { return f.Name }
func (f fakeFeedConfig) GetPath() string           { return f.Path }
func (f fakeFeedConfig) GetSymbol() string         { return f.Symbol }
func (f fakeFeedConfig) GetHeartbeatSec() int64    { return f.HeartbeatSec }
func (f fakeFeedConfig) GetContractType() string   { return f.ContractType }
func (f fakeFeedConfig) GetContractStatus() string { return f.ContractStatus }
func (f fakeFeedConfig) GetContractAddress() string {
	return base64.StdEncoding.EncodeToString(f.ContractAddress)
}
func (f fakeFeedConfig) GetContractAddressBytes() []byte { return f.ContractAddress }
func (f fakeFeedConfig) GetMultiply() uint64             { return f.Multiply }
func (f fakeFeedConfig) ToMapping() map[string]interface{} {
	return map[string]interface{}{
		"feed_name":        f.Name,
		"feed_path":        f.Path,
		"symbol":           f.Symbol,
		"heartbeat_sec":    int64(f.HeartbeatSec),
		"contract_type":    f.ContractType,
		"contract_status":  f.ContractStatus,
		"contract_address": f.ContractAddress,
		// These are solana specific but are kept here for backwards compatibility in Avro.
		"transmissions_account": []byte{},
		"state_account":         []byte{},
	}
}

func generateFeedConfig() FeedConfig {
	coins := []string{"btc", "eth", "matic", "link", "avax", "ftt", "srm", "usdc", "sol", "ray"}
	coin := coins[rand.Intn(len(coins))]
	contractAddress := generate32ByteArr()
	return fakeFeedConfig{
		ID:              hex.EncodeToString(contractAddress[:]),
		Name:            fmt.Sprintf("%s / usd", coin),
		Path:            fmt.Sprintf("%s-usd", coin),
		Symbol:          "$",
		HeartbeatSec:    1,
		ContractType:    "ocr2",
		ContractStatus:  "status",
		ContractAddress: contractAddress[:],
		Multiply:        1000,
	}
}

func generateNumericalMedianOffchainConfig() (*pb.NumericalMedianConfigProto, []byte, error) {
	out := &pb.NumericalMedianConfigProto{
		AlphaReportInfinite: ([]bool{true, false})[rand.Intn(2)],
		AlphaReportPpb:      rand.Uint64(),
		AlphaAcceptInfinite: ([]bool{true, false})[rand.Intn(2)],
		AlphaAcceptPpb:      rand.Uint64(),
		DeltaCNanoseconds:   rand.Uint64(),
	}
	buf, err := proto.Marshal(out)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal median plugin config: %w", err)
	}
	return out, buf, nil
}

func generateOffchainConfig(numOracles int) (
	*pb.OffchainConfigProto,
	*pb.NumericalMedianConfigProto,
	[]byte,
	error,
) {
	numericalMedianOffchainConfig, encodedNumericalMedianOffchainConfig, err := generateNumericalMedianOffchainConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	schedule := []uint32{}
	for i := 0; i < 10; i++ {
		schedule = append(schedule, 1)
	}
	offchainPublicKeys := [][]byte{}
	for i := 0; i < numOracles; i++ {
		randArr := generate32ByteArr()
		offchainPublicKeys = append(offchainPublicKeys, randArr[:])
	}
	peerIDs := []string{}
	for i := 0; i < numOracles; i++ {
		peerIDs = append(peerIDs, fmt.Sprintf("peer#%d", i))
	}
	config := &pb.OffchainConfigProto{
		DeltaProgressNanoseconds: rand.Uint64(),
		DeltaResendNanoseconds:   rand.Uint64(),
		DeltaRoundNanoseconds:    rand.Uint64(),
		DeltaGraceNanoseconds:    rand.Uint64(),
		DeltaStageNanoseconds:    rand.Uint64(),

		RMax:                  rand.Uint32(),
		S:                     schedule,
		OffchainPublicKeys:    offchainPublicKeys,
		PeerIds:               peerIDs,
		ReportingPluginConfig: encodedNumericalMedianOffchainConfig,

		MaxDurationQueryNanoseconds:       rand.Uint64(),
		MaxDurationObservationNanoseconds: rand.Uint64(),
		MaxDurationReportNanoseconds:      rand.Uint64(),

		MaxDurationShouldAcceptFinalizedReportNanoseconds:  rand.Uint64(),
		MaxDurationShouldTransmitAcceptedReportNanoseconds: rand.Uint64(),

		SharedSecretEncryptions: &pb.SharedSecretEncryptionsProto{
			DiffieHellmanPoint: []byte{'p', 'o', 'i', 'n', 't'},
			SharedSecretHash:   []byte{'h', 'a', 's', 'h'},
			Encryptions:        [][]byte{[]byte("encryption")},
		},
	}
	encodedConfig, err := proto.Marshal(config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to marshal offchain config: %w", err)
	}
	return config, numericalMedianOffchainConfig, encodedConfig, nil
}

func generateContractConfig(n int) (
	types.ContractConfig,
	median.OnchainConfig,
	*pb.OffchainConfigProto,
	*pb.NumericalMedianConfigProto,
	error,
) {
	signers := make([]types.OnchainPublicKey, n)
	transmitters := make([]types.Account, n)
	for i := 0; i < n; i++ {
		randArr := generate32ByteArr()
		signers[i] = types.OnchainPublicKey(randArr[:])
		transmitters[i] = types.Account(hexutil.Encode([]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint8(i),
		}))
	}
	onchainConfig := median.OnchainConfig{
		Min: big.NewInt(rand.Int63()),
		Max: big.NewInt(rand.Int63()),
	}
	onchainConfigEncoded, err := onchainConfig.Encode()
	if err != nil {
		return types.ContractConfig{}, median.OnchainConfig{}, nil, nil, err
	}
	offchainConfig, pluginOffchainConfig, offchainConfigEncoded, err := generateOffchainConfig(n)
	if err != nil {
		return types.ContractConfig{}, median.OnchainConfig{}, nil, nil, err
	}
	contractConfig := types.ContractConfig{
		ConfigDigest:          generate32ByteArr(),
		ConfigCount:           rand.Uint64(),
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     uint8(10),
		OnchainConfig:         onchainConfigEncoded,
		OffchainConfigVersion: rand.Uint64(),
		OffchainConfig:        offchainConfigEncoded,
	}
	return contractConfig, onchainConfig, offchainConfig, pluginOffchainConfig, nil
}

func generateEnvelope() (Envelope, error) {
	generated, _, _, _, err := generateContractConfig(31)
	if err != nil {
		return Envelope{}, err
	}
	return Envelope{
		ConfigDigest:    generated.ConfigDigest,
		Round:           uint8(rand.Intn(256)),
		Epoch:           rand.Uint32(),
		LatestAnswer:    big.NewInt(rand.Int63()),
		LatestTimestamp: time.Now(),

		ContractConfig: generated,

		BlockNumber: rand.Uint64(),
		Transmitter: types.Account(hexutil.Encode([]byte{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, uint8(rand.Intn(32)),
		})),

		LinkBalance: rand.Uint64(),
	}, nil
}

type fakeChainConfig struct {
	RPCEndpoint  string
	NetworkName  string
	NetworkID    string
	ChainID      string
	ReadTimeout  time.Duration
	PollInterval time.Duration
}

func generateChainConfig() ChainConfig {
	return fakeChainConfig{
		RPCEndpoint:  "http://some-chain-host:6666",
		NetworkName:  "mainnet-beta",
		NetworkID:    "1",
		ChainID:      "mainnet-beta",
		ReadTimeout:  100 * time.Millisecond,
		PollInterval: time.Duration(1+rand.Intn(5)) * time.Second,
	}
}

func (f fakeChainConfig) GetRPCEndpoint() string         { return f.RPCEndpoint }
func (f fakeChainConfig) GetNetworkName() string         { return f.NetworkName }
func (f fakeChainConfig) GetNetworkID() string           { return f.NetworkID }
func (f fakeChainConfig) GetChainID() string             { return f.ChainID }
func (f fakeChainConfig) GetReadTimeout() time.Duration  { return f.ReadTimeout }
func (f fakeChainConfig) GetPollInterval() time.Duration { return f.PollInterval }

func (f fakeChainConfig) ToMapping() map[string]interface{} {
	return map[string]interface{}{
		"network_name": f.NetworkName,
		"network_id":   f.NetworkID,
		"chain_id":     f.ChainID,
	}
}

// Metrics

type devnullMetrics struct{}

var _ Metrics = (*devnullMetrics)(nil)

func (d *devnullMetrics) SetHeadTrackerCurrentHead(blockNumber uint64, networkName, chainID, networkID string) {
}

func (d *devnullMetrics) SetFeedContractMetadata(chainID, contractAddress, feedID, contractStatus, contractType, feedName, feedPath, networkID, networkName, symbol string) {
}

func (d *devnullMetrics) SetFeedContractLinkBalance(balance uint64, chainID, contractAddress, feedID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
}

func (d *devnullMetrics) SetNodeMetadata(chainID, networkID, networkName, oracleName, sender string) {
}

func (d *devnullMetrics) SetOffchainAggregatorAnswersRaw(answer *big.Int, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
}

func (d *devnullMetrics) SetOffchainAggregatorAnswers(answer *big.Float, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
}

func (d *devnullMetrics) IncOffchainAggregatorAnswersTotal(contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
}

func (d *devnullMetrics) SetOffchainAggregatorSubmissionReceivedValues(value *big.Float, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
}

func (d *devnullMetrics) SetOffchainAggregatorAnswerStalled(isSet bool, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
}

func (d *devnullMetrics) Cleanup(networkName, networkID, chainID, oracleName, sender, feedName, feedPath, symbol, contractType, contractStatus, contractAddress, feedID string) {
}

func (d *devnullMetrics) HTTPHandler() http.Handler {
	return promhttp.Handler()
}

type keepLatestMetrics struct {
	*devnullMetrics

	latestTransmission *big.Float
	latestTransmitter  string
}

func (k *keepLatestMetrics) SetOffchainAggregatorAnswers(answer *big.Float, contractAddress, feedID, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	k.latestTransmission = &big.Float{}
	k.latestTransmission.Set(answer)
}

func (k *keepLatestMetrics) SetOffchainAggregatorSubmissionReceivedValues(value *big.Float, contractAddress, feedID, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	k.latestTransmission = &big.Float{}
	k.latestTransmission.Set(value)
	k.latestTransmitter = sender
}

// Producer

type producerMessage struct {
	key, value []byte
	topic      string
}

type fakeProducer struct {
	sendCh chan producerMessage
	ctx    context.Context
}

func (f fakeProducer) Produce(key, value []byte, topic string) error {
	select {
	case f.sendCh <- producerMessage{key, value, topic}:
	case <-f.ctx.Done():
	}
	return nil
}

// Schema

type fakeSchema struct {
	codec   *goavro.Codec
	subject string
}

func (f fakeSchema) ID() int {
	return 1
}

func (f fakeSchema) Version() int {
	return 1
}

func (f fakeSchema) Subject() string {
	return f.subject
}

func (f fakeSchema) Encode(value interface{}) ([]byte, error) {
	return f.codec.BinaryFromNative(nil, value)
}

func (f fakeSchema) Decode(buf []byte) (interface{}, error) {
	value, _, err := f.codec.NativeFromBinary(buf)
	return value, err
}

// Poller

type fakePoller struct {
	numUpdates int
	ch         chan interface{}
}

func (f *fakePoller) Run(ctx context.Context) {
	source := &fakeRddSource{1, 2}
	for i := 0; i < f.numUpdates; i++ {
		updates, _ := source.Fetch(ctx)
		select {
		case f.ch <- updates:
		case <-ctx.Done():
			return
		}
	}
}

func (f *fakePoller) Updates() <-chan interface{} {
	return f.ch
}

// Logger

type nullLogger struct{}

func newNullLogger() Logger {
	return &nullLogger{}
}

func (n *nullLogger) With(args ...interface{}) Logger {
	return n
}

func (n *nullLogger) Tracew(format string, values ...interface{})    {}
func (n *nullLogger) Debugw(format string, values ...interface{})    {}
func (n *nullLogger) Infow(format string, values ...interface{})     {}
func (n *nullLogger) Warnw(format string, values ...interface{})     {}
func (n *nullLogger) Errorw(format string, values ...interface{})    {}
func (n *nullLogger) Criticalw(format string, values ...interface{}) {}
func (n *nullLogger) Panicw(format string, values ...interface{})    {}
func (n *nullLogger) Fatalw(format string, values ...interface{})    {}

// This utilities are used primarely in tests but are present in the monitoring package because they are not inside a file ending in _test.go.
// This is done in order to expose NewRandomDataReader for use in cmd/monitoring.
// The following code is added to comply with the "unused" linter:
var (
	_ = generateChainConfig()
	_ = generateFeedConfig
	_ = fakeProducer{}
	_ = fakeSchema{}
	_ = keepLatestMetrics{}
	_ = fakePoller{}
	_ = newNullLogger()
	_ = fakeExporterFactory{}
	_ = fakeSourceWithWait{}
	_ = fakeSourceFactoryWithError{}
)
