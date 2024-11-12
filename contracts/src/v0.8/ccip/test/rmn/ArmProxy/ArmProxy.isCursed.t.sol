// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../../interfaces/IRMN.sol";

import {ARMProxy} from "../../../rmn/ARMProxy.sol";
import {MockRMN} from "../../mocks/MockRMN.sol";

import {ARMProxyTestSetup} from "./ARMProxyTestSetup.t.sol";

contract ARMProxy_isCursed is ARMProxyTestSetup {
  MockRMN internal s_mockRMN;

  function setUp() public virtual override {
    super.setUp();
    s_mockRMN = new MockRMN();
    s_armProxy = new ARMProxy(address(s_mockRMN));
  }

  function test_IsCursed_Success() public {
    s_armProxy.setARM(address(s_mockRMN));
    assertFalse(IRMN(address(s_armProxy)).isCursed());
    s_mockRMN.setGlobalCursed(true);
    assertTrue(IRMN(address(s_armProxy)).isCursed());
  }

  function test_isCursed_RevertReasonForwarded_Revert() public {
    bytes memory err = bytes("revert");
    s_mockRMN.setIsCursedRevert(err);
    s_armProxy.setARM(address(s_mockRMN));
    vm.expectRevert(abi.encodeWithSelector(MockRMN.CustomError.selector, err));
    IRMN(address(s_armProxy)).isCursed();
  }

  function test_call_ARMCallEmptyContract_Revert() public {
    s_armProxy.setARM(EMPTY_ADDRESS); // No code at address 1, should revert.
    vm.expectRevert();
    bytes memory b = new bytes(0);
    (bool success,) = address(s_armProxy).call(b);
    success;
  }
}
