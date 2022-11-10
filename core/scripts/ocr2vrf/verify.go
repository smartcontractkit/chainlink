package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/google"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
	"go.dedis.ch/kyber/v3/pairing"

	dkgContract "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func getDKGLatestConfigDetails(e helpers.Environment, dkgAddress string) dkgContract.LatestConfigDetails {
	dkg := newDKG(common.HexToAddress(dkgAddress), e.Ec)
	dkgConfig, err := dkg.LatestConfigDetails(nil)
	helpers.PanicErr(err)

	return dkgConfig
}

func getVRFLatestConfigDetails(e helpers.Environment, beaconAddress string) vrf_beacon.LatestConfigDetails {
	beacon := newVRFBeacon(common.HexToAddress(beaconAddress), e.Ec)
	beaconConfig, err := beacon.LatestConfigDetails(nil)
	helpers.PanicErr(err)

	return beaconConfig
}

func getDKGKeyData(e helpers.Environment, dkgAddress string, keyID, configDigest [32]byte) dkgContract.KeyDataStructKeyData {
	dkg := newDKG(common.HexToAddress(dkgAddress), e.Ec)
	keyData, err := dkg.GetKey(nil, keyID, configDigest)
	helpers.PanicErr(err)

	return keyData
}

func getKeyID(e helpers.Environment, beaconAddress string) [32]byte {
	beacon := newVRFBeacon(common.HexToAddress(beaconAddress), e.Ec)
	keyID, err := beacon.SKeyID(nil)
	helpers.PanicErr(err)
	return keyID
}

func getPublicKey(e helpers.Environment, dkgAddress string, keyID, configDigest [32]byte) kyber.Point {
	keyData := getDKGKeyData(e, dkgAddress, keyID, configDigest)
	kg := &altbn_128.G2{}
	pk := kg.Point()
	err := pk.UnmarshalBinary(keyData.PublicKey)
	helpers.PanicErr(err)
	return pk
}

func getHashToCurveMessage(e helpers.Environment, height uint64, confDelay uint32, vrfConfigDigest [32]byte, pk kyber.Point) *altbn_128.HashProof {
	blockNumber := big.NewInt(0).SetUint64(height)
	block, err := e.Ec.BlockByNumber(context.Background(), blockNumber)
	helpers.PanicErr(err)
	b := ocr2vrftypes.Block{
		Height:            height,
		ConfirmationDelay: confDelay,
		Hash:              block.Hash(),
	}
	h := b.VRFHash(vrfConfigDigest, pk)
	return altbn_128.NewHashProof(h)
}

func getVRFSignature(e helpers.Environment, coordinatorAddress string, height, confDelay, searchWindow uint64) (proofG1X, proofG1Y *big.Int) {
	// get transmission logs from requested block to requested block + search window blocks
	// TODO: index transmission logs by height and confirmation delay to
	// make the FilterQuery call more efficient
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(0).SetUint64(height),
		ToBlock:   big.NewInt(0).SetUint64(height + searchWindow),
		Addresses: []common.Address{
			common.HexToAddress(coordinatorAddress),
		},
		Topics: [][]common.Hash{
			{
				vrf_coordinator.VRFCoordinatorOutputsServed{}.Topic(),
			},
		},
	}
	logs, err := e.Ec.FilterLogs(context.Background(), query)
	helpers.PanicErr(err)

	coordinator := newVRFCoordinator(common.HexToAddress(coordinatorAddress), e.Ec)
	for _, log := range logs {
		t, err := coordinator.ParseOutputsServed(log)
		helpers.PanicErr(err)
		for _, o := range t.OutputsServed {
			if o.ConfirmationDelay.Uint64() == confDelay && o.Height == height {
				proofG1X = o.ProofG1X
				proofG1Y = o.ProofG1Y
			}
		}
	}
	return
}

func verifyBeaconRandomness(e helpers.Environment, dkgAddress, beaconAddress string, coordinatorAddress string, height, confDelay, searchWindow uint64) bool {
	dkgConfig := getDKGLatestConfigDetails(e, dkgAddress)
	vrfConfig := getVRFLatestConfigDetails(e, beaconAddress)
	keyID := getKeyID(e, beaconAddress)
	pk := getPublicKey(e, dkgAddress, keyID, dkgConfig.ConfigDigest)
	h := getHashToCurveMessage(e, height, uint32(confDelay), vrfConfig.ConfigDigest, pk)
	hpoint := h.HashPoint
	negHpoint := g1.Point()
	negHpoint.Neg(hpoint)
	g2Base := g2.Point().Base()

	// get BLS signature for the given height and confirmation delay
	proofG1X, proofG1Y := getVRFSignature(e, coordinatorAddress, height, confDelay, searchWindow)
	if proofG1X.Cmp(big.NewInt(0)) == 0 || proofG1Y.Cmp(big.NewInt(0)) == 0 {
		panic("signature not found")
	}
	g1Proof, err := altbn_128.CoordinatesToG1(mod.NewInt(proofG1X, bn256.P), mod.NewInt(proofG1Y, bn256.P))
	helpers.PanicErr(err)

	// Perform verification of BLS signature is done using pairing function
	isValid := validateSignature(suite, hpoint, pk, g1Proof)
	fmt.Println("Verification Result: ", isValid)

	// Perform the same verification as above using precompiled contract 0x8
	// This should always result in same result as validateSignature()
	// signature is valid iff contract0x8(-b_x, -b_y, pk_x, pk_y, p_x, p_y, g2_x, g2_y) == 1
	input := make([]byte, 384)
	hb := altbn_128.LongMarshal(negHpoint)
	if len(hb) != 64 {
		panic("wrong length of hpoint")
	}
	copy(input[:64], hb[:])

	pkb, err := pk.MarshalBinary()
	helpers.PanicErr(err)
	if len(pkb) != 128 {
		panic("wrong length of public key")
	}
	copy(input[64:192], pkb)

	if len(proofG1X.Bytes()) != 32 {
		panic("wrong length of VRF signature x-coordinator")
	}
	if len(proofG1Y.Bytes()) != 32 {
		panic("wrong length of VRF signature y-coordinator")
	}
	copy(input[192:224], proofG1X.Bytes())
	copy(input[224:256], proofG1Y.Bytes())

	g2b, err := g2Base.MarshalBinary()
	helpers.PanicErr(err)
	if len(g2b) != 128 {
		panic("wrong length of altbn_128 base points")
	}
	copy(input[256:384], g2b)

	contract := vm.PrecompiledContractsByzantium[common.HexToAddress("0x8")]
	res, err := contract.Run(input)
	helpers.PanicErr(err)
	isValidPrecompiledContract := big.NewInt(0).SetBytes(res).Uint64() == 1
	fmt.Println("Verification Result Using Precompiled Contract 0x8: ", isValidPrecompiledContract)

	if isValid && isValidPrecompiledContract {
		return true
	}
	return false
}

func validateSignature(p pairing.Suite, msg, publicKey, signature kyber.Point) bool {
	return p.Pair(msg, publicKey).Equal(p.Pair(signature, p.G2().Point().Base()))
}
