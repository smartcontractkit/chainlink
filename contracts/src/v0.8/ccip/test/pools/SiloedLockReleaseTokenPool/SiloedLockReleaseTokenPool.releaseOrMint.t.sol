// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../libraries/Pool.sol";
import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

contract SiloedLockReleaseTokenPool_releaseOrMint is SiloedLockReleaseTokenPoolSetup {
  function test_ReleaseOrMint_SiloedChain() public {
    uint256 amount = 10e18;

    deal(address(s_token), address(s_siloedLockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOnRamp);

    // Lock funds so that they can be released without underflowing the internal accounting
    s_siloedLockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: SILOED_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), amount);

    vm.startPrank(s_allowedOffRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_siloedLockReleaseTokenPool), OWNER, amount);

    s_siloedLockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SILOED_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_siloedDestPoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SILOED_CHAIN_SELECTOR), 0);
  }

  function test_ReleaseOrMint_UnsiloedChain() public {
    uint256 amount = 10e18;

    deal(address(s_token), address(s_siloedLockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOnRamp);

    // Lock funds for unsiloed chain so they can be released later
    s_siloedLockReleaseTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: STRANGER,
        receiver: bytes(""),
        amount: amount,
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SOURCE_CHAIN_SELECTOR), amount);
    assertEq(s_siloedLockReleaseTokenPool.getUnsiloedLiquidity(), amount);

    vm.startPrank(s_allowedOffRamp);

    vm.expectEmit();
    emit IERC20.Transfer(address(s_siloedLockReleaseTokenPool), OWNER, amount);

    s_siloedLockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_siloedDestPoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );

    assertEq(s_siloedLockReleaseTokenPool.getAvailableTokens(SOURCE_CHAIN_SELECTOR), 0);
    assertEq(s_siloedLockReleaseTokenPool.getUnsiloedLiquidity(), 0);
  }

  // Reverts

  function test_ReleaseOrMint_RevertsWhen_InsufficientLiquidity_SiloedChain() public {
    uint256 releaseAmount = 10e18;
    uint256 liquidityAmount = releaseAmount - 1;

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, liquidityAmount);

    // Since amount to release is greater than provided liquidity, the function should revert
    vm.expectRevert(
      abi.encodeWithSelector(SiloedLockReleaseTokenPool.InsufficientLiquidity.selector, liquidityAmount, releaseAmount)
    );

    vm.startPrank(s_allowedOffRamp);

    s_siloedLockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: releaseAmount,
        localToken: address(s_token),
        remoteChainSelector: SILOED_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_siloedDestPoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_ReleaseOrMint_RevertsWhen_InsufficientLiquidity_UnsiloedChain() public {
    uint256 releaseAmount = 10e18;
    uint256 liquidityAmount = releaseAmount - 1;

    // Call the provide liquidity function which provides to unsiloed chains.
    s_siloedLockReleaseTokenPool.provideLiquidity(liquidityAmount);

    // Since amount to release is greater than provided liquidity, the function should revert
    vm.expectRevert(
      abi.encodeWithSelector(SiloedLockReleaseTokenPool.InsufficientLiquidity.selector, liquidityAmount, releaseAmount)
    );

    vm.startPrank(s_allowedOffRamp);

    s_siloedLockReleaseTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: bytes(""),
        receiver: OWNER,
        amount: releaseAmount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: abi.encode(s_siloedDestPoolAddress),
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }
}
