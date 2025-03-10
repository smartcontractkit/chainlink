// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPReceiverLegacy} from "../../../../applications/internal/CCIPReceiverLegacy.sol";
import {EtherSenderReceiver} from "../../../../applications/internal/EtherSenderReceiver.sol";
import {EtherSenderReceiverTestSetup} from "./EtherSenderReceiverTestSetup.t.sol";

contract EtherSenderReceiverTest_constructor is EtherSenderReceiverTestSetup {
  function test_constructor() public view {
    assertEq(s_etherSenderReceiver.getRouter(), ROUTER, "router must be set correctly");
    uint256 allowance = s_weth.allowance(address(s_etherSenderReceiver), ROUTER);
    assertEq(allowance, type(uint256).max, "allowance must be set infinite");
  }

  function test_RevertWhen_InvalidRouter() public {
    vm.expectRevert(abi.encodeWithSelector(CCIPReceiverLegacy.InvalidRouter.selector, address(0)));

    new EtherSenderReceiver(address(0));
  }
}
