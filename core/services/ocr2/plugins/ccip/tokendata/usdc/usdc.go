package usdc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/http"
)

const (
	apiVersion                = "v1"
	attestationPath           = "attestations"
	defaultAttestationTimeout = 5 * time.Second

	// defaultCoolDownDurationSec defines the default time to wait after getting rate limited.
	// this value is only used if the 429 response does not contain the Retry-After header
	defaultCoolDownDuration = 5 * time.Minute

	// maxCoolDownDuration defines the maximum duration we can wait till firing the next request
	maxCoolDownDuration = 10 * time.Minute

	// defaultRequestInterval defines the rate in requests per second that the attestation API can be called.
	// this is set according to the APIs documentated 10 requests per second rate limit.
	defaultRequestInterval = 100 * time.Millisecond

	// APIIntervalRateLimitDisabled is a special value to disable the rate limiting.
	APIIntervalRateLimitDisabled = -1
	// APIIntervalRateLimitDefault is a special value to select the default rate limit interval.
	APIIntervalRateLimitDefault = 0
)

type attestationStatus string

const (
	attestationStatusSuccess attestationStatus = "complete"
	attestationStatusPending attestationStatus = "pending_confirmations"
)

var (
	ErrUnknownResponse = errors.New("unexpected response from attestation API")
)

// messageAndAttestation has to match the onchain struct `MessageAndAttestation` in the
// USDC token pool.
type messageAndAttestation struct {
	Message     []byte
	Attestation []byte
}

func (m messageAndAttestation) AbiString() string {
	return `
	[{
		"components": [
			{"name": "message", "type": "bytes"},
			{"name": "attestation", "type": "bytes"}
		],
		"type": "tuple"
	}]`
}

func (m messageAndAttestation) Validate() error {
	if len(m.Message) == 0 {
		return errors.New("message must be non-empty")
	}
	if len(m.Attestation) == 0 {
		return errors.New("attestation must be non-empty")
	}
	return nil
}

type TokenDataReader struct {
	lggr                  logger.Logger
	usdcReader            ccipdata.USDCReader
	httpClient            http.IHttpClient
	attestationApi        *url.URL
	attestationApiTimeout time.Duration
	usdcTokenAddress      common.Address
	rate                  *rate.Limiter

	// coolDownUntil defines whether requests are blocked or not.
	coolDownUntil time.Time
	coolDownMu    *sync.RWMutex
}

type attestationResponse struct {
	Status      attestationStatus `json:"status"`
	Attestation string            `json:"attestation"`
	Error       string            `json:"error"`
}

var _ tokendata.Reader = &TokenDataReader{}

func NewUSDCTokenDataReader(
	lggr logger.Logger,
	usdcReader ccipdata.USDCReader,
	usdcAttestationApi *url.URL,
	usdcAttestationApiTimeoutSeconds int,
	usdcTokenAddress common.Address,
	requestInterval time.Duration,
) *TokenDataReader {
	timeout := time.Duration(usdcAttestationApiTimeoutSeconds) * time.Second
	if usdcAttestationApiTimeoutSeconds == 0 {
		timeout = defaultAttestationTimeout
	}

	if requestInterval == APIIntervalRateLimitDisabled {
		requestInterval = 0
	} else if requestInterval == APIIntervalRateLimitDefault {
		requestInterval = defaultRequestInterval
	}

	return &TokenDataReader{
		lggr:                  lggr,
		usdcReader:            usdcReader,
		httpClient:            http.NewObservedUsdcIHttpClient(&http.HttpClient{}),
		attestationApi:        usdcAttestationApi,
		attestationApiTimeout: timeout,
		usdcTokenAddress:      usdcTokenAddress,
		coolDownMu:            &sync.RWMutex{},
		rate:                  rate.NewLimiter(rate.Every(requestInterval), 1),
	}
}

func NewUSDCTokenDataReaderWithHttpClient(
	origin TokenDataReader,
	httpClient http.IHttpClient,
	usdcTokenAddress common.Address,
	requestInterval time.Duration,
) *TokenDataReader {
	return &TokenDataReader{
		lggr:                  origin.lggr,
		usdcReader:            origin.usdcReader,
		httpClient:            httpClient,
		attestationApi:        origin.attestationApi,
		attestationApiTimeout: origin.attestationApiTimeout,
		coolDownMu:            origin.coolDownMu,
		usdcTokenAddress:      usdcTokenAddress,
		rate:                  rate.NewLimiter(rate.Every(requestInterval), 1),
	}
}

// ReadTokenData queries the USDC attestation API to construct a message and
// attestation response. When called back to back, or multiple times
// concurrently, responses are delayed according how the request interval is
// configured.
func (s *TokenDataReader) ReadTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) ([]byte, error) {
	if tokenIndex < 0 || tokenIndex >= len(msg.TokenAmounts) {
		return nil, errors.New("token index out of bounds")
	}

	if s.inCoolDownPeriod() {
		// rate limiting cool-down period, we prevent new requests from being sent
		return nil, tokendata.ErrRequestsBlocked
	}

	if s.rate != nil {
		// Wait blocks until it the attestation API can be called or the
		// context is Done.
		if waitErr := s.rate.Wait(ctx); waitErr != nil {
			return nil, fmt.Errorf("usdc rate limiting error: %w", waitErr)
		}
	}

	messageBody, err := s.getUSDCMessageBody(ctx, msg, tokenIndex)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed getting the USDC message body")
	}

	msgID := hexutil.Encode(msg.MessageID[:])
	msgBody := hexutil.Encode(messageBody)
	s.lggr.Infow("Calling attestation API", "messageBodyHash", msgBody, "messageID", msgID)

	// The attestation API expects the hash of the message body
	attestationResp, err := s.callAttestationApi(ctx, utils.Keccak256Fixed(messageBody))
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed calling usdc attestation API ")
	}

	s.lggr.Infow("Got response from attestation API", "messageID", msgID,
		"attestationStatus", attestationResp.Status, "attestation", attestationResp.Attestation,
		"attestationError", attestationResp.Error)

	switch attestationResp.Status {
	case attestationStatusSuccess:
		// The USDC pool needs a combination of the message body and the attestation
		messageAndAttestation, err := encodeMessageAndAttestation(messageBody, attestationResp.Attestation)
		if err != nil {
			return nil, fmt.Errorf("failed to encode messageAndAttestation : %w", err)
		}
		return messageAndAttestation, nil
	case attestationStatusPending:
		return nil, tokendata.ErrNotReady
	default:
		s.lggr.Errorw("Unexpected response from attestation API", "attestationResp", attestationResp)
		return nil, ErrUnknownResponse
	}
}

