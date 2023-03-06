// Package types contains the types and interfaces a consumer of the OCR library needs to be aware of
package types

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/commontypes"
	"golang.org/x/crypto/curve25519"
)

type ConfigDigest [16]byte

func (c ConfigDigest) Hex() string {
	return fmt.Sprintf("%x", c[:])
}

func BytesToConfigDigest(b []byte) (ConfigDigest, error) {
	configDigest := ConfigDigest{}

	if len(b) != len(configDigest) {
		return ConfigDigest{}, fmt.Errorf("Cannot convert bytes to ConfigDigest. bytes have wrong length %v", len(b))
	}

	if n := copy(configDigest[:], b); n != len(configDigest) {
		panic("copy returned wrong length")
	}

	return configDigest, nil
}

// BinaryNetworkEndpointFactory creates permissioned BinaryNetworkEndpoints.
//
// All its functions should be thread-safe.
type BinaryNetworkEndpointFactory interface {
	NewEndpoint(cd ConfigDigest, peerIDs []string,
		v1bootstrappers []string, v2bootstrappers []commontypes.BootstrapperLocator,
		failureThreshold int, tokenBucketRefillRate float64, tokenBucketSize int,
	) (commontypes.BinaryNetworkEndpoint, error)
	PeerID() string
}

// BootstrapperFactory creates permissioned Bootstrappers.
//
// All its functions should be thread-safe.
type BootstrapperFactory interface {
	NewBootstrapper(cd ConfigDigest, peerIDs []string,
		v1bootstrappers []string, v2bootstrappers []commontypes.BootstrapperLocator,
		failureThreshold int,
	) (commontypes.Bootstrapper, error)
}

// Observation is the type returned by the DataSource.Observe method. Represents
// an int192 at time of writing
type Observation *big.Int

// DataSource implementations must be thread-safe. Observe may be called by many different threads concurrently.
type DataSource interface {
	// Observe queries the data source. Returns a value or an error. Once the
	// context is expires, Observe may still do cheap computations and return a
	// result, but should return as quickly as possible.
	//
	// More details: In the current implementation, the context passed to
	// Observe will time out after LocalConfig.DataSourceTimeout. However,
	// Observe should *not* make any assumptions about context timeout behavior.
	// Once the context times out, Observe should prioritize returning as
	// quickly as possible, but may still perform fast computations to return a
	// result rather than error. For example, if Observe medianizes a number
	// of data sources, some of which already returned a result to Observe prior
	// to the context's expiry, Observe might still compute their median, and
	// return it instead of an error.
	//
	// Important: Observe should not perform any potentially time-consuming
	// actions like database access, once the context passed has expired.
	Observe(context.Context) (Observation, error)
}

type ConfigOverride struct {
	AlphaPPB uint64
	DeltaC   time.Duration
}

// ConfigOverrider allows overriding some OCR protocol configuration parameters.
//
// All its functions should be thread-safe.
type ConfigOverrider interface {
	// Enables locally overriding the configuration parameters in
	// ConfigOverride by returning a non-nil result.  If no override
	// is desired, return nil.
	//
	// This function is expected to return immediately.
	ConfigOverride() *ConfigOverride
}

// ContractTransmitter sends new reports to the OffchainAggregator smart contract.
//
// All its functions should be thread-safe.
type ContractTransmitter interface {

	// Transmit sends the report to the on-chain OffchainAggregator smart contract's Transmit method
	Transmit(
		ctx context.Context,
		report []byte, // wire-formatted report to transmit on-chain
		rs, ss [][32]byte, vs [32]byte, // Signatures; i'th elt's are i'th (v,r,s)
	) error

	LatestTransmissionDetails(
		ctx context.Context,
	) (
		configDigest ConfigDigest,
		epoch uint32,
		round uint8,
		latestAnswer Observation,
		latestTimestamp time.Time,
		err error,
	)

	// LatestRoundRequested returns the configDigest, epoch, and round from the latest
	// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
	// return a result if the latest such event was emitted in a block b such that
	// b.timestamp < tip.timestamp - lookback.
	//
	// If no event is found, LatestRoundRequested should return zero values, not an error.
	// An error should only be returned if an actual error occurred during execution,
	// e.g. because there was an error querying the blockchain or the database.
	//
	// As an optimization, this function may also return zero values, if no
	// RoundRequested event has been emitted after the latest NewTransmission event.
	LatestRoundRequested(
		ctx context.Context,
		lookback time.Duration,
	) (
		configDigest ConfigDigest,
		epoch uint32,
		round uint8,
		err error,
	)

	FromAddress() common.Address

	ChainID() *big.Int
}

