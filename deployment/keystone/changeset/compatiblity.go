package changeset

import "github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

type GetContractSetsRequest = internal.GetContractSetsRequest
type GetContractSetsResponse = internal.GetContractSetsResponse

var GetContractSets = internal.GetContractSets

// func GetContractSets  internal.GetContractSets
type OracleConfig = internal.OracleConfig

var FeedConsumer = internal.FeedConsumer
var KeystoneForwarder = internal.KeystoneForwarder

type RegisterCapabilitiesRequest = internal.RegisterCapabilitiesRequest
type RegisterCapabilitiesResponse = internal.RegisterCapabilitiesResponse

var RegisterCapabilities = internal.RegisterCapabilities

type RegisterNOPSRequest = internal.RegisterNOPSRequest
type RegisterNOPSResponse = internal.RegisterNOPSResponse

var RegisterNOPS = internal.RegisterNOPS

type RegisterNodesRequest = internal.RegisterNodesRequest
type RegisterNodesResponse = internal.RegisterNodesResponse

var RegisterNodes = internal.RegisterNodes

type RegisteredCapability = internal.RegisteredCapability

var FromCapabilitiesRegistryCapability = internal.FromCapabilitiesRegistryCapability

type RegisterDonsRequest = internal.RegisterDonsRequest
type RegisterDonsResponse = internal.RegisterDonsResponse

var RegisterDons = internal.RegisterDons

type DONToRegister = internal.DONToRegister

type ConfigureContractsRequest = internal.ConfigureContractsRequest
type ConfigureContractsResponse = internal.ConfigureContractsResponse

type DonCapabilities = internal.DonCapabilities

//var DeployForwarder = internal.DeployForwarder
