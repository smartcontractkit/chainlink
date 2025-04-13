package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/compile"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	"github.com/smartcontractkit/mcms/types"
)

const (
	ValidUntilHours     = 72
	MCMSProposalVersion = "v1"
)

func GenerateProposal(
	client aptos.AptosRpcClient,
	mcmsAddress aptos.AccountAddress,
	chainSel uint64,
	operations []types.BatchOperation,
	description string,
	role aptosmcms.TimelockRole,
) (*mcms.TimelockProposal, error) {
	// Create MCMS inspector
	inspector := aptosmcms.NewInspector(client, role)
	startingOpCount, err := inspector.GetOpCount(context.Background(), mcmsAddress.StringLong())
	if err != nil {
		return nil, fmt.Errorf("failed to get starting op count: %w", err)
	}
	opCount := startingOpCount

	action := types.TimelockActionSchedule
	if role == aptosmcms.TimelockRoleBypasser {
		action = types.TimelockActionBypass
	}
	jsonRole, _ := json.Marshal(aptosmcms.AdditionalFieldsMetadata{Role: role})

	// Create proposal builder
	validUntil := time.Now().Add(time.Hour * ValidUntilHours).Unix()
	proposalBuilder := mcms.NewTimelockProposalBuilder().
		SetVersion(MCMSProposalVersion).
		SetValidUntil(uint32(validUntil)).
		SetDescription(description).
		AddTimelockAddress(types.ChainSelector(chainSel), mcmsAddress.StringLong()).
		SetOverridePreviousRoot(true).
		AddChainMetadata(
			types.ChainSelector(chainSel),
			types.ChainMetadata{
				StartingOpCount:  opCount,
				MCMAddress:       mcmsAddress.StringLong(),
				AdditionalFields: jsonRole,
			},
		).
		SetAction(action).
		SetDelay(types.NewDuration(time.Second)) // TODO: set propper delay

	// Add operations and build
	for _, op := range operations {
		proposalBuilder.AddOperation(op)
	}
	proposal, err := proposalBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build proposal: %w", err)
	}

	return proposal, nil
}

func ToBatchOperations(ops []types.Operation) []types.BatchOperation {
	batchOps := []types.BatchOperation{}
	for _, op := range ops {
		batchOps = append(batchOps, types.BatchOperation{
			ChainSelector: op.ChainSelector,
			Transactions:  []types.Transaction{op.Transaction},
		})
	}
	return batchOps
}

// CreateChunksAndStage creates chunks from the compiled packages and build MCMS operations to stages them within the MCMS contract
func CreateChunksAndStage(
	payload compile.CompiledPackage,
	mcmsContract mcmsbind.MCMS,
	chainSel uint64,
	seed string,
	codeObjectAddress *aptos.AccountAddress,
) ([]types.Operation, error) {
	mcmsAddress := mcmsContract.Address()
	// Validate seed XOR codeObjectAddress, one and only one must be provided
	if (seed != "") == (codeObjectAddress != nil) {
		return nil, fmt.Errorf("either provide seed to publishToObject or objectAddress to upgradeObjectCode")
	}

	var operations []types.Operation

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
		if i != len(chunks)-1 {
			moduleInfo, function, _, args, err = mcmsContract.MCMSDeployer().Encoder().StageCodeChunk(
				chunk.Metadata,
				chunk.CodeIndices,
				chunk.Chunks,
			)
		} else if seed != "" {
			moduleInfo, function, _, args, err = mcmsContract.MCMSDeployer().Encoder().StageCodeChunkAndPublishToObject(
				chunk.Metadata,
				chunk.CodeIndices,
				chunk.Chunks,
				[]byte(seed),
			)
		} else {
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
		additionalFields := aptosmcms.AdditionalFields{
			PackageName: moduleInfo.PackageName,
			ModuleName:  moduleInfo.ModuleName,
			Function:    function,
		}
		afBytes, err := json.Marshal(additionalFields)
		if err != nil {
			return operations, fmt.Errorf("failed to marshal additional fields: %w", err)
		}
		operations = append(operations, types.Operation{
			ChainSelector: types.ChainSelector(chainSel),
			Transaction: types.Transaction{
				To:               mcmsAddress.StringLong(),
				Data:             aptosmcms.ArgsToData(args),
				AdditionalFields: afBytes,
			},
		})
	}

	return operations, nil
}

// GenerateMCMSTx is a helper function that generates a MCMS txs for the given parameters
func GenerateMCMSTx(toAddress aptos.AccountAddress, moduleInfo bind.ModuleInformation, function string, args [][]byte) (types.Transaction, error) {
	additionalFields := aptosmcms.AdditionalFields{
		PackageName: moduleInfo.PackageName,
		ModuleName:  moduleInfo.ModuleName,
		Function:    function,
	}
	afBytes, err := json.Marshal(additionalFields)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to marshal additional fields: %w", err)
	}
	return types.Transaction{
		To:               toAddress.StringLong(),
		Data:             aptosmcms.ArgsToData(args),
		AdditionalFields: afBytes,
	}, nil
}
