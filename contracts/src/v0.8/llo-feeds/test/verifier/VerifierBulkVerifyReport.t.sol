// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifierAndFeeManager} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierTestWithConfiguredVerifierAndFeeManager is BaseTestWithConfiguredVerifierAndFeeManager {
  uint256 internal constant DEFAULT_LINK_MINT_QUANTITY = 100 ether;
  uint256 internal constant DEFAULT_NATIVE_MINT_QUANTITY = 100 ether;

  function setUp() public virtual override {
    super.setUp();

    //mint some tokens to the user
    link.mint(USER, DEFAULT_LINK_MINT_QUANTITY);
    native.mint(USER, DEFAULT_NATIVE_MINT_QUANTITY);
    vm.deal(USER, DEFAULT_NATIVE_MINT_QUANTITY);
  }
}

contract VerifierBulkVerifyBillingReport is VerifierTestWithConfiguredVerifierAndFeeManager {
  function test_verifyWithBulkLink() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    //create an array containing 5 link reports
    bytes[] memory signedReports = new bytes[](5);
    for (uint256 i = 0; i < 5; i++) {
      signedReports[i] = signedReport;
    }

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 5, USER);

    _verify(signedReports, 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 5);
  }

  function test_verifyWithBulkNative() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    //create an array containing 5 link reports
    bytes[] memory signedReports = new bytes[](5);
    for (uint256 i = 0; i < 5; i++) {
      signedReports[i] = signedReport;
    }

    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 5, USER);

    _verify(signedReports, 0, USER);

    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
  }

  function test_verifyWithBulkNativeUnwrapped() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    //create an array containing 5 link reports
    bytes[] memory signedReports = new bytes[](5);
    for (uint256 i = 0; i < 5; i++) {
      signedReports[i] = signedReport;
    }

    _verify(signedReports, DEFAULT_REPORT_NATIVE_FEE * 5, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyWithNativeUnwrappedReturnsChange() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV2Report(),
      _generateReportContext(FEED_ID_V3),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    //create an array containing 5 link reports
    bytes[] memory signedReports = new bytes[](5);
    for (uint256 i = 0; i < 5; i++) {
      signedReports[i] = signedReport;
    }

    _verify(signedReports, DEFAULT_REPORT_NATIVE_FEE * 10, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(address(feeManager).balance, 0);
  }
}
