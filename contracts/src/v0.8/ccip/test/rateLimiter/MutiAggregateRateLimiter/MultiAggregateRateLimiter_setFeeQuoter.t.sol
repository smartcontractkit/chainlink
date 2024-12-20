// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {AuthorizedCallers} from "../../../../shared/access/AuthorizedCallers.sol";
import {Ownable2Step} from "../../../../shared/access/Ownable2Step.sol";
import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";

contract MultiAggregateRateLimiter_setFeeQuoter is MultiAggregateRateLimiterSetup {
  function test_Owner() public {
    address newAddress = address(42);

    vm.expectEmit();
    emit MultiAggregateRateLimiter.FeeQuoterSet(newAddress);

    s_rateLimiter.setFeeQuoter(newAddress);
    assertEq(newAddress, s_rateLimiter.getFeeQuoter());
  }

  // Reverts

  function test_RevertWhen_OnlyOwner() public {
    vm.startPrank(STRANGER);
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);

    s_rateLimiter.setFeeQuoter(STRANGER);
  }

  function test_RevertWhen_ZeroAddress() public {
    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);
    s_rateLimiter.setFeeQuoter(address(0));
  }
}
