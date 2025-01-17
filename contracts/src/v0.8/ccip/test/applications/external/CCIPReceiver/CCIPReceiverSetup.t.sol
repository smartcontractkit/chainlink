// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {OnRampSetup} from "../../../OnRamp/OnRamp/OnRampSetup.t.sol";
import {CCIPReceiverReverting} from "../../../helpers/receivers/CCIPReceiverReverting.sol";

contract CCIPReceiverSetup is OnRampSetup {
  event MessageFailed(bytes32 indexed messageId, bytes reason);
  event MessageSucceeded(bytes32 indexed messageId);
  event MessageRecovered(bytes32 indexed messageId);

  CCIPReceiverReverting internal s_receiver;
  uint64 internal s_sourceChainSelector = 7331;

  function setUp() public virtual override {
    OnRampSetup.setUp();

    s_receiver = new CCIPReceiverReverting(address(s_destRouter));

    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: s_sourceChainSelector,
      allowed: true,
      recipient: abi.encode(address(s_receiver)),
      extraArgsBytes: ""
    });
    s_receiver.applyChainUpdates(chainUpdates);

    CCIPBase.ApprovedSenderUpdate[] memory senderUpdates = new CCIPBase.ApprovedSenderUpdate[](1);
    senderUpdates[0] =
      CCIPBase.ApprovedSenderUpdate({destChainSelector: s_sourceChainSelector, sender: abi.encode(address(s_receiver))});

    s_receiver.updateApprovedSenders(senderUpdates, new CCIPBase.ApprovedSenderUpdate[](0));
  }
}
