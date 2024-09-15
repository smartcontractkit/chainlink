package streams

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	ocr2keepers "github.com/smartcontractkit/chainlink-common/pkg/types/automation"

	autov2common "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_automation_v21_plus_common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	v02 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/v02"
	v03 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury/v03"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/prommetrics"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	zeroAddress = "0x0000000000000000000000000000000000000000"
)

type Lookup interface {
	Lookup(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult
}

type latestBlockProvider interface {
	LatestBlock() *ocr2keepers.BlockKey
}

type streamRegistry interface {
	GetUpkeepPrivilegeConfig(opts *bind.CallOpts, upkeepId *big.Int) ([]byte, error)
	CheckCallback(opts *bind.CallOpts, id *big.Int, values [][]byte, extraData []byte) (autov2common.CheckCallback, error)
	Address() common.Address
}

type contextCaller interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type streams struct {
	services.StateMachine
	packer          mercury.Packer
	mercuryConfig   mercury.MercuryConfigProvider
	abi             abi.ABI
	blockSubscriber latestBlockProvider
	registry        streamRegistry
	client          contextCaller
	lggr            logger.Logger
	threadCtrl      utils.ThreadControl
	v02Client       mercury.MercuryClient
	v03Client       mercury.MercuryClient
}

// UpkeepPrivilegeConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepPrivilegeManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type UpkeepPrivilegeConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

func NewStreamsLookup(
	mercuryConfig mercury.MercuryConfigProvider,
	blockSubscriber latestBlockProvider,
	client contextCaller,
	registry streamRegistry,
	lggr logger.Logger) *streams {
	httpClient := http.DefaultClient
	threadCtrl := utils.NewThreadControl()
	packer := mercury.NewAbiPacker()

	return &streams{
		packer:          packer,
		mercuryConfig:   mercuryConfig,
		abi:             core.AutoV2CommonABI,
		blockSubscriber: blockSubscriber,
		registry:        registry,
		client:          client,
		lggr:            lggr,
		threadCtrl:      threadCtrl,
		v02Client:       v02.NewClient(mercuryConfig, httpClient, threadCtrl, lggr),
		v03Client:       v03.NewClient(mercuryConfig, httpClient, threadCtrl, lggr),
	}
}

// Lookup looks through check upkeep results to find any that needs off chain lookup
func (s *streams) Lookup(ctx context.Context, checkResults []ocr2keepers.CheckResult) []ocr2keepers.CheckResult {
	lookups := map[int]*mercury.StreamsLookup{}
	for i, checkResult := range checkResults {
		s.buildResult(ctx, i, checkResult, checkResults, lookups)
	}

	var wg sync.WaitGroup
	for i, lookup := range lookups {
		wg.Add(1)
		func(i int, lookup *mercury.StreamsLookup) {
			s.threadCtrl.GoCtx(ctx, func(ctx context.Context) {
				s.doLookup(ctx, &wg, lookup, i, checkResults)
			})
		}(i, lookup)
	}
	wg.Wait()

	// don't surface error to plugin bc StreamsLookup process should be self-contained.
	return checkResults
}

// buildResult checks if the upkeep is allowed by Mercury and builds a streams lookup request from the check result
func (s *streams) buildResult(ctx context.Context, i int, checkResult ocr2keepers.CheckResult, checkResults []ocr2keepers.CheckResult, lookups map[int]*mercury.StreamsLookup) {
	lookupLggr := logger.Sugared(s.lggr).With("where", "StreamsLookup")
	if checkResult.IneligibilityReason != uint8(encoding.UpkeepFailureReasonTargetCheckReverted) {
		// Streams Lookup only works when upkeep target check reverts
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorReasonNotReverted).Inc()
		return
	}

	block := big.NewInt(int64(checkResult.Trigger.BlockNumber))
	upkeepId := checkResult.UpkeepID

	// Try to decode the revert error into streams lookup format. User upkeeps can revert with any reason, see if they
	// tried to call mercury
	lookupLggr.Infof("at block %d upkeep %s trying to DecodeStreamsLookupRequest performData=%s", block, upkeepId, hexutil.Encode(checkResults[i].PerformData))
	streamsLookupErr, err := s.packer.DecodeStreamsLookupRequest(checkResult.PerformData)
	if err != nil {
		lookupLggr.Debugf("at block %d upkeep %s DecodeStreamsLookupRequest failed: %v", block, upkeepId, err)
		// user contract did not revert with StreamsLookup error
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorDecodeRequestFailed).Inc()
		return
	}
	streamsLookupResponse := &mercury.StreamsLookup{StreamsLookupError: streamsLookupErr}
	if s.mercuryConfig.Credentials() == nil {
		lookupLggr.Errorf("at block %d upkeep %s tries to access mercury server but mercury credential is not configured", block, upkeepId)
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorCredentialsNotConfigured).Inc()
		return
	}

	// mercury permission checking for v0.3 is done by mercury server, so no need to check here
	if streamsLookupResponse.IsMercuryV02() {
		// check permission on the registry for mercury v0.2
		opts := s.buildCallOpts(ctx, block)
		if state, reason, retryable, allowed, err := s.AllowedToUseMercury(opts, upkeepId.BigInt()); err != nil {
			lookupLggr.Warnf("at block %s upkeep %s failed to query mercury allow list: %s", block, upkeepId, err)
			checkResults[i].PipelineExecutionState = uint8(state)
			checkResults[i].IneligibilityReason = uint8(reason)
			checkResults[i].Retryable = retryable
			return
		} else if !allowed {
			lookupLggr.Debugf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
			checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonMercuryAccessNotAllowed)
			return
		}
	} else if !streamsLookupResponse.IsMercuryV03() {
		// if mercury version is not v02 or v03, set failure reason
		lookupLggr.Debugf("at block %d upkeep %s NOT allowed to query Mercury server", block, upkeepId)
		checkResults[i].IneligibilityReason = uint8(encoding.UpkeepFailureReasonInvalidRevertDataInput)
		return
	}

	streamsLookupResponse.UpkeepId = upkeepId.BigInt()
	// the block here is exclusively used to call checkCallback at this block, not to be confused with the block number
	// in the revert for mercury v0.2, which is denoted by time in the struct bc starting from v0.3, only timestamp will be supported
	streamsLookupResponse.Block = uint64(block.Int64())
	lookupLggr.Infof("at block %d upkeep %s DecodeStreamsLookupRequest feedKey=%s timeKey=%s feeds=%v time=%s extraData=%s", block, upkeepId, streamsLookupResponse.FeedParamKey, streamsLookupResponse.TimeParamKey, streamsLookupResponse.Feeds, streamsLookupResponse.Time, hexutil.Encode(streamsLookupResponse.ExtraData))
	lookups[i] = streamsLookupResponse
}

