// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

/*
 * @title ByteUtil
 * @author Michael Fletcher
 * @notice Byte utility functions for efficiently parsing and manipulating byte data
 */
library ByteUtil {
  // Error message when an offset is out of bounds
  error MalformedData();

  /**
   * @dev Reads a uint256 from a position within a byte array.
   * @param data Byte array to read from.
   * @param offset Position to start reading from.
   * @return result The uint256 read from the byte array.
   */
  function readUint256(bytes memory data, uint256 offset) internal pure returns (uint256 result) {
    //bounds check
    if (offset + 32 > data.length) revert MalformedData();

    assembly {
      //load 32 byte word accounting for 32 bit length and offset
      result := mload(add(add(data, 32), offset))
    }
  }

  /**
   * @dev Reads an address from a position within a byte array.
   * @param data Byte array to read from.
   * @param offset Position to start reading from.
   * @return result The uint32 read from the byte array.
   */
  function readAddress(bytes memory data, uint256 offset) internal pure returns (address result) {
    //bounds check
    if (offset + 20 > data.length) revert MalformedData();

    assembly {
      //load 32 byte word accounting for 32 bit length and offset
      let word := mload(add(add(data, 32), offset))
      //address is the last 20 bytes of the word, so shift right
      result := shr(mul(8, 12), word)
    }
  }
}
