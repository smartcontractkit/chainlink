// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {OwnerIsCreator} from "../shared/access/OwnerIsCreator.sol";

struct Capability {
  // Capability type, e.g. "data-streams-reports"
  // bytes32(string); validation regex: ^[a-z0-9_\-:]{1,32}$
  // Not "type" because that's a reserved keyword in Solidity.
  bytes32 capabilityType;
  // Semver, e.g., "1.2.3"
  // bytes32(string); must be valid Semver + max 32 characters.
  bytes32 version;
}

contract CapabilityRegistry is OwnerIsCreator, TypeAndVersionInterface {
  mapping(bytes32 => Capability) private s_capabilities;

  event CapabilityAdded(bytes32 indexed capabilityId);

  function typeAndVersion() external pure override returns (string memory) {
    return "CapabilityRegistry 1.0.0";
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
  /// @return A unique identifier for the capability
  function getCapabilityID(bytes32 capabilityType, bytes32 version) public pure returns (bytes32) {
    return keccak256(abi.encodePacked(capabilityType, version));
  }
}
