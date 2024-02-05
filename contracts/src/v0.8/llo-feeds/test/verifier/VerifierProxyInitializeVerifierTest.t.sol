// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";

contract VerifierProxyInitializeVerifierTest is BaseTest {
  bytes32 latestDigest;

  function setUp() public override {
    BaseTest.setUp();
  }

  function test_revertsIfNotOwner() public {
    changePrank(USER);
    vm.expectRevert("Only callable by owner");
    s_verifierProxy.initializeVerifier(address(s_verifier));
  }

  function test_revertsIfZeroAddress() public {
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.ZeroAddress.selector));
    s_verifierProxy.initializeVerifier(address(0));
  }

  function test_revertsIfVerifierAlreadyInitialized() public {
    s_verifierProxy.initializeVerifier(address(s_verifier));
    vm.expectRevert(abi.encodeWithSelector(VerifierProxy.VerifierAlreadyInitialized.selector, address(s_verifier)));
    s_verifierProxy.initializeVerifier(address(s_verifier));
  }
}
