// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Client} from "../../../libraries/Client.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OnRamp_withdrawFeeTokens is OnRampSetup {
  mapping(address => uint256) internal s_nopFees;

  function setUp() public virtual override {
    super.setUp();

    // Since we'll mostly be testing for valid calls from the router we'll
    // mock all calls to be originating from the router and re-mock in
    // tests that require failure.
    vm.startPrank(address(s_sourceRouter));

    uint256 feeAmount = 1234567890;

    // Send a bunch of messages, increasing the juels in the contract
    for (uint256 i = 0; i < s_sourceFeeTokens.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = s_sourceFeeTokens[i % s_sourceFeeTokens.length];
      uint256 newFeeTokenBalance = IERC20(message.feeToken).balanceOf(address(s_onRamp)) + feeAmount;
      deal(message.feeToken, address(s_onRamp), newFeeTokenBalance);
      s_nopFees[message.feeToken] = newFeeTokenBalance;
      s_onRamp.forwardFromRouter(DEST_CHAIN_SELECTOR, message, feeAmount, OWNER);
    }
  }

  function testFuzz_WithdrawFeeTokens_Success(
    uint256[5] memory amounts
  ) public {
    vm.startPrank(OWNER);
    address[] memory feeTokens = new address[](amounts.length);
    for (uint256 i = 0; i < amounts.length; ++i) {
      vm.assume(amounts[i] > 0);
      feeTokens[i] = _deploySourceToken("", amounts[i], 18);
      IERC20(feeTokens[i]).transfer(address(s_onRamp), amounts[i]);
    }

    s_feeQuoter.applyFeeTokensUpdates(new address[](0), feeTokens);

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      vm.expectEmit();
      emit OnRamp.FeeTokenWithdrawn(FEE_AGGREGATOR, feeTokens[i], amounts[i]);
    }

    s_onRamp.withdrawFeeTokens(feeTokens);

    for (uint256 i = 0; i < feeTokens.length; ++i) {
      assertEq(IERC20(feeTokens[i]).balanceOf(FEE_AGGREGATOR), amounts[i]);
      assertEq(IERC20(feeTokens[i]).balanceOf(address(s_onRamp)), 0);
    }
  }

  function test_WithdrawFeeTokens() public {
    vm.expectEmit();
    emit OnRamp.FeeTokenWithdrawn(FEE_AGGREGATOR, s_sourceFeeToken, s_nopFees[s_sourceFeeToken]);

    s_onRamp.withdrawFeeTokens(s_sourceFeeTokens);

    assertEq(IERC20(s_sourceFeeToken).balanceOf(FEE_AGGREGATOR), s_nopFees[s_sourceFeeToken]);
    assertEq(IERC20(s_sourceFeeToken).balanceOf(address(s_onRamp)), 0);
  }
}
