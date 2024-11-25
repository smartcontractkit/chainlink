// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../../libraries/Internal.sol";
import {Pool} from "../../../../libraries/Pool.sol";
import {TokenPool} from "../../../../pools/TokenPool.sol";
import {LOCK_RELEASE_FLAG} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {HybridLockReleaseUSDCTokenPool_releaseOrMint} from
  "../HybridLockReleaseUSDCTokenPool/HybridLockReleaseUSDCTokenPool.releaseOrMint.t.sol";

contract USDCBridgeMigrator_releaseOrMint is HybridLockReleaseUSDCTokenPool_releaseOrMint {
  function test_unstickManualTxAfterMigration_destChain() public {
    address recipient = address(1234);
    // Test the edge case where a tx is stuck in the manual tx queue and the destination chain is the one that
    // should process is after a migration. I.E the message will have the Lock-Release flag set in the OffChainData,
    // which should tell it to use the lock-release mechanism with the tokens provided.

    // We want the released amount to be 1e6, so to simulate the workflow, we sent those tokens to the contract as
    // liquidity
    uint256 amount = 1e6;
    // Add 1e12 liquidity so that there's enough to release
    vm.startPrank(s_usdcTokenPool.getLiquidityProvider(SOURCE_CHAIN_SELECTOR));

    s_token.approve(address(s_usdcTokenPool), type(uint256).max);
    s_usdcTokenPool.provideLiquidity(SOURCE_CHAIN_SELECTOR, amount);

    // By Default, the source chain will be indicated as use-CCTP so we need to change that. We create a message
    // that will use the Lock-Release flag in the offchain data to indicate that the tokens should be released
    // instead of minted since there's no attestation for us to use.

    vm.startPrank(s_routerAllowedOffRamp);

    vm.expectEmit();
    emit TokenPool.Released(s_routerAllowedOffRamp, recipient, amount);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    Pool.ReleaseOrMintOutV1 memory poolReturnDataV1 = s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    // By this point, the tx should have executed, with the Lock-Release taking over, and being forwaded to the
    // recipient

    assertEq(poolReturnDataV1.destinationAmount, amount, "destinationAmount and actual amount transferred differ");
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 0, "Tokens should be transferred out of the pool");
    assertEq(s_token.balanceOf(recipient), amount, "Tokens should be transferred to the recipient");

    // We also want to check that the system uses CCTP Burn/Mint for all other messages that don't have that flag
    // which after a migration will mean all new messages.

    // The message should fail without an error because it failed to decode a non-existent attestation which would
    // revert without an error
    vm.expectRevert();

    s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: "",
        offchainTokenData: ""
      })
    );
  }

  function test_unstickManualTxAfterMigration_homeChain() public {
    address CIRCLE = makeAddr("CIRCLE");
    address recipient = address(1234);

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = SOURCE_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    // Test the edge case where a tx is stuck in the manual tx queue and the source chain (mainnet) needs unsticking
    // In this test we want 1e6 worth of tokens to be stuck, so first we provide liquidity to the pool >1e6

    uint256 amount = 1e6;
    // Add 1e12 liquidity so that there's enough to release
    vm.startPrank(s_usdcTokenPool.getLiquidityProvider(SOURCE_CHAIN_SELECTOR));

    s_token.approve(address(s_usdcTokenPool), type(uint256).max);

    // I picked 3x the amount to be stuck so that we can have enough to release with a buffer
    s_usdcTokenPool.provideLiquidity(SOURCE_CHAIN_SELECTOR, amount * 3);

    // At this point in the process, the router will lock new messages, so we want to simulate excluding tokens
    // stuck coming back from the destination, to the home chain. This way they can be released and not minted
    // since there's no corresponding attestation to use for minting.
    vm.startPrank(OWNER);

    s_usdcTokenPool.proposeCCTPMigration(SOURCE_CHAIN_SELECTOR);

    // Exclude the tokens from being burned and check for the event
    vm.expectEmit();
    emit USDCBridgeMigrator.TokensExcludedFromBurn(SOURCE_CHAIN_SELECTOR, amount, (amount * 3) - amount);

    s_usdcTokenPool.excludeTokensFromBurn(SOURCE_CHAIN_SELECTOR, amount);

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(SOURCE_CHAIN_SELECTOR),
      (amount * 3),
      "Tokens locked minus ones excluded from the burn should be 2e6"
    );

    assertEq(
      s_usdcTokenPool.getExcludedTokensByChain(SOURCE_CHAIN_SELECTOR),
      1e6,
      "1e6 tokens should be excluded from the burn"
    );

    s_usdcTokenPool.setCircleMigratorAddress(CIRCLE);

    vm.startPrank(CIRCLE);

    s_usdcTokenPool.burnLockedUSDC();

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(SOURCE_CHAIN_SELECTOR), 0, "All tokens should be burned out of the pool"
    );

    assertEq(
      s_usdcTokenPool.getExcludedTokensByChain(SOURCE_CHAIN_SELECTOR),
      1e6,
      "There should still be 1e6 tokens excluded from the burn"
    );

    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 1e6, "All tokens minus the excluded should be in the pool");

    // Now that the burn is successful, we can release the tokens that were excluded from the burn
    vm.startPrank(s_routerAllowedOffRamp);

    vm.expectEmit();
    emit TokenPool.Released(s_routerAllowedOffRamp, recipient, amount);

    Internal.SourceTokenData memory sourceTokenData = Internal.SourceTokenData({
      sourcePoolAddress: abi.encode(SOURCE_CHAIN_USDC_POOL),
      destTokenAddress: abi.encode(address(s_usdcTokenPool)),
      extraData: abi.encode(USDCTokenPool.SourceTokenDataPayload({nonce: 1, sourceDomain: SOURCE_DOMAIN_IDENTIFIER})),
      destGasAmount: USDC_DEST_TOKEN_GAS
    });

    Pool.ReleaseOrMintOutV1 memory poolReturnDataV1 = s_usdcTokenPool.releaseOrMint(
      Pool.ReleaseOrMintInV1({
        originalSender: abi.encode(OWNER),
        receiver: recipient,
        amount: amount,
        localToken: address(s_token),
        remoteChainSelector: SOURCE_CHAIN_SELECTOR,
        sourcePoolAddress: sourceTokenData.sourcePoolAddress,
        sourcePoolData: abi.encode(LOCK_RELEASE_FLAG),
        offchainTokenData: ""
      })
    );

    assertEq(poolReturnDataV1.destinationAmount, amount, "destinationAmount and actual amount transferred differ");
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 0, "Tokens should be transferred out of the pool");
    assertEq(s_token.balanceOf(recipient), amount, "Tokens should be transferred to the recipient");
    assertEq(
      s_usdcTokenPool.getExcludedTokensByChain(SOURCE_CHAIN_SELECTOR),
      0,
      "All tokens should be released from the exclusion list"
    );

    // We also want to check that the system uses CCTP Burn/Mint for all other messages that don't have that flag
    test_incomingMessageWithPrimaryMechanism();
  }
}
