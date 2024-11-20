// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";
import {Common} from "../../../libraries/Common.sol";
import {MockConfigurator} from "../mocks/MockConfigurator.sol";

contract VerifierSetConfigTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    changePrank(USER);
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfSetWithTooManySigners() public {
    address[] memory signers = new address[](MAX_ORACLES + 1);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfFaultToleranceIsZero() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.FaultToleranceMustBePositive.selector));
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      0,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfNotEnoughSigners() public {
    address[] memory signers = new address[](2);
    signers[0] = address(1000);
    signers[1] = address(1001);

    vm.expectRevert(
      abi.encodeWithSelector(Verifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signers,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_revertsIfSignerContainsZeroAddress() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(Verifier.ZeroAddress.selector));
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      signerAddrs,
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_correctlyUpdatesTheConfig() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);

    s_verifierProxy.initializeVerifier(address(s_verifier));
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(signers),
      s_offchaintransmitters,
      FAULT_TOLERANCE,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

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

    (uint32 configCount, uint32 blockNumber) = s_verifier.latestConfigDetails(configDigest);
    assertEq(configCount, 1);
    assertEq(blockNumber, block.number);
  }
}

contract VerifierSetConfigWhenThereAreMultipleDigestsTest is BaseTestWithMultipleConfiguredDigests {
  function test_correctlyUpdatesTheDigestInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

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

    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));
  }

  function test_correctlyUpdatesDigestsOnMultipleVerifiersInTheProxy() public {
    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfigFromSource(
      FEED_ID_2,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

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

    address verifierAddr = s_verifierProxy.getVerifier(configDigest);
    assertEq(verifierAddr, address(s_verifier));

    s_verifier_2.setConfigFromSource(
      FEED_ID_3,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(newSigners),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

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
    
    address verifierAddr2 = s_verifierProxy.getVerifier(configDigest2);
    assertEq(verifierAddr2, address(s_verifier_2));
  }

  function test_correctlySetsConfigWhenDigestsAreRemoved() public {
    s_verifier.deactivateConfig(s_configDigestTwo);

    Signer[] memory newSigners = _getSigners(15);

    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS,
      1,
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

    (uint32 configCount, uint32 blockNumber) = s_verifier.latestConfigDetails(configDigest);

    assertEq(configCount, s_numConfigsSet + 1);
    assertEq(blockNumber, block.number);
    assertEq(configDigest, expectedConfigDigest);
  }

  function test_revertsIfDuplicateConfigIsSet() public {
    // Set initial config
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID, 
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

    // Try to set same config again
    vm.expectRevert(abi.encodeWithSelector(Verifier.NonUniqueSignatures.selector));
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS, 
      1,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }

  function test_incrementalConfigUpdates() public {
    // Set initial config
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID, 
      SOURCE_ADDRESS,
      1,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

    // Try to set same config again
    s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS, 
      2,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );

      s_verifier.setConfigFromSource(
      FEED_ID,
      SOURCE_CHAIN_ID,
      SOURCE_ADDRESS, 
      3,
      _getSignerAddresses(_getSigners(15)),
      s_offchaintransmitters,
      4,
      bytes(""),
      VERIFIER_VERSION,
      bytes(""),
      new Common.AddressAndWeight[](0)
    );
  }
  
  function test_configDigestMatchesConfiguratorDigest() public {
     MockConfigurator configurator = new MockConfigurator();

     // Convert addresses to bytes array
     Signer[] memory signers = _getSigners(15);
     bytes[] memory signersAsBytes = new bytes[](signers.length);
     for (uint i; i < signers.length; ++i){
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

    (,,bytes32 configDigest) = configurator.s_configurationStates(FEED_ID);

    assertEq(configDigest, expectedConfigDigest);
  }
}
