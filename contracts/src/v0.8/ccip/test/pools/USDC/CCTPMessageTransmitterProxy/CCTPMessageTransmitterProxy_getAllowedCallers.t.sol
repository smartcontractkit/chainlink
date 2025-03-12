// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCTPMessageTransmitterProxy} from "../../../../pools/USDC/CCTPMessageTransmitterProxy.sol";
import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_getAllowedCallers is CCTPMessageTransmitterProxySetup {
  function test_configureAllowedCallers() public {
    CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[] memory allowedCallerParams =
      new CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[](2);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: s_usdcTokenPool, allowed: true});
    allowedCallerParams[1] = CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: msg.sender, allowed: true});
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
    assertTrue(s_cctpMessageTransmitterProxy.isAllowedCaller(s_usdcTokenPool));

    // Get the allowed callers
    address[] memory allowedCallers = s_cctpMessageTransmitterProxy.getAllowedCallers();
    assertEq(allowedCallers.length, 2);
    assertEq(allowedCallers[0], s_usdcTokenPool);
    assertEq(allowedCallers[1], msg.sender);
  }
}
