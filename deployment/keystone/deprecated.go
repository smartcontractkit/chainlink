package keystone

import "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

//TODO: delete this after the downstream migration is done

// Deprecated: Use changeset package instead
// OracleConfig is the configuration for an oracle
type OracleConfig = changeset.OracleConfig

// Deprecated: Use changeset package instead
// OCR3OnchainConfig is the onchain configuration of an OCR2 contract
type OCR2OracleConfig = changeset.OCR3OnchainConfig

// Deprecated: Use changeset package instead
// NodeKeys is a set of public keys for a node
type NodeKeys = changeset.NodeKeys

// Deprecated: Use changeset package instead
// TopLevelConfigSource is the top level configuration source
type TopLevelConfigSource = changeset.TopLevelConfigSource

// Deprecated: Use changeset package instead
// GenerateOCR3Config generates an OCR3 config
var GenerateOCR3Config = changeset.GenerateOCR3Config

// Deprecated: Use changeset package instead
// FeedConsumer is a feed consumer contract type
var FeedConsumer = changeset.FeedConsumer

// Deprecated: Use changeset package instead
// KeystoneForwarder is a keystone forwarder contract type
var KeystoneForwarder = changeset.KeystoneForwarder

// Deprecated: Use changeset package instead
// GetContractSetsRequest is a request to get contract sets
type GetContractSetsRequest = changeset.GetContractSetsRequest

// Deprecated: Use changeset package instead
// GetContractSetsResponse is a response to get contract sets
type GetContractSetsResponse = changeset.GetContractSetsResponse

// Deprecated: Use changeset package instead
// GetContractSets gets contract sets
var GetContractSets = changeset.GetContractSets

// Deprecated: Use changeset package instead
// RegisterCapabilitiesRequest is a request to register capabilities
type RegisterCapabilitiesRequest = changeset.RegisterCapabilitiesRequest

// Deprecated: Use changeset package instead
// RegisterCapabilitiesResponse is a response to register capabilities
type RegisterCapabilitiesResponse = changeset.RegisterCapabilitiesResponse

// Deprecated: Use changeset package instead
// RegisterCapabilities registers capabilities
var RegisterCapabilities = changeset.RegisterCapabilities

// Deprecated: Use changeset package instead
// RegisterNOPSRequest is a request to register NOPS
type RegisterNOPSRequest = changeset.RegisterNOPSRequest

// Deprecated: Use changeset package instead
// RegisterNOPSResponse is a response to register NOPS
type RegisterNOPSResponse = changeset.RegisterNOPSResponse

// Deprecated: Use changeset package instead
// RegisterNOPS registers NOPS
var RegisterNOPS = changeset.RegisterNOPS

// Deprecated: Use changeset package instead
// RegisterNodesRequest is a request to register nodes with the capabilities registry
type RegisterNodesRequest = changeset.RegisterNodesRequest

// Deprecated: Use changeset package instead
// RegisterNodesResponse is a response to register nodes with the capabilities registry
type RegisterNodesResponse = changeset.RegisterNodesResponse

// Deprecated: Use changeset package instead
// RegisterNodes registers nodes with the capabilities registry
var RegisterNodes = changeset.RegisterNodes

// Deprecated: Use changeset package instead
// RegisteredCapability is a wrapper of a capability and its ID
type RegisteredCapability = changeset.RegisteredCapability

// Deprecated: Use changeset package instead
// FromCapabilitiesRegistryCapability converts a capabilities registry capability to a registered capability
var FromCapabilitiesRegistryCapability = changeset.FromCapabilitiesRegistryCapability

// Deprecated: Use changeset package instead
// RegisterDonsRequest is a request to register Dons with the capabilities registry
type RegisterDonsRequest = changeset.RegisterDonsRequest

// Deprecated: Use changeset package instead
// RegisterDonsResponse is a response to register Dons with the capabilities registry
type RegisterDonsResponse = changeset.RegisterDonsResponse

// Deprecated: Use changeset package instead
// RegisterDons registers Dons with the capabilities registry
var RegisterDons = changeset.RegisterDons

// Deprecated: Use changeset package instead
// DONToRegister is the minimal information needed to register a DON with the capabilities registry
type DONToRegister = changeset.DONToRegister

// Deprecated: Use changeset package instead
// ConfigureContractsRequest is a request to configure ALL the contracts
type ConfigureContractsRequest = changeset.ConfigureContractsRequest

// Deprecated: Use changeset package instead
// ConfigureContractsResponse is a response to configure ALL the contracts
type ConfigureContractsResponse = changeset.ConfigureContractsResponse

// Deprecated: Use changeset package instead
// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convenient form to handle the CLO representation of the nop data
type DonCapabilities = changeset.DonCapabilities

// Deprecated: Use changeset package instead
type DeployRequest = changeset.DeployRequest

// Deprecated: Use changeset package instead
type DeployResponse = changeset.DeployResponse

// Deprecated: Use changeset package instead
type ContractSet = changeset.ContractSet
