package dkg

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type DKGFactory struct {
	logger logger.Logger
}

func (fac DKGFactory) NewReportingPlugin(configuration types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	return &DKG{
			fac.logger,
		}, types.ReportingPluginInfo{
			Name:          "DKG",
			UniqueReports: false,
			Limits: types.ReportingPluginLimits{
				MaxQueryLength:       0,      // arbitrary
				MaxObservationLength: 100000, // arbitrary
				MaxReportLength:      100000, // arbitrary
			},
		}, nil
}

type DKG struct {
	logger logger.Logger
}

func (dkg *DKG) Query(ctx context.Context, repts types.ReportTimestamp) (types.Query, error) {
	dkg.logger.Info("OCR2 NODE IS QUERYING")
	time.Sleep(time.Second)
	return nil, nil
}

func (dkg *DKG) Observation(ctx context.Context, repts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	dkg.logger.Info("OCR2 NODE IS OBSERVING")
	time.Sleep(time.Second)
	return []byte{1, 2, 3, 4}, nil
}

func (dkg *DKG) Report(ctx context.Context, repts types.ReportTimestamp, query types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	dkg.logger.Info("OCR2 NODE IS REPORTING OBSERVATIONS: ", aos)
	time.Sleep(time.Second)

	report := []byte{}
	for _, observation := range aos {
		report = append(report, observation.Observation...)
	}

	dkg.logger.Info("REPORT IS CONSTRUCTED: ", report)
	return true, report, nil
}

func (dkg *DKG) shouldReport(ctx context.Context, repts types.ReportTimestamp, paos []ParsedAttributedObservation) (bool, error) {
	dkg.logger.Info("OCR2 NODE IS DECIDING IF IT WANTS TO REPORT OBSERVATIONS: ", paos)
	time.Sleep(time.Second)
	return true, nil
}

func (dkg *DKG) ShouldAcceptFinalizedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	dkg.logger.Info("OCR2 NODE IS DECIDING IF IT WANTS TO ACCEPT THE FINALIZED REPORT: ", report)
	time.Sleep(time.Second)
	return true, nil
}

func (dkg *DKG) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	dkg.logger.Info("OCR2 NODE IS NOT TRANSMITTING THE ACCEPTED REPORT: ", report)
	time.Sleep(time.Second)
	return false, nil
}

func (dkg *DKG) Close() error {
	fmt.Println("WE ARE CLOSING")
	time.Sleep(time.Second)
	return nil
}

type ParsedAttributedObservation struct {
	Timestamp       uint32
	Value           *big.Int
	JuelsPerFeeCoin *big.Int
	Observer        commontypes.OracleID
}
