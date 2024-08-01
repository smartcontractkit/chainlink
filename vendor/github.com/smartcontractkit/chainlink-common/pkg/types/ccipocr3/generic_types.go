package ccipocr3

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TokenPrice struct {
	TokenID types.Account `json:"tokenID"`
	Price   BigInt        `json:"price"`
}

func NewTokenPrice(tokenID types.Account, price *big.Int) TokenPrice {
	return TokenPrice{
		TokenID: tokenID,
		Price:   BigInt{price},
	}
}

type GasPriceChain struct {
	GasPrice BigInt        `json:"gasPrice"`
	ChainSel ChainSelector `json:"chainSel"`
}

func NewGasPriceChain(gasPrice *big.Int, chainSel ChainSelector) GasPriceChain {
	return GasPriceChain{
		GasPrice: NewBigInt(gasPrice),
		ChainSel: chainSel,
	}
}

type SeqNum uint64

func (s SeqNum) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func NewSeqNumRange(start, end SeqNum) SeqNumRange {
	return SeqNumRange{start, end}
}

// SeqNumRange defines an inclusive range of sequence numbers.
type SeqNumRange [2]SeqNum

func (s SeqNumRange) Start() SeqNum {
	return s[0]
}

func (s SeqNumRange) End() SeqNum {
	return s[1]
}

func (s *SeqNumRange) SetStart(v SeqNum) {
	s[0] = v
}

func (s *SeqNumRange) SetEnd(v SeqNum) {
	s[1] = v
}

// Overlaps returns true if the two ranges overlap.
func (s SeqNumRange) Overlaps(other SeqNumRange) bool {
	return s.Start() <= other.End() && other.Start() <= s.End()
}

// Contains returns true if the range contains the given sequence number.
func (s SeqNumRange) Contains(seq SeqNum) bool {
	return s.Start() <= seq && seq <= s.End()
}

func (s SeqNumRange) String() string {
	return fmt.Sprintf("[%d -> %d]", s[0], s[1])
}

type ChainSelector uint64

func (c ChainSelector) String() string {
	return fmt.Sprintf("ChainSelector(%d)", c)
}

// Message is the generic Any2Any message type for CCIP messages.
// It represents, in particular, a message emitted by a CCIP onramp.
// The message is expected to be consumed by a CCIP offramp after
// translating it into the appropriate format for the destination chain.
type Message struct {
	// Header is the family-agnostic header for OnRamp and OffRamp messages.
	// This is always set on all CCIP messages.
	Header RampMessageHeader `json:"header"`
	// Sender address on the source chain.
	// i.e if the source chain is EVM, this is an abi-encoded EVM address.
	Sender Bytes `json:"sender"`
	// Data is the arbitrary data payload supplied by the message sender.
	Data Bytes `json:"data"`
	// Receiver is the receiver address on the destination chain.
	// This is encoded in the destination chain family specific encoding.
	// i.e if the destination is EVM, this is abi.encode(receiver).
	Receiver Bytes `json:"receiver"`
	// ExtraArgs is destination-chain specific extra args, such as the gasLimit for EVM chains.
	ExtraArgs Bytes `json:"extraArgs"`
	// FeeToken is the fee token address.
	// i.e if the source chain is EVM, len(FeeToken) == 20 (i.e, is not abi-encoded).
	FeeToken Bytes `json:"feeToken"`
	// FeeTokenAmount is the amount of fee tokens paid.
	FeeTokenAmount BigInt `json:"feeTokenAmount"`
	// TokenAmounts is the array of tokens and amounts to transfer.
	TokenAmounts []RampTokenAmount `json:"tokenAmounts"`
}

func (c Message) String() string {
	js, _ := json.Marshal(c)
	return string(js)
}

// RampMessageHeader is the family-agnostic header for OnRamp and OffRamp messages.
// The MessageID is not expected to match MsgHash, since it may originate from a different
// ramp family.
type RampMessageHeader struct {
	// MessageID is a unique identifier for the message, it should be unique across all chains.
	// It is generated on the chain that the CCIP send is requested (i.e. the source chain of a message).
	MessageID Bytes32 `json:"messageId"`
	// SourceChainSelector is the chain selector of the chain that the message originated from.
	SourceChainSelector ChainSelector `json:"sourceChainSelector,string"`
	// DestChainSelector is the chain selector of the chain that the message is destined for.
	DestChainSelector ChainSelector `json:"destChainSelector,string"`
	// SequenceNumber is an auto-incrementing sequence number for the message.
	// Not unique across lanes.
	SequenceNumber SeqNum `json:"seqNum,string"`
	// Nonce is the nonce for this lane for this sender, not unique across senders/lanes
	Nonce uint64 `json:"nonce"`

	// MsgHash is the hash of all the message fields.
	// NOTE: The field is expected to be empty, and will be populated by the plugin using the MsgHasher interface.
	MsgHash Bytes32 `json:"msgHash"` // populated

	// OnRamp is the address of the onramp that sent the message.
	// NOTE: This is populated by the ccip reader. Not emitted explicitly onchain.
	OnRamp Bytes `json:"onRamp"`
}

// RampTokenAmount represents the family-agnostic token amounts used for both OnRamp & OffRamp messages.
type RampTokenAmount struct {
	// SourcePoolAddress is the source pool address, encoded according to source family native encoding scheme.
	// This value is trusted as it was obtained through the onRamp. It can be relied upon by the destination
	// pool to validate the source pool.
	SourcePoolAddress Bytes `json:"sourcePoolAddress"`
	// DestTokenAddress is the address of the destination token, abi encoded in the case of EVM chains.
	// This value is UNTRUSTED as any pool owner can return whatever value they want.
	DestTokenAddress Bytes `json:"destTokenAddress"`
	// ExtraData is optional pool data to be transferred to the destination chain. Be default this is capped at
	// CCIP_LOCK_OR_BURN_V1_RET_BYTES bytes. If more data is required, the TokenTransferFeeConfig.destBytesOverhead
	// has to be set for the specific token.
	ExtraData Bytes `json:"extraData"`
	// Amount is the amount of tokens to be transferred.
	Amount BigInt `json:"amount"`
}
