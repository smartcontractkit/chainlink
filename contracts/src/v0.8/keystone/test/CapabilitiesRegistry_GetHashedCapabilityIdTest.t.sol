// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";

import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_GetHashedCapabilityTest is BaseTest {
  string constant CAPABILITY_LABELLED_NAME = "ccip1";
  string constant CAPABILITY_VERSION = "1.0.0";

  function test_CorrectlyGeneratesHashedCapabilityId() public view {
    bytes32 expectedHashedCapabilityId = keccak256(abi.encode(CAPABILITY_LABELLED_NAME, CAPABILITY_VERSION));

    assertEq(
      s_CapabilitiesRegistry.getHashedCapabilityId(CAPABILITY_LABELLED_NAME, CAPABILITY_VERSION),
      expectedHashedCapabilityId
    );
  }

  function test_DoesNotCauseIncorrectClashes() public view {
    assertNotEq(
      s_CapabilitiesRegistry.getHashedCapabilityId(CAPABILITY_LABELLED_NAME, CAPABILITY_VERSION),
      s_CapabilitiesRegistry.getHashedCapabilityId("ccip", "11.0.0")
    );
  }
}
