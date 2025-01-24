// SPDX-License-Identifier: MIT

pragma solidity 0.8.24;

import {ICompressor, OPERATION_BITMASK, LENGTH_BITS_OFFSET, MAX_ENUMERATION_INDEX_SIZE} from "./interfaces/ICompressor.sol";
import {SystemContractBase} from "./abstract/SystemContractBase.sol";
import {Utils} from "./libraries/Utils.sol";
import {UnsafeBytesCalldata} from "./libraries/UnsafeBytesCalldata.sol";
import {EfficientCall} from "./libraries/EfficientCall.sol";
import {L1_MESSENGER_CONTRACT, STATE_DIFF_ENTRY_SIZE, KNOWN_CODE_STORAGE_CONTRACT} from "./Constants.sol";
import {DerivedKeyNotEqualToCompressedValue, EncodedAndRealBytecodeChunkNotEqual, DictionaryDividedByEightNotGreaterThanEncodedDividedByTwo, EncodedLengthNotFourTimesSmallerThanOriginal, IndexOutOfBounds, IndexSizeError, UnsupportedOperation, CompressorInitialWritesProcessedNotEqual, CompressorEnumIndexNotEqual, StateDiffLengthMismatch, CompressionValueTransformError, CompressionValueAddError, CompressionValueSubError} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice Contract with code pertaining to compression for zkEVM; at the moment this is used for bytecode compression
 * and state diff compression validation.
 * @dev Every deployed bytecode/published state diffs in zkEVM should be publicly restorable from the L1 data availability.
 * For this reason, the user may request the sequencer to publish the original bytecode and mark it as known.
 * Or the user may compress the bytecode and publish it instead (fewer data onchain!). At the end of every L1 Batch
 * we publish pubdata, part of which contains the state diffs that occurred within the batch.
 */
