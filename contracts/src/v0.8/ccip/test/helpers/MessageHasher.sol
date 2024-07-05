// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";

/// @notice MessageHasher is a contract that utility functions to hash an Any2EVMRampMessage
/// and encode various preimages for the final hash of the message.
contract MessageHasher {
  function hash(Internal.Any2EVMRampMessage memory message, bytes memory onRamp) public pure returns (bytes32) {
    return Internal._hash(message, onRamp);
  }

  function encodeTokenAmountsHashPreimage(Client.EVMTokenAmount[] memory tokenAmounts)
    public
    pure
    returns (bytes memory)
  {
    return abi.encode(tokenAmounts);
  }

  function encodeSourceTokenDataHashPreimage(bytes[] memory sourceTokenData) public pure returns (bytes memory) {
    return abi.encode(sourceTokenData);
  }

  function encodeMetadataHashPreimage(
    bytes32 any2EVMMessageHash,
    uint64 sourceChainSelector,
    uint64 destChainSelector,
    bytes memory onRamp
  ) public pure returns (bytes memory) {
    return abi.encode(any2EVMMessageHash, sourceChainSelector, destChainSelector, onRamp);
  }

  function encodeFixedSizeFieldsHashPreimage(
    bytes32 messageId,
    bytes memory sender,
    address receiver,
    uint64 sequenceNumber,
    uint256 gasLimit,
    uint64 nonce
  ) public pure returns (bytes memory) {
    return abi.encode(messageId, sender, receiver, sequenceNumber, gasLimit, nonce);
  }

  function encodeFinalHashPreimage(
    bytes32 leafDomainSeparator,
    bytes32 implicitMetadataHash,
    bytes32 fixedSizeFieldsHash,
    bytes32 dataHash,
    bytes32 tokenAmountsHash,
    bytes32 sourceTokenDataHash
  ) public pure returns (bytes memory) {
    return abi.encode(
      leafDomainSeparator, implicitMetadataHash, fixedSizeFieldsHash, dataHash, tokenAmountsHash, sourceTokenDataHash
    );
  }
}
