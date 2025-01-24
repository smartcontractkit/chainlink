// SPDX-License-Identifier: MIT

pragma solidity 0.8.24;

import {Utils} from "./libraries/Utils.sol";
import {EfficientCall} from "./libraries/EfficientCall.sol";
import {SystemContractBase} from "./abstract/SystemContractBase.sol";
import {SystemContractHelper} from "./libraries/SystemContractHelper.sol";
import {MSG_VALUE_SIMULATOR_IS_SYSTEM_BIT, REAL_BASE_TOKEN_SYSTEM_CONTRACT} from "./Constants.sol";
import {InvalidCall} from "./SystemContractErrors.sol";

/**
 * @author Matter Labs
 * @custom:security-contact security@matterlabs.dev
 * @notice The contract responsible for simulating transactions with `msg.value` inside zkEVM.
 * @dev It accepts value and whether the call should be system in the first extraAbi param and
 * the address to call in the second extraAbi param, transfers the funds and uses `mimicCall` to continue the
 * call with the same msg.sender.
 */
contract MsgValueSimulator is SystemContractBase {
    /// @notice Extract value, isSystemCall and to from the extraAbi params.
    /// @dev The contract accepts value, the callee and whether the call should be a system one via its ABI params.
    /// @dev The first ABI param contains the value in the [0..127] bits. The 128th contains
    /// the flag whether or not the call should be a system one.
    /// The second ABI params contains the callee.
    function _getAbiParams() internal view returns (uint256 value, bool isSystemCall, address to) {
        value = SystemContractHelper.getExtraAbiData(0);
        uint256 addressAsUint = SystemContractHelper.getExtraAbiData(1);
        uint256 mask = SystemContractHelper.getExtraAbiData(2);

        isSystemCall = (mask & MSG_VALUE_SIMULATOR_IS_SYSTEM_BIT) != 0;

        to = address(uint160(addressAsUint));
    }

    /// @notice The maximal number of gas out of the stipend that should be passed to the callee.
    uint256 private constant GAS_TO_PASS = 2300;

    /// @notice The amount of gas that is passed to the MsgValueSimulator as a stipend.
    /// This number servers to pay for the ETH transfer as well as to provide gas for the `GAS_TO_PASS` gas.
    /// It is equal to the following constant: https://github.com/matter-labs/era-zkevm_opcode_defs/blob/7bf8016f5bb13a73289f321ad6ea8f614540ece9/src/system_params.rs#L96.
    uint256 private constant MSG_VALUE_SIMULATOR_STIPEND_GAS = 27000;

    /// @notice The fallback function that is the main entry point for the MsgValueSimulator.
    /// @dev The contract accepts value, the callee and whether the call should be a system one via its ABI params.
    /// @param _data The calldata to be passed to the callee.
    /// @return The return data from the callee.
    // solhint-disable-next-line payable-fallback
    fallback(bytes calldata _data) external onlySystemCall returns (bytes memory) {
        // Firstly we calculate how much gas has been actually provided by the user to the inner call.
        // For that, we need to get the total gas available in this context and subtract the stipend from it.
        uint256 gasInContext = gasleft();
        // Note, that the `gasInContext` might be slightly less than the MSG_VALUE_SIMULATOR_STIPEND_GAS, since
        // by the time we retrieve it, some gas might have already been spent, e.g. on the `gasleft` opcode itself.
        uint256 userGas = gasInContext > MSG_VALUE_SIMULATOR_STIPEND_GAS
            ? gasInContext - MSG_VALUE_SIMULATOR_STIPEND_GAS
            : 0;

        (uint256 value, bool isSystemCall, address to) = _getAbiParams();

        // Prevent mimic call to the MsgValueSimulator to prevent an unexpected change of callee.
        if (to == address(this)) {
            revert InvalidCall();
        }

        if (value != 0) {
            (bool success, ) = address(REAL_BASE_TOKEN_SYSTEM_CONTRACT).call(
                abi.encodeCall(REAL_BASE_TOKEN_SYSTEM_CONTRACT.transferFromTo, (msg.sender, to, value))
            );

            // If the transfer of ETH fails, we do the most Ethereum-like behaviour in such situation: revert(0,0)
            if (!success) {
                assembly {
                    revert(0, 0)
                }
            }

            // If value is non-zero, we also provide additional gas to the callee.
            userGas += GAS_TO_PASS;
        }

        // For the next call this `msg.value` will be used.
        SystemContractHelper.setValueForNextFarCall(Utils.safeCastToU128(value));

        return
            EfficientCall.mimicCall({
                _gas: userGas,
                _address: to,
                _data: _data,
                _whoToMimic: msg.sender,
                _isConstructor: false,
                _isSystem: isSystemCall
            });
    }
}
