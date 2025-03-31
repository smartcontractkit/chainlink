pragma solidity ^0.8.24;

import {BurnMintERC20} from "../../../../shared/token/ERC20/BurnMintERC20.sol";
import {BurnMintWithLockReleaseFlagTokenPool} from "../../../pools/USDC/BurnMintWithLockReleaseFlagTokenPool.sol";
import {BurnMintSetup} from "../BurnMintTokenPool/BurnMintSetup.t.sol";

contract BurnMintWithLockReleaseFlagTokenPoolSetup is BurnMintSetup {
  BurnMintWithLockReleaseFlagTokenPool internal s_pool;

  function setUp() public virtual override {
    BurnMintSetup.setUp();

    // To simulate USDC we need to override the decimals to 6
    s_burnMintERC20 = new BurnMintERC20("Chainlink Token", "LINK", 6, 0, 0);

    s_pool = new BurnMintWithLockReleaseFlagTokenPool(
      s_burnMintERC20, 6, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );
    s_burnMintERC20.grantMintAndBurnRoles(address(s_pool));

    _applyChainUpdates(address(s_pool));
  }
}
