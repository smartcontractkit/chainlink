package action

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"golang.org/x/exp/maps"

	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/metrics"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/plugin"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/transmitter"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/blobconsensus/oracle/types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/requests"
	caperrors "github.com/smartcontractkit/chainlink-common/pkg/capabilities/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/consensus/server"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"

	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	valuespb "github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
)

const (
	defaultMaxRequestOutcomeSize        = 10000
	defaultKeyBundleIDForValueConsensus = "evm"
	KeyBundleIDEvm                      = "evm"
	KeyBundleIDAptos                    = "aptos"
	KeyBundleIDSolana                   = "solana"
	SigningAlgoEcdsa                    = "ecdsa"
	SigningAlgoEd25519                  = "ed25519"
	HashingAlgoKeccak256                = "keccak256"
	HashingAlgoBlake2b256               = "blake2b_256"
)

type ConsensusCapabilityConfig struct {
	RequestBatchSize             int
	MaxRequestSizeBytes          int
	KeyBundleIDForValueConsensus string
	MaxRequestOutcomeSize        int
}

var _ server.ConsensusCapability = &consensusCapability{}

type consensusCapability struct {
	lggr          logger.Logger
	oracle        core.Oracle
	reqStore      *requests.Store[*oracle.ConsensusRequest]
	reqHandler    *requests.Handler[*oracle.ConsensusRequest, oracle.ConsensusResponse]
	limitsFactory limits.Factory

	requestTimeout     time.Duration
	requestTimeoutLock sync.RWMutex

	maxRequestSizeBytes       limits.BoundLimiter[config.Size]
	valueConsensusKeyBundleID string
	maxRequestOutcomeSize     int

	metrics *metrics.Metrics
}

type storeStatsCollector struct {
	requestStoreRequests metric.Int64Gauge
}

func (s *storeStatsCollector) SetRequestCount(requestCount int) {
	s.requestStoreRequests.Record(context.Background(), int64(requestCount))
}

// NewConsensusCapability creates a new ConsensusCapability with the given logger, clock, and response cache expiry time.  The
// response cache expiry controls how long a response for a given request is cached before it is considered expired and evicted. This allows
// the capability to respond to slow requests sent after consensus has been reached.
func NewConsensusCapability(lggr logger.Logger, clock clockwork.Clock, responseCacheExpiry time.Duration,
	limitsFactory limits.Factory,
) (*consensusCapability, error) {
	metrics, err := metrics.NewMetrics()
	if err != nil {
		return nil, fmt.Errorf("error creating metrics: %w", err)
	}

	reqStore := requests.NewStoreWithStatsCollector[*oracle.ConsensusRequest](
		&storeStatsCollector{
			requestStoreRequests: metrics.PendingConsensusRequests,
		})

	return &consensusCapability{
		lggr:          lggr,
		reqStore:      reqStore,
		reqHandler:    requests.NewHandler(lggr, reqStore, clock, responseCacheExpiry),
		metrics:       metrics,
		limitsFactory: limitsFactory,
	}, nil
}

// SetRequestTimeout is used by the reporting plugin to set the request timeout for consensus requests.  The plugin
// receives the timeout after the Initialise method on the capability is called, thus it uses this method to set it on the capability.
func (c *consensusCapability) SetRequestTimeout(timeout time.Duration) {
	c.requestTimeoutLock.Lock()
	defer c.requestTimeoutLock.Unlock()
	c.requestTimeout = timeout
}

