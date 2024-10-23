// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {PrimaryAggregator} from "../PrimaryAggregator.sol";

contract PrimaryAggregatorTest is Test {
  PrimaryAggregator aggregator;

  function setUp() public virtual {
    // deploy the PrimaryAggregator contract here
    aggregator = new PrimaryAggregator(18);
  }

  function test_decimalsIs18() public {
    assertEq(aggregator.decimals(), 18);
  }

  function test_latestRound_IncrementsAfterTransmit() public {
    assertEq(aggregator.latestRound(), 0);
    // TODO: run a transmit
    assertEq(aggregator.latestRound(), 1);
  }
}

