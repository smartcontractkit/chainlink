// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifier} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../../../src/v0.8/Verifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";
import {AccessControllerInterface} from "../../../../src/v0.8/interfaces/AccessControllerInterface.sol";

contract VerifierVerifyTest is BaseTestWithConfiguredVerifier {
  bytes32[3] internal s_reportContext;

  event ReportVerified(bytes32 indexed feedId, bytes32 reportHash, address requester);

  Report internal s_testReportOne;

  function setUp() public virtual override {
    BaseTestWithConfiguredVerifier.setUp();
    (, , bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    s_reportContext[0] = configDigest;
    s_reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReportOne = _createReport(
      FEED_ID,
      OBSERVATIONS_TIMESTAMP,
      MEDIAN,
      BID,
      ASK,
      BLOCKNUMBER_UPPER_BOUND,
      blockhash(BLOCKNUMBER_UPPER_BOUND),
      BLOCKNUMBER_LOWER_BOUND
    );
  }

  function assertReportsEqual(bytes memory response, Report memory testReport) public {
    (
      bytes32 feedId,
      uint32 timestamp,
      int192 median,
      int192 bid,
      int192 ask,
      uint64 blockNumUB,
      bytes32 upperBlockhash,
      uint64 blockNumLB
    ) = abi.decode(response, (bytes32, uint32, int192, int192, int192, uint64, bytes32, uint64));
    assertEq(feedId, testReport.feedId);
    assertEq(timestamp, testReport.observationsTimestamp);
    assertEq(median, testReport.median);
    assertEq(bid, testReport.bid);
    assertEq(ask, testReport.ask);
    assertEq(blockNumLB, testReport.blocknumberLowerBound);
    assertEq(blockNumUB, testReport.blocknumberUpperBound);
    assertEq(upperBlockhash, testReport.upperBlockhash);
  }
}

contract VerifierProxyVerifyTest is VerifierVerifyTest {
  function test_revertsIfNoVerifierConfigured() public {
    s_reportContext[0] = bytes32("corrupt-digest");
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.VerifierNotFound.selector, bytes32("corrupt-digest")));
    s_verifierProxy.verify(signedReport);
  }

  function test_proxiesToTheCorrectVerifier() public {
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes memory response = s_verifierProxy.verify(signedReport);
    assertReportsEqual(response, s_testReportOne);
  }
}

contract VerifierProxyAccessControlledVerificationTest is VerifierVerifyTest {
  function setUp() public override {
    VerifierVerifyTest.setUp();
    AccessControllerInterface accessController = AccessControllerInterface(ACCESS_CONTROLLER_ADDRESS);

    s_verifierProxy.setAccessController(accessController);
  }

  function test_revertsIfNoAccess() public {
    vm.mockCall(
      ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(false)
    );
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.AccessForbidden.selector));

    changePrank(USER);
    s_verifierProxy.verify(signedReport);
  }

  function test_proxiesToTheVerifierIfHasAccess() public {
    vm.mockCall(
      ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(true)
    );

    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    changePrank(USER);
    bytes memory response = s_verifierProxy.verify(signedReport);
    assertReportsEqual(response, s_testReportOne);
  }
}

