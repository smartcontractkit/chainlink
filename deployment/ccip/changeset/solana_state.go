package changeset

import (
	"encoding/binary"
	"fmt"

	"errors"

	ag_binary "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink/deployment"
)

var (
	LinkToken     deployment.ContractType = "LinkToken"
	SolCcipRouter deployment.ContractType = "SolCcipRouter"
	SolReceiver   deployment.ContractType = "SolCcipReceiver"
	SolTokenPool  deployment.ContractType = "SolTokenPool"
)

type SolCCIPRouter struct {
	SolanaChainSelector             uint64
	DefaultGasLimit                 ag_binary.Uint128
	DefaultAllowOutOfOrderExecution bool
	EnableExecutionAfter            int64
	// Accounts:
	Config                  solana.PublicKey
	State                   solana.PublicKey
	Authority               solana.PublicKey
	SystemProgram           solana.PublicKey
	Program                 solana.PublicKey
	ProgramData             solana.PublicKey
	ExternalExecutionConfig solana.PublicKey
	TokenPoolsSigner        solana.PublicKey
}

// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type SolCCIPChainState struct {
	LinkToken       solana.PublicKey
	SolCcipRouter   solana.PublicKey
	SolCcipReceiver solana.PublicKey
	TokenPool       solana.PublicKey
	Weth9           solana.PublicKey
	Timelock        solana.PublicKey
}

func LoadOnchainStateSolana(e deployment.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		SolChains: make(map[uint64]SolCCIPChainState),
	}
	for chainSelector, chain := range e.SolChains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if errors.Is(err, deployment.ErrChainNotFound) {
				addresses = make(map[string]deployment.TypeAndVersion)
			} else {
				return state, err
			}
		}
		chainState, err := LoadChainStateSolana(chain, addresses)
		if err != nil {
			return state, err
		}
		state.SolChains[chainSelector] = chainState

	}
	return state, nil
}

// LoadChainStateSolana Loads all state for a SolChain into state
func LoadChainStateSolana(chain deployment.SolChain, addresses map[string]deployment.TypeAndVersion) (SolCCIPChainState, error) {
	var state SolCCIPChainState
	for address, tvStr := range addresses {
		switch tvStr.String() {
		case deployment.NewTypeAndVersion(SolCcipRouter, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			//ccip_router.SetProgramID(pub)
			state.SolCcipRouter = pub
		case deployment.NewTypeAndVersion(SolReceiver, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			//ccip_receiver.SetProgramID(pub)
			state.SolCcipReceiver = pub
		case deployment.NewTypeAndVersion(SolTokenPool, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			//token_pool.SetProgramID(pub)
			state.TokenPool = pub
		case deployment.NewTypeAndVersion(LinkToken, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			// token_pool.SetProgramID(pub)
			state.LinkToken = pub
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}

// GetRouterConfigPDA returns the PDA for the "config" account.
func GetRouterConfigPDA(CcipRouterProgram solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("config")},
		CcipRouterProgram,
	)
	return pda
}

// GetRouterStatePDA returns the PDA for the "state" account.
func GetRouterStatePDA(CcipRouterProgram solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("state")},
		CcipRouterProgram,
	)
	return pda
}

// GetExternalExecutionConfigPDA returns the PDA for the "external_execution_config" account.
func GetExternalExecutionConfigPDA(CcipRouterProgram solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("external_execution_config")},
		CcipRouterProgram,
	)
	return pda
}

// GetExternalTokenPoolsSignerPDA returns the PDA for the "external_token_pools_signer" account.
func GetExternalTokenPoolsSignerPDA(CcipRouterProgram solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("external_token_pools_signer")},
		CcipRouterProgram,
	)
	return pda
}

// GetSolanaSourceChainStatePDA returns the PDA for the "source_chain_state" account for Solana.
func GetSolanaSourceChainStatePDA(CcipRouterProgram solana.PublicKey, SolanaChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("source_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, SolanaChainSelector),
		},
		CcipRouterProgram,
	)
	return pda
}

// GetSolanaDestChainStatePDA returns the PDA for the "dest_chain_state" account for Solana.
func GetSolanaDestChainStatePDA(CcipRouterProgram solana.PublicKey, SolanaChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("dest_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, SolanaChainSelector),
		},
		CcipRouterProgram,
	)
	return pda
}

// GetEvmSourceChainStatePDA returns the PDA for the "source_chain_state" account for EVM.
func GetEvmSourceChainStatePDA(CcipRouterProgram solana.PublicKey, EvmChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("source_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, EvmChainSelector),
		},
		CcipRouterProgram,
	)
	return pda
}

// GetEvmDestChainStatePDA returns the PDA for the "dest_chain_state" account for EVM.
func GetEvmDestChainStatePDA(CcipRouterProgram solana.PublicKey, EvmChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("dest_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, EvmChainSelector),
		},
		CcipRouterProgram,
	)
	return pda
}
