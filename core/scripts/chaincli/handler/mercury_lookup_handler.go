package handler

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"crypto/sha256"

	"encoding/json"

	"github.com/avast/retry-go"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

// MercuryLookupHandler is responsible for initiating the calls to the Mercury server
// to determine whether the upkeeps are eligible
type MercuryLookupHandler struct {
	credentials *MercuryCredentials
	httpClient  HttpClient
	rpcClient   *rpc.Client
}

func NewMercuryLookupHandler(
	credentials *MercuryCredentials,
	rpcClient *rpc.Client,
) *MercuryLookupHandler {
	return &MercuryLookupHandler{
		credentials: credentials,
		httpClient:  http.DefaultClient,
		rpcClient:   rpcClient,
	}
}

type MercuryVersion string

type StreamsLookup struct {
	feedParamKey string
	feeds        []string
	timeParamKey string
	time         *big.Int
	extraData    []byte
	upkeepId     *big.Int
	block        uint64
}

//go:generate mockery --quiet --name HttpClient --output ./mocks/ --case=underscore
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MercuryCredentials struct {
	LegacyURL string
	URL       string
	ClientID  string
	ClientKey string
}

func (mc *MercuryCredentials) Validate() bool {
	return mc.URL != "" && mc.ClientID != "" && mc.ClientKey != ""
}

type MercuryData struct {
	Index     int
	Error     error
	Retryable bool
	Bytes     [][]byte
	State     PipelineExecutionState
}

// MercuryV02Response represents a JSON structure used by Mercury v0.2
type MercuryV02Response struct {
	ChainlinkBlob string `json:"chainlinkBlob"`
}

// MercuryV03Response represents a JSON structure used by Mercury v0.3
type MercuryV03Response struct {
	Reports []MercuryV03Report `json:"reports"`
}

type MercuryV03Report struct {
	FeedID                string `json:"feedID"` // feed id in hex encoded
	ValidFromTimestamp    uint32 `json:"validFromTimestamp"`
	ObservationsTimestamp uint32 `json:"observationsTimestamp"`
	FullReport            string `json:"fullReport"` // the actual hex encoded mercury report of this feed, can be sent to verifier
}

const (
	// DefaultAllowListExpiration decides how long an upkeep's allow list info will be valid for.
	DefaultAllowListExpiration = 20 * time.Minute
	// CleanupInterval decides when the expired items in cache will be deleted.
	CleanupInterval = 25 * time.Minute
)

const (
	ApplicationJson     = "application/json"
	BlockNumber         = "blockNumber" // valid for v0.2
	FeedIDs             = "feedIDs"     // valid for v0.3
	FeedIdHex           = "feedIdHex"   // valid for v0.2
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderTimestamp     = "X-Authorization-Timestamp"
	HeaderSignature     = "X-Authorization-Signature-SHA256"
	HeaderUpkeepId      = "X-Authorization-Upkeep-Id"
	MercuryPathV2       = "/client?"              // only used to access mercury v0.2 server
	MercuryBatchPathV3  = "/api/v1/reports/bulk?" // only used to access mercury v0.3 server
	RetryDelay          = 500 * time.Millisecond
	Timestamp           = "timestamp" // valid for v0.3
	TotalAttempt        = 3
	UserId              = "userId"
)

type UpkeepFailureReason uint8
type PipelineExecutionState uint8

