// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import "./BaseDestinationFeeManager.t.sol";
import {IDestinationRewardManager} from "../../interfaces/IDestinationRewardManager.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager processFee
 */
contract DestinationFeeManagerProcessFeeTest is BaseDestinationFeeManagerTest {
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

    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS, USER);

    processFee(poolIds, payloads, USER, address(link), DEFAULT_NATIVE_MINT_QUANTITY);

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

    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS, USER);

    processFee(poolIds, payloads, USER, address(native), 0);

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

    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    processFee(poolIds, payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS * 2);

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

    bytes32[] memory poolIds = new bytes32[](5);
    for (uint256 i = 0; i < 5; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    processFee(poolIds, payloads, USER, address(link), 0);

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

    bytes32[] memory poolIds = new bytes32[](5);
    for (uint256 i = 0; i < 5; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    processFee(poolIds, payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * 4);

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

    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    processFee(poolIds, payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * 5);

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

    IDestinationRewardManager.FeePayment[] memory payments = new IDestinationRewardManager.FeePayment[](5);
    payments[0] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[1] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[2] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[3] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));
    payments[4] = IDestinationRewardManager.FeePayment(DEFAULT_CONFIG_DIGEST, uint192(DEFAULT_REPORT_LINK_FEE));

    vm.expectEmit();

    bytes32[] memory poolIds = new bytes32[](5);
    for (uint256 i = 0; i < 5; ++i) {
      poolIds[i] = payments[i].poolId;
    }

    emit InsufficientLink(payments);

    processFee(poolIds, payloads, USER, address(native), 0);

    assertEq(getNativeBalance(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(getNativeBalance(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY);
  }

  function test_processPoolIdsPassedMismatched() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    // poolIds passed are different that number of reports in payload
    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS - 1);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS - 1; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    vm.expectRevert(POOLID_MISMATCH_ERROR);
    processFee(poolIds, payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS * 2);
  }

  function test_poolIdsCannotBeZeroAddress() public {
    mintLink(address(feeManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS + 1);

    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }

    poolIds[2] = 0x000;
    vm.expectRevert(INVALID_ADDRESS_ERROR);
    processFee(poolIds, payloads, USER, address(native), DEFAULT_REPORT_NATIVE_FEE * NUMBER_OF_REPORTS * 2);
  }

  function test_rewardsAreCorrectlySentToEachAssociatedPoolWhenVerifyingInBulk() public {
    bytes memory payload = getPayload(getV3Report(DEFAULT_FEED_1_V3));

    bytes[] memory payloads = new bytes[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      payloads[i] = payload;
    }

    bytes32[] memory poolIds = new bytes32[](NUMBER_OF_REPORTS);
    for (uint256 i = 0; i < NUMBER_OF_REPORTS - 1; ++i) {
      poolIds[i] = DEFAULT_CONFIG_DIGEST;
    }
    poolIds[NUMBER_OF_REPORTS - 1] = DEFAULT_CONFIG_DIGEST2;

    approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS, USER);

    // Checking no rewards yet for each pool
    for (uint256 i = 0; i < NUMBER_OF_REPORTS; ++i) {
      bytes32 p_id = poolIds[i];
      uint256 poolDeficit = rewardManager.s_totalRewardRecipientFees(p_id);
      assertEq(poolDeficit, 0);
    }

    processFee(poolIds, payloads, USER, address(link), DEFAULT_NATIVE_MINT_QUANTITY);

    assertEq(getLinkBalance(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);
    assertEq(getLinkBalance(address(feeManager)), 0);
    assertEq(getLinkBalance(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * NUMBER_OF_REPORTS);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY);
    assertEq(PROXY.balance, DEFAULT_NATIVE_MINT_QUANTITY);

    // Checking each pool got the correct rewards
    uint256 expectedRewards = DEFAULT_REPORT_LINK_FEE * (NUMBER_OF_REPORTS - 1);
    uint256 poolRewards = rewardManager.s_totalRewardRecipientFees(DEFAULT_CONFIG_DIGEST);
    assertEq(poolRewards, expectedRewards);

    expectedRewards = DEFAULT_REPORT_LINK_FEE;
    poolRewards = rewardManager.s_totalRewardRecipientFees(DEFAULT_CONFIG_DIGEST2);
    assertEq(poolRewards, expectedRewards);
  }
}
