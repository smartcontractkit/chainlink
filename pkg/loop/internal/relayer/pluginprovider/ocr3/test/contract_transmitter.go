package test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

var (
	configDigest = libocr.ConfigDigest([32]byte{1: 7, 13: 11, 31: 23})

	sigs = []libocr.AttributedOnchainSignature{{Signature: []byte{9: 8, 7: 6}, Signer: commontypes.OracleID(54)}}

	// ContractTransmitter is a static implementation of the ContractTransmitterTester interface for testing
	ContractTransmitter = staticOCR3ContractTransmitter{
		ocr3ContractTransmitterTestConfig: ocr3ContractTransmitterTestConfig{
			ConfigDigest: configDigest,
			SeqNr:        3,
			Report:       libocr.Report{41: 131},
			Info:         []byte("some-info"),
			Sigs:         sigs,
			Account:      libocr.Account("some-account"),
		},
	}
)

type ocr3ContractTransmitterTestConfig struct {
	ConfigDigest libocr.ConfigDigest
	SeqNr        uint64
	Report       libocr.Report
	Info         []byte
	Sigs         []libocr.AttributedOnchainSignature
	Account      libocr.Account
}

var _ testtypes.OCR3ContractTransmitterEvaluator = staticOCR3ContractTransmitter{}

type staticOCR3ContractTransmitter struct {
	ocr3ContractTransmitterTestConfig
}

func (s staticOCR3ContractTransmitter) Transmit(ctx context.Context, configDigest libocr.ConfigDigest, seqNr uint64, r ocr3types.ReportWithInfo[[]byte], ss []libocr.AttributedOnchainSignature) error {
	cd := [32]byte(configDigest)
	haveCd := [32]byte(s.ConfigDigest)
	if !bytes.Equal(cd[:], haveCd[:]) {
		return fmt.Errorf("expected configDigest %x but got %x", haveCd, cd)
	}
	if seqNr != s.SeqNr {
		return fmt.Errorf("expected seqNr %d but got %d", s.SeqNr, seqNr)
	}

	if !bytes.Equal(r.Report, s.Report) {
		return fmt.Errorf("expected report %x but got %x", s.Report, r.Report)
	}

	if !bytes.Equal(r.Info, s.Info) {
		return fmt.Errorf("expected info %x but got %x", s.Info, r.Report)
	}

	if !assert.ObjectsAreEqual(s.Sigs, ss) {
		return fmt.Errorf("expected signatures %v but got %v", s.Sigs, ss)
	}

	return nil
}

func (s staticOCR3ContractTransmitter) FromAccount() (libocr.Account, error) {
	return s.Account, nil
}

func (s staticOCR3ContractTransmitter) Evaluate(ctx context.Context, ct ocr3types.ContractTransmitter[[]byte]) error {
	gotAccount, err := ct.FromAccount()
	if err != nil {
		return fmt.Errorf("failed to get FromAccount: %w", err)
	}
	if gotAccount != s.Account {
		return fmt.Errorf("expectd FromAccount %s but got %s", s.Account, gotAccount)
	}
	err = ct.Transmit(ctx, s.ConfigDigest, s.SeqNr, ocr3types.ReportWithInfo[[]byte]{Report: s.Report, Info: s.Info}, s.Sigs)
	if err != nil {
		return fmt.Errorf("failed to Transmit")
	}
	return nil
}
