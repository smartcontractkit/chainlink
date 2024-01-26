package mercury

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/core"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
)

const (
	FeedIDs                  = "feedIDs"     // valid for v0.3
	FeedIdHex                = "feedIdHex"   // valid for v0.2
	BlockNumber              = "blockNumber" // valid for v0.2
	Timestamp                = "timestamp"   // valid for v0.3
	totalFastPluginRetries   = 5
	totalMediumPluginRetries = 10
)

var GenerateHMACFn = func(method string, path string, body []byte, clientId string, secret string, ts int64) string {
	bodyHash := sha256.New()
	bodyHash.Write(body)
	hashString := fmt.Sprintf("%s %s %s %s %d",
		method,
		path,
		hex.EncodeToString(bodyHash.Sum(nil)),
		clientId,
		ts)
	signedMessage := hmac.New(sha256.New, []byte(secret))
	signedMessage.Write([]byte(hashString))
	userHmac := hex.EncodeToString(signedMessage.Sum(nil))
	return userHmac
}

// CalculateRetryConfig returns plugin retry interval based on how many times plugin has retried this work
var CalculateRetryConfigFn = func(prk string, mercuryConfig MercuryConfigProvider) time.Duration {
	var retryInterval time.Duration
	var retries int
	totalAttempts, ok := mercuryConfig.GetPluginRetry(prk)
	if ok {
		retries = totalAttempts.(int)
		if retries < totalFastPluginRetries {
			retryInterval = 1 * time.Second
		} else if retries < totalMediumPluginRetries {
			retryInterval = 5 * time.Second
		}
		// if the core node has retried totalMediumPluginRetries times, do not set retry interval and plugin will use
		// the default interval
	} else {
		retryInterval = 1 * time.Second
	}
	mercuryConfig.SetPluginRetry(prk, retries+1, cache.DefaultExpiration)
	return retryInterval
}

type MercuryData struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     [][]byte
	State     encoding.PipelineExecutionState
}

type MercuryConfigProvider interface {
	Credentials() *types.MercuryCredentials
	IsUpkeepAllowed(string) (interface{}, bool)
	SetUpkeepAllowed(string, interface{}, time.Duration)
	GetPluginRetry(string) (interface{}, bool)
	SetPluginRetry(string, interface{}, time.Duration)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MercuryClient interface {
	DoRequest(ctx context.Context, streamsLookup *StreamsLookup, pluginRetryKey string) (encoding.PipelineExecutionState, encoding.UpkeepFailureReason, [][]byte, bool, time.Duration, error)
}

type StreamsLookupError struct {
	FeedParamKey string
	Feeds        []string
	TimeParamKey string
	Time         *big.Int
	ExtraData    []byte
}

type StreamsLookup struct {
	*StreamsLookupError
	UpkeepId *big.Int
	Block    uint64
}

func (l *StreamsLookup) IsMercuryV02() bool {
	return l.FeedParamKey == FeedIdHex && l.TimeParamKey == BlockNumber
}

func (l *StreamsLookup) IsMercuryV03() bool {
	return l.FeedParamKey == FeedIDs
}

// IsMercuryV03UsingBlockNumber is used to distinguish the batch path. It is used for Mercury V03 only
func (l *StreamsLookup) IsMercuryV03UsingBlockNumber() bool {
	return l.TimeParamKey == BlockNumber
}

type Packer interface {
	UnpackCheckCallbackResult(callbackResp []byte) (encoding.PipelineExecutionState, bool, []byte, encoding.UpkeepFailureReason, *big.Int, error)
	PackGetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error)
	UnpackGetUpkeepPrivilegeConfig(resp []byte) ([]byte, error)
	DecodeStreamsLookupRequest(data []byte) (*StreamsLookupError, error)
}

type abiPacker struct {
	registryABI abi.ABI
	streamsABI  abi.ABI
}

func NewAbiPacker() *abiPacker {
	return &abiPacker{registryABI: core.RegistryABI, streamsABI: core.StreamsCompatibleABI}
}

// DecodeStreamsLookupRequest decodes the revert error StreamsLookup(string feedParamKey, string[] feeds, string feedParamKey, uint256 time, byte[] extraData)
func (p *abiPacker) DecodeStreamsLookupRequest(data []byte) (*StreamsLookupError, error) {
	e := p.streamsABI.Errors["StreamsLookup"]
	unpack, err := e.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("unpack error: %w", err)
	}
	errorParameters := unpack.([]interface{})

	return &StreamsLookupError{
		FeedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
		Feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
		TimeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
		Time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
		ExtraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
	}, nil
}

func (p *abiPacker) UnpackCheckCallbackResult(callbackResp []byte) (encoding.PipelineExecutionState, bool, []byte, encoding.UpkeepFailureReason, *big.Int, error) {
	out, err := p.registryABI.Methods["checkCallback"].Outputs.UnpackValues(callbackResp)
	if err != nil {
		return encoding.PackUnpackDecodeFailed, false, nil, 0, nil, fmt.Errorf("%w: unpack checkUpkeep return: %s", err, hexutil.Encode(callbackResp))
	}

	upkeepNeeded := *abi.ConvertType(out[0], new(bool)).(*bool)
	rawPerformData := *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	failureReason := encoding.UpkeepFailureReason(*abi.ConvertType(out[2], new(uint8)).(*uint8))
	gasUsed := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return encoding.NoPipelineError, upkeepNeeded, rawPerformData, failureReason, gasUsed, nil
}

func (p *abiPacker) UnpackGetUpkeepPrivilegeConfig(resp []byte) ([]byte, error) {
	out, err := p.registryABI.Methods["getUpkeepPrivilegeConfig"].Outputs.UnpackValues(resp)
	if err != nil {
		return nil, fmt.Errorf("%w: unpack getUpkeepPrivilegeConfig return", err)
	}

	bts := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return bts, nil
}

func (p *abiPacker) PackGetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error) {
	return p.registryABI.Pack("getUpkeepPrivilegeConfig", upkeepId)
}
