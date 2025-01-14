// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTestWithConfiguredVerifierAndFeeManager, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";

contract VerifierActivateConfigTestV05 is BaseTestWithConfiguredVerifierAndFeeManager {
  function test_revertsIfNotOwner() public {
    vm.expectRevert("Only callable by owner");

    changePrank(address(s_verifierProxy));
    s_verifier.activateConfig(bytes32("mock"));
  }

  function test_revertsIfDigestIsEmpty() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestEmpty.selector));
    s_verifier.activateConfig(bytes32(""));
  }

  function test_revertsIfDigestNotSet() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestNotSet.selector, bytes32("non-existent-digest")));
    s_verifier.activateConfig(bytes32("non-existent-digest"));
  }
}

contract VerifierActivateConfigWithDeactivatedConfigTestV05 is BaseTestWithMultipleConfiguredDigests {
  bytes32[3] internal s_reportContext;

  event ConfigActivated(bytes32 configDigest);

  V1Report internal s_testReportOne;

  function setUp() public override {
    BaseTestWithMultipleConfiguredDigests.setUp();
    s_reportContext[0] = s_configDigestTwo;
    s_reportContext[1] = bytes32(abi.encode(uint32(5), uint8(1)));
    s_testReportOne = _createV1Report(
      FEED_ID,
      uint32(block.timestamp),
      MEDIAN,
      BID,
      ASK,
      uint64(block.number),
      blockhash(block.number + 3),
      uint64(block.number + 3),
      uint32(block.timestamp)
    );

    s_verifier.deactivateConfig(s_configDigestTwo);
  }

  function test_allowsVerification() public {
    s_verifier.activateConfig(s_configDigestTwo);
    changePrank(address(s_verifierProxy));

    bytes memory signedReport = _generateV1EncodedBlob(
      s_testReportOne,
      s_reportContext,
      _getSigners(FAULT_TOLERANCE_TWO + 1)
    );
    s_verifier.verify(signedReport, msg.sender);
  }
}
