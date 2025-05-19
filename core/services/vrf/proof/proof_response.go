package proof

// Contains logic/data for mandatorily mixing VRF seeds with the hash of the
// block in which a VRF request appeared

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2plus_interface"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
)

// ProofResponse is the data which is sent back to the VRFCoordinator, so that
// it can verify that the seed the oracle finally used is correct.
type ProofResponse struct {
	// Approximately the proof which will be checked on-chain. Note that this
	// contains the pre-seed in place of the final seed. That should be computed
	// as in FinalSeed.
	P        vrfkey.Proof
	PreSeed  Seed   // Seed received during VRF request
	BlockNum uint64 // Height of the block in which this request was made
}

// OnChainResponseLength is the length of the MarshaledOnChainResponse. The
// extra 32 bytes are for blocknumber (as a uint256), which goes at the end. The
// seed is rewritten with the preSeed. (See MarshalForVRFCoordinator and
// ProofResponse#ActualProof.)
const OnChainResponseLength = ProofLength +
	32 // blocknum

// MarshaledOnChainResponse is the flat bytes which are sent back to the
// VRFCoordinator.
type MarshaledOnChainResponse [OnChainResponseLength]byte

// MarshalForVRFCoordinator constructs the flat bytes which are sent to the
// VRFCoordinator.
func (p *ProofResponse) MarshalForVRFCoordinator() (
	response MarshaledOnChainResponse, err error) {
	solidityProof, err := SolidityPrecalculations(&p.P)
	if err != nil {
		return MarshaledOnChainResponse{}, errors.Wrap(err,
			"while marshaling proof for VRFCoordinator")
	}
	// Overwrite seed input to the VRF proof generator with the seed the
	// VRFCoordinator originally requested, so that it can identify the request
	// corresponding to this response, and compute the final seed itself using the
	// blockhash it infers from the block number.
	solidityProof.P.Seed = common.BytesToHash(p.PreSeed[:]).Big()
	mProof := solidityProof.MarshalForSolidityVerifier()
	wireBlockNum := utils.EVMWordUint64(p.BlockNum)
	rl := copy(response[:], append(mProof[:], wireBlockNum...))
	if rl != OnChainResponseLength {
		return MarshaledOnChainResponse{}, errors.Errorf(
			"wrong length for response to VRFCoordinator")
	}
	return response, nil
}

// UnmarshalProofResponse returns the ProofResponse represented by the bytes in m
func UnmarshalProofResponse(m MarshaledOnChainResponse) (*ProofResponse, error) {
	blockNum := common.BytesToHash(m[ProofLength : ProofLength+32]).Big().Uint64()
	proof, err := UnmarshalSolidityProof(m[:ProofLength])
	if err != nil {
		return nil, errors.Wrap(err, "while parsing ProofResponse")
	}
	preSeed, err := BigToSeed(proof.Seed)
	if err != nil {
		return nil, errors.Wrap(err, "while converting seed to bytes representation")
	}
	return &ProofResponse{P: proof, PreSeed: preSeed, BlockNum: blockNum}, nil
}

// CryptoProof returns the proof implied by p, with the correct seed
func (p ProofResponse) CryptoProof(s PreSeedData) (vrfkey.Proof, error) {
	proof := p.P // Copy P, which has wrong seed value
	proof.Seed = FinalSeed(s)
	valid, err := proof.VerifyVRFProof()
	if err != nil {
		return vrfkey.Proof{}, errors.Wrap(err,
			"could not validate proof implied by on-chain response")
	}
	if !valid {
		return vrfkey.Proof{}, errors.Errorf(
			"proof implied by on-chain response is invalid")
	}
	return proof, nil
}

func GenerateProofResponseFromProof(proof vrfkey.Proof, s PreSeedData) (MarshaledOnChainResponse, error) {
	p := ProofResponse{P: proof, PreSeed: s.PreSeed, BlockNum: s.BlockNum}
	rv, err := p.MarshalForVRFCoordinator()
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return rv, nil
}

