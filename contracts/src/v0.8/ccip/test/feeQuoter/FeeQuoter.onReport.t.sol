// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {KeystoneFeedsPermissionHandler} from "../../../keystone/KeystoneFeedsPermissionHandler.sol";

import {KeystoneForwarder} from "../../../keystone/KeystoneForwarder.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_onReport is FeeQuoterSetup {
  bytes32 internal constant EXECUTION_ID = hex"6d795f657865637574696f6e5f69640000000000000000000000000000000000";
  address internal constant TRANSMITTER = address(50);
  bytes32 internal constant WORKFLOW_ID_1 = hex"6d795f6964000000000000000000000000000000000000000000000000000000";
  address internal constant WORKFLOW_OWNER_1 = address(51);
  bytes10 internal constant WORKFLOW_NAME_1 = hex"000000000000DEADBEEF";
  bytes2 internal constant REPORT_NAME_1 = hex"0001";
  address internal s_onReportTestToken1;
  address internal s_onReportTestToken2;
  bytes public encodedPermissionsMetadata;
  KeystoneForwarder internal s_forwarder;

  function setUp() public virtual override {
    super.setUp();

    s_forwarder = new KeystoneForwarder();

    s_onReportTestToken1 = s_sourceTokens[0];
    s_onReportTestToken2 = _deploySourceToken("onReportTestToken2", 0, 20);

    KeystoneFeedsPermissionHandler.Permission[] memory permissions = new KeystoneFeedsPermissionHandler.Permission[](1);
    permissions[0] = KeystoneFeedsPermissionHandler.Permission({
      forwarder: address(s_forwarder),
      workflowOwner: WORKFLOW_OWNER_1,
      workflowName: WORKFLOW_NAME_1,
      reportName: REPORT_NAME_1,
      isAllowed: true
    });

    encodedPermissionsMetadata = abi.encodePacked(WORKFLOW_ID_1, WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);

    FeeQuoter.TokenPriceFeedUpdate[] memory tokenPriceFeeds = new FeeQuoter.TokenPriceFeedUpdate[](2);
    tokenPriceFeeds[0] = FeeQuoter.TokenPriceFeedUpdate({
      sourceToken: s_onReportTestToken1,
      feedConfig: FeeQuoter.TokenPriceFeedConfig({dataFeedAddress: address(0x0), tokenDecimals: 18, isEnabled: true})
    });
    tokenPriceFeeds[1] = FeeQuoter.TokenPriceFeedUpdate({
      sourceToken: s_onReportTestToken2,
      feedConfig: FeeQuoter.TokenPriceFeedConfig({dataFeedAddress: address(0x0), tokenDecimals: 20, isEnabled: true})
    });
    s_feeQuoter.setReportPermissions(permissions);
    s_feeQuoter.updateTokenPriceFeeds(tokenPriceFeeds);
  }

  function test_onReport() public {
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](2);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp)});
    report[1] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_onReportTestToken2, price: 4e18, timestamp: uint32(block.timestamp)});

    uint224 expectedStoredToken1Price = s_feeQuoter.calculateRebasedValue(18, 18, report[0].price);
    uint224 expectedStoredToken2Price = s_feeQuoter.calculateRebasedValue(18, 20, report[1].price);
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_onReportTestToken1, expectedStoredToken1Price, block.timestamp);
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_onReportTestToken2, expectedStoredToken2Price, block.timestamp);

    changePrank(address(s_forwarder));
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));

    vm.assertEq(s_feeQuoter.getTokenPrice(report[0].token).value, expectedStoredToken1Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(report[0].token).timestamp, report[0].timestamp);

    vm.assertEq(s_feeQuoter.getTokenPrice(report[1].token).value, expectedStoredToken2Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(report[1].token).timestamp, report[1].timestamp);
  }

  function test_onReport_withKeystoneForwarderContract() public {
    FeeQuoter.ReceivedCCIPFeedReport[] memory priceReportRaw = new FeeQuoter.ReceivedCCIPFeedReport[](2);
    priceReportRaw[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp)});
    priceReportRaw[1] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_onReportTestToken2, price: 4e18, timestamp: uint32(block.timestamp)});

    uint224 expectedStoredToken1Price = s_feeQuoter.calculateRebasedValue(18, 18, priceReportRaw[0].price);
    uint224 expectedStoredToken2Price = s_feeQuoter.calculateRebasedValue(18, 20, priceReportRaw[1].price);

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_onReportTestToken1, expectedStoredToken1Price, block.timestamp);
    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_onReportTestToken2, expectedStoredToken2Price, block.timestamp);

    changePrank(address(s_forwarder));
    s_forwarder.route(
      EXECUTION_ID, TRANSMITTER, address(s_feeQuoter), encodedPermissionsMetadata, abi.encode(priceReportRaw)
    );

    vm.assertEq(s_feeQuoter.getTokenPrice(priceReportRaw[0].token).value, expectedStoredToken1Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(priceReportRaw[0].token).timestamp, priceReportRaw[0].timestamp);

    vm.assertEq(s_feeQuoter.getTokenPrice(priceReportRaw[1].token).value, expectedStoredToken2Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(priceReportRaw[1].token).timestamp, priceReportRaw[1].timestamp);
  }

  function test_OnReport_SkipPriceUpdateWhenStaleUpdateReceived() public {
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp)});

    uint224 expectedStoredTokenPrice = s_feeQuoter.calculateRebasedValue(18, 18, report[0].price);

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_onReportTestToken1, expectedStoredTokenPrice, block.timestamp);

    changePrank(address(s_forwarder));
    //setting the correct price and time with the correct report
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));

    //create a stale report
    report[0] = FeeQuoter.ReceivedCCIPFeedReport({
      token: s_onReportTestToken1,
      price: 4e18,
      timestamp: uint32(block.timestamp - 1)
    });

    //record logs to check no events were emitted
    vm.recordLogs();

    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));

    //no logs should have been emitted
    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_RevertWhen_onReportWhen_TokenNotSupported() public {
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_sourceTokens[1], price: 4e18, timestamp: uint32(block.timestamp)});

    // Revert due to token config not being set with the isEnabled flag
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, s_sourceTokens[1]));
    changePrank(address(s_forwarder));
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }

  function test_RevertWhen_onReportWhen_InvalidForwarder() public {
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_sourceTokens[0], price: 4e18, timestamp: uint32(block.timestamp)});

    vm.expectRevert(
      abi.encodeWithSelector(
        KeystoneFeedsPermissionHandler.ReportForwarderUnauthorized.selector,
        STRANGER,
        WORKFLOW_OWNER_1,
        WORKFLOW_NAME_1,
        REPORT_NAME_1
      )
    );
    changePrank(STRANGER);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }
}
