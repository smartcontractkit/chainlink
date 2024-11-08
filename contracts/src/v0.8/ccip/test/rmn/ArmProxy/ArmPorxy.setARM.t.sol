// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ARMProxy} from "../../../rmn/ARMProxy.sol";

import {ARMProxyTestSetup} from "./ARMProxyTestSetup.t.sol";

contract ARMProxy_setARM is ARMProxyTestSetup {
  function test_SetARM() public {
    vm.expectEmit();
    emit ARMProxy.ARMSet(MOCK_RMN_ADDRESS);
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(MOCK_RMN_ADDRESS);
    assertEq(s_armProxy.getARM(), MOCK_RMN_ADDRESS);
  }

  function test_SetARMzero() public {
    vm.expectRevert(abi.encodeWithSelector(ARMProxy.ZeroAddressNotAllowed.selector));
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(address(0x0));
  }
}
