package ocr3

import (
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jonboulle/clockwork"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	pbtypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
)

func TestTransmitter(t *testing.T) {
	wid := "consensus-workflow-test-id-1"
	wowner := "foo-owner"
	repId := []byte{0xf0, 0xe0}
	ctx := tests.Context(t)
	lggr := logger.Test(t)
	s := newStore()

	weid := uuid.New().String()

	cp := newCapability(
		s,
		clockwork.NewFakeClock(),
		10*time.Second,
		mockAggregatorFactory,
		func(config *values.Map) (pbtypes.Encoder, error) { return &encoder{}, nil },
		lggr,
		10,
	)
	servicetest.Run(t, cp)

	payload, err := values.NewMap(map[string]any{"observations": []string{"something happened"}})
	require.NoError(t, err)
	config, err := values.NewMap(map[string]any{
		"aggregation_method": "data_feeds",
		"aggregation_config": map[string]any{},
		"encoder":            "",
		"encoder_config":     map[string]any{},
		"report_id":          hex.EncodeToString(repId),
	})
	require.NoError(t, err)
	gotCh, err := cp.Execute(ctx, capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: weid,
			WorkflowID:          wid,
		},
		Config: config,
		Inputs: payload,
	})
	require.NoError(t, err)

	r := mocks.NewCapabilitiesRegistry(t)
	r.On("Get", mock.Anything, ocrCapabilityID).Return(cp, nil)

	info := &pbtypes.ReportInfo{
		Id: &pbtypes.Id{
			WorkflowExecutionId: weid,
			WorkflowId:          wid,
			WorkflowOwner:       wowner,
			ReportId:            hex.EncodeToString(repId),
		},
		ShouldReport: true,
	}
	infob, err := proto.Marshal(info)
	require.NoError(t, err)

	sp := values.Proto(values.NewString("hello"))
	spb, err := proto.Marshal(sp)
	require.NoError(t, err)
	rep := ocr3types.ReportWithInfo[[]byte]{
		Info:   infob,
		Report: spb,
	}

	transmitter := NewContractTransmitter(lggr, r, "fromAccountString")

	var sqNr uint64
	sigs := []types.AttributedOnchainSignature{
		{Signature: []byte("a-signature")},
	}
	err = transmitter.Transmit(ctx, types.ConfigDigest{}, sqNr, rep, sigs)
	require.NoError(t, err)

	resp := <-gotCh
	assert.Nil(t, resp.Err)

	signedReport := pbtypes.SignedReport{}
	require.NoError(t, resp.Value.UnwrapTo(&signedReport))

	assert.Equal(t, spb, signedReport.Report)
	assert.Len(t, signedReport.Signatures, 1)
	assert.Len(t, signedReport.Context, 96)
	assert.Equal(t, repId, signedReport.ID)
}

func TestTransmitter_ShouldReportFalse(t *testing.T) {
	wid := "consensus-workflow-test-id-1"
	wowner := "foo-owner"
	ctx := tests.Context(t)
	lggr := logger.Test(t)
	s := newStore()

	weid := uuid.New().String()

	cp := newCapability(
		s,
		clockwork.NewFakeClock(),
		10*time.Second,
		mockAggregatorFactory,
		func(config *values.Map) (pbtypes.Encoder, error) { return &encoder{}, nil },
		lggr,
		10,
	)
	servicetest.Run(t, cp)

	payload, err := values.NewMap(map[string]any{"observations": []string{"something happened"}})
	require.NoError(t, err)
	config, err := values.NewMap(map[string]any{
		"aggregation_method": "data_feeds",
		"aggregation_config": map[string]any{},
		"encoder":            "",
		"encoder_config":     map[string]any{},
		"report_id":          "aaff",
	})
	require.NoError(t, err)
	gotCh, err := cp.Execute(ctx, capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: weid,
			WorkflowID:          wid,
		},
		Inputs: payload,
		Config: config,
	})
	require.NoError(t, err)

	r := mocks.NewCapabilitiesRegistry(t)
	r.On("Get", mock.Anything, ocrCapabilityID).Return(cp, nil)

	info := &pbtypes.ReportInfo{
		Id: &pbtypes.Id{
			WorkflowExecutionId: weid,
			WorkflowId:          wid,
			WorkflowOwner:       wowner,
		},
		ShouldReport: false,
	}
	infob, err := proto.Marshal(info)
	require.NoError(t, err)

	sp := values.Proto(values.NewString("hello"))
	spb, err := proto.Marshal(sp)
	require.NoError(t, err)
	rep := ocr3types.ReportWithInfo[[]byte]{
		Info:   infob,
		Report: spb,
	}

	transmitter := NewContractTransmitter(lggr, r, "fromAccountString")

	var sqNr uint64
	sigs := []types.AttributedOnchainSignature{
		{Signature: []byte("a-signature")},
	}
	err = transmitter.Transmit(ctx, types.ConfigDigest{}, sqNr, rep, sigs)
	require.NoError(t, err)

	resp := <-gotCh
	assert.NotNil(t, resp.Err)
	assert.True(t, errors.Is(resp.Err, capabilities.ErrStopExecution))
}
