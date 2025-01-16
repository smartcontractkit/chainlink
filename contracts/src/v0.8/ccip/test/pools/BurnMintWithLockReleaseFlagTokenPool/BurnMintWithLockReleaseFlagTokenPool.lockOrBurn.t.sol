// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {BurnMintWithLockReleaseFlagTokenPool} from "../../../pools/USDC/BurnMintWithLockReleaseFlagTokenPool.sol";

import {LOCK_RELEASE_FLAG} from "../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {BurnMintSetup} from "../BurnMintTokenPool/BurnMintSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract BurnMintWithLockReleaseFlagTokenPoolSetup is BurnMintSetup {
  BurnMintWithLockReleaseFlagTokenPool internal s_pool;

  function setUp() public virtual override {
    BurnMintSetup.setUp();

    s_pool = new BurnMintWithLockReleaseFlagTokenPool(
      s_burnMintERC20, DEFAULT_TOKEN_DECIMALS, new address[](0), address(s_mockRMNRemote), address(s_sourceRouter)
    );
    s_burnMintERC20.grantMintAndBurnRoles(address(s_pool));

    _applyChainUpdates(address(s_pool));
  }
}

contract BurnMintWithLockReleaseFlagTokenPool_lockOrBurn is BurnMintWithLockReleaseFlagTokenPoolSetup {
  function test_LockOrBurn_CorrectReturnData() public {
    uint256 burnAmount = 20_000e18;

    deal(address(s_burnMintERC20), address(s_pool), burnAmount);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(burnAmount);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_pool), address(0), burnAmount);

    vm.expectEmit();
    emit TokenPool.Burned(address(s_burnMintOnRamp), burnAmount);

    bytes4 expectedSignature = bytes4(keccak256("burn(uint256)"));
    vm.expectCall(address(s_burnMintERC20), abi.encodeWithSelector(expectedSignature, burnAmount));

    Pool.LockOrBurnOutV1 memory lockOrBurnOut = s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: bytes(""),
        amount: burnAmount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_burnMintERC20)
      })
    );

    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), 0);

    assertEq(bytes4(lockOrBurnOut.destPoolData), LOCK_RELEASE_FLAG);
  }
}
