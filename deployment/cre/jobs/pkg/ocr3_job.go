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

type OCR3JobConfigInput struct {
	TemplateName         string        `yaml:"templateName"`
	ContractQualifier    string        `yaml:"contractQualifier"`
	ChainSelectorEVM     ChainSelector `yaml:"chainSelectorEVM"`
	ChainSelectorAptos   ChainSelector `yaml:"chainSelectorAptos"`
	BootstrapperOCR3Urls []string      `yaml:"bootstrapperOCR3Urls"`

	// Optionals: specific to the worker vault OCR3 Job spec
	DKGContractQualifier string `yaml:"dkgContractQualifier"`
}

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
	TemplateName         string

	DKGContractAddress string
}

func (c OCR3JobConfig) Validate() error {
	if c.TemplateName == "" {
		return errors.New("TemplateName is empty")
	}
	if c.JobName == "" {
		return errors.New("JobName is empty")
	}
	if c.ChainID == "" {
		return errors.New("ChainID is empty")
	}
	if c.P2PID == "" {
		return errors.New("P2PID is empty")
	}
	if c.OCR2EVMKeyBundleID == "" && c.OCR2AptosKeyBundleID == "" {
		return errors.New("OCR2EVMKeyBundleID and OCR2AptosKeyBundleID are empty, one must be set")
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

	return nil
}

func (c OCR3JobConfig) ResolveJob() (string, error) {
	templateName := c.TemplateName + ".tmpl"
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
	donName, jobName, templateName string,
	dkgContractAddress string,
) ([]OCR3JobConfigSpec, error) {
	nodesLen := len(nodes)
	if nodesLen == 0 {
		return nil, errors.New("no nodes to build OCR3 job configs")
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

	extJobID, err := ExternalJobID(donName, contractID, templateName, evmChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to get external job ID: %w", err)
	}

	jobConfigByNode := make(map[string]*OCR3JobConfig)
	for _, node := range nodeInfos {
		if node.IsBootstrap {
			lggr.Infow("Skipping bootstrap node", "nodeID", node.NodeID, "chainSelector", evmChainSel)
			continue
		}
		evmConfig, ok := node.OCRConfigForChainSelector(evmChainSel)
		if !ok {
			return nil, fmt.Errorf("no evm ocr2 config for node %s", node.NodeID)
		}

		aptosKeyBundleID := ""
		aptosConfig, ok := node.OCRConfigForChainSelector(aptosChainSel)
		if ok {
			aptosKeyBundleID = aptosConfig.KeyBundleID
		}

		jbName := "OCR3 Multichain Capability (" + node.Name + ")"
		if jobName != "" {
			jbName = jobName + " (" + node.Name + ")"
		}
		jobConfig := &OCR3JobConfig{
			JobName:              jbName,
			ChainID:              chainID,
			P2PID:                node.PeerID.String(),
			OCR2EVMKeyBundleID:   evmConfig.KeyBundleID,
			OCR2AptosKeyBundleID: aptosKeyBundleID,
			ContractID:           contractID,
			TransmitterID:        string(evmConfig.TransmitAccount),
			P2Pv2Bootstrappers:   btURLs,
			ExternalJobID:        extJobID,
			TemplateName:         templateName,
			DKGContractAddress:   dkgContractAddress,
		}

		err1 := jobConfig.Validate()
		if err1 != nil {
			return nil, fmt.Errorf("failed to validate ocr3 job config: %w", err1)
		}
		jobConfigByNode[node.NodeID] = jobConfig
	}

	specs := make([]OCR3JobConfigSpec, 0)
	for nodeID, jobConfig := range jobConfigByNode {
		spec, err1 := jobConfig.ResolveJob()
		if err1 != nil {
			return nil, fmt.Errorf("failed to resolve ocr3 job: %w", err)
		}
		specs = append(specs, OCR3JobConfigSpec{
			Spec:    spec,
			NodeID:  nodeID,
			JobName: jobConfig.JobName,
		})
	}

	return specs, nil
}

func ExternalJobID(donName, contractID, templateName string, evmChainSel uint64) (string, error) {
	in := []byte(donName + "-" + contractID + "-" + templateName + "-ocr3-capability-job-spec")
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
