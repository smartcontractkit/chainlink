// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IRouterClient} from "../../../interfaces/IRouterClient.sol";
import {Client} from "../../../libraries/Client.sol";
import {OnRampSetup} from "../../onRamp/OnRamp/OnRampSetup.t.sol";

contract Router_getFee is OnRampSetup {
  function test_GetFeeSupportedChain() public view {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();
    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_SELECTOR, message);
    assertGt(expectedFee, 10e9);
  }

  // Reverts
  function test_RevertWhen_UnsupportedDestinationChain() public {
    Client.EVM2AnyMessage memory message = _generateEmptyMessage();

    vm.expectRevert(abi.encodeWithSelector(IRouterClient.UnsupportedDestinationChain.selector, 999));
    s_sourceRouter.getFee(999, message);
  }
}