// ContractConfigTracker tracks OffchainAggregator.ConfigSet events emitted from blockchain.
//
// All its functions should be thread-safe.
type ContractConfigTracker interface {
	SubscribeToNewConfigs(ctx context.Context) (ContractConfigSubscription, error)
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ConfigDigest, err error)
	ConfigFromLogs(ctx context.Context, changedInBlock uint64) (ContractConfig, error)

	// LatestBlockHeight returns the height of the most recent block in the chain.
	LatestBlockHeight(ctx context.Context) (blockheight uint64, err error)
}

type ContractConfigSubscription interface {
	// May be closed by sender at any time
	Configs() <-chan ContractConfig
	// Calling this multiple times may return an error, but must not panic.
	Close()
}

type ContractConfig struct {
	ConfigDigest         ConfigDigest
	Signers              []common.Address // TODO: use OnChainSigningAddress?
	Transmitters         []common.Address
	Threshold            uint8
	EncodedConfigVersion uint64
	Encoded              []byte
}

// OffChainPublicKey is the public key used to cryptographically identify an
// oracle in inter-oracle communications.
type OffchainPublicKey ed25519.PublicKey

// OnChainSigningAddress is the public key used to cryptographically identify an
// oracle to the on-chain smart contract.
type OnChainSigningAddress common.Address

// SharedSecretEncryptionPublicKey is the public key used to receive an encrypted
// version of the secret shared amongst all oracles on a common contract.
type SharedSecretEncryptionPublicKey [curve25519.PointSize]byte // X25519

// PrivateKeys contains the secret keys needed for the OCR protocol, and methods
// which use those keys without exposing them to the rest of the application.
// There are three key pairs to track, here:
//
// - The on-chain signing key pair (secp256k1), used to sign contract reports
//
// - The off-chain key signing key pair (Ed25519), used to sign observations
//
// - The config encryption key (X25519), used to decrypt the symmetric
// key which encrypts the offchain configuration data passed through the OffchainAggregator
// smart contract.
//
// All its functions should be thread-safe.
type PrivateKeys interface {

	// SignOnChain returns an ethereum-style ECDSA secp256k1 signature on msg. See
	// signature.OnChainPrivateKey.Sign for the logic it needs to implement
	SignOnChain(msg []byte) (signature []byte, err error)

	// SignOffChain returns an EdDSA-Ed25519 signature on msg. See
	// signature.OffChainPrivateKey.Sign for the logic it needs to implement
	SignOffChain(msg []byte) (signature []byte, err error)

	// ConfigDiffieHellman multiplies base, as a representative of a Curve 25519
	// point, by a secret scalar, which is also the scalar to multiply
	// curve25519.BasePoint to, in order to get PublicKeyConfig
	ConfigDiffieHellman(base *[curve25519.ScalarSize]byte) (sharedPoint *[curve25519.PointSize]byte, err error)

	// PublicKeyAddressOnChain returns the address corresponding to the
	// public component of the keypair used in SignOnChain
	PublicKeyAddressOnChain() OnChainSigningAddress

	// PublicKeyOffChain returns the pbulic component of the keypair used in SignOffChain
	PublicKeyOffChain() OffchainPublicKey

	// PublicKeyConfig returns the public component of the keypair used in ConfigKeyShare
	PublicKeyConfig() [curve25519.PointSize]byte
}
