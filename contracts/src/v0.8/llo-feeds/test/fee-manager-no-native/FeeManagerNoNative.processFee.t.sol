// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {Common} from "../../libraries/Common.sol";

import "./BaseFeeManagerNoNative.t.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

/**
 * @title FeeManagerNoNativeProcessFeedTest
 * @author Michael Fletcher
 * @author ad0ll
 * @notice This contract will test the functionality of the FeeManagerNoNative processFee function
 */
contract FeeManagerNoNativeProcessFeeTest is BaseFeeManagerNoNativeTest {
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
    processFee(payload, USER, address(link));

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_processFeeIfSubscriberIsSelf() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert due to the feeManagerNoNative being the subscriber
    vm.expectRevert(INVALID_ADDRESS_ERROR);

    //process the fee will fail due to assertion
    processFee(payload, address(feeManagerNoNative), address(native));
  }

  function test_processFeeWithWithEmptyQuotePayload() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as the quote is invalid
    vm.expectRevert(); //TODO check what this revert this should be

    //processing the fee will transfer the link by default
    processFee(payload, USER, address(0));
  }

  function test_processFeeWithWithZeroQuotePayload() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as the quote is invalid
    vm.expectRevert(INVALID_QUOTE_ERROR);

    //processing the fee will transfer the link by default
    processFee(payload, USER, INVALID_ADDRESS);
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
    processFee(payload, USER, address(link));
  }

  function test_processFeeDefaultReportsStillVerifiesWithEmptyQuote() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V1));

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(0));
  }

  function test_processFeeWithDefaultReportPayloadAndQuoteStillVerifies() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V1));

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, address(link));
  }

  function test_processFeeNative() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the native to be transferred from the user
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE, USER);

    //processing the fee will transfer the native from the user to the feeManagerNoNative contract
    processFee(payload, USER, address(native));

    //check the native has been transferred
    assertEq(getNativeBalance(address(feeManagerNoNative)), DEFAULT_REPORT_NATIVE_FEE);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManagerNoNative has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManagerNoNative)), 0);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeEmitsEventIfNotEnoughLink() public {
    //simulate a deposit of half the link required for the fee
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE / 2);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the native to be transferred from the user
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE, USER);

    //expect an emit as there's not enough link
    vm.expectEmit();

    IRewardManager.FeePayment[] memory contractFees = new IRewardManager.FeePayment[](1);
    contractFees[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    //emit the event that is expected to be emitted
    emit InsufficientLink(contractFees);

    //processing the fee will transfer the native from the user to the feeManagerNoNative contract
    processFee(payload, USER, address(native));

    //check the native has been transferred
    assertEq(getNativeBalance(address(feeManagerNoNative)), DEFAULT_REPORT_NATIVE_FEE);

    //check no link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), 0);
    assertEq(getLinkBalance(address(feeManagerNoNative)), DEFAULT_REPORT_LINK_FEE / 2);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeWithWrappedNative() public {
    // Mint and approve LINK for feeManagerNoNative transfer to rewardManager following successful verification
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, address(feeManagerNoNative));

    // Mint and approve ERC20 representing native for the verification fee. USER given native ERC20 in setup
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE, USER);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //only the proxy or admin can call processFee, they will pass in the native value on the users behalf
    processFee(payload, USER, address(native));

    //Check that fee manager has received wrapped native
    assertEq(getNativeBalance(address(feeManagerNoNative)), DEFAULT_REPORT_NATIVE_FEE);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManagerNoNative has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManagerNoNative)), 0);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeWithUnwrappedNativeLinkAddress() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //expect a revert as not enough funds
    vm.expectRevert(INSUFFICIENT_ALLOWANCE_ERROR);

    //the change will be returned and the user will attempted to be billed in LINK
    processFee(payload, USER, address(link));
  }

  function test_processFeeWithWrappedNativeLinkAddressExcessiveFee() public {
    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, PROXY);

    // Mint native so we can check balance at end, no need to approve since transfer shouldn't happen
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //call processFee from the proxy to test whether the funds are returned to the subscriber. In reality, the funds would be returned to the caller of the proxy.
    processFee(payload, PROXY, address(link));

    //check the native unwrapped is no longer in the account
    assertEq(getNativeBalance(address(feeManagerNoNative)), 0);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManagerNoNative has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManagerNoNative)), 0);

    //native should not be deducted, note that native was minted for PROXY in setup
    assertEq(getNativeBalance(PROXY), DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processFeeWithUnwrappedNativeWithExcessiveFee() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE, PROXY);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //call processFee from the proxy to test whether the funds are returned to the subscriber. In reality, the funds would be returned to the caller of the proxy.
    processFee(payload, PROXY, address(native));

    //check the native has been transferred and converted to wrapped native
    assertEq(getNativeBalance(address(feeManagerNoNative)), DEFAULT_REPORT_NATIVE_FEE);
    assertEq(getNativeUnwrappedBalance(address(feeManagerNoNative)), 0);

    //check the link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the feeManagerNoNative has had the link deducted, the remaining balance should be 0
    assertEq(getLinkBalance(address(feeManagerNoNative)), 0);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(PROXY), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeUsesCorrectDigest() public {
    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link));

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
    processFee(payload, USER, address(0));
  }

  function test_V2PayloadVerifies() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V2));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link));

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
    processFee(payload, USER, address(0));
  }

  function test_V2PayloadWithoutZeroFee() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V2));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, address(link));
  }

  function test_processFeeWithInvalidReportVersionFailsToDecode() public {
    bytes memory data = abi.encode(0x0000100000000000000000000000000000000000000000000000000000000000);

    //get the default payload
    bytes memory payload = getPayload(data);

    //serialization will fail as there is no report to decode
    vm.expectRevert();

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, address(link));
  }

  function test_processFeeWithZeroNativeNonZeroLinkWithNativeQuote() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, DEFAULT_REPORT_LINK_FEE, 0)
    );

    //call processFee should not revert as the fee is 0
    processFee(payload, PROXY, address(native));
  }

  function test_processFeeWithZeroNativeNonZeroLinkWithLinkQuote() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, DEFAULT_REPORT_LINK_FEE, 0)
    );

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link to the rewardManager from the user
    processFee(payload, USER, address(link));

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_processFeeWithZeroLinkNonZeroNativeWithNativeQuote() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, 0, DEFAULT_REPORT_NATIVE_FEE)
    );

    //approve the native to be transferred from the user
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE, USER);

    //processing the fee will transfer the native from the user to the feeManagerNoNative contract
    processFee(payload, USER, address(native));

    //check the native has been transferred
    assertEq(getNativeBalance(address(feeManagerNoNative)), DEFAULT_REPORT_NATIVE_FEE);

    //check no link has been transferred to the rewardManager
    assertEq(getLinkBalance(address(rewardManager)), 0);

    //check the feeManagerNoNative has had no link deducted
    assertEq(getLinkBalance(address(feeManagerNoNative)), DEFAULT_REPORT_LINK_FEE);

    //check the subscriber has had the native deducted
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
  }

  function test_processFeeWithZeroLinkNonZeroNativeWithLinkQuote() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, 0, DEFAULT_REPORT_NATIVE_FEE)
    );

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(link));
  }

  function test_processFeeWithZeroNativeNonZeroLinkReturnsChange() public {
    //get the default payload
    bytes memory payload = getPayload(
      getV3ReportWithCustomExpiryAndFee(DEFAULT_FEED_1_V3, block.timestamp, 0, DEFAULT_REPORT_NATIVE_FEE)
    );

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(link));

    //check the change has been returned
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_blocksNativeBilling() public {
    //emulate a V1 payload with no quote
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));
    //record the current address and switch to the recipient
    changePrank(USER);

    //Expect revert since no native allowed
    vm.expectRevert(NATIVE_BILLING_DISALLOWED);

    //process the fee
    feeManagerProxy.processFee{value: 1 wei}(payload, abi.encode(address(0)));
  }

  function test_processFeeWithDiscountEmitsEvent() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);

    //set the subscriber discount to 50%
    setSubscriberDiscount(USER, DEFAULT_FEED_1_V3, address(native), FEE_SCALAR / 2, ADMIN);

    //approve the native to be transferred from the user
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE / 2, USER);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    Common.Asset memory fee = getFee(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);
    Common.Asset memory reward = getReward(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);
    uint256 appliedDiscount = getAppliedDiscount(getV3Report(DEFAULT_FEED_1_V3), getNativeQuote(), USER);

    vm.expectEmit();

    emit DiscountApplied(DEFAULT_CONFIG_DIGEST, USER, fee, reward, appliedDiscount);

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(native));
  }

  function test_processFeeWithNoDiscountDoesNotEmitEvent() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManagerNoNative), DEFAULT_REPORT_LINK_FEE);

    //approve the native to be transferred from the user
    approveNative(address(feeManagerNoNative), DEFAULT_REPORT_NATIVE_FEE, USER);

    //get the default payload
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    //call processFee should not revert as the fee is 0
    processFee(payload, USER, address(native));

    //no logs should have been emitted
    assertEq(vm.getRecordedLogs().length, 0);
  }
}