contract Compressor is ICompressor, SystemContractBase {
    using UnsafeBytesCalldata for bytes;

    /// @notice Verify the compressed bytecode and publish it on the L1.
    /// @param _bytecode The original bytecode to be verified against.
    /// @param _rawCompressedData The compressed bytecode in a format of:
    ///    - 2 bytes: the length of the dictionary
    ///    - N bytes: the dictionary
    ///    - M bytes: the encoded data
    /// @return bytecodeHash The hash of the original bytecode.
    /// @dev The dictionary is a sequence of 8-byte chunks, each of them has the associated index.
    /// @dev The encoded data is a sequence of 2-byte chunks, each of them is an index of the dictionary.
    /// @dev The compression algorithm works as follows:
    ///     1. The original bytecode is split into 8-byte chunks.
    ///     Since the bytecode size is always a multiple of 32, this is always possible.
    ///     2. For each 8-byte chunk in the original bytecode:
    ///         * If the chunk is not already in the dictionary, it is added to the dictionary array.
    ///         * If the dictionary becomes overcrowded (2^16 + 1 elements), the compression process will fail.
    ///         * The 2-byte index of the chunk in the dictionary is added to the encoded data.
    /// @dev Currently, the method may be called only from the bootloader because the server is not ready to publish bytecodes
    /// in internal transactions. However, in the future, we will allow everyone to publish compressed bytecodes.
    /// @dev Read more about the compression: https://github.com/matter-labs/zksync-era/blob/main/docs/guides/advanced/compression.md
    function publishCompressedBytecode(
        bytes calldata _bytecode,
        bytes calldata _rawCompressedData
    ) external onlyCallFromBootloader returns (bytes32 bytecodeHash) {
        unchecked {
            (bytes calldata dictionary, bytes calldata encodedData) = _decodeRawBytecode(_rawCompressedData);

            if (encodedData.length * 4 != _bytecode.length) {
                revert EncodedLengthNotFourTimesSmallerThanOriginal();
            }

            if (dictionary.length / 8 > encodedData.length / 2) {
                revert DictionaryDividedByEightNotGreaterThanEncodedDividedByTwo();
            }

            // We disable this check because calldata array length is cheap.
            // solhint-disable-next-line gas-length-in-loops
            for (uint256 encodedDataPointer = 0; encodedDataPointer < encodedData.length; encodedDataPointer += 2) {
                uint256 indexOfEncodedChunk = uint256(encodedData.readUint16(encodedDataPointer)) * 8;
                if (indexOfEncodedChunk > dictionary.length - 1) {
                    revert IndexOutOfBounds();
                }

                uint64 encodedChunk = dictionary.readUint64(indexOfEncodedChunk);
                uint64 realChunk = _bytecode.readUint64(encodedDataPointer * 4);

                if (encodedChunk != realChunk) {
                    revert EncodedAndRealBytecodeChunkNotEqual(realChunk, encodedChunk);
                }
            }
        }

        bytecodeHash = Utils.hashL2Bytecode(_bytecode);
        L1_MESSENGER_CONTRACT.sendToL1(_rawCompressedData);
        KNOWN_CODE_STORAGE_CONTRACT.markBytecodeAsPublished(bytecodeHash);
    }

    /// @notice Verifies that the compression of state diffs has been done correctly for the {_stateDiffs} param.
    /// @param _numberOfStateDiffs The number of state diffs being checked.
    /// @param _enumerationIndexSize Number of bytes used to represent an enumeration index for repeated writes.
    /// @param _stateDiffs Encoded full state diff structs. See the first dev comment below for encoding.
    /// @param _compressedStateDiffs The compressed state diffs
    /// @dev We don't verify that the size of {_stateDiffs} is equivalent to {_numberOfStateDiffs} * STATE_DIFF_ENTRY_SIZE since that check is
    ///      done within the L1Messenger calling contract.
    /// @return stateDiffHash Hash of the encoded (uncompressed) state diffs to be committed to via system log.
    /// @dev This check assumes that the ordering of state diffs are sorted by (address, key) for the encoded state diffs and
    ///      then the compressed are sorted the same but with all the initial writes coming before the repeated writes.
    /// @dev state diff:   [20bytes address][32bytes key][32bytes derived key][8bytes enum index][32bytes initial value][32bytes final value]
    /// @dev The compression format:
    ///     - 2 bytes: number of initial writes
    ///     - N bytes initial writes
    ///         - 32 bytes derived key
    ///         - 1 byte metadata:
    ///             - first 5 bits: length in bytes of compressed value
    ///             - last 3 bits: operation
    ///                 - 0 -> Nothing (32 bytes)
    ///                 - 1 -> Add
    ///                 - 2 -> Subtract
    ///                 - 3 -> Transform (< 32 bytes)
    ///         - Len Bytes: Compressed Value
    ///     - M bytes repeated writes
    ///         - {_enumerationIndexSize} bytes for enumeration index
    ///         - 1 byte metadata:
    ///             - first 5 bits: length in bytes of compressed value
    ///             - last 3 bits: operation
    ///                 - 0 -> Nothing (32 bytes)
    ///                 - 1 -> Add
    ///                 - 2 -> Subtract
    ///                 - 3 -> Transform (< 32 bytes)
    ///         - Len Bytes: Compressed Value
    function verifyCompressedStateDiffs(
        uint256 _numberOfStateDiffs,
        uint256 _enumerationIndexSize,
        bytes calldata _stateDiffs,
        bytes calldata _compressedStateDiffs
    ) external onlyCallFrom(address(L1_MESSENGER_CONTRACT)) returns (bytes32 stateDiffHash) {
        // We do not enforce the operator to use the optimal, i.e. the minimally possible _enumerationIndexSize.
        // We do enforce however, that the _enumerationIndexSize is not larger than 8 bytes long, which is the
        // maximal ever possible size for enumeration index.
        if (_enumerationIndexSize > MAX_ENUMERATION_INDEX_SIZE) {
            revert IndexSizeError();
        }

        uint256 numberOfInitialWrites = uint256(_compressedStateDiffs.readUint16(0));

        uint256 stateDiffPtr = 2;
        uint256 numInitialWritesProcessed = 0;

        // Process initial writes
        for (uint256 i = 0; i < _numberOfStateDiffs * STATE_DIFF_ENTRY_SIZE; i += STATE_DIFF_ENTRY_SIZE) {
            bytes calldata stateDiff = _stateDiffs[i:i + STATE_DIFF_ENTRY_SIZE];
            uint64 enumIndex = stateDiff.readUint64(84);
            if (enumIndex != 0) {
                // It is a repeated write, so we skip it.
                continue;
            }

            ++numInitialWritesProcessed;

            bytes32 derivedKey = stateDiff.readBytes32(52);
            uint256 initValue = stateDiff.readUint256(92);
            uint256 finalValue = stateDiff.readUint256(124);
            bytes32 compressedDerivedKey = _compressedStateDiffs.readBytes32(stateDiffPtr);
            if (derivedKey != compressedDerivedKey) {
                revert DerivedKeyNotEqualToCompressedValue(derivedKey, compressedDerivedKey);
            }
            stateDiffPtr += 32;

            uint8 metadata = uint8(bytes1(_compressedStateDiffs[stateDiffPtr]));
            ++stateDiffPtr;
            uint8 operation = metadata & OPERATION_BITMASK;
            uint8 len = operation == 0 ? 32 : metadata >> LENGTH_BITS_OFFSET;
            _verifyValueCompression(
                initValue,
                finalValue,
                operation,
                _compressedStateDiffs[stateDiffPtr:stateDiffPtr + len]
            );
            stateDiffPtr += len;
        }

        if (numInitialWritesProcessed != numberOfInitialWrites) {
            revert CompressorInitialWritesProcessedNotEqual(numberOfInitialWrites, numInitialWritesProcessed);
        }

        // Process repeated writes
        for (uint256 i = 0; i < _numberOfStateDiffs * STATE_DIFF_ENTRY_SIZE; i += STATE_DIFF_ENTRY_SIZE) {
            bytes calldata stateDiff = _stateDiffs[i:i + STATE_DIFF_ENTRY_SIZE];
            uint64 enumIndex = stateDiff.readUint64(84);
            if (enumIndex == 0) {
                continue;
            }

            uint256 initValue = stateDiff.readUint256(92);
            uint256 finalValue = stateDiff.readUint256(124);
            uint256 compressedEnumIndex = _sliceToUint256(
                _compressedStateDiffs[stateDiffPtr:stateDiffPtr + _enumerationIndexSize]
            );
            if (enumIndex != compressedEnumIndex) {
                revert CompressorEnumIndexNotEqual(enumIndex, compressedEnumIndex);
            }
            stateDiffPtr += _enumerationIndexSize;

            uint8 metadata = uint8(bytes1(_compressedStateDiffs[stateDiffPtr]));
            ++stateDiffPtr;
            uint8 operation = metadata & OPERATION_BITMASK;
            uint8 len = operation == 0 ? 32 : metadata >> LENGTH_BITS_OFFSET;
            _verifyValueCompression(
                initValue,
                finalValue,
                operation,
                _compressedStateDiffs[stateDiffPtr:stateDiffPtr + len]
            );
            stateDiffPtr += len;
        }

        if (stateDiffPtr != _compressedStateDiffs.length) {
            revert StateDiffLengthMismatch();
        }

        stateDiffHash = EfficientCall.keccak(_stateDiffs);
    }

    /// @notice Decode the raw compressed data into the dictionary and the encoded data.
    /// @param _rawCompressedData The compressed bytecode in a format of:
    ///    - 2 bytes: the bytes length of the dictionary
    ///    - N bytes: the dictionary
    ///    - M bytes: the encoded data
    function _decodeRawBytecode(
        bytes calldata _rawCompressedData
    ) internal pure returns (bytes calldata dictionary, bytes calldata encodedData) {
        unchecked {
            // The dictionary length can't be more than 2^16, so it fits into 2 bytes.
            uint256 dictionaryLen = uint256(_rawCompressedData.readUint16(0));
            dictionary = _rawCompressedData[2:2 + dictionaryLen * 8];
            encodedData = _rawCompressedData[2 + dictionaryLen * 8:];
        }
    }

    /// @notice Verify value compression was done correct given initial value, final value, operation, and compressed value
    /// @param _initialValue Previous value of key/enumeration index.
    /// @param _finalValue Updated value of key/enumeration index.
    /// @param _operation The operation that was performed on value.
    /// @param _compressedValue The slice of calldata with compressed value either representing the final
    /// value or difference between initial and final value. It should be of arbitrary length less than or equal to 32 bytes.
    /// @dev It is the responsibility of the caller of this function to ensure that the `_compressedValue` has length no longer than 32 bytes.
    /// @dev Operation id mapping:
    /// 0 -> Nothing (32 bytes)
    /// 1 -> Add
    /// 2 -> Subtract
    /// 3 -> Transform (< 32 bytes)
    function _verifyValueCompression(
        uint256 _initialValue,
        uint256 _finalValue,
        uint256 _operation,
        bytes calldata _compressedValue
    ) internal pure {
        uint256 convertedValue = _sliceToUint256(_compressedValue);

        unchecked {
            if (_operation == 0 || _operation == 3) {
                if (convertedValue != _finalValue) {
                    revert CompressionValueTransformError(_finalValue, convertedValue);
                }
            } else if (_operation == 1) {
                if (_initialValue + convertedValue != _finalValue) {
                    revert CompressionValueAddError(_finalValue, _initialValue + convertedValue);
                }
            } else if (_operation == 2) {
                if (_initialValue - convertedValue != _finalValue) {
                    revert CompressionValueSubError(_finalValue, _initialValue - convertedValue);
                }
            } else {
                revert UnsupportedOperation();
            }
        }
    }

    /// @notice Converts a calldata slice into uint256. It is the responsibility of the caller to ensure that
    /// the _calldataSlice has length no longer than 32 bytes
    /// @param _calldataSlice The calldata slice to convert to uint256
    /// @return number The uint256 representation of the calldata slice
    function _sliceToUint256(bytes calldata _calldataSlice) internal pure returns (uint256 number) {
        number = uint256(bytes32(_calldataSlice));
        number >>= (256 - (_calldataSlice.length * 8));
    }
}
