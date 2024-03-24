package pluginprovider

import (
	"bytes"
	"context"
	"fmt"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"

	testtypes "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/test/types"
)

type contractTransmitterTestConfig struct {
	ReportContext libocr.ReportContext
	Report        libocr.Report
	Sigs          []libocr.AttributedOnchainSignature

	ConfigDigest libocr.ConfigDigest
	Account      libocr.Account
	Epoch        uint32
}

var _ testtypes.ContractTransmitterEvaluator = staticContractTransmitter{}

type staticContractTransmitter struct {
	contractTransmitterTestConfig
}

func (s staticContractTransmitter) Transmit(ctx context.Context, rc libocr.ReportContext, r libocr.Report, ss []libocr.AttributedOnchainSignature) error {
	if !assert.ObjectsAreEqual(s.ReportContext, rc) {
		return fmt.Errorf("expected report context %v but got %v", s.ReportContext, rc)
	}
	if !bytes.Equal(s.Report, r) {
		return fmt.Errorf("expected report %x but got %x", s.Report, r)
	}
	if !assert.ObjectsAreEqual(s.Sigs, ss) {
		return fmt.Errorf("expected signatures %v but got %v", s.Sigs, ss)
	}
	return nil
}

func (s staticContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (libocr.ConfigDigest, uint32, error) {
	return s.ConfigDigest, s.Epoch, nil
}

func (s staticContractTransmitter) FromAccount(ctx context.Context) (libocr.Account, error) {
	return s.Account, nil
}

func (s staticContractTransmitter) Evaluate(ctx context.Context, ct libocr.ContractTransmitter) error {
	gotAccount, err := ct.FromAccount(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FromAccount: %w", err)
	}
	if gotAccount != s.Account {
		return fmt.Errorf("expectd FromAccount %s but got %s", s.Account, gotAccount)
	}
	gotConfigDigest, gotEpoch, err := ct.LatestConfigDigestAndEpoch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get LatestConfigDigestAndEpoch: %w", err)
	}
	if gotConfigDigest != s.ConfigDigest {
		return fmt.Errorf("expected ConfigDigest %s but got %s", s.ConfigDigest, gotConfigDigest)
	}
	if gotEpoch != epoch {
		return fmt.Errorf("expected Epoch %d but got %d", epoch, gotEpoch)
	}
	err = ct.Transmit(ctx, s.ReportContext, s.Report, sigs)
	if err != nil {
		return fmt.Errorf("failed to Transmit")
	}
	return nil
}
