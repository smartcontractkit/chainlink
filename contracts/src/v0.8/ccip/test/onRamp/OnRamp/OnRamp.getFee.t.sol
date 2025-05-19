// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {FeeQuoter} from "../../../FeeQuoter.sol";
import {Client} from "../../../libraries/Client.sol";
import {USDPriceWith18Decimals} from "../../../libraries/USDPriceWith18Decimals.sol";
import {OnRamp} from "../../../onRamp/OnRamp.sol";
import {OnRampSetup} from "./OnRampSetup.t.sol";

contract OnRamp_getFee is OnRampSetup {
  using USDPriceWith18Decimals for uint224;

  function test_EmptyMessage() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateEmptyMessage();
      message.feeToken = testTokens[i];

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
      uint256 expectedFeeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      assertEq(expectedFeeAmount, feeAmount);
    }
  }

  function test_SingleTokenMessage() public view {
    address[2] memory testTokens = [s_sourceFeeToken, s_sourceRouter.getWrappedNative()];
    uint224[2] memory feeTokenPrices = [s_feeTokenPrice, s_wrappedTokenPrice];

    uint256 tokenAmount = 10000e18;
    for (uint256 i = 0; i < feeTokenPrices.length; ++i) {
      Client.EVM2AnyMessage memory message = _generateSingleTokenMessage(s_sourceFeeToken, tokenAmount);
      message.feeToken = testTokens[i];

      uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
      uint256 expectedFeeAmount = s_feeQuoter.getValidatedFee(DEST_CHAIN_SELECTOR, message);

      assertEq(expectedFeeAmount, feeAmount);
    }
  }

  function test_GetFeeOfZeroForTokenMessage() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    uint256 feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
    assertTrue(feeAmount > 0);

    FeeQuoter.PremiumMultiplierWeiPerEthArgs[] memory tokenMults = new FeeQuoter.PremiumMultiplierWeiPerEthArgs[](1);
    tokenMults[0] = FeeQuoter.PremiumMultiplierWeiPerEthArgs({token: message.feeToken, premiumMultiplierWeiPerEth: 0});
    s_feeQuoter.applyPremiumMultiplierWeiPerEthUpdates(tokenMults);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.destDataAvailabilityMultiplierBps = 0;
    destChainConfigArgs[0].destChainConfig.gasMultiplierWeiPerEth = 0;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);

    feeAmount = s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);

    assertEq(0, feeAmount);
  }

  // Reverts

  function test_RevertWhen_Unhealthy() public {
    _setMockRMNChainCurse(DEST_CHAIN_SELECTOR, true);
    vm.expectRevert(abi.encodeWithSelector(OnRamp.CursedByRMN.selector, DEST_CHAIN_SELECTOR));
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, _generateEmptyMessage());
  }

  function test_RevertWhen_EnforceOutOfOrder() public {
    // Update dynamic config to enforce allowOutOfOrderExecution = true.
    vm.stopPrank();
    vm.startPrank(OWNER);

    FeeQuoter.DestChainConfigArgs[] memory destChainConfigArgs = _generateFeeQuoterDestChainConfigArgs();
    destChainConfigArgs[0].destChainConfig.enforceOutOfOrder = true;
    s_feeQuoter.applyDestChainConfigUpdates(destChainConfigArgs);
    vm.stopPrank();

    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    // Empty extraArgs to should revert since it enforceOutOfOrder is true.
    message.extraArgs = "";

    vm.expectRevert(FeeQuoter.ExtraArgOutOfOrderExecutionMustBeTrue.selector);
    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }

  function test_RevertWhen_NotAFeeTokenButPricedToken() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.FeeTokenNotSupported.selector, message.feeToken));

    s_onRamp.getFee(DEST_CHAIN_SELECTOR, message);
  }
}
