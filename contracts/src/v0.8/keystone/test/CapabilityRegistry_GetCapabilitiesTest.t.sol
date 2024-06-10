// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetCapabilitiesTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;
    s_capabilityRegistry.addCapabilities(capabilities);
  }

  function test_ReturnsCapabilities() public view {
    (bytes32[] memory hashedCapabilityIds, CapabilityRegistry.Capability[] memory capabilities) = s_capabilityRegistry
      .getCapabilities();

    assertEq(hashedCapabilityIds.length, 2);
    assertEq(hashedCapabilityIds[0], keccak256(abi.encode(capabilities[0].labelledName, capabilities[0].version)));
    assertEq(hashedCapabilityIds[1], keccak256(abi.encode(capabilities[1].labelledName, capabilities[1].version)));

    assertEq(capabilities.length, 2);

    assertEq(capabilities[0].labelledName, "data-streams-reports");
    assertEq(capabilities[0].version, "1.0.0");
    assertEq(uint256(capabilities[0].responseType), uint256(CapabilityRegistry.CapabilityResponseType.REPORT));
    assertEq(uint256(capabilities[0].capabilityType), uint256(CapabilityRegistry.CapabilityType.TRIGGER));
    assertEq(capabilities[0].configurationContract, address(0));

    assertEq(capabilities[1].labelledName, "read-ethereum-mainnet-gas-price");
    assertEq(capabilities[1].version, "1.0.2");
    assertEq(
      uint256(capabilities[1].responseType),
      uint256(CapabilityRegistry.CapabilityResponseType.OBSERVATION_IDENTICAL)
    );
    assertEq(uint256(capabilities[1].capabilityType), uint256(CapabilityRegistry.CapabilityType.ACTION));
    assertEq(capabilities[1].configurationContract, address(s_capabilityConfigurationContract));
  }

  function test_ExcludesDeprecatedCapabilities() public {
    bytes32 hashedCapabilityId = s_capabilityRegistry.getHashedCapabilityId(
      s_basicCapability.labelledName,
      s_basicCapability.version
    );
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = hashedCapabilityId;
    s_capabilityRegistry.deprecateCapabilities(deprecatedCapabilities);

    (bytes32[] memory hashedCapabilityIds, CapabilityRegistry.Capability[] memory capabilities) = s_capabilityRegistry
      .getCapabilities();

    assertEq(hashedCapabilityIds.length, 1);
    assertEq(hashedCapabilityIds[0], keccak256(abi.encode(capabilities[0].labelledName, capabilities[0].version)));

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
