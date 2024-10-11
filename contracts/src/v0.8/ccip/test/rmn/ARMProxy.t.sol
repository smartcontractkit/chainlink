// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../interfaces/IRMN.sol";

import {ARMProxy} from "../../rmn/ARMProxy.sol";
import {MockRMN} from "../mocks/MockRMN.sol";
import {Test} from "forge-std/Test.sol";

contract ARMProxyTest is Test {
  MockRMN internal s_mockRMN;
  ARMProxy internal s_armProxy;

  function setUp() public virtual {
    s_mockRMN = new MockRMN();
    s_armProxy = new ARMProxy(address(s_mockRMN));
  }

  function test_ARMIsCursed_Success() public {
    s_armProxy.setARM(address(s_mockRMN));
    assertFalse(IRMN(address(s_armProxy)).isCursed());
    s_mockRMN.setGlobalCursed(true);
    assertTrue(IRMN(address(s_armProxy)).isCursed());
  }

  function test_ARMCallRevertReasonForwarded() public {
    bytes memory err = bytes("revert");
    s_mockRMN.setIsCursedRevert(err);
    s_armProxy.setARM(address(s_mockRMN));
    vm.expectRevert(abi.encodeWithSelector(MockRMN.CustomError.selector, err));
    IRMN(address(s_armProxy)).isCursed();
  }
}

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

  function test_ARMCallEmptyContractRevert() public {
    vm.prank(OWNER_ADDRESS);
    s_armProxy.setARM(EMPTY_ADDRESS); // No code at address 1, should revert.
    vm.expectRevert();
    bytes memory b = new bytes(0);
    (bool success,) = address(s_armProxy).call(b);
    success;
  }
}
