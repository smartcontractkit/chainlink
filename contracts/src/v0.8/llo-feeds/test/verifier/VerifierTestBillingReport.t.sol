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

    //mint some link tokens to the feeManager pool
    link.mint(address(feeManager), DEFAULT_REPORT_LINK_FEE);
  }
}

contract VerifierTestBillingReport is VerifierTestWithConfiguredVerifierAndFeeManager {
  function test_verifyWithLink() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);

    _verify(signedReport, 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
  }

  function test_verifyWithNative() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    _verify(signedReport, 0, USER);

    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE);
  }

  function test_verifyWithNativeUnwrapped() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    _verify(signedReport, DEFAULT_REPORT_NATIVE_FEE, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyWithNativeUnwrappedReturnsChange() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    _verify(signedReport, DEFAULT_REPORT_NATIVE_FEE * 2, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(address(feeManager).balance, 0);
  }
}

contract VerifierBulkVerifyBillingReport is VerifierTestWithConfiguredVerifierAndFeeManager {
  uint256 internal constant NUMBERS_OF_REPORTS = 5;

  function test_verifyWithBulkLink() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i = 0; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * NUMBERS_OF_REPORTS, USER);

    _verifyBulk(signedReports, 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * NUMBERS_OF_REPORTS);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * NUMBERS_OF_REPORTS);
  }

  function test_verifyWithBulkNative() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i = 0; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS, USER);

    _verifyBulk(signedReports, 0, USER);

    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS);
  }

  function test_verifyWithBulkNativeUnwrapped() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _verifyBulk(signedReports, DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 5);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyWithBulkNativeUnwrappedReturnsChange() public {
    bytes memory signedReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](NUMBERS_OF_REPORTS);
    for (uint256 i = 0; i < NUMBERS_OF_REPORTS; i++) {
      signedReports[i] = signedReport;
    }

    _verifyBulk(signedReports, DEFAULT_REPORT_NATIVE_FEE * (NUMBERS_OF_REPORTS * 2), USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * NUMBERS_OF_REPORTS);
    assertEq(address(feeManager).balance, 0);
  }

  function test_verifyBulkWithLinkAndWrappedNative() public {
    bytes memory linkReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    bytes memory nativeReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](5);

    signedReports[0] = linkReport;
    signedReports[1] = linkReport;
    signedReports[2] = linkReport;
    signedReports[3] = nativeReport;
    signedReports[4] = nativeReport;

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 3, USER);
    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE * 2, USER);

    _verifyBulk(signedReports, 0, USER);

    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 3);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 3);
  }

  function test_verifyBulkWithLinkAndUnwrappedNative() public {
    bytes memory linkReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    bytes memory nativeReport = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](5);

    signedReports[0] = linkReport;
    signedReports[1] = linkReport;
    signedReports[2] = linkReport;
    signedReports[3] = nativeReport;
    signedReports[4] = nativeReport;

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE * 3, USER);

    _verifyBulk(signedReports, DEFAULT_REPORT_NATIVE_FEE * 40, USER);

    //user should have some link returned as twice the amount was passed into msg.value
    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(native.balanceOf(address(feeManager)), DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE * 3);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 3);
  }

  function test_verifyBulkReportWithUnwrappedAndWrappedNativeDefaultsToUnwrapped() public {
    bytes memory nativeReport1 = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes memory nativeReport2 = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](2);

    signedReports[0] = nativeReport1;
    signedReports[1] = nativeReport2;

    _verifyBulk(signedReports, DEFAULT_REPORT_NATIVE_FEE * 2, USER);

    assertEq(USER.balance, DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE * 2);
    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY);
  }

  function test_verifyMultiVersions() public {
    bytes memory signedReportV1 = _generateEncodedBlob(
      _generateV1Report(),
      _generateReportContext(v1ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes memory signedReportV3Link = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    bytes memory signedReportV3Native = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](3);

    signedReports[0] = signedReportV1;
    signedReports[1] = signedReportV3Link;
    signedReports[2] = signedReportV3Native;

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    _verifyBulk(signedReports, 0, USER);

    assertEq(link.balanceOf(USER), DEFAULT_LINK_MINT_QUANTITY - DEFAULT_REPORT_LINK_FEE);
    assertEq(native.balanceOf(USER), DEFAULT_NATIVE_MINT_QUANTITY - DEFAULT_REPORT_NATIVE_FEE);
    assertEq(link.balanceOf(address(rewardManager)), DEFAULT_REPORT_LINK_FEE * 2);
  }

  function test_verifyMultiVersionsReturnsVerifiedReports() public {
    bytes memory signedReportV1 = _generateEncodedBlob(
      _generateV1Report(),
      _generateReportContext(v1ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes memory signedReportV3Link = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(link))
    );

    bytes memory signedReportV3Native = _generateEncodedBlobWithQuote(
      _generateV3Report(),
      _generateReportContext(v3ConfigDigest),
      _getSigners(FAULT_TOLERANCE + 1),
      _generateQuote(address(native))
    );

    bytes[] memory signedReports = new bytes[](3);

    signedReports[0] = signedReportV1;
    signedReports[1] = signedReportV3Link;
    signedReports[2] = signedReportV3Native;

    _approveLink(address(rewardManager), DEFAULT_REPORT_LINK_FEE, USER);
    _approveNative(address(feeManager), DEFAULT_REPORT_NATIVE_FEE, USER);

    address originalAddr = msg.sender;
    changePrank(USER);

    bytes[] memory verifierReports = s_verifierProxy.verifyBulk{value: 0}(signedReports);

    changePrank(originalAddr);

    assertEq(verifierReports[0], _encodeReport(_generateV1Report()));
    assertEq(verifierReports[1], _encodeReport(_generateV3Report()));
    assertEq(verifierReports[2], _encodeReport(_generateV3Report()));
  }
}