func (c *consensusCapability) Initialise(ctx context.Context, dependencies core.StandardCapabilitiesDependencies) error {
	c.lggr.Debugf("Initialising Consensus Capability")

	if err := c.setConfiguration(dependencies.Config); err != nil {
		return fmt.Errorf("error setting consensus capability configuration: %w", err)
	}

	reportingPlugin, err := plugin.NewReportingPluginFactory(c.lggr, c.metrics, c.reqStore, c.SetRequestTimeout,
		defaultKeyBundleIDForValueConsensus, c.maxRequestOutcomeSize)
	if err != nil {
		return fmt.Errorf("error when creating reporting plugin factory: %w", err)
	}

	contractTransmitter := transmitter.NewContractTransmitter(c.lggr, c.SendResponse)

	// These values set to the maximum permitted, response time for config update is not critical
	localOcrConfig := ocrtypes.LocalConfig{
		BlockchainTimeout:                  time.Second * 20,
		ContractConfigTrackerPollInterval:  time.Second * 60,
		ContractConfigConfirmations:        1,
		ContractTransmitterTransmitTimeout: time.Second * 60,
		DatabaseTimeout:                    time.Second * 10,
		ContractConfigLoadTimeout:          time.Second * 60,
		DefaultMaxDurationInitialization:   time.Second * 60,
	}

	oracle, err := dependencies.OracleFactory.NewOracle(ctx, core.OracleArgs{
		LocalConfig:                   localOcrConfig,
		ReportingPluginFactoryService: reportingPlugin,
		ContractTransmitter:           contractTransmitter,
	})
	if err != nil {
		return fmt.Errorf("error when creating oracle: %w", err)
	}

	c.oracle = oracle
	err = c.reqHandler.Start(context.Background())
	if err != nil {
		return fmt.Errorf("error when starting request handler: %w", err)
	}

	err = c.oracle.Start(context.Background())
	if err != nil {
		return fmt.Errorf("error when starting oracle: %w", err)
	}

	c.lggr.Debug("Initialised Consensus Capability")

	return nil
}

func (c *consensusCapability) setConfiguration(cfg string) error {
	c.valueConsensusKeyBundleID = defaultKeyBundleIDForValueConsensus

	var capabilityConfig ConsensusCapabilityConfig
	if len(cfg) > 0 {
		err := json.Unmarshal([]byte(cfg), &capabilityConfig)
		if err != nil {
			return fmt.Errorf("failed to deserialize config into ConsensusCapabilityConfig: %w", err)
		}
	}

	if capabilityConfig.MaxRequestOutcomeSize > 0 {
		c.maxRequestOutcomeSize = capabilityConfig.MaxRequestOutcomeSize
	} else {
		c.maxRequestOutcomeSize = defaultMaxRequestOutcomeSize
	}

	c.lggr.Debugw("Parsed Consensus Capability config", "config", capabilityConfig)

	if len(capabilityConfig.KeyBundleIDForValueConsensus) > 0 {
		c.valueConsensusKeyBundleID = capabilityConfig.KeyBundleIDForValueConsensus
	}

	// Get the limit from CRE settings or capability config
	requestSizeLimit := cresettings.Default.PerWorkflow.Consensus.ObservationSizeLimit // make a copy

	var configuredLimit config.Size
	if capabilityConfig.MaxRequestSizeBytes > 0 {
		// Use capability config if explicitly set
		configuredLimit = config.Size(capabilityConfig.MaxRequestSizeBytes)
	} else {
		// Otherwise use CRE settings default
		configuredLimit = requestSizeLimit.DefaultValue
	}

	// Cap the limit at libOCR's MaxMaxObservationLength to ensure it never exceeds OCR protocol limits
	libOCRLimit := config.Size(int(ocr3types.MaxMaxObservationLength))
	if configuredLimit > libOCRLimit {
		c.lggr.Warnw("Request size limit exceeds libOCR maximum, capping at libOCR limit",
			"configuredLimit", configuredLimit,
			"libOCRLimit", libOCRLimit)
		configuredLimit = libOCRLimit
	}

	requestSizeLimit.DefaultValue = configuredLimit
	maxRequestSizeBytes, err := limits.MakeBoundLimiter(c.limitsFactory, requestSizeLimit)
	if err != nil {
		return err
	}
	c.maxRequestSizeBytes = maxRequestSizeBytes
	return nil
}