// Does the requested lookup and sets appropriate fields on checkResult[i]
func (s *streams) doLookup(ctx context.Context, wg *sync.WaitGroup, lookup *mercury.StreamsLookup, i int, checkResults []ocr2keepers.CheckResult) {
	defer wg.Done()

	values, errCode, err := s.DoMercuryRequest(ctx, lookup, checkResults, i)
	if err != nil {
		s.lggr.Errorf("at block %d upkeep %s requested time %s DoMercuryRequest err: %s", lookup.Block, lookup.UpkeepId, lookup.Time, err.Error())
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorDoMercuryRequest).Inc()
		return
	}

	if errCode != encoding.ErrCodeNil {
		err = s.CheckErrorHandler(ctx, errCode, lookup, checkResults, i)
		if err != nil {
			s.lggr.Errorf("at block %d upkeep %s requested time %s CheckErrorHandler err: %s", lookup.Block, lookup.UpkeepId, lookup.Time, err.Error())
		}
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorCodeNotNil).Inc()
		return
	}

	// Mercury request returned values or user's checkErrorhandler didn't return error, call checkCallback
	err = s.CheckCallback(ctx, values, lookup, checkResults, i)
	if err != nil {
		s.lggr.Errorf("at block %d upkeep %s requested time %s CheckCallback err: %s", lookup.Block, lookup.UpkeepId, lookup.Time, err.Error())
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorCheckCallback).Inc()
	}
}

func (s *streams) CheckCallback(ctx context.Context, values [][]byte, lookup *mercury.StreamsLookup, checkResults []ocr2keepers.CheckResult, i int) error {
	prommetrics.AutomationStreamsLookupStep.WithLabelValues(prommetrics.StreamsLookupStepCheckCallback).Inc()
	payload, err := s.abi.Pack("checkCallback", lookup.UpkeepId, values, lookup.ExtraData)
	if err != nil {
		checkResults[i].Retryable = false
		checkResults[i].PipelineExecutionState = uint8(encoding.PackUnpackDecodeFailed)
		return err
	}

	return s.makeCallbackEthCall(ctx, payload, lookup, checkResults, i)
}

