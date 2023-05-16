// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../../../src/v0.8/Verifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";

contract VerifierSetConfigTest is BaseTest {
  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    changePrank(USER);
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
  }

  function test_revertsIfSetWithTooManySigners() public {
    address[] memory signers = new address[](MAX_ORACLES + 1);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_verifier.setConfig(
      FEED_ID,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
  }

  function test_revertsIfFaultToleranceIsZero() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.FaultToleranceMustBePositive.selector));
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      0,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
  }

  function test_revertsIfNotEnoughSigners() public {
    address[] memory signers = new address[](2);
    signers[0] = address(1000);
    signers[1] = address(1001);

    vm.expectRevert(
      abi.encodeWithSelector(Verifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_verifier.setConfig(
      FEED_ID,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
  }

  function test_revertsIfDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    s_verifier.setConfig(
      FEED_ID,
      signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
  }

  function test_revertsIfSignerContainsZeroAddress() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    s_verifier.setConfig(
      FEED_ID,
      signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );
  }

  function test_correctlyUpdatesTheConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    bytes32 expectedConfigDigest = _configDigestFromConfigData(
      FEED_ID,
      address(s_verifier),
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    assertEq(configCount, 1);
    assertEq(blockNumber, block.number);
    assertEq(configDigest, expectedConfigDigest);

    (bool scanLogs, bytes32 configDigestTwo, uint32 epoch) = s_verifier.latestConfigDigestAndEpoch(FEED_ID);
    assertEq(scanLogs, false);
    assertEq(configDigestTwo, expectedConfigDigest);
    assertEq(epoch, 0);
  }
}

contract VerifierSetConfigWhenThereAreMultipleDigestsTest is BaseTestWithMultipleConfiguredDigests {
  function test_correctlyUpdatesTheDigestInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    (, , bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));
  }

  function test_correctlyUpdatesDigestsOnMultipleVerifiersInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfig(
      FEED_ID_2,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    (, , bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID_2);
    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));

    s_verifier_2.setConfig(
      FEED_ID_3,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    (, , bytes32 configDigest2) = s_verifier_2.latestConfigDetails(FEED_ID_3);
    address verifierAddr2 = s_verifierProxy.getVerifier(configDigest2);
    assertEq(verifierAddr2, address(s_verifier_2));
  }

  function test_correctlySetsConfigWhenDigestsAreRemoved() public {
    s_verifier.deactivateConfig(FEED_ID, s_configDigestTwo);

    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    bytes32 expectedConfigDigest = _configDigestFromConfigData(
      FEED_ID,
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