func (c *consensusCapability) Simple(ctx context.Context, metadata capabilities.RequestMetadata, input *sdk.SimpleConsensusInputs) (*capabilities.ResponseAndMetadata[*valuespb.Value], caperrors.Error) {
	ctx = metadata.ContextWithCRE(ctx)
	lggr := c.requestLggr(metadata)

	if err := decodeObservationType(lggr, input); err != nil {
		responseAndMetadata := capabilities.ResponseAndMetadata[*valuespb.Value]{}
		return &responseAndMetadata, caperrors.NewPublicSystemError(fmt.Errorf("failed to decode observation: %s", err), caperrors.InvalidArgument)
	}

	lggr.Debugw("received simple consensus request", "metadata", metadata)

	consensusRequestMetaData := oracle.ConsensusRequestMetadata{
		RequestMetadata: metadata,
		KeyBundleID:     c.valueConsensusKeyBundleID,
		ReportID:        "0000", // Report ID is not used for value consensus, so we can use a dummy value
		RequestType:     types.RequestType_VALUE_CONSENSUS,
	}

	requestSize, reqSizeErr := c.validateRequestSize(ctx, lggr, consensusRequestMetaData, input)
	if reqSizeErr != nil {
		var capErr caperrors.Error
		if errors.As(reqSizeErr, &capErr) && capErr.Origin() != caperrors.OriginUser {
			// Return system errors immediately
			lggr.Errorw("failed to validate request size", "err", reqSizeErr)
			return nil, reqSizeErr
		}

		lggr.Debugw("request size validation failed with user error, sending as error input to consensus",
			"err", reqSizeErr,
			"requestSize", requestSize)

		input = &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Error{
				Error: reqSizeErr.Error(),
			},
			Descriptors: input.Descriptors,
			Default:     input.Default,
		}
	}

	c.metrics.RecordRequestSize(ctx, float64(requestSize))

	callbackChan := c.sendRequest(ctx, input, consensusRequestMetaData)

	select {
	case <-ctx.Done():
		return nil, caperrors.NewPublicSystemError(ctx.Err(), caperrors.Canceled)
	case response := <-callbackChan:
		if response.Err != nil {
			return nil, response.Err
		}

		// Remove the metadata prefix from the raw report to get the serialised value
		serialisedValue := response.RawReport[plugin.ReportMetaDataPrependLength:]

		valueProto := &valuespb.Value{}
		if err := proto.Unmarshal(serialisedValue, valueProto); err != nil {
			return nil, caperrors.NewPublicSystemError(fmt.Errorf("failed to unmarshal value for request %s: %w", consensusRequestMetaData.RequestID(), err),
				caperrors.Internal)
		}

		c.lggr.Debugw("returning consensus response", "metadata", metadata)

		responseAndMetadata := capabilities.ResponseAndMetadata[*valuespb.Value]{
			Response:         valueProto,
			ResponseMetadata: capabilities.ResponseMetadata{},
		}
		return &responseAndMetadata, nil
	}
}

