package common

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth_types "github.com/ethereum/go-ethereum/core/types"
)

// RoninHeader is a copy of Header struct from ronin-chain repository
// https://github.com/ronin-chain/ronin/blob/8a3b1ada700a48c78c178c56c8a2b0a0b9bfd3e6/core/types/block.go#L302
type RoninHeader struct {
	ParentHash  common.Hash          `json:"parentHash"       gencodec:"required"`
	UncleHash   common.Hash          `json:"sha3Uncles"       gencodec:"required"`
	Coinbase    common.Address       `json:"miner"            gencodec:"required"`
	Root        common.Hash          `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash          `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash          `json:"receiptsRoot"     gencodec:"required"`
	Bloom       eth_types.Bloom      `json:"logsBloom"        gencodec:"required"`
	Difficulty  *big.Int             `json:"difficulty"       gencodec:"required"`
	Number      *big.Int             `json:"number"           gencodec:"required"`
	GasLimit    uint64               `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64               `json:"gasUsed"          gencodec:"required"`
	Time        uint64               `json:"timestamp"        gencodec:"required"`
	Extra       []byte               `json:"extraData"        gencodec:"required"`
	MixDigest   common.Hash          `json:"mixHash"`
	Nonce       eth_types.BlockNonce `json:"nonce"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`

	// These fields were added by EIP-4844 and is ignored in legacy headers.
	BlobGasUsed   *uint64 `json:"blobGasUsed" rlp:"optional"`
	ExcessBlobGas *uint64 `json:"excessBlobGas" rlp:"optional"`
}

// Unmarshalling from RPC-JSON response
func (h *RoninHeader) UnmarshalJSON(input []byte) error {
	type Header struct {
		ParentHash    *common.Hash          `json:"parentHash"       gencodec:"required"`
		UncleHash     *common.Hash          `json:"sha3Uncles"       gencodec:"required"`
		Coinbase      *common.Address       `json:"miner"            gencodec:"required"`
		Root          *common.Hash          `json:"stateRoot"        gencodec:"required"`
		TxHash        *common.Hash          `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash   *common.Hash          `json:"receiptsRoot"     gencodec:"required"`
		Bloom         *eth_types.Bloom      `json:"logsBloom"        gencodec:"required"`
		Difficulty    *hexutil.Big          `json:"difficulty"       gencodec:"required"`
		Number        *hexutil.Big          `json:"number"           gencodec:"required"`
		GasLimit      *hexutil.Uint64       `json:"gasLimit"         gencodec:"required"`
		GasUsed       *hexutil.Uint64       `json:"gasUsed"          gencodec:"required"`
		Time          *hexutil.Uint64       `json:"timestamp"        gencodec:"required"`
		Extra         *hexutil.Bytes        `json:"extraData"        gencodec:"required"`
		MixDigest     *common.Hash          `json:"mixHash"`
		Nonce         *eth_types.BlockNonce `json:"nonce"`
		BaseFee       *hexutil.Big          `json:"baseFeePerGas" rlp:"optional"`
		BlobGasUsed   *hexutil.Uint64       `json:"blobGasUsed" rlp:"optional"`
		ExcessBlobGas *hexutil.Uint64       `json:"excessBlobGas" rlp:"optional"`
	}
	var dec Header
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for Header")
	}
	h.ParentHash = *dec.ParentHash
	if dec.UncleHash == nil {
		return errors.New("missing required field 'sha3Uncles' for Header")
	}
	h.UncleHash = *dec.UncleHash
	if dec.Coinbase == nil {
		return errors.New("missing required field 'miner' for Header")
	}
	h.Coinbase = *dec.Coinbase
	if dec.Root == nil {
		return errors.New("missing required field 'stateRoot' for Header")
	}
	h.Root = *dec.Root
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionsRoot' for Header")
	}
	h.TxHash = *dec.TxHash
	if dec.ReceiptHash == nil {
		return errors.New("missing required field 'receiptsRoot' for Header")
	}
	h.ReceiptHash = *dec.ReceiptHash
	if dec.Bloom == nil {
		return errors.New("missing required field 'logsBloom' for Header")
	}
	h.Bloom = *dec.Bloom
	if dec.Difficulty == nil {
		return errors.New("missing required field 'difficulty' for Header")
	}
	h.Difficulty = (*big.Int)(dec.Difficulty)
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	h.Number = (*big.Int)(dec.Number)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Header")
	}
	h.GasLimit = uint64(*dec.GasLimit)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Header")
	}
	h.GasUsed = uint64(*dec.GasUsed)
	if dec.Time == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	h.Time = uint64(*dec.Time)
	if dec.Extra == nil {
		return errors.New("missing required field 'extraData' for Header")
	}
	h.Extra = *dec.Extra
	if dec.MixDigest != nil {
		h.MixDigest = *dec.MixDigest
	}
	if dec.Nonce != nil {
		h.Nonce = *dec.Nonce
	}
	if dec.BaseFee != nil {
		h.BaseFee = (*big.Int)(dec.BaseFee)
	}
	if dec.BlobGasUsed != nil {
		h.BlobGasUsed = (*uint64)(dec.BlobGasUsed)
	}
	if dec.ExcessBlobGas != nil {
		h.ExcessBlobGas = (*uint64)(dec.ExcessBlobGas)
	}
	return nil
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *RoninHeader) Hash() common.Hash {
	return rlpHash(h)
}
