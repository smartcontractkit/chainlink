package operation

import (
	"fmt"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	aptosmcms "github.com/smartcontractkit/mcms/sdk/aptos"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
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

// OP: InitializeAllowedCursersOp generates a main-MCMS transaction that
// calls rmn_remote::initialize_allowed_cursers_v2 with the CurseMCMS address.
type InitializeAllowedCursersInput struct {
	CCIPAddress      aptos.AccountAddress
	CurseMCMSAddress aptos.AccountAddress
}

var InitializeAllowedCursersOp = operations.NewOperation(
	"initialize-allowed-cursers-op",
	Version1_0_0,
	"Generates MCMS transaction to register CurseMCMS as an allowed curser on RMN Remote",
	initializeAllowedCursers,
)

func initializeAllowedCursers(b operations.Bundle, deps dependency.AptosDeps, in InitializeAllowedCursersInput) (mcmstypes.Transaction, error) {
	ccipBind := ccip.Bind(in.CCIPAddress, deps.AptosChain.Client)

	moduleInfo, function, _, args, err := ccipBind.RMNRemote().Encoder().InitializeAllowedCursersV2(
		[]aptos.AccountAddress{in.CurseMCMSAddress},
	)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to encode InitializeAllowedCursersV2: %w", err)
	}

	tx, err := utils.GenerateMCMSTx(in.CCIPAddress, moduleInfo, function, args)
	if err != nil {
		return mcmstypes.Transaction{}, fmt.Errorf("failed to generate MCMS transaction: %w", err)
	}
	return tx, nil
}
