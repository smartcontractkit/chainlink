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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
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
	defaultCoolDownDuration = 60 * time.Second

	// maxCoolDownDuration defines the maximum duration we can wait till firing the next request
	maxCoolDownDuration = 10 * time.Minute
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

func NewUSDCTokenDataReader(lggr logger.Logger, usdcReader ccipdata.USDCReader, usdcAttestationApi *url.URL, usdcAttestationApiTimeoutSeconds int) *TokenDataReader {
	timeout := time.Duration(usdcAttestationApiTimeoutSeconds) * time.Second
	if usdcAttestationApiTimeoutSeconds == 0 {
		timeout = defaultAttestationTimeout
	}
	return &TokenDataReader{
		lggr:                  lggr,
		usdcReader:            usdcReader,
		httpClient:            http.NewObservedIHttpClient(&http.HttpClient{}),
		attestationApi:        usdcAttestationApi,
		attestationApiTimeout: timeout,
		coolDownMu:            &sync.RWMutex{},
	}
}

func NewUSDCTokenDataReaderWithHttpClient(origin TokenDataReader, httpClient http.IHttpClient) *TokenDataReader {
	return &TokenDataReader{
		lggr:                  origin.lggr,
		usdcReader:            origin.usdcReader,
		httpClient:            httpClient,
		attestationApi:        origin.attestationApi,
		attestationApiTimeout: origin.attestationApiTimeout,
		coolDownMu:            origin.coolDownMu,
	}
}

func (s *TokenDataReader) ReadTokenData(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta, tokenIndex int) (messageAndAttestation []byte, err error) {
	if tokenIndex < 0 || tokenIndex >= len(msg.TokenAmounts) {
		return nil, fmt.Errorf("token index out of bounds")
	}

	if s.inCoolDownPeriod() {
		// rate limiting cool-down period, we prevent new requests from being sent
		return nil, tokendata.ErrRequestsBlocked
	}

	messageBody, err := s.getUSDCMessageBody(ctx, msg)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed getting the USDC message body")
	}

	s.lggr.Infow("Calling attestation API", "messageBodyHash", hexutil.Encode(messageBody[:]), "messageID", hexutil.Encode(msg.MessageID[:]))

	// The attestation API expects the hash of the message body
	attestationResp, err := s.callAttestationApi(ctx, utils.Keccak256Fixed(messageBody))
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed calling usdc attestation API ")
	}
	switch attestationResp.Status {
	case attestationStatusSuccess:
		// The USDC pool needs a combination of the message body and the attestation
		messageAndAttestation, err = encodeMessageAndAttestation(messageBody, attestationResp.Attestation)
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

func (s *TokenDataReader) getUSDCMessageBody(ctx context.Context, msg cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([]byte, error) {
	parsedMsgBody, err := s.usdcReader.GetLastUSDCMessagePriorToLogIndexInTx(ctx, int64(msg.LogIndex), msg.TxHash)
	if err != nil {
		return []byte{}, err
	}
	s.lggr.Infow("Got USDC message body", "messageBody", hexutil.Encode(parsedMsgBody), "messageID", hexutil.Encode(msg.MessageID[:]))
	return parsedMsgBody, nil
}

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
