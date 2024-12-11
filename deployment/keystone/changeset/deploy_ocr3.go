package changeset

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
)

var _ deployment.ChangeSet[uint64] = DeployOCR3

func DeployOCR3(env deployment.Environment, registryChainSel uint64) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	// ocr3 only deployed on registry chain
	c, ok := env.Chains[registryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain not found in environment")
	}
	ocr3Resp, err := kslib.DeployOCR3(c, ab)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy OCR3Capability: %w", err)
	}
	lggr.Infof("Deployed %s chain selector %d addr %s", ocr3Resp.Tv.String(), c.Selector, ocr3Resp.Address.String())
	return deployment.ChangesetOutput{AddressBook: ab}, nil
}

var _ deployment.ChangeSet[ConfigureOCR3Config] = ConfigureOCR3Contract

type ConfigureOCR3Config struct {
	ChainSel             uint64
	NodeIDs              []string
	OCR3Config           *kslib.OracleConfig
	DryRun               bool
	WriteGeneratedConfig io.Writer // if not nil, write the generated config to this writer as JSON [OCR2OracleConfig]

	UseMCMS bool
}

func ConfigureOCR3Contract(env deployment.Environment, cfg ConfigureOCR3Config) (deployment.ChangesetOutput, error) {
	resp, err := kslib.ConfigureOCR3ContractFromJD(&env, kslib.ConfigureOCR3Config{
		ChainSel:   cfg.ChainSel,
		NodeIDs:    cfg.NodeIDs,
		OCR3Config: cfg.OCR3Config,
		DryRun:     cfg.DryRun,
		UseMCMS:    cfg.UseMCMS,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to configure OCR3Capability: %w", err)
	}
	if w := cfg.WriteGeneratedConfig; w != nil {
		b, err := json.MarshalIndent(&resp.OCR2OracleConfig, "", "  ")
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to marshal response output: %w", err)
		}
		env.Logger.Infof("Generated OCR3 config: %s", string(b))
		n, err := w.Write(b)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to write response output: %w", err)
		}
		if n != len(b) {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to write all bytes")
		}
	}
	// does not create any new addresses
	var proposals []timelock.MCMSWithTimelockProposal
	if cfg.UseMCMS {
		proposals = append(proposals, *resp.Proposal)
	}
	return deployment.ChangesetOutput{
		Proposals: proposals,
	}, nil
}
