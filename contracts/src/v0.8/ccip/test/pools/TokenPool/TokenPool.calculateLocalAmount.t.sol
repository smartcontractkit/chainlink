// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

contract TokenPool_calculateLocalAmount is TokenPoolSetup {
  function test_calculateLocalAmount() public view {
    uint8 localDecimals = s_tokenPool.getTokenDecimals();
    uint256 remoteAmount = 123e18;

    // Zero decimals should return amount * 10^localDecimals
    assertEq(s_tokenPool.calculateLocalAmount(remoteAmount, 0), remoteAmount * 10 ** localDecimals);

    // Equal decimals should return the same amount
    assertEq(s_tokenPool.calculateLocalAmount(remoteAmount, localDecimals), remoteAmount);

    // Remote amount with more decimals should return less local amount
    uint256 expectedAmount = remoteAmount;
    for (uint8 remoteDecimals = localDecimals + 1; remoteDecimals < 36; ++remoteDecimals) {
      expectedAmount /= 10;
      assertEq(s_tokenPool.calculateLocalAmount(remoteAmount, remoteDecimals), expectedAmount);
    }

    // Remote amount with less decimals should return more local amount
    expectedAmount = remoteAmount;
    for (uint8 remoteDecimals = localDecimals - 1; remoteDecimals > 0; --remoteDecimals) {
      expectedAmount *= 10;
      assertEq(s_tokenPool.calculateLocalAmount(remoteAmount, remoteDecimals), expectedAmount);
    }
  }
}