const (
	// upkeep failure onchain reasons
	UpkeepFailureReasonNone                    UpkeepFailureReason = 0
	UpkeepFailureReasonUpkeepCancelled         UpkeepFailureReason = 1
	UpkeepFailureReasonUpkeepPaused            UpkeepFailureReason = 2
	UpkeepFailureReasonTargetCheckReverted     UpkeepFailureReason = 3
	UpkeepFailureReasonUpkeepNotNeeded         UpkeepFailureReason = 4
	UpkeepFailureReasonPerformDataExceedsLimit UpkeepFailureReason = 5
	UpkeepFailureReasonInsufficientBalance     UpkeepFailureReason = 6
	UpkeepFailureReasonMercuryCallbackReverted UpkeepFailureReason = 7
	UpkeepFailureReasonRevertDataExceedsLimit  UpkeepFailureReason = 8
	UpkeepFailureReasonRegistryPaused          UpkeepFailureReason = 9
	// leaving a gap here for more onchain failure reasons in the future
	// upkeep failure offchain reasons
	UpkeepFailureReasonMercuryAccessNotAllowed UpkeepFailureReason = 32
	UpkeepFailureReasonTxHashNoLongerExists    UpkeepFailureReason = 33
	UpkeepFailureReasonInvalidRevertDataInput  UpkeepFailureReason = 34
	UpkeepFailureReasonSimulationFailed        UpkeepFailureReason = 35
	UpkeepFailureReasonTxHashReorged           UpkeepFailureReason = 36

	// pipeline execution error
	NoPipelineError        PipelineExecutionState = 0
	CheckBlockTooOld       PipelineExecutionState = 1
	CheckBlockInvalid      PipelineExecutionState = 2
	RpcFlakyFailure        PipelineExecutionState = 3
	MercuryFlakyFailure    PipelineExecutionState = 4
	PackUnpackDecodeFailed PipelineExecutionState = 5
	MercuryUnmarshalError  PipelineExecutionState = 6
	InvalidMercuryRequest  PipelineExecutionState = 7
	InvalidMercuryResponse PipelineExecutionState = 8 // this will only happen if Mercury server sends bad responses
	UpkeepNotAuthorized    PipelineExecutionState = 9
)

// UpkeepPrivilegeConfig represents the administrative offchain config for each upkeep. It can be set by s_upkeepPrivilegeManager
// role on the registry. Upkeeps allowed to use Mercury server will have this set to true.
type UpkeepPrivilegeConfig struct {
	MercuryEnabled bool `json:"mercuryEnabled"`
}

// generateHMAC calculates a user HMAC for Mercury server authentication.
func (mlh *MercuryLookupHandler) generateHMAC(method string, path string, body []byte, clientId string, secret string, ts int64) string {
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

// singleFeedRequest sends a v0.2 Mercury request for a single feed report.
func (mlh *MercuryLookupHandler) singleFeedRequest(ctx context.Context, ch chan<- MercuryData, index int, ml *StreamsLookup) {
	q := url.Values{
		ml.feedParamKey: {ml.feeds[index]},
		ml.timeParamKey: {ml.time.String()},
	}
	mercuryURL := mlh.credentials.LegacyURL
	reqUrl := fmt.Sprintf("%s%s%s", mercuryURL, MercuryPathV2, q.Encode())
	// mlh.logger.Debugf("request URL for upkeep %s feed %s: %s", ml.upkeepId.String(), ml.feeds[index], reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: index, Error: err, Retryable: false, State: InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := mlh.generateHMAC(http.MethodGet, MercuryPathV2+q.Encode(), []byte{}, mlh.credentials.ClientID, mlh.credentials.ClientKey, ts)
	req.Header.Set(HeaderContentType, ApplicationJson)
	req.Header.Set(HeaderAuthorization, mlh.credentials.ClientID)
	req.Header.Set(HeaderTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(HeaderSignature, signature)

	// in the case of multiple retries here, use the last attempt's data
	state := NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := mlh.httpClient.Do(req)
			if err1 != nil {
				// mlh.logger.Errorw("StreamsLookup GET request failed", "upkeepID", ml.upkeepId.String(), "time", ml.time.String(), "feed", ml.feeds[index], "error", err1)
				retryable = true
				state = MercuryFlakyFailure
				return err1
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					// mlh.logger.Errorf("Encountered error when closing the body of the response in single feed: %s", err)
				}
			}(resp.Body)

			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = InvalidMercuryResponse
				return err1
			}

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				// mlh.logger.Errorw("StreamsLookup received retryable status code", "upkeepID", ml.upkeepId.String(), "time", ml.time.String(), "statusCode", resp.StatusCode, "feed", ml.feeds[index])
				retryable = true
				state = MercuryFlakyFailure
				return errors.New(strconv.FormatInt(int64(resp.StatusCode), 10))
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = InvalidMercuryRequest
				return fmt.Errorf("StreamsLookup upkeep %s block %s received status code %d for feed %s", ml.upkeepId.String(), ml.time.String(), resp.StatusCode, ml.feeds[index])
			}

			// mlh.logger.Debugf("at block %s upkeep %s received status code %d from mercury v0.2 with BODY=%s", ml.time.String(), ml.upkeepId.String(), resp.StatusCode, hexutil.Encode(body))

			var m MercuryV02Response
			err1 = json.Unmarshal(body, &m)
			if err1 != nil {
				// mlh.logger.Errorw("StreamsLookup failed to unmarshal body to MercuryResponse", "upkeepID", ml.upkeepId.String(), "time", ml.time.String(), "feed", ml.feeds[index], "error", err1)
				retryable = false
				state = MercuryUnmarshalError
				return err1
			}
			blobBytes, err1 := hexutil.Decode(m.ChainlinkBlob)
			if err1 != nil {
				// mlh.logger.Errorw("StreamsLookup failed to decode chainlinkBlob for feed", "upkeepID", ml.upkeepId.String(), "time", ml.time.String(), "blob", m.ChainlinkBlob, "feed", ml.feeds[index], "error", err1)
				retryable = false
				state = InvalidMercuryResponse
				return err1
			}
			ch <- MercuryData{
				Index:     index,
				Bytes:     [][]byte{blobBytes},
				Retryable: false,
				State:     NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	if !sent {
		md := MercuryData{
			Index:     index,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     fmt.Errorf("failed to request feed for %s: %w", ml.feeds[index], retryErr),
			State:     state,
		}
		ch <- md
	}
}

