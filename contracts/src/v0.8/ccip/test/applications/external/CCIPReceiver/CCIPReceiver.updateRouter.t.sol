// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

contract CCIPReceiver_updateRouter is CCIPReceiverSetup {
  function test_updateRouter() public {
    address newRouter = address(0x1234);

    vm.expectEmit();
    emit CCIPBase.CCIPRouterModified(address(s_destRouter), newRouter);

    s_receiver.updateRouter(newRouter);

    assertEq(s_receiver.getRouter(), newRouter, "Router Address not set correctly to the new router");
  }

  function test_updateRouter_RevertWhen_ZeroAddressNotAllowed() public {
    vm.expectRevert(abi.encodeWithSelector(CCIPBase.ZeroAddressNotAllowed.selector));

    s_receiver.updateRouter(address(0));
  }
}
