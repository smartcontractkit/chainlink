// SPDX-License-Identifier: MIT
// We use a floating point pragma here so it can be used within other projects that interact with the ZKsync ecosystem without using our exact pragma version.
pragma solidity ^0.8.20;

import {SystemContractHelper, ADDRESS_MASK} from "./SystemContractHelper.sol";
import {SystemContractsCaller, CalldataForwardingMode, RAW_FAR_CALL_BY_REF_CALL_ADDRESS, SYSTEM_CALL_BY_REF_CALL_ADDRESS, MSG_VALUE_SIMULATOR_IS_SYSTEM_BIT, MIMIC_CALL_BY_REF_CALL_ADDRESS} from "./SystemContractsCaller.sol";
import {Utils} from "./Utils.sol";
import {SHA256_SYSTEM_CONTRACT, KECCAK256_SYSTEM_CONTRACT, MSG_VALUE_SYSTEM_CONTRACT} from "../Constants.sol";
import {Keccak256InvalidReturnData, ShaInvalidReturnData} from "../SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice This library is used to perform ultra-efficient calls using zkEVM-specific features.
 * @dev EVM calls always accept a memory slice as input and return a memory slice as output.
 * Therefore, even if the user has a ready-made calldata slice, they still need to copy it to memory
 * before calling. This is especially inefficient for large inputs (proxies, multi-calls, etc.).
 * In turn, zkEVM operates over a fat pointer, which is a set of (memory page, offset, start, length) in the memory/calldata/returndata.
 * This allows forwarding the calldata slice as is, without copying it to memory.
 * @dev Fat pointer is not just an integer, it is an extended data type supported on the VM level.
 * zkEVM creates the wellformed fat pointers for all the calldata/returndata regions, later
 * the contract may manipulate the already created fat pointers to forward a slice of the data, but not
 * to create new fat pointers!
 * @dev The allowed operation on fat pointers are:
 * 1. `ptr.add` - Transforms `ptr.offset` into `ptr.offset + u32(_value)`. If overflow happens then it panics.
 * 2. `ptr.sub` - Transforms `ptr.offset` into `ptr.offset - u32(_value)`. If underflow happens then it panics.
 * 3. `ptr.pack` - Do the concatenation between the lowest 128 bits of the pointer itself and the highest 128 bits of `_value`. It is typically used to prepare the ABI for external calls.
 * 4. `ptr.shrink` - Transforms `ptr.length` into `ptr.length - u32(_shrink)`. If underflow happens then it panics.
 * @dev The call opcodes accept the fat pointer and change it to its canonical form before passing it to the child call
 * 1. `ptr.start` is transformed into `ptr.offset + ptr.start`
 * 2. `ptr.length` is transformed into `ptr.length - ptr.offset`
 * 3. `ptr.offset` is transformed into `0`
 */
