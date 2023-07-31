// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../Verifier.sol";
import {VerifierProxy} from "../VerifierProxy.sol";

contract VerifierSetConfigTest is BaseTest {
  function setUp() public virtual override {
    persistConfig = true;
    BaseTest.setUp();
  }

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
    bytes memory onchainConfigBytes = bytes("onchain config");
    bytes memory offchainConfigBytes = bytes("offchain config");

    s_verifierProxy.initializeVerifier(address(s_verifier));
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      onchainConfigBytes,
      VERIFIER_VERSION,
      offchainConfigBytes
    );

    bytes32 expectedConfigDigest = _configDigestFromConfigData(
      FEED_ID,
      address(s_verifier),
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      onchainConfigBytes,
      VERIFIER_VERSION,
      offchainConfigBytes
    );

    //change to an EOA to request the config
    changePrank(tx.origin);
    Verifier.ActiveConfig memory latestConfig = s_verifier.latestConfig(FEED_ID);
    changePrank(msg.sender);

    //check the latest config matches
    assertEq(latestConfig.previousConfigBlockNumber, 0);
    assertEq(latestConfig.currentConfigBlockNumber, 12345);
    assertEq(latestConfig.configDigest, expectedConfigDigest);
    assertEq(latestConfig.configCount, 1);
    for (uint256 i; i < signers.length; i++) {
      assertEq(latestConfig.signers[i], _getSignerAddresses(signers)[i]);
    }
    for (uint256 i; i < s_offchaintransmitters.length; i++) {
      assertEq(latestConfig.transmitters[i], s_offchaintransmitters[i]);
    }
    assertEq(latestConfig.f, FAULT_TOLERANCE);
    assertEq(latestConfig.onchainConfig, onchainConfigBytes);
    assertEq(latestConfig.offchainConfigVersion, VERIFIER_VERSION);
    assertEq(latestConfig.offchainConfig, offchainConfigBytes);

    (uint32 configCount, uint32 blockNumber, bytes32 configDigest) = s_verifier.latestConfigDetails(FEED_ID);
    assertEq(configCount, 1);
    assertEq(blockNumber, block.number);
    assertEq(configDigest, expectedConfigDigest);

    (bool scanLogs, bytes32 configDigestTwo, uint32 epoch) = s_verifier.latestConfigDigestAndEpoch(FEED_ID);
    assertEq(scanLogs, false);
    assertEq(configDigestTwo, expectedConfigDigest);
    assertEq(epoch, 0);
  }

  function test_setConfigUpdatesLatest() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    bytes memory onchainConfigBytes = bytes("onchain config");
    bytes memory offchainConfigBytes = bytes("offchain config");

    s_verifierProxy.initializeVerifier(address(s_verifier));
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      onchainConfigBytes,
      VERIFIER_VERSION,
      offchainConfigBytes
    );

    bytes32 expectedConfigDigest = _configDigestFromConfigData(
      FEED_ID,
      address(s_verifier),
      2,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      onchainConfigBytes,
      VERIFIER_VERSION,
      offchainConfigBytes
    );

    //change to an EOA to request the config
    changePrank(tx.origin);
    Verifier.ActiveConfig memory latestConfig = s_verifier.latestConfig(FEED_ID);
    changePrank(msg.sender);

    //check the latest config matches
    assertEq(latestConfig.previousConfigBlockNumber, 12345);
    assertEq(latestConfig.currentConfigBlockNumber, 12345);
    assertEq(latestConfig.configDigest, expectedConfigDigest);
    assertEq(latestConfig.configCount, 2);
    for (uint256 i; i < signers.length; i++) {
      assertEq(latestConfig.signers[i], _getSignerAddresses(signers)[i]);
    }
    for (uint256 i; i < s_offchaintransmitters.length; i++) {
      assertEq(latestConfig.transmitters[i], s_offchaintransmitters[i]);
    }
    assertEq(latestConfig.f, FAULT_TOLERANCE);
    assertEq(latestConfig.onchainConfig, onchainConfigBytes);
    assertEq(latestConfig.offchainConfigVersion, VERIFIER_VERSION);
    assertEq(latestConfig.offchainConfig, offchainConfigBytes);
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

contract VerifierProxyTestSetConfigWithoutPersist is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
  }

  function test_setConfigDoesNotPersist() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    bytes memory onchainConfigBytes = bytes("onchain config");
    bytes memory offchainConfigBytes = bytes("offchain config");

    s_verifierProxy.initializeVerifier(address(s_verifier));
    s_verifier.setConfig(
      FEED_ID,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      onchainConfigBytes,
      VERIFIER_VERSION,
      offchainConfigBytes
    );

    //change to an EOA to request the config
    changePrank(tx.origin);
    Verifier.ActiveConfig memory latestConfig = s_verifier.latestConfig(FEED_ID);
    changePrank(msg.sender);

    //check the latest config is null
    assertEq(latestConfig.previousConfigBlockNumber, 0);
    assertEq(latestConfig.currentConfigBlockNumber, 0);
    assertEq(latestConfig.configDigest, bytes32(""));
    assertEq(latestConfig.configCount, 0);
    assertEq(latestConfig.signers.length, 0);
    assertEq(latestConfig.transmitters.length, 0);
    assertEq(latestConfig.f, 0);
    assertEq(latestConfig.onchainConfig, bytes(""));
    assertEq(latestConfig.offchainConfigVersion, 0);
    assertEq(latestConfig.offchainConfig, bytes(""));
  }
}
