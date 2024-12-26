// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {RateLimiter} from "../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../pools/TokenPool.sol";
import {BurnToAddressMintTokenPoolSetup} from "./BurnToAddressMintTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

contract BurnToAddressMintTokenPool_lockOrBurn is BurnToAddressMintTokenPoolSetup {
  function test_LockOrBurn() public {
    uint256 burnAmount = s_initialTokenAmount;

    deal(address(s_burnMintERC20), address(s_pool), burnAmount);
    assertEq(s_burnMintERC20.balanceOf(address(s_pool)), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(burnAmount);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_pool), address(0xdead), burnAmount);

    vm.expectEmit();
    emit TokenPool.Burned(address(s_burnMintOnRamp), burnAmount);

    bytes4 expectedSignature = bytes4(keccak256("transfer(address,uint256)"));
    vm.expectCall(
      address(s_burnMintERC20), abi.encodeWithSelector(expectedSignature, s_pool.getBurnAddress(), burnAmount)
    );

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
  }

  // Reverts

  function test_LockOrBurn_RevertWhen_LockedTokensUnderflows() public {
    uint256 burnAmount = s_initialTokenAmount + 1;

    deal(address(s_burnMintERC20), address(s_pool), burnAmount);

    vm.startPrank(s_burnMintOnRamp);

    // Call should revert due to underflow error due to trying to burn more tokens than are locked via CCIP.
    vm.expectRevert(abi.encodeWithSignature("Panic(uint256)", 0x11));

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
