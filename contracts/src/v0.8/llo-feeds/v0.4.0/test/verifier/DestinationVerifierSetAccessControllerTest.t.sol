// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseTest} from "./BaseDestinationVerifierTest.t.sol";

contract DestinationVerifierSetAccessControllerTest is BaseTest {
  event AccessControllerSet(address oldAccessController, address newAccessController);

  function test_revertsIfCalledByNonOwner() public {
    vm.expectRevert("Only callable by owner");

    changePrank(USER);
    s_verifier.setAccessController(ACCESS_CONTROLLER_ADDRESS);
  }

  function test_successfullySetsNewAccessController() public {
    s_verifier.setAccessController(ACCESS_CONTROLLER_ADDRESS);
    address ac = s_verifier.s_accessController();
    assertEq(ac, ACCESS_CONTROLLER_ADDRESS);
  }

  function test_successfullySetsNewAccessControllerIsEmpty() public {
    s_verifier.setAccessController(address(0));
    address ac = s_verifier.s_accessController();
    assertEq(ac, address(0));
  }

  function test_emitsTheCorrectEvent() public {
    vm.expectEmit(true, false, false, false);
    emit AccessControllerSet(address(0), ACCESS_CONTROLLER_ADDRESS);
    s_verifier.setAccessController(ACCESS_CONTROLLER_ADDRESS);
  }
}
