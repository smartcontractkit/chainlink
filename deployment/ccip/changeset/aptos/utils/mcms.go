package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/compile"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

const MCMSProposalVersion = "v1"

func GenerateProposal(
	env cldf.Environment,
	mcmsAddress aptos.AccountAddress,
	chainSel uint64,
	operations []mcmstypes.BatchOperation,
	description string,
	mcmsCfg proposalutils.TimelockConfig,
) (*mcms.TimelockProposal, error) {
	// Get role from action
	role, err := proposalutils.GetAptosRoleFromAction(mcmsCfg.MCMSAction)
	if err != nil {
		return nil, fmt.Errorf("failed to get role from action: %w", err)
	}
	inspector := aptosmcms.NewInspector(env.BlockChains.AptosChains()[chainSel].Client, role)

	return proposalutils.BuildProposalFromBatchesV2(
		env,
		map[uint64]string{chainSel: mcmsAddress.StringLong()},
		map[uint64]string{chainSel: mcmsAddress.StringLong()},
		map[uint64]mcmssdk.Inspector{chainSel: inspector},
		operations,
		description,
		mcmsCfg,
	)
}

// GenerateCurseMCMSProposal creates a TimelockProposal targeting the CurseMCMS
// contract. It uses NewInspectorWithMCMSType(MCMSTypeCurse) so the inspector
// reads config from the "curse_mcms" module instead of "mcms". It also sets
// MCMSType=MCMSTypeCurse in the chain metadata so downstream tools
// (mcms-tools set-signers, executor, etc.) use the correct binding.
func GenerateCurseMCMSProposal(
	env cldf.Environment,
	curseMCMSAddress aptos.AccountAddress,
	chainSel uint64,
	operations []mcmstypes.BatchOperation,
	description string,
	mcmsCfg proposalutils.TimelockConfig,
) (*mcms.TimelockProposal, error) {
	role, err := proposalutils.GetAptosRoleFromAction(mcmsCfg.MCMSAction)
	if err != nil {
		return nil, fmt.Errorf("failed to get role from action: %w", err)
	}
	inspector := aptosmcms.NewInspectorWithMCMSType(env.BlockChains.AptosChains()[chainSel].Client, role, aptosmcms.MCMSTypeCurse)

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		env,
		map[uint64]string{chainSel: curseMCMSAddress.StringLong()},
		map[uint64]string{chainSel: curseMCMSAddress.StringLong()},
		map[uint64]mcmssdk.Inspector{chainSel: inspector},
		operations,
		description,
		mcmsCfg,
	)
	if err != nil {
		return nil, err
	}

	if err := markChainMetadataAsCurseMCMS(proposal, mcmstypes.ChainSelector(chainSel)); err != nil {
		return nil, fmt.Errorf("failed to set MCMSType in chain metadata: %w", err)
	}

	return proposal, nil
}

// markChainMetadataAsCurseMCMS sets MCMSType=MCMSTypeCurse in the Aptos
// AdditionalFieldsMetadata for the given chain selector.
func markChainMetadataAsCurseMCMS(proposal *mcms.TimelockProposal, cs mcmstypes.ChainSelector) error {
	meta, ok := proposal.ChainMetadata[cs]
	if !ok {
		return fmt.Errorf("chain selector %d not found in proposal chain metadata", cs)
	}
	var afm aptosmcms.AdditionalFieldsMetadata
	if len(meta.AdditionalFields) > 0 {
		if err := json.Unmarshal(meta.AdditionalFields, &afm); err != nil {
			return fmt.Errorf("unmarshal additional fields metadata: %w", err)
		}
	}
	afm.MCMSType = aptosmcms.MCMSTypeCurse
	b, err := json.Marshal(afm)
	if err != nil {
		return fmt.Errorf("marshal additional fields metadata: %w", err)
	}
	meta.AdditionalFields = b
	proposal.ChainMetadata[cs] = meta
	return nil
}

