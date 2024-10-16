// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.4;

import {Client} from "../../libraries/Client.sol";
import {Internal} from "../../libraries/Internal.sol";

/// @notice MessageHasher is a contract that utility functions to hash an Any2EVMRampMessage
/// and encode various preimages for the final hash of the message.
/// @dev This is only deployed in tests and is not part of the production contracts.
contract MessageHasher {
  function hash(Internal.Any2EVMRampMessage memory message, bytes memory onRamp) public pure returns (bytes32) {
    return Internal._hash(
      message,
      keccak256(
        abi.encode(
          Internal.ANY_2_EVM_MESSAGE_HASH,
          message.header.sourceChainSelector,
          message.header.destChainSelector,
          keccak256(onRamp)
        )
      )
    );
  }

  function encodeAny2EVMTokenAmountsHashPreimage(
    Internal.Any2EVMTokenTransfer[] memory tokenAmounts
  ) public pure returns (bytes memory) {
    return abi.encode(tokenAmounts);
  }

  function encodeEVM2AnyTokenAmountsHashPreimage(
    Internal.EVM2AnyTokenTransfer[] memory tokenAmount
  ) public pure returns (bytes memory) {
    return abi.encode(tokenAmount);
  }

  function encodeMetadataHashPreimage(
    bytes32 any2EVMMessageHash,
    uint64 sourceChainSelector,
    uint64 destChainSelector,
    bytes32 onRampHash
  ) public pure returns (bytes memory) {
    return abi.encode(any2EVMMessageHash, sourceChainSelector, destChainSelector, onRampHash);
  }

  function encodeFixedSizeFieldsHashPreimage(
    bytes32 messageId,
    address receiver,
    uint64 sequenceNumber,
    uint256 gasLimit,
    uint64 nonce
  ) public pure returns (bytes memory) {
    return abi.encode(messageId, receiver, sequenceNumber, gasLimit, nonce);
  }

  function encodeFinalHashPreimage(
    bytes32 leafDomainSeparator,
    bytes32 metaDataHash,
    bytes32 fixedSizeFieldsHash,
    bytes32 senderHash,
    bytes32 dataHash,
    bytes32 tokenAmountsHash
  ) public pure returns (bytes memory) {
    return abi.encode(leafDomainSeparator, metaDataHash, fixedSizeFieldsHash, senderHash, dataHash, tokenAmountsHash);
  }

  function encodeEVMExtraArgsV1(
    Client.EVMExtraArgsV1 memory extraArgs
  ) public pure returns (bytes memory) {
    return Client._argsToBytes(extraArgs);
  }

  function encodeEVMExtraArgsV2(
    Client.EVMExtraArgsV2 memory extraArgs
  ) public pure returns (bytes memory) {
    return Client._argsToBytes(extraArgs);
  }

  function decodeEVMExtraArgsV1(
    uint256 gasLimit
  ) public pure returns (Client.EVMExtraArgsV1 memory) {
    return Client.EVMExtraArgsV1(gasLimit);
  }

  function decodeEVMExtraArgsV2(
    uint256 gasLimit,
    bool allowOutOfOrderExecution
  ) public pure returns (Client.EVMExtraArgsV2 memory) {
    return Client.EVMExtraArgsV2(gasLimit, allowOutOfOrderExecution);
  }
}
