package changeset

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	// TestXXXMCMSSigner is a throwaway private key used for signing MCMS proposals.
	// in tests.
	TestXXXMCMSSigner *ecdsa.PrivateKey
)

func init() {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	TestXXXMCMSSigner = key
}

func SingleGroupMCMS(t *testing.T) config.Config {
	publicKey := TestXXXMCMSSigner.Public().(*ecdsa.PublicKey)
	// Convert the public key to an Ethereum address
	address := crypto.PubkeyToAddress(*publicKey)
	c, err := config.NewConfig(1, []common.Address{address}, []config.Config{})
	require.NoError(t, err)
	return *c
}

func SignProposal(env deployment.Environment, proposal *timelock.MCMSWithTimelockProposal) (*mcms.Executor, error) {
	executorClients := make(map[mcms.ChainIdentifier]mcms.ContractDeployBackend)
	for _, chain := range env.Chains {
		chainselc, exists := chainsel.ChainBySelector(chain.Selector)
		if !exists {
			return nil, fmt.Errorf("failed to find chain by selector: %d", chain.Selector)
		}
		chainSel := mcms.ChainIdentifier(chainselc.Selector)
		executorClients[chainSel] = chain.Client
	}
	executor, err := proposal.ToExecutor(true)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proposal to executor: %w", err)
	}
	payload, err := executor.SigningHash()
	if err != nil {
		return nil, fmt.Errorf("failed to get signing hash: %w", err)
	}
	// Sign the payload
	sig, err := crypto.Sign(payload.Bytes(), TestXXXMCMSSigner)
	if err != nil {
		return nil, fmt.Errorf("failed to sign payload: %w", err)
	}
	mcmSig, err := mcms.NewSignatureFromBytes(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature from bytes: %w", err)
	}
	executor.Proposal.AddSignature(mcmSig)
	if err := executor.Proposal.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate proposal: %w", err)
	}
	return executor, nil
}

func ExecuteProposal(env deployment.Environment, executor *mcms.Executor, timelock *owner_helpers.RBACTimelock, sel uint64) error {
	env.Logger.Infow("Executing proposal on chain", "selector", sel)
	// Set the root.
	tx, err2 := executor.SetRootOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, mcms.ChainIdentifier(sel))
	if err2 != nil {
		return fmt.Errorf("failed to set root: %w", deployment.MaybeDataErr(err2))
	}
	_, err2 = env.Chains[sel].Confirm(tx)
	if err2 != nil {
		return fmt.Errorf("failed to confirm set-root: %w", err2)
	}

	// TODO: This sort of helper probably should move to the MCMS lib.
	// Execute all the transactions in the proposal which are for this chain.
	for _, chainOp := range executor.Operations[mcms.ChainIdentifier(sel)] {
		for idx, op := range executor.ChainAgnosticOps {
			if bytes.Equal(op.Data, chainOp.Data) && op.To == chainOp.To {
				opTx, err3 := executor.ExecuteOnChain(env.Chains[sel].Client, env.Chains[sel].DeployerKey, idx)
				if err3 != nil {
					return fmt.Errorf("failed to execute operation: %w", deployment.MaybeDataErr(err3))
				}
				block, err3 := env.Chains[sel].Confirm(opTx)
				if err3 != nil {
					return fmt.Errorf("failed to confirm operation ExecuteOnChain: %w", err3)
				}
				env.Logger.Infow("executed", "chainOp", chainOp)
				it, err3 := timelock.FilterCallScheduled(&bind.FilterOpts{
					Start:   block,
					End:     &block,
					Context: context.Background(),
				}, nil, nil)
				if err3 != nil {
					return fmt.Errorf("failed to filter call scheduled: %w", err3)
				}
				var calls []owner_helpers.RBACTimelockCall
				var pred, salt [32]byte
				for it.Next() {
					// Note these are the same for the whole batch, can overwrite
					pred = it.Event.Predecessor
					salt = it.Event.Salt
					env.Logger.Infow("scheduled", "Event", it.Event)
					calls = append(calls, owner_helpers.RBACTimelockCall{
						Target: it.Event.Target,
						Data:   it.Event.Data,
						Value:  it.Event.Value,
					})
				}
				tx, err := timelock.ExecuteBatch(
					env.Chains[sel].DeployerKey, calls, pred, salt)
				if err != nil {
					return fmt.Errorf("failed to execute batch: %w", deployment.MaybeDataErr(err))
				}
				_, err = env.Chains[sel].Confirm(tx)
				if err != nil {
					return fmt.Errorf("failed to confirm execute batch: %w", err)
				}
			}
		}
	}
	return nil
}
