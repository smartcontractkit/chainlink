package changeset

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	commonTypes "github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

func ValidateCacheForChain(env deployment.Environment, chainSelector uint64, cacheAddress common.Address) error {
	state, err := LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load on chain state %w", err)
	}
	_, ok := env.Chains[chainSelector]
	if !ok {
		return errors.New("chain not found in environment")
	}
	chainState, ok := state.Chains[chainSelector]
	if !ok {
		return errors.New("chain not found in on chain state")
	}
	if chainState.DataFeedsCache == nil {
		return errors.New("DataFeedsCache not found in on chain state")
	}
	_, ok = chainState.DataFeedsCache[cacheAddress]
	if !ok {
		return errors.New("contract not found in on chain state")
	}
	return nil
}

func ValidateMCMSAddresses(ab deployment.AddressBook, chainSelector uint64) error {
	if _, err := deployment.SearchAddressBook(ab, chainSelector, commonTypes.RBACTimelock); err != nil {
		return fmt.Errorf("timelock not present on the chain %w", err)
	}
	if _, err := deployment.SearchAddressBook(ab, chainSelector, commonTypes.ProposerManyChainMultisig); err != nil {
		return fmt.Errorf("mcms proposer not present on the chain %w", err)
	}
	return nil
}

// spec- https://docs.google.com/document/d/13ciwTx8lSUfyz1IdETwpxlIVSn1lwYzGtzOBBTpl5Vg/edit?tab=t.0#heading=h.dxx2wwn1dmoz
func ValidateFeedID(feedID string) error {
	// Check for "0x" prefix and remove it
	if feedID[:2] != "0x" {
		return errors.New("invalid feed ID")
	}
	feedID = feedID[2:]

	if len(feedID) != 32 {
		return errors.New("invalid feed ID length")
	}

	bytes, err := hex.DecodeString(feedID)
	if err != nil {
		return fmt.Errorf("invalid feed ID format: %w", err)
	}

	// Validate format byte
	format := bytes[0]
	if format != 0x01 && format != 0x02 {
		return errors.New("invalid format byte")
	}

	// bytes (1-4) are random bytes, so no validation needed

	// Validate attribute bucket bytes (5-6)
	attributeBucket := [2]byte{bytes[5], bytes[6]}
	attributeBucketHex := hex.EncodeToString(attributeBucket[:])
	if attributeBucketHex != "0003" && attributeBucketHex != "0700" {
		return errors.New("invalid attribute bucket bytes")
	}

	// Validate data type byte (7)
	dataType := bytes[7]
	validDataTypes := map[byte]bool{
		0x00: true, 0x01: true, 0x02: true, 0x03: true, 0x04: true,
	}
	for i := 0x20; i <= 0x60; i++ {
		validDataTypes[byte(i)] = true
	}
	if !validDataTypes[dataType] {
		return errors.New("invalid data type byte")
	}

	// Validate reserved bytes (8-15) are zero
	for i := 8; i < 16; i++ {
		if bytes[i] != 0x00 {
			return errors.New("reserved bytes must be zero")
		}
	}

	return nil
}
