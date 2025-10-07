package ccip_attestation

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	chainsel "github.com/smartcontractkit/chain-selectors"

	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry_solana"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	cs_solana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
)

// use this changeset to upload the IDL for the attestation program
var _ cldf.ChangeSet[BaseIDLConfig] = BaseUploadIDLChangeset

// use this changeset to set the authority for the IDL of the attestation program (timelock)
var _ cldf.ChangeSet[BaseIDLConfig] = BaseSetAuthorityIDLChangeset

const IdlIxTag uint64 = 0x0a69e9a778bcf440

type BaseIDLConfig struct {
	ChainSelector uint64
	WorkflowRun   string
	ArtifactID    string
}

// resolve artifacts based on workflow run and write anchor.toml file to simulate anchor workspace
func repoSetup(e cldf.Environment, chain cldf_solana.Chain, run string, artifactID string) error {
	programName := deployment.BaseSignerRegistryProgramName
	idlFileName := programName + ".json"
	idlFilePath := filepath.Join(chain.ProgramsPath, idlFileName)
	if _, err := os.Stat(idlFilePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error checking existing IDL artifact: %w", err)
		}
		if strings.TrimSpace(run) == "" || strings.TrimSpace(artifactID) == "" {
			return fmt.Errorf("IDL artifact %s not found in %s and workflow run/artifact ID not provided", idlFileName, chain.ProgramsPath)
		}
		e.Logger.Debug("Downloading artifacts from workflow run...")
		if err := DownloadReleaseArtifactsFromGithubWorkflowRun(context.Background(), run, artifactID, chain.ProgramsPath); err != nil {
			return fmt.Errorf("error downloading program artifacts: %w", err)
		}
	}

	// get anchor version
	output, err := cs_solana.RunCommand("anchor", []string{"--version"}, ".")
	if err != nil {
		return errors.New("anchor-cli not installed in path")
	}
	e.Logger.Debugw("Anchor version command output", "output", output)
	anchorVersion, err := cs_solana.ParseAnchorVersion(output)
	if err != nil {
		return fmt.Errorf("error parsing anchor version: %w", err)
	}
	// create Anchor.toml
	// this creates anchor workspace with cluster and wallet configured
	if err := cs_solana.WriteAnchorToml(e, filepath.Join(chain.ProgramsPath, "Anchor.toml"), anchorVersion, chain.URL, chain.KeypairPath); err != nil {
		return fmt.Errorf("error writing Anchor.toml: %w", err)
	}

	return nil
}

func (c BaseIDLConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	return nil
}

func BaseUploadIDLChangeset(e cldf.Environment, c BaseIDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]
	if err := repoSetup(e, chain, c.WorkflowRun, c.ArtifactID); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error setting up repo: %w", err)
	}

	idlFile := filepath.Join(chain.ProgramsPath, deployment.BaseSignerRegistryProgramName+".json")
	if _, err := os.Stat(idlFile); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("idl file not found: %w", err)
	}

	e.Logger.Infow("Uploading IDL", "programName", deployment.BaseSignerRegistryProgramName)
	args := []string{"idl", "init", "--filepath", idlFile, signer_registry.ProgramID.String()}
	e.Logger.Info(args)
	output, err := cs_solana.RunCommand("anchor", args, chain.ProgramsPath)
	e.Logger.Debugw("IDL init output", "output", output)
	if err != nil {
		e.Logger.Debugw("IDL init error", "error", err)
		return cldf.ChangesetOutput{}, fmt.Errorf("error uploading idl: %w", err)
	}
	e.Logger.Infow("IDL uploaded", "programName", deployment.BaseSignerRegistryProgramName)
	return cldf.ChangesetOutput{}, nil
}

func BaseSetAuthorityIDLChangeset(e cldf.Environment, c BaseIDLConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(e); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error validating idl config: %w", err)
	}
	chain := e.BlockChains.SolanaChains()[c.ChainSelector]

	timelockSignerPDA, err := cs_solana.FetchTimelockSigner(e, c.ChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("error loading timelockSignerPDA: %w", err)
	}

	err = cs_solana.SetAuthorityIDLByCLI(e, timelockSignerPDA.String(), chain.ProgramsPath, signer_registry.ProgramID.String(), deployment.BaseSignerRegistryProgramName, "")
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{}, nil
}
