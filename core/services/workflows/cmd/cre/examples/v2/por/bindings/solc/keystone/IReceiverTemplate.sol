// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IERC165} from "./IERC165.sol";

import {IReceiver} from "./IReceiver.sol";

/// @title IReceiverTemplate - Abstract receiver with workflow validation and metadata decoding
abstract contract IReceiverTemplate is IReceiver {

    // Immutable expected values
    address public EXPECTED_AUTHOR;
    bytes10 public EXPECTED_WORKFLOW_NAME;

    // Custom errors
    error InvalidAuthor(address received, address expected);
    error InvalidWorkflowName(bytes10 received, bytes10 expected);

    constructor(
        address expectedAuthor,
        bytes10 expectedWorkflowName
    ) {
        EXPECTED_AUTHOR = expectedAuthor;
        EXPECTED_WORKFLOW_NAME = expectedWorkflowName;
    }

    /// @inheritdoc IReceiver
    function onReport(
        bytes calldata metadata,
        bytes calldata report
    ) external override {
         (address workflowOwner, bytes10 workflowName) =  _decodeMetadata(metadata);

        if (workflowOwner != EXPECTED_AUTHOR) {
            revert InvalidAuthor(workflowOwner, EXPECTED_AUTHOR);
        }
        if (workflowName != EXPECTED_WORKFLOW_NAME) {
            revert InvalidWorkflowName(workflowName, EXPECTED_WORKFLOW_NAME);
        }

        _processReport(report);
    }

    /// @notice Extracts the workflow name and the workflow owner from the metadata parameter of onReport
    /// @param metadata The metadata in bytes format
    /// @return workflowOwner The owner of the workflow
    /// @return workflowName  The name of the workflow
    function _decodeMetadata(
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

    /// @notice Abstract function to process the report
    /// @param report The report calldata
    function _processReport(
        bytes calldata report
    ) internal virtual;

    /// @inheritdoc IERC165
    function supportsInterface(
        bytes4 interfaceId
    ) public pure virtual override returns (bool) {
        return
            interfaceId == type(IReceiver).interfaceId ||
            interfaceId == type(IERC165).interfaceId;
    }
}
