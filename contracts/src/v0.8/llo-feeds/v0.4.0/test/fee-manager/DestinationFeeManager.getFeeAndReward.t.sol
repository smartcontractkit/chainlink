// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Common} from "../../../libraries/Common.sol";
import "./BaseDestinationFeeManager.t.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager's getFeeAndReward
 */
contract DestinationFeeManagerProcessFeeTest is BaseDestinationFeeManagerTest {
  function test_baseFeeIsAppliedForNative() public view {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_baseFeeIsAppliedForLink() public view {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_discountAIsNotAppliedWhenSetForOtherUsers() public {
    //set the subscriber discount for another user
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), INVALID_ADDRESS);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_discountIsNotAppliedForInvalidTokenAddress() public {
    //should revert with invalid address as it's not a configured token
    vm.expectRevert(INVALID_ADDRESS_ERROR);

    //set the subscriber discount for another user
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, INVALID_ADDRESS, FEE_SCALAR / 2, ADMIN);
  }

  function test_discountIsAppliedForLink() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE / 2);
  }

  function test_DiscountIsAppliedForNative() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE / 2);
  }

  function test_discountIsNoLongerAppliedAfterRemoving() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE / 2);

    //remove the discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), 0, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_surchargeIsApplied() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 5;

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected surcharge
    uint256 expectedSurcharge = ((DEFAULT_REPORT_NATIVE_FEE * nativeSurcharge) / FEE_SCALAR);

    //expected fee should the base fee offset by the surcharge and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge);
  }

  function test_surchargeIsNotAppliedForLinkFee() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 5;

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_surchargeIsNoLongerAppliedAfterRemoving() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 5;

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected surcharge
    uint256 expectedSurcharge = ((DEFAULT_REPORT_NATIVE_FEE * nativeSurcharge) / FEE_SCALAR);

    //expected fee should be the base fee offset by the surcharge and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge);

    //remove the surcharge
    setNativeSurcharge(0, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be the default
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_feeIsUpdatedAfterNewSurchargeIsApplied() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 5;

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected surcharge
    uint256 expectedSurcharge = ((DEFAULT_REPORT_NATIVE_FEE * nativeSurcharge) / FEE_SCALAR);

    //expected fee should the base fee offset by the surcharge and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge);

    //change the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected surcharge
    expectedSurcharge = ((DEFAULT_REPORT_NATIVE_FEE * nativeSurcharge) / FEE_SCALAR);

    //expected fee should the base fee offset by the surcharge and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge);
  }

  function test_surchargeIsAppliedForNativeFeeWithDiscount() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 5;

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected surcharge quantity
    uint256 expectedSurcharge = ((DEFAULT_REPORT_NATIVE_FEE * nativeSurcharge) / FEE_SCALAR);

    //calculate the expected discount quantity
    uint256 expectedDiscount = ((DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge) / 2);

    //expected fee should the base fee offset by the surcharge and discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge - expectedDiscount);
  }

  function test_emptyQuoteRevertsWithError() public {
    //expect a revert
    vm.expectRevert(INVALID_QUOTE_ERROR);

    //get the fee required by the feeManager
    getFee(getV3Report(DEFAULT_FEED_1_V3), address(0), USER);
  }

  function test_nativeSurcharge100Percent() public {
    //set the surcharge
    setNativeSurcharge(FEE_SCALAR, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be twice the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE * 2);
  }

  function test_nativeSurcharge0Percent() public {
    //set the surcharge
    setNativeSurcharge(0, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_nativeSurchargeCannotExceed100Percent() public {
    //should revert if surcharge is greater than 100%
    vm.expectRevert(INVALID_SURCHARGE_ERROR);

    //set the surcharge above the max
    setNativeSurcharge(FEE_SCALAR + 1, ADMIN);
  }

  function test_discountIsAppliedWith100PercentSurcharge() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //set the surcharge
    setNativeSurcharge(FEE_SCALAR, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected discount quantity
    uint256 expectedDiscount = DEFAULT_REPORT_NATIVE_FEE;

    //fee should be twice the surcharge minus the discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE * 2 - expectedDiscount);
  }

  function test_feeIsZeroWith100PercentDiscount() public {
    //set the subscriber discount to 100%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_feeIsUpdatedAfterDiscountIsRemoved() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected discount quantity
    uint256 expectedDiscount = DEFAULT_REPORT_NATIVE_FEE / 2;

    //fee should be 50% of the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);

    //remove the discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), 0, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_feeIsUpdatedAfterNewDiscountIsApplied() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected discount quantity
    uint256 expectedDiscount = DEFAULT_REPORT_NATIVE_FEE / 2;

    //fee should be 50% of the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);

    //change the discount to 25%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 4, ADMIN);

    //get the fee required by the feeManager
    fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //expected discount is now 25%
    expectedDiscount = DEFAULT_REPORT_NATIVE_FEE / 4;

    //fee should be the base fee minus the expected discount
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);
  }

  function test_setDiscountOver100Percent() public {
    //should revert with invalid discount
    vm.expectRevert(INVALID_DISCOUNT_ERROR);

    //set the subscriber discount to over 100%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR + 1, ADMIN);
  }

  function test_surchargeIsNotAppliedWith100PercentDiscount() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 5;

    //set the subscriber discount to 100%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR, ADMIN);

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_nonAdminUserCanNotSetDiscount() public {
    //should revert with unauthorized
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR, USER);
  }

  function test_surchargeFeeRoundsUpWhenUneven() public {
    //native surcharge
    uint256 nativeSurcharge = FEE_SCALAR / 3;

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected surcharge quantity
    uint256 expectedSurcharge = (DEFAULT_REPORT_NATIVE_FEE * nativeSurcharge) / FEE_SCALAR;

    //expected fee should the base fee offset by the expected surcharge
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE + expectedSurcharge + 1);
  }

  function test_discountFeeRoundsDownWhenUneven() public {
    //native surcharge
    uint256 discount = FEE_SCALAR / 3;

    //set the subscriber discount to 33.333%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), discount, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected quantity
    uint256 expectedDiscount = ((DEFAULT_REPORT_NATIVE_FEE * discount) / FEE_SCALAR);

    //expected fee should the base fee offset by the expected surcharge
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscount);
  }

  function test_reportWithNoExpiryOrFeeReturnsZero() public view {
    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV1Report(DEFAULT_FEED_1_V1), getNativeQuote(), USER);

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_correctDiscountIsAppliedWhenBothTokensAreDiscounted() public {
    //set the subscriber and native discounts
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 4, ADMIN);
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager for both tokens
    Common.Asset memory linkFee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);
    Common.Asset memory nativeFee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //calculate the expected discount quantity for each token
    uint256 expectedDiscountLink = (DEFAULT_REPORT_LINK_FEE * FEE_SCALAR) / 4 / FEE_SCALAR;
    uint256 expectedDiscountNative = (DEFAULT_REPORT_NATIVE_FEE * FEE_SCALAR) / 2 / FEE_SCALAR;

    //check the fee calculation for each token
    assertEq(linkFee.amount, DEFAULT_REPORT_LINK_FEE - expectedDiscountLink);
    assertEq(nativeFee.amount, DEFAULT_REPORT_NATIVE_FEE - expectedDiscountNative);
  }

  function test_discountIsNotAppliedToOtherFeeds() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_2_V3), getNativeQuote(), USER);

    //fee should be the base fee
    assertEq(fee.amount, DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_noFeeIsAppliedWhenReportHasZeroFee() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, uint32(block.timestamp), 0, 0),
      getNativeQuote(),
      USER
    );

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_noFeeIsAppliedWhenReportHasZeroFeeAndDiscountAndSurchargeIsSet() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //set the surcharge
    setNativeSurcharge(FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, uint32(block.timestamp), 0, 0),
      getNativeQuote(),
      USER
    );

    //fee should be zero
    assertEq(fee.amount, 0);
  }

  function test_nativeSurchargeEventIsEmittedOnUpdate() public {
    //native surcharge
    uint64 nativeSurcharge = FEE_SCALAR / 3;

    //an event should be emitted
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit NativeSurchargeUpdated(nativeSurcharge);

    //set the surcharge
    setNativeSurcharge(nativeSurcharge, ADMIN);
  }

  function test_getBaseRewardWithLinkQuote() public view {
    //get the fee required by the feeManager
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //the reward should equal the base fee
    assertEq(reward.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_getRewardWithLinkQuoteAndLinkDiscount() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //the reward should equal the discounted base fee
    assertEq(reward.amount, DEFAULT_REPORT_LINK_FEE / 2);
  }

  function test_getRewardWithNativeQuote() public view {
    //get the fee required by the feeManager
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //the reward should equal the base fee in link
    assertEq(reward.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_getRewardWithNativeQuoteAndSurcharge() public {
    //set the native surcharge
    setNativeSurcharge(FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //the reward should equal the base fee in link regardless of the surcharge
    assertEq(reward.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_getRewardWithLinkDiscount() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 2, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //the reward should equal the discounted base fee
    assertEq(reward.amount, DEFAULT_REPORT_LINK_FEE / 2);
  }

  function test_getLinkFeeIsRoundedUp() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 3, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //the reward should equal .66% + 1 of the base fee due to a 33% discount rounded up
    assertEq(fee.amount, (DEFAULT_REPORT_LINK_FEE * 2) / 3 + 1);
  }

  function test_getLinkRewardIsSameAsFee() public {
    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 3, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //check the reward is in link
    assertEq(fee.assetAddress, address(link));

    //the reward should equal .66% of the base fee due to a 33% discount rounded down
    assertEq(reward.amount, fee.amount);
  }

  function test_getLinkRewardWithNativeQuoteAndSurchargeWithLinkDiscount() public {
    //set the native surcharge
    setNativeSurcharge(FEE_SCALAR / 2, ADMIN);

    //set the link discount
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 3, ADMIN);

    //get the fee required by the feeManager
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //the reward should equal the base fee in link regardless of the surcharge
    assertEq(reward.amount, DEFAULT_REPORT_LINK_FEE);
  }

  function test_testRevertIfReportHasExpired() public {
    //expect a revert
    vm.expectRevert(EXPIRED_REPORT_ERROR);

    //get the fee required by the feeManager
    getFee(
      getV3ReportWithCustomExpiryAndFee(
        DEFAULT_FEED_1_V3,
        block.timestamp - 1,
        DEFAULT_REPORT_LINK_FEE,
        DEFAULT_REPORT_NATIVE_FEE
      ),
      getNativeQuote(),
      USER
    );
  }

  function test_discountIsReturnedForLink() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 2, ADMIN);

    //get the fee applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);
  }

  function test_DiscountIsReturnedForNative() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);
  }

  function test_DiscountIsReturnedForNativeWithSurcharge() public {
    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //set the surcharge
    setNativeSurcharge(FEE_SCALAR / 5, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);
  }

  function test_GlobalDiscountWithNative() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(native), FEE_SCALAR / 2, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);
  }

  function test_GlobalDiscountWithLink() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(link), FEE_SCALAR / 2, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);
  }

  function test_GlobalDiscountWithNativeAndLink() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(native), FEE_SCALAR / 2, ADMIN);
    setSubscriberGlobalDiscount(USER, address(link), FEE_SCALAR / 2, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);
  }

  function test_GlobalDiscountIsOverridenByIndividualDiscountNative() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(native), FEE_SCALAR / 2, ADMIN);
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 4, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 4);
  }

  function test_GlobalDiscountIsOverridenByIndividualDiscountLink() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(link), FEE_SCALAR / 2, ADMIN);
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(link), FEE_SCALAR / 4, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 4);
  }

  function test_GlobalDiscountIsUpdatedAfterBeingSetToZeroLink() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(link), FEE_SCALAR / 2, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);

    //set the global discount to zero
    setSubscriberGlobalDiscount(USER, address(link), 0, ADMIN);

    //get the discount applied
    discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getLinkQuote(), USER);

    //fee should be zero
    assertEq(discount, 0);
  }

  function test_GlobalDiscountIsUpdatedAfterBeingSetToZeroNative() public {
    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(native), FEE_SCALAR / 2, ADMIN);

    //get the discount applied
    uint256 discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be half the default
    assertEq(discount, FEE_SCALAR / 2);

    //set the global discount to zero
    setSubscriberGlobalDiscount(USER, address(native), 0, ADMIN);

    //get the discount applied
    discount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    //fee should be zero
    assertEq(discount, 0);
  }

  function test_GlobalDiscountCantBeSetToMoreThanMaximum() public {
    //should revert with invalid discount
    vm.expectRevert(INVALID_DISCOUNT_ERROR);

    //set the global discount to 101%
    setSubscriberGlobalDiscount(USER, address(native), FEE_SCALAR + 1, ADMIN);
  }

  function test_onlyOwnerCanSetGlobalDiscount() public {
    //should revert with unauthorized
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //set the global discount to 50%
    setSubscriberGlobalDiscount(USER, address(native), FEE_SCALAR / 2, USER);
  }
}
