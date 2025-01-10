package lbtc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/http"
)

const (
	apiVersion                = "v1"
	attestationPath           = "deposits/getByHash"
	defaultAttestationTimeout = 5 * time.Second

	// defaultCoolDownDurationSec defines the default time to wait after getting rate limited.
	// this value is only used if the 429 response does not contain the Retry-After header
	defaultCoolDownDuration = 30 * time.Second

	// defaultRequestInterval defines the rate in requests per second that the attestation API can be called.
	// this is set according to the APIs recommended 5 requests per second rate limit.
	defaultRequestInterval = 200 * time.Millisecond

	// APIIntervalRateLimitDisabled is a special value to disable the rate limiting.
	APIIntervalRateLimitDisabled = -1
	// APIIntervalRateLimitDefault is a special value to select the default rate limit interval.
	APIIntervalRateLimitDefault = 0
)

type attestationStatus string

const (
	attestationStatusUnspecified     attestationStatus = "NOTARIZATION_STATUS_UNSPECIFIED"
	attestationStatusPending         attestationStatus = "NOTARIZATION_STATUS_PENDING"
	attestationStatusSubmitted       attestationStatus = "NOTARIZATION_STATUS_SUBMITTED"
	attestationStatusSessionApproved attestationStatus = "NOTARIZATION_STATUS_SESSION_APPROVED"
	attestationStatusFailed          attestationStatus = "NOTARIZATION_STATUS_FAILED"
)

var (
	ErrUnknownResponse = errors.New("unexpected response from attestation API")
)

type TokenDataReader struct {
	lggr                  logger.Logger
	httpClient            http.IHttpClient
	attestationApi        *url.URL
	attestationApiTimeout time.Duration
	lbtcTokenAddress      common.Address
	rate                  *rate.Limiter

	// coolDownUntil defines whether requests are blocked or not.
	coolDownUntil time.Time
	coolDownMu    *sync.RWMutex
}

type messageAttestationResponse struct {
	MessageHash string            `json:"message_hash"`
	Status      attestationStatus `json:"status"`
	Attestation string            `json:"attestation,omitempty"` // Attestation represented by abi.encode(payload, proof)
}

type attestationRequest struct {
	PayloadHashes []string `json:"messageHash"`
}

type attestationResponse struct {
	Attestations []messageAttestationResponse `json:"attestations"`
}

type sourceTokenData struct {
	SourcePoolAddress []byte
	DestTokenAddress  []byte
	ExtraData         []byte
	DestGasAmount     uint32
}

func (m sourceTokenData) AbiString() string {
	return `[{
		"components": [
			{"name": "sourcePoolAddress", "type": "bytes"},
			{"name": "destTokenAddress", "type": "bytes"},
			{"name": "extraData", "type": "bytes"},
			{"name": "destGasAmount", "type": "uint32"}
		],
		"type": "tuple"
	}]`
}

func (m sourceTokenData) Validate() error {
	if len(m.SourcePoolAddress) == 0 {
		return errors.New("sourcePoolAddress must be non-empty")
	}
	if len(m.DestTokenAddress) == 0 {
		return errors.New("destTokenAddress must be non-empty")
	}
	if len(m.ExtraData) == 0 {
		return errors.New("extraData must be non-empty")
	}
	return nil
}

var _ tokendata.Reader = &TokenDataReader{}

func NewLBTCTokenDataReader(
	lggr logger.Logger,
	lbtcAttestationApi *url.URL,
	lbtcAttestationApiTimeoutSeconds int,
	lbtcTokenAddress common.Address,
	requestInterval time.Duration,
) *TokenDataReader {
	timeout := time.Duration(lbtcAttestationApiTimeoutSeconds) * time.Second
	if lbtcAttestationApiTimeoutSeconds == 0 {
		timeout = defaultAttestationTimeout
	}

	if requestInterval == APIIntervalRateLimitDisabled {
		requestInterval = 0
	} else if requestInterval == APIIntervalRateLimitDefault {
		requestInterval = defaultRequestInterval
	}

	return &TokenDataReader{
		lggr:                  lggr,
		httpClient:            http.NewObservedLbtcIHttpClient(&http.HttpClient{}),
		attestationApi:        lbtcAttestationApi,
		attestationApiTimeout: timeout,
		lbtcTokenAddress:      lbtcTokenAddress,
		coolDownMu:            &sync.RWMutex{},
		rate:                  rate.NewLimiter(rate.Every(requestInterval), 1),
	}
}

func NewLBTCTokenDataReaderWithHttpClient(
	origin TokenDataReader,
	httpClient http.IHttpClient,
	lbtcTokenAddress common.Address,
	requestInterval time.Duration,
) *TokenDataReader {
	return &TokenDataReader{
		lggr:                  origin.lggr,
		httpClient:            httpClient,
		attestationApi:        origin.attestationApi,
		attestationApiTimeout: origin.attestationApiTimeout,
		coolDownMu:            origin.coolDownMu,
		lbtcTokenAddress:      lbtcTokenAddress,
		rate:                  rate.NewLimiter(rate.Every(requestInterval), 1),
	}
}

