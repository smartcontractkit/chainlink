// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
  CapabilityRegistry.Capability private validCapability =
    CapabilityRegistry.Capability({
      capabilityType: "data-streams-reports",
      version: "1.0.0",
      responseType: CapabilityRegistry.CapabilityResponseType.REPORT,
      configurationContract: address(0)
    });

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);

    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.addCapability(validCapability);
  }

  function test_RevertWhen_CapabilityExists() public {
    // Successfully add the capability the first time
    s_capabilityRegistry.addCapability(validCapability);

    // Try to add the same capability again
    vm.expectRevert(CapabilityRegistry.CapabilityAlreadyExists.selector);
    s_capabilityRegistry.addCapability(validCapability);
  }

  function test_AddCapability() public {
    s_capabilityRegistry.addCapability(validCapability);

    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(bytes32("data-streams-reports"), bytes32("1.0.0"));
    CapabilityRegistry.Capability memory storedCapability = s_capabilityRegistry.getCapability(capabilityId);

    assertEq(storedCapability.capabilityType, "data-streams-reports");
    assertEq(storedCapability.version, "1.0.0");
    assertEq(uint256(storedCapability.responseType), uint256(CapabilityRegistry.CapabilityResponseType.REPORT));
    assertEq(storedCapability.configurationContract, address(0));
  }
}
