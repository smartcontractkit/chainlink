// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

contract CapabilityRegistry is OwnerIsCreator, TypeAndVersionInterface {
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

    mapping(bytes32 => Capability) private s_capabilities;

    /// @notice Mapping of node operators
    mapping(uint256 nodeOperatorId => NodeOperator) private s_nodeOperators;

    /// @notice The latest node operator ID
    /// @dev No getter for this as this is an implementation detail
    uint256 private s_nodeOperatorId;

    function typeAndVersion() external pure override returns (string memory) {
        return "CapabilityRegistry 1.0.0";
    }

    /// @notice Adds a new node operator
    /// @param admin The address of the admin that can manage the node
    /// operator
    /// @param name The human readable name of the node operator
    function addNodeOperator(address admin, string calldata name) external onlyOwner {
        if (admin == address(0)) revert InvalidNodeOperatorAdmin();
        uint256 nodeOperatorId = s_nodeOperatorId;
        s_nodeOperators[nodeOperatorId] = NodeOperator({id: nodeOperatorId, admin: admin, name: name});
        ++s_nodeOperatorId;
        emit NodeOperatorAdded(nodeOperatorId, admin, name);
    }

    /// @notice Gets a node operator's data
    /// @param nodeOperatorId The ID of the node operator to query for
    /// @return NodeOperator The node operator data
    function getNodeOperator(uint256 nodeOperatorId) external view returns (NodeOperator memory) {
        return s_nodeOperators[nodeOperatorId];
    }

    function addCapability(Capability calldata capability) external onlyOwner {
        bytes32 capabilityId = getCapabilityID(capability.capabilityType, capability.version);
        s_capabilities[capabilityId] = capability;
        emit CapabilityAdded(capabilityId);
    }

    function getCapability(bytes32 capabilityID) public view returns (Capability memory) {
        return s_capabilities[capabilityID];
    }

    /// @notice This functions returns a Capability ID packed into a bytes32 for cheaper access
    /// @return bytes32 A unique identifier for the capability
    function getCapabilityID(bytes32 capabilityType, bytes32 version) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(capabilityType, version));
    }
}
