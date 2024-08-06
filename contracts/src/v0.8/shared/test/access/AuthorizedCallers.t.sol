// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {AuthorizedCallers} from "../../access/AuthorizedCallers.sol";
import {BaseTest} from "../BaseTest.t.sol";

contract AuthorizedCallers_setup is BaseTest {
  address[] s_callers;

  AuthorizedCallers s_authorizedCallers;

  function setUp() public override {
    super.setUp();
    s_callers.push(makeAddr("caller1"));
    s_callers.push(makeAddr("caller2"));

    s_authorizedCallers = new AuthorizedCallers(s_callers);
  }
}

contract AuthorizedCallers_constructor is AuthorizedCallers_setup {
  event AuthorizedCallerAdded(address caller);

  function test_constructor_Success() public {
    for (uint256 i = 0; i < s_callers.length; ++i) {
      vm.expectEmit();
      emit AuthorizedCallerAdded(s_callers[i]);
    }

    s_authorizedCallers = new AuthorizedCallers(s_callers);

    assertEq(s_callers, s_authorizedCallers.getAllAuthorizedCallers());
  }

  function test_ZeroAddressNotAllowed_Revert() public {
    s_callers[0] = address(0);

    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);

    new AuthorizedCallers(s_callers);
  }
}

contract AuthorizedCallers_applyAuthorizedCallerUpdates is AuthorizedCallers_setup {
  event AuthorizedCallerAdded(address caller);
  event AuthorizedCallerRemoved(address caller);

  function test_OnlyAdd_Success() public {
    address[] memory addedCallers = new address[](2);
    addedCallers[0] = vm.addr(3);
    addedCallers[1] = vm.addr(4);

    address[] memory removedCallers = new address[](0);

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), s_callers);

    vm.expectEmit();
    emit AuthorizedCallerAdded(addedCallers[0]);
    vm.expectEmit();
    emit AuthorizedCallerAdded(addedCallers[1]);

    s_authorizedCallers.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    address[] memory expectedCallers = new address[](4);
    expectedCallers[0] = s_callers[0];
    expectedCallers[1] = s_callers[1];
    expectedCallers[2] = addedCallers[0];
    expectedCallers[3] = addedCallers[1];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_OnlyRemove_Success() public {
    address[] memory addedCallers = new address[](0);
    address[] memory removedCallers = new address[](1);
    removedCallers[0] = s_callers[0];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), s_callers);

    vm.expectEmit();
    emit AuthorizedCallerRemoved(removedCallers[0]);

    s_authorizedCallers.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    address[] memory expectedCallers = new address[](1);
    expectedCallers[0] = s_callers[1];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_AddAndRemove_Success() public {
    address[] memory addedCallers = new address[](2);
    addedCallers[0] = address(42);
    addedCallers[1] = address(43);

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = s_callers[0];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), s_callers);

    vm.expectEmit();
    emit AuthorizedCallerRemoved(removedCallers[0]);
    vm.expectEmit();
    emit AuthorizedCallerAdded(addedCallers[0]);
    vm.expectEmit();
    emit AuthorizedCallerAdded(addedCallers[1]);

    s_authorizedCallers.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    // Order of the set changes on removal
    address[] memory expectedCallers = new address[](3);
    expectedCallers[0] = s_callers[1];
    expectedCallers[1] = addedCallers[0];
    expectedCallers[2] = addedCallers[1];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_RemoveThenAdd_Success() public {
    address[] memory addedCallers = new address[](1);
    addedCallers[0] = s_callers[0];

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = s_callers[0];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), s_callers);

    vm.expectEmit();
    emit AuthorizedCallerRemoved(removedCallers[0]);

    vm.expectEmit();
    emit AuthorizedCallerAdded(addedCallers[0]);

    s_authorizedCallers.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    address[] memory expectedCallers = new address[](2);
    expectedCallers[0] = s_callers[1];
    expectedCallers[1] = s_callers[0];

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), expectedCallers);
  }

  function test_SkipRemove_Success() public {
    address[] memory addedCallers = new address[](0);

    address[] memory removedCallers = new address[](1);
    removedCallers[0] = address(42);

    vm.recordLogs();
    s_authorizedCallers.applyAuthorizedCallerUpdates(
      AuthorizedCallers.AuthorizedCallerArgs({addedCallers: addedCallers, removedCallers: removedCallers})
    );

    assertEq(s_authorizedCallers.getAllAuthorizedCallers(), s_callers);
    assertEq(vm.getRecordedLogs().length, 0);
  }

  function test_OnlyCallableByOwner_Revert() public {
    vm.stopPrank();

    AuthorizedCallers.AuthorizedCallerArgs memory authorizedCallerArgs = AuthorizedCallers.AuthorizedCallerArgs({
      addedCallers: new address[](0),
      removedCallers: new address[](0)
    });

    vm.expectRevert("Only callable by owner");

    s_authorizedCallers.applyAuthorizedCallerUpdates(authorizedCallerArgs);
  }

  function test_ZeroAddressNotAllowed_Revert() public {
    s_callers[0] = address(0);

    vm.expectRevert(AuthorizedCallers.ZeroAddressNotAllowed.selector);

    new AuthorizedCallers(s_callers);
  }
}
