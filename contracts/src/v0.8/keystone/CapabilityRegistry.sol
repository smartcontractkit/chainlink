// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";
import {ICapabilityRegistry} from "./interfaces/ICapabilityRegistry.sol";

contract CapabilityRegistry is ICapabilityRegistry, OwnerIsCreator, TypeAndVersionInterface {
    mapping(bytes32 => Capability) private s_capabilities;

    /// @notice Mapping of node operators
    mapping(uint256 nodeOperatorId => NodeOperator) private s_nodeOperators;

    /// @notice The latest node operator ID
    /// @dev No getter for this as this is an implementation detail
    uint256 private s_nodeOperatorId;

    function typeAndVersion() external pure override returns (string memory) {
        return "CapabilityRegistry 1.0.0";
    }

    /// @inheritdoc ICapabilityRegistry
    function addNodeOperator(address admin, string calldata name) external onlyOwner {
        if (admin == address(0)) revert InvalidNodeOperatorAdmin();
        uint256 nodeOperatorId = s_nodeOperatorId;
        s_nodeOperators[nodeOperatorId] = NodeOperator({id: nodeOperatorId, admin: admin, name: name});
        ++s_nodeOperatorId;
        emit NodeOperatorAdded(nodeOperatorId, admin, name);
    }

    /// @inheritdoc ICapabilityRegistry
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

    /// @inheritdoc ICapabilityRegistry
    function getCapabilityID(bytes32 capabilityType, bytes32 version) public pure returns (bytes32) {
        return keccak256(abi.encodePacked(capabilityType, version));
    }
}
