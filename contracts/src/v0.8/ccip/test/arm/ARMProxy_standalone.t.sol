// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {ARMProxy} from "../../ARMProxy.sol";
import {Test} from "forge-std/Test.sol";

contract ARMProxyStandaloneTest is Test {
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

  function test_Constructor() public {
    vm.expectEmit();
    emit ARMProxy.ARMSet(MOCK_RMN_ADDRESS);
    ARMProxy proxy = new ARMProxy(MOCK_RMN_ADDRESS);
    assertEq(proxy.getARM(), MOCK_RMN_ADDRESS);
  }

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

  /*
  function test_Fuzz_ARMCall(bool expectedSuccess, bytes memory call, bytes memory ret) public {
    // filter out calls to functions that will be handled on the ARMProxy instead
    // of the underlying ARM contract
    vm.assume(
      call.length < 4 ||
        (bytes4(call) != s_armProxy.getARM.selector &&
          bytes4(call) != s_armProxy.setARM.selector &&
          bytes4(call) != s_armProxy.owner.selector &&
          bytes4(call) != s_armProxy.acceptOwnership.selector &&
          bytes4(call) != s_armProxy.transferOwnership.selector &&
          bytes4(call) != s_armProxy.typeAndVersion.selector)
    );

    if (expectedSuccess) {
      vm.mockCall(MOCK_RMN_ADDRESS, 0, call, ret);
    } else {
      vm.mockCallRevert(MOCK_RMN_ADDRESS, 0, call, ret);
    }
    (bool actualSuccess, bytes memory result) = address(s_armProxy).call(call);
    vm.clearMockedCalls();

    assertEq(result, ret);
    assertEq(expectedSuccess, actualSuccess);
  }
  */

  function test_ARMCallEmptyContractRevert() public {
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(EMPTY_ADDRESS); // No code at address 1, should revert.
    vm.expectRevert();
    bytes memory b = new bytes(0);
    (bool success,) = address(s_armProxy).call(b);
    success;
  }
}
