// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {IVerifier} from "../../interfaces/IVerifier.sol";
import {VerifierProxy} from "../../VerifierProxy.sol";
import {AccessControllerInterface} from "../../../interfaces/AccessControllerInterface.sol";

contract VerifierProxyConstructorTest is BaseTest {
  function test_correctlySetsTheOwner() public {
    VerifierProxy proxy = new VerifierProxy(
      AccessControllerInterface(address(0)),
      s_feeManager,
      s_rewardManager,
      address(s_wrappedNative)
    );
    assertEq(proxy.owner(), ADMIN);
  }

  function test_correctlySetsTheCorrectAccessControllerInterface() public {
    address accessControllerAddr = address(1234);
    VerifierProxy proxy = new VerifierProxy(
      AccessControllerInterface(accessControllerAddr),
      s_feeManager,
      s_rewardManager,
      address(s_wrappedNative)
    );
    assertEq(address(proxy.getAccessController()), accessControllerAddr);
  }

  function test_correctlySetsVersion() public {
    string memory version = s_verifierProxy.typeAndVersion();
    assertEq(version, "VerifierProxy 1.0.0");
  }
}
