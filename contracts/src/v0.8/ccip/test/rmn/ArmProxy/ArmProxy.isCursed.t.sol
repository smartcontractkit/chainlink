// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {IRMN} from "../../../interfaces/IRMN.sol";

import {ARMProxy} from "../../../rmn/ARMProxy.sol";
import {GLOBAL_CURSE_SUBJECT, RMNRemote} from "../../../rmn/RMNRemote.sol";
import {ARMProxyTestSetup} from "./ARMProxyTestSetup.t.sol";

contract ARMProxy_isCursed is ARMProxyTestSetup {
  RMNRemote internal s_mockRMNRemote;

  function setUp() public virtual override {
    super.setUp();
    s_mockRMNRemote = new RMNRemote(1, IRMN(address(0)));
    s_armProxy = new ARMProxy(address(s_mockRMNRemote));
  }

  function test_IsCursed_GlobalCurseSubject() public {
    assertFalse(IRMN(address(s_armProxy)).isCursed());

    s_mockRMNRemote.curse(GLOBAL_CURSE_SUBJECT);
    vm.assertTrue(IRMN(address(s_armProxy)).isCursed());
  }

  error CustomError(bytes err);

  function test_isCursed_RevertWhen_isCursedReasonForwarded() public {
    bytes memory err = bytes("revert");
    vm.mockCallRevert(
      address(s_mockRMNRemote), abi.encodeWithSignature("isCursed()"), abi.encodeWithSelector(CustomError.selector, err)
    );

    s_armProxy.setARM(address(s_mockRMNRemote));
    vm.expectRevert(abi.encodeWithSelector(CustomError.selector, err));
    IRMN(address(s_armProxy)).isCursed();
  }

  function test_RevertWhen_call_ARMCallEmptyContract() public {
    s_armProxy.setARM(EMPTY_ADDRESS); // No code at address 1, should revert.
    vm.expectRevert();
    (bool success,) = address(s_armProxy).call(new bytes(0));
    success;
  }
}
