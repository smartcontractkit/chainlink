// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "forge-std/Test.sol";
import {StructFactory} from "./StructFactory.sol";

contract BaseTest is Test, StructFactory{
    function setUp() public virtual {
        vm.startPrank(OWNER);
        deal(OWNER, 1e20);
    }
}