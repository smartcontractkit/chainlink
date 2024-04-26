// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {CBOR} from "../../vendor/solidity-cborutils/v2.0.0/CBOR.sol";

//struct UpkeepOffchainConfig {
//  uint256 maxGasPrice;
//}

contract CBOREncoder {
  using CBOR for CBOR.CBORBuffer;

  bytes public data;

  /// @notice encodes a max gas price to CBOR encoded bytes
  /// @param maxGasPrice The max gas price
  /// @return CBOR encoded bytes and the struct depth
  function encode(uint256 maxGasPrice, uint256 capacity) external pure returns (bytes memory, uint256) {
    CBOR.CBORBuffer memory buffer = CBOR.create(capacity);
    buffer.writeString("maxGasPrice");
    buffer.writeUInt256(maxGasPrice);
    return (buffer.buf.buf, buffer.depth);
  }

  /// @notice encodes a max gas price to CBOR encoded bytes
  /// @param maxGasPrice The max gas price
  function encodeWrite(uint256 maxGasPrice, uint256 capacity) external {
    CBOR.CBORBuffer memory buffer = CBOR.create(capacity);
    buffer.writeString("maxGasPrice");
    buffer.writeUInt256(maxGasPrice);
    data = buffer.buf.buf;
  }
}
