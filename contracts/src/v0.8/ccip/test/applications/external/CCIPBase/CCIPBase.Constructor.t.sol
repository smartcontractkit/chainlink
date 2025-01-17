// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";

import {OnRampSetup} from "../../../onRamp/OnRamp/OnRampSetup.t.sol";
import {CCIPReceiverReverting} from "../../../helpers/receivers/CCIPReceiverReverting.sol";

contract CCIPBase_Constructor is OnRampSetup {
  function test_Constructor() public {
    CCIPReceiverReverting revertingReceiver = new CCIPReceiverReverting(address(s_destRouter));

    // Check that the router is set correctly
    assertEq(address(revertingReceiver.getRouter()), address(s_destRouter));
  }

  function test_Constructor_RevertWhen_ZeroAddressNotAllowed() public {
    vm.expectRevert(CCIPBase.ZeroAddressNotAllowed.selector);

    new CCIPReceiverReverting(address(0));
  }
}
