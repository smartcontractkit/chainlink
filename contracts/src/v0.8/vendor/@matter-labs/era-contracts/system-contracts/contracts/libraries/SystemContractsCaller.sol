// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

import {MSG_VALUE_SYSTEM_CONTRACT, MSG_VALUE_SIMULATOR_IS_SYSTEM_BIT} from "../Constants.sol";
import {Utils} from "./Utils.sol";

// Addresses used for the compiler to be replaced with the
// ZKsync-specific opcodes during the compilation.
// IMPORTANT: these are just compile-time constants and are used
// only if used in-place by Yul optimizer.
address constant TO_L1_CALL_ADDRESS = address((1 << 16) - 1);
address constant CODE_ADDRESS_CALL_ADDRESS = address((1 << 16) - 2);
address constant PRECOMPILE_CALL_ADDRESS = address((1 << 16) - 3);
address constant META_CALL_ADDRESS = address((1 << 16) - 4);
address constant MIMIC_CALL_CALL_ADDRESS = address((1 << 16) - 5);
address constant SYSTEM_MIMIC_CALL_CALL_ADDRESS = address((1 << 16) - 6);
address constant MIMIC_CALL_BY_REF_CALL_ADDRESS = address((1 << 16) - 7);
address constant SYSTEM_MIMIC_CALL_BY_REF_CALL_ADDRESS = address((1 << 16) - 8);
address constant RAW_FAR_CALL_CALL_ADDRESS = address((1 << 16) - 9);
address constant RAW_FAR_CALL_BY_REF_CALL_ADDRESS = address((1 << 16) - 10);
address constant SYSTEM_CALL_CALL_ADDRESS = address((1 << 16) - 11);
address constant SYSTEM_CALL_BY_REF_CALL_ADDRESS = address((1 << 16) - 12);
address constant SET_CONTEXT_VALUE_CALL_ADDRESS = address((1 << 16) - 13);
address constant SET_PUBDATA_PRICE_CALL_ADDRESS = address((1 << 16) - 14);
address constant INCREMENT_TX_COUNTER_CALL_ADDRESS = address((1 << 16) - 15);
address constant PTR_CALLDATA_CALL_ADDRESS = address((1 << 16) - 16);
address constant CALLFLAGS_CALL_ADDRESS = address((1 << 16) - 17);
address constant PTR_RETURNDATA_CALL_ADDRESS = address((1 << 16) - 18);
address constant EVENT_INITIALIZE_ADDRESS = address((1 << 16) - 19);
address constant EVENT_WRITE_ADDRESS = address((1 << 16) - 20);
address constant LOAD_CALLDATA_INTO_ACTIVE_PTR_CALL_ADDRESS = address((1 << 16) - 21);
address constant LOAD_LATEST_RETURNDATA_INTO_ACTIVE_PTR_CALL_ADDRESS = address((1 << 16) - 22);
address constant PTR_ADD_INTO_ACTIVE_CALL_ADDRESS = address((1 << 16) - 23);
address constant PTR_SHRINK_INTO_ACTIVE_CALL_ADDRESS = address((1 << 16) - 24);
address constant PTR_PACK_INTO_ACTIVE_CALL_ADDRESS = address((1 << 16) - 25);
address constant MULTIPLICATION_HIGH_ADDRESS = address((1 << 16) - 26);
address constant GET_EXTRA_ABI_DATA_ADDRESS = address((1 << 16) - 27);

// All the offsets are in bits
uint256 constant META_PUBDATA_PUBLISHED_OFFSET = 0 * 8;
uint256 constant META_HEAP_SIZE_OFFSET = 8 * 8;
uint256 constant META_AUX_HEAP_SIZE_OFFSET = 12 * 8;
uint256 constant META_SHARD_ID_OFFSET = 28 * 8;
uint256 constant META_CALLER_SHARD_ID_OFFSET = 29 * 8;
uint256 constant META_CODE_SHARD_ID_OFFSET = 30 * 8;

