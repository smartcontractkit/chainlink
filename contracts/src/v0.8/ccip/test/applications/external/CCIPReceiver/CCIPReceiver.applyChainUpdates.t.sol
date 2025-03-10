// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIPBase} from "../../../../applications/external/CCIPBase.sol";
import {CCIPReceiverSetup} from "./CCIPReceiverSetup.t.sol";

contract CCIPReceiver_ApplyChainUpdates is CCIPReceiverSetup {
  function test_applyChainUpdate() public {
    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] = CCIPBase.ChainUpdate({
      chainSelector: s_sourceChainSelector,
      allowed: true,
      recipient: "RECEIVER",
      extraArgsBytes: ""
    });

    // Revert because the recipient of an allowed chain is the zero address, which is prohibited
    s_receiver.applyChainUpdates(chainUpdates);
  }

  function test_applyChainUpdates_RevertWhen_ZeroAddressNotAllowed() public {
    CCIPBase.ChainUpdate[] memory chainUpdates = new CCIPBase.ChainUpdate[](1);
    chainUpdates[0] =
      CCIPBase.ChainUpdate({chainSelector: s_sourceChainSelector, allowed: true, recipient: "", extraArgsBytes: ""});

    // Revert because the recipient of an allowed chain is the zero address, which is prohibited
    vm.expectRevert(abi.encodeWithSelector(CCIPBase.ZeroAddressNotAllowed.selector));
    s_receiver.applyChainUpdates(chainUpdates);
  }
}
