package executable

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
	"github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
)

func TestWriteReportExcludeSignaturesHasher_Hash(t *testing.T) {
	req1a := getRequest(t, []byte("testdata"), [][]byte{[]byte("sig1"), []byte("sig2")})
	req1b := getRequest(t, []byte("testdata"), [][]byte{[]byte("sig3"), []byte("sig4")})
	req2 := getRequest(t, []byte("otherdata"), [][]byte{[]byte("sig1"), []byte("sig2")})

	hasher := NewWriteReportExcludeSignaturesHasher()
	hash1a, err := hasher.Hash(req1a)
	require.NoError(t, err)
	hash1b, err := hasher.Hash(req1b)
	require.NoError(t, err)
	hash2, err := hasher.Hash(req2)
	require.NoError(t, err)

	require.Equal(t, hash1a, hash1b)   // same data, different signatures
	require.NotEqual(t, hash1a, hash2) // different data, same signatures
}

func TestWriteReportExcludeSignaturesHasher_Hash_NilPayload(t *testing.T) {
	nilReq := capabilities.CapabilityRequest{Payload: nil}
	nilReqBytes, err := pb.MarshalCapabilityRequest(nilReq)
	require.NoError(t, err)

	msgBody := &types.MessageBody{Payload: nilReqBytes}

	hasher := NewWriteReportExcludeSignaturesHasher()
	_, err = hasher.Hash(msgBody)
	require.Error(t, err)
	require.Contains(t, err.Error(), "capability request payload is nil")
}

func TestWriteReportExcludeSignaturesHasher_Hash_NilReport(t *testing.T) {
	nilReq := &evmcappb.WriteReportRequest{Report: nil}
	nilPb, err := anypb.New(nilReq)
	require.NoError(t, err)

	capReq := capabilities.CapabilityRequest{Payload: nilPb}
	capReqBytes, err := pb.MarshalCapabilityRequest(capReq)
	require.NoError(t, err)

	msgBody := &types.MessageBody{Payload: capReqBytes}

	hasher := NewWriteReportExcludeSignaturesHasher()
	_, err = hasher.Hash(msgBody)
	require.Error(t, err)
	require.Contains(t, err.Error(), "WriteReportRequest.Report is nil")
}

func TestWriteReportExcludeSignaturesHasher_Hash_InvalidPayload(t *testing.T) {
	// Test with completely invalid payload that cannot be unmarshaled
	msgBody := &types.MessageBody{
		Payload: []byte("invalid protobuf data"),
	}

	hasher := NewWriteReportExcludeSignaturesHasher()
	_, err := hasher.Hash(msgBody)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to unmarshal capability request")
}

func TestSimpleHasher_ExcludesSpendLimits(t *testing.T) {
	// Create two requests with identical payloads but different SpendLimits
	req1 := getRequestWithSpendLimits(t, []byte("testdata"), []capabilities.SpendLimit{
		{SpendType: "gas", Limit: "1000"},
		{SpendType: "compute", Limit: "500"},
	})
	req2 := getRequestWithSpendLimits(t, []byte("testdata"), []capabilities.SpendLimit{
		{SpendType: "gas", Limit: "2000"},
		{SpendType: "compute", Limit: "1000"},
	})
	req3 := getRequestWithSpendLimits(t, []byte("otherdata"), []capabilities.SpendLimit{
		{SpendType: "gas", Limit: "1000"},
		{SpendType: "compute", Limit: "500"},
	})

	hasher := NewSimpleHasher()
	hash1, err := hasher.Hash(req1)
	require.NoError(t, err)
	hash2, err := hasher.Hash(req2)
	require.NoError(t, err)
	hash3, err := hasher.Hash(req3)
	require.NoError(t, err)

	require.Equal(t, hash1, hash2)    // same data, different SpendLimits should produce same hash
	require.NotEqual(t, hash1, hash3) // different data should produce different hash
}