func (c *consensusCapability) Report(ctx context.Context, metadata capabilities.RequestMetadata, reportRequest *sdk.ReportRequest) (*capabilities.ResponseAndMetadata[*sdk.ReportResponse], caperrors.Error) {
	ctx = metadata.ContextWithCRE(ctx)
	lggr := c.requestLggr(metadata)

	lggr.Debug("received reporting request", "metadata", metadata)

	stepReferenceID := metadata.ReferenceID
	reportID, err := toReportID(stepReferenceID)
	if err != nil {
		return nil, caperrors.NewPublicSystemError(fmt.Errorf("failed to convert step reference ID to report ID: %w", err), caperrors.Internal)
	}

	keyBundleID, err := validateReportRequest(reportRequest)
	if err != nil {
		return nil, caperrors.NewPublicUserError(fmt.Errorf("report request validation failed: %w", err), caperrors.InvalidArgument)
	}

	consensusRequestMetaData := oracle.ConsensusRequestMetadata{
		RequestMetadata: metadata,
		KeyBundleID:     keyBundleID,
		ReportID:        reportID,
		RequestType:     types.RequestType_REPORT_GENERATION,
	}

	input := &sdk.SimpleConsensusInputs{
		Observation: &sdk.SimpleConsensusInputs_Value{
			Value: values.Proto(values.NewBytes(reportRequest.EncodedPayload)),
		},
		Descriptors: &sdk.ConsensusDescriptor{
			Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
				Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL,
			},
		},
		Default: nil,
	}

	requestSize, reqSizeErr := c.validateRequestSize(ctx, lggr, consensusRequestMetaData, reportRequest)
	if reqSizeErr != nil {
		var capErr caperrors.Error
		if errors.As(reqSizeErr, &capErr) && capErr.Origin() != caperrors.OriginUser {
			// Return system errors immediately
			lggr.Errorw("failed to validate request size", "err", reqSizeErr)
			return nil, reqSizeErr
		}
		lggr.Debugw("request size validation failed with user error, sending as error input to consensus",
			"err", reqSizeErr,
			"requestSize", requestSize)

		input = &sdk.SimpleConsensusInputs{
			Observation: &sdk.SimpleConsensusInputs_Error{
				Error: reqSizeErr.Error(),
			},
			Descriptors: &sdk.ConsensusDescriptor{
				Descriptor_: &sdk.ConsensusDescriptor_Aggregation{
					Aggregation: sdk.AggregationType_AGGREGATION_TYPE_IDENTICAL,
				},
			},
			Default: nil,
		}
	}

	c.metrics.RecordRequestSize(ctx, float64(requestSize))

	callbackChan := c.sendRequest(ctx, input, consensusRequestMetaData)

	select {
	case <-ctx.Done():
		return nil, caperrors.NewPublicSystemError(ctx.Err(), caperrors.Canceled)
	case response := <-callbackChan:
		if response.Err != nil {
			return nil, response.Err
		}

		var sigs []*sdk.AttributedSignature

		for _, s := range response.Sigs {
			sigs = append(sigs, &sdk.AttributedSignature{
				Signature: s.Signature,
				SignerId:  uint32(s.Signer),
			})
		}

		c.lggr.Debugw("returning report", "metadata", metadata)

		reportResponse := &sdk.ReportResponse{
			ConfigDigest:  response.ConfigDigest[:],
			SeqNr:         response.SeqNr,
			ReportContext: response.ReportContext,
			RawReport:     response.RawReport,
			Sigs:          sigs,
		}
		responseAndMetadata := capabilities.ResponseAndMetadata[*sdk.ReportResponse]{
			Response:         reportResponse,
			ResponseMetadata: capabilities.ResponseMetadata{},
		}
		return &responseAndMetadata, nil
	}
}

func (c *consensusCapability) sendRequest(ctx context.Context, input *sdk.SimpleConsensusInputs, consensusRequestMetaData oracle.ConsensusRequestMetadata) <-chan oracle.ConsensusResponse {
	c.requestTimeoutLock.RLock()
	requestTimeout := c.requestTimeout
	c.requestTimeoutLock.RUnlock()

	callbackChan := make(chan oracle.ConsensusResponse, 1)

	c.reqHandler.SendRequest(ctx,
		oracle.NewConsensusRequest(input, time.Now(), time.Now().Add(requestTimeout), callbackChan,
			consensusRequestMetaData,
		))
	return callbackChan
}

