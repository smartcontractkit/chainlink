package usdc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TokenDataReader struct {
	lggr               logger.Logger
	sourceChainEvents  ccipdata.Reader
	attestationApi     *url.URL
	messageTransmitter common.Address
	sourceToken        common.Address
	onRampAddress      common.Address

	// Cache of sequence number -> usdc message body
	usdcMessageHashCache      map[uint64][32]byte
	usdcMessageHashCacheMutex sync.Mutex
}

type attestationResponse struct {
	Status      attestationStatus `json:"status"`
	Attestation string            `json:"attestation"`
}

const (
	version                  = "v1"
	attestationPath          = "attestations"
	MESSAGE_SENT_FILTER_NAME = "USDC message sent"
)

type attestationStatus string

const (
	attestationStatusSuccess attestationStatus = "complete"
	attestationStatusPending attestationStatus = "pending_confirmations"
)

var _ tokendata.Reader = &TokenDataReader{}

func NewUSDCTokenDataReader(lggr logger.Logger, sourceChainEvents ccipdata.Reader, usdcTokenAddress, usdcMessageTransmitterAddress, onRampAddress common.Address, usdcAttestationApi *url.URL) *TokenDataReader {
	return &TokenDataReader{
		lggr:                 lggr.With("tokenDataProvider", "usdc"),
		sourceChainEvents:    sourceChainEvents,
		attestationApi:       usdcAttestationApi,
		messageTransmitter:   usdcMessageTransmitterAddress,
		onRampAddress:        onRampAddress,
		sourceToken:          usdcTokenAddress,
		usdcMessageHashCache: make(map[uint64][32]byte),
	}
}

func (s *TokenDataReader) ReadTokenData(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (attestation []byte, err error) {
	response, err := s.getUpdatedAttestation(ctx, msg)
	if err != nil {
		return []byte{}, err
	}

	if response.Status == attestationStatusSuccess {
		attestationBytes, err := hex.DecodeString(response.Attestation)
		if err != nil {
			return nil, fmt.Errorf("decode response attestation: %w", err)
		}
		return attestationBytes, nil
	}
	return []byte{}, tokendata.ErrNotReady
}

func (s *TokenDataReader) getUpdatedAttestation(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) (attestationResponse, error) {
	messageBodyHash, err := s.getUSDCMessageBodyHash(ctx, msg)
	if err != nil {
		return attestationResponse{}, errors.Wrap(err, "failed getting the USDC message body")
	}

	s.lggr.Infow("Calling attestation API", "messageBody", messageBodyHash, "messageID", msg.MessageId)

	response, err := s.callAttestationApi(ctx, messageBodyHash)
	if err != nil {
		return attestationResponse{}, errors.Wrap(err, "failed calling usdc attestation API ")
	}

	return response, nil
}

func (s *TokenDataReader) getUSDCMessageBodyHash(ctx context.Context, msg internal.EVM2EVMOnRampCCIPSendRequestedWithMeta) ([32]byte, error) {
	s.usdcMessageHashCacheMutex.Lock()
	defer s.usdcMessageHashCacheMutex.Unlock()

	if body, ok := s.usdcMessageHashCache[msg.SequenceNumber]; ok {
		return body, nil
	}

	usdcMessageBody, err := s.sourceChainEvents.GetLastUSDCMessagePriorToLogIndexInTx(ctx, int64(msg.LogIndex), msg.TxHash)
	if err != nil {
		return [32]byte{}, err
	}

	s.lggr.Infow("Got USDC message body", "messageBody", usdcMessageBody, "messageID", msg.MessageId)

	msgBodyHash := utils.Keccak256Fixed(usdcMessageBody)

	// Save the attempt in the cache in case the external call fails
	s.usdcMessageHashCache[msg.SequenceNumber] = msgBodyHash
	return msgBodyHash, nil
}

func (s *TokenDataReader) callAttestationApi(ctx context.Context, usdcMessageHash [32]byte) (attestationResponse, error) {
	fullAttestationUrl := fmt.Sprintf("%s/%s/%s/0x%x", s.attestationApi, version, attestationPath, usdcMessageHash)
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
