// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
  event CapabilityDeprecated(bytes32 indexed capabilityId);

  function setUp() public override {
    BaseTest.setUp();

    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(
      s_basicCapability.capabilityType,
      s_basicCapability.version
    );

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.deprecateCapability(capabilityId);
  }

  function test_RevertWhen_CapabilityDoesNotExist() public {
    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID("non-existent-capability", "1.0.0");

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.CapabilityDoesNotExist.selector, capabilityId));
    s_capabilityRegistry.deprecateCapability(capabilityId);
  }

  function test_RevertWhen_CapabilityAlreadyDeprecated() public {
    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(
      s_basicCapability.capabilityType,
      s_basicCapability.version
    );

    s_capabilityRegistry.deprecateCapability(capabilityId);

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.CapabilityAlreadyDeprecated.selector, capabilityId));
    s_capabilityRegistry.deprecateCapability(capabilityId);
  }

  function test_DeprecatesCapability() public {
    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(
      s_basicCapability.capabilityType,
      s_basicCapability.version
    );

    s_capabilityRegistry.deprecateCapability(capabilityId);

    assertEq(s_capabilityRegistry.isCapabilityDeprecated(capabilityId), true);
  }

  function test_EmitsEvent() public {
    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(
      s_basicCapability.capabilityType,
      s_basicCapability.version
    );

    vm.expectEmit(address(s_capabilityRegistry));
    emit CapabilityDeprecated(capabilityId);
    s_capabilityRegistry.deprecateCapability(capabilityId);
  }
}
