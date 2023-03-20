// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {IVerifier} from "../../../../src/v0.8/interfaces/IVerifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";
import {AccessControllerInterface} from "../../../../src/v0.8/interfaces/AccessControllerInterface.sol";

contract VerifierProxyInitializeVerifierTest is BaseTest {
    bytes32 latestDigest;

    function setUp() public override {
        BaseTest.setUp();
        Signer[] memory signers = _getSigners(MAX_ORACLES);
        vm.prank(ADMIN);
        s_verifier.setConfig(
            FEED_ID,
            _getSignerAddresses(signers),
            s_offchaintransmitters,
            FAULT_TOLERANCE,
            bytes(""),
            VERIFIER_VERSION,
            bytes("")
        );
        (, , latestDigest) = s_verifier.latestConfigDetails(FEED_ID);
    }

    function test_revertsIfNotOwner() public {
        vm.expectRevert("Only callable by owner");
        s_verifierProxy.initializeVerifier(latestDigest, address(s_verifier));
    }

    function test_revertsIfZeroAddress() public {
        vm.expectRevert(
            abi.encodeWithSelector(VerifierProxy.ZeroAddress.selector)
        );
        vm.prank(ADMIN);
        s_verifierProxy.initializeVerifier(latestDigest, address(0));
    }

    function test_revertsIfDigestAlreadySet() public {
        vm.prank(ADMIN);
        s_verifierProxy.initializeVerifier(latestDigest, address(s_verifier));
        vm.expectRevert(
            abi.encodeWithSelector(
                VerifierProxy.ConfigDigestAlreadySet.selector,
                latestDigest,
                address(s_verifier)
            )
        );
        vm.prank(ADMIN);
        s_verifierProxy.initializeVerifier(latestDigest, address(s_verifier));
    }

    function test_correctlySetsVerifier() public {
        vm.prank(ADMIN);
        s_verifierProxy.initializeVerifier(latestDigest, address(s_verifier));
        address verifier = s_verifierProxy.getVerifier(latestDigest);
        assertEq(verifier, address(s_verifier));
    }
}