// multiFeedsRequest sends a Mercury v0.3 request for a multi-feed report
func (mlh *MercuryLookupHandler) multiFeedsRequest(ctx context.Context, ch chan<- MercuryData, ml *StreamsLookup) {
	params := fmt.Sprintf("%s=%s&%s=%s", FeedIDs, strings.Join(ml.feeds, ","), Timestamp, ml.time.String())
	reqUrl := fmt.Sprintf("%s%s%s", mlh.credentials.URL, MercuryBatchPathV3, params)
	// mlh.logger.Debugf("request URL for upkeep %s userId %s: %s", ml.upkeepId.String(), mlh.credentials.ClientID, reqUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		ch <- MercuryData{Index: 0, Error: err, Retryable: false, State: InvalidMercuryRequest}
		return
	}

	ts := time.Now().UTC().UnixMilli()
	signature := mlh.generateHMAC(http.MethodGet, MercuryBatchPathV3+params, []byte{}, mlh.credentials.ClientID, mlh.credentials.ClientKey, ts)
	req.Header.Set(HeaderContentType, ApplicationJson)
	// username here is often referred to as user id
	req.Header.Set(HeaderAuthorization, mlh.credentials.ClientID)
	req.Header.Set(HeaderTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(HeaderSignature, signature)
	// mercury will inspect authorization headers above to make sure this user (in automation's context, this node) is eligible to access mercury
	// and if it has an automation role. it will then look at this upkeep id to check if it has access to all the requested feeds.
	req.Header.Set(HeaderUpkeepId, ml.upkeepId.String())

	// in the case of multiple retries here, use the last attempt's data
	state := NoPipelineError
	retryable := false
	sent := false
	retryErr := retry.Do(
		func() error {
			retryable = false
			resp, err1 := mlh.httpClient.Do(req)
			if err1 != nil {
				// mlh.logger.Errorw("StreamsLookup GET request fails for multi feed", "upkeepID", ml.upkeepId.String(), "time", ml.time.String(), "error", err1)
				retryable = true
				state = MercuryFlakyFailure
				return err1
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					// mlh.logger.Errorf("Encountered error when closing the body of the response in the multi feed: %s", err)
				}
			}(resp.Body)
			body, err1 := io.ReadAll(resp.Body)
			if err1 != nil {
				retryable = false
				state = InvalidMercuryResponse
				return err1
			}

			// mlh.logger.Infof("at timestamp %s upkeep %s received status code %d from mercury v0.3", ml.time.String(), ml.upkeepId.String(), resp.StatusCode)
			if resp.StatusCode == http.StatusUnauthorized {
				retryable = false
				state = UpkeepNotAuthorized
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by unauthorized upkeep", ml.time.String(), ml.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode == http.StatusBadRequest {
				retryable = false
				state = InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by invalid format of timestamp", ml.time.String(), ml.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode == http.StatusInternalServerError {
				retryable = true
				state = MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusInternalServerError)
			} else if resp.StatusCode == 420 {
				// in 0.3, this will happen when missing/malformed query args, missing or bad required headers, non-existent feeds, or no permissions for feeds
				retryable = false
				state = InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3, most likely this is caused by missing/malformed query args, missing or bad required headers, non-existent feeds, or no permissions for feeds", ml.time.String(), ml.upkeepId.String(), resp.StatusCode)
			} else if resp.StatusCode != http.StatusOK {
				retryable = false
				state = InvalidMercuryRequest
				return fmt.Errorf("at timestamp %s upkeep %s received status code %d from mercury v0.3", ml.time.String(), ml.upkeepId.String(), resp.StatusCode)
			}

			var response MercuryV03Response
			err1 = json.Unmarshal(body, &response)
			if err1 != nil {
				// mlh.logger.Errorw("StreamsLookup failed to unmarshal body to MercuryResponse for multi feed", "upkeepID", ml.upkeepId.String(), "time", ml.time.String(), "error", err1)
				retryable = false
				state = MercuryUnmarshalError
				return err1
			}
			// in v0.3, if some feeds are not available, the server will only return available feeds, but we need to make sure ALL feeds are retrieved before calling user contract
			// hence, retry in this case. retry will help when we send a very new timestamp and reports are not yet generated
			if len(response.Reports) != len(ml.feeds) {
				// TODO: AUTO-5044: calculate what reports are missing and log a warning
				retryable = true
				state = MercuryFlakyFailure
				return fmt.Errorf("%d", http.StatusNotFound)
			}
			var reportBytes [][]byte
			for _, rsp := range response.Reports {
				b, err := hexutil.Decode(rsp.FullReport)
				if err != nil {
					retryable = false
					state = InvalidMercuryResponse
					return err
				}
				reportBytes = append(reportBytes, b)
			}
			ch <- MercuryData{
				Index:     0,
				Bytes:     reportBytes,
				Retryable: false,
				State:     NoPipelineError,
			}
			sent = true
			return nil
		},
		// only retry when the error is 404 Not Found or 500 Internal Server Error
		retry.RetryIf(func(err error) bool {
			return err.Error() == fmt.Sprintf("%d", http.StatusNotFound) || err.Error() == fmt.Sprintf("%d", http.StatusInternalServerError)
		}),
		retry.Context(ctx),
		retry.Delay(RetryDelay),
		retry.Attempts(TotalAttempt))

	if !sent {
		md := MercuryData{
			Index:     0,
			Bytes:     [][]byte{},
			Retryable: retryable,
			Error:     retryErr,
			State:     state,
		}
		ch <- md
	}
}

