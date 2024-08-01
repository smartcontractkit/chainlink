package managed

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// Convenience wrapper around OffchainConfigDigester
type prefixCheckConfigDigester struct {
	offchainConfigDigester types.OffchainConfigDigester
}

// ConfigDigest method that checks that the computed ConfigDigest's prefix is
// consistent with OffchainConfigDigester.ConfigDigestPrefix
func (d prefixCheckConfigDigester) ConfigDigest(cc types.ContractConfig) (types.ConfigDigest, error) {
	prefix, err := d.offchainConfigDigester.ConfigDigestPrefix()
	if err != nil {
		return types.ConfigDigest{}, err
	}

	cd, err := d.offchainConfigDigester.ConfigDigest(cc)
	if err != nil {
		return types.ConfigDigest{}, err
	}

	if !prefix.IsPrefixOf(cd) {
		return types.ConfigDigest{}, fmt.Errorf("ConfigDigest has prefix %s, but wanted prefix %s", types.ConfigDigestPrefixFromConfigDigest(cd), prefix)
	}

	return cd, nil
}

// Check that the ContractConfig's ConfigDigest matches the one computed
// offchain
func (d prefixCheckConfigDigester) CheckContractConfig(cc types.ContractConfig) error {
	goodConfigDigest, err := d.ConfigDigest(cc)
	if err != nil {
		return err
	}

	if goodConfigDigest != cc.ConfigDigest {
		return fmt.Errorf("ConfigDigest mismatch. Expected %s but got %s", goodConfigDigest, cc.ConfigDigest)
	}

	return nil
}
