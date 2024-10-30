// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "../BaseTest.t.sol";
import {Ownable2Step} from "../../access/Ownable2Step.sol";

contract Ownable2Step_setup is BaseTest {
  Ownable2StepHelper internal s_ownable2Step;

  function setUp() public override {
    super.setUp();
    s_ownable2Step = new Ownable2StepHelper(OWNER, address(0));
  }
}

contract Ownable2Step_constructor is Ownable2Step_setup {
  function test_constructor_success() public view {
    assertEq(OWNER, s_ownable2Step.owner());
  }

  function test_constructor_OwnerCannotBeZero_reverts() public {
    vm.expectRevert(Ownable2Step.OwnerCannotBeZero.selector);
    new Ownable2Step(address(0), address(0));
  }
}

contract Ownable2Step_transferOwnership is Ownable2Step_setup {
  function test_transferOwnership_success() public {
    vm.expectEmit();
    emit Ownable2Step.OwnershipTransferRequested(OWNER, STRANGER);

    s_ownable2Step.transferOwnership(STRANGER);

    assertTrue(STRANGER != s_ownable2Step.owner());

    vm.startPrank(STRANGER);
    s_ownable2Step.acceptOwnership();
  }

  function test_transferOwnership_CannotTransferToSelf_reverts() public {
    vm.expectRevert(Ownable2Step.CannotTransferToSelf.selector);
    s_ownable2Step.transferOwnership(OWNER);
  }
}

contract Ownable2Step_acceptOwnership is Ownable2Step_setup {
  function test_acceptOwnership_success() public {
    s_ownable2Step.transferOwnership(STRANGER);

    assertTrue(STRANGER != s_ownable2Step.owner());

    vm.startPrank(STRANGER);

    vm.expectEmit();
    emit Ownable2Step.OwnershipTransferred(OWNER, STRANGER);

    s_ownable2Step.acceptOwnership();

    assertEq(STRANGER, s_ownable2Step.owner());
  }

  function test_acceptOwnership_MustBeProposedOwner_reverts() public {
    vm.expectRevert(Ownable2Step.MustBeProposedOwner.selector);
    s_ownable2Step.acceptOwnership();
  }
}

contract Ownable2StepHelper is Ownable2Step {
  constructor(address newOwner, address pendingOwner) Ownable2Step(newOwner, pendingOwner) {}

  function validateOwnership() external view {
    _validateOwnership();
  }
}

contract Ownable2Step_onlyOwner is Ownable2Step_setup {
  function test_onlyOwner_success() public view {
    s_ownable2Step.validateOwnership();
  }

  function test_onlyOwner_OnlyCallableByOwner_reverts() public {
    vm.stopPrank();

    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_ownable2Step.validateOwnership();
  }
}
