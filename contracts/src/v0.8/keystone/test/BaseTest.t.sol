// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {Constants} from "./Constants.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract BaseTest is Test, Constants {
    CapabilityRegistry internal s_capabilityRegistry;

    function setUp() public virtual {
        vm.startPrank(ADMIN);
        s_capabilityRegistry = new CapabilityRegistry();
    }
}
