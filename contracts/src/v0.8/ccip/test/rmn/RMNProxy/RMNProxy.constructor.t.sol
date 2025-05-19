// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNProxy} from "../../../rmn/RMNProxy.sol";
import {RMNProxyTestSetup} from "./RMNProxyTestSetup.t.sol";

contract RMNProxy_constructor is RMNProxyTestSetup {
  function test_Constructor() public {
    vm.expectEmit();
    emit RMNProxy.ARMSet(MOCK_RMN_ADDRESS);
    RMNProxy proxy = new RMNProxy(MOCK_RMN_ADDRESS);
    assertEq(proxy.getARM(), MOCK_RMN_ADDRESS);
  }
}
