// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ARMProxy} from "../../../rmn/ARMProxy.sol";
import {ARMProxyTestSetup} from "./ARMProxyTestSetup.t.sol";

contract ARMProxy_constructor is ARMProxyTestSetup {
  function test_Constructor() public {
    vm.expectEmit();
    emit ARMProxy.ARMSet(MOCK_RMN_ADDRESS);
    ARMProxy proxy = new ARMProxy(MOCK_RMN_ADDRESS);
    assertEq(proxy.getARM(), MOCK_RMN_ADDRESS);
  }
}
