// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
  event CapabilityDeprecated(bytes32 indexed hashedCapabilityId);

  function setUp() public override {
    BaseTest.setUp();

    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.deprecateCapability(hashedCapabilityId);
  }

  function test_RevertWhen_CapabilityDoesNotExist() public {
    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.CapabilityDoesNotExist.selector, s_nonExistentHashedCapabilityId)
    );
    s_capabilityRegistry.deprecateCapability(s_nonExistentHashedCapabilityId);
  }

  function test_RevertWhen_CapabilityIsDeprecated() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    s_capabilityRegistry.deprecateCapability(hashedCapabilityId);

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.CapabilityIsDeprecated.selector, hashedCapabilityId));
    s_capabilityRegistry.deprecateCapability(hashedCapabilityId);
  }

  function test_DeprecatesCapability() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    s_capabilityRegistry.deprecateCapability(hashedCapabilityId);

    assertEq(s_capabilityRegistry.isCapabilityDeprecated(hashedCapabilityId), true);
  }

  function test_EmitsEvent() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );

    vm.expectEmit(address(s_capabilityRegistry));
    emit CapabilityDeprecated(hashedCapabilityId);
    s_capabilityRegistry.deprecateCapability(hashedCapabilityId);
  }
}