// doMercuryRequest sends requests to Mercury API to retrieve ChainlinkBlob.
func (mlh *MercuryLookupHandler) doMercuryRequest(ctx context.Context, ml *StreamsLookup) (PipelineExecutionState, UpkeepFailureReason, [][]byte, bool, error) {
	var isMercuryV03 bool
	resultLen := len(ml.feeds)
	ch := make(chan MercuryData, resultLen)
	if len(ml.feeds) == 0 {
		return NoPipelineError, UpkeepFailureReasonInvalidRevertDataInput, nil, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", ml.feedParamKey, ml.timeParamKey, ml.feeds)
	}
	if ml.feedParamKey == FeedIdHex && ml.timeParamKey == BlockNumber {
		// only v0.2
		for i := range ml.feeds {
			go mlh.singleFeedRequest(ctx, ch, i, ml)
		}
	} else if ml.feedParamKey == FeedIDs && ml.timeParamKey == Timestamp {
		// only v0.3
		resultLen = 1
		isMercuryV03 = true
		ch = make(chan MercuryData, resultLen)
		go mlh.multiFeedsRequest(ctx, ch, ml)
	} else {
		return NoPipelineError, UpkeepFailureReasonInvalidRevertDataInput, nil, false, fmt.Errorf("invalid revert data input: feed param key %s, time param key %s, feeds %s", ml.feedParamKey, ml.timeParamKey, ml.feeds)
	}

	var reqErr error
	results := make([][]byte, len(ml.feeds))
	retryable := true
	allSuccess := true
	// in v0.2, use the last execution error as the state, if no execution errors, state will be no error
	state := NoPipelineError
	for i := 0; i < resultLen; i++ {
		m := <-ch
		if m.Error != nil {
			if reqErr == nil {
				reqErr = errors.New(m.Error.Error())
			} else {
				reqErr = errors.New(reqErr.Error() + m.Error.Error())
			}
			retryable = retryable && m.Retryable
			allSuccess = false
			if m.State != NoPipelineError {
				state = m.State
			}
			continue
		}
		if isMercuryV03 {
			results = m.Bytes
		} else {
			results[m.Index] = m.Bytes[0]
		}
	}
	// only retry when not all successful AND none are not retryable
	return state, UpkeepFailureReasonNone, results, retryable && !allSuccess, reqErr
}

