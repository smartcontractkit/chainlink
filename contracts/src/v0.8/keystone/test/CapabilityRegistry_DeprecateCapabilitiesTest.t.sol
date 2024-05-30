// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_DeprecateCapabilitiesTest is BaseTest {
  event CapabilityDeprecated(bytes32 indexed hashedCapabilityId);

  function setUp() public override {
    BaseTest.setUp();
    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_capabilityRegistry.addCapabilities(capabilities);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);
  }

  function test_RevertWhen_CapabilityDoesNotExist() public {
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = s_nonExistentHashedCapabilityId;

    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.CapabilityDoesNotExist.selector, s_nonExistentHashedCapabilityId)
    );
    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);
  }

  function test_RevertWhen_CapabilityIsDeprecated() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.CapabilityIsDeprecated.selector, hashedCapabilityId));
    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);
  }

  function test_DeprecatesCapability() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);
    assertEq(s_capabilityRegistry.isCapabilityDeprecated(hashedCapabilityId), true);
  }

  function test_EmitsEvent() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;

    vm.expectEmit(address(s_capabilityRegistry));
    emit CapabilityDeprecated(hashedCapabilityId);
    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);
  }
}
