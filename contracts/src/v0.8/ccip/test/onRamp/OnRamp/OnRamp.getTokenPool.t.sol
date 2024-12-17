// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {OnRampSetup} from "./OnRampSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract OnRamp_getTokenPool is OnRampSetup {
  function test_GetTokenPool() public view {
    assertEq(
      s_sourcePoolByToken[s_sourceTokens[0]],
      address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_sourceTokens[0])))
    );
    assertEq(
      s_sourcePoolByToken[s_sourceTokens[1]],
      address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(s_sourceTokens[1])))
    );

    address wrongToken = address(123);
    address nonExistentPool = address(s_onRamp.getPoolBySourceToken(DEST_CHAIN_SELECTOR, IERC20(wrongToken)));

    assertEq(address(0), nonExistentPool);
  }
}
