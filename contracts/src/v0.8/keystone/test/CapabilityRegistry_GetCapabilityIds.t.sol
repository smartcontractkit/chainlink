// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetCapabilitiesTest is BaseTest {
  function test_ReturnsCapabilities() public {
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

    CapabilityRegistry.Capability[] memory capabilities = s_capabilityRegistry.getCapabilities();

    assertEq(capabilities.length, 2);

    assertEq(capabilities[0].capabilityType, "data-streams-reports");
    assertEq(capabilities[0].version, "1.0.0");
    assertEq(uint256(capabilities[0].responseType), uint256(CapabilityRegistry.CapabilityResponseType.REPORT));
    assertEq(capabilities[0].configurationContract, address(0));

    assertEq(capabilities[1].capabilityType, "read-ethereum-mainnet-gas-price");
    assertEq(capabilities[1].version, "1.0.2");
    assertEq(
      uint256(capabilities[1].responseType),
      uint256(CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL)
    );
    assertEq(capabilities[1].configurationContract, address(s_capabilityConfigurationContract));
  }
}
