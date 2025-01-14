// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import "./BaseFeeManager.t.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager processFee
 */
contract FeeManagerProcessFeeTestV05 is BaseFeeManagerTest {
  uint256 internal constant NUMBER_OF_REPORTS = 5;

  function setUp() public override {
    super.setUp();
  }

  function test_processMultipleLinkReports() public {
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS, USER);

    processFee(payloads, USER, address(link), DEFAULT_NATIVE_MINT_QUANTITY);

    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 0);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);

    //the subscriber (user) should receive funds back and not the proxy, although when live the proxy will forward the funds sent and not cover it seen here
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_processMultipleWrappedNativeReports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS, USER);

    processFee(payloads, USER, address(native), 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 1);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
  }

  function test_processMultipleUnwrappedNativeReports() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    processFee(payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS * 2);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS);
  }

  function test_processV1V2V3Reports() public {
    mintLink(address(feeManager), 1);

    bytes memory payloadV1 = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes memory linkPayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2));
    bytes memory linkPayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = payloadV1;
    payloads[1] = linkPayloadV2;
    payloads[2] = linkPayloadV2;
    payloads[3] = linkPayloadV3;
    payloads[4] = linkPayloadV3;

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 4, USER);

    processFee(payloads, USER, address(link), 0);

    assertEq(getNativeBalance(address(feeManager)), 0);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 4);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 4);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - 0);
  }

  function test_processV1V2V3ReportsWithUnwrapped() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * 4 + 1);

    bytes memory payloadV1 = abi.encode(
      [DEFAULT_CONFIG_DIGEST, 0, 0],
      getV1Report(DEFAULT_FEED_1_V1),
      new bytes32[](1),
      new bytes32[](1),
      bytes32("")
    );

    bytes memory nativePayloadV2 = getPayload(getV2Report(DEFAULT_FEED_1_V2));
    bytes memory nativePayloadV3 = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = payloadV1;
    payloads[1] = nativePayloadV2;
    payloads[2] = nativePayloadV2;
    payloads[3] = nativePayloadV3;
    payloads[4] = nativePayloadV3;

    processFee(payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * 4);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 4);
    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 4);
    assertEq(getLinkBalance(address(feeManager)), 1);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 4);
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

    processFee(payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * 5);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_eventIsEmittedIfNotEnoughLink() public {
    bytes memory nativePayload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](5);
    payloads[0] = nativePayload;
    payloads[1] = nativePayload;
    payloads[2] = nativePayload;
    payloads[3] = nativePayload;
    payloads[4] = nativePayload;

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 5, USER);

    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](5);
    payments[0] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[1] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[2] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[3] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[4] = IRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    vm.expectEmit();

    emit InsufficientLink(payments);

    processFee(payloads, USER, address(native), 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY);
  }
}
