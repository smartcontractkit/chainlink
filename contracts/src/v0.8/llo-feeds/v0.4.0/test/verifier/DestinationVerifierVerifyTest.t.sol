// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationVerifierProxy} from "../../../v0.4.0/DestinationVerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierVerifyTest is BaseTest {
  bytes32[3] internal s_reportContext;
  V3Report internal s_testReportThree;

  function setUp() public virtual override {
    BaseTest.setUp();

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

  function test_verifyReport() public {
    // Simple use case just setting a config and verifying a report
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);

    bytes memory verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, s_testReportThree);
  }

  function test_verifyTooglingActiveFlagsDonConfigs() public {
    // sets config
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signers);
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    // verifies report
    bytes memory verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, s_testReportThree);

    // test verifying via a config that is deactivated
    s_verifier.setConfigActive(0, false);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));

    // test verifying via a reactivated config
    s_verifier.setConfigActive(0, true);
    verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, s_testReportThree);
  }

  function test_failToVerifyReportIfNotEnoughSigners() public {
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
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    // only one signer, signers < MINIMAL_FAULT_TOLERANCE
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](1);
    signersSubset2[0] = signers[4];

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signersSubset2);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_failToVerifyReportIfNoSigners() public {
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
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    // No signers for this report
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](0);
    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signersSubset2);

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.NoSigners.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_failToVerifyReportIfDupSigners() public {
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
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    // One signer is repeated
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](4);
    signersSubset2[0] = signers[0];
    signersSubset2[1] = signers[1];
    // repeated signers
    signersSubset2[2] = signers[2];
    signersSubset2[3] = signers[2];

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, signersSubset2);

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_failToVerifyReportIfSignerNotInConfig() public {
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
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    // one report whose signer is not in the config
    BaseTest.Signer[] memory reportSigners = new BaseTest.Signer[](4);
    // these signers are part ofm the config
    reportSigners[0] = signers[4];
    reportSigners[1] = signers[5];
    reportSigners[2] = signers[6];
    // this single signer is not in the config
    reportSigners[3] = signers[7];

    bytes memory signedReport = _generateV3EncodedBlob(s_testReportThree, s_reportContext, reportSigners);

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_canVerifyOlderV3ReportsWithOlderConfigs() public {
    /*
          This test is checking we can use historical Configs to verify reports:
          - DonConfigA has signers {A, B, C, E} is set at time T1
          - DonConfigB has signers {A, B, C, D} is set at time T2
          - checks we can verify a report with {B, C, D} signers (via DonConfigB)
          - checks we can verify a report with {B, C, E} signers and timestamp below T2 (via DonConfigA historical config)
          - checks we can't verify a report with {B, C, E} signers and timestamp above T2 (it gets verivied via DonConfigB)
          - sets DonConfigA as deactivated
          - checks we can't verify a report with {B, C, E} signers and timestamp below T2 (via DonConfigA)
         */
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

    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](7);
    signersSubset2[0] = signers[0];
    signersSubset2[1] = signers[1];
    signersSubset2[2] = signers[2];
    signersSubset2[3] = signers[3];
    signersSubset2[4] = signers[4];
    signersSubset2[5] = signers[5];
    signersSubset2[6] = signers[29];
    address[] memory signersAddrSubset2 = _getSignerAddresses(signersSubset2);

    V3Report memory reportAtSetConfig1Timestmap = V3Report({
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

    vm.warp(block.timestamp + 100);

    // Config2
    s_verifier.setConfig(signersAddrSubset2, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    V3Report memory reportAtSetConfig2Timestmap = V3Report({
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

    BaseTest.Signer[] memory reportSigners = new BaseTest.Signer[](5);
    reportSigners[0] = signers[0];
    reportSigners[1] = signers[1];
    reportSigners[2] = signers[2];
    reportSigners[3] = signers[3];
    reportSigners[4] = signers[29];

    bytes memory signedReport = _generateV3EncodedBlob(reportAtSetConfig2Timestmap, s_reportContext, reportSigners);

    // this report is verified via Config2
    bytes memory verifierResponse = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(verifierResponse, reportAtSetConfig2Timestmap);

    BaseTest.Signer[] memory reportSigners2 = new BaseTest.Signer[](5);
    reportSigners2[0] = signers[0];
    reportSigners2[1] = signers[1];
    reportSigners2[2] = signers[2];
    reportSigners2[3] = signers[3];
    reportSigners2[4] = signers[6];

    bytes memory signedReport2 = _generateV3EncodedBlob(reportAtSetConfig1Timestmap, s_reportContext, reportSigners2);

    // this report is verified via Config1 (using a historical config)
    bytes memory verifierResponse2 = s_verifierProxy.verify(signedReport2, abi.encode(native));
    assertReportsEqual(verifierResponse2, reportAtSetConfig1Timestmap);

    // same report with same signers but with a higher timestamp gets verified via Config2
    // which means verification fails
    bytes memory signedReport3 = _generateV3EncodedBlob(reportAtSetConfig2Timestmap, s_reportContext, reportSigners2);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport3, abi.encode(native));

    // deactivating Config1 and trying a reverifications ends in failure
    s_verifier.setConfigActive(0, false);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport2, abi.encode(native));
  }

  function test_revertsVerifyIfNoAccess() public {
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

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.AccessForbidden.selector));

    changePrank(USER);
    s_verifier.verify(signedReport, abi.encode(native), msg.sender);
  }

  function test_canVerifyNewerReportsWithNewerConfigs() public {
    /*
          This test is checking that we use prefer verifiying via newer configs instead of old ones.
          - DonConfigA has signers {A, B, C, E} is set at time T1
          - DonConfigB has signers {F, G, H, I} is set at time T2
          - DonConfigC has signers {J, K, L, M } is set at time T3
          - checks we can verify a report with {K, L, M} signers (via DonConfigC)
         */
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
    vm.warp(block.timestamp + 1);

    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](7);
    signersSubset2[0] = signers[7];
    signersSubset2[1] = signers[8];
    signersSubset2[2] = signers[9];
    signersSubset2[3] = signers[10];
    signersSubset2[4] = signers[11];
    signersSubset2[5] = signers[12];
    signersSubset2[6] = signers[13];

    address[] memory signersAddrSubset2 = _getSignerAddresses(signersSubset2);
    // Config2
    s_verifier.setConfig(signersAddrSubset2, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    vm.warp(block.timestamp + 1);

    BaseTest.Signer[] memory signersSubset3 = new BaseTest.Signer[](7);
    signersSubset3[0] = signers[30];
    signersSubset3[1] = signers[29];
    signersSubset3[2] = signers[28];
    signersSubset3[3] = signers[27];
    signersSubset3[4] = signers[26];
    signersSubset3[5] = signers[25];
    signersSubset3[6] = signers[24];

    address[] memory signersAddrSubset3 = _getSignerAddresses(signersSubset3);
    // Config3
    s_verifier.setConfig(signersAddrSubset3, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    vm.warp(block.timestamp + 1);

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
    reportSigners[0] = signers[30];
    reportSigners[1] = signers[29];
    reportSigners[2] = signers[28];

    bytes memory signedReport = _generateV3EncodedBlob(report, s_reportContext, reportSigners);

    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_rollingOutConfiguration() public {
    /*
          This test is checking that we can roll out to a new DON without downtime using a transition configuration
          - DonConfigA has signers {A, B, C} is set at time T1
          - DonConfigB (transition config) has signers {A, B, C, D, E, F} is set at time T2
          - DonConfigC has signers {D, E, F} is set at time T3
          
          - checks we can verify a report with {A, B, C} signers (via DonConfigA) at time between T1 and T2
          - checks we can verify a report with {A, B, C} signers (via DonConfigB) at time between T2 and T3
          - checks we can verify a report with {D, E, F} signers (via DonConfigB) at time between T2 and T3
          - checks we can verify a report with {D, E, F} signers (via DonConfigC) at time > T3
          - checks we can't verify a report with {A, B, C} signers (via DonConfigC) and timestamp >T3 at time > T3
          - checks we can verify a report with {A, B, C} signers (via DonConfigC) and timestamp between T2 and T3  at time > T3 (historical check)

         */

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

    // ConfigA
    address[] memory signersAddrSubset1 = _getSignerAddresses(signersSubset1);
    s_verifier.setConfig(signersAddrSubset1, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    V3Report memory reportT1 = V3Report({
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

    BaseTest.Signer[] memory reportSignersConfigA = new BaseTest.Signer[](3);
    reportSignersConfigA[0] = signers[0];
    reportSignersConfigA[1] = signers[1];
    reportSignersConfigA[2] = signers[2];

    // just testing ConfigA
    bytes memory signedReport = _generateV3EncodedBlob(reportT1, s_reportContext, reportSignersConfigA);
    s_verifierProxy.verify(signedReport, abi.encode(native));

    vm.warp(block.timestamp + 100);

    BaseTest.Signer[] memory signersSuperset = new BaseTest.Signer[](14);
    // signers in ConfigA
    signersSuperset[0] = signers[0];
    signersSuperset[1] = signers[1];
    signersSuperset[2] = signers[2];
    signersSuperset[3] = signers[3];
    signersSuperset[4] = signers[4];
    signersSuperset[5] = signers[5];
    signersSuperset[6] = signers[6];
    // new signers
    signersSuperset[7] = signers[7];
    signersSuperset[8] = signers[8];
    signersSuperset[9] = signers[9];
    signersSuperset[10] = signers[10];
    signersSuperset[11] = signers[11];
    signersSuperset[12] = signers[12];
    signersSuperset[13] = signers[13];

    BaseTest.Signer[] memory reportSignersConfigC = new BaseTest.Signer[](3);
    reportSignersConfigC[0] = signers[7];
    reportSignersConfigC[1] = signers[8];
    reportSignersConfigC[2] = signers[9];

    // ConfigB (transition Config)
    address[] memory signersAddrsSuperset = _getSignerAddresses(signersSuperset);
    s_verifier.setConfig(signersAddrsSuperset, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    V3Report memory reportT2 = V3Report({
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

    // testing we can verify a fresh (block timestamp) report with ConfigA signers. This should use ConfigB
    signedReport = _generateV3EncodedBlob(reportT2, s_reportContext, reportSignersConfigA);
    s_verifierProxy.verify(signedReport, abi.encode(native));

    // testing we can verify an old ( non fresh block timestamp) report with ConfigA signers. This should use ConfigA
    signedReport = _generateV3EncodedBlob(reportT1, s_reportContext, reportSignersConfigA);
    s_verifierProxy.verify(signedReport, abi.encode(native));
    // deactivating to make sure we are really verifiying via ConfigA
    s_verifier.setConfigActive(0, false);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
    s_verifier.setConfigActive(0, true);

    // testing we can verify a fresh  (block timestamp) report with the new signers.  This should use ConfigB
    signedReport = _generateV3EncodedBlob(reportT2, s_reportContext, reportSignersConfigC);
    s_verifierProxy.verify(signedReport, abi.encode(native));

    vm.warp(block.timestamp + 100);

    // Adding ConfigC
    BaseTest.Signer[] memory signersSubset2 = new BaseTest.Signer[](7);
    signersSubset2[0] = signers[7];
    signersSubset2[1] = signers[8];
    signersSubset2[2] = signers[9];
    signersSubset2[3] = signers[10];
    signersSubset2[4] = signers[11];
    signersSubset2[5] = signers[12];
    signersSubset2[6] = signers[13];
    address[] memory signersAddrsSubset2 = _getSignerAddresses(signersSubset2);
    s_verifier.setConfig(signersAddrsSubset2, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    V3Report memory reportT3 = V3Report({
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

    // testing we can verify reports with ConfigC signers
    signedReport = _generateV3EncodedBlob(reportT3, s_reportContext, reportSignersConfigC);
    s_verifierProxy.verify(signedReport, abi.encode(native));

    //  testing an old report (block timestamp) with ConfigC signers should  verify via ConfigB
    signedReport = _generateV3EncodedBlob(reportT2, s_reportContext, reportSignersConfigC);
    s_verifierProxy.verify(signedReport, abi.encode(native));
    // deactivating to make sure we are really verifiying via ConfigB
    s_verifier.setConfigActive(1, false);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
    s_verifier.setConfigActive(1, true);

    // testing a recent report with ConfigA signers should not verify
    signedReport = _generateV3EncodedBlob(reportT3, s_reportContext, reportSignersConfigA);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));

    // testing an old report (block timestamp) with ConfigA signers should  verify via ConfigB
    signedReport = _generateV3EncodedBlob(reportT2, s_reportContext, reportSignersConfigA);
    s_verifierProxy.verify(signedReport, abi.encode(native));
    // deactivating to make sure we are really verifiying via ConfigB
    s_verifier.setConfigActive(1, false);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
    s_verifier.setConfigActive(1, true);

    // testing an old report (block timestamp) with ConfigA signers should  verify via ConfigA
    signedReport = _generateV3EncodedBlob(reportT1, s_reportContext, reportSignersConfigA);
    s_verifierProxy.verify(signedReport, abi.encode(native));
    // deactivating to make sure we are really verifiying via ConfigB
    s_verifier.setConfigActive(0, false);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
    s_verifier.setConfigActive(0, true);
  }

  function test_verifyFailsWhenReportIsOlderThanConfig() public {
    /*
          - SetConfig A at time T0
          - SetConfig B at time T1
          - tries verifing report issued at blocktimestmap < T0
          
          this test is failing: ToDo Ask Michael
         */
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    s_reportContext[0] = bytes32(abi.encode(uint32(5), uint8(1)));

    vm.warp(block.timestamp + 100);

    V3Report memory reportAtTMinus100 = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(block.timestamp - 100),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    vm.warp(block.timestamp + 100);
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0));

    bytes memory signedReport = _generateV3EncodedBlob(reportAtTMinus100, s_reportContext, signers);

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_scenarioRollingNewChainWithHistoricConfigs() public {
    /*
       This test is checking that we can roll out in a new network and set historic configurations :
       - Stars with a chain at blocktimestamp 1000
       - SetConfigA with teimstamp 100 
       - SetConfigB with timesmtap 200
       - SetConfigC with timestamp current 
       - tries verifying reports for all the configs
    */

    vm.warp(block.timestamp + 1000);

    Signer[] memory signers = _getSigners(MAX_ORACLES);

    uint8 MINIMAL_FAULT_TOLERANCE = 2;
    BaseTest.Signer[] memory signersA = new BaseTest.Signer[](7);
    signersA[0] = signers[0];
    signersA[1] = signers[1];
    signersA[2] = signers[2];
    signersA[3] = signers[3];
    signersA[4] = signers[4];
    signersA[5] = signers[5];
    signersA[6] = signers[6];

    // ConfigA (historical config)
    uint32 configATimestmap = 100;
    address[] memory signersAddrA = _getSignerAddresses(signersA);
    s_verifier.setConfigWithActivationTime(
      signersAddrA,
      MINIMAL_FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0),
      configATimestmap
    );

    // ConfigB (historical config)
    uint32 configBTimestmap = 200;
    // Config B
    BaseTest.Signer[] memory signersB = new BaseTest.Signer[](7);
    // signers in ConfigA
    signersB[0] = signers[8];
    signersB[1] = signers[9];
    signersB[2] = signers[10];
    signersB[3] = signers[11];
    signersB[4] = signers[12];
    signersB[5] = signers[13];
    signersB[6] = signers[14];
    address[] memory signersAddrsB = _getSignerAddresses(signersB);
    s_verifier.setConfigWithActivationTime(
      signersAddrsB,
      MINIMAL_FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0),
      configBTimestmap
    );

    // ConfigC (config at current timestamp)
    //    BaseTest.Signer[] memory signersC = new BaseTest.Signer[](7);
    // signers in ConfigA
    signersB[6] = signers[15];
    address[] memory signersAddrsC = _getSignerAddresses(signersB);
    s_verifier.setConfig(signersAddrsC, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    vm.warp(block.timestamp + 10);

    // historical report
    V3Report memory s_testReportA = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(101),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp + 1000),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });

    // historical report
    V3Report memory s_testReportB = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(201),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp + 1000),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });

    // report at recent timestamp
    V3Report memory s_testReportC = V3Report({
      feedId: FEED_ID_V3,
      observationsTimestamp: OBSERVATIONS_TIMESTAMP,
      validFromTimestamp: uint32(block.timestamp),
      nativeFee: uint192(DEFAULT_REPORT_NATIVE_FEE),
      linkFee: uint192(DEFAULT_REPORT_LINK_FEE),
      expiresAt: uint32(block.timestamp + 1000),
      benchmarkPrice: MEDIAN,
      bid: BID,
      ask: ASK
    });

    BaseTest.Signer[] memory reportSignersA = new BaseTest.Signer[](3);
    reportSignersA[0] = signers[0];
    reportSignersA[1] = signers[1];
    reportSignersA[2] = signers[2];

    BaseTest.Signer[] memory reportSignersB = new BaseTest.Signer[](3);
    reportSignersB[0] = signers[8];
    reportSignersB[1] = signers[9];
    reportSignersB[2] = signers[14];

    BaseTest.Signer[] memory reportSignersC = new BaseTest.Signer[](3);
    reportSignersC[0] = signers[15];
    reportSignersC[1] = signers[13];
    reportSignersC[2] = signers[12];

    bytes memory signedReportA = _generateV3EncodedBlob(s_testReportA, s_reportContext, reportSignersA);
    bytes memory signedReportB = _generateV3EncodedBlob(s_testReportB, s_reportContext, reportSignersB);
    bytes memory signedReportC = _generateV3EncodedBlob(s_testReportC, s_reportContext, reportSignersC);

    // verifying historical reports
    s_verifierProxy.verify(signedReportA, abi.encode(native));
    s_verifierProxy.verify(signedReportB, abi.encode(native));
    // verifiying a current report
    s_verifierProxy.verify(signedReportC, abi.encode(native));

    // current report verified by historical report fails
    bytes memory signedNewReportWithOldSignatures = _generateV3EncodedBlob(
      s_testReportC,
      s_reportContext,
      reportSignersA
    );
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadVerification.selector));
    s_verifierProxy.verify(signedNewReportWithOldSignatures, abi.encode(native));
  }
}
