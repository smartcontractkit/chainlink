// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetCapabilitiesTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();

    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);
  }

  function test_ReturnsCapabilities() public view {
    CapabilityRegistry.Capability[] memory capabilities = s_capabilityRegistry.getCapabilities();

    assertEq(capabilities.length, 2);

    assertEq(capabilities[0].labelledName, "data-streams-reports");
    assertEq(capabilities[0].version, "1.0.0");
    assertEq(uint256(capabilities[0].responseType), uint256(CapabilityRegistry.CapabilityResponseType.REPORT));
    assertEq(capabilities[0].configurationContract, address(0));

    assertEq(capabilities[1].labelledName, "read-ethereum-mainnet-gas-price");
    assertEq(capabilities[1].version, "1.0.2");
    assertEq(
      uint256(capabilities[1].responseType),
      uint256(CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL)
    );
    assertEq(capabilities[1].configurationContract, address(s_capabilityConfigurationContract));
  }

  function test_ExcludesDeprecatedCapabilities() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    s_capabilityRegistry.deprecateCapability(hashedCapabilityId);

    CapabilityRegistry.Capability[] memory capabilities = s_capabilityRegistry.getCapabilities();
    assertEq(capabilities.length, 1);

    assertEq(capabilities[0].labelledName, "read-ethereum-mainnet-gas-price");
    assertEq(capabilities[0].version, "1.0.2");
    assertEq(
      uint256(capabilities[0].responseType),
      uint256(CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL)
    );
    assertEq(capabilities[0].configurationContract, address(s_capabilityConfigurationContract));
  }
}
