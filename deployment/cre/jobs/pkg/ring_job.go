package pkg

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"text/template"

	"github.com/google/uuid"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg/templates"
)

// RingJobConfigInput contains the input configuration for a Ring job proposal.
// This is used when proposing Ring jobs via the job distributor.
type RingJobConfigInput struct {
	// ContractQualifier is the qualifier for the OCR3 capability contract.
	ContractQualifier string `yaml:"contractQualifier"`
	// ChainSelectorEVM is the EVM chain selector for the target chain.
	ChainSelectorEVM ChainSelector `yaml:"chainSelectorEVM"`
	// ShardConfigAddr is the address of the ShardConfig contract on-chain.
	ShardConfigAddr string `yaml:"shardConfigAddr"`
	// BootstrapperRingUrls are the P2P URLs of the bootstrapper nodes for Ring.
	BootstrapperRingUrls []string `yaml:"bootstrapperRingUrls"`
}

// RingJobConfig contains the configuration for rendering a Ring job spec template.
// This is used internally to generate the TOML job specification.
type RingJobConfig struct {
	// JobName is the human-readable name for the job.
	JobName string
	// ChainID is the numeric chain ID for the target EVM chain.
	ChainID string
	// P2PID is the node's P2P peer ID.
	P2PID string
	// OCR2EVMKeyBundleID is the OCR2 key bundle ID for EVM signing.
	OCR2EVMKeyBundleID string
	// TransmitterID is the address used for transmitting reports.
	TransmitterID string
	// ContractID is the OCR3 capability contract address.
	ContractID string
	// P2Pv2Bootstrappers are the bootstrapper URLs for P2P networking.
	P2Pv2Bootstrappers []string
	// ExternalJobID is the deterministic UUID for the job (generated from DON name and contract).
	ExternalJobID string
	// ShardConfigAddr is the address of the ShardConfig contract.
	ShardConfigAddr string
}

// Validate validates the Ring job configuration.
func (c RingJobConfig) Validate() error {
	if c.JobName == "" {
		return errors.New("JobName is empty")
	}
	if c.ChainID == "" {
		return errors.New("ChainID is empty")
	}
	if c.P2PID == "" {
		return errors.New("P2PID is empty")
	}
	if c.OCR2EVMKeyBundleID == "" {
		return errors.New("OCR2EVMKeyBundleID is empty")
	}
	if c.TransmitterID == "" {
		return errors.New("TransmitterID is empty")
	}
	if c.ContractID == "" {
		return errors.New("ContractID is empty")
	}
	if len(c.P2Pv2Bootstrappers) == 0 {
		return errors.New("P2Pv2Bootstrappers is empty")
	}
	if c.ShardConfigAddr == "" {
		return errors.New("ShardConfigAddr is empty")
	}

	return nil
}

// ResolveJob resolves the Ring job spec from the template.
func (c RingJobConfig) ResolveJob() (string, error) {
	templateName := "ring.tmpl"
	t, err := template.New("s").ParseFS(templates.FS, templateName)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", templateName, err)
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, templateName, c)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}

// RingJobConfigSpec contains the resolved job spec for a specific node.
// This is the output of BuildRingJobConfigSpecs and contains the rendered TOML spec.
type RingJobConfigSpec struct {
	// NodeID is the unique identifier of the node this spec is for.
	NodeID string
	// JobName is the human-readable name assigned to this job.
	JobName string
	// Spec is the fully rendered TOML job specification.
	Spec string
}

// BuildRingJobConfigSpecs builds Ring job specs for all plugin nodes.
func BuildRingJobConfigSpecs(
	client deployment.NodeChainConfigsLister,
	lggr logger.Logger,
	contractID string,
	evmChainSel uint64,
	nodes []*nodev1.Node,
	btURLs []string,
	donName, jobName string,
	shardConfigAddr string,
) ([]RingJobConfigSpec, error) {
	nodesLen := len(nodes)
	if nodesLen == 0 {
		return nil, errors.New("no nodes to build Ring job configs")
	}

	nodeIDs := make([]string, 0, nodesLen)
	for _, node := range nodes {
		nodeIDs = append(nodeIDs, node.Id)
	}

	nodeInfos, err := deployment.NodeInfo(nodeIDs, client)
	if err != nil {
		return nil, fmt.Errorf("failed to get node info: %w", err)
	}

	chainID, err := chainsel.GetChainIDFromSelector(evmChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID from selector: %w", err)
	}

	extJobID, err := RingExternalJobID(donName, contractID, evmChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to get external job ID: %w", err)
	}

	jobConfigByNode := make(map[string]*RingJobConfig)
	for _, node := range nodeInfos {
		if node.IsBootstrap {
			lggr.Infow("Skipping bootstrap node for Ring job", "nodeID", node.NodeID, "chainSelector", evmChainSel)
			continue
		}
		evmConfig, ok := node.OCRConfigForChainSelector(evmChainSel)
		if !ok {
			return nil, fmt.Errorf("no evm ocr2 config for node %s", node.NodeID)
		}

		jbName := "Ring Capability (" + node.Name + ")"
		if jobName != "" {
			jbName = jobName + " (" + node.Name + ")"
		}
		jobConfig := &RingJobConfig{
			JobName:            jbName,
			ChainID:            chainID,
			P2PID:              node.PeerID.String(),
			OCR2EVMKeyBundleID: evmConfig.KeyBundleID,
			ContractID:         contractID,
			TransmitterID:      string(evmConfig.TransmitAccount),
			P2Pv2Bootstrappers: btURLs,
			ExternalJobID:      extJobID,
			ShardConfigAddr:    shardConfigAddr,
		}

		err1 := jobConfig.Validate()
		if err1 != nil {
			return nil, fmt.Errorf("failed to validate ring job config: %w", err1)
		}
		jobConfigByNode[node.NodeID] = jobConfig
	}

	specs := make([]RingJobConfigSpec, 0)
	for nodeID, jobConfig := range jobConfigByNode {
		spec, err1 := jobConfig.ResolveJob()
		if err1 != nil {
			return nil, fmt.Errorf("failed to resolve ring job: %w", err1)
		}
		specs = append(specs, RingJobConfigSpec{
			Spec:    spec,
			NodeID:  nodeID,
			JobName: jobConfig.JobName,
		})
	}

	return specs, nil
}

// RingExternalJobID generates a deterministic external job ID for Ring jobs.
func RingExternalJobID(donName, contractID string, evmChainSel uint64) (string, error) {
	in := []byte(donName + "-" + contractID + "-ring-job-spec")
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, evmChainSel)
	in = append(in, b...)
	sha256Hash := sha256.New()
	sha256Hash.Write(in)
	in = sha256Hash.Sum(nil)[:16]
	// tag as valid UUID v4 https://github.com/google/uuid/blob/0f11ee6918f41a04c201eceeadf612a377bc7fbc/version4.go#L53-L54
	in[6] = (in[6] & 0x0f) | 0x40 // Version 4
	in[8] = (in[8] & 0x3f) | 0x80 // Variant is 10

	id, err := uuid.FromBytes(in)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
