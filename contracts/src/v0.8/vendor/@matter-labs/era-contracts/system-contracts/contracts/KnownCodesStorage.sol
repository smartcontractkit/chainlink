// SPDX-License-Identifier: MIT

pragma solidity 0.8.24;

import {IKnownCodesStorage} from "./interfaces/IKnownCodesStorage.sol";
import {SystemContractBase} from "./abstract/SystemContractBase.sol";
import {Utils} from "./libraries/Utils.sol";
import {COMPRESSOR_CONTRACT, L1_MESSENGER_CONTRACT} from "./Constants.sol";
import {Unauthorized, MalformedBytecode, BytecodeError} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice The storage of this contract will basically serve as a mapping for the known code hashes.
 * @dev Code hash is not strictly a hash, it's a structure where the first byte denotes the version of the hash,
 * the second byte denotes whether the contract is constructed, and the next two bytes denote the length in 32-byte words.
 * And then the next 28 bytes is the truncated hash.
 */
contract KnownCodesStorage is IKnownCodesStorage, SystemContractBase {
    modifier onlyCompressor() {
        if (msg.sender != address(COMPRESSOR_CONTRACT)) {
            revert Unauthorized(msg.sender);
        }
        _;
    }

    /// @notice The method that is used by the bootloader to mark several bytecode hashes as known.
    /// @param _shouldSendToL1 Whether the bytecode should be sent on L1.
    /// @param _hashes Hashes of the bytecodes to be marked as known.
    function markFactoryDeps(bool _shouldSendToL1, bytes32[] calldata _hashes) external onlyCallFromBootloader {
        unchecked {
            uint256 hashesLen = _hashes.length;
            for (uint256 i = 0; i < hashesLen; ++i) {
                _markBytecodeAsPublished(_hashes[i], _shouldSendToL1);
            }
        }
    }

    /// @notice The method used to mark a single bytecode hash as known.
    /// @dev Only trusted contacts can call this method, currently only the bytecode compressor.
    /// @param _bytecodeHash The hash of the bytecode that is marked as known.
    function markBytecodeAsPublished(bytes32 _bytecodeHash) external onlyCompressor {
        _markBytecodeAsPublished(_bytecodeHash, false);
    }

    /// @notice The method used to mark a single bytecode hash as known
    /// @param _bytecodeHash The hash of the bytecode that is marked as known
    /// @param _shouldSendToL1 Whether the bytecode should be sent on L1
    function _markBytecodeAsPublished(bytes32 _bytecodeHash, bool _shouldSendToL1) internal {
        if (getMarker(_bytecodeHash) == 0) {
            _validateBytecode(_bytecodeHash);

            if (_shouldSendToL1) {
                L1_MESSENGER_CONTRACT.requestBytecodeL1Publication(_bytecodeHash);
            }

            // Save as known, to not resend the log to L1
            assembly {
                sstore(_bytecodeHash, 1)
            }

            emit MarkedAsKnown(_bytecodeHash, _shouldSendToL1);
        }
    }

    /// @notice Returns the marker stored for a bytecode hash. 1 means that the bytecode hash is known
    /// and can be used for deploying contracts. 0 otherwise.
    function getMarker(bytes32 _hash) public view override returns (uint256 marker) {
        assembly {
            marker := sload(_hash)
        }
    }

    /// @notice Validates the format of bytecodehash
    /// @dev zk-circuit accepts & handles only valid format of bytecode hash, other input has undefined behavior
    /// That's why we need to validate it
    function _validateBytecode(bytes32 _bytecodeHash) internal pure {
        uint8 version = uint8(_bytecodeHash[0]);
        if (version != 1 || _bytecodeHash[1] != bytes1(0)) {
            revert MalformedBytecode(BytecodeError.Version);
        }

        if (Utils.bytecodeLenInWords(_bytecodeHash) % 2 == 0) {
            revert MalformedBytecode(BytecodeError.NumberOfWords);
        }
    }
}