// encodeMessageAndAttestation encodes the message body and attestation into a single byte array
// that is readable onchain.
func encodeMessageAndAttestation(messageBody []byte, attestation string) ([]byte, error) {
	attestationBytes, err := hex.DecodeString(strings.TrimPrefix(attestation, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to decode response attestation: %w", err)
	}

	return abihelpers.EncodeAbiStruct[messageAndAttestation](messageAndAttestation{
		Message:     messageBody,
		Attestation: attestationBytes,
	})
}

func (s *TokenDataReader) getUSDCMessageBody(
	ctx context.Context,
	msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta,
	tokenIndex int,
) ([]byte, error) {
	usdcTokenEndOffset, err := s.getUsdcTokenEndOffset(msg, tokenIndex)
	if err != nil {
		return nil, fmt.Errorf("get usdc token %d end offset: %w", tokenIndex, err)
	}

	parsedMsgBody, err := s.usdcReader.GetUSDCMessagePriorToLogIndexInTx(ctx, int64(msg.LogIndex), usdcTokenEndOffset, msg.TxHash)
	if err != nil {
		return []byte{}, err
	}

	s.lggr.Infow("Got USDC message body", "messageBody", hexutil.Encode(parsedMsgBody), "messageID", hexutil.Encode(msg.MessageID[:]))
	return parsedMsgBody, nil
}

func (s *TokenDataReader) getUsdcTokenEndOffset(msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) (int, error) {
	if tokenIndex >= len(msg.TokenAmounts) || tokenIndex < 0 {
		return 0, fmt.Errorf("invalid token index %d for msg with %d tokens", tokenIndex, len(msg.TokenAmounts))
	}

	if msg.TokenAmounts[tokenIndex].Token != ccipcalc.EvmAddrToGeneric(s.usdcTokenAddress) {
		return 0, fmt.Errorf("the specified token index %d is not a usdc token", tokenIndex)
	}

	usdcTokenEndOffset := 0
	for i := tokenIndex + 1; i < len(msg.TokenAmounts); i++ {
		evmTokenAddr, err := ccipcalc.GenericAddrToEvm(msg.TokenAmounts[i].Token)
		if err != nil {
			continue
		}

		if evmTokenAddr == s.usdcTokenAddress {
			usdcTokenEndOffset++
		}
	}

	return usdcTokenEndOffset, nil
}

// callAttestationApi calls the USDC attestation API with the given USDC message hash.
// The attestation service rate limit is 10 requests per second. If you exceed 10 requests
// per second, the service blocks all API requests for the next 5 minutes and returns an
// HTTP 429 response.
//
// Documentation:
//
//	https://developers.circle.com/stablecoins/reference/getattestation
//	https://developers.circle.com/stablecoins/docs/transfer-usdc-on-testnet-from-ethereum-to-avalanche
func (s *TokenDataReader) callAttestationApi(ctx context.Context, usdcMessageHash [32]byte) (attestationResponse, error) {
	body, _, headers, err := s.httpClient.Get(
		ctx,
		fmt.Sprintf("%s/%s/%s/0x%x", s.attestationApi, apiVersion, attestationPath, usdcMessageHash),
		s.attestationApiTimeout,
	)
	switch {
	case errors.Is(err, tokendata.ErrRateLimit):
		coolDownDuration := defaultCoolDownDuration
		if retryAfterHeader, exists := headers["Retry-After"]; exists && len(retryAfterHeader) > 0 {
			if retryAfterSec, errParseInt := strconv.ParseInt(retryAfterHeader[0], 10, 64); errParseInt == nil {
				coolDownDuration = time.Duration(retryAfterSec) * time.Second
			}
		}
		s.setCoolDownPeriod(coolDownDuration)

		// Explicitly signal if the API is being rate limited
		return attestationResponse{}, tokendata.ErrRateLimit
	case err != nil:
		return attestationResponse{}, fmt.Errorf("request error: %w", err)
	}

	var response attestationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return attestationResponse{}, err
	}
	if response.Error != "" {
		return attestationResponse{}, fmt.Errorf("attestation API error: %s", response.Error)
	}
	if response.Status == "" {
		return attestationResponse{}, fmt.Errorf("invalid attestation response: %s", string(body))
	}
	return response, nil
}

func (s *TokenDataReader) setCoolDownPeriod(d time.Duration) {
	s.coolDownMu.Lock()
	if d > maxCoolDownDuration {
		d = maxCoolDownDuration
	}
	s.coolDownUntil = time.Now().Add(d)
	s.coolDownMu.Unlock()
}

func (s *TokenDataReader) inCoolDownPeriod() bool {
	s.coolDownMu.RLock()
	defer s.coolDownMu.RUnlock()
	return time.Now().Before(s.coolDownUntil)
}

func (s *TokenDataReader) Close() error {
	return nil
}
