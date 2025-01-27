// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPReceiver} from "../../../../applications/external/CCIPReceiver.sol";
import {OnRampSetup} from "../../../onRamp/OnRamp/OnRampSetup.t.sol";

contract CCIPBase_Constructor is OnRampSetup {
  function test_Constructor() public {
    CCIPReceiver revertingReceiver = new CCIPReceiver(address(s_destRouter));

    // Check that the router is set correctly
    assertEq(address(revertingReceiver.getRouter()), address(s_destRouter));
  }

  function test_Constructor_RevertWhen_ZeroAddressNotAllowed() public {
    vm.expectRevert(CCIPBase.ZeroAddressNotAllowed.selector);

    new CCIPReceiver(address(0));
  }
}
