// Command vault-dkg-decode decodes an exported vault DKG result package and
// prints its contents, including the TDH2 master public key in the exact form
// CapReg stores it (the VaultPublicKey `stringValue`).
//
// It is meant to be run whenever a vault DKG reshare (or fresh dealing) happens
// so you can (a) confirm what the result package contains and (b) get the
// CapReg-ready public-key value to publish. On a reshare the master key
// (G_bar/H) is preserved but the per-recipient verification shares (HArray) are
// redistributed to the new committee, so CapReg's VaultPublicKey MUST be
// refreshed to the value this tool prints or `secrets.get` decryption-share
// verification will fail.
//
// Usage:
//
//	# from a file containing the hex result package (no 0x / whitespace ok):
//	go run . ./resultpkg.hex
//
//	# from stdin:
//	cat resultpkg.hex | go run . -
//	curl ... | jq -r .hexDKGResultPackage | go run . -
//
// Flags:
//
//	--capreg   print ONLY the CapReg VaultPublicKey stringValue (hex) and exit.
//	           This is what you paste into the capability_registry_update_don
//	           input's VaultPublicKey.stringValue.
//	--json     print ONLY the decoded TDH2 public key JSON and exit.
//
// With no flag it prints a full human-readable summary followed by the CapReg
// stringValue.
package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/smartcontractkit/smdkg/dkgocr"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
)

func main() {
	capregOnly := flag.Bool("capreg", false, "print only the CapReg VaultPublicKey stringValue (hex) and exit")
	jsonOnly := flag.Bool("json", false, "print only the decoded TDH2 public key JSON and exit")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: vault-dkg-decode [--capreg|--json] <hexfile|->")
		fmt.Fprintln(os.Stderr, "  reads a hex-encoded DKG result package from a file, or from stdin if the arg is '-'")
		os.Exit(2)
	}

	hexStr, err := readHexInput(flag.Arg(0))
	must(err)
	raw, err := hex.DecodeString(hexStr)
	must(err)

	rp := dkgocr.NewResultPackage()
	must(rp.UnmarshalBinary(raw))

	// The CapReg VaultPublicKey value: TDH2 pubkey marshaled to JSON, then the
	// UTF-8 JSON bytes hex-encoded (this is exactly what CapReg stores as
	// stringValue).
	tdh2Pub, err := tdh2shim.TDH2PublicKeyFromDKGResult(rp)
	must(err)
	pkJSON, err := tdh2Pub.Marshal()
	must(err)
	capregStringValue := hex.EncodeToString(pkJSON)

	if *capregOnly {
		fmt.Println(capregStringValue)
		return
	}
	if *jsonOnly {
		fmt.Println(string(pkJSON))
		return
	}

	cfg := rp.ReportingPluginConfig()
	fmt.Println("== DKG result package ==")
	fmt.Println("InstanceID          :", rp.InstanceID())
	fmt.Println("Config.T (threshold):", cfg.T)
	fmt.Println("Config #dealers     :", len(cfg.DealerPublicKeys))
	fmt.Println("Config #recipients  :", len(cfg.RecipientPublicKeys))
	if cfg.PreviousInstanceID != nil {
		fmt.Println("PreviousInstanceID  :", *cfg.PreviousInstanceID, "  (=> RESHARE: master key preserved)")
	} else {
		fmt.Println("PreviousInstanceID  : <nil>   (=> FRESH dealing: new master key)")
	}

	mpk := rp.MasterPublicKey()
	fmt.Printf("MasterPublicKey(raw): %x\n", []byte(mpk))
	shares := rp.MasterPublicKeyShares()
	fmt.Println("#MasterPubKeyShares :", len(shares), " (== HArray length == committee size)")

	// Pretty-print the decoded pubkey fields so G_bar/H can be eyeballed
	// against the previous instance (they must match on a reshare).
	var m struct {
		Group  string   `json:"Group"`
		GBar   string   `json:"G_bar"`
		H      string   `json:"H"`
		HArray []string `json:"HArray"`
	}
	if json.Unmarshal(pkJSON, &m) == nil {
		fmt.Println()
		fmt.Println("== TDH2 public key (CapReg VaultPublicKey, decoded) ==")
		fmt.Println("Group               :", m.Group)
		fmt.Println("G_bar               :", m.GBar, " (stable across resharings)")
		fmt.Println("H                   :", m.H, " (master public key; stable across resharings)")
		fmt.Println("HArray entries      :", len(m.HArray), " (per-recipient verification shares; CHANGES on reshare)")
	}

	fmt.Println()
	fmt.Println("== CapReg VaultPublicKey stringValue (paste into capability_registry_update_don) ==")
	fmt.Println(capregStringValue)
}

// readHexInput returns the hex string from the given path, or from stdin if
// path is "-". Any surrounding whitespace/newlines and a leading 0x are removed.
func readHexInput(path string) (string, error) {
	var b []byte
	var err error
	if path == "-" {
		b, err = io.ReadAll(os.Stdin)
	} else {
		b, err = os.ReadFile(path)
	}
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(b))
	s = strings.TrimPrefix(s, "0x")
	s = strings.Join(strings.Fields(s), "") // strip any internal whitespace/newlines
	return s, nil
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}