// ReadTokenData queries the LBTC attestation API.
func (s *TokenDataReader) ReadTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) ([]byte, error) {
	if tokenIndex < 0 || tokenIndex >= len(msg.TokenAmounts) {
		return nil, fmt.Errorf("token index out of bounds")
	}

	if s.inCoolDownPeriod() {
		// rate limiting cool-down period, we prevent new requests from being sent
		return nil, tokendata.ErrRequestsBlocked
	}

	if s.rate != nil {
		// Wait blocks until it the attestation API can be called or the
		// context is Done.
		if waitErr := s.rate.Wait(ctx); waitErr != nil {
			return nil, fmt.Errorf("lbtc rate limiting error: %w", waitErr)
		}
	}

	decodedSourceTokenData, err := abihelpers.DecodeAbiStruct[sourceTokenData](msg.SourceTokenData[tokenIndex])
	if err != nil {
		return []byte{}, err
	}
	destTokenData := decodedSourceTokenData.ExtraData
	// We don't have better way to determine if the extraData is a payload or sha256(payload)
	// Last parameter of the payload struct is 32-bytes nonce (see Lombard's Bridge._deposit(...) method),
	// so we can assume that payload always exceeds 32 bytes
	if len(destTokenData) != 32 {
		s.lggr.Infow("SourceTokenData.extraData size is not 32. This is deposit payload, not sha256(payload). Attestation is disabled onchain",
			"destTokenData", hexutil.Encode(destTokenData))
		return destTokenData, nil
	}
	payloadHash := [32]byte(destTokenData)

	msgID := hexutil.Encode(msg.MessageID[:])
	payloadHashHex := hexutil.Encode(payloadHash[:])
	s.lggr.Infow("Calling attestation API", "messageBodyHash", payloadHashHex, "messageID", msgID)

	attestationResp, err := s.callAttestationApi(ctx, payloadHash)
	if err != nil {
		return nil, errors.Wrap(err, "failed calling lbtc attestation API")
	}
	if attestationResp.Attestations == nil || len(attestationResp.Attestations) == 0 {
		return nil, errors.New("attestation response is empty")
	}
	if len(attestationResp.Attestations) > 1 {
		s.lggr.Warnw("Multiple attestations received, expected one", "attestations", attestationResp.Attestations)
	}
	var attestation messageAttestationResponse
	for _, attestationCandidate := range attestationResp.Attestations {
		if attestationCandidate.MessageHash == payloadHashHex {
			attestation = attestationCandidate
		}
	}
	if attestation == (messageAttestationResponse{}) {
		return nil, fmt.Errorf("requested attestation %s not found in response", payloadHashHex)
	}
	s.lggr.Infow("Got response from attestation API", "messageID", msgID,
		"attestationStatus", attestation.Status, "attestation", attestation)
	switch attestation.Status {
	case attestationStatusSessionApproved:
		payloadAndProof, err := hexutil.Decode(attestation.Attestation)
		if err != nil {
			return nil, err
		}
		return payloadAndProof, nil
	case attestationStatusPending:
		return nil, tokendata.ErrNotReady
	case attestationStatusSubmitted:
		return nil, tokendata.ErrNotReady
	default:
		s.lggr.Errorw("Unexpected response from attestation API", "attestation", attestation)
		return nil, ErrUnknownResponse
	}
}

func (s *TokenDataReader) callAttestationApi(ctx context.Context, lbtcMessageHash [32]byte) (attestationResponse, error) {
	attestationUrl := fmt.Sprintf("%s/bridge/%s/%s", s.attestationApi.String(), apiVersion, attestationPath)
	request := attestationRequest{PayloadHashes: []string{hexutil.Encode(lbtcMessageHash[:])}}
	encodedRequest, err := json.Marshal(request)
	requestBuffer := bytes.NewBuffer(encodedRequest)
	if err != nil {
		return attestationResponse{}, err
	}
	respRaw, _, _, err := s.httpClient.Post(ctx, attestationUrl, requestBuffer, s.attestationApiTimeout)
	switch {
	case errors.Is(err, tokendata.ErrRateLimit):
		s.setCoolDownPeriod(defaultCoolDownDuration)
		return attestationResponse{}, tokendata.ErrRateLimit
	case err != nil:
		return attestationResponse{}, err
	}
	var attestationResp attestationResponse
	err = json.Unmarshal(respRaw, &attestationResp)
	return attestationResp, err
}

func (s *TokenDataReader) setCoolDownPeriod(d time.Duration) {
	s.coolDownMu.Lock()
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
