// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.addCapability(
      CapabilityRegistry.Capability({
        capabilityType: "data-streams-reports",
        version: "1.0.0",
        responseType: CapabilityRegistry.CapabilityResponseType.REPORT,
        configurationContract: address(0)
      })
    );
  }

  function test_AddCapability() public {
    s_capabilityRegistry.addCapability(
      CapabilityRegistry.Capability({
        capabilityType: "data-streams-reports",
        version: "1.0.0",
        responseType: CapabilityRegistry.CapabilityResponseType.REPORT,
        configurationContract: address(0)
      })
    );

    bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(bytes32("data-streams-reports"), bytes32("1.0.0"));
    CapabilityRegistry.Capability memory capability = s_capabilityRegistry.getCapability(capabilityId);

    assertEq(capability.capabilityType, "data-streams-reports");
    assertEq(capability.version, "1.0.0");
    assertEq(uint256(capability.responseType), uint256(CapabilityRegistry.CapabilityResponseType.REPORT));
    assertEq(capability.configurationContract, address(0));
  }
}
