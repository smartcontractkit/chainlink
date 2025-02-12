package changeset

import "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

//TODO: KS-673 refactor internal package to reduce and remove the duplication

// OracleConfig is the configuration for an oracle
type OracleConfig = internal.OracleConfig

// OCR3OnchainConfig is the onchain configuration of an OCR2 contract
type OCR3OnchainConfig = internal.OCR2OracleConfig

// NodeKeys is a set of public keys for a node
type NodeKeys = internal.NodeKeys

// TopLevelConfigSource is the top level configuration source
type TopLevelConfigSource = internal.TopLevelConfigSource

// GenerateOCR3Config generates an OCR3 config
var GenerateOCR3Config = internal.GenerateOCR3Config

// FeedConsumer is a feed consumer contract type
var FeedConsumer = internal.FeedConsumer

// KeystoneForwarder is a keystone forwarder contract type
var KeystoneForwarder = internal.KeystoneForwarder

// GetContractSetsRequest is a request to get contract sets
type GetContractSetsRequest = internal.GetContractSetsRequest

// GetContractSetsResponse is a response to get contract sets
type GetContractSetsResponse = internal.GetContractSetsResponse

// GetContractSets gets contract sets
var GetContractSets = internal.GetContractSets

// RegisterNOPSRequest is a request to register NOPS
type RegisterNOPSRequest = internal.RegisterNOPSRequest

// RegisterNOPSResponse is a response to register NOPS
type RegisterNOPSResponse = internal.RegisterNOPSResponse

// RegisterNOPS registers NOPS
var RegisterNOPS = internal.RegisterNOPS

// RegisterNodesRequest is a request to register nodes with the capabilities registry
type RegisterNodesRequest = internal.RegisterNodesRequest

// RegisterNodesResponse is a response to register nodes with the capabilities registry
type RegisterNodesResponse = internal.RegisterNodesResponse

// RegisterNodes registers nodes with the capabilities registry
var RegisterNodes = internal.RegisterNodes

// RegisteredCapability is a wrapper of a capability and its ID
type RegisteredCapability = internal.RegisteredCapability

// FromCapabilitiesRegistryCapability converts a capabilities registry capability to a registered capability
var FromCapabilitiesRegistryCapability = internal.FromCapabilitiesRegistryCapability

// RegisterDonsRequest is a request to register Dons with the capabilities registry
type RegisterDonsRequest = internal.RegisterDonsRequest

// RegisterDonsResponse is a response to register Dons with the capabilities registry
type RegisterDonsResponse = internal.RegisterDonsResponse

// RegisterDons registers Dons with the capabilities registry
var RegisterDons = internal.RegisterDons

// DONToRegister is the minimal information needed to register a DON with the capabilities registry
type DONToRegister = internal.DONToRegister

// NOP is a node operator profile, required to register a node with the capabilities registry
type NOP = internal.NOP

// ConfigureContractsRequest is a request to configure ALL the contracts
type ConfigureContractsRequest = internal.ConfigureContractsRequest

// ConfigureContractsResponse is a response to configure ALL the contracts
type ConfigureContractsResponse = internal.ConfigureContractsResponse

// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convenient form to handle the CLO representation of the nop data
type DonCapabilities = internal.DonCapabilities

type DONCapabilityWithConfig = internal.DONCapabilityWithConfig

type DeployRequest = internal.DeployRequest
type DeployResponse = internal.DeployResponse

type ContractSet = internal.ContractSet
