// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {BurnToAddressMintTokenPool} from "../../../pools/BurnToAddressMintTokenPool.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

contract BurnToAddressMintTokenPool_lockOrBurn is BurnToAddressMintTokenPoolSetup {
  function test_LockOrBurn() public {
    uint256 burnAmount = 1e24;

    s_pool.setMintedTokens(burnAmount);

    deal(address(s_burnMintERC20), address(s_pool), burnAmount);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_pool), BURN_ADDRESS, burnAmount);

    vm.expectCall(address(s_burnMintERC20), abi.encodeWithSelector(IERC20.transfer.selector, BURN_ADDRESS, burnAmount));

    s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: bytes(""),
        amount: burnAmount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_burnMintERC20)
      })
    );

    assertEq(s_burnMintERC20.balanceOf(s_pool.getBurnAddress()), burnAmount);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), 0);
    assertEq(s_pool.getMintedTokens(), 0);
  }

  // Reverts

  function test_LockOrBurn_RevertWhen_InsufficientMintedTokens() public {
    uint256 burnAmount = 1e24;

    s_pool.setMintedTokens(burnAmount - 1);

    deal(address(s_burnMintERC20), address(s_pool), burnAmount);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    vm.expectRevert(BurnToAddressMintTokenPool.InsufficientMintedTokens.selector);

    s_pool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: bytes(""),
        amount: burnAmount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_burnMintERC20)
      })
    );
  }
}
