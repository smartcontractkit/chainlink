// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../FeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import "./BaseFeeManager.t.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager processFee
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  function setUp() public override {
    super.setUp();
  }

  function test_nonAdminProxyUserCannotProcessFee() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //should revert as the user is not the owner
    vm.expectRevert(UNAUTHORIZED_ERROR);

    //process the fee
    ProcessFeeAsUser(payload, USER, address(link), 0, USER);
  }

  function test_processFeeAsProxy() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link), 0);

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_processFeeIfSubscriberIsSelf() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert due to the feeManager being the subscriber
    vm.expectRevert(INVALID_ADDRESS_ERROR);

    //process the fee will fail due to assertion
    processFee(payload, address(feeManager), address(native), 0);
  }

  function test_processFeeWithWithEmptyQuotePayload() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link by default
    processFee(payload, USER, address(0), 0);
  }

  function test_processFeeWithWithZeroQuotePayload() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as the quote is invalid
    vm.expectRevert(INVALID_QUOTE_ERROR);

    //processing the fee will transfer the link by default
    processFee(payload, USER, INVALID_ADDRESS, 0);
  }

  function test_processFeeWithWithCorruptQuotePayload() public {
    //get the default payload
    bytes memory payload = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV3Report(DEFAULT_FEED_1_V3),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    //expect an evm revert as the quote is corrupt
    vm.expectRevert();

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, address(link), 0);
  }

  function test_processFeeDefaultReportsStillVerifiesWithEmptyQuote() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V1));

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(0), 0);
  }

  function test_processFeeWithDefaultReportPayloadAndQuoteStillVerifies() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V1));

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, address(link), 0);
  }

  function test_processFeeNative() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(payload, USER, address(native), 0);

    //check the native has been transferred
    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManager has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManager)), 0);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeEmitsEventIfNotEnoughLink() public {
    //simulate a deposit of half the link required for the fee
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE / 2);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //expect an emit as there's not enough link
    vm.expectEmit();

    IRewardManager.FeePayment[] memory contractFees = new IRewardManager.FeePayment[](1);
    contractFees[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    //emit the event that is expected to be emitted
    emit InsufficientLink(contractFees);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(payload, USER, address(native), 0);

    //check the native has been transferred
    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE);

    //check no link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), 0);
    assertEq(getLinkBalance(address(feeManager)), DEFAULT_REPORT_LINK_FEE / 2);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeWithUnwrappedNative() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //only the proxy or admin can call processFee, they will pass in the native value on the users behalf
    processFee(payload, USER, address(native), DEFAULT_REPORT_NATIVE_FEE);

    //check the native has been transferred and converted to wrapped native
    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE);
    assertEq(getNativeUnwrappedBalance(address(feeManager)), 0);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManager has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManager)), 0);

    //check the subscriber has had the native deducted
    assertEq(getNativeUnwrappedBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeWithUnwrappedNativeShortFunds() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as not enough funds
    vm.expectRevert(INVALID_DEPOSIT_ERROR);

    //only the proxy or admin can call processFee, they will pass in the native value on the users behalf
    processFee(payload, USER, address(native), DEFAULT_REPORT_NATIVE_FEE - 1);
  }

  function test_processFeeWithUnwrappedNativeLinkAddress() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as not enough funds
    vm.expectRevert(INSUFFICIENT_ALLOWANCE_ERROR);

    //the change will be returned and the user will attempted to be billed in LINK
    processFee(payload, USER, address(link), DEFAULT_REPORT_NATIVE_FEE - 1);
  }

  function test_processFeeWithUnwrappedNativeLinkAddressExcessiveFee() public {
    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, PROXY);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //call processFee from the proxy to test whether the funds are returned to the subscriber. In reality, the funds would be returned to the caller of the proxy.
    processFee(payload, PROXY, address(link), DEFAULT_REPORT_NATIVE_FEE);

    //check the native unwrapped is no longer in the account
    assertEq(getNativeBalance(address(feeManager)), 0);
    assertEq(getNativeUnwrappedBalance(address(feeManager)), 0);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManager has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManager)), 0);

    //native should not be deducted
    assertEq(getNativeUnwrappedBalance(PROXY), DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processFeeWithUnwrappedNativeWithExcessiveFee() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //call processFee from the proxy to test whether the funds are returned to the subscriber. In reality, the funds would be returned to the caller of the proxy.
    processFee(payload, PROXY, address(native), DEFAULT_REPORT_NATIVE_FEE * 2);

    //check the native has been transferred and converted to wrapped native
    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE);
    assertEq(getNativeUnwrappedBalance(address(feeManager)), 0);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManager has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManager)), 0);

    //check the subscriber has had the native deducted
    assertEq(getNativeUnwrappedBalance(PROXY), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeUsesCorrectDigest() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link), 0);

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);

    //check funds have been paid to the reward manager
    assertEq(rewardManager.s_totalRewardRecipientFees(DEFAULT_CONFIG_DIGEST), DEFAULT_REPORT_LINK_FEE);
  }

  function test_V1PayloadVerifies() public {
    //replicate a default payload
    bytes memory payload = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV2Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(0), 0);
  }

  function test_V2PayloadVerifies() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V2));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link), 0);

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_V2PayloadWithoutQuoteFails() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V2));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(0), 0);
  }

  function test_V2PayloadWithoutZeroFee() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V2));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link), 0);
  }

  function test_processFeeWithInvalidReportVersionFailsToDecode() public {
    bytes memory data = abi.encode(0x0000100000000000000000000000000000000000000000000000000000000000);

    //get the default payload
    bytes memory payload = getPayload(data);

    //serialization will fail as there is no report to decode
    vm.expectRevert();

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, address(link), 0);
  }

  function test_processFeeWithZeroNativeNonZeroLinkWithNativeQuote() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, DEFAULT_REPORT_LINK_FEE, 0)
    );

    //call processFee should not revert as the fee is 0
    processFee(payload, PROXY, address(native), 0);
  }

  function test_processFeeWithZeroNativeNonZeroLinkWithLinkQuote() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, DEFAULT_REPORT_LINK_FEE, 0)
    );

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link to the rewardManager from the user
    processFee(payload, USER, address(link), 0);

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_processFeeWithZeroLinkNonZeroNativeWithNativeQuote() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, 0, DEFAULT_REPORT_NATIVE_FEE)
    );

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(payload, USER, address(native), 0);

    //check the native has been transferred
    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE);

    //check no link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), 0);

    //check the feeManager has had no link deducted
    assertEq(getLinkBalance(address(feeManager)), DEFAULT_REPORT_LINK_FEE);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeWithZeroLinkNonZeroNativeWithLinkQuote() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, 0, DEFAULT_REPORT_NATIVE_FEE)
    );

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(link), 0);
  }

  function test_processFeeWithZeroNativeNonZeroLinkReturnsChange() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, 0, DEFAULT_REPORT_NATIVE_FEE)
    );

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(link), DEFAULT_REPORT_NATIVE_FEE);

    //check the change has been returned
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_V1PayloadVerifiesAndReturnsChange() public {
    //emulate a V1 payload with no quote
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V1));

    processFee(payload, USER, address(0), DEFAULT_REPORT_NATIVE_FEE);

    //Fee manager should not contain any native
    assertEq(address(feeManager).balance, 0);
    assertEq(getNativeBalance(address(feeManager)), 0);

    //check the unused native passed in is returned
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processFeeWithDiscountEmitsEvent() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE / 2, USER);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);
    uint256 appliedDiscount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    vm.expectEmit();

    emit DiscountApplied(DEFAULT_CONFIG_DIGEST, USER, fee, reward, appliedDiscount);

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(native), 0);
  }

  function test_processFeeWithNoDiscountDoesNotEmitEvent() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(native), 0);

    //no logs should have been emitted
    assertEq(vm.getRecordedLogs().length, 0);
  }
}
