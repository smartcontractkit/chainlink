// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNProxy} from "../../../rmn/RMNProxy.sol";

import {RMNProxyTestSetup} from "./RMNProxyTestSetup.t.sol";

contract RMNProxy_setARM is RMNProxyTestSetup {
  function test_SetARM() public {
    vm.expectEmit();
    emit RMNProxy.ARMSet(MOCK_RMN_ADDRESS);
    vm.prank(OWNER_ADDRESS);
    s_rmnProxy.setARM(MOCK_RMN_ADDRESS);
    assertEq(s_rmnProxy.getARM(), MOCK_RMN_ADDRESS);
  }

  function test_SetARMzero() public {
    vm.expectRevert(abi.encodeWithSelector(RMNProxy.ZeroAddressNotAllowed.selector));
    vm.prank(OWNER_ADDRESS);
    s_rmnProxy.setARM(address(0x0));
  }
}
