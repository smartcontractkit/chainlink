// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC165} from "./keystone/IERC165.sol";
import {IReceiver} from "./keystone/IReceiver.sol";
import {IReserveManager, UpdateReserves} from "./IReserveManager.sol";

contract UpdateReservesProxy is IReceiver {
    IReserveManager public reserveManager;

    address[] internal s_allowedWorkflowOwnersList;
    mapping(address => bool) internal s_allowedWorkflowOwners;
    bytes10[] internal s_allowedWorkflowNamesList;
    mapping(bytes10 => bool) internal s_allowedWorkflowNames;

    error UnauthorizedWorkflowOwner(address workflowOwner);
    error UnauthorizedWorkflowName(bytes10 workflowName);

    constructor(address _reserveManager) {
        reserveManager = IReserveManager(_reserveManager);
    }

    /// @inheritdoc IReceiver
    function onReport(
        bytes calldata metadata,
        bytes calldata report
    ) external override {
        (address workflowOwner, bytes10 workflowName) = _getWorkflowMetaData(
            metadata
        );
        if (!s_allowedWorkflowNames[workflowName]) {
            revert UnauthorizedWorkflowName(workflowName);
        }
        if (!s_allowedWorkflowOwners[workflowOwner]) {
            revert UnauthorizedWorkflowOwner(workflowOwner);
        }

        // Decode the report to get the UpdateReserves struct
        UpdateReserves memory updateReservesData = abi.decode(
            report,
            (UpdateReserves)
        );

        // Call updateReserves on the reserveManager contract
        reserveManager.updateReserves(updateReservesData);
    }

    /// @inheritdoc IERC165
    function supportsInterface(
        bytes4 interfaceId
    ) public pure override returns (bool) {
        return
            interfaceId == type(IReceiver).interfaceId ||
            interfaceId == type(IERC165).interfaceId;
    }

    /// @notice Extracts the workflow name and the workflow owner from the metadata parameter of onReport
    /// @param metadata The metadata in bytes format
    /// @return workflowOwner The owner of the workflow
    /// @return workflowName  The name of the workflow
    function _getWorkflowMetaData(
        bytes memory metadata
    ) internal pure returns (address, bytes10) {
        address workflowOwner;
        bytes10 workflowName;
        // (first 32 bytes contain length of the byte array)
        // workflow_cid             // offset 32, size 32
        // workflow_name            // offset 64, size 10
        // workflow_owner           // offset 74, size 20
        // report_name              // offset 94, size  2
        assembly {
            // no shifting needed for bytes10 type
            workflowName := mload(add(metadata, 64))
            // shift right by 12 bytes to get the actual value
            workflowOwner := shr(mul(12, 8), mload(add(metadata, 74)))
        }
        return (workflowOwner, workflowName);
    }
}
