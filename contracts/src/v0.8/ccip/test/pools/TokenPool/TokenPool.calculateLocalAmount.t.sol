// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolHelper} from "../../helpers/TokenPoolHelper.sol";
import {TokenPoolSetup} from "./TokenPoolSetup.t.sol";

import {IERC20Metadata} from
  "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol";

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

  // Reverts

  function test_RevertWhen_calculateLocalAmountWhen_LowRemoteDecimalsOverflows() public {
    uint8 remoteDecimals = 0;
    uint8 localDecimals = 78;
    uint256 remoteAmount = 1;

    vm.mockCall(address(s_token), abi.encodeWithSelector(IERC20Metadata.decimals.selector), abi.encode(localDecimals));

    s_tokenPool =
      new TokenPoolHelper(s_token, localDecimals, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter));

    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.OverflowDetected.selector, remoteDecimals, localDecimals, remoteAmount)
    );

    s_tokenPool.calculateLocalAmount(remoteAmount, remoteDecimals);
  }

  function test_RevertWhen_calculateLocalAmountWhen_HighLocalDecimalsOverflows() public {
    uint8 remoteDecimals = 18;
    uint8 localDecimals = 18 + 78;
    uint256 remoteAmount = 1;

    vm.mockCall(address(s_token), abi.encodeWithSelector(IERC20Metadata.decimals.selector), abi.encode(localDecimals));

    s_tokenPool =
      new TokenPoolHelper(s_token, localDecimals, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter));

    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.OverflowDetected.selector, remoteDecimals, localDecimals, remoteAmount)
    );

    s_tokenPool.calculateLocalAmount(remoteAmount, remoteDecimals);
  }

  function test_RevertWhen_calculateLocalAmountWhen_HighRemoteDecimalsOverflows() public {
    uint8 remoteDecimals = 18 + 78;
    uint8 localDecimals = 18;
    uint256 remoteAmount = 1;

    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.OverflowDetected.selector, remoteDecimals, localDecimals, remoteAmount)
    );

    s_tokenPool.calculateLocalAmount(remoteAmount, remoteDecimals);
  }

  function test_RevertWhen_calculateLocalAmountWhen_HighAmountOverflows() public {
    uint8 remoteDecimals = 18;
    uint8 localDecimals = 18 + 28;
    uint256 remoteAmount = 1e50;

    vm.mockCall(address(s_token), abi.encodeWithSelector(IERC20Metadata.decimals.selector), abi.encode(localDecimals));

    s_tokenPool =
      new TokenPoolHelper(s_token, localDecimals, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter));

    vm.expectRevert(
      abi.encodeWithSelector(TokenPool.OverflowDetected.selector, remoteDecimals, localDecimals, remoteAmount)
    );

    s_tokenPool.calculateLocalAmount(remoteAmount, remoteDecimals);
  }
}
