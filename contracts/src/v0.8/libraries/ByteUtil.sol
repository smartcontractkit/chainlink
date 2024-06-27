// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

/*
 * @title ByteUtil
 * @author Michael Fletcher
 * @notice Byte utility functions for efficiently parsing and manipulating packed byte data
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
   * @dev Reads a uint192 from a position within a byte array.
   * @param data Byte array to read from.
   * @param offset Position to start reading from.
   * @return result The uint192 read from the byte array.
   */
  function readUint192(bytes memory data, uint256 offset) internal pure returns (uint256 result) {
    //bounds check
    if (offset + 24 > data.length) revert MalformedData();

    assembly {
      //load 32 byte word accounting for 32 bit length and offset
      result := mload(add(add(data, 32), offset))
      //shift the result right 64 bits
      result := shr(64, result)
    }
  }

  /**
   * @dev Reads a uint32 from a position within a byte array.
   * @param data Byte array to read from.
   * @param offset Position to start reading from.
   * @return result The uint32 read from the byte array.
   */
  function readUint32(bytes memory data, uint256 offset) internal pure returns (uint256 result) {
    //bounds check
    if (offset + 4 > data.length) revert MalformedData();

    assembly {
      //load 32 byte word accounting for 32 bit length and offset
      result := mload(add(add(data, 32), offset))
      //shift the result right 224 bits
      result := shr(224, result)
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
      result := shr(96, word)
    }
  }
}