contract VerifierVerifySingleConfigDigestTest is VerifierVerifyTest {
  function test_revertsIfVerifiedByNonProxy() public {
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectRevert(abi.encodeWithSelector(Verifier.AccessForbidden.selector));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfVerifiedWithIncorrectAddresses() public {
    Signer[] memory signers = _getSigners(FAULT_TOLERANCE + 1);
    signers[10].mockPrivateKey = 1234;
    bytes memory signedReport = _generateEncodedBlob(s_testReportOne, s_reportContext, signers);
    changePrank(address(s_verifierProxy));
    vm.expectRevert(abi.encodeWithSelector(Verifier.AccessForbidden.selector));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfMismatchedSignatureLength() public {
    bytes32[] memory rs = new bytes32[](FAULT_TOLERANCE + 1);
    bytes32[] memory ss = new bytes32[](FAULT_TOLERANCE + 3);
    bytes32 rawVs = bytes32("");
    bytes memory signedReport = abi.encode(s_reportContext, abi.encode(s_testReportOne), rs, ss, rawVs);
    changePrank(address(s_verifierProxy));
    vm.expectRevert(abi.encodeWithSelector(Verifier.MismatchedSignatures.selector, rs.length, ss.length));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfConfigDigestNotSet() public {
    bytes32[3] memory reportContext = s_reportContext;
    reportContext[0] = bytes32("wrong-context-digest");
    bytes memory signedReport = _generateEncodedBlob(s_testReportOne, reportContext, _getSigners(FAULT_TOLERANCE + 1));
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestNotSet.selector, FEED_ID, reportContext[0]));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfReportHasUnconfiguredFeedID() public {
    Report memory report = _createReport(
      FEED_ID_2,
      OBSERVATIONS_TIMESTAMP,
      MEDIAN,
      BID,
      ASK,
      BLOCKNUMBER_UPPER_BOUND,
      blockhash(BLOCKNUMBER_UPPER_BOUND),
      BLOCKNUMBER_LOWER_BOUND
    );
    bytes memory signedReport = _generateEncodedBlob(report, s_reportContext, _getSigners(FAULT_TOLERANCE + 1));
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestNotSet.selector, FEED_ID_2, s_reportContext[0]));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfWrongNumberOfSigners() public {
    bytes memory signedReport = _generateEncodedBlob(s_testReportOne, s_reportContext, _getSigners(10));
    vm.expectRevert(abi.encodeWithSelector(Verifier.IncorrectSignatureCount.selector, 10, FAULT_TOLERANCE + 1));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfDuplicateSignersHaveSigned() public {
    Signer[] memory signers = _getSigners(FAULT_TOLERANCE + 1);
    // Duplicate signer at index 1
    signers[0] = signers[1];
    bytes memory signedReport = _generateEncodedBlob(s_testReportOne, s_reportContext, signers);
    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_returnsThePriceAndBlockNumIfReportVerified() public {
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    bytes memory response = s_verifier.verify(signedReport, msg.sender);

    assertReportsEqual(response, s_testReportOne);
  }

  function test_setsTheCorrectEpoch() public {
    s_reportContext[1] = bytes32(uint256(5 << 8));
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);

    (, , uint32 latestEpoch) = s_verifier.latestConfigDigestAndEpoch(FEED_ID);
    assertEq(latestEpoch, 5);
  }

  function test_emitsAnEventIfReportVerified() public {
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectEmit(true, true, true, true, address(s_verifier));
    emit ReportVerified(s_testReportOne.feedId, keccak256(abi.encode(s_testReportOne)), msg.sender);
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }
}

contract VerifierVerifyMultipleConfigDigestTest is VerifierVerifyTest {
  bytes32 internal s_oldConfigDigest;
  bytes32 internal s_newConfigDigest;

  uint8 internal constant FAULT_TOLERANCE_TWO = 5;

  function setUp() public override {
    VerifierVerifyTest.setUp();
    (, , s_oldConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(_getSigners(20)),
      s_offchaintransmitters,
      FAULT_TOLERANCE_TWO,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
    (, , s_newConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
  }

  function test_revertsIfVerifyingWithAnUnsetDigest() public {
    s_verifier.deactivateConfig(FEED_ID, (s_oldConfigDigest));

    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestInactive.selector, FEED_ID, s_reportContext[0]));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_canVerifyOlderReportsWithOlderConfigs() public {
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    bytes memory response = s_verifier.verify(signedReport, msg.sender);
    assertReportsEqual(response, s_testReportOne);
  }

  function test_canVerifyNewerReportsWithNewerConfigs() public {
    s_reportContext[0] = s_newConfigDigest;
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE_TWO + 1)
    );
    changePrank(address(s_verifierProxy));
    bytes memory response = s_verifier.verify(signedReport, msg.sender);
    assertReportsEqual(response, s_testReportOne);
  }

  function test_revertsIfAReportIsVerifiedWithAnExistingButIncorrectDigest() public {
    // Try sending the older digest signed with the new set of signers
    s_reportContext[0] = s_oldConfigDigest;
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE_TWO + 1)
    );
    vm.expectRevert(
      abi.encodeWithSelector(Verifier.IncorrectSignatureCount.selector, FAULT_TOLERANCE_TWO + 1, FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }
}