// eth_call to checkCallback and checkErrorHandler and update checkResults[i] accordingly
func (s *streams) makeCallbackEthCall(ctx context.Context, payload []byte, lookup *mercury.StreamsLookup, checkResults []ocr2keepers.CheckResult, i int) error {
	var responseBytes hexutil.Bytes
	args := map[string]interface{}{
		"from": zeroAddress,
		"to":   s.registry.Address().Hex(),
		"data": hexutil.Bytes(payload),
	}

	if err := s.client.CallContext(ctx, &responseBytes, "eth_call", args, hexutil.EncodeUint64(lookup.Block)); err != nil {
		checkResults[i].Retryable = true
		checkResults[i].PipelineExecutionState = uint8(encoding.RpcFlakyFailure)
		return err
	}

	s.lggr.Infof("at block %d upkeep %s requested time %s responseBytes: %s", lookup.Block, lookup.UpkeepId, lookup.Time, hexutil.Encode(responseBytes))

	unpackCallBackState, needed, performData, failureReason, _, err := s.packer.UnpackCheckCallbackResult(responseBytes)
	if err != nil {
		checkResults[i].PipelineExecutionState = uint8(unpackCallBackState)
		return err
	}

	s.lggr.Infof("at block %d upkeep %s requested time %s returns needed: %v, failure reason: %d, perform data: %s", lookup.Block, lookup.UpkeepId, lookup.Time, needed, failureReason, hexutil.Encode(performData))

	checkResults[i].IneligibilityReason = uint8(failureReason)
	checkResults[i].Eligible = needed
	checkResults[i].PerformData = performData

	return nil
}

// Does the mercury request for the checkResult. Returns either the looked up values or an error code if something is wrong with mercury
// In case of any pipeline processing issues, returns an error and also sets approriate state on the checkResult itself
func (s *streams) DoMercuryRequest(ctx context.Context, lookup *mercury.StreamsLookup, checkResults []ocr2keepers.CheckResult, i int) ([][]byte, encoding.ErrCode, error) {
	prommetrics.AutomationStreamsLookupStep.WithLabelValues(prommetrics.StreamsLookupStepDoMercuryRequest).Inc()
	var state, values, errCode, retryable, retryInterval = encoding.NoPipelineError, [][]byte{}, encoding.ErrCodeNil, false, 0 * time.Second
	var err error
	pluginRetryKey := generatePluginRetryKey(checkResults[i].WorkID, lookup.Block)
	upkeepType := core.GetUpkeepType(checkResults[i].UpkeepID)

	if lookup.IsMercuryV02() {
		state, values, errCode, retryable, retryInterval, err = s.v02Client.DoRequest(ctx, lookup, upkeepType, pluginRetryKey)
	} else if lookup.IsMercuryV03() {
		state, values, errCode, retryable, retryInterval, err = s.v03Client.DoRequest(ctx, lookup, upkeepType, pluginRetryKey)
	}

	if err != nil {
		// Something went wrong in the pipeline processing, set the state, retry reason and return
		checkResults[i].Retryable = retryable
		checkResults[i].RetryInterval = retryInterval
		checkResults[i].PipelineExecutionState = uint8(state)
		s.lggr.Debugf("at block %d upkeep %s requested time %s doMercuryRequest err: %s", lookup.Block, lookup.UpkeepId, lookup.Time, err.Error())
		return nil, encoding.ErrCodeNil, err
	}

	if errCode != encoding.ErrCodeNil {
		s.lggr.Infof("at block %d upkeep %s requested time %s doMercuryRequest error code: %d", lookup.Block, lookup.UpkeepId, lookup.Time, errCode)
		return nil, errCode, nil
	}

	for j, v := range values {
		s.lggr.Infof("at block %d upkeep %s requested time %s doMercuryRequest values[%d]: %s", lookup.Block, lookup.UpkeepId, lookup.Time, j, hexutil.Encode(v))
	}
	return values, encoding.ErrCodeNil, nil
}

