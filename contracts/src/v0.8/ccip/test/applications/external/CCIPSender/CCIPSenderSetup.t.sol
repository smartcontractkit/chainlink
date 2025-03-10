// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IRouterClient} from "../../../../interfaces/IRouterClient.sol";

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPSender} from "../../../../applications/external/CCIPSender.sol";
import {OnRampSetup} from "../../../onRamp/OnRamp/OnRampSetup.t.sol";

contract CCIPSenderSetup is OnRampSetup {
  bytes32 public constant MESSAGE_ID = keccak256("FAKE_MESSAGE_ID");

  CCIPSender internal s_sender;

  function setUp() public virtual override {
    OnRampSetup.setUp();

    vm.mockCall(
      address(s_sourceRouter), abi.encodeWithSelector(IRouterClient.ccipSend.selector), abi.encode(MESSAGE_ID)
    );

    s_sender = new CCIPSender(address(s_sourceRouter));

    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: DEST_CHAIN_SELECTOR,
      allowed: true,
      recipient: abi.encode(address(s_sender)),
      extraArgsBytes: ""
    });
    s_sender.applyChainUpdates(chainUpdates);
  }
}
