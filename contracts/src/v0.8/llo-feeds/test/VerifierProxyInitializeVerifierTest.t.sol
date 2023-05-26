// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {IVerifier} from "../interfaces/IVerifier.sol";
import {VerifierProxy} from "..//VerifierProxy.sol";
import {AccessControllerInterface} from "../../interfaces/AccessControllerInterface.sol";

contract VerifierProxyInitializeVerifierTest is BaseTest {
  bytes32 latestDigest;

  function setUp() public override {
    BaseTest.setUp();
    Signer[] memory signers = _getSigners(MAX_ORACLES);
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
