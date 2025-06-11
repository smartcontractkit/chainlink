package sequence

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms/sequence/operation"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/wsrpc/logger"
)

var (
	DeployMCMSWithTimelockSeq = operations.NewSequence(
		"deploy-access-controller-seq",
		&deployment.Version1_0_0,
		"Deploys AccessController,MCM and Timelock contracts and initializes them",
		deployMCMSWithTimelock,
	)
)

type (
	DeployMCMSWithTimelockInput struct {
		MCMConfig commontypes.MCMSWithTimelockConfigV2
	}

	DeployMCMSWithTimelockOutput struct{}
)

func deployMCMSWithTimelock(b operations.Bundle, deps operation.Deps, in DeployMCMSWithTimelockInput) (DeployMCMSWithTimelockOutput, error) {
	var out DeployMCMSWithTimelockOutput

	//  access controller
	err := deployAccessController(b, deps)
	if err != nil {
		return out, err
	}

	err = initAccessController(b, deps)
	if err != nil {
		return out, err
	}

	// mcm
	err = deployMCM(b, deps)
	if err != nil {
		return out, err
	}

	err = initMCM(b, deps, in.MCMConfig)
	if err != nil {
		return out, err
	}
	// TODO timelock

	return out, nil
}

func deployAccessController(b operations.Bundle, deps operation.Deps) error {
	typeAndVersion := cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	log := logger.With(b.Logger, "chain", deps.Chain.String(), "contract", typeAndVersion.String())

	programID, _, err := deps.State.GetStateFromType(commontypes.AccessControllerProgram)
	if err != nil {
		return fmt.Errorf("failed to get access controller program state: %w", err)
	}

	if !programID.IsZero() {
		log.Infow("using existing AccessController program", "programId", programID)
		return nil
	}

	opOut, err := operations.ExecuteOperation(b, operation.DeployAccessControllerOp, commonOps.Deps{Chain: deps.Chain},
		commonOps.DeployInput{
			ProgramName:  deployment.AccessControllerProgramName,
			Overallocate: true,
			Size:         deployment.SolanaProgramBytes[deployment.AccessControllerProgramName],
		},
	)
	if err != nil {
		return fmt.Errorf("failed to deploy access controller: %w", err)
	}
	programID = opOut.Output.ProgramID

	log.Infow("deployed access controller contract", "programId", programID)

	err = deps.Datastore.Addresses().Add(datastore.AddressRef{
		ChainSelector: deps.Chain.ChainSelector(),
		Address:       programID.String(),
		Version:       &deployment.Version1_0_0,
		Type:          datastore.ContractType(commontypes.AccessControllerProgram),
	})
	if err != nil {
		return fmt.Errorf("failed to add access controller to datastore: %w", err)
	}

	err = deps.State.SetState(commontypes.AccessControllerProgram, programID, state.PDASeed{})
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initAccessController(b operations.Bundle, deps operation.Deps) error {
	roles := []cldf.ContractType{commontypes.ProposerAccessControllerAccount, commontypes.ExecutorAccessControllerAccount,
		commontypes.CancellerAccessControllerAccount, commontypes.BypasserAccessControllerAccount}
	for _, role := range roles {
		_, err := operations.ExecuteOperation(b, operation.InitAccessControllerOp, operation.Deps{State: deps.State, Chain: deps.Chain},
			operation.InitAccessControllerInput{
				ContractType: role,
			})
		if err != nil {
			return fmt.Errorf("failed to init access controller account role %q: %w", role, err)
		}
	}

	return nil
}

func deployMCM(b operations.Bundle, deps operation.Deps) error {
	typeAndVersion := cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	log := logger.With(b.Logger, "chain", deps.Chain.String(), "contract", typeAndVersion.String())

	programID, _, err := deps.State.GetStateFromType(commontypes.ManyChainMultisigProgram)
	if err != nil {
		return fmt.Errorf("failed to get mcm state: %w", err)
	}
	if !programID.IsZero() {
		log.Infow("using existing MCM program", "programId", programID.String())
	}

	opOut, err := operations.ExecuteOperation(b, operation.DeployMCMProgramOp, commonOps.Deps{Chain: deps.Chain},
		commonOps.DeployInput{
			ProgramName:  deployment.McmProgramName,
			Overallocate: true,
			Size:         deployment.SolanaProgramBytes[deployment.McmProgramName],
		},
	)
	if err != nil {
		return fmt.Errorf("failed to deploy mcm program : %w", err)
	}
	programID = opOut.Output.ProgramID

	log.Infow("deployed mcm contract", "programId", programID)

	err = deps.Datastore.Addresses().Add(datastore.AddressRef{
		ChainSelector: deps.Chain.ChainSelector(),
		Address:       programID.String(),
		Version:       &deployment.Version1_0_0,
		Type:          datastore.ContractType(commontypes.ManyChainMultisig),
	})
	if err != nil {
		return fmt.Errorf("failed to add mcm to datastore: %w", err)
	}

	err = deps.State.SetState(commontypes.ManyChainMultisig, programID, state.PDASeed{})
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initMCM(b operations.Bundle, deps operation.Deps, cfg commontypes.MCMSWithTimelockConfigV2) error {
	configs := []struct {
		ctype cldf.ContractType
		cfg   mcmsTypes.Config
	}{
		{
			commontypes.BypasserManyChainMultisig,
			cfg.Bypasser,
		},
		{
			commontypes.CancellerManyChainMultisig,
			cfg.Canceller,
		},
		{
			commontypes.ProposerManyChainMultisig,
			cfg.Proposer,
		},
	}

	for _, cfg := range configs {
		_, err := operations.ExecuteOperation(b, operation.InitMCMOp, deps, operation.InitMCMInput{ContractType: cfg.ctype, MCMConfig: cfg.cfg})
		if err != nil {
			return fmt.Errorf("failed to init config type:%q, err:%w", cfg.ctype, err)
		}
	}
	return nil
}
