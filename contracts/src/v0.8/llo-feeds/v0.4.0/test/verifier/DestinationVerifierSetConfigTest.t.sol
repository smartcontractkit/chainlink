// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationRewardManager} from "../../../v0.4.0/DestinationRewardManager.sol";
import {Common} from "../../../libraries/Common.sol";

contract DestinationVerifierSetConfigTest is BaseTest {
  function setUp() public virtual override {
    BaseTest.setUp();
  }

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    changePrank(USER);
    s_verifier.setConfig(_getSignerAddresses(signers), FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfSetWithTooManySigners() public {
    address[] memory signers = new address[](MAX_ORACLES + 1);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.ExcessSigners.selector, signers.length, MAX_ORACLES));
    s_verifier.setConfig(signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfFaultToleranceIsZero() public {
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.FaultToleranceMustBePositive.selector));
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    s_verifier.setConfig(_getSignerAddresses(signers), 0, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfNotEnoughSigners() public {
    address[] memory signers = new address[](2);
    signers[0] = address(1000);
    signers[1] = address(1001);

    vm.expectRevert(
      abi.encodeWithSelector(DestinationVerifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1)
    );
    s_verifier.setConfig(signers, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfDuplicateSigners() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = signerAddrs[1];
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.NonUniqueSignatures.selector));
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_revertsIfSignerContainsZeroAddress() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    signerAddrs[0] = address(0);
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.ZeroAddress.selector));
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_donConfigIdIsSameForSignersInDifferentOrder() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);

    bytes24 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    vm.warp(block.timestamp + 1);

    address temp = signerAddrs[0];
    signerAddrs[0] = signerAddrs[1];
    signerAddrs[1] = temp;

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.DonConfigAlreadyExists.selector, expectedDonConfigId));

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }

  function test_NoDonConfigAlreadyExists() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

    vm.warp(block.timestamp + 1);

    // testing adding same set of Signers but different FAULT_TOLERENCE does not result in DonConfigAlreadyExists revert
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0));

    vm.warp(block.timestamp + 1);

    // testing adding a different set of Signers with same FAULT_TOLERENCE does not result in DonConfigAlreadyExists revert
    address[] memory signerAddrsMinusOne = new address[](signerAddrs.length - 1);
    for (uint256 i = 0; i < signerAddrs.length - 1; i++) {
      signerAddrsMinusOne[i] = signerAddrs[i];
    }
    s_verifier.setConfig(signerAddrsMinusOne, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0));
  }

  function test_addressesAndWeightsDoNotProduceSideEffectsInDonConfigIds() public {
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);

    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    vm.warp(block.timestamp + 1);

    bytes24 expectedDonConfigId = _donConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.DonConfigAlreadyExists.selector, expectedDonConfigId));

    // Same call to setConfig with different addressAndWeights do not entail a new DonConfigID
    // Resulting in a DonConfigAlreadyExists error
    Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
    weights[0] = Common.AddressAndWeight(signers[0].signerAddress, 1);
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
  }

  function test_setConfigActiveUnknownDonConfigId() public {
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.DonConfigDoesNotExist.selector));
    s_verifier.setConfigActive(3, true);
  }

  function test_setConfigWithActivationTime() public {
    // simple case setting a config with specific activation time
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    uint32 activationTime = 10;
    s_verifier.setConfigWithActivationTime(
      signerAddrs,
      FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0),
      activationTime
    );
  }

  function test_setConfigWithActivationTimeNoFutureTimeShouldFail() public {
    // calling setConfigWithActivationTime with a future timestamp should fail
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    uint32 activationTime = uint32(block.timestamp) + 100;
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadActivationTime.selector));
    s_verifier.setConfigWithActivationTime(
      signerAddrs,
      FAULT_TOLERANCE,
      new Common.AddressAndWeight[](0),
      activationTime
    );
  }

  function test_setConfigWithActivationTimeEarlierThanLatestConfigShouldFail() public {
    // setting a config older than the latest current config should fail
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    uint32 oldActivationTime = uint32(block.timestamp) - 1;
    // sets a config with timestamp = block.timestamp
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    // setting a config with ealier timestamp retuls in failure
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadActivationTime.selector));
    s_verifier.setConfigWithActivationTime(
      signerAddrs,
      FAULT_TOLERANCE - 1,
      new Common.AddressAndWeight[](0),
      oldActivationTime
    );
  }

  function test_setConfigWithActivationTimeTheSameAsLatestConfigShouldFail() public {
    // setting a config older than the latest current config should fail
    Signer[] memory signers = _getSigners(MAX_ORACLES);
    address[] memory signerAddrs = _getSignerAddresses(signers);
    // sets a config with timestamp = block.timestamp
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    // setting a config with ealier timestamp retuls in failure
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.BadActivationTime.selector));
    s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
  }
}
