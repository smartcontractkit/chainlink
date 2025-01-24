// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @dev The library provides a set of functions that help read data from calldata bytes.
 * @dev Each of the functions accepts the `bytes calldata` and the offset where data should be read and returns a value of a certain type.
 *
 * @dev WARNING!
 * 1) Functions don't check the length of the bytes array, so it can go out of bounds.
 * The user of the library must check for bytes length before using any functions from the library!
 *
 * 2) Read variables are not cleaned up - https://docs.soliditylang.org/en/v0.8.16/internals/variable_cleanup.html.
 * Using data in inline assembly can lead to unexpected behavior!
 */
library UnsafeBytesCalldata {
    function readUint16(bytes calldata _bytes, uint256 _start) internal pure returns (uint16 result) {
        assembly {
            let offset := sub(_bytes.offset, 30)
            result := calldataload(add(offset, _start))
        }
    }

    function readUint32(bytes calldata _bytes, uint256 _start) internal pure returns (uint32 result) {
        assembly {
            let offset := sub(_bytes.offset, 28)
            result := calldataload(add(offset, _start))
        }
    }

    function readUint64(bytes calldata _bytes, uint256 _start) internal pure returns (uint64 result) {
        assembly {
            let offset := sub(_bytes.offset, 24)
            result := calldataload(add(offset, _start))
        }
    }

    function readBytes32(bytes calldata _bytes, uint256 _start) internal pure returns (bytes32 result) {
        assembly {
            result := calldataload(add(_bytes.offset, _start))
        }
    }

    function readUint256(bytes calldata _bytes, uint256 _start) internal pure returns (uint256 result) {
        assembly {
            result := calldataload(add(_bytes.offset, _start))
        }
    }
}
