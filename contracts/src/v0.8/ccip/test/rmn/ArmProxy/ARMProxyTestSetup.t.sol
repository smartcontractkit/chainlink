// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ARMProxy} from "../../../rmn/ARMProxy.sol";
import {Test} from "forge-std/Test.sol";

contract ARMProxyTestSetup is Test {
  address internal constant EMPTY_ADDRESS = address(0x1);
  address internal constant OWNER_ADDRESS = 0xC0ffeeEeC0fFeeeEc0ffeEeEc0ffEEEEC0FfEEee;
  address internal constant MOCK_RMN_ADDRESS = 0x1337133713371337133713371337133713371337;
  ARMProxy internal s_armProxy;

  function setUp() public virtual {
    // needed so that the extcodesize check in ARMProxy.fallback doesn't revert
    vm.etch(MOCK_RMN_ADDRESS, bytes("fake bytecode"));

    vm.prank(OWNER_ADDRESS);
    s_armProxy = new ARMProxy(MOCK_RMN_ADDRESS);
  }
}
