// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_GetCapabilitiesTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;
    s_CapabilitiesRegistry.addCapabilities(capabilities);
  }

  function test_ReturnsCapabilities() public {
    bytes32 hashedCapabilityId = s_CapabilitiesRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;
    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);

    CapabilitiesRegistry.CapabilityInfo[] memory capabilities = s_CapabilitiesRegistry.getCapabilities();

    assertEq(capabilities.length, 2);

    assertEq(capabilities[0].labelledName, "data-streams-reports");
    assertEq(capabilities[0].version, "1.0.0");
    assertEq(uint256(capabilities[0].responseType), uint256(CapabilitiesRegistry.CapabilityResponseType.REPORT));
    assertEq(uint256(capabilities[0].capabilityType), uint256(CapabilitiesRegistry.CapabilityType.TRIGGER));
    assertEq(capabilities[0].configurationContract, address(0));
    assertEq(capabilities[0].hashedId, keccak256(abi.encode(capabilities[0].labelledName, capabilities[0].version)));
    assertEq(capabilities[0].isDeprecated, true);

    assertEq(capabilities[1].labelledName, "read-ethereum-mainnet-gas-price");
    assertEq(capabilities[1].version, "1.0.2");
    assertEq(
      uint256(capabilities[1].responseType),
      uint256(CapabilitiesRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL)
    );
    assertEq(uint256(capabilities[1].capabilityType), uint256(CapabilitiesRegistry.CapabilityType.ACTION));
    assertEq(capabilities[1].configurationContract, address(s_capabilityConfigurationContract));
    assertEq(capabilities[1].hashedId, keccak256(abi.encode(capabilities[1].labelledName, capabilities[1].version)));
    assertEq(capabilities[1].isDeprecated, false);
  }
}
