package engine

import (
	"testing"
)

func TestVersionSHA(t *testing.T) {
	t.Parallel()

	cases := []struct {
		version string
		want    string
	}{
		{"v0.0.0-20260714122420-7b2200a59a79", "7b2200a59a79"},
		{"v0.1.1-solana.0.20260625091148-e5618f5682ee", "e5618f5682ee"},
		{"v0.3.4-0.20260715161014-611d8ac32364", "611d8ac32364"},
		{"v1.0.5-0.20260629213843-c52e07523035", "c52e07523035"},
		{"v1.3.1-0.20260605202330-b5a89c32fdc1", "b5a89c32fdc1"},
		{"a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		{"v2.55.0", ""},         // clean release tag: no SHA
		{"v0.1.1-solana.0", ""}, // prerelease tag: no SHA
		{"", ""},
	}
	for _, c := range cases {
		if got := VersionSHA(c.version); got != c.want {
			t.Errorf("VersionSHA(%q) = %q, want %q", c.version, got, c.want)
		}
	}
}

const testGoMod = `module github.com/smartcontractkit/chainlink/v2

go 1.26.4

require (
	github.com/smartcontractkit/chainlink-aptos/codec v0.0.0-20260714122420-7b2200a59a79
	github.com/smartcontractkit/chainlink-ccip v0.1.1-solana.0.20260625091148-e5618f5682ee
	github.com/smartcontractkit/chainlink-ccip/chains/evm v0.0.0-20260624154507-ea7ff77a0ddb
	github.com/smartcontractkit/chainlink-ccip/chains/solana v0.0.0-20260415165642-49f23e4d76cc
	github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings v0.0.0-20260415165642-49f23e4d76cc
	github.com/smartcontractkit/chainlink-evm v0.3.4-0.20260715161014-611d8ac32364
	github.com/smartcontractkit/chainlink-sui/codec v0.0.0-20260714120433-7667cad5ff5c
	github.com/smartcontractkit/chainlink-ton v1.0.5-0.20260629213843-c52e07523035
	github.com/smartcontractkit/chainlink-common v0.0.0-20260101000000-aaaaaaaaaaaa
)
`

func TestParseGoMod(t *testing.T) {
	t.Parallel()

	mods, err := ParseGoMod([]byte(testGoMod), testProduct.Repos)
	if err != nil {
		t.Fatal(err)
	}
	// Tracked modules extracted.
	if got := mods["github.com/smartcontractkit/chainlink-ccip"]; got != "v0.1.1-solana.0.20260625091148-e5618f5682ee" {
		t.Errorf("chainlink-ccip version = %q", got)
	}
	if got := mods["github.com/smartcontractkit/chainlink-ton"]; got != "v1.0.5-0.20260629213843-c52e07523035" {
		t.Errorf("chainlink-ton version = %q", got)
	}
	// Untracked modules ignored.
	if _, ok := mods["github.com/smartcontractkit/chainlink-common"]; ok {
		t.Error("chainlink-common should not be tracked")
	}
	// solana main module is not in root go.mod.
	if _, ok := mods["github.com/smartcontractkit/chainlink-solana"]; ok {
		t.Error("chainlink-solana should not be present")
	}
	if len(mods) != 8 {
		t.Errorf("expected 8 tracked modules, got %d: %v", len(mods), mods)
	}
}

const testPluginsYAML = `defaults:
  goflags: "-ldflags=-s"

plugins:
  aptos:
    - moduleURI: "github.com/smartcontractkit/chainlink-aptos"
      gitRef: "v0.0.0-20260708114855-e953eeb028a7"
      installPath: "./cmd/chainlink-aptos"
  sui:
    - moduleURI: "github.com/smartcontractkit/chainlink-sui"
      # Must track the gRPC-migrated chainlink-sui used in go.mod
      gitRef: "v0.0.0-20260707125635-abec997b6eae"
      installPath: "./relayer/cmd/chainlink-sui"
  solana:
    - moduleURI: "github.com/smartcontractkit/chainlink-solana"
      gitRef: "v1.3.1-0.20260605202330-b5a89c32fdc1"
      installPath: "./pkg/solana/cmd/chainlink-solana"
  ton:
    - moduleURI: "github.com/smartcontractkit/chainlink-ton"
      gitRef: "v1.0.5-0.20260629213843-c52e07523035"
      installPath: "./cmd/chainlink-ton"
  evm:
    - moduleURI: "github.com/smartcontractkit/chainlink-evm"
      gitRef: "v0.3.4-0.20260715161014-611d8ac32364"
      installPath: "./pkg/cmd/chainlink-evm"
  starknet:
    - moduleURI: "github.com/smartcontractkit/chainlink-starknet/relayer"
      gitRef: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
      installPath: "./cmd/chainlink-starknet"
`

func TestParsePluginsYAML(t *testing.T) {
	t.Parallel()

	plugins, err := ParsePluginsYAML([]byte(testPluginsYAML), testProduct.Repos)
	if err != nil {
		t.Fatal(err)
	}
	if got := plugins["aptos"].GitRef; got != "v0.0.0-20260708114855-e953eeb028a7" {
		t.Errorf("aptos gitRef = %q", got)
	}
	if got := plugins["ton"].ModuleURI; got != "github.com/smartcontractkit/chainlink-ton" {
		t.Errorf("ton moduleURI = %q", got)
	}
	if _, ok := plugins["starknet"]; ok {
		t.Error("starknet should not be tracked")
	}
	if len(plugins) != 5 {
		t.Errorf("expected 5 tracked plugins, got %d: %v", len(plugins), plugins)
	}
}

func TestPinFor_PrimarySelection(t *testing.T) {
	t.Parallel()

	mods, _ := ParseGoMod([]byte(testGoMod), testProduct.Repos)
	plugins, _ := ParsePluginsYAML([]byte(testPluginsYAML), testProduct.Repos)
	snap := DepSnapshot{Ref: "v1", SHA: "sha", Modules: mods, Plugins: plugins}

	// ccip: no plugin -> primary go.mod module.
	ccip := pinFor(testProduct.Repos[0], snap)
	if ccip.PrimarySHA != "e5618f5682ee" {
		t.Errorf("ccip primary SHA = %q", ccip.PrimarySHA)
	}

	// aptos: plugin gitRef is primary, not the codec module.
	aptos := pinFor(testProduct.Repos[1], snap)
	if aptos.PrimarySHA != "e953eeb028a7" {
		t.Errorf("aptos primary SHA = %q, want plugin gitRef SHA", aptos.PrimarySHA)
	}

	// ton: plugin primary even though go.mod module exists (same SHA here).
	ton := pinFor(testProduct.Repos[4], snap)
	if ton.PrimarySHA != "c52e07523035" {
		t.Errorf("ton primary SHA = %q", ton.PrimarySHA)
	}
}
