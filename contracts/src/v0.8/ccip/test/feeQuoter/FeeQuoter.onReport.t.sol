// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {KeystoneFeedsPermissionHandler} from "../../../keystone/KeystoneFeedsPermissionHandler.sol";
import {FeeQuoter} from "../../FeeQuoter.sol";
import {FeeQuoterSetup} from "./FeeQuoterSetup.t.sol";

contract FeeQuoter_onReport is FeeQuoterSetup {
  address internal constant FORWARDER_1 = address(0x1);
  address internal constant WORKFLOW_OWNER_1 = address(0x3);
  bytes10 internal constant WORKFLOW_NAME_1 = "workflow1";
  bytes2 internal constant REPORT_NAME_1 = "01";
  address internal s_onReportTestToken1;
  address internal s_onReportTestToken2;

  function setUp() public virtual override {
    super.setUp();
    s_onReportTestToken1 = s_sourceTokens[0];
    s_onReportTestToken2 = _deploySourceToken("onReportTestToken2", 0, 20);

    KeystoneFeedsPermissionHandler.Permission[] memory permissions = new KeystoneFeedsPermissionHandler.Permission[](1);
    permissions[0] = KeystoneFeedsPermissionHandler.Permission({
      forwarder: FORWARDER_1,
      workflowOwner: WORKFLOW_OWNER_1,
      workflowName: WORKFLOW_NAME_1,
      reportName: REPORT_NAME_1,
      isAllowed: true
    });
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

  function test_onReport_Success() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);

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

    changePrank(FORWARDER_1);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));

    vm.assertEq(s_feeQuoter.getTokenPrice(report[0].token).value, expectedStoredToken1Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(report[0].token).timestamp, report[0].timestamp);

    vm.assertEq(s_feeQuoter.getTokenPrice(report[1].token).value, expectedStoredToken2Price);
    vm.assertEq(s_feeQuoter.getTokenPrice(report[1].token).timestamp, report[1].timestamp);
  }

  function test_OnReport_StaleUpdate_SkipPriceUpdate_Success() public {
    //Creating a correct report
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);

    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_onReportTestToken1, price: 4e18, timestamp: uint32(block.timestamp)});

    uint224 expectedStoredTokenPrice = s_feeQuoter.calculateRebasedValue(18, 18, report[0].price);

    vm.expectEmit();
    emit FeeQuoter.UsdPerTokenUpdated(s_onReportTestToken1, expectedStoredTokenPrice, block.timestamp);

    changePrank(FORWARDER_1);
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

  function test_onReport_TokenNotSupported_Revert() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_sourceTokens[1], price: 4e18, timestamp: uint32(block.timestamp)});

    // Revert due to token config not being set with the isEnabled flag
    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, s_sourceTokens[1]));
    vm.startPrank(FORWARDER_1);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }

  function test_onReport_InvalidForwarder_Reverts() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);
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

  function test_onReport_UnsupportedToken_Reverts() public {
    bytes memory encodedPermissionsMetadata =
      abi.encodePacked(keccak256(abi.encode("workflowCID")), WORKFLOW_NAME_1, WORKFLOW_OWNER_1, REPORT_NAME_1);
    FeeQuoter.ReceivedCCIPFeedReport[] memory report = new FeeQuoter.ReceivedCCIPFeedReport[](1);
    report[0] =
      FeeQuoter.ReceivedCCIPFeedReport({token: s_sourceTokens[1], price: 4e18, timestamp: uint32(block.timestamp)});

    vm.expectRevert(abi.encodeWithSelector(FeeQuoter.TokenNotSupported.selector, s_sourceTokens[1]));
    changePrank(FORWARDER_1);
    s_feeQuoter.onReport(encodedPermissionsMetadata, abi.encode(report));
  }
}
