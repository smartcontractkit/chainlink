package loop_test

import (
	"bytes"
	"context"
	"fmt"

	libocr "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
)

type staticReportingPlugin struct{}

func (s staticReportingPlugin) Query(ctx context.Context, timestamp libocr.ReportTimestamp) (libocr.Query, error) {
	if timestamp != reportContext.ReportTimestamp {
		return nil, fmt.Errorf("expected %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	return query, nil
}

func (s staticReportingPlugin) Observation(ctx context.Context, timestamp libocr.ReportTimestamp, q libocr.Query) (libocr.Observation, error) {
	if timestamp != reportContext.ReportTimestamp {
		return nil, fmt.Errorf("expected %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(q, query) {
		return nil, fmt.Errorf("expected %x but got %x", query, q)
	}
	return observation, nil
}

func (s staticReportingPlugin) Report(ctx context.Context, timestamp libocr.ReportTimestamp, q libocr.Query, observations []libocr.AttributedObservation) (bool, libocr.Report, error) {
	if timestamp != reportContext.ReportTimestamp {
		return false, nil, fmt.Errorf("expected %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(q, query) {
		return false, nil, fmt.Errorf("expected %x but got %x", query, q)
	}
	if !assert.ObjectsAreEqual(obs, observations) {
		return false, nil, fmt.Errorf("expected %v but got %v", obs, observations)
	}
	return shouldReport, report, nil
}

func (s staticReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp libocr.ReportTimestamp, r libocr.Report) (bool, error) {
	if timestamp != reportContext.ReportTimestamp {
		return false, fmt.Errorf("expected %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(r, report) {
		return false, fmt.Errorf("expected %x but got %x", report, r)
	}
	return shouldAccept, nil
}

func (s staticReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp libocr.ReportTimestamp, r libocr.Report) (bool, error) {
	if timestamp != reportContext.ReportTimestamp {
		return false, fmt.Errorf("expected %v but got %v", reportContext.ReportTimestamp, timestamp)
	}
	if !bytes.Equal(r, report) {
		return false, fmt.Errorf("expected %x but got %x", report, r)
	}
	return shouldTransmit, nil
}

func (s staticReportingPlugin) Close() error { return nil }
