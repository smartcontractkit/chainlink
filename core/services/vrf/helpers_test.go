package vrf

import "math/big"

func GenerateProofWithNonce(secretKey, seed, nonce *big.Int) (Proof, error) {
	return generateProofWithNonce(secretKey, seed, nonce)
}

func GenerateProofResponseWithNonce(secretKey *big.Int, s PreSeedData,
	nonce *big.Int) (
	MarshaledOnChainResponse, error) {
	seed := FinalSeed(s)
	proof, err := generateProofWithNonce(secretKey, seed, nonce)
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	p := ProofResponse{P: proof, PreSeed: s.PreSeed, BlockNum: s.BlockNum}
	rv, err := p.MarshalForVRFCoordinator()
	if err != nil {
		return MarshaledOnChainResponse{}, err
	}
	return rv, nil
}
