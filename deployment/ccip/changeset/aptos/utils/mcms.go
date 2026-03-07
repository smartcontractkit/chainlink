package utils

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/compile"
	curse_mcms "github.com/smartcontractkit/chainlink-aptos/bindings/curse_mcms"
	mcmsbind "github.com/smartcontractkit/chainlink-aptos/bindings/mcms"
	module_mcms "github.com/smartcontractkit/chainlink-aptos/bindings/mcms/mcms"
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

// CurseMCMSInspector implements sdk.Inspector for the CurseMCMS contract.
// The standard aptosmcms.Inspector uses mcms.Bind which targets the "mcms" module,
// but CurseMCMS exposes the same view functions under the "curse_mcms" module.
var _ mcmssdk.Inspector = &CurseMCMSInspector{}

type CurseMCMSInspector struct {
	configTransformer aptosmcms.ConfigTransformer
	client            aptos.AptosRpcClient
	role              aptosmcms.TimelockRole
}

func NewCurseMCMSInspector(client aptos.AptosRpcClient, role aptosmcms.TimelockRole) *CurseMCMSInspector {
	return &CurseMCMSInspector{client: client, role: role}
}

func (i CurseMCMSInspector) GetConfig(ctx context.Context, addr string) (*mcmstypes.Config, error) {
	address, err := parseAptosAddress(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CurseMCMS address: %w", err)
	}
	binding := curse_mcms.Bind(address, i.client)
	cfg, err := binding.CurseMCMS().GetConfig(nil, i.role.Byte())
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}
	mcmsCfg := module_mcms.Config{
		GroupQuorums: cfg.GroupQuorums,
		GroupParents: cfg.GroupParents,
	}
	for _, s := range cfg.Signers {
		mcmsCfg.Signers = append(mcmsCfg.Signers, module_mcms.Signer{Addr: s.Addr, Index: s.Index, Group: s.Group})
	}
	return i.configTransformer.ToConfig(mcmsCfg)
}

func (i CurseMCMSInspector) GetOpCount(ctx context.Context, addr string) (uint64, error) {
	address, err := parseAptosAddress(addr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CurseMCMS address: %w", err)
	}
	binding := curse_mcms.Bind(address, i.client)
	opCount, err := binding.CurseMCMS().GetOpCount(nil, i.role.Byte())
	if err != nil {
		return 0, fmt.Errorf("get op count: %w", err)
	}
	return opCount, nil
}

func (i CurseMCMSInspector) GetRoot(ctx context.Context, addr string) (common.Hash, uint32, error) {
	address, err := parseAptosAddress(addr)
	if err != nil {
		return common.Hash{}, 0, fmt.Errorf("failed to parse CurseMCMS address: %w", err)
	}
	binding := curse_mcms.Bind(address, i.client)
	root, validUntil, err := binding.CurseMCMS().GetRoot(nil, i.role.Byte())
	if err != nil {
		return common.Hash{}, 0, fmt.Errorf("get root: %w", err)
	}

	if validUntil > math.MaxUint32 {
		return common.Hash{}, 0, fmt.Errorf("validUntil %d overflows uint32", validUntil)
	}
	return common.BytesToHash(root), uint32(validUntil), nil
}

func (i CurseMCMSInspector) GetRootMetadata(ctx context.Context, addr string) (mcmstypes.ChainMetadata, error) {
	address, err := parseAptosAddress(addr)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to parse CurseMCMS address: %w", err)
	}
	binding := curse_mcms.Bind(address, i.client)
	rootMetadata, err := binding.CurseMCMS().GetRootMetadata(nil, i.role.Byte())
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("get root metadata: %w", err)
	}
	return mcmstypes.ChainMetadata{
		StartingOpCount: rootMetadata.PreOpCount,
		MCMAddress:      rootMetadata.Multisig.StringLong(),
	}, nil
}

// GenerateCurseMCMSProposal creates a TimelockProposal targeting the CurseMCMS
// contract. It uses CurseMCMSInspector instead of the standard MCMS inspector
// because CurseMCMS exposes view functions under "curse_mcms" not "mcms".
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
	inspector := NewCurseMCMSInspector(env.BlockChains.AptosChains()[chainSel].Client, role)

	return proposalutils.BuildProposalFromBatchesV2(
		env,
		map[uint64]string{chainSel: curseMCMSAddress.StringLong()},
		map[uint64]string{chainSel: curseMCMSAddress.StringLong()},
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

func parseAptosAddress(addr string) (aptos.AccountAddress, error) {
	var address aptos.AccountAddress
	if err := address.ParseStringRelaxed(addr); err != nil {
		return aptos.AccountAddress{}, err
	}
	return address, nil
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
