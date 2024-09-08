// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierSetConfigFromSourceTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    changePrank(USER);
    s_verifier.setConfigFromSource(
      FEED_ID,
      12345,
      address(12345),
      0,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }
}

contract VerifierSetConfigFromSourceMultipleDigestsTest is BaseTestWithMultipleConfiguredDigests {
  function test_correctlyUpdatesTheDigestInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfigFromSource(
      FEED_ID,
      12345,
      address(12345),
      0,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

    (, , bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));
  }

  function test_correctlyUpdatesDigestsOnMultipleVerifiersInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfigFromSource(
      FEED_ID_2,
      12345,
      address(12345),
      0,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

    (, , bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID_2);
    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));

    s_verifier_2.setConfigFromSource(
      FEED_ID_3,
      12345,
      address(12345),
      0,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

    (, , bytes32 configDigest2) = s_verifier_2.latestConfigDetails(FEED_ID_3);
    address verifierAddr2 = s_verifierProxy.getVerifier(configDigest2);
    assertEq(verifierAddr2, address(s_verifier_2));
  }

  function test_correctlySetsConfigWhenDigestsAreRemoved() public {
    s_verifier.deactivateConfig(FEED_ID, s_configDigestTwo);

    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfigFromSource(
      FEED_ID,
      12345,
      address(s_verifier),
      0,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

    bytes32 expectedConfigDigest = _configDigestFromConfigData(
      FEED_ID,
      12345,
      address(s_verifier),
      s_numConfigsSet + 1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);

    assertEq(configCount, s_numConfigsSet + 1);
    assertEq(blockNumber, block.number);
    assertEq(configDigest, expectedConfigDigest);
  }
}
