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

	"github.com/patrickmn/go-cache"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/models"
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
	State     MercuryUpkeepState
}

type Packer interface {
	UnpackCheckCallbackResult(callbackResp []byte) (uint8, bool, []byte, uint8, *big.Int, error)
	PackGetUpkeepPrivilegeConfig(upkeepId *big.Int) ([]byte, error)
	UnpackGetUpkeepPrivilegeConfig(resp []byte) ([]byte, error)
	DecodeStreamsLookupRequest(data []byte) (*StreamsLookupError, error)
}

type MercuryConfigProvider interface {
	Credentials() *models.MercuryCredentials
	IsUpkeepAllowed(string) (interface{}, bool)
	SetUpkeepAllowed(string, interface{}, time.Duration)
	GetPluginRetry(string) (interface{}, bool)
	SetPluginRetry(string, interface{}, time.Duration)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MercuryClient interface {
	DoRequest(ctx context.Context, streamsLookup *StreamsLookup, pluginRetryKey string) (MercuryUpkeepState, MercuryUpkeepFailureReason, [][]byte, bool, time.Duration, error)
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

func (l *StreamsLookup) IsMercuryV03UsingBlockNumber() bool {
	return l.TimeParamKey == BlockNumber
}
