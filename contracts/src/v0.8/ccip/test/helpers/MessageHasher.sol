// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {Internal} from "../../libraries/Internal.sol";

// MessageHasher is a contract that provides a function to hash an EVM2EVMMessage.
contract MessageHasher {
  function hash(Internal.EVM2EVMMessage memory msg, bytes32 metadataHash) public pure returns (bytes32) {
    return Internal._hash(msg, metadataHash);
  }
}
