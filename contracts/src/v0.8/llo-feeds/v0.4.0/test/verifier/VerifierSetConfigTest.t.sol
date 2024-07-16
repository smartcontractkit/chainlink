// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";
import {DestinationRewardManager} from "../../../v0.4.0/DestinationRewardManager.sol";
import {Common} from "../../../libraries/Common.sol";

contract VerifierSetConfigTest is BaseTest {
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
            abi.encodeWithSelector(
                DestinationVerifier.InsufficientSigners.selector, signers.length, FAULT_TOLERANCE * 3 + 1
            )
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

    function test_DONConfigIDIsSameForSignersInDifferentOrder() public {
        Signer[] memory signers = _getSigners(MAX_ORACLES);
        address[] memory signerAddrs = _getSignerAddresses(signers);

        bytes24 expectedDonConfigID = _DONConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

        address temp = signerAddrs[0];
        signerAddrs[0] = signerAddrs[1];
        signerAddrs[1] = temp;

        vm.expectRevert(
            abi.encodeWithSelector(DestinationVerifier.DONConfigAlreadyExists.selector, expectedDonConfigID)
        );

        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
    }

    function test_NoDonConfigAlreadyExists() public {
        Signer[] memory signers = _getSigners(MAX_ORACLES);
        address[] memory signerAddrs = _getSignerAddresses(signers);

        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));

        // testing adding same set of Signers but different FAULT_TOLERENCE does not result in DONConfigAlreadyExists revert
        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE - 1, new Common.AddressAndWeight[](0));

        // testing adding a different set of Signers with same FAULT_TOLERENCE does not result in DONConfigAlreadyExists revert
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

        bytes24 expectedDonConfigID = _DONConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

        vm.expectRevert(
            abi.encodeWithSelector(DestinationVerifier.DONConfigAlreadyExists.selector, expectedDonConfigID)
        );

        // Same call to setConfig with different addressAndWeights do not entail a new DONConfigId
        // Resulting in a DONConfigAlreadyExists error
        Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
        weights[0] = Common.AddressAndWeight(signers[0].signerAddress, 1);
        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    }

    function test_correctlyUpdatesInternalState() public {
        Signer[] memory signers = _getSigners(MAX_ORACLES);
        address[] memory signerAddrs = _getSignerAddresses(signers);

        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
        uint256 t1 = block.timestamp;

        bytes24 expectedDonConfigID1 = _DONConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);

        // Checking expected internal state is updated accordingly

        // 1. check internal state of: s_DONConfigByID
        DestinationVerifier.DONConfig[] memory donConfigs = s_verifier.getDONConfigs();

        assertEq(donConfigs.length, 1);
        DestinationVerifier.DONConfig memory donConfigT1 = donConfigs[0];
        assertEq(donConfigT1.f, FAULT_TOLERANCE);
        assertEq(donConfigT1.isActive, true);
        assertEq(donConfigT1.DONConfigID, expectedDonConfigID1);
        assertEq(donConfigT1.activationTime, t1);

        // 2. check internal state of: s_SignerByAddressAndDONConfigId

        for (uint256 i; i < signers.length; ++i) {
            bytes32 signerToDonConfigKey = _signerAddressAndDonConfigKey(signers[i].signerAddress, expectedDonConfigID1);
            bool c = s_verifier.getSignerConfigByAddressAndDONConfigId(signerToDonConfigKey);
            assertEq(c, true);
        }

        // setConfig again but  this config contains a subset of the signers
        BaseTest.Signer[] memory signers2 = new BaseTest.Signer[](4);
        signers2[0] = signers[0];
        signers2[1] = signers[1];
        signers2[2] = signers[2];
        signers2[3] = signers[3];
        uint8 MINIMAL_FAULT_TOLERANCE = 1;

        address[] memory signerAddrs2 = _getSignerAddresses(signers2);
        s_verifier.setConfig(signerAddrs2, MINIMAL_FAULT_TOLERANCE, new Common.AddressAndWeight[](0));
        uint256 t2 = block.timestamp;

        bytes24 expectedDonConfigID2 = _DONConfigIdFromConfigData(signerAddrs2, 1);
        donConfigs = s_verifier.getDONConfigs();

        // 1. check internal state of s_DONConfigByID
        donConfigT1 = donConfigs[0];
        assertEq(donConfigT1.f, FAULT_TOLERANCE);
        assertEq(donConfigT1.isActive, true);
        assertEq(donConfigT1.DONConfigID, expectedDonConfigID1);
        assertEq(donConfigT1.activationTime, t1);

        assertEq(donConfigs.length, 2);
        DestinationVerifier.DONConfig memory donConfigT2 = donConfigs[1];
        assertEq(donConfigT2.f, MINIMAL_FAULT_TOLERANCE);
        assertEq(donConfigT2.isActive, true);
        assertEq(donConfigT2.DONConfigID, expectedDonConfigID2);
        assertEq(donConfigT2.activationTime, t2);

        // 2. check state of s_SignerByAddressAndDONConfigId
        // checking first DONConfig
        for (uint256 i; i < signers.length; ++i) {
            bytes32 signerToDonConfigKey1 =
                _signerAddressAndDonConfigKey(signers[i].signerAddress, expectedDonConfigID1);
            bool c1 = s_verifier.getSignerConfigByAddressAndDONConfigId(signerToDonConfigKey1);
            assertEq(c1, true);
        }

        // checking second DONConfig
        for (uint256 i; i < signers.length; ++i) {
            bytes32 signerToDonConfigKey2 =
                _signerAddressAndDonConfigKey(signers[i].signerAddress, expectedDonConfigID2);
            bool c2 = s_verifier.getSignerConfigByAddressAndDONConfigId(signerToDonConfigKey2);
            if (i < 4) {
                assertEq(c2, true);
            } else {
                assertEq(c2, false);
            }
        }

        //  setting DONConfig2 as activated false
        s_verifier.setConfigActive(1, false);

        donConfigs = s_verifier.getDONConfigs();
        assertEq(donConfigs.length, 2);

        DestinationVerifier.DONConfig memory donConfig2AtT3 = donConfigs[1];
        assertEq(donConfig2AtT3.f, MINIMAL_FAULT_TOLERANCE);
        assertEq(donConfig2AtT3.isActive, false);
        assertEq(donConfig2AtT3.DONConfigID, expectedDonConfigID2);

        //  setting DONConfig2 as activated false (again)
        s_verifier.setConfigActive(1, false);

        donConfigs = s_verifier.getDONConfigs();
        assertEq(donConfigs.length, 2);
        DestinationVerifier.DONConfig memory donConfig2AtT4 = donConfigs[1];
        assertEq(donConfig2AtT4.f, MINIMAL_FAULT_TOLERANCE);
        assertEq(donConfig2AtT4.isActive, false);
        assertEq(donConfig2AtT4.DONConfigID, expectedDonConfigID2);

        // checking other DONConfigs were not affected
        donConfigs = s_verifier.getDONConfigs();
        assertEq(donConfigs.length, 2);
        DestinationVerifier.DONConfig memory donConfig1AtT5 = donConfigs[0];
        assertEq(donConfig1AtT5.f, FAULT_TOLERANCE);
        assertEq(donConfig1AtT5.isActive, true);
        assertEq(donConfig1AtT5.DONConfigID, expectedDonConfigID1);

        // setting DONConfig2 as activated true
        s_verifier.setConfigActive(1, true);

        donConfigs = s_verifier.getDONConfigs();
        assertEq(donConfigs.length, 2);
        DestinationVerifier.DONConfig memory donConfig2AtT5 = donConfigs[1];
        assertEq(donConfig2AtT5.f, MINIMAL_FAULT_TOLERANCE);
        assertEq(donConfig2AtT5.isActive, true);
        assertEq(donConfig2AtT5.DONConfigID, expectedDonConfigID2);
    }

    function test_setConfigActiveUnknownDONConfigID() public {
        vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.DONConfigDoesNotExist.selector));
        s_verifier.setConfigActive(3, true);
    }

    function test_setConfigWithAddressesAndWeightsAreSetCorrectly() public {
        Signer[] memory signers = _getSigners(MAX_ORACLES);
        address[] memory signerAddrs = _getSignerAddresses(signers);
        address recipient = signers[0].signerAddress;
        Common.AddressAndWeight[] memory weights = new Common.AddressAndWeight[](1);
        weights[0] = Common.AddressAndWeight(recipient, ONE_PERCENT * 100);
        bytes32 expectedDonConfigID = _DONConfigIdFromConfigData(signerAddrs, FAULT_TOLERANCE);
        vm.expectEmit();
        emit RewardRecipientsUpdated(expectedDonConfigID, weights);
        s_verifier.setConfig(signerAddrs, FAULT_TOLERANCE, weights);
    }
}
