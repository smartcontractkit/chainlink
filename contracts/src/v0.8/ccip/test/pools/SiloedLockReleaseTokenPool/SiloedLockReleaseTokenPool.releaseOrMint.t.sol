// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

import {Pool} from "../../../libraries/Pool.sol";
import {SiloedLockReleaseTokenPool} from "../../../pools/SiloedLockReleaseTokenPool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_releaseOrMint is SiloedLockReleaseTokenPoolSetup {
  function test_ReleaseOrMint_SiloedFunds() public {
    uint256 amount = 10e18;

    deal(address(s_token), address(s_siloedLockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOnRamp);

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

  function test_ReleaseOrMint_UnsiloedFunds() public {
    uint256 amount = 10e18;

    deal(address(s_token), address(s_siloedLockReleaseTokenPool), amount);
    vm.startPrank(s_allowedOnRamp);

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

  function test_ReleaseOrMint_RevertsWhen_InsufficientLiquidity_SiloedTokenPool() public {
    uint256 releaseAmount = 10e18;
    uint256 liquidityAmount = releaseAmount - 1;

    s_siloedLockReleaseTokenPool.provideSiloedLiquidity(SILOED_CHAIN_SELECTOR, liquidityAmount);

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

  function test_ReleaseOrMint_RevertsWhen_InsufficientLiquidity_UnsiloedTokenPool() public {
    uint256 releaseAmount = 10e18;
    uint256 liquidityAmount = releaseAmount - 1;

    s_siloedLockReleaseTokenPool.provideLiquidity(liquidityAmount);

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