func (s *streams) CheckErrorHandler(ctx context.Context, errCode encoding.ErrCode, lookup *mercury.StreamsLookup, checkResults []ocr2keepers.CheckResult, i int) error {
	s.lggr.Debugf("at block %d upkeep %s requested time %s CheckErrorHandler error code: %d", lookup.Block, lookup.UpkeepId, lookup.Time, errCode)
	prommetrics.AutomationStreamsLookupStep.WithLabelValues(prommetrics.StreamsLookupStepCheckErrorHandler).Inc()

	userPayload, err := s.packer.PackUserCheckErrorHandler(errCode, lookup.ExtraData)
	if err != nil {
		checkResults[i].Retryable = false
		checkResults[i].PipelineExecutionState = uint8(encoding.PackUnpackDecodeFailed)
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorPackUserCheckErrorHandler).Inc()
		return err
	}

	payload, err := s.abi.Pack("executeCallback", lookup.UpkeepId, userPayload)
	if err != nil {
		checkResults[i].Retryable = false
		checkResults[i].PipelineExecutionState = uint8(encoding.PackUnpackDecodeFailed)
		prommetrics.AutomationStreamsLookupError.WithLabelValues(prommetrics.StreamsLookupErrorPackExecuteCallback).Inc()
		return err
	}

	return s.makeCallbackEthCall(ctx, payload, lookup, checkResults, i)
}

// AllowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
func (s *streams) AllowedToUseMercury(opts *bind.CallOpts, upkeepId *big.Int) (state encoding.PipelineExecutionState, reason encoding.UpkeepFailureReason, retryable bool, allow bool, err error) {
	allowed, ok := s.mercuryConfig.IsUpkeepAllowed(upkeepId.String())
	if ok {
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, allowed.(bool), nil
	}

	payload, err := s.packer.PackGetUpkeepPrivilegeConfig(upkeepId)
	if err != nil {
		// pack error, no retryable
		s.lggr.Warnf("failed to pack getUpkeepPrivilegeConfig data for upkeepId %s: %s", upkeepId, err)

		return encoding.PackUnpackDecodeFailed, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to pack upkeepId: %w", err)
	}

	var resultBytes hexutil.Bytes
	args := map[string]interface{}{
		"to":   s.registry.Address().Hex(),
		"data": hexutil.Bytes(payload),
	}

	if err = s.client.CallContext(opts.Context, &resultBytes, "eth_call", args, hexutil.EncodeBig(opts.BlockNumber)); err != nil {
		return encoding.RpcFlakyFailure, encoding.UpkeepFailureReasonNone, true, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}

	var upkeepPrivilegeConfigBytes []byte
	upkeepPrivilegeConfigBytes, err = s.packer.UnpackGetUpkeepPrivilegeConfig(resultBytes)

	if err != nil {
		return encoding.PackUnpackDecodeFailed, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to get upkeep privilege config: %v", err)
	}

	if len(upkeepPrivilegeConfigBytes) == 0 {
		s.mercuryConfig.SetUpkeepAllowed(upkeepId.String(), false, cache.DefaultExpiration)
		return encoding.NoPipelineError, encoding.UpkeepFailureReasonMercuryAccessNotAllowed, false, false, fmt.Errorf("upkeep privilege config is empty")
	}

	var privilegeConfig UpkeepPrivilegeConfig
	if err = json.Unmarshal(upkeepPrivilegeConfigBytes, &privilegeConfig); err != nil {
		return encoding.PrivilegeConfigUnmarshalError, encoding.UpkeepFailureReasonNone, false, false, fmt.Errorf("failed to unmarshal privilege config: %v", err)
	}

	s.mercuryConfig.SetUpkeepAllowed(upkeepId.String(), privilegeConfig.MercuryEnabled, cache.DefaultExpiration)

	return encoding.NoPipelineError, encoding.UpkeepFailureReasonNone, false, privilegeConfig.MercuryEnabled, nil
}

func (s *streams) buildCallOpts(ctx context.Context, block *big.Int) *bind.CallOpts {
	opts := bind.CallOpts{
		Context: ctx,
	}

	if block == nil || block.Int64() == 0 {
		if latestBlock := s.blockSubscriber.LatestBlock(); latestBlock != nil && latestBlock.Number != 0 {
			opts.BlockNumber = big.NewInt(int64(latestBlock.Number))
		}
	} else {
		opts.BlockNumber = block
	}

	return &opts
}

// generatePluginRetryKey returns a plugin retry cache key
func generatePluginRetryKey(workID string, block uint64) string {
	return workID + "|" + fmt.Sprintf("%d", block)
}

func (s *streams) Close() error {
	return s.StopOnce("streams_lookup", func() error {
		s.threadCtrl.Close()
		return nil
	})
}
