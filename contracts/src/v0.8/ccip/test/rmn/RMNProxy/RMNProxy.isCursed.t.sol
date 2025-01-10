// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRMN} from "../../../interfaces/IRMN.sol";

import {RMNProxy} from "../../../rmn/RMNProxy.sol";
import {GLOBAL_CURSE_SUBJECT, RMNRemote} from "../../../rmn/RMNRemote.sol";
import {RMNProxyTestSetup} from "./RMNProxyTestSetup.t.sol";

contract RMNProxy_isCursed is RMNProxyTestSetup {
  RMNRemote internal s_mockRMNRemote;

  function setUp() public virtual override {
    super.setUp();
    s_mockRMNRemote = new RMNRemote(1, IRMN(address(0)));
    s_rmnProxy = new RMNProxy(address(s_mockRMNRemote));
  }

  function test_IsCursed_GlobalCurseSubject() public {
    assertFalse(IRMN(address(s_rmnProxy)).isCursed());

    s_mockRMNRemote.curse(GLOBAL_CURSE_SUBJECT);
    vm.assertTrue(IRMN(address(s_rmnProxy)).isCursed());
  }

  error CustomError(bytes err);

  function test_isCursed_RevertWhen_isCursedReasonForwarded() public {
    bytes memory err = bytes("revert");
    vm.mockCallRevert(
      address(s_mockRMNRemote), abi.encodeWithSignature("isCursed()"), abi.encodeWithSelector(CustomError.selector, err)
    );

    s_rmnProxy.setARM(address(s_mockRMNRemote));
    vm.expectRevert(abi.encodeWithSelector(CustomError.selector, err));
    IRMN(address(s_rmnProxy)).isCursed();
  }

  function test_RevertWhen_call_ARMCallEmptyContract() public {
    s_rmnProxy.setARM(EMPTY_ADDRESS); // No code at address 1, should revert.
    vm.expectRevert();
    (bool success,) = address(s_rmnProxy).call(new bytes(0));
    success;
  }
}
