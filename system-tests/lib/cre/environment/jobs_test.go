package environment

import (
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/tunnel"
)

func TestDescribeJDEndpointsUsesExternalWSRPC(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1:14231",
		ExternalWSRPCUrl: "127.0.0.1:8080",
		InternalWSRPCUrl: "job-distributor:8080",
	}

	refs, err := describeJDEndpoints(output)
	if err != nil {
		t.Fatalf("describeJDEndpoints returned error: %v", err)
	}
	if len(refs) != 2 {
		t.Fatalf("expected 2 endpoint refs, got %d", len(refs))
	}

	var wsrpcRef *tunnel.EndpointRef
	for i := range refs {
		if refs[i].EndpointName == "wsrpc" {
			wsrpcRef = &refs[i]
			break
		}
	}
	if wsrpcRef == nil {
		t.Fatal("missing wsrpc endpoint ref")
	}
	if wsrpcRef.Host != "127.0.0.1" || wsrpcRef.Port != 8080 {
		t.Fatalf("expected wsrpc endpoint to use external address 127.0.0.1:8080, got %s:%d", wsrpcRef.Host, wsrpcRef.Port)
	}
}

func TestRewriteJDWithBindingsRewritesNodeFacingWSRPC(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1:14231",
		ExternalWSRPCUrl: "127.0.0.1:8080",
		InternalWSRPCUrl: "job-distributor:8080",
	}
	bindings := []tunnel.TunnelBinding{
		{
			EndpointRef: tunnel.EndpointRef{EndpointName: "grpc"},
			LocalPort:   61001,
		},
		{
			EndpointRef: tunnel.EndpointRef{EndpointName: "wsrpc"},
			LocalPort:   61002,
		},
	}

	if err := rewriteJDWithBindings(output, bindings, true); err != nil {
		t.Fatalf("rewriteJDWithBindings returned error: %v", err)
	}

	if output.ExternalWSRPCUrl != "127.0.0.1:61002" {
		t.Fatalf("expected external wsrpc url to be rewritten to 127.0.0.1:61002, got %s", output.ExternalWSRPCUrl)
	}
	if !strings.HasSuffix(output.InternalWSRPCUrl, ":61002") {
		t.Fatalf("expected internal wsrpc url to use tunneled port 61002, got %s", output.InternalWSRPCUrl)
	}
}
