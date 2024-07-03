// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_DeprecateCapabilitiesTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_CapabilitiesRegistry.addCapabilities(capabilities);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    vm.expectRevert("Only callable by owner");
    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);
  }

  function test_RevertWhen_CapabilityDoesNotExist() public {
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = s_nonExistentHashedCapabilityId;

    vm.expectRevert(
      abi.encodeWithSelector(CapabilitiesRegistry.CapabilityDoesNotExist.selector, s_nonExistentHashedCapabilityId)
    );
    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);
  }

  function test_RevertWhen_CapabilityIsDeprecated() public {
    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);
    vm.expectRevert(abi.encodeWithSelector(CapabilitiesRegistry.CapabilityIsDeprecated.selector, hashedCapabilityId));
    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);
  }

  function test_DeprecatesCapability() public {
    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);
    assertEq(s_CapabilitiesRegistry.isCapabilityDeprecated(hashedCapabilityId), true);
  }

  function test_EmitsEvent() public {
    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    vm.expectEmit(address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.CapabilityDeprecated(hashedCapabilityId);
    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);
  }
}
