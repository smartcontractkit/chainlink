// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";
import {DestinationVerifier} from "../../../v0.4.0/DestinationVerifier.sol";

contract DestinationVerifierConstructorTest is BaseTest {
  bytes32[3] internal s_reportContext;

  function test_revertsIfInitializedWithEmptyVerifierProxy() public {
    vm.expectRevert(abi.encodeWithSelector(DestinationVerifier.ZeroAddress.selector));
    new DestinationVerifier(address(0));
  }

  function test_typeAndVersion() public {
    DestinationVerifier v = new DestinationVerifier(address(s_verifierProxy));
    assertEq(v.owner(), ADMIN);
    string memory typeAndVersion = s_verifier.typeAndVersion();
    assertEq(typeAndVersion, "DestinationVerifier 1.0.0");
  }

  function test_falseIfIsNotCorrectInterface() public view {
    bool isInterface = s_verifier.supportsInterface(bytes4("abcd"));
    assertEq(isInterface, false);
  }

  function test_trueIfIsCorrectInterface() public view {
    bool isInterface = s_verifier.supportsInterface(DestinationVerifier.verify.selector);
    assertEq(isInterface, true);
  }
}
