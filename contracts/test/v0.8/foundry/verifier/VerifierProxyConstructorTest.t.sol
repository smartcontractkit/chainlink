// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {IVerifier} from "../../../../src/v0.8/interfaces/IVerifier.sol";
import {VerifierProxy} from "../../../../src/v0.8/VerifierProxy.sol";
import {AccessControllerInterface} from "../../../../src/v0.8/interfaces/AccessControllerInterface.sol";

contract VerifierProxyConstructorTest is BaseTest {
  function test_correctlySetsTheOwner() public {
    VerifierProxy proxy = new VerifierProxy(AccessControllerInterface(address(0)));
    assertEq(proxy.owner(), ADMIN);
  }

  function test_correctlySetsTheCorrectAccessControllerInterface() public {
    address accessControllerAddr = address(1234);
    VerifierProxy proxy = new VerifierProxy(AccessControllerInterface(accessControllerAddr));
    assertEq(address(proxy.getAccessController()), accessControllerAddr);
  }

  function test_correctlySetsVersion() public {
    string memory version = s_verifierProxy.typeAndVersion();
    assertEq(version, "VerifierProxy 0.0.1");
  }
}
