// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {MultiAggregateRateLimiter} from "../../../MultiAggregateRateLimiter.sol";
import {Client} from "../../../libraries/Client.sol";
import {MultiAggregateRateLimiterSetup} from "./MultiAggregateRateLimiterSetup.t.sol";

contract MultiAggregateRateLimiter_getTokenValue is MultiAggregateRateLimiterSetup {
  function test_GetTokenValue() public view {
    uint256 numberOfTokens = 10;
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: TOKEN, amount: 10});
    uint256 value = s_rateLimiter.getTokenValue(tokenAmount);
    assertEq(value, (numberOfTokens * TOKEN_PRICE) / 1e18);
  }

  // Reverts
  function test_RevertWhen_NoTokenPrices() public {
    address tokenWithNoPrice = makeAddr("Token with no price");
    Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: tokenWithNoPrice, amount: 10});

    vm.expectRevert(abi.encodeWithSelector(MultiAggregateRateLimiter.PriceNotFoundForToken.selector, tokenWithNoPrice));
    s_rateLimiter.getTokenValue(tokenAmount);
  }
}
