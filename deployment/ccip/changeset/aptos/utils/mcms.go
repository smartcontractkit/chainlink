package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/smartcontractkit/mcms"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/compile"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

const MCMSProposalVersion = "v1"

func GenerateProposal(
	client aptos.AptosRpcClient,
	mcmsAddress aptos.AccountAddress,
	chainSel uint64,
	operations []mcmstypes.BatchOperation,
	description string,
	mcmsCfg proposalutils.TimelockConfig,
) (*mcms.TimelockProposal, error) {
	// Get role from action
	role, err := roleFromAction(mcmsCfg.MCMSAction)
	if err != nil {
		return nil, fmt.Errorf("failed to get role from action: %w", err)
	}
	jsonRole, _ := json.Marshal(aptosmcms.AdditionalFieldsMetadata{Role: role})
	var action = mcmsCfg.MCMSAction
	if action == "" {
		action = mcmstypes.TimelockActionSchedule
	}
	// Create MCMS inspector
	inspector := aptosmcms.NewInspector(client, role)
	startingOpCount, err := inspector.GetOpCount(context.Background(), mcmsAddress.StringLong())
	if err != nil {
		return nil, fmt.Errorf("failed to get starting op count: %w", err)
	}
	opCount := startingOpCount

	// Create proposal builder
	validUntil := time.Now().Unix() + int64(proposalutils.DefaultValidUntil.Seconds())
	if validUntil < 0 || validUntil > math.MaxUint32 {
		return nil, fmt.Errorf("validUntil value out of range for uint32: %d", validUntil)
	}

	proposalBuilder := mcms.NewTimelockProposalBuilder().
		SetVersion(MCMSProposalVersion).
		SetValidUntil(uint32(validUntil)).
		SetDescription(description).
		AddTimelockAddress(mcmstypes.ChainSelector(chainSel), mcmsAddress.StringLong()).
		SetOverridePreviousRoot(mcmsCfg.OverrideRoot).
		AddChainMetadata(
			mcmstypes.ChainSelector(chainSel),
			mcmstypes.ChainMetadata{
				StartingOpCount:  opCount,
				MCMAddress:       mcmsAddress.StringLong(),
				AdditionalFields: jsonRole,
			},
		).
		SetAction(action).
		SetDelay(mcmstypes.NewDuration(mcmsCfg.MinDelay))

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

func roleFromAction(action mcmstypes.TimelockAction) (aptosmcms.TimelockRole, error) {
	switch action {
	case mcmstypes.TimelockActionSchedule:
		return aptosmcms.TimelockRoleProposer, nil
	case mcmstypes.TimelockActionBypass:
		return aptosmcms.TimelockRoleBypasser, nil
	case mcmstypes.TimelockActionCancel:
		return aptosmcms.TimelockRoleCanceller, nil
	case "":
		return aptosmcms.TimelockRoleProposer, nil
	default:
		return aptosmcms.TimelockRoleProposer, fmt.Errorf("invalid action: %s", action)
	}
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
