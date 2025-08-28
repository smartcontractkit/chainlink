package ccipbase

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	signerRegistry "github.com/smartcontractkit/ccip-base/chains/solana/go_bindings"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
)

// Temporary reference to prevent import removal
var _ = signerRegistry.SIGNERS_SEED

type DeploySignerRegistryContractConfig struct {
	HomeChainSelector uint64 // eth mainnet/testnet
	ChainSelector     uint64
}

func DeploySignerRegistryContractChangeset(e cldf.Environment, c DeploySignerRegistryContractConfig) (cldf.ChangesetOutput, error) {
	existingState, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load existing onchain state: %w", err)
	}
	if err := c.Validate(e, existingState); err != nil {
		// TODO
		return cldf.ChangesetOutput{}, fmt.Errorf("Invalid DeploySignerRegistryContractConfig: %w", err)
	}

}

func DeploySignerRegistryContract(e cldf.Environment, chain cldf_solana.Chain, ab cldf.AddressBook, config DeploySignerRegistryContractConfig,
	version semver.Version,
) (solana.PublicKey, error) {
	contractType := shared.BaseSignerRegistry
	programName := deployment.BaseSignerRegistryProgramName

	programID, err := chain.DeployProgram(e.Logger, cldf_solana.ProgramInfo{
		Name:  programName,
		Bytes: deployment.SolanaProgramBytes[programName],
	}, false, true)

	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to deploy program: %w", err)
	}
	address := solana.MustPublicKeyFromBase58(programID)

	e.Logger.Infow("Deployed program", "Program", contractType, "addr", programID, "chain", chain.String())
	tv := cldf.NewTypeAndVersion(contractType, version)
	err = ab.Save(chain.Selector, programID, tv)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to save address: %w", err)
	}

	return address, nil

}

func (c DeploySignerRegistryContractConfig) Validate(e cldf.Environment, existingState stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(c.HomeChainSelector); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSelector, err)
	}
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}
	if _, exists := existingState.SupportedChains()[c.ChainSelector]; !exists {
		return fmt.Errorf("chain %d not supported", c.ChainSelector)
	}

	return nil
}

// DownloadReleaseArtifactsFromGithub makes a GET request to the URL, retrieves the release artifacts zip file,
// and extracts it to a destinationDir
func DownloadReleaseArtifactsFromGithub(
	ctx context.Context,
	e cldf.Environment,
	name string,
	tag string,
	destinationDir string,
) error {
	url := fmt.Sprintf(
		"https://github.com/smartcontractkit/ccip-base/releases/download/%s/%s",
		tag,
		name,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download release asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download release asset: HTTP %d", resp.StatusCode)
	}

	// Read the entire zip file into memory
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Extract each file from the zip archive
	for _, file := range zipReader.File {
		targetPath := filepath.Join(destinationDir, file.Name)

		// Check if it's a directory
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
		}

		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip %s: %w", file.Name, err)
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			targetFile.Close()
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}
		targetFile.Close()
		
		e.Logger.Infow("Downloaded file", "filename", file.Name, "targetPath", targetPath)
	}

	return nil
}
