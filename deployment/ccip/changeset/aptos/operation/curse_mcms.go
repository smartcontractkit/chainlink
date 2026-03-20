package operation

import (
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	curse_mcms "github.com/smartcontractkit/chainlink-aptos/bindings/curse_mcms"
	"github.com/smartcontractkit/chainlink-aptos/contracts"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/dependency"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/utils"
)

// OP: DeployCurseMCMSOp deploys CurseMCMS to a new resource account.
type DeployCurseMCMSInput struct {
	MCMSAddress aptos.AccountAddress
	CCIPAddress aptos.AccountAddress
}

var DeployCurseMCMSOp = operations.NewOperation(
	"deploy-curse-mcms-op",
	Version1_0_0,
	"Deploys CurseMCMS Contract to a resource account",
	deployCurseMCMS,
)

func deployCurseMCMS(b operations.Bundle, deps dependency.AptosDeps, in DeployCurseMCMSInput) (aptos.AccountAddress, error) {
	seed := curse_mcms.DefaultSeed + time.Now().String()
	address, tx, err := bind.DeployPackageToResourceAccount(
		deps.AptosChain.DeployerSigner,
		deps.AptosChain.Client,
		contracts.CurseMCMS,
		seed,
		map[string]aptos.AccountAddress{
			"curse_mcms_owner":          deps.AptosChain.DeployerSigner.AccountAddress(),
			"ccip":                      in.CCIPAddress,
			"mcms":                      in.MCMSAddress,
			"mcms_register_entrypoints": aptos.AccountOne,
		},
		aptos.MaxGasAmount(1_000_000),
	)
	if err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("failed to deploy CurseMCMS: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return aptos.AccountAddress{}, fmt.Errorf("failed to confirm CurseMCMS deployment: %w", err)
	}
	return address, nil
}

// OP: ConfigureCurseMCMSOp configures a role on the CurseMCMS contract.
// Unlike the main MCMS configurer, this calls curse_mcms::set_config directly.
type ConfigureCurseMCMSInput struct {
	CurseMCMSAddress aptos.AccountAddress
	MCMSConfigs      mcmstypes.Config
	MCMSRole         aptosmcms.TimelockRole
}

var ConfigureCurseMCMSOp = operations.NewOperation(
	"configure-curse-mcms-op",
	Version1_0_0,
	"Configure CurseMCMS role (bypasser, canceller, or proposer)",
	configureCurseMCMS,
)

func configureCurseMCMS(b operations.Bundle, deps dependency.AptosDeps, in ConfigureCurseMCMSInput) (any, error) {
	curseMcmsBinding := curse_mcms.Bind(in.CurseMCMSAddress, deps.AptosChain.Client)
	opts := &bind.TransactOpts{Signer: deps.AptosChain.DeployerSigner}

	groupQuorum, groupParents, signerAddresses, signerGroups, err := mcmssdk.ExtractSetConfigInputs(&in.MCMSConfigs)
	if err != nil {
		return nil, fmt.Errorf("failed to extract config inputs: %w", err)
	}
	signers := make([][]byte, len(signerAddresses))
	for i, addr := range signerAddresses {
		signers[i] = addr.Bytes()
	}

	tx, err := curseMcmsBinding.CurseMCMS().SetConfig(
		opts,
		in.MCMSRole.Byte(),
		signers,
		signerGroups,
		groupQuorum[:],
		groupParents[:],
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to configure CurseMCMS: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return nil, fmt.Errorf("CurseMCMS configure transaction failed: %w", err)
	}
	return nil, nil
}

// OP: TransferCurseMCMSOwnershipToSelfOp transfers CurseMCMS ownership to
// the CurseMCMS resource account itself. The deployer is still the owner at
// this point so the call is signed directly.
var TransferCurseMCMSOwnershipToSelfOp = operations.NewOperation(
	"transfer-curse-mcms-ownership-to-self-op",
	Version1_0_0,
	"Transfer CurseMCMS ownership to self",
	transferCurseMCMSOwnershipToSelf,
)

func transferCurseMCMSOwnershipToSelf(b operations.Bundle, deps dependency.AptosDeps, curseMCMSAddress aptos.AccountAddress) (any, error) {
	opts := &bind.TransactOpts{Signer: deps.AptosChain.DeployerSigner}
	contractCurseMCMS := curse_mcms.Bind(curseMCMSAddress, deps.AptosChain.Client)
	tx, err := contractCurseMCMS.CurseMCMSAccount().TransferOwnershipToSelf(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to TransferOwnershipToSelf in CurseMCMS contract: %w", err)
	}
	if err := deps.AptosChain.Confirm(tx.Hash); err != nil {
		return nil, fmt.Errorf("CurseMCMS TransferOwnershipToSelf transaction failed: %w", err)
	}
	return nil, nil
}

// OP: AcceptCurseMCMSOwnershipOp encodes an AcceptOwnership transaction for
// the CurseMCMS contract. This must be submitted as a CurseMCMS proposal so
// the CurseMCMS resource account signer can call it.
var AcceptCurseMCMSOwnershipOp = operations.NewOperation(
	"accept-curse-mcms-ownership-op",
	Version1_0_0,
	"Generate Accept Ownership transaction for CurseMCMS",
	acceptCurseMCMSOwnership,
)

func acceptCurseMCMSOwnership(b operations.Bundle, deps dependency.AptosDeps, curseMCMSAddress aptos.AccountAddress) (mcmstypes.Transaction, error) {
	contractCurseMCMS := curse_mcms.Bind(curseMCMSAddress, deps.AptosChain.Client)
	moduleInfo, function, _, args, err := contractCurseMCMS.CurseMCMSAccount().Encoder().AcceptOwnership()
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode CurseMCMS AcceptOwnership: %w", err)
	}
	return utils.GenerateMCMSTx(curseMCMSAddress, moduleInfo, function, args)
}

// OP: SetCurseMCMSMinDelayOp encodes a CurseMCMS timelock transaction for
// timelock_update_min_delay. This must be submitted as a CurseMCMS proposal
// because the function requires TIMELOCK_ROLE.
type CurseMCMSMinDelayInput struct {
	CurseMCMSAddress aptos.AccountAddress
	TimelockMinDelay uint64
}

var SetCurseMCMSMinDelayOp = operations.NewOperation(
	"set-curse-mcms-min-delay-op",
	Version1_0_0,
	"Generate set timelock min delay transaction for CurseMCMS",
	setCurseMCMSMinDelay,
)

func setCurseMCMSMinDelay(b operations.Bundle, deps dependency.AptosDeps, in CurseMCMSMinDelayInput) (mcmstypes.Transaction, error) {
	curseMcmsBinding := curse_mcms.Bind(in.CurseMCMSAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := curseMcmsBinding.CurseMCMS().Encoder().TimelockUpdateMinDelay(in.TimelockMinDelay)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode CurseMCMS TimelockUpdateMinDelay: %w", err)
	}
	return utils.GenerateMCMSTx(in.CurseMCMSAddress, moduleInfo, function, args)
}
