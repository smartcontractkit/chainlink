// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest, BaseTestWithConfiguredVerifier} from "../BaseVerifierTest.t.sol";
import {Verifier} from "../../../../../src/v0.8/Verifier.sol";
import {SimpleWriteAccessController} from "../../../../../src/v0.8/SimpleWriteAccessController.sol";

contract Gas_SetConfigTest is BaseTest {
    address[] internal s_signerAddrs;

    function setUp() public override {
        BaseTest.setUp();
        vm.prank(ADMIN);
        Signer[] memory signers = _getSigners(MAX_ORACLES);
        s_signerAddrs = _getSignerAddresses(signers);
    }

    function testSetConfigSuccess() public {
        vm.prank(ADMIN);
        s_verifier.setConfig(
            FEED_ID,
            s_signerAddrs,
            s_offchaintransmitters,
            FAULT_TOLERANCE,
            bytes(""),
            VERIFIER_VERSION,
            bytes("")
        );
    }
}

contract Gas_VerifyTest is BaseTestWithConfiguredVerifier {
    bytes internal s_signedReport;
    bytes32 internal s_configDigest;

    function setUp() public override {
        BaseTestWithConfiguredVerifier.setUp();
        BaseTest.Report memory s_testReportOne = _createReport(
            FEED_ID,
            OBSERVATIONS_TIMESTAMP,
            MEDIAN,
            BID,
            ASK,
            BLOCKNUMBER_UPPER_BOUND,
            blockhash(BLOCKNUMBER_UPPER_BOUND),
            BLOCKNUMBER_LOWER_BOUND
        );
        (, , s_configDigest) = s_verifier.latestConfigDetails(FEED_ID);
        bytes32[3] memory reportContext;
        reportContext[0] = s_configDigest;
        reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
        s_signedReport = _generateEncodedBlob(
            s_testReportOne,
            reportContext,
            _getSigners(FAULT_TOLERANCE + 1)
        );
        vm.prank(ADMIN);
    }

    function testVerifySuccess() public {
        vm.prank(address(s_verifierProxy));
        s_verifier.verify(s_signedReport, msg.sender);
    }

    function testVerifyProxySuccess() public {
        s_verifierProxy.verify(s_signedReport);
    }
}

contract GasAccessControlledVerifyTest is BaseTestWithConfiguredVerifier {
    bytes internal s_signedReport;
    bytes32 internal s_configDigest;
    SimpleWriteAccessController s_accessController;

    address internal constant CLIENT = address(9000);
    address internal constant ACCESS_CONTROLLER_ADDR = address(10000);

    function setUp() public override {
        BaseTestWithConfiguredVerifier.setUp();
        BaseTest.Report memory s_testReportOne = _createReport(
            FEED_ID,
            OBSERVATIONS_TIMESTAMP,
            MEDIAN,
            BID,
            ASK,
            BLOCKNUMBER_UPPER_BOUND,
            blockhash(BLOCKNUMBER_UPPER_BOUND),
            BLOCKNUMBER_LOWER_BOUND
        );
        (, , s_configDigest) = s_verifier.latestConfigDetails(FEED_ID);
        bytes32[3] memory reportContext;
        reportContext[0] = s_configDigest;
        reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
        s_signedReport = _generateEncodedBlob(
            s_testReportOne,
            reportContext,
            _getSigners(FAULT_TOLERANCE + 1)
        );
        vm.prank(ADMIN);
        s_accessController = new SimpleWriteAccessController();
        vm.prank(ADMIN);
        s_verifierProxy.setAccessController(s_accessController);
        vm.prank(ADMIN);
        s_accessController.addAccess(CLIENT);
    }

    function testVerifyWithAccessControl() public {
        vm.prank(CLIENT);
        s_verifierProxy.verify(s_signedReport);
    }
}
