// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {VerifierProxy} from "../../../v0.3.0/VerifierProxy.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";

contract VerifierProxyConstructorTest is BaseTest {
  function test_correctlySetsTheOwner() public {
    VerifierProxy proxy = new VerifierProxy(AccessControllerInterface(address(0)));
    assertEq(proxy.owner(), ADMIN);
  }

  function test_correctlySetsTheCorrectAccessControllerInterface() public {
    address accessControllerAddr = address(1234);
    VerifierProxy proxy = new VerifierProxy(AccessControllerInterface(accessControllerAddr));
    assertEq(address(proxy.s_accessController()), accessControllerAddr);
  }

  function test_correctlySetsVersion() public view {
    string memory version = s_verifierProxy.typeAndVersion();
    assertEq(version, "VerifierProxy 2.0.0");
  }
}
