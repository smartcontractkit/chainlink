// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {Test} from "forge-std/Test.sol";
import {ARMProxy} from "../../ARMProxy.sol";

contract ARMProxyStandaloneTest is Test {
  event ARMSet(address arm);

  address internal constant EMPTY_ADDRESS = address(0x1);
  address internal constant OWNER_ADDRESS = 0xC0ffeeEeC0fFeeeEc0ffeEeEc0ffEEEEC0FfEEee;
  address internal constant MOCK_ARM_ADDRESS = 0x1337133713371337133713371337133713371337;

  ARMProxy internal s_armProxy;

  function setUp() public virtual {
    // needed so that the extcodesize check in ARMProxy.fallback doesn't revert
    vm.etch(MOCK_ARM_ADDRESS, bytes("fake bytecode"));

    vm.prank(OWNER_ADDRESS);
    s_armProxy = new ARMProxy(MOCK_ARM_ADDRESS);
  }

  function testConstructor() public {
    vm.expectEmit();
    emit ARMSet(MOCK_ARM_ADDRESS);
    ARMProxy proxy = new ARMProxy(MOCK_ARM_ADDRESS);
    assertEq(proxy.getARM(), MOCK_ARM_ADDRESS);
  }

  function testSetARM() public {
    vm.expectEmit();
    emit ARMSet(MOCK_ARM_ADDRESS);
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(MOCK_ARM_ADDRESS);
    assertEq(s_armProxy.getARM(), MOCK_ARM_ADDRESS);
  }

  function testSetARMzero() public {
    vm.expectRevert(abi.encodeWithSelector(ARMProxy.ZeroAddressNotAllowed.selector));
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(address(0x0));
  }

  /*
  function testARMCall_fuzz(bool expectedSuccess, bytes memory call, bytes memory ret) public {
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
      vm.mockCall(MOCK_ARM_ADDRESS, 0, call, ret);
    } else {
      vm.mockCallRevert(MOCK_ARM_ADDRESS, 0, call, ret);
    }
    (bool actualSuccess, bytes memory result) = address(s_armProxy).call(call);
    vm.clearMockedCalls();

    assertEq(result, ret);
    assertEq(expectedSuccess, actualSuccess);
  }
  */

  function testARMCallEmptyContractRevert() public {
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(EMPTY_ADDRESS); // No code at address 1, should revert.
    vm.expectRevert();
    bytes memory b = new bytes(0);
    (bool success, ) = address(s_armProxy).call(b);
    success;
  }
}
