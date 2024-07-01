// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.16;

import {IAutomationRegistryConsumer} from "./interfaces/IAutomationRegistryConsumer.sol";

uint256 constant PERFORM_GAS_CUSHION = 5_000;

interface ISystemContext {
    function gasPerPubdataByte() external view returns (uint256 gasPerPubdataByte);
    function getCurrentPubdataSpent() external view returns (uint256 currentPubdataSpent);
}

interface IGasBoundCaller {
    function gasBoundCall(address _to, uint256 _maxTotalGas, bytes calldata _data) external payable;
}

/**
 * @title AutomationZKSyncForwarder is a relayer that sits between the registry and the customer's target contract
 * @dev The purpose of the forwarder is to give customers a consistent address to authorize against,
 * which stays consistent between migrations. The Forwarder also exposes the registry address, so that users who
 * want to programatically interact with the registry (ie top up funds) can do so.
 */
contract AutomationZKSyncForwarder {
    ISystemContext public constant SYSTEM_CONTEXT_CONTRACT = ISystemContext(address(0x800b));
    address public constant GAS_BOUND_CALLER = address(0xc706EC7dfA5D4Dc87f29f859094165E8290530f5);

    /// @notice the user's target contract address
    address public immutable i_target;

    /// @notice the shared logic address
    address public immutable i_logic;

    IAutomationRegistryConsumer public s_registry;

    event GasDetails(uint256 indexed pubdataUsed, uint256 indexed gasPerPubdataByte, uint256 indexed executionGasUsed, uint256 gasprice);

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
        bytes memory returnData;

        (success, returnData) = GAS_BOUND_CALLER.call{gas: gasAmount}(abi.encodeWithSelector(IGasBoundCaller.gasBoundCall.selector, target, gasAmount, data));
        uint256 pubdataGasSpent;
        if (success) {
            (, pubdataGasSpent) = abi.decode(returnData, (bytes, uint256));
        }

        uint256 g2 = gasleft();
        gasUsed = g1 - g2;
        emit GasDetails(pubdataGasSpent, SYSTEM_CONTEXT_CONTRACT.gasPerPubdataByte(), gasUsed, tx.gasprice);
        return (success, gasUsed, pubdataGasSpent);
    }

/*
0x000000000000000000000000000000000000000000000000000000000000000a
0x0000000000000000000000000000000000000000000000000000000000000000
0x0000000000000000000000000000000000000000000000000000000000d89056
0x0000000000000000000000000000000000000000000000000000000000d093ec
*/
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
