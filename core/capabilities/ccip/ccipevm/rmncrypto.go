package ccipevm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// encodingUtilsAbi is the ABI for the EncodingUtils contract.
// Should be imported when gethwrappers are moved from ccip repo to core.
// nolint:lll
const encodingUtilsAbiRaw = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"DoNotDeploy","type":"error"},{"inputs":[{"internalType":"bytes32","name":"rmnReportVersion","type":"bytes32"},{"components":[{"internalType":"uint256","name":"destChainId","type":"uint256"},{"internalType":"uint64","name":"destChainSelector","type":"uint64"},{"internalType":"address","name":"rmnRemoteContractAddress","type":"address"},{"internalType":"address","name":"offrampAddress","type":"address"},{"internalType":"bytes32","name":"rmnHomeContractConfigDigest","type":"bytes32"},{"components":[{"internalType":"uint64","name":"sourceChainSelector","type":"uint64"},{"internalType":"bytes","name":"onRampAddress","type":"bytes"},{"internalType":"uint64","name":"minSeqNr","type":"uint64"},{"internalType":"uint64","name":"maxSeqNr","type":"uint64"},{"internalType":"bytes32","name":"merkleRoot","type":"bytes32"}],"internalType":"struct Internal.MerkleRoot[]","name":"destLaneUpdates","type":"tuple[]"}],"internalType":"struct RMNRemote.Report","name":"rmnReport","type":"tuple"}],"name":"_rmnReport","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
const addressEncodeAbiRaw = `[{"name":"method","type":"function","inputs":[{"name": "", "type": "address"}]}]`

var (
	encodingUtilsABI abi.ABI
	addressEncodeABI abi.ABI
)

func init() {
	var err error

	encodingUtilsABI, err = abi.JSON(strings.NewReader(encodingUtilsAbiRaw))
	if err != nil {
		panic(fmt.Errorf("failed to parse encoding utils ABI: %v", err))
	}

	addressEncodeABI, err = abi.JSON(strings.NewReader(addressEncodeAbiRaw))
	if err != nil {
		panic(fmt.Errorf("failed to parse address encode ABI: %v", err))
	}
}

const (
	// v is the recovery ID for ECDSA signatures. This implementation assumes that v is always 27.
	v = 27
)

// EVMRMNCrypto is the RMNCrypto implementation for EVM chains.
type EVMRMNCrypto struct{}

// Interface compliance check
var _ cciptypes.RMNCrypto = (*EVMRMNCrypto)(nil)

func NewEVMRMNCrypto() *EVMRMNCrypto {
	return &EVMRMNCrypto{}
}

// Should be replaced by gethwrapper types when they're available
type evmRMNRemoteReport struct {
	DestChainID                 *big.Int `abi:"destChainId"`
	DestChainSelector           uint64
	RmnRemoteContractAddress    common.Address
	OfframpAddress              common.Address
	RmnHomeContractConfigDigest [32]byte
	DestLaneUpdates             []evmInternalMerkleRoot
}

type evmInternalMerkleRoot struct {
	SourceChainSelector uint64
	OnRampAddress       []byte
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          [32]byte
}

func (r *EVMRMNCrypto) VerifyReportSignatures(
	_ context.Context,
	sigs []cciptypes.RMNECDSASignature,
	report cciptypes.RMNReport,
	signerAddresses []cciptypes.Bytes,
) error {
	if sigs == nil {
		return fmt.Errorf("no signatures provided")
	}
	if report.LaneUpdates == nil {
		return fmt.Errorf("no lane updates provided")
	}

	rmnVersionHash := crypto.Keccak256Hash([]byte(report.ReportVersion))

	evmLaneUpdates := make([]evmInternalMerkleRoot, len(report.LaneUpdates))
	for i, lu := range report.LaneUpdates {
		onRampAddress := common.BytesToAddress(lu.OnRampAddress)
		onRampAddrAbi, err := abiEncodeMethodInputs(addressEncodeABI, onRampAddress)
		if err != nil {
			return fmt.Errorf("ΑΒΙ encode onRampAddress: %w", err)
		}
		evmLaneUpdates[i] = evmInternalMerkleRoot{
			SourceChainSelector: uint64(lu.SourceChainSelector),
			OnRampAddress:       onRampAddrAbi,
			MinSeqNr:            uint64(lu.MinSeqNr),
			MaxSeqNr:            uint64(lu.MaxSeqNr),
			MerkleRoot:          lu.MerkleRoot,
		}
	}

	evmReport := evmRMNRemoteReport{
		DestChainID:                 report.DestChainID.Int,
		DestChainSelector:           uint64(report.DestChainSelector),
		RmnRemoteContractAddress:    common.HexToAddress(report.RmnRemoteContractAddress.String()),
		OfframpAddress:              common.HexToAddress(report.OfframpAddress.String()),
		RmnHomeContractConfigDigest: report.RmnHomeContractConfigDigest,
		DestLaneUpdates:             evmLaneUpdates,
	}

	abiEnc, err := encodingUtilsABI.Methods["_rmnReport"].Inputs.Pack(rmnVersionHash, evmReport)
	if err != nil {
		return fmt.Errorf("failed to ABI encode args: %w", err)
	}

	signedHash := crypto.Keccak256Hash(abiEnc)

	// keep track of the previous signer for validating signers ordering
	prevSignerAddr := common.Address{}

	for _, sig := range sigs {
		recoveredAddress, err := recoverAddressFromSig(
			v,
			sig.R,
			sig.S,
			signedHash[:],
		)
		if err != nil {
			return fmt.Errorf("failed to recover public key from signature: %w", err)
		}

		// make sure that signers are ordered correctly (ASC addresses).
		if bytes.Compare(prevSignerAddr.Bytes(), recoveredAddress.Bytes()) == 1 {
			return fmt.Errorf("signers are not ordered correctly")
		}
		prevSignerAddr = recoveredAddress

		// Check if the public key is in the list of the provided RMN nodes
		found := false
		for _, signerAddr := range signerAddresses {
			signerAddrEvm := common.BytesToAddress(signerAddr)
			if signerAddrEvm == recoveredAddress {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("the recovered public key does not match any signer address, verification failed")
		}
	}

	return nil
}

// recoverAddressFromSig Recovers a public address from an ECDSA signature using r, s, v, and the hash of the message.
func recoverAddressFromSig(v int, r, s [32]byte, hash []byte) (common.Address, error) {
	// Ensure v is either 27 or 28 (as used in Ethereum)
	if v != 27 && v != 28 {
		return common.Address{}, errors.New("v must be 27 or 28")
	}

	// Construct the signature by concatenating r, s, and the recovery ID (v - 27 to convert to 0/1)
	sig := append(r[:], s[:]...)
	sig = append(sig, byte(v-27))

	// Recover the public key bytes from the signature and message hash
	pubKeyBytes, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Convert the recovered public key to an ECDSA public key
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to unmarshal public key: %v", err)
	} // or SigToPub

	return crypto.PubkeyToAddress(*pubKey), nil
}
