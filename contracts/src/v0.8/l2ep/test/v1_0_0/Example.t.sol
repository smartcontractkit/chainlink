// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";

// Use the following command to run this test file:
//
//  FOUNDRY_PROFILE=l2ep forge test -vvv --mp ./src/v0.8/l2ep/test/v1_0_0/Example.t.sol
//
contract ExampleTest is Test {
  function test_example() public {
    assertEq(true, true);
  }
}
