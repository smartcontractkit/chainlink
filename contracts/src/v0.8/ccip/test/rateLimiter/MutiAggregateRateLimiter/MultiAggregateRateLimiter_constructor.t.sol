// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {MultiAggregateRateLimiterHelper} from "../../helpers/MultiAggregateRateLimiterHelper.sol";
import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";
import {Vm} from "forge-std/Vm.sol";

contract MultiAggregateRateLimiter_constructor is MultiAggregateRateLimiterSetup {
  function test_ConstructorNoAuthorizedCallers() public {
    address[] memory authorizedCallers = new address[](0);

    vm.recordLogs();
    s_rateLimiter = new MultiAggregateRateLimiterHelper(address(s_feeQuoter), authorizedCallers);

    // FeeQuoterSet
    Vm.Log[] memory logEntries = vm.getRecordedLogs();
    assertEq(logEntries.length, 1);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_feeQuoter), s_rateLimiter.getFeeQuoter());
  }

  function test_Constructor() public {
    address[] memory authorizedCallers = new address[](2);
    authorizedCallers[0] = MOCK_OFFRAMP;
    authorizedCallers[1] = MOCK_ONRAMP;

    vm.expectEmit();
    emit MultiAggregateRateLimiter.FeeQuoterSet(address(s_feeQuoter));

    s_rateLimiter = new MultiAggregateRateLimiterHelper(address(s_feeQuoter), authorizedCallers);

    assertEq(OWNER, s_rateLimiter.owner());
    assertEq(address(s_feeQuoter), s_rateLimiter.getFeeQuoter());
    assertEq(s_rateLimiter.typeAndVersion(), "MultiAggregateRateLimiter 1.6.0");
  }
}
