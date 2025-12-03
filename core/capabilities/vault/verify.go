package vault

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/smdkg/dkgocr"
	"github.com/smartcontractkit/smdkg/dkgocr/dkgocrtypes"
	"github.com/smartcontractkit/smdkg/dkgocr/tdh2shim"
)

func VerifyDKGResult(resultPackage []byte, masterPublicKey string, key dkgocrtypes.P256Keyring) error {
	rp := dkgocr.NewResultPackage()
	err := rp.UnmarshalBinary(resultPackage)
	if err != nil {
		return fmt.Errorf("could not unmarshal result package: %w", err)
	}

	tdh2PubKey, err := tdh2shim.TDH2PublicKeyFromDKGResult(rp)
	if err != nil {
		return fmt.Errorf("could not derive TDH2 public key from DKG result: %w", err)
	}

	pubKeyBytes, err := tdh2PubKey.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal TDH2 public key: %w", err)
	}

	mpk, err := hex.DecodeString(masterPublicKey)
	if err != nil {
		return fmt.Errorf("could not hex decode master public key from request: %w", err)
	}

	if !bytes.Equal(pubKeyBytes, mpk) {
		return fmt.Errorf("master public key does not match: got %x, want %x", pubKeyBytes, mpk)
	}

	_, err = rp.MasterSecretKeyShare(key)
	if err != nil {
		return fmt.Errorf("could not decrypt share with DKG recipient key: %w", err)
	}

	return nil
}
