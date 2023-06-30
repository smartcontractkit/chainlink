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
     * @dev Reads n bytes from a position within a byte array
     * @param data Byte array to read from.
     * @param offset Position to start reading from.
     * @param n Number of bytes to read.
     * @return result The bytes read from the byte array.
     */
    function readBytes(bytes memory data, uint256 offset, uint256 n) internal pure returns (bytes memory) {
        //bounds check
        if(offset + n > data.length) revert MalformedData();

        //allocate n bytes of memory
        bytes memory result = new bytes(n);

        assembly {
            //first 32 bytes is length, so offset by 32 + offset
            let dataPtr := add(add(data, 32), offset)

            //first 32 bytes of local array is length, so offset by 32
            let resultPtr := add(result, 32)

            //calculate the end pointer and the remaining bytes as saving in 32 byte chunks is more efficient than using mstore8, i.e 68 bytes would be saved as, 32,32,4, and endPointer would be 64
            let remainingBytes := mod(n, 32)
            let endPtr := add(dataPtr, sub(n, remainingBytes))

            //copy 32-byte chunks at a time from the offset
            for { } lt(dataPtr, endPtr) { dataPtr := add(dataPtr, 32) resultPtr := add(resultPtr, 32) } {
                mstore(resultPtr, mload(dataPtr))
            }

            //handle the remaining bytes, which would be 4 bytes in the example above
            if gt(remainingBytes, 0) {
                //number of bytes in the mask
                let maskSizeBytes := sub(32, remainingBytes)
                //convert number of bytes in the mask to bits
                let maskSizeBits := mul(maskSizeBytes, 8)
                //create the mask by raising it to the power of the number of bits and subtracting 1
                let mask := sub(exp(2, maskSizeBits), 1)
                //the current mask would apply to the least significant bits (i.e 0x00FF), when we want it to apply to the most significant bits (i.e 0xFF00) due to it being a dynamic array, so we need to negate it
                let negatedMask := not(mask)
                //apply the mask to the data e.g if we have 2 bytes of data loaded, but only want 1, then the mask would be 0xFF00 and the data would be 0xFFFF, so we would do 0xFFFF & 0xFF00 = 0xFF, dropping the extra bits we don't want to load
                let maskedData := and(mload(dataPtr), negatedMask)
                //save the data
                mstore(resultPtr, maskedData)
            }
        }

        return result;
    }

    /**
     * @dev Reads a uint256 from a position within a byte array.
     * @param data Byte array to read from.
     * @param offset Position to start reading from.
     * @return result The uint256 read from the byte array.
     */
    function readUint256(bytes memory data, uint256 offset) internal pure returns (uint256 result) {
        //bounds check
        if(offset + 32 > data.length) revert MalformedData();

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
        if(offset + 20 > data.length) revert MalformedData();

        assembly {
            //load 32 byte word accounting for 32 bit length and offset
            let word := mload(add(add(data, 32), offset))
            //address is the last 20 bytes of the word, so shift right
            result := shr(mul(8, 12), word)
        }
    }

    /**
     * @dev Reads a uint32 from a position within a byte array.
     * @param data Byte array to read from.
     * @param offset Position to start reading from.
     * @return result The uint32 read from the byte array.
     */
    function readBytes32(bytes memory data, uint256 offset) internal pure returns (bytes32 result) {
        //bounds check
        if(offset + 32 > data.length) revert MalformedData();

        assembly {
            //load 32 byte word accounting for 32 bit length and offset
            result := mload(add(add(data, 32), offset))
        }
    }

}

