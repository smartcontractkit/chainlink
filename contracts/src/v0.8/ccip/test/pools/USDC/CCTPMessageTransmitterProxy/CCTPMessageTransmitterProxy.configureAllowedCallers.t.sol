// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../../shared/access/Ownable2Step.sol";
import {CCTPMessageTransmitterProxy} from "../../../../pools/USDC/CCTPMessageTransmitterProxy.sol";
import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_configureAllowedCallers is CCTPMessageTransmitterProxySetup {
  function test_configureAllowedCallers() public {
    CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[] memory allowedCallerParams =
      new CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[](2);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: s_usdcTokenPool, allowed: true});
    allowedCallerParams[1] = CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: msg.sender, allowed: true});
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
    assertTrue(s_cctpMessageTransmitterProxy.isAllowedCaller(s_usdcTokenPool));

    // Remove the allowed caller
    allowedCallerParams = new CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[](1);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: s_usdcTokenPool, allowed: false});
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
    assertFalse(s_cctpMessageTransmitterProxy.isAllowedCaller(s_usdcTokenPool));

    address[] memory allowedCallers = s_cctpMessageTransmitterProxy.getAllowedCallers();
    assertEq(allowedCallers.length, 1);
    assertEq(allowedCallers[0], msg.sender);
  }

  // Revert cases
  function test_configureAllowedCallers_RevertWhen_NotOwner() public {
    changePrank(makeAddr("RANDOM"));
    CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[] memory allowedCallerParams =
      new CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[](1);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: s_usdcTokenPool, allowed: true});
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
  }
}
