// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {BurnToAddressMintTokenPool} from "../../../pools/BurnToAddressMintTokenPool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

contract BurnToAddressMintTokenPool_setOutstandingokens is BurnToAddressMintTokenPoolSetup {
  function test_setOutstandingTokens() public {
    uint256 amount = 1e18;

    assertEq(s_pool.getOutstandingTokens(), 0);

    vm.expectEmit();
    emit BurnToAddressMintTokenPool.OutstandingTokensSet(amount, 0);

    s_pool.setOutstandingTokens(amount);

    assertEq(s_pool.getOutstandingTokens(), amount);
  }
}
