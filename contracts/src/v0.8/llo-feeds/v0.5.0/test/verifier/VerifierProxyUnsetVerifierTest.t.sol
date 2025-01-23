// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest, BaseTestWithConfiguredVerifierAndFeeManager} from "./BaseVerifierTest.t.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";

contract VerifierProxyUnsetVerifierTest is BaseTest {
  function test_revertsIfNotAdmin() public {
    vm.expectRevert("Only callable by owner");

    changePrank(USER);
    s_verifierProxy.unsetVerifier(bytes32(""));
  }

  function test_revertsIfDigestDoesNotExist() public {
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.VerifierNotFound.selector, bytes32("")));
    s_verifierProxy.unsetVerifier(bytes32(""));
  }
}

contract VerifierProxyUnsetVerifierWithPreviouslySetVerifierTest is BaseTestWithConfiguredVerifierAndFeeManager {
  bytes32 internal s_configDigest;

  event VerifierUnset(bytes32 configDigest, address verifierAddr);

  function setUp() public override {
    BaseTestWithConfiguredVerifierAndFeeManager.setUp();
    s_configDigest = v3ConfigDigest;
  }

  function test_correctlyUnsetsVerifier() public {
    s_verifierProxy.unsetVerifier(s_configDigest);
    address verifierAddr = s_verifierProxy.getVerifier(s_configDigest);
    assertEq(verifierAddr, address(0));
  }

  function test_emitsAnEventAfterUnsettingVerifier() public {
    vm.expectEmit(true, false, false, false);
    emit VerifierUnset(s_configDigest, address(s_verifier));
    s_verifierProxy.unsetVerifier(s_configDigest);
  }
}
