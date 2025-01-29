package changeset

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/deployment"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	AddressLookupTable deployment.ContractType = "AddressLookupTable"
	TokenPool          deployment.ContractType = "TokenPool"
	Receiver           deployment.ContractType = "Receiver"
)

// SolChainState holds a Go binding for all the currently deployed CCIP programs
// on a chain. If a binding is nil, it means here is no such contract on the chain.
type SolCCIPChainState struct {
	LinkToken          solana.PublicKey
	Router             solana.PublicKey
	Timelock           solana.PublicKey
	AddressLookupTable solana.PublicKey // for chain writer
	Receiver           solana.PublicKey // for tests only
}

func LoadOnchainStateSolana(e deployment.Environment) (CCIPOnChainState, error) {
	state := CCIPOnChainState{
		SolChains: make(map[uint64]SolCCIPChainState),
	}
	for chainSelector, chain := range e.SolChains {
		addresses, err := e.ExistingAddresses.AddressesForChain(chainSelector)
		if err != nil {
			// Chain not found in address book, initialize empty
			if !errors.Is(err, deployment.ErrChainNotFound) {
				return state, err
			}
			addresses = make(map[string]deployment.TypeAndVersion)
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
		case deployment.NewTypeAndVersion(commontypes.LinkToken, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.LinkToken = pub
		case deployment.NewTypeAndVersion(Router, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.Router = pub
		case deployment.NewTypeAndVersion(AddressLookupTable, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.AddressLookupTable = pub
		case deployment.NewTypeAndVersion(Receiver, deployment.Version1_0_0).String():
			pub := solana.MustPublicKeyFromBase58(address)
			state.Receiver = pub
		default:
			return state, fmt.Errorf("unknown contract %s", tvStr)
		}
	}
	return state, nil
}

// GetRouterConfigPDA returns the PDA for the "config" account.
func GetRouterConfigPDA(ccipRouterProgramID solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("config")},
		ccipRouterProgramID,
	)
	return pda
}

// GetRouterStatePDA returns the PDA for the "state" account.
func GetRouterStatePDA(ccipRouterProgramID solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("state")},
		ccipRouterProgramID,
	)
	return pda
}

// GetExternalExecutionConfigPDA returns the PDA for the "external_execution_config" account.
func GetExternalExecutionConfigPDA(ccipRouterProgramID solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("external_execution_config")},
		ccipRouterProgramID,
	)
	return pda
}

// GetExternalTokenPoolsSignerPDA returns the PDA for the "external_token_pools_signer" account.
func GetExternalTokenPoolsSignerPDA(ccipRouterProgramID solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{[]byte("external_token_pools_signer")},
		ccipRouterProgramID,
	)
	return pda
}

// GetSolanaSourceChainStatePDA returns the PDA for the "source_chain_state" account for Solana.
func GetSolanaSourceChainStatePDA(ccipRouterProgramID solana.PublicKey, solanaChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("source_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, solanaChainSelector),
		},
		ccipRouterProgramID,
	)
	return pda
}

// GetSolanaDestChainStatePDA returns the PDA for the "dest_chain_state" account for Solana.
func GetSolanaDestChainStatePDA(ccipRouterProgramID solana.PublicKey, solanaChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("dest_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, solanaChainSelector),
		},
		ccipRouterProgramID,
	)
	return pda
}

// GetEvmSourceChainStatePDA returns the PDA for the "source_chain_state" account for EVM.
func GetEvmSourceChainStatePDA(ccipRouterProgramID solana.PublicKey, evmChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("source_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, evmChainSelector),
		},
		ccipRouterProgramID,
	)
	return pda
}

// GetEvmDestChainStatePDA returns the PDA for the "dest_chain_state" account for EVM.
func GetEvmDestChainStatePDA(ccipRouterProgramID solana.PublicKey, evmChainSelector uint64) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress(
		[][]byte{
			[]byte("dest_chain_state"),
			binary.LittleEndian.AppendUint64([]byte{}, evmChainSelector),
		},
		ccipRouterProgramID,
	)
	return pda
}

func GetReceiverTargetAccountPDA(ccipReceiverProgram solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress([][]byte{[]byte("counter")}, ccipReceiverProgram)
	return pda
}

func GetReceiverExternalExecutionConfigPDA(ccipReceiverProgram solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress([][]byte{[]byte("external_execution_config")}, ccipReceiverProgram)
	return pda
}

func GetTokenAdminRegistryPDA(ccipRouterProgramID, tokenMint solana.PublicKey) solana.PublicKey {
	pda, _, _ := solana.FindProgramAddress([][]byte{[]byte("token_admin_registry"), tokenMint.Bytes()}, ccipRouterProgramID)
	return pda
}
