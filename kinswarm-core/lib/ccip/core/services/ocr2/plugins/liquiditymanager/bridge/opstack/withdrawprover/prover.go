package withdrawprover

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_dispute_game_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_l2_output_oracle"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_portal"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/optimism_portal_2"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/abiutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/opstack/merkleutils"
)

type bedrockOutput struct {
	OutputRoot    [32]byte
	L1Timestamp   *big.Int
	L2BlockNumber *big.Int
	L2OutputIndex *big.Int
}

// Prover is able to prove Optimism withdrawal transactions from L2, on L1.
type Prover interface {
	// Prove returns all the information needed to prove a withdrawal from L2 to L1 on L1.
	// See docstrings of BedrockMessageProof and OutputRootProof for more details.
	Prove(ctx context.Context, withdrawalTxHash common.Hash) (
		messageProof BedrockMessageProof,
		err error,
	)
}

type EthClient interface {
	bind.ContractBackend
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*gethtypes.Receipt, error)
}

type prover struct {
	l1Client EthClient
	l2Client EthClient

	optimismPortal  optimism_portal.OptimismPortalInterface
	optimismPortal2 optimism_portal_2.OptimismPortal2Interface
	l2OutputOracle  optimism_l2_output_oracle.OptimismL2OutputOracleInterface
}

var (
	_ Prover = &prover{}
)

func New(
	l1Client,
	l2Client EthClient,
	optimismPortalAddress,
	l2OutputOracleAddress common.Address,
) (*prover, error) {
	if l1Client == nil || l2Client == nil {
		return nil, fmt.Errorf("args l1Client and l2Client must be non-nil")
	}

	optimismPortal, err := optimism_portal.NewOptimismPortal(optimismPortalAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("new optimism portal: %w", err)
	}

	optimismPortal2, err := optimism_portal_2.NewOptimismPortal2(optimismPortalAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("new optimism portal 2: %w", err)
	}

	l2OutputOracle, err := optimism_l2_output_oracle.NewOptimismL2OutputOracle(l2OutputOracleAddress, l1Client)
	if err != nil {
		return nil, fmt.Errorf("new l2 output oracle: %w", err)
	}

	return &prover{
		l1Client:        l1Client,
		l2Client:        l2Client,
		optimismPortal:  optimismPortal,
		optimismPortal2: optimismPortal2,
		l2OutputOracle:  l2OutputOracle,
	}, nil
}

