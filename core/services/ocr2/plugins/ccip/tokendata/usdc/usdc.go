package usdc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	apiVersion               = "v1"
	attestationPath          = "attestations"
	MESSAGE_SENT_FILTER_NAME = "USDC message sent"
)

type attestationStatus string

const (
	attestationStatusSuccess attestationStatus = "complete"
	attestationStatusPending attestationStatus = "pending_confirmations"
)

// usdcPayload has to match the onchain event emitted by the USDC message transmitter
type usdcPayload []byte

func (d usdcPayload) AbiString() string {
	return `[{"type": "bytes"}]`
}

func (d usdcPayload) Validate() error {
	if len(d) == 0 {
		return errors.New("must be non-empty")
	}
	return nil
}

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
	lggr               logger.Logger
	sourceChainEvents  ccipdata.Reader
	attestationApi     *url.URL
	messageTransmitter common.Address
	sourceToken        common.Address
	onRampAddress      common.Address

	// Cache of sequence number -> usdc message body
	usdcMessageHashCache      map[uint64][]byte
	usdcMessageHashCacheMutex sync.Mutex
}

type attestationResponse struct {
	Status      attestationStatus `json:"status"`
	Attestation string            `json:"attestation"`
}

var _ tokendata.Reader = &TokenDataReader{}

func NewUSDCTokenDataReader(lggr logger.Logger, sourceChainEvents ccipdata.Reader, usdcTokenAddress, usdcMessageTransmitterAddress, onRampAddress common.Address, usdcAttestationApi *url.URL) *TokenDataReader {
	return &TokenDataReader{
		lggr:                 lggr.With("tokenDataProvider", "usdc"),
		sourceChainEvents:    sourceChainEvents,
		attestationApi:       usdcAttestationApi,
		messageTransmitter:   usdcMessageTransmitterAddress,
		onRampAddress:        onRampAddress,
		sourceToken:          usdcTokenAddress,
		usdcMessageHashCache: make(map[uint64][]byte),
	}
}

func (s *TokenDataReader) ReadTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (messageAndAttestation []byte, err error) {
	messageBody, err := s.getUSDCMessageBody(ctx, msg)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed getting the USDC message body")
	}

	s.lggr.Infow("Calling attestation API", "messageBodyHash", hexutil.Encode(messageBody[:]), "messageID", hexutil.Encode(msg.MessageId[:]))

	// The attestation API expects the hash of the message body
	attestationResp, err := s.callAttestationApi(ctx, utils.Keccak256Fixed(messageBody))
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed calling usdc attestation API ")
	}

	if attestationResp.Status != attestationStatusSuccess {
		return []byte{}, tokendata.ErrNotReady
	}

	// The USDC pool needs a combination of the message body and the attestation
	messageAndAttestation, err = encodeMessageAndAttestation(messageBody, attestationResp.Attestation)
	if err != nil {
		return nil, fmt.Errorf("failed to encode messageAndAttestation : %w", err)
	}

	return messageAndAttestation, nil
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

func (s *TokenDataReader) getUSDCMessageBody(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([]byte, error) {
	s.usdcMessageHashCacheMutex.Lock()
	defer s.usdcMessageHashCacheMutex.Unlock()

	if body, ok := s.usdcMessageHashCache[msg.SequenceNumber]; ok {
		return body, nil
	}

	usdcMessageBody, err := s.sourceChainEvents.GetLastUSDCMessagePriorToLogIndexInTx(ctx, int64(msg.LogIndex), msg.TxHash)
	if err != nil {
		return []byte{}, err
	}

	parsedMsgBody, err := decodeUSDCMessageSent(usdcMessageBody)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed parsing solidity struct")
	}

	s.lggr.Infow("Got USDC message body", "messageBody", hexutil.Encode(parsedMsgBody), "messageID", hexutil.Encode(msg.MessageId[:]))

	// Save the attempt in the cache in case the external call fails
	s.usdcMessageHashCache[msg.SequenceNumber] = parsedMsgBody
	return parsedMsgBody, nil
}

func decodeUSDCMessageSent(logData []byte) ([]byte, error) {
	decodeAbiStruct, err := abihelpers.DecodeAbiStruct[usdcPayload](logData)
	if err != nil {
		return nil, err
	}
	return decodeAbiStruct, nil
}

func (s *TokenDataReader) callAttestationApi(ctx context.Context, usdcMessageHash [32]byte) (attestationResponse, error) {
	fullAttestationUrl := fmt.Sprintf("%s/%s/%s/0x%x", s.attestationApi, apiVersion, attestationPath, usdcMessageHash)
	req, err := http.NewRequestWithContext(ctx, "GET", fullAttestationUrl, nil)
	if err != nil {
		return attestationResponse{}, err
	}
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return attestationResponse{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return attestationResponse{}, err
	}

	var response attestationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return attestationResponse{}, err
	}

	if response.Status == "" {
		return attestationResponse{}, fmt.Errorf("invalid attestation response: %s", string(body))
	}

	return response, nil
}

func (s *TokenDataReader) GetSourceLogPollerFilters() []logpoller.Filter {
	return []logpoller.Filter{
		{
			Name:      logpoller.FilterName(MESSAGE_SENT_FILTER_NAME, s.messageTransmitter.Hex()),
			EventSigs: []common.Hash{abihelpers.EventSignatures.USDCMessageSent},
			Addresses: []common.Address{s.messageTransmitter},
		},
	}
}
