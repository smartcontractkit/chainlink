// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./IVRFCoordinatorV2Plus.sol";

// IVRFCoordinatorV2PlusInternal is the interface used by chainlink node
// for backwards-compatibility for future versions of V2Plus
// This interface should not be used by consumer conracts
interface IVRFCoordinatorV2PlusInternal is IVRFCoordinatorV2Plus {
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint256 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bytes extraArgs,
    address indexed sender
  );

  event RandomWordsFulfilled(
    uint256 indexed requestId,
    uint256 outputSeed,
    uint256 indexed subId,
    uint96 payment,
    bool success
  );

  struct RequestCommitment {
    uint64 blockNum;
    uint256 subId;
    uint32 callbackGasLimit;
    uint32 numWords;
    address sender;
    bytes extraArgs;
  }

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

  function s_requestCommitments(uint256 requestID) external view returns (bytes32);

  function fulfillRandomWords(Proof memory proof, RequestCommitment memory rc) external returns (uint96);

  function LINK_NATIVE_FEED() external view returns (address);
}