func (p *prover) Prove(ctx context.Context, withdrawalTxHash common.Hash) (
	messageProof BedrockMessageProof,
	err error,
) {
	receipt, err := p.l2Client.TransactionReceipt(ctx, withdrawalTxHash)
	if err != nil {
		return messageProof, fmt.Errorf("get receipt for withdrawal tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	messagePassedLog := GetMessagePassedLog(receipt.Logs)
	if messagePassedLog == nil {
		return messageProof, fmt.Errorf("no message passed log found for withdrawal tx hash %s", withdrawalTxHash.Hex())
	}

	messagePassed, err := ParseMessagePassedLog(messagePassedLog)
	if err != nil {
		return messageProof, fmt.Errorf("parse message passed log tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	llmHash, err := hashLowLevelMessage(messagePassed)
	if err != nil {
		return messageProof, fmt.Errorf("hash low level message tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	messageSlot, err := hashMessageHash(llmHash)
	if err != nil {
		return messageProof, fmt.Errorf("hash message hash tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	messageBedrockOutput, err := p.getMessageBedrockOutput(
		ctx,
		receipt.BlockNumber,
	)
	if err != nil {
		return messageProof, fmt.Errorf("get message bedrock output tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	stateTrieProof, err := p.makeStateTrieProof(
		ctx,
		messageBedrockOutput.L2BlockNumber,
		messagePassed.Raw.Address,
		messageSlot,
	)
	if err != nil {
		return messageProof, fmt.Errorf("make state trie proof tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	header, err := p.l2Client.HeaderByNumber(ctx, messageBedrockOutput.L2BlockNumber)
	if err != nil {
		return messageProof, fmt.Errorf("get header by number tx hash %s: %w", withdrawalTxHash.Hex(), err)
	}

	return BedrockMessageProof{
		LowLevelMessage: messagePassed,
		OutputRootProof: OutputRootProof{
			Version:                  [32]byte{},
			StateRoot:                header.Root,
			MessagePasserStorageRoot: stateTrieProof.StorageRoot,
			LatestBlockHash:          header.Hash(),
		},
		WithdrawalProof: stateTrieProof.StorageProof,
		L2OutputIndex:   messageBedrockOutput.L2OutputIndex,
	}, nil
}

func (p *prover) makeStateTrieProof(
	ctx context.Context,
	l2BlockNumber *big.Int,
	address common.Address,
	slot [32]byte,
) (stateTrieProof, error) {
	var resp getProofResponse
	fmt.Println("calling eth_getProof, args: address:", address.String(), "slot:", hexutil.Encode(slot[:]), "l2BlockNumber:", l2BlockNumber.String())
	err := p.l2Client.CallContext(ctx, &resp, "eth_getProof",
		address, []string{hexutil.Encode(slot[:])}, hexutil.EncodeBig(l2BlockNumber))
	if err != nil {
		return stateTrieProof{}, fmt.Errorf("call eth_getProof with address %s, slot %s, l2BlockNumber %s: %w", address.String(), hexutil.Encode(slot[:]), l2BlockNumber.String(), err)
	}

	updatedProof, err := merkleutils.MaybeAddProofNode(
		crypto.Keccak256Hash(slot[:]), toProofBytes(resp.StorageProof[0].Proof))
	if err != nil {
		return stateTrieProof{}, fmt.Errorf("maybe add proof node: %w", err)
	}

	return stateTrieProof{
		AccountProof: toProofBytes(resp.AccountProof),
		StorageProof: updatedProof,
		StorageValue: resp.StorageProof[0].Value.ToInt(),
		StorageRoot:  resp.StorageHash,
	}, nil
}

func (p *prover) getMessageBedrockOutput(
	ctx context.Context,
	l2BlockNumber *big.Int,
) (bedrockOutput, error) {
	fpacEnabled, err := p.GetFPAC(ctx)
	if err != nil {
		return bedrockOutput{}, fmt.Errorf("get FPAC: %w", err)
	}
	if fpacEnabled {
		gameType, err2 := p.optimismPortal2.RespectedGameType(&bind.CallOpts{Context: ctx})
		if err2 != nil {
			return bedrockOutput{}, fmt.Errorf("get respected game type from portal: %w", err2)
		}

		disputeGameFactoryAddress, err2 := p.optimismPortal2.DisputeGameFactory(&bind.CallOpts{Context: ctx})
		if err2 != nil {
			return bedrockOutput{}, fmt.Errorf("get dispute game factory: %w", err2)
		}

		disputeGameFactory, err2 := optimism_dispute_game_factory.NewOptimismDisputeGameFactory(disputeGameFactoryAddress, p.l1Client)
		if err2 != nil {
			return bedrockOutput{}, fmt.Errorf("new dispute game factory: %w", err2)
		}

		gameCount, err2 := disputeGameFactory.GameCount(&bind.CallOpts{Context: ctx})
		if err2 != nil {
			return bedrockOutput{}, fmt.Errorf("get game count: %w", err2)
		}

		start := int64(0)
		if gameCount.Int64()-1 > start {
			start = gameCount.Int64() - 1
		}
		end := int64(100)
		if gameCount.Int64() < end {
			end = gameCount.Int64()
		}

		latestGames, err2 := disputeGameFactory.FindLatestGames(
			&bind.CallOpts{Context: ctx},
			gameType,
			big.NewInt(start),
			big.NewInt(end))
		if err2 != nil {
			return bedrockOutput{}, fmt.Errorf("find latest games: %w", err2)
		}

		for _, game := range latestGames {
			blockNumber, err2 := abiutils.UnpackUint256(game.ExtraData)
			if err2 != nil {
				return bedrockOutput{}, fmt.Errorf("unpack block number from dispute game: %w", err2)
			}

			if blockNumber.Cmp(l2BlockNumber) >= 0 {
				return bedrockOutput{
					OutputRoot:    game.RootClaim,
					L1Timestamp:   new(big.Int).SetUint64(game.Timestamp),
					L2BlockNumber: blockNumber,
					L2OutputIndex: game.Index,
				}, nil
			}
		}

		// if there's no match then we can't prove the message to the portal.
		return bedrockOutput{}, fmt.Errorf("no game found for block number %s", l2BlockNumber.String())
	}
	// Try to find the output index that corresponds to the block number attached to the message.
	// We'll explicitly handle "cannot get output" errors as a null return value, but anything else
	// needs to get thrown. Might need to revisit this in the future to be a little more robust
	// when connected to RPCs that don't return nice error messages.
	l2OutputIndex, err := p.l2OutputOracle.GetL2OutputIndexAfter(&bind.CallOpts{Context: ctx}, l2BlockNumber)
	if err != nil {
		return bedrockOutput{}, fmt.Errorf("[FPAC not enabled] get l2 output index after block number %s: %w", l2BlockNumber.String(), err)
	}

	// Now pull the proposal out given the output index. Should always work as long as the above
	// codepath completed successfully.
	proposal, err := p.l2OutputOracle.GetL2Output(&bind.CallOpts{Context: ctx}, l2OutputIndex)
	if err != nil {
		return bedrockOutput{}, fmt.Errorf("[FPAC not enabled] get l2 output for index %s from oracle: %w", l2OutputIndex.String(), err)
	}

	return bedrockOutput{
		OutputRoot:    proposal.OutputRoot,
		L1Timestamp:   proposal.Timestamp,
		L2BlockNumber: proposal.L2BlockNumber,
		L2OutputIndex: l2OutputIndex,
	}, nil
}

// GetFPAC returns whether FPAC (fault proof upgrade) is enabled on the optimism portal.
func (p *prover) GetFPAC(ctx context.Context) (bool, error) {
	semVer, err := p.optimismPortal.Version(&bind.CallOpts{Context: ctx})
	if err != nil {
		return false, fmt.Errorf("get version from portal: %w", err)
	}

	version := semver.MustParse(semVer)
	return version.GreaterThan(semver.MustParse("3.0.0")) || version.Equal(semver.MustParse("3.0.0")), nil
}

// StorageEntry represents a single entry in the state trie.
// See https://eips.ethereum.org/EIPS/eip-1186#specification
type StorageEntry struct {
	// Key is the storage key.
	Key hexutil.Bytes

	// Value is the value of the storage.
	Value hexutil.Big

	// Proof is an array of rlp-serialized MerkleTree-Nodes, starting with the storageHash-Node,
	// following the path of the keccak256(key) as path.
	Proof []hexutil.Bytes
}

// getProofResponse is the response from the eth_getProof JSON-RPC method.
// See https://eips.ethereum.org/EIPS/eip-1186#specification for more details
// on the response format.
// We only include the fields we need for our use case.
type getProofResponse struct {
	// AccountProof is an array of rlp-serialized MerkleTree-Nodes, starting with the stateRoot-Node,
	// following the path of the keccak256(address) as key.
	AccountProof []hexutil.Bytes `json:"accountProof"`

	// Storage hash is the keccak256 of the StorageRoot. All storage will deliver a MerkleProof
	// starting with this rootHash.
	StorageHash common.Hash `json:"storageHash"`

	// StorageProof is an array of StorageEntry objects.
	// See the StorageEntry struct for more details.
	StorageProof []StorageEntry `json:"storageProof"`
}

// stateTrieProof is an intermediate data structure that holds the state trie proof of a single storage entry.
type stateTrieProof struct {
	// AccountProof is an array of rlp-serialized MerkleTree-Nodes, starting with the stateRoot-Node,
	// following the path of the keccak256(address) as key.
	AccountProof [][]byte

	// StorageProof houses the proof of the storage entry.
	StorageProof [][]byte

	// StorageValue is the value of the storage entry.
	StorageValue *big.Int

	// StorageRoot is the root of the storage trie of the account in question.
	StorageRoot [32]byte
}
