package sequence

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
	"github.com/smartcontractkit/wsrpc/logger"

	sollegacy "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/solana"
	pdasol "github.com/smartcontractkit/cld-changesets/pkg/family/solana"

	timelockBindings "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/timelock"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms/sequence/operation"
	commonOps "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/operations"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

var (
	DeployMCMSWithTimelockSeq = operations.NewSequence(
		"deploy-access-controller-seq",
		&deployment.Version1_0_0,
		"Deploy AccessController,MCM and Timelock programs, Initialize them, set up role",
		deployMCMSWithTimelock,
	)
)

type (
	DeployMCMSWithTimelockInput struct {
		MCMConfig        cldfproposalutils.MCMSWithTimelockConfig
		TimelockMinDelay *big.Int
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

	// timelock
	err = deployTimelock(b, deps)
	if err != nil {
		return out, err
	}

	err = initTimelock(b, deps, in.TimelockMinDelay)
	if err != nil {
		return out, err
	}

	// roles
	err = setupRoles(b, deps)

	return out, err
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
			ProgramName:  solutils.ProgAccessController,
			Overallocate: true,
			Size:         solutils.GetProgramBufferBytes(solutils.ProgAccessController),
			ChainSel:     deps.Chain.ChainSelector(),
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

	err = deps.State.SetState(commontypes.AccessControllerProgram, programID, sollegacy.PDASeed{})
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initAccessController(b operations.Bundle, deps operation.Deps) error {
	roles := []cldf.ContractType{commontypes.ProposerAccessControllerAccount, commontypes.ExecutorAccessControllerAccount,
		commontypes.CancellerAccessControllerAccount, commontypes.BypasserAccessControllerAccount}
	for _, role := range roles {
		_, err := operations.ExecuteOperation(b, operation.InitAccessControllerOp, deps,
			operation.InitAccessControllerInput{
				ContractType: role,
				ChainSel:     deps.Chain.ChainSelector(),
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
		return nil
	}

	opOut, err := operations.ExecuteOperation(b, operation.DeployMCMProgramOp, commonOps.Deps{Chain: deps.Chain},
		commonOps.DeployInput{
			ProgramName:  solutils.ProgMCM,
			Overallocate: true,
			Size:         solutils.GetProgramBufferBytes(solutils.ProgMCM),
			ChainSel:     deps.Chain.ChainSelector(),
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

	err = deps.State.SetState(commontypes.ManyChainMultisigProgram, programID, sollegacy.PDASeed{})
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initMCM(b operations.Bundle, deps operation.Deps, cfg cldfproposalutils.MCMSWithTimelockConfig) error {
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
		_, err := operations.ExecuteOperation(b, operation.InitMCMOp, deps,
			operation.InitMCMInput{ContractType: cfg.ctype, MCMConfig: cfg.cfg, ChainSel: deps.Chain.ChainSelector()})
		if err != nil {
			return fmt.Errorf("failed to init config type:%q, err:%w", cfg.ctype, err)
		}
	}
	return nil
}

func deployTimelock(b operations.Bundle, deps operation.Deps) error {
	typeAndVersion := cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	log := logger.With(b.Logger, "chain", deps.Chain.String(), "contract", typeAndVersion.String())

	programID, _, err := deps.State.GetStateFromType(commontypes.RBACTimelock)
	if err != nil {
		return fmt.Errorf("failed to get timelock state: %w", err)
	}

	if !programID.IsZero() {
		log.Infow("using existing Timelock program", "programId", programID.String())
		if err := waitForProgramExecutable(b.GetContext(), deps.Chain.Client, programID, 30*time.Second); err != nil {
			return fmt.Errorf("timelock program not ready: %w", err)
		}
		return nil
	}

	opOut, err := operations.ExecuteOperation(b, operation.DeployTimelockOp, commonOps.Deps{Chain: deps.Chain},
		commonOps.DeployInput{
			ProgramName:  solutils.ProgTimelock,
			Overallocate: true,
			Size:         solutils.GetProgramBufferBytes(solutils.ProgTimelock),
			ChainSel:     deps.Chain.ChainSelector(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to deploy timelock program: %w", err)
	}

	programID = opOut.Output.ProgramID

	log.Infow("deployed timelock program", "programId", programID)

	err = deps.Datastore.Addresses().Add(datastore.AddressRef{
		ChainSelector: deps.Chain.ChainSelector(),
		Address:       programID.String(),
		Version:       &deployment.Version1_0_0,
		Type:          datastore.ContractType(commontypes.RBACTimelockProgram),
	})
	if err != nil {
		return fmt.Errorf("failed to add timelock to datastore: %w", err)
	}

	err = deps.State.SetState(commontypes.RBACTimelockProgram, programID, sollegacy.PDASeed{})
	if err != nil {
		return fmt.Errorf("failed to save onchain state: %w", err)
	}

	return nil
}

func initTimelock(b operations.Bundle, deps operation.Deps, minDelay *big.Int) error {
	if deps.State.TimelockProgram.IsZero() {
		return errors.New("mcm program is not deployed")
	}

	_, err := operations.ExecuteOperation(b, operation.InitTimelockOp, deps, operation.InitTimelockInput{
		ContractType: commontypes.RBACTimelock,
		ChainSel:     deps.Chain.ChainSelector(),
		MinDelay:     minDelay,
	})

	return err
}

// waitForProgramExecutable polls until the program account is executable (validator may still be loading BPF).
func waitForProgramExecutable(ctx context.Context, client *rpc.Client, programID solana.PublicKey, maxWait time.Duration) error {
	timeout := time.After(maxWait)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for program %s to be deployed", programID.String())
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resp, err := client.GetAccountInfo(ctx, programID)
			if err != nil {
				continue
			}
			if resp != nil && resp.Value != nil && resp.Value.Executable {
				return nil
			}
		}
	}
}

func setupRoles(b operations.Bundle, deps operation.Deps) error {
	proposerPDA := pdasol.GetMCMSignerPDA(deps.State.McmProgram, deps.State.ProposerMcmSeed)
	cancellerPDA := pdasol.GetMCMSignerPDA(deps.State.McmProgram, deps.State.CancellerMcmSeed)
	bypasserPDA := pdasol.GetMCMSignerPDA(deps.State.McmProgram, deps.State.BypasserMcmSeed)
	roles := []struct {
		pdas []solana.PublicKey
		role timelockBindings.Role
	}{
		{
			role: timelockBindings.Proposer_Role,
			pdas: []solana.PublicKey{proposerPDA},
		},
		{
			role: timelockBindings.Executor_Role,
			pdas: []solana.PublicKey{deps.Chain.DeployerKey.PublicKey()},
		},
		{
			role: timelockBindings.Canceller_Role,
			pdas: []solana.PublicKey{cancellerPDA, proposerPDA, bypasserPDA},
		},
		{
			role: timelockBindings.Bypasser_Role,
			pdas: []solana.PublicKey{bypasserPDA},
		},
	}
	for _, role := range roles {
		_, err := operations.ExecuteOperation(b, operation.AddAccessOp, deps, operation.AddAccessInput{
			Role:     role.role,
			Accounts: role.pdas,
		})
		if err != nil {
			return fmt.Errorf("failed to add access for role %d: %w", role.role, err)
		}
	}

	return nil
}
