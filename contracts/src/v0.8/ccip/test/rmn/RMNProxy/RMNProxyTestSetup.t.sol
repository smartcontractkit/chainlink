// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {RMNProxy} from "../../../rmn/RMNProxy.sol";
import {Test} from "forge-std/Test.sol";

contract RMNProxyTestSetup is Test {
  address internal constant EMPTY_ADDRESS = address(0x1);
  address internal constant OWNER_ADDRESS = 0xC0ffeeEeC0fFeeeEc0ffeEeEc0ffEEEEC0FfEEee;
  address internal constant MOCK_RMN_ADDRESS = 0x1337133713371337133713371337133713371337;
  RMNProxy internal s_rmnProxy;

  function setUp() public virtual {
    // needed so that the extcodesize check in RMNProxy.fallback doesn't revert
    vm.etch(MOCK_RMN_ADDRESS, bytes("fake bytecode"));

    vm.prank(OWNER_ADDRESS);
    s_rmnProxy = new RMNProxy(MOCK_RMN_ADDRESS);
  }
}
