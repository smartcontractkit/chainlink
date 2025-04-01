// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {LOCK_RELEASE_FLAG} from "../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from
  "../USDC/HybridLockReleaseUSDCTokenPool/HybridLockReleaseUSDCTokenPoolSetup.t.sol";
import {BurnMintWithLockReleaseFlagTokenPoolSetup} from "./BurnMintWithLockReleaseFlagTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";

contract BurnMintWithLockReleaseFlagTokenPool_releaseOrMint is BurnMintWithLockReleaseFlagTokenPoolSetup {
  function test_releaseOrMint_LockReleaseFlagInSourcePoolData() public {
    uint256 amount = 1e19;
    address receiver = makeAddr("receiver_address");

    vm.startPrank(s_burnMintOffRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), receiver, amount);

    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: receiver,
        amount: amount,
        localToken: address(s_burnMintERC20),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    assertEq(s_burnMintERC20.balanceOf(receiver), amount);
  }

  function test_releaseOrMint_EmptySourcePoolData() public {
    uint256 amount = 1e19;
    address receiver = makeAddr("receiver_address");

    vm.startPrank(s_burnMintOffRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(0), receiver, amount);

    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: receiver,
        amount: amount,
        localToken: address(s_burnMintERC20),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    assertEq(s_burnMintERC20.balanceOf(receiver), amount);
  }
}

contract BurnMintWithLockReleaseFlagTokenPool_releaseOrMint_e2eTest is
  BurnMintWithLockReleaseFlagTokenPoolSetup,
  HybridLockReleaseUSDCTokenPoolSetup
{
  function setUp() public override(BurnMintWithLockReleaseFlagTokenPoolSetup, HybridLockReleaseUSDCTokenPoolSetup) {
    // Set up the BurnMintWithLockReleaseFlagTokenPool and source chain hybrid Pool
    BurnMintWithLockReleaseFlagTokenPoolSetup.setUp();
    HybridLockReleaseUSDCTokenPoolSetup.setUp();
  }

  function test_releaseOrMint_SourcePoolDataFromHybridUSDCPool() public {
    bytes memory receiver = abi.encode(STRANGER);
    uint256 amount = 1e6;

    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR),
      "Lock/Release mech not configured for outgoing message to DEST_CHAIN_SELECTOR"
    );

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    // Get the output value from the hybrid pool which will be passed to the destination pool
    Pool.LockOrBurnOutV1 memory lockOrBurnOut = s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), amount, "Incorrect token amount in the tokenPool");

    vm.startPrank(s_burnMintOffRamp);

    // Call the releaseOrMint function on the destination pool with the output value from the source pool
    s_pool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: STRANGER,
        amount: amount,
        localToken: address(s_burnMintERC20),
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_remoteBurnMintPool),
        sourcePoolData: lockOrBurnOut.destPoolData, // Use the output value from the source pool destData
        offchainTokenData: ""
      })
    );

    assertEq(s_burnMintERC20.balanceOf(STRANGER), amount);
  }
}
