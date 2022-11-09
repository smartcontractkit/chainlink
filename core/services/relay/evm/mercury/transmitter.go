package mercury

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ ocrtypes.ContractTransmitter = &MercuryTransmitter{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type MercuryTransmitter struct {
	lggr       logger.Logger
	httpClient HTTPClient

	fromAccount common.Address

	reportURL string
	username  string
	password  string
}

var payloadTypes = getPayloadTypes()

func getPayloadTypes() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "reportContext", Type: mustNewType("bytes32[3]")},
		{Name: "report", Type: mustNewType("bytes")},
		{Name: "rawRs", Type: mustNewType("bytes32[]")},
		{Name: "rawSs", Type: mustNewType("bytes32[]")},
		{Name: "rawVs", Type: mustNewType("bytes32")},
	})
}

func NewTransmitter(lggr logger.Logger, httpClient HTTPClient, fromAccount common.Address, reportURL, username, password string) *MercuryTransmitter {
	return &MercuryTransmitter{lggr.Named("Mercury"), httpClient, fromAccount, reportURL, username, password}
}

type MercuryReport struct {
	Payload     hexutil.Bytes
	FromAccount common.Address
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (mt *MercuryTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	for i, as := range signatures {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := evmutil.RawReportContext(reportCtx)

	payload, err := payloadTypes.Pack(rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	mr := MercuryReport{
		Payload:     payload,
		FromAccount: mt.fromAccount,
	}
	mt.lggr.Infow("Transmitting report", "mercuryReport", mr, "report", report, "reportCtx", reportCtx, "signatures", signatures)

	b, err := json.Marshal(mr)
	if err != nil {
		return errors.Wrap(err, "failed to marshal mercury report JSON")
	}

	req, err := http.NewRequest(http.MethodPost, mt.reportURL, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "failed to instantiate mercury server http request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(mt.username, mt.password)
	req = req.WithContext(ctx)

	res, err := mt.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to POST to mercury server")
	}
	defer res.Body.Close()

	// It's only used for logging, keep it short
	safeLimitResp := http.MaxBytesReader(nil, res.Body, 1024)
	respBody, err := io.ReadAll(safeLimitResp)
	if err != nil {
		mt.lggr.Errorw("Failed to read response body", "err", err)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		mt.lggr.Infow("Transmit report success", "responseStatus", res.Status, "reponseBody", string(respBody), "reportCtx", reportCtx)
	} else {
		mt.lggr.Errorw("Transmit report failed", "responseStatus", res.Status, "reponseBody", string(respBody), "reportCtx", reportCtx)

	}

	return nil
}

func (mt *MercuryTransmitter) FromAccount() ocrtypes.Account {
	return ocrtypes.Account(mt.fromAccount.Hex())
}

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
// It is plugin independent, in particular avoids use of the plugin specific generated evm wrappers
// by using the evm client Call directly for functions/events that are part of OCR2Abstract.
func (mt *MercuryTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (cd ocrtypes.ConfigDigest, epoch uint32, err error) {
	// ConfigDigest and epoch are not stored on the contract in mercury mode
	// TODO: Do we need to support retrieving it from the server? Does it matter?
	// https://app.shortcut.com/chainlinklabs/story/57500/return-the-actual-latest-transmission-details
	err = errors.New("Retrieving config digest/epoch is not supported in Mercury mode")
	return
}
