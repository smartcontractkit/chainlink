// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../dev/FeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import "./BaseFeeManager.t.sol";

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
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    //should revert as the user is not the owner
    vm.expectRevert(UNAUTHORIZED_ERROR);

    //process the fee
    processFee(payload, USER, 0, USER);
  }

  function test_processFeeAsProxy() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, PROXY);

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_processFeeIfSubscriberIsSelf() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    //expect a revert due to funds being unapproved
    vm.expectRevert(INVALID_ADDRESS_ERROR);

    //process the fee will attempt to transfer link from the contract to the rewardManager, which won't be approved
    processFee(payload, address(feeManager), 0, ADMIN);
  }

  function test_processFeeWithWithEmptyQuotePayload() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), bytes(""));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link by default
    processFee(payload, USER, 0, ADMIN);
  }

  function test_processFeeWithWithZeroQuotePayload() public {
    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(INVALID_ADDRESS));

    //expect a revert as the quote is invalid
    vm.expectRevert(INVALID_QUOTE_ERROR);

    //processing the fee will transfer the link by default
    processFee(payload, USER, 0, ADMIN);
  }

  function test_processFeeWithWithCorruptQuotePayload() public {
    //get the default payload
    bytes memory payload = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV2Report(DEFAULT_FEED_1_V3),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    //expect an evm revert as the quote is corrupt
    vm.expectRevert();

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, 0, ADMIN);
  }

  function test_processFeeDefaultReportsStillVerifiesWithEmptyQuote() public {
    //get the default payload
    bytes memory payload = getPayload(getV0Report(DEFAULT_FEED_1_V1), bytes(""));

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, PROXY);
  }

  function test_processFeeWithDefaultReportPayloadAndQuoteStillVerifies() public {
    //get the default payload
    bytes memory payload = getPayload(getV0Report(DEFAULT_FEED_1_V1), getQuotePayload(getLinkAddress()));

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, 0, ADMIN);
  }

  function test_processFeeNative() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(payload, USER, 0, ADMIN);

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
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    //approve the native to be transferred from the user
    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    //expect an emit as there's not enough link
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit InsufficientLink(DEFAULT_CONFIG_DIGEST, DEFAULT_REPORT_LINK_FEE, DEFAULT_REPORT_NATIVE_FEE);

    //processing the fee will transfer the native from the user to the feeManager
    processFee(payload, USER, 0, ADMIN);

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
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    //only the proxy or admin can call processFee, they will pass in the native value on the users behalf
    processFee(payload, USER, DEFAULT_REPORT_NATIVE_FEE, PROXY);

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

  function test_processFeeWithUnwrappedNativeShortFunds() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    //expect a revert as not enough funds
    vm.expectRevert(INVALID_DEPOSIT_ERROR);

    //only the proxy or admin can call processFee, they will pass in the native value on the users behalf
    processFee(payload, USER, DEFAULT_REPORT_NATIVE_FEE - 1, PROXY);
  }

  function test_processFeeWithUnwrappedNativeLinkAddress() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    //expect a revert as not enough funds
    vm.expectRevert(INVALID_DEPOSIT_ERROR);

    //only the proxy or admin can call processFee, they will pass in the native value on the users behalf
    processFee(payload, USER, DEFAULT_REPORT_NATIVE_FEE - 1, PROXY);
  }

  function test_processFeeWithUnwrappedNativeWithExcessiveFee() public {
    //simulate a deposit of link for the conversion pool
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE);

    //get the default payload
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    //call processFee from the proxy to test whether the funds are returned to the subscriber. In reality, the funds would be returned to the caller of the proxy.
    processFee(payload, PROXY, DEFAULT_REPORT_NATIVE_FEE * 2, PROXY);

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
    bytes memory payload = getPayload(getV2Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, PROXY);

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
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, PROXY);
  }

  function test_V2PayloadVerifies() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V2), getQuotePayload(getLinkAddress()));

    //approve the link to be transferred from the from the subscriber to the rewardManager
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, ADMIN);

    //check the link has been transferred
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);

    //check the user has had the link fee deducted
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_V2PayloadWithoutQuoteFails() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V2), bytes(""));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, ADMIN);
  }

  function test_V2PayloadWithoutZeroFee() public {
    //get the default payload
    bytes memory payload = getPayload(getV1Report(DEFAULT_FEED_1_V2), getQuotePayload(getLinkAddress()));

    //expect a revert as the quote is invalid
    vm.expectRevert();

    //processing the fee will transfer the link from the user to the rewardManager
    processFee(payload, USER, 0, ADMIN);
  }

  function test_processFeeWithInvalidReportVersion() public {
    bytes memory data = abi.encode(0x0000100000000000000000000000000000000000000000000000000000000000);

    //get the default payload
    bytes memory payload = getPayload(data, getQuotePayload(getLinkAddress()));

    //version is invalid as feedId is 0
    vm.expectRevert(INVALID_REPORT_VERSION_ERROR);

    //processing the fee will not withdraw anything as there is no fee to collect
    processFee(payload, USER, 0, ADMIN);
  }
}
