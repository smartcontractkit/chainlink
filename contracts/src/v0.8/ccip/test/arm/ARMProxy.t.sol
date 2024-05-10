// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IRMN} from "../../interfaces/IRMN.sol";

import {ARMProxy} from "../../ARMProxy.sol";
import {RMN} from "../../RMN.sol";
import {MockRMN} from "../mocks/MockRMN.sol";
import {RMNSetup} from "./RMNSetup.t.sol";

contract ARMProxyTest is RMNSetup {
  event ARMSet(address arm);

  ARMProxy internal s_armProxy;

  function setUp() public virtual override {
    RMNSetup.setUp();
    s_armProxy = new ARMProxy(address(s_rmn));
  }

  function test_ARMIsCursed_Success() public {
    s_armProxy.setARM(address(s_mockRMN));
    assertFalse(IRMN(address(s_armProxy)).isCursed());
    RMN(address(s_armProxy)).voteToCurse(bytes32(0));
    assertTrue(IRMN(address(s_armProxy)).isCursed());
  }

  function test_ARMIsBlessed_Success() public {
    s_armProxy.setARM(address(s_mockRMN));
    assertTrue(IRMN(address(s_armProxy)).isBlessed(IRMN.TaggedRoot({commitStore: address(0), root: bytes32(0)})));
    RMN(address(s_armProxy)).voteToCurse(bytes32(0));
    assertFalse(IRMN(address(s_armProxy)).isBlessed(IRMN.TaggedRoot({commitStore: address(0), root: bytes32(0)})));
  }

  function test_ARMCallRevertReasonForwarded() public {
    bytes memory err = bytes("revert");
    s_mockRMN.setRevert(err);
    s_armProxy.setARM(address(s_mockRMN));
    vm.expectRevert(abi.encodeWithSelector(MockRMN.CustomError.selector, err));
    IRMN(address(s_armProxy)).isCursed();
  }
}
