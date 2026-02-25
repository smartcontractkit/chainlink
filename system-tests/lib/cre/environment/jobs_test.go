package environment

import (
	"testing"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/jd"
)

func TestRewriteJDForDirectAccessRewritesExternalEndpoints(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1:14231",
		ExternalWSRPCUrl: "127.0.0.1:9080",
		InternalWSRPCUrl: "job-distributor:8080",
	}

	if err := rewriteJDForDirectAccess(output, "10.20.30.40"); err != nil {
		t.Fatalf("rewriteJDForDirectAccess returned error: %v", err)
	}
	if output.ExternalGRPCUrl != "10.20.30.40:14231" {
		t.Fatalf("expected external grpc url to be rewritten, got %s", output.ExternalGRPCUrl)
	}
	if output.ExternalWSRPCUrl != "10.20.30.40:9080" {
		t.Fatalf("expected external wsrpc url to be rewritten, got %s", output.ExternalWSRPCUrl)
	}
	if output.InternalWSRPCUrl != "job-distributor:8080" {
		t.Fatalf("expected internal wsrpc url to remain unchanged, got %s", output.InternalWSRPCUrl)
	}
}

func TestRewriteJDForDirectAccessFallsBackToInternalWSRPCSource(t *testing.T) {
	output := &jd.Output{
		ExternalGRPCUrl:  "127.0.0.1:14231",
		ExternalWSRPCUrl: "",
		InternalWSRPCUrl: "job-distributor:8080",
	}

	if err := rewriteJDForDirectAccess(output, "10.20.30.40"); err != nil {
		t.Fatalf("rewriteJDForDirectAccess returned error: %v", err)
	}
	if output.ExternalWSRPCUrl != "10.20.30.40:8080" {
		t.Fatalf("expected external wsrpc url to be derived from internal source, got %s", output.ExternalWSRPCUrl)
	}
	if output.InternalWSRPCUrl != "job-distributor:8080" {
		t.Fatalf("expected internal wsrpc url to remain unchanged, got %s", output.InternalWSRPCUrl)
	}
}
