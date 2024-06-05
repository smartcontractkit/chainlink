// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./KeystoneForwarderBaseTest.t.sol";
import {KeystoneForwarder} from "../KeystoneForwarder.sol";

contract KeystoneForwarder_SetConfigTest is BaseTest {
  address internal constant STRANGER = address(2);

  function test_RevertWhen_NotOwner() public {
    vm.stopPrank(); // BaseTest sets ADMIN
    vm.prank(STRANGER);
    vm.expectRevert();
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, _getSignerAddresses());
  }

  function test_RevertWhen_FaultToleranceIsZero() public {
    vm.expectRevert(KeystoneForwarder.FaultToleranceMustBePositive.selector);
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, 0, _getSignerAddresses());
  }

  function test_RevertWhen_InsufficientSigners() public {
    address[] memory signers = new address[](1);

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.InsufficientSigners.selector, 1, 4));
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, signers);
  }

  function test_RevertWhen_ExcessSigners() public {
    address[] memory signers = new address[](64);

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.ExcessSigners.selector, 64, 31));
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, signers);
  }

  function test_RevertWhen_ProvidingDuplicateSigners() public {
    address[] memory signers = _getSignerAddresses();
    signers[1] = signers[0];

    vm.expectRevert(abi.encodeWithSelector(KeystoneForwarder.DuplicateSigner.selector, signers[0]));
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, signers);
  }

  function test_SetConfig_FirstTime() public {
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, _getSignerAddresses());
  }

  function test_SetConfig_WhenSignersAreRemoved() public {
    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, _getSignerAddresses());

    s_forwarder.setConfig(DON_ID, CONFIG_VERSION, F, _getSignerAddresses(16));
  }
}
