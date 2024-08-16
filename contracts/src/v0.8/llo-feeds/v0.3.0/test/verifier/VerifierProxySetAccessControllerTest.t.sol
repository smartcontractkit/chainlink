// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseVerifierTest.t.sol";
import {AccessControllerInterface} from "../../../../shared/interfaces/AccessControllerInterface.sol";

contract VerifierProxySetAccessControllerTest is BaseTest {
  event AccessControllerSet(address oldAccessController, address newAccessController);

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");

    changePrank(USER);
    s_verifierProxy.setAccessController(AccessControllerInterface(ACCESS_CONTROLLER_ADDRESS));
  }

  function test_successfullySetsNewAccessController() public {
    s_verifierProxy.setAccessController(AccessControllerInterface(ACCESS_CONTROLLER_ADDRESS));
    AccessControllerInterface ac = s_verifierProxy.s_accessController();
    assertEq(address(ac), ACCESS_CONTROLLER_ADDRESS);
  }

  function test_successfullySetsNewAccessControllerIsEmpty() public {
    s_verifierProxy.setAccessController(AccessControllerInterface(address(0)));
    AccessControllerInterface ac = s_verifierProxy.s_accessController();
    assertEq(address(ac), address(0));
  }

  function test_emitsTheCorrectEvent() public {
    vm.expectEmit(true, false, false, false);
    emit AccessControllerSet(address(0), ACCESS_CONTROLLER_ADDRESS);
    s_verifierProxy.setAccessController(AccessControllerInterface(ACCESS_CONTROLLER_ADDRESS));
  }
}
