// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../dev/FeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import "./BaseFeeManager.t.sol";
import {IRewardManager} from "../../dev/interfaces/IRewardManager.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager processFee
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  uint256 internal constant NUMBER_OF_REPORTS = 5;

  function setUp() public override {
    super.setUp();
  }

  function test_processMultipleLinkReports() public {
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS, USER);

    processFee(payloads, USER, DEFAULT_NATIVE_MINT_QUANTITY);

    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 0);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);

    //the subscriber (user) should receive funds back and not the proxy, although when live the proxy will forward the funds sent and not cover it seen here
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processMultipleWrappedNativeReports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS, USER);

    processFee(payloads, USER, 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 1);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
  }

  function test_processMultipleUnwrappedNativeReports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    processFee(payloads, USER, DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS * 2);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
  }

  function test_processMultipleLinkAndNativeWrappedReports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * 2 + 1);

    bytes memory nativePayload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));
    bytes memory linkPayload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = linkPayload;
    payloads[1] = linkPayload;
    payloads[2] = linkPayload;
    payloads[3] = nativePayload;
    payloads[4] = nativePayload;

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 2, USER);
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 3, USER);

    processFee(payloads, USER, 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 5);
    assertEq(getLinkBalance(address(feeManager)), 1);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 3);
  }

  function test_processMultipleLinkAndNativeUnwrappedReports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * 2 + 1);

    bytes memory nativePayload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));
    bytes memory linkPayload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = linkPayload;
    payloads[1] = linkPayload;
    payloads[2] = linkPayload;
    payloads[3] = nativePayload;
    payloads[4] = nativePayload;

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 3, USER);

    processFee(payloads, USER, DEFAULT_REPORT_NATIVE_FEE * 4);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 5);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 3);
  }

  function test_processV1V2V3Reports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * 2 + 1);

    bytes memory payloadV1 = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes memory nativePayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2), getQuotePayload(getNativeAddress()));
    bytes memory linkPayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2), getQuotePayload(getLinkAddress()));

    bytes memory nativePayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));
    bytes memory linkPayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = payloadV1;
    payloads[1] = nativePayloadV2;
    payloads[2] = linkPayloadV2;
    payloads[3] = nativePayloadV3;
    payloads[4] = linkPayloadV3;

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 2, USER);
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 2, USER);

    processFee(payloads, USER, 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 4);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 2);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
  }

  function test_processV1V2V3ReportsWithUnwrapped() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * 2 + 1);

    bytes memory payloadV1 = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes memory nativePayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2), getQuotePayload(getNativeAddress()));
    bytes memory linkPayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2), getQuotePayload(getLinkAddress()));

    bytes memory nativePayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));
    bytes memory linkPayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = payloadV1;
    payloads[1] = nativePayloadV2;
    payloads[2] = linkPayloadV2;
    payloads[3] = nativePayloadV3;
    payloads[4] = linkPayloadV3;

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 2, USER);

    processFee(payloads, USER, DEFAULT_REPORT_NATIVE_FEE * 4);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 4);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 2);
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processMultipleV1Reports() public {
    bytes memory payload = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    processFee(payloads, USER, DEFAULT_REPORT_NATIVE_FEE * 5);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_eventIsEmittedIfNotEnoughLink() public {
    bytes memory nativePayload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getNativeAddress()));
    bytes memory linkPayload = getPayload(getV3Report(DEFAULT_FEED_1_V3), getQuotePayload(getLinkAddress()));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = linkPayload;
    payloads[1] = linkPayload;
    payloads[2] = linkPayload;
    payloads[3] = nativePayload;
    payloads[4] = nativePayload;

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 2, USER);
    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 3, USER);

    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](2);
    payments[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[1] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    vm.expectEmit();

    emit InsufficientLink(payments);

    processFee(payloads, USER, 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 3);
  }
}
