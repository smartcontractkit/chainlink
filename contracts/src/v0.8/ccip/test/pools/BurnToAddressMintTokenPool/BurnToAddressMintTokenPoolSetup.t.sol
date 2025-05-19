// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {BurnToAddressMintTokenPool} from "../../../pools/BurnToAddressMintTokenPool.sol";
import {BurnMintSetup} from "../BurnMintTokenPool/BurnMintSetup.t.sol";

contract BurnToAddressMintTokenPoolSetup is BurnMintSetup {
  BurnToAddressMintTokenPool internal s_pool;

  address public constant BURN_ADDRESS = address(0xdead);

  function setUp() public virtual override {
    BurnMintSetup.setUp();

    s_pool = new BurnToAddressMintTokenPool(
      s_burnMintERC20,
      DEFAULT_TOKEN_DECIMALS,
      new address[](0),
      address(s_mockRMNRemote),
      address(s_sourceRouter),
      BURN_ADDRESS
    );

    s_burnMintERC20.grantMintAndBurnRoles(address(s_pool));

    _applyChainUpdates(address(s_pool));
  }
}
