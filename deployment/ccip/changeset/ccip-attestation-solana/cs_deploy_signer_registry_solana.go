package ccip_attestation

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	chainsel "github.com/smartcontractkit/chain-selectors"

	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry_solana"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"

	sol_binary "github.com/gagliardetto/binary"
	sol_rpc "github.com/gagliardetto/solana-go/rpc"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// use this changeset to deploy the base signer registry contract
var _ cldf.ChangeSet[DeployBaseSignerRegistryContractConfig] = DeployBaseSignerRegistryContractChangeset

// use this changeset to initialize the base signer registry contract and set an initial owner
var _ cldf.ChangeSet[InitalizeBaseSignerRegistryContractConfig] = InitializeBaseSignerRegistryContractChangeset

type DeployBaseSignerRegistryContractConfig struct {
	ChainSelector uint64
	Version       semver.Version
	WorkflowRun   string
	ArtifactID    string
	IsUpgrade     bool
}

type InitalizeBaseSignerRegistryContractConfig struct {
	ChainSelector uint64
}

func DeployBaseSignerRegistryContractChangeset(e cldf.Environment, c DeployBaseSignerRegistryContractConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Deploying base signer registry", "chain_selector", c.ChainSelector)
	err := c.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy base signer registry contract: %w", err)
	}
	chainSel := c.ChainSelector
	chain := e.BlockChains.SolanaChains()[chainSel]

	newAddresses := cldf.NewMemoryAddressBook()

	programFileName := deployment.BaseSignerRegistryProgramName + ".so"
	programFilePath := filepath.Join(chain.ProgramsPath, programFileName)
	if _, err := os.Stat(programFilePath); err != nil {
		if !os.IsNotExist(err) {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to check existing program artifact: %w", err)
		}
		if strings.TrimSpace(c.WorkflowRun) == "" || strings.TrimSpace(c.ArtifactID) == "" {
			return cldf.ChangesetOutput{}, fmt.Errorf("program artifact %s not found in %s and workflow run/artifact ID not provided", programFileName, chain.ProgramsPath)
		}
		if err := DownloadReleaseArtifactsFromGithubWorkflowRun(context.Background(), c.WorkflowRun, c.ArtifactID, chain.ProgramsPath); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to download release artifacts: %w", err)
		}
	}
	_, err = deployBaseSignerRegistryContract(e, chain, newAddresses, c)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy base signer registry contract: %w", err)
	}

	return cldf.ChangesetOutput{
		AddressBook: newAddresses,
	}, nil
}

func InitializeBaseSignerRegistryContractChangeset(e cldf.Environment, c InitalizeBaseSignerRegistryContractConfig) (cldf.ChangesetOutput, error) {
	e.Logger.Infow("Initializing base signer registry", "chain_selector", c.ChainSelector)
	err := c.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize base signer registry contract: %w", err)
	}
	chainSel := c.ChainSelector
	chain := e.BlockChains.SolanaChains()[chainSel]
	authority := chain.DeployerKey.PublicKey()

	configPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("config")}, signer_registry.ProgramID)
	signersPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("signers")}, signer_registry.ProgramID)
	eventAuthorityPda, _, _ := solana.FindProgramAddress([][]byte{[]byte("__event_authority")}, signer_registry.ProgramID)
	programData, err := getSolProgramData(e, chain, signer_registry.ProgramID)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ix, err := signer_registry.NewInitializeInstruction(
		authority,
		solana.SystemProgramID,
		configPda,
		signersPda,
		signer_registry.ProgramID,
		programData.Address,
		eventAuthorityPda,
		signer_registry.ProgramID,
	)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize base signer registry contract: %w", err)
	}

	if err := chain.Confirm([]solana.Instruction{ix}); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize base signer registry contract: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func deployBaseSignerRegistryContract(e cldf.Environment, chain cldf_solana.Chain, ab cldf.AddressBook, config DeployBaseSignerRegistryContractConfig,
) (solana.PublicKey, error) {
	contractType := shared.SVMSignerRegistry
	programName := deployment.BaseSignerRegistryProgramName

	programID, err := chain.DeployProgram(e.Logger, cldf_solana.ProgramInfo{
		Name:  programName,
		Bytes: deployment.SolanaProgramBytes[programName],
	}, config.IsUpgrade, true)

	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to deploy program: %w", err)
	}
	address := solana.MustPublicKeyFromBase58(programID)

	e.Logger.Infow("Deployed program", "Program", contractType, "addr", programID, "chain", chain.String())
	tv := cldf.NewTypeAndVersion(contractType, config.Version)
	err = ab.Save(chain.Selector, programID, tv)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to save address: %w", err)
	}

	return address, nil
}

func (c DeployBaseSignerRegistryContractConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}

	return nil
}

func (c InitalizeBaseSignerRegistryContractConfig) Validate(e cldf.Environment) error {
	if err := cldf.IsValidChainSelector(c.ChainSelector); err != nil {
		return fmt.Errorf("invalid chain selector: %d - %w", c.ChainSelector, err)
	}
	family, _ := chainsel.GetSelectorFamily(c.ChainSelector)
	if family != chainsel.FamilySolana {
		return fmt.Errorf("chain %d is not a solana chain", c.ChainSelector)
	}

	return nil
}

func DownloadReleaseArtifactsFromGithubWorkflowRun(
	ctx context.Context,
	run string,
	artifactID string,
	targetPath string,
) error {
	url := fmt.Sprintf(
		"https://github.com/smartcontractkit/ccip-base/actions/runs/%s/artifacts/%s",
		run,
		artifactID,
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
		// Clean the file path to prevent directory traversal
		cleanedName := filepath.Clean(file.Name)
		// Ensure the file path doesn't escape the target directory
		if strings.Contains(cleanedName, "..") {
			return fmt.Errorf("invalid file path in archive: %s", file.Name)
		}
		filePath := filepath.Join(targetPath, cleanedName)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filePath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for %s: %w", filePath, err)
		}

		fileReader, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip %s: %w", file.Name, err)
		}

		targetFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}

		// Limit the amount of data to copy to prevent decompression bombs
		const maxFileSize = 100 * 1024 * 1024 // 100MB limit per file
		limitedReader := io.LimitReader(fileReader, maxFileSize)
		n, err := io.Copy(targetFile, limitedReader)
		fileReader.Close()
		targetFile.Close()
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
		if n == maxFileSize {
			return fmt.Errorf("file %s exceeds maximum allowed size of %d bytes", filePath, maxFileSize)
		}
	}

	return nil
}

func getSolProgramData(e cldf.Environment, chain cldf_solana.Chain, programID solana.PublicKey) (struct {
	DataType uint32
	Address  solana.PublicKey
}, error) {
	var programData struct {
		DataType uint32
		Address  solana.PublicKey
	}
	data, err := chain.Client.GetAccountInfoWithOpts(e.GetContext(), programID, &sol_rpc.GetAccountInfoOpts{
		Commitment: sol_rpc.CommitmentConfirmed,
	})
	if err != nil {
		return programData, fmt.Errorf("failed to deploy program: %w", err)
	}

	err = sol_binary.UnmarshalBorsh(&programData, data.Bytes())
	if err != nil {
		return programData, fmt.Errorf("failed to unmarshal program data: %w", err)
	}
	return programData, nil
}
