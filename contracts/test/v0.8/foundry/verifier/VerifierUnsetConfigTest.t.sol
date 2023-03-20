// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTestWithConfiguredVerifier, BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../../../src/v0.8/Verifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";

contract VerificationdeactivateConfigWhenThereAreMultipleDigestsTest is
    BaseTestWithMultipleConfiguredDigests
{
    function test_revertsIfCalledByNonOwner() public {
        vm.expectRevert("Only callable by owner");
        s_verifier.deactivateConfig(FEED_ID, bytes32(""));
    }

    function test_revertsIfRemovingAnEmptyDigest() public {
        vm.prank(ADMIN);
        vm.expectRevert(abi.encodeWithSelector(Verifier.DigestEmpty.selector));
        s_verifier.deactivateConfig(FEED_ID, bytes32(""));
    }

    function test_revertsIfRemovingAnNonExistentDigest() public {
        vm.prank(ADMIN);
        vm.expectRevert(
            abi.encodeWithSelector(
                Verifier.DigestNotSet.selector,
                FEED_ID,
                bytes32("mock-digest")
            )
        );
        s_verifier.deactivateConfig(FEED_ID, bytes32("mock-digest"));
    }

    function test_correctlyRemovesAMiddleDigest() public {
        vm.prank(ADMIN);
        s_verifier.deactivateConfig(FEED_ID, s_configDigestTwo);
        (, , bytes32 lastConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
        assertEq(lastConfigDigest, s_configDigestThree);
    }

    function test_correctlyRemovesTheFirstDigest() public {
        vm.prank(ADMIN);
        s_verifier.deactivateConfig(FEED_ID, s_configDigestOne);
        (, , bytes32 lastConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
        assertEq(lastConfigDigest, s_configDigestThree);
    }

    function test_correctlyUnsetsDigestsInSequence() public {
        // Delete config digest 2
        vm.prank(ADMIN);
        s_verifier.deactivateConfig(FEED_ID, s_configDigestTwo);
        (, , bytes32 lastConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
        assertEq(lastConfigDigest, s_configDigestThree);

        // Delete config digest 1
        vm.prank(ADMIN);
        s_verifier.deactivateConfig(FEED_ID, s_configDigestOne);
        (, , lastConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
        assertEq(lastConfigDigest, s_configDigestThree);

        // Delete config digest 3
        vm.prank(ADMIN);
        vm.expectRevert(
            abi.encodeWithSelector(
                Verifier.CannotDeactivateLatestConfig.selector,
                FEED_ID,
                s_configDigestThree
            )
        );
        s_verifier.deactivateConfig(FEED_ID, s_configDigestThree);
        (, , lastConfigDigest) = s_verifier.latestConfigDetails(FEED_ID);
        assertEq(lastConfigDigest, s_configDigestThree);
    }
}
