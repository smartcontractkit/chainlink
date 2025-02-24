// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Ownable2Step} from "../../../../../shared/access/Ownable2Step.sol";
import {CCTPMessageTransmitterProxy} from "../../../../pools/USDC/CCTPMessageTransmitterProxy.sol";
import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_configureAllowedCallers is CCTPMessageTransmitterProxySetup {
  function test_configureAllowedCallers() public {
    CCTPMessageTransmitterProxy.AllowedCallerConfigParam[] memory allowedCallerParams =
      new CCTPMessageTransmitterProxy.AllowedCallerConfigParam[](1);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigParam({caller: s_usdcTokenPool, allowed: true});
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
    assertTrue(s_cctpMessageTransmitterProxy.isAllowedCaller(s_usdcTokenPool));

    // Remove the allowed caller
    allowedCallerParams[0].allowed = false;
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
    assertFalse(s_cctpMessageTransmitterProxy.isAllowedCaller(s_usdcTokenPool));
  }

  // Revert cases
  function test_configureAllowedCallers_RevertWhen_NotOwner() public {
    changePrank(makeAddr("RANDOM"));
    CCTPMessageTransmitterProxy.AllowedCallerConfigParam[] memory allowedCallerParams =
      new CCTPMessageTransmitterProxy.AllowedCallerConfigParam[](1);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigParam({caller: s_usdcTokenPool, allowed: true});
    vm.expectRevert(Ownable2Step.OnlyCallableByOwner.selector);
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
  }
}
