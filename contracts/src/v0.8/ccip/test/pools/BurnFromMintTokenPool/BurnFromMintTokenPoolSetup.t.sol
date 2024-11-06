// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {BurnFromMintTokenPool} from "../../../pools/BurnFromMintTokenPool.sol";
import {BurnMintSetup} from "../BurnMintTokenPool/BurnMintSetup.t.sol";

contract BurnFromMintTokenPoolSetup is BurnMintSetup {
  BurnFromMintTokenPool internal s_pool;

  function setUp() public virtual override {
    BurnMintSetup.setUp();

    s_pool = new BurnFromMintTokenPool(s_burnMintERC677, new address[](0), address(s_mockRMN), address(s_sourceRouter));
    s_burnMintERC677.grantMintAndBurnRoles(address(s_pool));

    _applyChainUpdates(address(s_pool));
  }
}
