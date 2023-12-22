// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "foundry-lib/forge-std/src/Test.sol";

contract L2EPTest is Test {
  function assertGasUsageIsCloseTo(
    uint256 expectedGasUsage,
    uint256 startGasUsage,
    uint256 finalGasUsage,
    uint256 deviation
  ) public {
    uint256 gasUsed = (startGasUsage - finalGasUsage) * tx.gasprice;
    assertLe(gasUsed, expectedGasUsage + deviation);
    assertGe(gasUsed, expectedGasUsage - deviation);
  }
}