// ToBatchOperations converts Operations into BatchOperations with a single transaction each
func ToBatchOperations(ops []mcmstypes.Operation) []mcmstypes.BatchOperation {
	var batchOps []mcmstypes.BatchOperation
	for _, op := range ops {
		batchOps = append(batchOps, mcmstypes.BatchOperation{
			ChainSelector: op.ChainSelector,
			Transactions:  []mcmstypes.Transaction{op.Transaction},
		})
	}
	return batchOps
}

// IsMCMSStagingAreaClean checks if the MCMS staging area is clean
func IsMCMSStagingAreaClean(client aptos.AptosRpcClient, aptosMCMSObjAddr aptos.AccountAddress) (bool, error) {
	resources, err := client.AccountResources(aptosMCMSObjAddr)
	if err != nil {
		return false, err
	}
	for _, resource := range resources {
		if strings.Contains(resource.Type, "StagingArea") {
			return false, nil
		}
	}
	return true, nil
}

// CreateChunksAndStage creates chunks from the compiled packages and build MCMS operations to stages them within the MCMS contract
func CreateChunksAndStage(
	payload compile.CompiledPackage,
	mcmsContract mcmsbind.MCMS,
	chainSel uint64,
	seed string,
	codeObjectAddress *aptos.AccountAddress,
) ([]mcmstypes.Operation, error) {
	mcmsAddress := mcmsContract.Address()
	// Validate seed XOR codeObjectAddress, one and only one must be provided
	if (seed != "") == (codeObjectAddress != nil) {
		return nil, errors.New("either provide seed to publishToObject or objectAddress to upgradeObjectCode")
	}

	var operations []mcmstypes.Operation

	// Create chunks
	chunks, err := bind.CreateChunks(payload, bind.ChunkSizeInBytes)
	if err != nil {
		return operations, fmt.Errorf("failed to create chunks: %w", err)
	}

	// Stage chunks with mcms_deployer module and execute with the last one
	for i, chunk := range chunks {
		var (
			moduleInfo bind.ModuleInformation
			function   string
			args       [][]byte
			err        error
		)

		// First chunks get staged, the last one gets published or upgraded
		switch {
		case i != len(chunks)-1:
			moduleInfo, function, _, args, err = mcmsContract.MCMSDeployer().Encoder().StageCodeChunk(
				chunk.Metadata,
				chunk.CodeIndices,
				chunk.Chunks,
			)
		case seed != "":
			moduleInfo, function, _, args, err = mcmsContract.MCMSDeployer().Encoder().StageCodeChunkAndPublishToObject(
				chunk.Metadata,
				chunk.CodeIndices,
				chunk.Chunks,
				[]byte(seed),
			)
		default:
			moduleInfo, function, _, args, err = mcmsContract.MCMSDeployer().Encoder().StageCodeChunkAndUpgradeObjectCode(
				chunk.Metadata,
				chunk.CodeIndices,
				chunk.Chunks,
				*codeObjectAddress,
			)
		}
		if err != nil {
			return operations, fmt.Errorf("failed to encode chunk %d: %w", i, err)
		}

		tx, err := GenerateMCMSTx(mcmsAddress, moduleInfo, function, args)
		if err != nil {
			return operations, fmt.Errorf("failed to create transaction: %w", err)
		}

		operations = append(operations, mcmstypes.Operation{
			ChainSelector: mcmstypes.ChainSelector(chainSel),
			Transaction:   tx,
		})
	}

	return operations, nil
}

// GenerateMCMSTx is a helper function that generates a MCMS txs for the given parameters
func GenerateMCMSTx(toAddress aptos.AccountAddress, moduleInfo bind.ModuleInformation, function string, args [][]byte) (mcmstypes.Transaction, error) {
	return aptosmcms.NewTransaction(
		moduleInfo.PackageName,
		moduleInfo.ModuleName,
		function,
		toAddress,
		aptosmcms.ArgsToData(args),
		"",
		nil,
	)
}
