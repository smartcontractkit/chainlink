// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";

import {OnRampSetup} from "../../../onRamp/OnRamp/OnRampSetup.t.sol";
import {CCIPReceiverReverting} from "../../../helpers/receivers/CCIPReceiverReverting.sol";

contract CCIPBase_getRemoteChainConfig is OnRampSetup {
  CCIPReceiverReverting internal s_receiver;
  uint64 internal s_sourceChainSelector = 7331;
  bytes public constant RANDOM_BYTES = "RANDOM_BYTES";
  bytes public constant RANDOM_ADDR = "RANDOM_ADDR";

  function setUp() public virtual override {
    OnRampSetup.setUp();

    s_receiver = new CCIPReceiverReverting(address(this));

    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: s_sourceChainSelector,
      allowed: true,
      recipient: RANDOM_ADDR,
      extraArgsBytes: RANDOM_BYTES
    });
    s_receiver.applyChainUpdates(chainUpdates);

    CCIPBase.ApprovedSenderUpdate[] memory senderUpdates = new CCIPBase.ApprovedSenderUpdate[](1);
    senderUpdates[0] =
      CCIPBase.ApprovedSenderUpdate({destChainSelector: s_sourceChainSelector, sender: abi.encode(address(s_receiver))});

    s_receiver.updateApprovedSenders(senderUpdates, new CCIPBase.ApprovedSenderUpdate[](0));
  }

  function test_getRemoteChainConfig() public view {
    (bytes memory recipient, bytes memory extraArgsBytes) = s_receiver.getRemoteChainConfig(s_sourceChainSelector);

    assertEq(recipient, RANDOM_ADDR);
    assertEq(extraArgsBytes, RANDOM_BYTES);
  }
}
