// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Test} from "forge-std/Test.sol";
import {PrimaryAggregator} from "../PrimaryAggregator.sol";
import {LinkTokenInterface} from "../../shared/interfaces/LinkTokenInterface.sol";
import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";

contract PrimaryAggregatorTest is Test {
  address constant LINK_TOKEN_ADDRESS = address(1);
  address constant BILLING_ACCESS_CONTROLLER_ADDRESS = address(100);
  address constant REQUESTER_ACCESS_CONTROLLER_ADDRESS = address(101);

  PrimaryAggregator aggregator;

  function setUp() public virtual {
    LinkTokenInterface _link = LinkTokenInterface(LINK_TOKEN_ADDRESS);
    AccessControllerInterface _billingAccessController = AccessControllerInterface(BILLING_ACCESS_CONTROLLER_ADDRESS);
    AccessControllerInterface _requesterAccessController = AccessControllerInterface(REQUESTER_ACCESS_CONTROLLER_ADDRESS);

    aggregator = new PrimaryAggregator(_link, 0, 100, _billingAccessController, _requesterAccessController, 18, "TEST");
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

