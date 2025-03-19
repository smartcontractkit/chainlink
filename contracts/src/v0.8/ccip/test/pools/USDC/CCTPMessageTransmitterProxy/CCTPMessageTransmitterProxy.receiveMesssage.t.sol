// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IMessageTransmitter} from "../../../../pools/USDC/IMessageTransmitter.sol";

import {CCTPMessageTransmitterProxy} from "../../../../pools/USDC/CCTPMessageTransmitterProxy.sol";
import {CCTPMessageTransmitterProxySetup} from "./CCTPMessageTransmitterProxySetup.t.sol";

contract CCTPMessageTransmitterProxy_receiveMesssage is CCTPMessageTransmitterProxySetup {
  function test_receiveMesssage() public {
    bytes memory message = bytes("message");
    bytes memory attestation = bytes("attestation");

    // Mocking the call to the IMessageTransmitter to return true
    vm.mockCall(
      s_cctpMessageTransmitter,
      abi.encodeWithSelector(IMessageTransmitter.receiveMessage.selector, message, attestation),
      abi.encode(true)
    );

    changePrank(s_usdcTokenPool);
    assertTrue(s_cctpMessageTransmitterProxy.receiveMessage(message, attestation));

    // Mocking the call to the IMessageTransmitter to return false
    vm.mockCall(
      s_cctpMessageTransmitter,
      abi.encodeWithSelector(IMessageTransmitter.receiveMessage.selector, message, attestation),
      abi.encode(false)
    );

    changePrank(s_usdcTokenPool);
    assertFalse(s_cctpMessageTransmitterProxy.receiveMessage(message, attestation));
  }

  // Revert cases
  function test_receiveMessage_RevertWhen_UnAuthorizedCaller() public {
    bytes memory message = bytes("message");
    bytes memory attestation = bytes("attestation");

    changePrank(makeAddr("RANDOM"));
    vm.expectRevert(abi.encodeWithSelector(CCTPMessageTransmitterProxy.Unauthorized.selector, makeAddr("RANDOM")));
    s_cctpMessageTransmitterProxy.receiveMessage(message, attestation);
  }
}
