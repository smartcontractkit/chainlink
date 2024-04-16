// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice This library contains various token pool functions to aid constructing the return data.
library Pool {
  error InvalidTag(bytes4 tag);
  error MalformedPoolReturnData(bytes data);

  // bytes4(keccak256("POOL_RETURN_DATA_V1_TAG"))
  bytes4 public constant POOL_RETURN_DATA_V1_TAG = 0x179fa694;

  struct PoolReturnDataV1 {
    bytes destPoolAddress;
    bytes destPoolData;
  }

  ///  @notice Generates the return dataV1 for the burnOrMint pool call.
  ///  @param remotePoolAddress The address of the remote pool.
  ///  @param destPoolData The data to send to the remote pool.
  ///  @return The return data for the burnOrMint pool call.
  function _generatePoolReturnDataV1(
    bytes memory remotePoolAddress,
    bytes memory destPoolData
  ) internal pure returns (bytes memory) {
    return abi.encodeWithSelector(
      POOL_RETURN_DATA_V1_TAG, PoolReturnDataV1({destPoolAddress: remotePoolAddress, destPoolData: destPoolData})
    );
  }

  /// @notice Decodes the PoolReturnDataV1 struct from the given data. Also checks if the tag is correct.
  /// @param encodedData The data to decode.
  /// @dev Can revert. Since this is only used on the sending side, this is acceptable.
  /// @return The decoded PoolReturnDataV1 struct.
  function _decodePoolReturnDataV1(bytes memory encodedData) internal pure returns (PoolReturnDataV1 memory) {
    if (bytes4(encodedData) != POOL_RETURN_DATA_V1_TAG) {
      revert InvalidTag(bytes4(encodedData));
    }

    return abi.decode(_removeFirstFourBytes(encodedData), (PoolReturnDataV1));
  }

  uint256 private constant SELECTOR_LENGTH = 4;

  /// @notice Removes the first four bytes from the given bytes. This can be used to undo `encodeWithSelector`.
  /// @param _bytes The bytes to remove the first four bytes from.
  /// @dev Can revert if the given bytes are less than four bytes long.
  /// @return trimmedBytes The bytes with the first four bytes removed.
  function _removeFirstFourBytes(bytes memory _bytes) internal pure returns (bytes memory trimmedBytes) {
    if (_bytes.length < SELECTOR_LENGTH) {
      revert MalformedPoolReturnData(_bytes);
    }

    uint256 newSliceLength = _bytes.length - SELECTOR_LENGTH;
    assembly {
      // Get a location of some free memory and store it in trimmedBytes as Solidity does for memory variables.
      trimmedBytes := mload(0x40)

      // Calculate length mod 32 to handle slices that are not a multiple of 32 in size.
      let lengthmod := and(newSliceLength, 31)

      // trimmedBytes will have the following format in memory: <length><data>
      // When copying data we will offset the start forward to avoid allocating additional memory
      // Therefore part of the length area will be written, but this will be overwritten later anyways.
      // In case no offset is require, the start is set to the data region (0x20 from the trimmedBytes)
      // mc will be used to keep track where to copy the data to.
      let mc := add(add(trimmedBytes, lengthmod), mul(0x20, iszero(lengthmod)))
      let end := add(mc, newSliceLength)

      for {
        // Same logic as for mc is applied and additionally the start offset specified for the method is added
        let cc := add(add(add(_bytes, lengthmod), mul(0x20, iszero(lengthmod))), SELECTOR_LENGTH)
      } lt(mc, end) {
        // increase `mc` and `cc` to read the next word from memory
        mc := add(mc, 0x20)
        cc := add(cc, 0x20)
      } {
        // Copy the data from source (cc location) to the slice data (mc location)
        mstore(mc, mload(cc))
      }

      // Store the length of the slice. This will overwrite any partial data that
      // was copied when having slices that are not a multiple of 32.
      mstore(trimmedBytes, newSliceLength)

      // update free-memory pointer
      // allocating the array padded to 32 bytes like the compiler does now
      // To set the used memory as a multiple of 32, add 31 to the actual memory usage (mc)
      // and remove the modulo 32 (the `and` with `not(31)`)
      mstore(0x40, and(add(mc, 31), not(31)))
    }

    return trimmedBytes;
  }
}
