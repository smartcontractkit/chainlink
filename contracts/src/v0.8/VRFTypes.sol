// SPDX-License-Identifier: MIT
pragma solidity 0.8.13;

/**
 * @title VRFTypes
 * @notice The VRFTypes library is a collection of types that is required to fulfill VRF requests
 * 	on-chain. They must be ABI-compatible with the types used by the coordinator contracts.
 */
library VRFTypes {
  // ABI-compatible with VRF.Proof.
  // This proof is used for VRF V2.
  struct Proof {
    uint256[2] pk;
    uint256[2] gamma;
    uint256 c;
    uint256 s;
    uint256 seed;
    address uWitness;
    uint256[2] cGammaWitness;
    uint256[2] sHashWitness;
    uint256 zInv;
  }

  // ABI-compatible with VRFCoordinatorV2.RequestCommitment.
  // This is only used for VRF V2.
  struct RequestCommitment {
    uint64 blockNum;
    uint64 subId;
    uint32 callbackGasLimit;
    uint32 numWords;
    address sender;
  }
}