func TestWriteReportExcludeSignaturesHasher_ExcludesSpendLimits(t *testing.T) {
	// Create two requests with identical payloads but different SpendLimits
	req1 := getWriteReportRequestWithSpendLimits(t, []byte("testdata"), [][]byte{[]byte("sig1"), []byte("sig2")}, []capabilities.SpendLimit{
		{SpendType: "gas", Limit: "1000"},
		{SpendType: "compute", Limit: "500"},
	})
	req2 := getWriteReportRequestWithSpendLimits(t, []byte("testdata"), [][]byte{[]byte("sig3"), []byte("sig4")}, []capabilities.SpendLimit{
		{SpendType: "gas", Limit: "2000"},
		{SpendType: "compute", Limit: "1000"},
	})
	req3 := getWriteReportRequestWithSpendLimits(t, []byte("otherdata"), [][]byte{[]byte("sig1"), []byte("sig2")}, []capabilities.SpendLimit{
		{SpendType: "gas", Limit: "1000"},
		{SpendType: "compute", Limit: "500"},
	})

	hasher := NewWriteReportExcludeSignaturesHasher()
	hash1, err := hasher.Hash(req1)
	require.NoError(t, err)
	hash2, err := hasher.Hash(req2)
	require.NoError(t, err)
	hash3, err := hasher.Hash(req3)
	require.NoError(t, err)

	require.Equal(t, hash1, hash2)    // same data, different SpendLimits and signatures should produce same hash
	require.NotEqual(t, hash1, hash3) // different data should produce different hash
}

func getRequest(t *testing.T, data []byte, sigs [][]byte) *types.MessageBody {
	attrSigs := []*sdk.AttributedSignature{}
	for i, sig := range sigs {
		attrSigs = append(attrSigs, &sdk.AttributedSignature{
			Signature: sig,
			SignerId:  uint32(i), //nolint:gosec // G115
		})
	}
	report := &sdk.ReportResponse{
		RawReport: data,
		Sigs:      attrSigs,
	}
	wrReq := &evmcappb.WriteReportRequest{
		Report: report,
	}
	wrAny, err := anypb.New(wrReq)
	require.NoError(t, err)
	capReq := capabilities.CapabilityRequest{
		Payload: wrAny,
	}
	capReqBytes, err := pb.MarshalCapabilityRequest(capReq)
	require.NoError(t, err)
	return &types.MessageBody{
		Payload: capReqBytes,
	}
}

func getRequestWithSpendLimits(t *testing.T, data []byte, spendLimits []capabilities.SpendLimit) *types.MessageBody {
	report := &sdk.ReportResponse{
		RawReport: data,
		Sigs:      []*sdk.AttributedSignature{},
	}
	wrReq := &evmcappb.WriteReportRequest{
		Report: report,
	}
	wrAny, err := anypb.New(wrReq)
	require.NoError(t, err)
	capReq := capabilities.CapabilityRequest{
		Payload: wrAny,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
			SpendLimits:         spendLimits,
		},
	}
	capReqBytes, err := pb.MarshalCapabilityRequest(capReq)
	require.NoError(t, err)
	return &types.MessageBody{
		Payload: capReqBytes,
	}
}

func getWriteReportRequestWithSpendLimits(t *testing.T, data []byte, sigs [][]byte, spendLimits []capabilities.SpendLimit) *types.MessageBody {
	attrSigs := []*sdk.AttributedSignature{}
	for i, sig := range sigs {
		attrSigs = append(attrSigs, &sdk.AttributedSignature{
			Signature: sig,
			SignerId:  uint32(i), //nolint:gosec // G115
		})
	}
	report := &sdk.ReportResponse{
		RawReport: data,
		Sigs:      attrSigs,
	}
	wrReq := &evmcappb.WriteReportRequest{
		Report: report,
	}
	wrAny, err := anypb.New(wrReq)
	require.NoError(t, err)
	capReq := capabilities.CapabilityRequest{
		Payload: wrAny,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          "test-workflow",
			WorkflowExecutionID: "test-execution",
			SpendLimits:         spendLimits,
		},
	}
	capReqBytes, err := pb.MarshalCapabilityRequest(capReq)
	require.NoError(t, err)
	return &types.MessageBody{
		Payload: capReqBytes,
	}
}
