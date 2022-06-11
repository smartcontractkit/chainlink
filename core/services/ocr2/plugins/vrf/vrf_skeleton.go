package vrf

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

type VRFFactory struct {
	logger logger.Logger
}

func (fac VRFFactory) NewReportingPlugin(configuration types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	return &VRF{
			fac.logger,
		}, types.ReportingPluginInfo{
			Name:          "VRF",
			UniqueReports: false,
			Limits: types.ReportingPluginLimits{
				MaxQueryLength:       0,      // arbitrary
				MaxObservationLength: 100000, // arbitrary
				MaxReportLength:      100000, // arbitrary
			},
		}, nil
}

type VRF struct {
	logger logger.Logger
}

func (vrf *VRF) Query(ctx context.Context, repts types.ReportTimestamp) (types.Query, error) {
	vrf.logger.Info("OCR2 NODE IS QUERYING")
	time.Sleep(time.Second)
	return nil, nil
}

func (vrf *VRF) Observation(ctx context.Context, repts types.ReportTimestamp, query types.Query) (types.Observation, error) {
	vrf.logger.Info("OCR2 NODE IS OBSERVING")
	time.Sleep(time.Second)
	return []byte{1, 2, 3}, nil
}

func (vrf *VRF) Report(ctx context.Context, repts types.ReportTimestamp, query types.Query, aos []types.AttributedObservation) (bool, types.Report, error) {
	vrf.logger.Info("OCR2 NODE IS REPORTING OBSERVATIONS: ", aos)
	time.Sleep(time.Second)

	report := []byte{}
	for _, observation := range aos {
		report = append(report, observation.Observation...)
	}

	vrf.logger.Info("REPORT IS CONSTRUCTED: ", report)
	return true, report, nil
}

func (vrf *VRF) shouldReport(ctx context.Context, repts types.ReportTimestamp, paos []ParsedAttributedObservation) (bool, error) {
	vrf.logger.Info("OCR2 NODE IS DECIDING IF IT WANTS TO REPORT OBSERVATIONS: ", paos)
	time.Sleep(time.Second)
	return true, nil
}

func (vrf *VRF) ShouldAcceptFinalizedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	vrf.logger.Info("OCR2 NODE IS DECIDING IF IT WANTS TO ACCEPT THE FINALIZED REPORT: ", report)
	time.Sleep(time.Second)
	return true, nil
}

func (vrf *VRF) ShouldTransmitAcceptedReport(ctx context.Context, repts types.ReportTimestamp, report types.Report) (bool, error) {
	vrf.logger.Info("OCR2 NODE IS NOT TRANSMITTING THE ACCEPTED REPORT: ", report)
	time.Sleep(time.Second)
	return false, nil
}

func (vrf *VRF) Close() error {
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
