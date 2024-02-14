// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Test, stdError} from "forge-std/Test.sol";
import {MockARM} from "./mocks/MockARM.sol";
import {StructFactory} from "./StructFactory.sol";

contract BaseTest is Test, StructFactory {
  bool private s_baseTestInitialized;

  MockARM internal s_mockARM;

  function setUp() public virtual {
    // BaseTest.setUp is often called multiple times from tests' setUp due to inheritance.
    if (s_baseTestInitialized) return;
    s_baseTestInitialized = true;

    // Set the sender to OWNER permanently
    vm.startPrank(OWNER);
    deal(OWNER, 1e20);
    vm.label(OWNER, "Owner");
    vm.label(STRANGER, "Stranger");

    // Set the block time to a constant known value
    vm.warp(BLOCK_TIME);

    s_mockARM = new MockARM();
  }
}