// decodeStreamsLookup decodes the revert error StreamsLookup(string feedParamKey, string[] feeds, string timeParamKey, uint256 time, byte[] extraData)
// func (mlh *MercuryLookupHandler) decodeStreamsLookup(data []byte) (*StreamsLookup, error) {
// 	e := mlh.mercuryConfig.Abi.Errors["StreamsLookup"]
// 	unpack, err := e.Unpack(data)
// 	if err != nil {
// 		return nil, fmt.Errorf("unpack error: %w", err)
// 	}
// 	errorParameters := unpack.([]interface{})

// 	return &StreamsLookup{
// 		feedParamKey: *abi.ConvertType(errorParameters[0], new(string)).(*string),
// 		feeds:        *abi.ConvertType(errorParameters[1], new([]string)).(*[]string),
// 		timeParamKey: *abi.ConvertType(errorParameters[2], new(string)).(*string),
// 		time:         *abi.ConvertType(errorParameters[3], new(*big.Int)).(**big.Int),
// 		extraData:    *abi.ConvertType(errorParameters[4], new([]byte)).(*[]byte),
// 	}, nil
// }

// allowedToUseMercury retrieves upkeep's administrative offchain config and decode a mercuryEnabled bool to indicate if
// this upkeep is allowed to use Mercury service.
// func (mlh *MercuryLookupHandler) allowedToUseMercury(upkeep models.Upkeep) (bool, error) {
// 	allowed, ok := mlh.mercuryConfig.AllowListCache.Get(upkeep.Admin.Hex())
// 	if ok {
// 		return allowed.(bool), nil
// 	}

// 	if upkeep.UpkeepPrivilegeConfig == nil {
// 		return false, fmt.Errorf("the upkeep privilege config was not retrieved for upkeep with ID %s", upkeep.UpkeepID)
// 	}

// 	if len(upkeep.UpkeepPrivilegeConfig) == 0 {
// 		return false, fmt.Errorf("the upkeep privilege config is empty")
// 	}

// 	var a UpkeepPrivilegeConfig
// 	err := json.Unmarshal(upkeep.UpkeepPrivilegeConfig, &a)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to unmarshal privilege config for upkeep ID %s: %v", upkeep.UpkeepID, err)
// 	}

// 	mlh.mercuryConfig.AllowListCache.Set(upkeep.Admin.Hex(), a.MercuryEnabled, cache.DefaultExpiration)
// 	return a.MercuryEnabled, nil
// }

func (mlh *MercuryLookupHandler) CheckCallback(ctx context.Context, values [][]byte, lookup *StreamsLookup, registryABI ethabi.ABI, registryAddress common.Address) (hexutil.Bytes, error) {
	payload, err := registryABI.Pack("checkCallback", lookup.upkeepId, values, lookup.extraData)
	if err != nil {
		return nil, err
	}

	var theBytes hexutil.Bytes
	args := map[string]interface{}{
		"to":   registryAddress.Hex(),
		"data": hexutil.Bytes(payload),
	}

	// call checkCallback function at the block which OCR3 has agreed upon
	err = mlh.rpcClient.CallContext(ctx, &theBytes, "eth_call", args, hexutil.EncodeUint64(lookup.block))
	if err != nil {
		return nil, err
	}
	return theBytes, nil
}
