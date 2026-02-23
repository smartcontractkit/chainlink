package cre

import "testing"

func TestResolveNodeFacingJDUriForDON_LocalDonToLocalJD_UsesInternal(t *testing.T) {
	donMeta := &DonMetadata{
		Name: "workflow",
		ns: &NodeSet{
			Placement: "local",
		},
	}

	got, err := resolveNodeFacingJDUriForDON(donMeta, "local", "jd:8080", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("resolveNodeFacingJDUriForDON returned error: %v", err)
	}
	if got != "jd:8080" {
		t.Fatalf("expected internal JD URI jd:8080, got %s", got)
	}
}

func TestResolveNodeFacingJDUriForDON_RemoteDonToLocalJD_RewritesForBridge(t *testing.T) {
	donMeta := &DonMetadata{
		Name: "workflow",
		ns: &NodeSet{
			Placement: "remote",
		},
	}

	got, err := resolveNodeFacingJDUriForDON(donMeta, "local", "jd:8080", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("resolveNodeFacingJDUriForDON returned error: %v", err)
	}
	if got != "host.docker.internal:8080" {
		t.Fatalf("expected bridged JD URI host.docker.internal:8080, got %s", got)
	}
}
