// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {Common} from "../../../libraries/Common.sol";
import {MockConfigurator} from "../mocks/MockConfigurator.sol";

contract VerifierSetConfigTestV05 is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    changePrank(USER);
    s_verifier.setConfig(configDigest, _getSignerAddresses(signers), FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfSetWithTooManySigners() public {
    address[] memory signers = new address[](MAX_ORACLES + 1);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    vm.expectRevert(abi.encodeWithSelector(Verifier.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_verifier.setConfig(configDigest, signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfFaultToleranceIsZero() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      0,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    vm.expectRevert(abi.encodeWithSelector(Verifier.FaultToleranceMustBePositive.selector));
    s_verifier.setConfig(configDigest, _getSignerAddresses(signers), 0, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfNotEnoughSigners() public {
    address[] memory signers = new address[](2);
    signers[0] = address(1000);
    signers[1] = address(1001);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    vm.expectRevert(
      abi.encodeWithSelector(Verifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_verifier.setConfig(configDigest, signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = signerAddrs[1];

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    s_verifier.setConfig(configDigest, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfSignerContainsZeroAddress() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = address(0);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    s_verifier.setConfig(configDigest, signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_correctlyUpdatesTheConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    s_verifierProxy.initializeVerifier(address(s_verifier));

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest, _getSignerAddresses(signers), FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    uint32 blockNumber = s_verifier.latestConfigDetails(configDigest);
    assertEq(blockNumber, block.number);
  }
}

contract VerifierUpdateConfigTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();

    s_verifierProxy.initializeVerifier(address(s_verifier));
  }

  function test_updateConfig() public {
    // Get initial signers and config digest
    address[] memory signerAddresses = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test222");

    // Set initial config
    s_verifier.setConfig(configDigest, signerAddresses, 4, new Common.AddressAndWeight[](0));

    // Unset the config
    s_verifier.updateConfig(configDigest, signerAddresses, signerAddresses, 4);
  }

  function test_updateConfigRevertsIfFIsZero() public {
    // Get initial signers and config digest
    address[] memory signerAddresses = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    // Set initial config
    s_verifier.setConfig(configDigest, signerAddresses, 4, new Common.AddressAndWeight[](0));

    // Try to update with f=0
    vm.expectRevert(Verifier.FaultToleranceMustBePositive.selector);
    s_verifier.updateConfig(configDigest, signerAddresses, signerAddresses, 0);
  }

  function test_updateConfigRevertsIfFTooHigh() public {
    // Get initial signers and config digest
    address[] memory signerAddresses = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    // Set initial config
    s_verifier.setConfig(configDigest, signerAddresses, 4, new Common.AddressAndWeight[](0));

    // Try to update with f too high
    vm.expectRevert(abi.encodeWithSelector(Verifier.InsufficientSigners.selector, signerAddresses.length, 46));
    s_verifier.updateConfig(configDigest, signerAddresses, signerAddresses, 15);
  }

  function test_updateConfigWithDifferentSigners() public {
    // Get initial signers and config digest
    address[] memory initialSigners = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    // Set initial config
    s_verifier.setConfig(configDigest, initialSigners, 4, new Common.AddressAndWeight[](0));

    // Get new signers
    address[] memory newSigners = _getSignerAddresses(_getSigners(20));

    // Update config with new signers
    s_verifier.updateConfig(configDigest, initialSigners, newSigners, 6);

    // Verify config was updated
    uint32 blockNumber = s_verifier.latestConfigDetails(configDigest);
    assertEq(blockNumber, block.number);
  }

  function test_updateConfigRevertsIfDigestNotSet() public {
    address[] memory signerAddresses = _getSignerAddresses(_getSigners(15));
    bytes32 nonExistentDigest = keccak256("nonexistent");

    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestNotSet.selector, nonExistentDigest));
    s_verifier.updateConfig(nonExistentDigest, signerAddresses, signerAddresses, 4);
  }

  function test_updateConfigRevertsIfPrevSignersLengthMismatch() public {
    // Get initial signers and config digest
    address[] memory initialSigners = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    // Set initial config
    s_verifier.setConfig(configDigest, initialSigners, 4, new Common.AddressAndWeight[](0));

    // Try to update with wrong number of previous signers
    address[] memory wrongPrevSigners = _getSignerAddresses(_getSigners(10));
    address[] memory newSigners = _getSignerAddresses(_getSigners(15));

    vm.expectRevert(Verifier.NonUniqueSignatures.selector);
    s_verifier.updateConfig(configDigest, wrongPrevSigners, newSigners, 4);
  }

  function test_updateConfigRevertsIfCalledByNonOwner() public {
    address[] memory signerAddresses = _getSignerAddresses(_getSigners(15));
    bytes32 configDigest = keccak256("test");

    // Set initial config
    s_verifier.setConfig(configDigest, signerAddresses, 4, new Common.AddressAndWeight[](0));

    // Try to update as non-owner
    changePrank(USER);
    vm.expectRevert("Only callable by owner");
    s_verifier.updateConfig(configDigest, signerAddresses, signerAddresses, 4);
  }
}

contract VerifierSetConfigWhenThereAreMultipleDigestsTest05 is BaseTestWithMultipleConfiguredDigests {
  function test_correctlyUpdatesTheDigestInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest, _getSignerAddresses(newSigners), 4, new Common.AddressAndWeight[](0));

    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));
  }

  function test_correctlyUpdatesDigestsOnMultipleVerifiersInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID_2,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest, _getSignerAddresses(newSigners), 4, new Common.AddressAndWeight[](0));

    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));

    bytes32 configDigest2 = _configDigestFromConfigData(
      FEED_ID_3,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier_2.setConfig(configDigest2, _getSignerAddresses(newSigners), 4, new Common.AddressAndWeight[](0));

    address verifierAddr2 = s_verifierProxy.getVerifier(configDigest2);
    assertEq(verifierAddr2, address(s_verifier_2));
  }

  function test_correctlySetsConfigWhenDigestsAreRemoved() public {
    s_verifier.deactivateConfig(s_configDigestTwo);

    Signer[] memory newSigners = _getSigners(15);

    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest, _getSignerAddresses(newSigners), 4, new Common.AddressAndWeight[](0));

    uint32 blockNumber = s_verifier.latestConfigDetails(configDigest);

    assertEq(blockNumber, block.number);
  }

  function test_revertsIfDuplicateConfigIsSet() public {
    // Set initial config
    bytes32 configDigest = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest, _getSignerAddresses(_getSigners(15)), 4, new Common.AddressAndWeight[](0));

    // Try to set same config again
    vm.expectRevert(abi.encodeWithSelector(Verifier.ConfigDigestAlreadySet.selector));
    s_verifier.setConfig(configDigest, _getSignerAddresses(_getSigners(15)), 4, new Common.AddressAndWeight[](0));
  }

  function test_incrementalConfigUpdates() public {
    // Set initial config
    bytes32 configDigest1 = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest1, _getSignerAddresses(_getSigners(15)), 4, new Common.AddressAndWeight[](0));

    // Set second config
    bytes32 configDigest2 = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      2,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest2, _getSignerAddresses(_getSigners(15)), 4, new Common.AddressAndWeight[](0));

    // Set third config
    bytes32 configDigest3 = _configDigestFromConfigData(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      3,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    s_verifier.setConfig(configDigest3, _getSignerAddresses(_getSigners(15)), 4, new Common.AddressAndWeight[](0));
  }

  function test_configDigestMatchesConfiguratorDigest() public {
    MockConfigurator configurator = new MockConfigurator();

    // Convert addresses to bytes array
    Signer[] memory signers = _getSigners(15);
    bytes[] memory signersAsBytes = new bytes[](signers.length);
    for (uint i; i < signers.length; ++i) {
      signersAsBytes[i] = abi.encodePacked(signers[i].signerAddress);
    }

    configurator.setStagingConfig(
      FEED_ID,
      signersAsBytes,
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    bytes32 expectedConfigDigest = _configDigestFromConfigData(
      FEED_ID,
      block.chainid,
      address(configurator),
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes("")
    );

    (, , bytes32 configDigest) = configurator.s_configurationStates(FEED_ID);

    assertEq(configDigest, expectedConfigDigest);
  }
}