func GenerateProofResponseFromProofV2(p vrfkey.Proof, s PreSeedDataV2) (vrf_coordinator_v2.VRFProof, vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment, error) {
	var proof vrf_coordinator_v2.VRFProof
	var rc vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment
	solidityProof, err := SolidityPrecalculations(&p)
	if err != nil {
		return proof, rc, errors.Wrap(err,
			"while marshaling proof for VRFCoordinatorV2")
	}
	solidityProof.P.Seed = common.BytesToHash(s.PreSeed[:]).Big()
	x, y := secp256k1.Coordinates(solidityProof.P.PublicKey)
	gx, gy := secp256k1.Coordinates(solidityProof.P.Gamma)
	cgx, cgy := secp256k1.Coordinates(solidityProof.CGammaWitness)
	shx, shy := secp256k1.Coordinates(solidityProof.SHashWitness)
	return vrf_coordinator_v2.VRFProof{
			Pk:            [2]*big.Int{x, y},
			Gamma:         [2]*big.Int{gx, gy},
			C:             solidityProof.P.C,
			S:             solidityProof.P.S,
			Seed:          common.BytesToHash(s.PreSeed[:]).Big(),
			UWitness:      solidityProof.UWitness,
			CGammaWitness: [2]*big.Int{cgx, cgy},
			SHashWitness:  [2]*big.Int{shx, shy},
			ZInv:          solidityProof.ZInv,
		}, vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment{
			BlockNum:         s.BlockNum,
			SubId:            s.SubId,
			CallbackGasLimit: s.CallbackGasLimit,
			NumWords:         s.NumWords,
			Sender:           s.Sender,
		}, nil
}

func GenerateProofResponseFromProofV2Plus(
	p vrfkey.Proof,
	s PreSeedDataV2Plus) (
	vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof,
	vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment,
	error) {
	var proof vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof
	var rc vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment
	solidityProof, err := SolidityPrecalculations(&p)
	if err != nil {
		return proof, rc, errors.Wrap(err,
			"while marshaling proof for VRFCoordinatorV2Plus")
	}
	solidityProof.P.Seed = common.BytesToHash(s.PreSeed[:]).Big()
	x, y := secp256k1.Coordinates(solidityProof.P.PublicKey)
	gx, gy := secp256k1.Coordinates(solidityProof.P.Gamma)
	cgx, cgy := secp256k1.Coordinates(solidityProof.CGammaWitness)
	shx, shy := secp256k1.Coordinates(solidityProof.SHashWitness)
	return vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof{
			Pk:            [2]*big.Int{x, y},
			Gamma:         [2]*big.Int{gx, gy},
			C:             solidityProof.P.C,
			S:             solidityProof.P.S,
			Seed:          common.BytesToHash(s.PreSeed[:]).Big(),
			UWitness:      solidityProof.UWitness,
			CGammaWitness: [2]*big.Int{cgx, cgy},
			SHashWitness:  [2]*big.Int{shx, shy},
			ZInv:          solidityProof.ZInv,
		}, vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment{
			BlockNum:         s.BlockNum,
			SubId:            s.SubId,
			CallbackGasLimit: s.CallbackGasLimit,
			NumWords:         s.NumWords,
			Sender:           s.Sender,
			ExtraArgs:        s.ExtraArgs,
		}, nil
}

func GenerateProofResponse(keystore keystore.VRF, id string, s PreSeedData) (
	MarshaledOnChainResponse, error) {
	seed := FinalSeed(s)
	proof, err := keystore.GenerateProof(id, seed)
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return GenerateProofResponseFromProof(proof, s)
}

func GenerateProofResponseV2(keystore keystore.VRF, id string, s PreSeedDataV2) (
	vrf_coordinator_v2.VRFProof, vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment, error) {
	seedHashMsg := append(s.PreSeed[:], s.BlockHash.Bytes()...)
	seed := utils.MustHash(string(seedHashMsg)).Big()
	proof, err := keystore.GenerateProof(id, seed)
	if err != nil {
		return vrf_coordinator_v2.VRFProof{}, vrf_coordinator_v2.VRFCoordinatorV2RequestCommitment{}, err
	}
	return GenerateProofResponseFromProofV2(proof, s)
}

func GenerateProofResponseV2Plus(keystore keystore.VRF, id string, s PreSeedDataV2Plus) (
	vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof, vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment, error) {
	seedHashMsg := append(s.PreSeed[:], s.BlockHash.Bytes()...)
	seed := utils.MustHash(string(seedHashMsg)).Big()
	proof, err := keystore.GenerateProof(id, seed)
	if err != nil {
		return vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalProof{}, vrf_coordinator_v2plus_interface.IVRFCoordinatorV2PlusInternalRequestCommitment{}, err
	}
	return GenerateProofResponseFromProofV2Plus(proof, s)
}
