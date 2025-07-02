package utils

import (
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
	// Create MCMS inspector
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
