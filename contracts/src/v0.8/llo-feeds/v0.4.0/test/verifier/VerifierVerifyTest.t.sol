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

    function test_verifyTooglingActiveFlagsDONConfigs() public {
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

        bytes memory signedReport2 =
            _generateV3EncodedBlob(reportAtSetConfig1Timestmap, s_reportContext, reportSigners2);

        // this report is verified via Config1 (using a historical config)
        bytes memory verifierResponse2 = s_verifierProxy.verify(signedReport2, abi.encode(native));
        assertReportsEqual(verifierResponse2, reportAtSetConfig1Timestmap);

        // same report with same signers but with a higher timestamp gets verified via Config2
        // which means verification fails
        bytes memory signedReport3 =
            _generateV3EncodedBlob(reportAtSetConfig2Timestmap, s_reportContext, reportSigners2);
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
        bytes memory signedReport =
            _generateV3EncodedBlob(s_testReportThree, s_reportContext, _getSigners(FAULT_TOLERANCE + 1));

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

}