/// @notice The way to forward the calldata:
/// - Use the current heap (i.e. the same as on EVM).
/// - Use the auxiliary heap.
/// - Forward via a pointer
/// @dev Note, that currently, users do not have access to the auxiliary
/// heap and so the only type of forwarding that will be used by the users
/// are UseHeap and ForwardFatPointer for forwarding a slice of the current calldata
/// to the next call.
enum CalldataForwardingMode {
    UseHeap,
    ForwardFatPointer,
    UseAuxHeap
}

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice A library that allows calling contracts with the `isSystem` flag.
 * @dev It is needed to call ContractDeployer and NonceHolder.
 */
library SystemContractsCaller {
    /// @notice Makes a call with the `isSystem` flag.
    /// @param gasLimit The gas limit for the call.
    /// @param to The address to call.
    /// @param value The value to pass with the transaction.
    /// @param data The calldata.
    /// @return success Whether the transaction has been successful.
    /// @dev Note, that the `isSystem` flag can only be set when calling system contracts.
    function systemCall(uint32 gasLimit, address to, uint256 value, bytes memory data) internal returns (bool success) {
        address callAddr = SYSTEM_CALL_CALL_ADDRESS;

        uint32 dataStart;
        assembly {
            dataStart := add(data, 0x20)
        }
        uint32 dataLength = Utils.safeCastToU32(data.length);

        uint256 farCallAbi = SystemContractsCaller.getFarCallABI({
            dataOffset: 0,
            memoryPage: 0,
            dataStart: dataStart,
            dataLength: dataLength,
            gasPassed: gasLimit,
            // Only rollup is supported for now
            shardId: 0,
            forwardingMode: CalldataForwardingMode.UseHeap,
            isConstructorCall: false,
            isSystemCall: true
        });

        if (value == 0) {
            // Doing the system call directly
            assembly {
                success := call(to, callAddr, 0, 0, farCallAbi, 0, 0)
            }
        } else {
            address msgValueSimulator = MSG_VALUE_SYSTEM_CONTRACT;
            // We need to supply the mask to the MsgValueSimulator to denote
            // that the call should be a system one.
            uint256 forwardMask = MSG_VALUE_SIMULATOR_IS_SYSTEM_BIT;

            assembly {
                success := call(msgValueSimulator, callAddr, value, to, farCallAbi, forwardMask, 0)
            }
        }
    }

    /// @notice Makes a call with the `isSystem` flag.
    /// @param gasLimit The gas limit for the call.
    /// @param to The address to call.
    /// @param value The value to pass with the transaction.
    /// @param data The calldata.
    /// @return success Whether the transaction has been successful.
    /// @return returnData The returndata of the transaction (revert reason in case the transaction has failed).
    /// @dev Note, that the `isSystem` flag can only be set when calling system contracts.
    function systemCallWithReturndata(
        uint32 gasLimit,
        address to,
        uint128 value,
        bytes memory data
    ) internal returns (bool success, bytes memory returnData) {
        success = systemCall(gasLimit, to, value, data);

        uint256 size;
        assembly {
            size := returndatasize()
        }

        returnData = new bytes(size);
        assembly {
            returndatacopy(add(returnData, 0x20), 0, size)
        }
    }

    /// @notice Makes a call with the `isSystem` flag.
    /// @param gasLimit The gas limit for the call.
    /// @param to The address to call.
    /// @param value The value to pass with the transaction.
    /// @param data The calldata.
    /// @return returnData The returndata of the transaction. In case the transaction reverts, the error
    /// bubbles up to the parent frame.
    /// @dev Note, that the `isSystem` flag can only be set when calling system contracts.
    function systemCallWithPropagatedRevert(
        uint32 gasLimit,
        address to,
        uint128 value,
        bytes memory data
    ) internal returns (bytes memory returnData) {
        bool success;
        (success, returnData) = systemCallWithReturndata(gasLimit, to, value, data);

        if (!success) {
            assembly {
                let size := mload(returnData)
                revert(add(returnData, 0x20), size)
            }
        }
    }

    /// @notice Calculates the packed representation of the FarCallABI.
    /// @param dataOffset Calldata offset in memory. Provide 0 unless using custom pointer.
    /// @param memoryPage Memory page to use. Provide 0 unless using custom pointer.
    /// @param dataStart The start of the calldata slice. Provide the offset in memory
    /// if not using custom pointer.
    /// @param dataLength The calldata length. Provide the length of the calldata in bytes
    /// unless using custom pointer.
    /// @param gasPassed The gas to pass with the call.
    /// @param shardId Of the account to call. Currently only 0 is supported.
    /// @param forwardingMode The forwarding mode to use:
    /// - provide CalldataForwardingMode.UseHeap when using your current memory
    /// - provide CalldataForwardingMode.ForwardFatPointer when using custom pointer.
    /// @param isConstructorCall Whether the call will be a call to the constructor
    /// (ignored when the caller is not a system contract).
    /// @param isSystemCall Whether the call will have the `isSystem` flag.
    /// @return farCallAbi The far call ABI.
    /// @dev The `FarCallABI` has the following structure:
    /// pub struct FarCallABI {
    ///     pub memory_quasi_fat_pointer: FatPointer,
    ///     pub gas_passed: u32,
    ///     pub shard_id: u8,
    ///     pub forwarding_mode: FarCallForwardPageType,
    ///     pub constructor_call: bool,
    ///     pub to_system: bool,
    /// }
    ///
    /// The FatPointer struct:
    ///
    /// pub struct FatPointer {
    ///     pub offset: u32, // offset relative to `start`
    ///     pub memory_page: u32, // memory page where slice is located
    ///     pub start: u32, // absolute start of the slice
    ///     pub length: u32, // length of the slice
    /// }
    ///
    /// @dev Note, that the actual layout is the following:
    ///
    /// [0..32) bits -- the calldata offset
    /// [32..64) bits -- the memory page to use. Can be left blank in most of the cases.
    /// [64..96) bits -- the absolute start of the slice
    /// [96..128) bits -- the length of the slice.
    /// [128..192) bits -- empty bits.
    /// [192..224) bits -- gasPassed.
    /// [224..232) bits -- forwarding_mode
    /// [232..240) bits -- shard id.
    /// [240..248) bits -- constructor call flag
    /// [248..256] bits -- system call flag
    function getFarCallABI(
        uint32 dataOffset,
        uint32 memoryPage,
        uint32 dataStart,
        uint32 dataLength,
        uint32 gasPassed,
        uint8 shardId,
        CalldataForwardingMode forwardingMode,
        bool isConstructorCall,
        bool isSystemCall
    ) internal pure returns (uint256 farCallAbi) {
        // Fill in the call parameter fields
        farCallAbi = getFarCallABIWithEmptyFatPointer({
            gasPassed: gasPassed,
            shardId: shardId,
            forwardingMode: forwardingMode,
            isConstructorCall: isConstructorCall,
            isSystemCall: isSystemCall
        });

        // Fill in the fat pointer fields
        farCallAbi |= dataOffset;
        farCallAbi |= (uint256(memoryPage) << 32);
        farCallAbi |= (uint256(dataStart) << 64);
        farCallAbi |= (uint256(dataLength) << 96);
    }

    /// @notice Calculates the packed representation of the FarCallABI with zero fat pointer fields.
    /// @param gasPassed The gas to pass with the call.
    /// @param shardId Of the account to call. Currently only 0 is supported.
    /// @param forwardingMode The forwarding mode to use:
    /// - provide CalldataForwardingMode.UseHeap when using your current memory
    /// - provide CalldataForwardingMode.ForwardFatPointer when using custom pointer.
    /// @param isConstructorCall Whether the call will be a call to the constructor
    /// (ignored when the caller is not a system contract).
    /// @param isSystemCall Whether the call will have the `isSystem` flag.
    /// @return farCallAbiWithEmptyFatPtr The far call ABI with zero fat pointer fields.
    function getFarCallABIWithEmptyFatPointer(
        uint32 gasPassed,
        uint8 shardId,
        CalldataForwardingMode forwardingMode,
        bool isConstructorCall,
        bool isSystemCall
    ) internal pure returns (uint256 farCallAbiWithEmptyFatPtr) {
        farCallAbiWithEmptyFatPtr |= (uint256(gasPassed) << 192);
        farCallAbiWithEmptyFatPtr |= (uint256(forwardingMode) << 224);
        farCallAbiWithEmptyFatPtr |= (uint256(shardId) << 232);
        if (isConstructorCall) {
            farCallAbiWithEmptyFatPtr |= (1 << 240);
        }
        if (isSystemCall) {
            farCallAbiWithEmptyFatPtr |= (1 << 248);
        }
    }
}
