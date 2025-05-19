// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTestWithMultipleConfiguredDigests} from "./BaseVerifierTest.t.sol";
import {Verifier} from "../../Verifier.sol";

contract VerificationdeactivateConfigWhenThereAreMultipleDigestsTestV05 is BaseTestWithMultipleConfiguredDigests {
  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");

    changePrank(USER);
    s_verifier.deactivateConfig(bytes32(""));
  }

  function test_revertsIfRemovingAnEmptyDigest() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestEmpty.selector));
    s_verifier.deactivateConfig(bytes32(""));
  }

  function test_revertsIfRemovingAnNonExistentDigest() public {
    vm.expectRevert(abi.encodeWithSelector(Verifier.DigestNotSet.selector, bytes32("mock-digest")));
    s_verifier.deactivateConfig(bytes32("mock-digest"));
  }
}
