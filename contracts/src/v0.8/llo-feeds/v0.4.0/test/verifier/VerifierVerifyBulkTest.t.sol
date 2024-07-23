// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierVerifyBulkTest is BaseTest {
  bytes32[3] internal s_reportContext;
  V3Report internal s_testReportThree;

  function setUp() public virtual override {
    BaseTest.setUp();
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));

    s_testReportThree = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(block.timestamp),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });
  }

  function test_revertsVerifyBulkIfNoAccess() public {
    vm.mockCall(
      ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    bytes memory signedReport = _generateV3EncodedBlob(
      s_testReportThree,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes[] memory signedReports = new bytes[](2);
    signedReports[0] = signedReport;
    signedReports[1] = signedReport;
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.AccessForbidden.selector));
    changePrank(USER);
    s_verifier.verifyBulk(signedReports, abi.encode(native), msg.sender);
  }

  function test_verifyBulkSingleCaseWithSingleConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersSubset1 = new BaseTest.Signer[](7);
    signersSubset1[0] = signers[0];
    signersSubset1[1] = signers[1];
    signersSubset1[2] = signers[2];
    signersSubset1[3] = signers[3];
    signersSubset1[4] = signers[4];
    signersSubset1[5] = signers[5];
    signersSubset1[6] = signers[6];

    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    // Config1
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    V3Report memory report = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(block.timestamp),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });

    BaseTest.Signer[] memory reportSigners = new BaseTest.Signer[](3);
    reportSigners[0] = signers[0];
    reportSigners[1] = signers[1];
    reportSigners[2] = signers[2];

    bytes[] memory signedReports = new bytes[](10);

    bytes memory signedReport = _generateV3EncodedBlob(report, s_reportContext, reportSigners);

    for (uint256 i = 0; i < signedReports.length; i++) {
      signedReports[i] = signedReport;
    }

    bytes[] memory verifierResponses = s_verifierProxy.verifyBulk(signedReports, abi.encode(native));

    for (uint256 i = 0; i < verifierResponses.length; i++) {
      bytes memory verifierResponse = verifierResponses[i];
      assertReportsEqual(verifierResponse, report);
    }
  }

  function test_verifyBulkWithSingleConfigOneVerifyFails() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersSubset1 = new BaseTest.Signer[](7);
    signersSubset1[0] = signers[0];
    signersSubset1[1] = signers[1];
    signersSubset1[2] = signers[2];
    signersSubset1[3] = signers[3];
    signersSubset1[4] = signers[4];
    signersSubset1[5] = signers[5];
    signersSubset1[6] = signers[6];

    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    // Config1
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    BaseTest.Signer[] memory reportSigners = new BaseTest.Signer[](3);
    reportSigners[0] = signers[0];
    reportSigners[1] = signers[1];
    reportSigners[2] = signers[2];

    bytes[] memory signedReports = new bytes[](11);
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, reportSigners);

    for (uint256 i = 0; i < 10; i++) {
      signedReports[i] = signedReport;
    }

    // Making the last report in  this batch not verifiable
    BaseTest.Signer[] memory reportSigners2 = new BaseTest.Signer[](3);
    reportSigners2[0] = signers[30];
    reportSigners2[1] = signers[29];
    reportSigners2[2] = signers[28];
    signedReports[10] = _generateV3EncodedBlob(s_testReportThree, s_reportContext, reportSigners2);

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verifyBulk(signedReports, abi.encode(native));
  }
}
