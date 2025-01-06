package keystone

import "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"

//TODO: delete this after the downstream migration is done

// DEPRECATED: Use changeset package instead
// OracleConfig is the configuration for an oracle
type OracleConfig = changeset.OracleConfig

// DEPRECATED: Use changeset package instead
// OCR3OnchainConfig is the onchain configuration of an OCR2 contract
type OCR2OracleConfig = changeset.OCR3OnchainConfig

// DEPRECATED: Use changeset package instead
// NodeKeys is a set of public keys for a node
type NodeKeys = changeset.NodeKeys

// DEPRECATED: Use changeset package instead
// TopLevelConfigSource is the top level configuration source
type TopLevelConfigSource = changeset.TopLevelConfigSource

// DEPRECATED: Use changeset package instead
// GenerateOCR3Config generates an OCR3 config
var GenerateOCR3Config = changeset.GenerateOCR3Config

// DEPRECATED: Use changeset package instead
// FeedConsumer is a feed consumer contract type
var FeedConsumer = changeset.FeedConsumer

// DEPRECATED: Use changeset package instead
// KeystoneForwarder is a keystone forwarder contract type
var KeystoneForwarder = changeset.KeystoneForwarder

// DEPRECATED: Use changeset package instead
// GetContractSetsRequest is a request to get contract sets
type GetContractSetsRequest = changeset.GetContractSetsRequest

// DEPRECATED: Use changeset package instead
// GetContractSetsResponse is a response to get contract sets
type GetContractSetsResponse = changeset.GetContractSetsResponse

// DEPRECATED: Use changeset package instead
// GetContractSets gets contract sets
var GetContractSets = changeset.GetContractSets

// DEPRECATED: Use changeset package instead
// RegisterCapabilitiesRequest is a request to register capabilities
type RegisterCapabilitiesRequest = changeset.RegisterCapabilitiesRequest

// DEPRECATED: Use changeset package instead
// RegisterCapabilitiesResponse is a response to register capabilities
type RegisterCapabilitiesResponse = changeset.RegisterCapabilitiesResponse

// DEPRECATED: Use changeset package instead
// RegisterCapabilities registers capabilities
var RegisterCapabilities = changeset.RegisterCapabilities

// DEPRECATED: Use changeset package instead
// RegisterNOPSRequest is a request to register NOPS
type RegisterNOPSRequest = changeset.RegisterNOPSRequest

// DEPRECATED: Use changeset package instead
// RegisterNOPSResponse is a response to register NOPS
type RegisterNOPSResponse = changeset.RegisterNOPSResponse

// DEPRECATED: Use changeset package instead
// RegisterNOPS registers NOPS
var RegisterNOPS = changeset.RegisterNOPS

// DEPRECATED: Use changeset package instead
// RegisterNodesRequest is a request to register nodes with the capabilities registry
type RegisterNodesRequest = changeset.RegisterNodesRequest

// DEPRECATED: Use changeset package instead
// RegisterNodesResponse is a response to register nodes with the capabilities registry
type RegisterNodesResponse = changeset.RegisterNodesResponse

// DEPRECATED: Use changeset package instead
// RegisterNodes registers nodes with the capabilities registry
var RegisterNodes = changeset.RegisterNodes

// DEPRECATED: Use changeset package instead
// RegisteredCapability is a wrapper of a capability and its ID
type RegisteredCapability = changeset.RegisteredCapability

// DEPRECATED: Use changeset package instead
// FromCapabilitiesRegistryCapability converts a capabilities registry capability to a registered capability
var FromCapabilitiesRegistryCapability = changeset.FromCapabilitiesRegistryCapability

// DEPRECATED: Use changeset package instead
// RegisterDonsRequest is a request to register Dons with the capabilities registry
type RegisterDonsRequest = changeset.RegisterDonsRequest

// DEPRECATED: Use changeset package instead
// RegisterDonsResponse is a response to register Dons with the capabilities registry
type RegisterDonsResponse = changeset.RegisterDonsResponse

// DEPRECATED: Use changeset package instead
// RegisterDons registers Dons with the capabilities registry
var RegisterDons = changeset.RegisterDons

// DEPRECATED: Use changeset package instead
// DONToRegister is the minimal information needed to register a DON with the capabilities registry
type DONToRegister = changeset.DONToRegister

// DEPRECATED: Use changeset package instead
// ConfigureContractsRequest is a request to configure ALL the contracts
type ConfigureContractsRequest = changeset.ConfigureContractsRequest

// DEPRECATED: Use changeset package instead
// ConfigureContractsResponse is a response to configure ALL the contracts
type ConfigureContractsResponse = changeset.ConfigureContractsResponse

// DEPRECATED: Use changeset package instead
// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convenient form to handle the CLO representation of the nop data
type DonCapabilities = changeset.DonCapabilities

// DEPRECATED: Use changeset package instead
type DeployRequest = changeset.DeployRequest

// DEPRECATED: Use changeset package instead
type DeployResponse = changeset.DeployResponse

// DEPRECATED: Use changeset package instead
type ContractSet = changeset.ContractSet
