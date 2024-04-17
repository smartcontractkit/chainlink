// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {Capability, CapabilityRegistry} from "./CapabilityRegistry.sol";

contract CapabilityRegistryTest is Test {
    function setUp() public virtual {}

    function testAddCapability() public {
        CapabilityRegistry capabilityRegistry = new CapabilityRegistry();

        capabilityRegistry.addCapability(Capability("data-streams-reports", "1.0.0"));

        bytes32 capabilityId = capabilityRegistry.getCapabilityID(Capability("data-streams-reports", "1.0.0"));
        Capability memory capability = capabilityRegistry.getCapability(capabilityId);

        assertEq(capability.capabilityType, "data-streams-reports");
        assertEq(capability.version, "1.0.0");
    }
}
