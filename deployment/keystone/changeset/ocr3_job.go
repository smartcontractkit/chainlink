package changeset

import (
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink/deployment"
)

// ErrOCR3JobDeprecated is returned by the OCR3 capability job builders below.
// The chainlink-ocr3-capability plugin and its job spec have been removed;
// use deployment/cre/jobs instead.
var ErrOCR3JobDeprecated = errors.New("the OCR3 capability job is deprecated and no longer supported")

type OCR3JobConfig struct {
	JobName              string
	ChainID              string
	P2PID                string
	OCR2EVMKeyBundleID   string
	TransmitterID        string
	OCR2AptosKeyBundleID string
	ContractID           string // contract ID of the ocr3 contract
	P2Pv2Bootstrappers   []string
	ExternalJobID        string
}

func (c OCR3JobConfig) Validate() error {
	return ErrOCR3JobDeprecated
}

func ResolveOCR3Job(cfg OCR3JobConfig) (string, error) {
	return "", ErrOCR3JobDeprecated
}

type OCR3JobConfigSpec struct {
	NodeID  string
	JobName string
	Spec    string
}

func BuildOCR3JobConfigSpecs(
	client deployment.NodeChainConfigsLister,
	lggr logger.Logger,
	contractID string,
	evmChainSel, aptosChainSel uint64,
	nodes []*nodev1.Node,
	btURLs []string,
	donName string,
) ([]OCR3JobConfigSpec, error) {
	return nil, ErrOCR3JobDeprecated
}

// NOTE: consider adding contract address to the hash
func ExternalJobID(donName string, evmChainSel uint64) (string, error) {
	return "", ErrOCR3JobDeprecated
}
