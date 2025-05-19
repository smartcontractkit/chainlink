package streams

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/datastreams"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v3/reportcodec"
)

type codec struct {
	lggr logger.Logger
}

var _ datastreams.ReportCodec = &codec{}

func (c *codec) Unwrap(wrapped values.Value) ([]datastreams.FeedReport, error) {
	dest, err := datastreams.UnwrapStreamsTriggerEventToFeedReportList(wrapped)
	if err != nil {
		return nil, fmt.Errorf("failed to unwrap: %w", err)
	}
	for i, report := range dest {
		// decoding fields
		id, err2 := datastreams.NewFeedID(report.FeedID)
		if err2 != nil {
			return nil, fmt.Errorf("malformed feed ID: %w", err2)
		}
		v3Codec := reportcodec.NewReportCodec(id.Bytes(), nil)
		decoded, err2 := v3Codec.Decode(report.FullReport)
		if err2 != nil {
			return nil, fmt.Errorf("failed to decode: %w", err2)
		}
		if decoded.FeedId != id.Bytes() {
			return nil, fmt.Errorf("feed ID mismatch: FeedID: %s, FullReport.FeedId: %s", id, hex.EncodeToString(decoded.FeedId[:]))
		}
		dest[i].BenchmarkPrice = decoded.BenchmarkPrice.Bytes()
		dest[i].ObservationTimestamp = int64(decoded.ObservationsTimestamp)
	}
	return dest, nil
}

func (c *codec) Wrap(reports []datastreams.FeedReport) (values.Value, error) {
	return values.Wrap(&datastreams.StreamsTriggerEvent{
		Payload: reports,
	})
}

func (c *codec) Validate(report datastreams.FeedReport, allowedSigners [][]byte, minRequiredSignatures int) error {
	signersMap := make(map[common.Address]struct{})
	for _, signer := range allowedSigners {
		signersMap[common.BytesToAddress(signer)] = struct{}{}
	}
	// signatures (report and context are signed together)
	sigData := append(crypto.Keccak256(report.FullReport), report.ReportContext...)
	fullHash := crypto.Keccak256(sigData)
	validated := map[common.Address]struct{}{}
	for _, sig := range report.Signatures {
		signerPubkey, err2 := crypto.SigToPub(fullHash, sig)
		if err2 != nil {
			return fmt.Errorf("malformed signer: %w", err2)
		}
		signerAddr := crypto.PubkeyToAddress(*signerPubkey)
		if _, ok := signersMap[signerAddr]; !ok {
			c.lggr.Debugw("invalid signer", "signerAddr", signerAddr)
			continue
		}
		validated[signerAddr] = struct{}{}
		if len(validated) >= minRequiredSignatures {
			break // early exit
		}
	}
	if len(validated) < minRequiredSignatures {
		return fmt.Errorf("not enough valid signatures %d, needed %d", len(validated), minRequiredSignatures)
	}
	return nil
}

func NewCodec(lggr logger.Logger) *codec {
	return &codec{lggr: lggr}
}
