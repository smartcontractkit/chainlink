package changeset

import (
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

//TODO: KS-673 refactor internal package to reduce and remove the duplication

// OracleConfig is the configuration for an oracle
type OracleConfig = ocr3.OracleConfig

// OCR3OnchainConfig is the onchain configuration of an OCR2 contract
type OCR3OnchainConfig = ocr3.OCR2OracleConfig

// NodeKeys is a set of public keys for a node
type NodeKeys = ocr3.NodeKeys

// TopLevelConfigSource is the top level configuration source
type TopLevelConfigSource = ocr3.TopLevelConfigSource

// GenerateOCR3Config generates an OCR3 config
var GenerateOCR3Config = ocr3.GenerateOCR3Config

// NOP is a node operator profile, required to register a node with the capabilities registry
type NOP = internal.NOP

// DonCapabilities is a set of capabilities hosted by a set of node operators
// in is in a convenient form to handle the CLO representation of the nop data
type DonCapabilities = internal.DonCapabilities

type DONCapabilityWithConfig = internal.DONCapabilityWithConfig

type DeployRequest = internal.DeployRequest
type DeployResponse = internal.DeployResponse
