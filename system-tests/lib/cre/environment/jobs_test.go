package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/stretchr/testify/require"
)

func TestRewriteJDForDirectAccess_NilOutputNoop(t *testing.T) {
	var output *jd.Output
	err := rewriteJDForDirectAccess(output, "10.20.30.40")
	require.NoError(t, err, "expected nil output rewrite to be a no-op")
}

func TestRewriteJDForDirectAccessRewritesExternalEndpoints(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1:14231",
		ExternalWSRPCUrl: "127.0.0.1:9080",
		InternalWSRPCUrl: "job-distributor:8080",
	}

	err := rewriteJDForDirectAccess(output, "10.20.30.40")
	require.NoError(t, err, "rewriteJDForDirectAccess should succeed")
	require.Equal(t, "10.20.30.40:14231", output.ExternalGRPCUrl, "external grpc url should be rewritten")
	require.Equal(t, "10.20.30.40:9080", output.ExternalWSRPCUrl, "external wsrpc url should be rewritten")
	require.Equal(t, "job-distributor:8080", output.InternalWSRPCUrl, "internal wsrpc url should remain unchanged")
}

func TestRewriteJDForDirectAccess_MixedFallsBackToInternalWSRPCSource(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1:14231",
		ExternalWSRPCUrl: "",
		InternalWSRPCUrl: "job-distributor:8080",
	}

	err := rewriteJDForDirectAccess(output, "10.20.30.40")
	require.NoError(t, err, "rewriteJDForDirectAccess should succeed")
	require.Equal(t, "10.20.30.40:8080", output.ExternalWSRPCUrl, "external wsrpc url should be derived from internal source")
	require.Equal(t, "job-distributor:8080", output.InternalWSRPCUrl, "internal wsrpc url should remain unchanged")
}

func TestRewriteJDForDirectAccess_InvalidAddressFails(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1",
		ExternalWSRPCUrl: "127.0.0.1:9080",
	}

	err := rewriteJDForDirectAccess(output, "10.20.30.40")
	require.Error(t, err, "expected invalid host:port to fail rewrite")
	require.Contains(t, err.Error(), "failed to parse host:port", "expected parse failure context")
}

func TestRewriteAddressHost_UnsupportedURLWithoutPortFails(t *testing.T) {
	_, err := rewriteAddressHost("http://job-distributor", "10.20.30.40")
	require.Error(t, err, "expected address without port to fail")
	require.Contains(t, err.Error(), "must include a port", "expected missing-port context")
}

func TestRewriteAddressHost_EmptyInputNoop(t *testing.T) {
	rewritten, err := rewriteAddressHost("   ", "10.20.30.40")
	require.NoError(t, err, "expected empty input to be a no-op")
	require.Equal(t, "", rewritten, "expected empty output for empty input")
}
