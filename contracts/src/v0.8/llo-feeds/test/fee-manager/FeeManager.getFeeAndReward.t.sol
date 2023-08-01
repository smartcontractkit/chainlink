// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../FeeManager.sol";
import {IFeeManager} from "../../interfaces/IFeeManager.sol";
import {Common} from "../../../libraries/internal/Common.sol";
import "./BaseFeeManager.t.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the fee manager's getFeeAndReward
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  function test_baseFeeIsAppliedForNative() public {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_baseFeeIsAppliedForLink() public {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_discountAIsNotAppliedWhenSetForOtherUsers() public {
    //set the subscriber discount for another user
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), INVALID_ADDRESS);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_discountIsNotAppliedForInvalidTokenAddress() public {
    //should revert with invalid address as it's not a configured token
    vm.expectRevert(INVALID_ADDRESS_ERROR);

    //set the subscriber discount for another user
    setSubscriberDiscount(USER, DEFAULT_FEED_1, INVALID_ADDRESS, FEE_SCALAR / 2, ADMIN);
  }

  function test_discountIsAppliedForLink() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE / 2);
  }

  function test_DiscountIsAppliedForNative() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE / 2);
  }

  function test_discountIsNoLongerAppliedAfterRemoving() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE / 2);

    //remove the discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), 0, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_premiumIsApplied() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 5;

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium
    uint256 expectedPremium = ((DEFAULT_REPORT_NATIVE_FEE * nativePremium) / FEE_SCALAR);

    //expected fee should the base fee offset by the premium and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedPremium);
  }

  function test_premiumIsNotAppliedForLinkFee() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 5;

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_premiumIsNoLongerAppliedAfterRemoving() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 5;

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium
    uint256 expectedPremium = ((DEFAULT_REPORT_NATIVE_FEE * nativePremium) / FEE_SCALAR);

    //expected fee should the base fee offset by the premium and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedPremium);

    //remove the premium
    setNativePremium(0, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_feeIsUpdatedAfterNewPremiumIsApplied() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 5;

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium
    uint256 expectedPremium = ((DEFAULT_REPORT_NATIVE_FEE * nativePremium) / FEE_SCALAR);

    //expected fee should the base fee offset by the premium and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedPremium);

    //change the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium
    expectedPremium = ((DEFAULT_REPORT_NATIVE_FEE * nativePremium) / FEE_SCALAR);

    //expected fee should the base fee offset by the premium and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedPremium);
  }

  function test_premiumIsAppliedForNativeFeeWithDiscount() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 5;

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium quantity
    uint256 expectedPremium = ((DEFAULT_REPORT_NATIVE_FEE * nativePremium) / FEE_SCALAR);

    //calculate the expected discount quantity
    uint256 expectedDiscount = ((DEFAULT_REPORT_NATIVE_FEE + expectedPremium) / 2);

    //expected fee should the base fee offset by the premium and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedPremium - expectedDiscount);
  }

  function test_emptyQuoteReturnsLinkBaseFee() public {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), IFeeManager.Quote(address(0)), USER);

    //fee should be the base link fee
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_nativePremium100Percent() public {
    //set the premium
    setNativePremium(FEE_SCALAR, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE * 2);
  }

  function test_nativePremium0Percent() public {
    //set the premium
    setNativePremium(0, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_nativePremiumCannotExceed100Percent() public {
    //should revert if premium is greater than 100%
    vm.expectRevert(INVALID_PREMIUM_ERROR);

    //set the premium
    setNativePremium(FEE_SCALAR + 1, ADMIN);
  }

  function test_discountIsAppliedWith100PercentPremium() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //set the premium
    setNativePremium(FEE_SCALAR, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected discount quantity
    uint256 expectedDiscount = DEFAULT_REPORT_NATIVE_FEE;

    //fee should be zero
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE * 2 - expectedDiscount);
  }

  function test_feeIsZeroWith100PercentDiscount() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_feeIsUpdatedAfterDiscountIsRemoved() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected discount quantity
    uint256 expectedDiscount = DEFAULT_REPORT_NATIVE_FEE / 2;

    //fee should be 50% of the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);

    //remove the discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), 0, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_feeIsUpdatedAfterNewDiscountIsApplied() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected discount quantity
    uint256 expectedDiscount = DEFAULT_REPORT_NATIVE_FEE / 2;

    //fee should be 50% of the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);

    //change the discount to 25%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 4, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //expected discount is now 25%
    expectedDiscount = DEFAULT_REPORT_NATIVE_FEE / 4;

    //fee should be the base fee minus the expected discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);
  }

  function test_setDiscountOver100Percent() public {
    //should revert with invalid discount
    vm.expectRevert(INVALID_DISCOUNT_ERROR);

    //set the subscriber discount to over 100%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR + 1, ADMIN);
  }

  function test_premiumIsNotAppliedWith100PercentDiscount() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 5;

    //set the subscriber discount to 100%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR, ADMIN);

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_nonAdminUserCanNotSetDiscount() public {
    //should revert with unauthorized
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //change to the user prank
    changePrank(ADMIN);

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR, USER);
  }

  function test_premiumRoundsDownWhenUneven() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 3;

    //set the premium
    setNativePremium(nativePremium, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium quantity
    uint256 expectedPremium = ((DEFAULT_REPORT_NATIVE_FEE * nativePremium) / FEE_SCALAR);

    //expected fee should the base fee offset by the premium and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedPremium);
  }

  function test_discountRoundsUpWhenUneven() public {
    //native premium
    uint256 discount = FEE_SCALAR / 3;

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), discount, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected premium quantity
    uint256 expectedDiscount = ((DEFAULT_REPORT_NATIVE_FEE * discount) / FEE_SCALAR);

    //fee should be zero
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);
  }

  function test_reportWithNoExpiryOrFeeReturnsZero() public {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReport(DEFAULT_FEED_1), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_correctDiscountIsAppliedWhenBothTokensAreDiscounted() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 4, ADMIN);
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager for both tokens
    Common.Asset memory linkFee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);
    Common.Asset memory nativeFee = getFee(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //calculate the expected discount quantity for each token
    uint256 expectedDiscountLink = (DEFAULT_REPORT_LINK_FEE * FEE_SCALAR) / 4 / FEE_SCALAR;
    uint256 expectedDiscountNative = (DEFAULT_REPORT_NATIVE_FEE * FEE_SCALAR) / 2 / FEE_SCALAR;

    //check the fee calculation for each token
    assertEq(linkFee.amount, DEFAULT_REPORT_LINK_FEE - expectedDiscountLink);
    assertEq(nativeFee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscountNative);
  }

  function test_discountIsNotAppliedToOtherFeeds() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_2), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_noFeeIsAppliedWhenReportHasZeroFee() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(
      getReportWithCustomExpiryAndFee(
        DEFAULT_FEED_1,
        uint32(block.timestamp + DEFAULT_REPORT_EXPIRY_OFFSET_SECONDS),
        0,
        0
      ),
      getNativeQuote(),
      USER
    );

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_noFeeIsAppliedWhenReportHasZeroFeeAndDiscountAndPremiumIsSet() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getNativeAddress(), FEE_SCALAR / 2, ADMIN);

    //set the premium
    setNativePremium(FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(
      getReportWithCustomExpiryAndFee(
        DEFAULT_FEED_1,
        uint32(block.timestamp + DEFAULT_REPORT_EXPIRY_OFFSET_SECONDS),
        0,
        0
      ),
      getNativeQuote(),
      USER
    );

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_nativePremiumEventIsEmittedOnUpdate() public {
    //native premium
    uint256 nativePremium = FEE_SCALAR / 3;

    //an event should be emitted
    vm.expectEmit();

    //emit the event which we expect to be emitted
    emit NativePremiumSet(nativePremium);

    //set the premium
    setNativePremium(nativePremium, ADMIN);
  }

  function test_getBaseRewardWithLinkQuote() public {
    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //the reward should equal the base fee
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_getRewardWithLinkQuoteAndLinkDiscount() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //the reward should equal the discounted base fee
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE / 2);
  }

  function test_getRewardWithNativeQuote() public {
    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //the reward should equal the base fee in link
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_getRewardWithNativeQuoteAndPremium() public {
    //set the native premium
    setNativePremium(FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //the reward should equal the base fee in link regardless of the premium
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_getRewardWithLinkDiscount() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //the reward should equal the discounted base fee
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE / 2);
  }

  function test_getLinkFeeIsRoundedUp() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 3, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //the reward should equal .66% + 1 of the base fee due to a 33% discount rounded up
    assertEq(fee.amount, (DEFAULT_REPORT_LINK_FEE * 2) / 3 + 1);
  }

  function test_getLinkRewardIsRoundedDown() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 3, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getLinkQuote(), USER);

    //the reward should equal .66% of the base fee due to a 33% discount rounded down
    assertEq(fee.amount, (DEFAULT_REPORT_LINK_FEE * 2) / 3);
  }

  function test_getLinkRewardWithNativeQuoteAndPremiumWithLinkDiscount() public {
    //set the native premium
    setNativePremium(FEE_SCALAR / 2, ADMIN);

    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1, getLinkAddress(), FEE_SCALAR / 3, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getReward(getReportWithFee(DEFAULT_FEED_1), getNativeQuote(), USER);

    //the reward should equal the base fee in link regardless of the premium
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_testRevertIfReportHasExpired() public {
    //expect a revert
    vm.expectRevert(EXPIRED_REPORT_ERROR);

    //get the fee required by the feeManager
    getFee(
      getReportWithCustomExpiryAndFee(
        DEFAULT_FEED_1,
        block.timestamp - 1,
        DEFAULT_REPORT_LINK_FEE,
        DEFAULT_REPORT_NATIVE_FEE
      ),
      getNativeQuote(),
      USER
    );
  }
}
