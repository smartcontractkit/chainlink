// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {TokenPool} from "../../../pools/TokenPool.sol";
import {TokenPoolWithAllowListSetup} from "./TokenPoolWithAllowListSetup.t.sol";

contract TokenPoolWithAllowList_setRouter is TokenPoolWithAllowListSetup {
  function test_SetRouter_Success() public {
    assertEq(address(s_sourceRouter), s_tokenPool.getRouter());

    address newRouter = makeAddr("newRouter");

    vm.expectEmit();
    emit TokenPool.RouterUpdated(address(s_sourceRouter), newRouter);

    s_tokenPool.setRouter(newRouter);

    assertEq(newRouter, s_tokenPool.getRouter());
  }
}
