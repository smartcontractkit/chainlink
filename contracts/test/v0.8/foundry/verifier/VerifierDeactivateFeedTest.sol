// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifier, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../../../src/v0.8/Verifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";

contract VerifierActivateFeedTest is BaseTestWithConfiguredVerifier {
  function test_revertsIfNotOwnerActivateFeed() public {
    changePrank(address(s_verifierProxy));
    vm.expectRevert("Only callable by owner");
    s_verifier.activateFeed(FEED_ID);
  }

  function test_revertsIfNotOwnerDeactivateFeed() public {
    changePrank(address(s_verifierProxy));
    vm.expectRevert("Only callable by owner");
    s_verifier.deactivateFeed(FEED_ID);
  }

  function test_revertsIfNoFeedExistsActivate() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.InvalidFeed.selector, INVALID_FEED));
    s_verifier.activateFeed(INVALID_FEED);
  }

  function test_revertsIfNoFeedExistsDeactivate() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.InvalidFeed.selector, INVALID_FEED));
    s_verifier.deactivateFeed(INVALID_FEED);
  }
}

contract VerifierDeactivateFeedWithVerifyTest is BaseTestWithMultipleConfiguredDigests {
  bytes32[3] internal s_reportContext;

  event ConfigActivated(bytes32 configDigest);

  Report internal s_testReportOne;

  function setUp() public override {
    BaseTestWithMultipleConfiguredDigests.setUp();
    s_reportContext[0] = s_configDigestOne;
    s_reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReportOne = _createReport(
      FEED_ID,
      uint32(block.timestamp),
      MEDIAN,
      BID,
      ASK,
      uint64(block.number),
      blockhash(block.number + 3),
      uint64(block.number + 3)
    );

    s_verifier.deactivateFeed(FEED_ID);
  }

  function test_currentReportAllowsVerification() public {
    s_verifier.activateFeed(FEED_ID);
    changePrank(address(s_verifierProxy));

    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_previousReportAllowsVerification() public {
    s_verifier.activateFeed(FEED_ID);
    changePrank(address(s_verifierProxy));

    s_reportContext[0] = s_configDigestTwo;
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE_TWO + 1)
    );
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_currentReportFailsVerification() public {
    changePrank(address(s_verifierProxy));

    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE + 1)
    );

    vm.expectRevert(abi.encodeWithSelector(Verifier.InactiveFeed.selector, FEED_ID));
    s_verifier.verify(signedReport, msg.sender);
  }

  function test_previousReportFailsVerification() public {
    changePrank(address(s_verifierProxy));

    s_reportContext[0] = s_configDigestTwo;
    bytes memory signedReport = _generateEncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE_TWO + 1)
    );

    vm.expectRevert(abi.encodeWithSelector(Verifier.InactiveFeed.selector, FEED_ID));
    s_verifier.verify(signedReport, msg.sender);
  }
}
