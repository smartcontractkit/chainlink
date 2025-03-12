// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_getCCTPTransmitter is CCTPMessageTransmitterProxySetup {
  function test_getCCTPTransmitter() public {
    assertEq(address(s_cctpMessageTransmitterProxy.i_cctpTransmitter()), address(s_cctpMessageTransmitter));
  }
}
