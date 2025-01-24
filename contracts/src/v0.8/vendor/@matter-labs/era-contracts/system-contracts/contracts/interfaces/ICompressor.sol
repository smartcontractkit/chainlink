// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

// The bitmask by applying which to the compressed state diff metadata we retrieve its operation.
uint8 constant OPERATION_BITMASK = 7;
// The number of bits shifting the compressed state diff metadata by which we retrieve its length.
uint8 constant LENGTH_BITS_OFFSET = 3;
// The maximal length in bytes that an enumeration index can have.
uint8 constant MAX_ENUMERATION_INDEX_SIZE = 8;

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice The interface for the Compressor contract, responsible for verifying the correctness of
 * the compression of the state diffs and bytecodes.
 */
interface ICompressor {
    function publishCompressedBytecode(
        bytes calldata _bytecode,
        bytes calldata _rawCompressedData
    ) external returns (bytes32 bytecodeHash);

    function verifyCompressedStateDiffs(
        uint256 _numberOfStateDiffs,
        uint256 _enumerationIndexSize,
        bytes calldata _stateDiffs,
        bytes calldata _compressedStateDiffs
    ) external returns (bytes32 stateDiffHash);
}
