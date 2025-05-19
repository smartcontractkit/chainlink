// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {ITokenMessenger} from "../../../../pools/USDC/ITokenMessenger.sol";

import {Pool} from "../../../../libraries/Pool.sol";
import {RateLimiter} from "../../../../libraries/RateLimiter.sol";
import {TokenPool} from "../../../../pools/TokenPool.sol";
import {HybridLockReleaseUSDCTokenPool} from "../../../../pools/USDC/HybridLockReleaseUSDCTokenPool.sol";
import {USDCTokenPool} from "../../../../pools/USDC/USDCTokenPool.sol";
import {HybridLockReleaseUSDCTokenPoolSetup} from "./HybridLockReleaseUSDCTokenPoolSetup.t.sol";

contract HybridLockReleaseUSDCTokenPool_lockOrBurn is HybridLockReleaseUSDCTokenPoolSetup {
  function test_onLockReleaseMechanism() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));

    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR),
      "Lock/Release mech not configured for outgoing message to DEST_CHAIN_SELECTOR"
    );

    uint256 amount = 1e6;

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    vm.expectEmit();
    emit TokenPool.Locked(s_routerAllowedOnRamp, amount);

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), amount, "Incorrect token amount in the tokenPool");
  }

  function test_PrimaryMechanism() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));
    uint256 amount = 1;

    vm.startPrank(OWNER);

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    USDCTokenPool.Domain memory expectedDomain = s_usdcTokenPool.getDomain(DEST_CHAIN_SELECTOR);

    vm.expectEmit();
    emit RateLimiter.TokensConsumed(amount);

    vm.expectEmit();
    emit ITokenMessenger.DepositForBurn(
      s_mockUSDC.s_nonce(),
      address(s_token),
      amount,
      address(s_usdcTokenPool),
      receiver,
      expectedDomain.domainIdentifier,
      s_mockUSDC.DESTINATION_TOKEN_MESSENGER(),
      expectedDomain.allowedCaller
    );

    vm.expectEmit();
    emit TokenPool.Burned(s_routerAllowedOnRamp, amount);

    Pool.LockOrBurnOutV1 memory poolReturnDataV1 = s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );

    uint64 nonce = abi.decode(poolReturnDataV1.destPoolData, (uint64));
    assertEq(s_mockUSDC.s_nonce() - 1, nonce);
  }

  function test_onLockReleaseMechanism_thenSwitchToPrimary() public {
    // Test Enabling the LR mechanism and sending an outgoing message
    test_PrimaryMechanism();

    // Disable the LR mechanism so that primary CCTP is used and then attempt to send a message
    uint64[] memory destChainRemoves = new uint64[](1);
    destChainRemoves[0] = DEST_CHAIN_SELECTOR;

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit HybridLockReleaseUSDCTokenPool.LockReleaseDisabled(DEST_CHAIN_SELECTOR);

    s_usdcTokenPool.updateChainSelectorMechanisms(destChainRemoves, new uint64[](0));

    // Send an outgoing message
    test_PrimaryMechanism();
  }

  function test_RevertWhen_WhileMigrationPause() public {
    // Mark the destination chain as supporting CCTP, so use L/R instead.
    uint64[] memory destChainAdds = new uint64[](1);
    destChainAdds[0] = DEST_CHAIN_SELECTOR;

    s_usdcTokenPool.updateChainSelectorMechanisms(new uint64[](0), destChainAdds);

    // Create a fake migration proposal
    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    assertEq(s_usdcTokenPool.getCurrentProposedCCTPChainMigration(), DEST_CHAIN_SELECTOR);

    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));

    assertTrue(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR),
      "Lock Release mech not configured for outgoing message to DEST_CHAIN_SELECTOR"
    );

    uint256 amount = 1e6;

    s_token.transfer(address(s_usdcTokenPool), amount);

    vm.startPrank(s_routerAllowedOnRamp);

    // Expect the lockOrBurn to fail because a pending CCTP-Migration has paused outgoing messages on CCIP
    vm.expectRevert(
      abi.encodeWithSelector(HybridLockReleaseUSDCTokenPool.LanePausedForCCTPMigration.selector, DEST_CHAIN_SELECTOR)
    );

    s_usdcTokenPool.lockOrBurn(
      Pool.LockOrBurnInV1({
        originalSender: OWNER,
        receiver: abi.encodePacked(receiver),
        amount: amount,
        remoteChainSelector: DEST_CHAIN_SELECTOR,
        localToken: address(s_token)
      })
    );
  }
}
