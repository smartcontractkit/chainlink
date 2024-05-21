// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract VRFLogEmitter {
  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint64 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    address indexed sender
  );
  event RandomWordsFulfilled(uint256 indexed requestId, uint256 outputSeed, uint96 payment, bool success);

  function emitRandomWordsRequested(
    bytes32 keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint64 subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    address sender
  ) public {
    emit RandomWordsRequested(
      keyHash,
      requestId,
      preSeed,
      subId,
      minimumRequestConfirmations,
      callbackGasLimit,
      numWords,
      sender
    );
  }

  function emitRandomWordsFulfilled(uint256 requestId, uint256 outputSeed, uint96 payment, bool success) public {
    emit RandomWordsFulfilled(requestId, outputSeed, payment, success);
  }
}
