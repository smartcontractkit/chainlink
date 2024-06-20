// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.16;

import {IAutomationRegistryConsumer} from "./interfaces/IAutomationRegistryConsumer.sol";

uint256 constant PERFORM_GAS_CUSHION = 5_000;

interface ISystemContext {
    function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);
    function getCurrentPubdataSpent() external view returns (uint256 currentPubdataSpent);
}

ISystemContext constant SYSTEM_CONTEXT_CONTRACT = ISystemContext(address(0x800b));

/**
 * @title AutomationZKSyncForwarder is a relayer that sits between the registry and the customer's target contract
 * @dev The purpose of the forwarder is to give customers a consistent address to authorize against,
 * which stays consistent between migrations. The Forwarder also exposes the registry address, so that users who
 * want to programatically interact with the registry (ie top up funds) can do so.
 */
contract AutomationZKSyncForwarder {
    /// @notice the user's target contract address
    address private immutable i_target;

    /// @notice the shared logic address
    address private immutable i_logic;

    IAutomationRegistryConsumer private s_registry;

    event GasDetails(uint256 indexed pubdataUsed, uint256 indexed gasPerPubdataByte, uint256 indexed executionGasUsed, uint256 p1, uint256 p2, uint256 g1, uint256 g2);

    constructor(address target, address registry, address logic) {
        s_registry = IAutomationRegistryConsumer(registry);
        i_target = target;
        i_logic = logic;
    }

    /**
     * @notice forward is called by the registry and forwards the call to the target
   * @param gasAmount is the amount of gas to use in the call
   * @param data is the 4 bytes function selector + arbitrary function data
   * @return success indicating whether the target call succeeded or failed
   */
    function forward(uint256 gasAmount, bytes memory data) external returns (bool success, uint256 gasUsed, uint256 l1GasUsed) {
        if (msg.sender != address(s_registry)) revert();
        address target = i_target;
        uint256 g1 = gasleft();
        uint256 p1 = SYSTEM_CONTEXT_CONTRACT.getCurrentPubdataSpent();
        assembly {
            let g := gas()
        // Compute g -= PERFORM_GAS_CUSHION and check for underflow
            if lt(g, PERFORM_GAS_CUSHION) {
                revert(0, 0)
            }
            g := sub(g, PERFORM_GAS_CUSHION)
        // if g - g//64 <= gasAmount, revert
        // (we subtract g//64 because of EIP-150)
            if iszero(gt(sub(g, div(g, 64)), gasAmount)) {
                revert(0, 0)
            }
        // solidity calls check that a contract actually exists at the destination, so we do the same
            if iszero(extcodesize(target)) {
                revert(0, 0)
            }
        // call with exact gas
            success := call(gasAmount, target, 0, add(data, 0x20), mload(data), 0, 0)
        }
        uint256 p2 = SYSTEM_CONTEXT_CONTRACT.getCurrentPubdataSpent();
        // pubdata size can be less than 0
        uint256 pubdataUsed;
        if (p2 - p1 > 0) {
            pubdataUsed = p2 - p1;
        }

        uint256 gasPerPubdataByte = SYSTEM_CONTEXT_CONTRACT.gasPerPubdataByte();
        uint256 g2 = gasleft();
        gasUsed = g1 - g2 + pubdataUsed * gasPerPubdataByte;
        emit GasDetails(pubdataUsed, gasPerPubdataByte, g1 - g2, p1, p2, g1, g2);
        return (success, gasUsed, pubdataUsed * gasPerPubdataByte);
    }

    function getTarget() external view returns (address) {
        return i_target;
    }

    fallback() external {
        // copy to memory for assembly access
        address logic = i_logic;
        // copied directly from OZ's Proxy contract
        assembly {
        // Copy msg.data. We take full control of memory in this inline assembly
        // block because it will not return to Solidity code. We overwrite the
        // Solidity scratch pad at memory position 0.
            calldatacopy(0, 0, calldatasize())

        // out and outsize are 0 because we don't know the size yet.
            let result := delegatecall(gas(), logic, 0, calldatasize(), 0, 0)

        // Copy the returned data.
            returndatacopy(0, 0, returndatasize())

            switch result
            // delegatecall returns 0 on error.
            case 0 {
                revert(0, returndatasize())
            }
            default {
                return(0, returndatasize())
            }
        }
    }
}