func validateReportRequest(reportRequest *sdk.ReportRequest) (string, error) {
	supportedKeyBundleIDs := map[string]struct{}{
		KeyBundleIDEvm:    {},
		KeyBundleIDAptos:  {},
		KeyBundleIDSolana: {},
	}

	keyBundleID := strings.ToLower(reportRequest.EncoderName)
	if _, ok := supportedKeyBundleIDs[keyBundleID]; !ok {
		return "", fmt.Errorf("unsupported encoder name '%s' for report request. Supported encoder names are: %v", reportRequest.EncoderName, maps.Keys(supportedKeyBundleIDs))
	}

	signingAlgo := strings.ToLower(reportRequest.SigningAlgo)
	supportedSigningAlgorithms := map[string]struct{}{
		SigningAlgoEcdsa:   {},
		SigningAlgoEd25519: {},
	}

	if _, ok := supportedSigningAlgorithms[signingAlgo]; !ok {
		return "", fmt.Errorf("unsupported signing algorithm '%s' for report request. Supported signing algorithms are: %v", reportRequest.SigningAlgo, maps.Keys(supportedSigningAlgorithms))
	}

	hashingAlgo := strings.ToLower(reportRequest.HashingAlgo)

	supportedHashingAlgorithms := map[string]struct{}{
		HashingAlgoKeccak256:  {},
		HashingAlgoBlake2b256: {},
	}

	if _, ok := supportedHashingAlgorithms[hashingAlgo]; !ok {
		return "", fmt.Errorf("unsupported hashing algorithm '%s' for report request. Supported hashing algorithms are: %v", reportRequest.HashingAlgo, maps.Keys(supportedHashingAlgorithms))
	}
	return keyBundleID, nil
}

func toReportID(stepReferenceID string) (string, error) {
	// Parse the step reference to an int
	if stepReferenceID == "" {
		return "", fmt.Errorf("step reference ID is required for report request")
	}

	stepReferenceAsInt, err := strconv.Atoi(stepReferenceID)
	if err != nil {
		return "", fmt.Errorf("failed to parse step reference ID '%s' to int: %w", stepReferenceID, err)
	}

	stepRefAsHex := fmt.Sprintf("%x", stepReferenceAsInt)

	// Pad it with leading zeros so the string is at least 4 characters long
	if len(stepRefAsHex) < 4 {
		padding := 4
		stepRefAsHex = fmt.Sprintf("%0*s", padding, stepRefAsHex) // Pad with zeros to the left
	}

	b, err := hex.DecodeString(stepRefAsHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode step reference ID '%s' to bytes: %w", stepReferenceID, err)
	}

	if len(b) != 2 {
		return "", fmt.Errorf("step reference ID '%s' must be exactly 2 bytes long when encoded as hex, got %d bytes", stepReferenceID, len(b))
	}
	return stepRefAsHex, nil
}

func (c *consensusCapability) requestLggr(metadata capabilities.RequestMetadata) logger.Logger {
	lggr := logger.With(
		c.lggr,
		"workflowID", metadata.WorkflowID,
		"executionID", metadata.WorkflowExecutionID,
		"stepReferenceID", metadata.ReferenceID,
		"workflowName", metadata.DecodedWorkflowName,
	)
	return lggr
}

func (c *consensusCapability) SendResponse(ctx context.Context, response oracle.ConsensusResponse) {
	c.reqHandler.SendResponse(ctx, response)
}

// Start is not called when running as remote standard capability, instead Initialise is called  (note Close is called)
func (c *consensusCapability) Start(ctx context.Context) error {
	return nil
}

func (c *consensusCapability) Close() error {
	err := c.reqHandler.Close()
	if err != nil {
		c.lggr.Errorw("error closing request handler", "err", err)
	}

	if c.oracle != nil {
		if err := c.oracle.Close(context.Background()); err != nil {
			return fmt.Errorf("error when closing oracle: %w", err)
		}
		c.oracle = nil
	}

	return nil
}

func (c *consensusCapability) Ready() error {
	return nil
}

func (c *consensusCapability) HealthReport() map[string]error {
	return map[string]error{c.Name(): nil}
}

func (c *consensusCapability) Name() string {
	return c.lggr.Name()
}

func (c *consensusCapability) Description() string {
	return "Consensus Capability"
}

