// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITokenMessenger} from "../../../../pools/USDC/ITokenMessenger.sol";

import {CCTPMessageTransmitterProxy} from "../../../../pools/USDC/CCTPMessageTransmitterProxy.sol";
import {BaseTest} from "../../../BaseTest.t.sol";

contract CCTPMessageTransmitterProxySetup is BaseTest {
  address internal s_tokenMessenger = makeAddr("TOKEN_MESSENGER");
  address internal s_cctpMessageTransmitter = makeAddr("CCTP_MT");
  address internal s_usdcTokenPool = makeAddr("USDC_TP");
  CCTPMessageTransmitterProxy internal s_cctpMessageTransmitterProxy;

  function setUp() public virtual override {
    super.setUp();
    vm.mockCall(
      s_tokenMessenger,
      abi.encodeWithSelector(ITokenMessenger.localMessageTransmitter.selector),
      abi.encode(s_cctpMessageTransmitter)
    );
    s_cctpMessageTransmitterProxy = new CCTPMessageTransmitterProxy(ITokenMessenger(s_tokenMessenger));
    CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[] memory allowedCallerParams =
      new CCTPMessageTransmitterProxy.AllowedCallerConfigArgs[](1);
    allowedCallerParams[0] =
      CCTPMessageTransmitterProxy.AllowedCallerConfigArgs({caller: s_usdcTokenPool, allowed: true});
    s_cctpMessageTransmitterProxy.configureAllowedCallers(allowedCallerParams);
  }
}
