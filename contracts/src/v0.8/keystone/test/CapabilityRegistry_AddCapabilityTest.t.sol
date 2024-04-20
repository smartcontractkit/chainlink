// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {ICapabilityRegistry} from "../interfaces/ICapabilityRegistry.sol";

contract CapabilityRegistry_AddCapabilityTest is BaseTest {
    function test_AddCapability() public {
        s_capabilityRegistry.addCapability(ICapabilityRegistry.Capability("data-streams-reports", "1.0.0"));

        bytes32 capabilityId = s_capabilityRegistry.getCapabilityID(bytes32("data-streams-reports"), bytes32("1.0.0"));
        ICapabilityRegistry.Capability memory capability = s_capabilityRegistry.getCapability(capabilityId);

        assertEq(capability.capabilityType, "data-streams-reports");
        assertEq(capability.version, "1.0.0");
    }
}