library EfficientCall {
    /// @notice Call the `keccak256` without copying calldata to memory.
    /// @param _data The preimage data.
    /// @return The `keccak256` hash.
    function keccak(bytes calldata _data) internal view returns (bytes32) {
        bytes memory returnData = staticCall(gasleft(), KECCAK256_SYSTEM_CONTRACT, _data);
        if (returnData.length != 32) {
            revert Keccak256InvalidReturnData();
        }
        return bytes32(returnData);
    }

    /// @notice Call the `sha256` precompile without copying calldata to memory.
    /// @param _data The preimage data.
    /// @return The `sha256` hash.
    function sha(bytes calldata _data) internal view returns (bytes32) {
        bytes memory returnData = staticCall(gasleft(), SHA256_SYSTEM_CONTRACT, _data);
        if (returnData.length != 32) {
            revert ShaInvalidReturnData();
        }
        return bytes32(returnData);
    }

    /// @notice Perform a `call` without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _value The `msg.value` to send.
    /// @param _data The calldata to use for the call.
    /// @param _isSystem Whether the call should contain the `isSystem` flag.
    /// @return returnData The copied to memory return data.
    function call(
        uint256 _gas,
        address _address,
        uint256 _value,
        bytes calldata _data,
        bool _isSystem
    ) internal returns (bytes memory returnData) {
        bool success = rawCall({_gas: _gas, _address: _address, _value: _value, _data: _data, _isSystem: _isSystem});
        returnData = _verifyCallResult(success);
    }

    /// @notice Perform a `staticCall` without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _data The calldata to use for the call.
    /// @return returnData The copied to memory return data.
    function staticCall(
        uint256 _gas,
        address _address,
        bytes calldata _data
    ) internal view returns (bytes memory returnData) {
        bool success = rawStaticCall(_gas, _address, _data);
        returnData = _verifyCallResult(success);
    }

    /// @notice Perform a `delegateCall` without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _data The calldata to use for the call.
    /// @return returnData The copied to memory return data.
    function delegateCall(
        uint256 _gas,
        address _address,
        bytes calldata _data
    ) internal returns (bytes memory returnData) {
        bool success = rawDelegateCall(_gas, _address, _data);
        returnData = _verifyCallResult(success);
    }

    /// @notice Perform a `mimicCall` (a call with custom msg.sender) without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _data The calldata to use for the call.
    /// @param _whoToMimic The `msg.sender` for the next call.
    /// @param _isConstructor Whether the call should contain the `isConstructor` flag.
    /// @param _isSystem Whether the call should contain the `isSystem` flag.
    /// @return returnData The copied to memory return data.
    function mimicCall(
        uint256 _gas,
        address _address,
        bytes calldata _data,
        address _whoToMimic,
        bool _isConstructor,
        bool _isSystem
    ) internal returns (bytes memory returnData) {
        bool success = rawMimicCall({
            _gas: _gas,
            _address: _address,
            _data: _data,
            _whoToMimic: _whoToMimic,
            _isConstructor: _isConstructor,
            _isSystem: _isSystem
        });

        returnData = _verifyCallResult(success);
    }

    /// @notice Perform a `call` without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _value The `msg.value` to send.
    /// @param _data The calldata to use for the call.
    /// @param _isSystem Whether the call should contain the `isSystem` flag.
    /// @return success whether the call was successful.
    function rawCall(
        uint256 _gas,
        address _address,
        uint256 _value,
        bytes calldata _data,
        bool _isSystem
    ) internal returns (bool success) {
        if (_value == 0) {
            _loadFarCallABIIntoActivePtr(_gas, _data, false, _isSystem);

            address callAddr = RAW_FAR_CALL_BY_REF_CALL_ADDRESS;
            assembly {
                success := call(_address, callAddr, 0, 0, 0xFFFF, 0, 0)
            }
        } else {
            _loadFarCallABIIntoActivePtr(_gas, _data, false, true);

            // If there is provided `msg.value` call the `MsgValueSimulator` to forward ether.
            address msgValueSimulator = MSG_VALUE_SYSTEM_CONTRACT;
            address callAddr = SYSTEM_CALL_BY_REF_CALL_ADDRESS;
            // We need to supply the mask to the MsgValueSimulator to denote
            // that the call should be a system one.
            uint256 forwardMask = _isSystem ? MSG_VALUE_SIMULATOR_IS_SYSTEM_BIT : 0;

            assembly {
                success := call(msgValueSimulator, callAddr, _value, _address, 0xFFFF, forwardMask, 0)
            }
        }
    }

    /// @notice Perform a `staticCall` without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _data The calldata to use for the call.
    /// @return success whether the call was successful.
    function rawStaticCall(uint256 _gas, address _address, bytes calldata _data) internal view returns (bool success) {
        _loadFarCallABIIntoActivePtr(_gas, _data, false, false);

        address callAddr = RAW_FAR_CALL_BY_REF_CALL_ADDRESS;
        assembly {
            success := staticcall(_address, callAddr, 0, 0xFFFF, 0, 0)
        }
    }

    /// @notice Perform a `delegatecall` without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _data The calldata to use for the call.
    /// @return success whether the call was successful.
    function rawDelegateCall(uint256 _gas, address _address, bytes calldata _data) internal returns (bool success) {
        _loadFarCallABIIntoActivePtr(_gas, _data, false, false);

        address callAddr = RAW_FAR_CALL_BY_REF_CALL_ADDRESS;
        assembly {
            success := delegatecall(_address, callAddr, 0, 0xFFFF, 0, 0)
        }
    }

    /// @notice Perform a `mimicCall` (call with custom msg.sender) without copying calldata to memory.
    /// @param _gas The gas to use for the call.
    /// @param _address The address to call.
    /// @param _data The calldata to use for the call.
    /// @param _whoToMimic The `msg.sender` for the next call.
    /// @param _isConstructor Whether the call should contain the `isConstructor` flag.
    /// @param _isSystem Whether the call should contain the `isSystem` flag.
    /// @return success whether the call was successful.
    /// @dev If called not in kernel mode, it will result in a revert (enforced by the VM)
    function rawMimicCall(
        uint256 _gas,
        address _address,
        bytes calldata _data,
        address _whoToMimic,
        bool _isConstructor,
        bool _isSystem
    ) internal returns (bool success) {
        _loadFarCallABIIntoActivePtr(_gas, _data, _isConstructor, _isSystem);

        address callAddr = MIMIC_CALL_BY_REF_CALL_ADDRESS;
        uint256 cleanupMask = ADDRESS_MASK;
        assembly {
            // Clearing values before usage in assembly, since Solidity
            // doesn't do it by default
            _whoToMimic := and(_whoToMimic, cleanupMask)

            success := call(_address, callAddr, 0, 0, _whoToMimic, 0, 0)
        }
    }

    /// @dev Verify that a low-level call was successful, and revert if it wasn't, by bubbling the revert reason.
    /// @param _success Whether the call was successful.
    /// @return returnData The copied to memory return data.
    function _verifyCallResult(bool _success) private pure returns (bytes memory returnData) {
        if (_success) {
            uint256 size;
            assembly {
                size := returndatasize()
            }

            returnData = new bytes(size);
            assembly {
                returndatacopy(add(returnData, 0x20), 0, size)
            }
        } else {
            propagateRevert();
        }
    }

    /// @dev Propagate the revert reason from the current call to the caller.
    function propagateRevert() internal pure {
        assembly {
            let size := returndatasize()
            returndatacopy(0, 0, size)
            revert(0, size)
        }
    }

    /// @dev Load the far call ABI into active ptr, that will be used for the next call by reference.
    /// @param _gas The gas to be passed to the call.
    /// @param _data The calldata to be passed to the call.
    /// @param _isConstructor Whether the call is a constructor call.
    /// @param _isSystem Whether the call is a system call.
    function _loadFarCallABIIntoActivePtr(
        uint256 _gas,
        bytes calldata _data,
        bool _isConstructor,
        bool _isSystem
    ) private view {
        SystemContractHelper.loadCalldataIntoActivePtr();

        uint256 dataOffset;
        assembly {
            dataOffset := _data.offset
        }

        // Safe to cast, offset is never bigger than `type(uint32).max`
        SystemContractHelper.ptrAddIntoActive(uint32(dataOffset));
        // Safe to cast, `data.length` is never bigger than `type(uint32).max`
        uint32 shrinkTo = uint32(msg.data.length - (_data.length + dataOffset));
        SystemContractHelper.ptrShrinkIntoActive(shrinkTo);

        uint32 gas = Utils.safeCastToU32(_gas);
        uint256 farCallAbi = SystemContractsCaller.getFarCallABIWithEmptyFatPointer({
            gasPassed: gas,
            // Only rollup is supported for now
            shardId: 0,
            forwardingMode: CalldataForwardingMode.ForwardFatPointer,
            isConstructorCall: _isConstructor,
            isSystemCall: _isSystem
        });
        SystemContractHelper.ptrPackIntoActivePtr(farCallAbi);
    }
}
