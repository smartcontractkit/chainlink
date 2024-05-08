package capreg

import (
	"bytes"
	"context"
	"math/big"
)

// CapabilityID is the unique identifier of the capability in the CR.
// It is calculated as keccak256(abi.encode(capabilityType, capabilityVersion)).
type CapabilityID = string

// DON represents a DON in the Capability Registry.
type DON struct {
	// ID is the unique identifier of the DON in the CR.
	ID uint32

	// IsPublic indicates whether this DON's capabilities can be accessed publicly.
	IsPublic bool

	// Nodes is the list of nodes in this DON, represented by their RageP2P public keys/IDs.
	Nodes [][]byte

	// CapabilityConfigurations are the configurations of the various capabilities of this DON.
	CapabilityConfigurations []CapabilityConfiguration
}

// CapabilityConfiguration represents the configuration of a capability in the Capability Registry.
// This can be considered shared configuration among all DONs that have this capability.
type CapabilityConfiguration struct {
	// CapabilityID is the unique identifier of the capability in the CR that this configuration is for.
	CapabilityID          CapabilityID
	OnchainConfigVersion  uint64
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

func (c CapabilityConfiguration) Equal(other CapabilityConfiguration) bool {
	if c.CapabilityID != other.CapabilityID ||
		c.OnchainConfigVersion != other.OnchainConfigVersion ||
		c.OffchainConfigVersion != other.OffchainConfigVersion {
		return false
	}
	if !bytes.Equal(c.OnchainConfig, other.OnchainConfig) ||
		!bytes.Equal(c.OffchainConfig, other.OffchainConfig) {
		return false
	}
	return true
}

type NodeOperator struct {
	// ID is the unique identifier of the Node Operator in the CR.
	ID *big.Int

	// Admin is the address of the admin of the Node Operator.
	Admin string

	// Name is the name of the Node Operator.
	Name string
}

type Node struct {
	NodeOperatorID *big.Int
	// ID is the unique identifier of the Node in the CR,
	// which is a RageP2P ID.
	ID []byte

	// CapabilityIDs are the unique identifiers of the capabilities of this Node.
	CapabilityIDs []CapabilityID
}

type ResponseType int

const (
	ResponseTypeReport ResponseType = iota
	ResponseTypeObservationIdentical
)

// Capability represents a capability in the Capability Registry.
// These can be thought of roughly as "things DONs can do".
type Capability struct {
	// Type is the type of the capability, e.g trigger, action, consensus, etc.
	Type string

	// Version is the version of the capability, in semantic versioning notation, e.g 1.2.3.
	Version string

	// ID is the unique identifier of the capability in the CR.
	// It is calculated as keccak256(abi.encode(capabilityType, capabilityVersion)).
	ID CapabilityID

	// ResponseType indicates whether remote response requires aggregation
	// or is an OCR report. There are multiple ways to aggregate.
	ResponseType ResponseType

	// ConfigurationContractAddress is the address of the configuration contract for this capability.
	ConfigurationContractAddress string
}

// State mirrors the state in the onchain capability registry.
type State struct {
	Capabilities   map[CapabilityID]Capability
	CapabilityIDs  []CapabilityID
	DONs           []DON
	CapabilityDONs map[string][]DON
}

// Local is an interface for updating local capability registry state.
// Implementations of this can choose to filter out capabilities or DONS
// as they see fit.
// For example, for a capability router, it would need to know all DONs that
// support a particular capability.
// However, for a local capability launcher, it may only need to know about
// capabilities and DONs that the local node participates in.
type Local interface {
	// Sync is called when the onchain state of the capability registry is fetched.
	// It is up to the implementation to decide how to handle the new state.
	Sync(ctx context.Context, s State) error

	// Close is called when the local registry is shutting down.
	Close() error
}

// CapabilityFactory is an interface that must be implemented by capability factories that get started
// via the capability registry synchronization process.
// Capability factories are constructed outside of the capability registry syncer context and are expected
// to have all of the dependencies required in order to create the capability service, e.g relayer config,
// pipeline runner, etc.
//
//go:generate mockery --name CapabilityFactory --inpackage --inpackage-suffix --case underscore
type CapabilityFactory interface {
	// CapabilityID returns the unique identifier of the capability being implemented.
	CapabilityID() CapabilityID

	// Start is called on startup of the local registry or when a new DON is added
	// to the capability registry while the local registry is running.
	Start(ctx context.Context, d DON) error

	// Stop is called upon a DON being removed from the
	// capability registry while the local registry is running.
	Stop(ctx context.Context, d DON) error

	// Close is called upon shutdown of the local registry.
	Close() error

	// Update is called when the local registry is updated with new state while it is running.
	// This is separated from the Start/Stop flow so as to allow downstream services more granular
	// control of their own lifecycles, e.g using OCR's ContractConfigTracker to update OCR instances.
	Update(ctx context.Context, d DON) error
}
