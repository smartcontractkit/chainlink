// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {LockReleaseTokenPool} from "../../../pools/LockReleaseTokenPool.sol";
import {LockReleaseTokenPoolSetup} from "./LockReleaseTokenPoolSetup.t.sol";

contract LockReleaseTokenPool_canAcceptLiquidity is LockReleaseTokenPoolSetup {
  function test_CanAcceptLiquidity() public {
    assertEq(true, s_lockReleaseTokenPool.canAcceptLiquidity());

    s_lockReleaseTokenPool = new LockReleaseTokenPool(
      s_token, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), false, address(s_sourceRouter)
    );
    assertEq(false, s_lockReleaseTokenPool.canAcceptLiquidity());
  }
}
