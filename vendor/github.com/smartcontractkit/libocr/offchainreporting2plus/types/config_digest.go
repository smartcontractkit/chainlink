package types

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// ConfigDigestPrefix acts as a domain separator between different (typically
// chain-specific) methods of computing a ConfigDigest. The Prefix is encoded
// in big-endian.
type ConfigDigestPrefix uint16

// This acts as the canonical "registry" of ConfigDigestPrefixes. Pick an unused
// prefix and add it to this list before you build an OffchainConfigDigester for
// whatever chain you're targeting.
const (
	_                                        ConfigDigestPrefix = 0 // reserved to prevent errors where a zero-default creeps through somewhere
	ConfigDigestPrefixEVMSimple              ConfigDigestPrefix = 0x0001
	ConfigDigestPrefixTerra                  ConfigDigestPrefix = 0x0002
	ConfigDigestPrefixSolana                 ConfigDigestPrefix = 0x0003
	ConfigDigestPrefixStarknet               ConfigDigestPrefix = 0x0004
	_                                                           = 0x0005 // reserved, not sure for what
	ConfigDigestPrefixMercuryV02             ConfigDigestPrefix = 0x0006 // Mercury v0.2 and v0.3
	ConfigDigestPrefixEVMThresholdDecryption ConfigDigestPrefix = 0x0007 // Run Threshold/S4 plugins as part of another product under one contract.
	ConfigDigestPrefixEVMS4                  ConfigDigestPrefix = 0x0008 // Run Threshold/S4 plugins as part of another product under one contract.
	ConfigDigestPrefixLLO                    ConfigDigestPrefix = 0x0009 // Mercury v1
	ConfigDigestPrefixCCIPMultiRole          ConfigDigestPrefix = 0x000a // CCIP multi role

	ConfigDigestPrefixOCR1 ConfigDigestPrefix = 0xEEEE // we translate ocr1 config digest to ocr2 config digests in the networking layer
	_                      ConfigDigestPrefix = 0xFFFF // reserved for future use

	// Deprecated: Use equivalent ConfigDigestPrefixEVMSimple instead
	ConfigDigestPrefixEVM ConfigDigestPrefix = ConfigDigestPrefixEVMSimple
)

func ConfigDigestPrefixFromConfigDigest(configDigest ConfigDigest) ConfigDigestPrefix {
	return ConfigDigestPrefix(binary.BigEndian.Uint16(configDigest[:2]))
}

// Checks whether a ConfigDigestPrefix is actually a prefix of a ConfigDigest.
func (prefix ConfigDigestPrefix) IsPrefixOf(configDigest ConfigDigest) bool {
	return prefix == ConfigDigestPrefixFromConfigDigest(configDigest)
}

var _ fmt.Stringer = ConfigDigestPrefix(0)

func (prefix ConfigDigestPrefix) String() string {
	var encoded [2]byte
	binary.BigEndian.PutUint16(encoded[:], uint16(prefix))
	return fmt.Sprintf("%x", encoded)
}

// Digest of the configuration for a OCR2 protocol instance. The first two
// bytes indicate which config digester (typically specific to a targeted
// blockchain) was used to compute a ConfigDigest. This value is used as a
// domain separator between different protocol instances and is thus security
// critical. It should be the output of a cryptographic hash function over all
// relevant configuration fields as well as e.g. the address of the target
// contract/state accounts/...
type ConfigDigest [32]byte

func (c ConfigDigest) Hex() string {
	return fmt.Sprintf("%x", c[:])
}

func BytesToConfigDigest(b []byte) (ConfigDigest, error) {
	configDigest := ConfigDigest{}

	if len(b) != len(configDigest) {
		return ConfigDigest{}, fmt.Errorf("cannot convert bytes to ConfigDigest. bytes have wrong length %v", len(b))
	}

	if len(configDigest) != copy(configDigest[:], b) {
		// assertion
		panic("copy returned wrong length")
	}

	return configDigest, nil
}

var _ fmt.Stringer = ConfigDigest{}

func (c ConfigDigest) String() string {
	return c.Hex()
}

var _ encoding.TextMarshaler = ConfigDigest{}

func (c ConfigDigest) MarshalText() (text []byte, err error) {
	s := c.String()
	return []byte(s), nil
}

// An OffchainConfigDigester computes a ConfigDigest the same way as the
// contract, but *offchain*. This is used to ensure that the ConfigDigest
// returned from the contract was computed correctly and to prevent a malicious
// blockchain node from breaking domain separation between different protocol
// instances.
//
// All its functions should be thread-safe.
type OffchainConfigDigester interface {
	// Compute ConfigDigest for the given ContractConfig. The first two bytes of the
	// ConfigDigest must be the big-endian encoding of ConfigDigestPrefix!
	ConfigDigest(ContractConfig) (ConfigDigest, error)

	// This should return the same constant value on every invocation
	ConfigDigestPrefix() (ConfigDigestPrefix, error)
}
