// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTestWithConfiguredVerifierAndFeeManager} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierVerifyTest is BaseTestWithConfiguredVerifierAndFeeManager {
  bytes32[3] internal s_reportContext;

  event ReportVerified(bytes32 indexed feedId, address requester);

  V1Report internal s_testReportOne;

  function setUp() public virtual override {
    BaseTestWithConfiguredVerifierAndFeeManager.setUp();
    s_reportContext[0] = v1ConfigDigest;
    s_reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReportOne = _createV1Report(
      FEED_ID,
      OBSERVATIONS_TIMESTAMP,
      MEDIAN,
      BID,
      ASK,
      BLOCKNUMBER_UPPER_BOUND,
      blockhash(BLOCKNUMBER_UPPER_BOUND),
      BLOCKNUMBER_LOWER_BOUND,
      uint32(block.timestamp)
    );
  }

  function assertReportsEqual(bytes memory response, V1Report memory testReport) public pure {
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

contract VerifierProxyVerifyTestV05 is VerifierVerifyTest {
  function test_revertsIfNoVerifierConfigured() public {
    s_reportContext[0] = bytes32("corrupt-digest");
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.VerifierNotFound.selector, bytes32("corrupt-digest")));
    s_verifierProxy.verify(signedReport, bytes(""));
  }

  function test_proxiesToTheCorrectVerifier() public {
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    bytes memory response = s_verifierProxy.verify(signedReport, abi.encode(native));
    assertReportsEqual(response, s_testReportOne);
  }
}

contract VerifierProxyAccessControlledVerificationTestV05 is VerifierVerifyTest {
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
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.AccessForbidden.selector));

    changePrank(USER);
    s_verifierProxy.verify(signedReport, abi.encode(native));
  }

  function test_proxiesToTheVerifierIfHasAccess() public {
    vm.mockCall(
      ACCESS_CONTROLLER_ADDRESS,
      abi.encodeWithSelector(AccessControllerInterface.hasAccess.selector, USER),
      abi.encode(true)
    );

    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    changePrank(USER);
    bytes memory response = s_verifierProxy.verify(signedReport, bytes(""));
    assertReportsEqual(response, s_testReportOne);
  }
}

contract VerifierVerifySingleConfigDigestTestV05 is VerifierVerifyTest {
  function test_revertsIfVerifiedByNonProxy() public {
    bytes memory signedReport = _generateV1EncodedBlob(
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
    bytes memory signedReport = _generateV1EncodedBlob(s_testReportOne, s_reportContext, signers);
    changePrank(address(s_verifierProxy));
    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
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
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestInactive.selector, reportContext[0]));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfReportHasUnconfiguredConfigDigest() public {
    V1Report memory report = _createV1Report(
      FEED_ID_2,
      OBSERVATIONS_TIMESTAMP,
      MEDIAN,
      BID,
      ASK,
      BLOCKNUMBER_UPPER_BOUND,
      blockhash(BLOCKNUMBER_UPPER_BOUND),
      BLOCKNUMBER_LOWER_BOUND,
      uint32(block.timestamp)
    );
    s_reportContext[0] = keccak256("unconfigured-digesty");
    bytes memory signedReport = _generateV1EncodedBlob(report, s_reportContext, _getSigners(FAULT_TOLERANCE + 1));
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestInactive.selector, s_reportContext[0]));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfWrongNumberOfSigners() public {
    bytes memory signedReport = _generateV1EncodedBlob(s_testReportOne, s_reportContext, _getSigners(10));
    vm.expectRevert(abi.encodeWithSelector(Verifier.IncorrectSignatureCount.selector, 10, FAULT_TOLERANCE + 1));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_revertsIfDuplicateSignersHaveSigned() public {
    Signer[] memory signers = _getSigners(FAULT_TOLERANCE + 1);
    // Duplicate signer at index 1
    signers[0] = signers[1];
    bytes memory signedReport = _generateV1EncodedBlob(s_testReportOne, s_reportContext, signers);
    vm.expectRevert(abi.encodeWithSelector(Verifier.BadVerification.selector));
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_returnsThePriceAndBlockNumIfReportVerified() public {
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    bytes memory response = s_verifier.verify(signedReport, msg.sender);

    assertReportsEqual(response, s_testReportOne);
  }

  function test_emitsAnEventIfReportVerified() public {
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    vm.expectEmit(true, true, true, true, address(s_verifier));
    emit ReportVerified(s_testReportOne.feedId, msg.sender);
    changePrank(address(s_verifierProxy));
    s_verifier.verify(signedReport, msg.sender);
  }
}

contract VerifierVerifyMultipleConfigDigestTestV05 is VerifierVerifyTest {
  bytes32 internal s_oldConfigDigest;
  bytes32 internal s_newConfigDigest;

  uint8 internal constant FAULT_TOLERANCE_TWO = 5;

  function setUp() public override {
    VerifierVerifyTest.setUp();
    s_oldConfigDigest = v1ConfigDigest;

    s_newConfigDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      2,
      _getSignerAddresses(_getSigners(20)),
      s_offchaintransmitters,
      FAULT_TOLERANCE_TWO,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(
      s_newConfigDigest,
      _getSignerAddresses(_getSigners(20)),
      FAULT_TOLERANCE_TWO,
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfVerifyingWithAnUnsetDigest() public {
    s_verifier.deactivateConfig(s_oldConfigDigest);

    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    changePrank(address(s_verifierProxy));
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestInactive.selector, s_reportContext[0]));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_canVerifyOlderReportsWithOlderConfigs() public {
    bytes memory signedReport = _generateV1EncodedBlob(
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
    bytes memory signedReport = _generateV1EncodedBlob(
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
    bytes memory signedReport = _generateV1EncodedBlob(
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

  function test_verifyAfterConfigUpdate() public {
    // Get initial signers and set initial config
    address[] memory initialSigners = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    s_verifier.setConfig(configDigest, initialSigners, 4, new Common.AddressAndWeight[](0));

    // Get new signers and update config
    address[] memory newSigners = _getSignerAddresses(_getSigners(20));
    s_verifier.updateConfig(configDigest, initialSigners, newSigners, 6);

    // Verify report with new signers should pass
    s_reportContext[0] = configDigest;
    bytes memory signedReportNewSigners = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(7) // More than f=6 signers
    );

    bytes memory response = s_verifierProxy.verify(signedReportNewSigners, bytes(""));
    assertReportsEqual(response, s_testReportOne);

    // Verify report with old signers should fail
    bytes memory signedReportOldSigners = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(5) // Old number of signers
    );
    vm.expectRevert(abi.encodeWithSelector(Verifier.IncorrectSignatureCount.selector, 5, 7));

    s_verifierProxy.verify(signedReportOldSigners, bytes(""));
  }

  function test_verifyAfterConfigUpdateWithExistingSigners() public {
    // Get initial signers and set initial config
    address[] memory signers = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    s_verifier.setConfig(configDigest, signers, 4, new Common.AddressAndWeight[](0));

    // Update config with same signers and f
    s_verifier.updateConfig(configDigest, signers, signers, 4);

    // Verify report should pass
    s_reportContext[0] = configDigest;
    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(5) // More than f=4 signers
    );

    bytes memory response = s_verifierProxy.verify(signedReport, bytes(""));
    assertReportsEqual(response, s_testReportOne);
  }
}
