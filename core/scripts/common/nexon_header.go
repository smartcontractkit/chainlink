package common

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// NexonHeader mirrors subnet-evm customtypes.HeaderSerializable; Nexon chains run Granite,
// whose extra header fields must be RLP-encoded for the hash to match the real blockhash.
type NexonHeader struct {
	ParentHash       common.Hash
	UncleHash        common.Hash
	Coinbase         common.Address
	Root             common.Hash
	TxHash           common.Hash
	ReceiptHash      common.Hash
	Bloom            types.Bloom
	Difficulty       *big.Int
	Number           *big.Int
	GasLimit         uint64
	GasUsed          uint64
	Time             uint64
	Extra            []byte
	MixDigest        common.Hash
	Nonce            types.BlockNonce
	BaseFee          *big.Int     `rlp:"optional"`
	BlockGasCost     *big.Int     `rlp:"optional"`
	BlobGasUsed      *uint64      `rlp:"optional"`
	ExcessBlobGas    *uint64      `rlp:"optional"`
	ParentBeaconRoot *common.Hash `rlp:"optional"`
	TimeMilliseconds *uint64      `rlp:"optional"`
	MinDelayExcess   *uint64      `rlp:"optional"`
}

func (h *NexonHeader) UnmarshalJSON(input []byte) error {
	var dec struct {
		ParentHash       common.Hash      `json:"parentHash"`
		UncleHash        common.Hash      `json:"sha3Uncles"`
		Coinbase         common.Address   `json:"miner"`
		Root             common.Hash      `json:"stateRoot"`
		TxHash           common.Hash      `json:"transactionsRoot"`
		ReceiptHash      common.Hash      `json:"receiptsRoot"`
		Bloom            types.Bloom      `json:"logsBloom"`
		Difficulty       *hexutil.Big     `json:"difficulty"`
		Number           *hexutil.Big     `json:"number"`
		GasLimit         hexutil.Uint64   `json:"gasLimit"`
		GasUsed          hexutil.Uint64   `json:"gasUsed"`
		Time             hexutil.Uint64   `json:"timestamp"`
		Extra            hexutil.Bytes    `json:"extraData"`
		MixDigest        common.Hash      `json:"mixHash"`
		Nonce            types.BlockNonce `json:"nonce"`
		BaseFee          *hexutil.Big     `json:"baseFeePerGas"`
		BlockGasCost     *hexutil.Big     `json:"blockGasCost"`
		BlobGasUsed      *hexutil.Uint64  `json:"blobGasUsed"`
		ExcessBlobGas    *hexutil.Uint64  `json:"excessBlobGas"`
		ParentBeaconRoot *common.Hash     `json:"parentBeaconBlockRoot"`
		TimeMilliseconds *hexutil.Uint64  `json:"timestampMilliseconds"`
		MinDelayExcess   *hexutil.Uint64  `json:"minDelayExcess"`
	}
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	*h = NexonHeader{
		ParentHash:       dec.ParentHash,
		UncleHash:        dec.UncleHash,
		Coinbase:         dec.Coinbase,
		Root:             dec.Root,
		TxHash:           dec.TxHash,
		ReceiptHash:      dec.ReceiptHash,
		Bloom:            dec.Bloom,
		Difficulty:       (*big.Int)(dec.Difficulty),
		Number:           (*big.Int)(dec.Number),
		GasLimit:         uint64(dec.GasLimit),
		GasUsed:          uint64(dec.GasUsed),
		Time:             uint64(dec.Time),
		Extra:            dec.Extra,
		MixDigest:        dec.MixDigest,
		Nonce:            dec.Nonce,
		BaseFee:          (*big.Int)(dec.BaseFee),
		BlockGasCost:     (*big.Int)(dec.BlockGasCost),
		BlobGasUsed:      optionalUint64(dec.BlobGasUsed),
		ExcessBlobGas:    optionalUint64(dec.ExcessBlobGas),
		ParentBeaconRoot: dec.ParentBeaconRoot,
		TimeMilliseconds: optionalUint64(dec.TimeMilliseconds),
		MinDelayExcess:   optionalUint64(dec.MinDelayExcess),
	}
	return nil
}

func (h *NexonHeader) Hash() common.Hash {
	enc, err := rlp.EncodeToBytes(h)
	if err != nil {
		panic(err)
	}
	return crypto.Keccak256Hash(enc)
}

func optionalUint64(v *hexutil.Uint64) *uint64 {
	if v == nil {
		return nil
	}
	u := uint64(*v)
	return &u
}
