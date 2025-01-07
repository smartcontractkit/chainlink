// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {IERC20} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";

import {Pool} from "../../../libraries/Pool.sol";
import {SiloedLockReleaseTokenPoolSetup} from "./SiloedLockReleaseTokenPoolSetup.t.sol";

contract SiloedLockReleaseTokenPool_releaseOrMint is SiloedLockReleaseTokenPoolSetup {
  function test_ReleaseOrMint_SiloedFunds_Success() public {
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

    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), amount);

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

    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SILOED_CHAIN_SELECTOR), 0);
  }

  function test_ReleaseOrMint_UnsiloedFunds_Success() public {
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

    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SOURCE_CHAIN_SELECTOR), amount);
    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), amount);

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

    assertEq(s_siloedLockReleaseTokenPool.getSiloedTokensByChain(SOURCE_CHAIN_SELECTOR), 0);
    assertEq(s_siloedLockReleaseTokenPool.getliquidityForUnsiloedChains(), 0);
  }
}