// validateRequestSize ensures the combined size of input and metadata does not exceed the allowed limit.
// This prevents oversized requests that could disrupt the consensus process.
func (c *consensusCapability) validateRequestSize(ctx context.Context, lggr logger.Logger, consensusRequestMetaData oracle.ConsensusRequestMetadata, input proto.Message) (int, caperrors.Error) {
	lggr.Debugw("starting request size validation",
		"requestID", consensusRequestMetaData.RequestID(),
		"workflowID", consensusRequestMetaData.WorkflowID)

	requestMetaData := plugin.ToRequestMetaData(consensusRequestMetaData)

	serialisedInput, err := proto.Marshal(input)
	if err != nil {
		return 0, caperrors.NewPublicSystemError(fmt.Errorf("failed to serialise input: %w", err), caperrors.Internal)
	}
	inputSize := len(serialisedInput)

	serialisedMetadata, err := proto.Marshal(requestMetaData)
	if err != nil {
		return 0, caperrors.NewPublicSystemError(fmt.Errorf("failed to serialise metadata: %w", err), caperrors.Internal)
	}
	metadataSize := len(serialisedMetadata)

	totalSize := config.Size(inputSize + metadataSize)
	lggr.Debugw("calculated total size",
		"inputSizeBytes", inputSize,
		"metadataSizeBytes", metadataSize,
		"totalSizeBytes", int(totalSize))

	// Get the limit to log it
	limit, limitErr := c.maxRequestSizeBytes.Limit(ctx)
	if limitErr != nil {
		lggr.Warnw("failed to get limit, proceeding with check", "err", limitErr)
	} else {
		lggr.Debugw("retrieved limit", "limitBytes", int(limit))
	}

	err = c.maxRequestSizeBytes.Check(ctx, totalSize)
	if err != nil {
		if errors.Is(err, limits.ErrorBoundLimited[config.Size]{}) {
			var limitErr limits.ErrorBoundLimited[config.Size]
			if errors.As(err, &limitErr) {
				lggr.Warnw("request size exceeds limit",
					"totalSizeBytes", int(totalSize),
					"limitBytes", int(limitErr.Limit),
					"excessBytes", int(totalSize-limitErr.Limit),
					"inputSizeBytes", inputSize,
					"metadataSizeBytes", metadataSize)
			} else {
				lggr.Warnw("request size exceeds limit",
					"totalSizeBytes", int(totalSize),
					"err", err)
			}
			return int(totalSize), caperrors.NewPublicUserError(fmt.Errorf("request size %d bytes exceeds maximum allowed size: %w", totalSize, err), caperrors.InvalidArgument)
		}
		return int(totalSize), caperrors.NewPublicSystemError(fmt.Errorf("unexpected error checking request size limit: %w", err), caperrors.Internal)
	}

	lggr.Debugw("request size validation passed",
		"totalSizeBytes", int(totalSize),
		"limitBytes", int(limit))
	return int(totalSize), nil
}

func decodeObservationType(lggr logger.Logger, input *sdk.SimpleConsensusInputs) error {
	switch obs := input.GetObservation().(type) {
	case *sdk.SimpleConsensusInputs_Value:
		_, err := values.FromProto(obs.Value)
		if err != nil {
			lggr.Debugw("failed to decode observation value", "err", err)
			return fmt.Errorf("failed to decode observation value: %w", err)
		}
		lggr.Debugw("received observation value")
	case *sdk.SimpleConsensusInputs_Error:
		lggr.Debugw("observation is an error")
	default:
		if input.Default != nil {
			val, err := values.FromProto(input.Default)
			if err != nil {
				lggr.Debugw("failed to decode default observation value", "err", err)
				return fmt.Errorf("failed to decode default observation value: %w", err)
			}
			lggr.Debugw("serialised and using default value", "default_value", val)
		} else {
			lggr.Debugw("neither value, error or default is set in the observation input for request")
		}
	}
	return nil
}
