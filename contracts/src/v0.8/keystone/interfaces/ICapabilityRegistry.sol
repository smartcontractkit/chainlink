// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../../shared/access/OwnerIsCreator.sol";

interface ICapabilityRegistry {
    struct NodeOperator {
        /// @notice Unique identifier for the node operator
        uint256 id;
        /// @notice The address of the admin that can manage a node
        /// operator
        address admin;
        /// @notice Human readable name of a Node Operator managing the node
        string name;
    }

    struct Capability {
        // Capability type, e.g. "data-streams-reports"
        // bytes32(string); validation regex: ^[a-z0-9_\-:]{1,32}$
        // Not "type" because that's a reserved keyword in Solidity.
        bytes32 capabilityType;
        // Semver, e.g., "1.2.3"
        // bytes32(string); must be valid Semver + max 32 characters.
        bytes32 version;
    }

    /// @notice This error is thrown when trying to set a node operator's
    /// admin address to the zero address
    error InvalidNodeOperatorAdmin();

    /// @notice This event is emitted when a new node operator is added
    /// @param nodeOperatorId The ID of the newly added node operator
    /// @param admin The address of the admin that can manage the node
    /// operator
    /// @param name The human readable name of the node operator
    event NodeOperatorAdded(uint256 nodeOperatorId, address indexed admin, string name);

    /// @notice This event is emitted when a new capability is added
    /// @param capabilityId The ID of the newly added capability
    event CapabilityAdded(bytes32 indexed capabilityId);

    /// @notice Adds a new node operator
    /// @param admin The address of the admin that can manage the node
    /// operator
    /// @param name The human readable name of the node operator
    function addNodeOperator(address admin, string calldata name) external;

    /// @notice Gets a node operator's data
    /// @param nodeOperatorId The ID of the node operator to query for
    /// @return NodeOperator The node operator data
    function getNodeOperator(uint256 nodeOperatorId) external view returns (NodeOperator memory);

    function addCapability(Capability calldata capability) external;

    function getCapability(bytes32 capabilityID) external view returns (Capability memory);

    /// @notice This functions returns a Capability ID packed into a bytes32 for cheaper access
    /// @return bytes32 A unique identifier for the capability
    function getCapabilityID(bytes32 capabilityType, bytes32 version) external pure returns (bytes32);
}
