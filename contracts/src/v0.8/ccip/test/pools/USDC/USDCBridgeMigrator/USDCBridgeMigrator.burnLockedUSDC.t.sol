// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Pool} from "../../../../libraries/Pool.sol";
import {TokenPool} from "../../../../pools/TokenPool.sol";
import {USDCBridgeMigrator} from "../../../../pools/USDC/USDCBridgeMigrator.sol";
import {HybridLockReleaseUSDCTokenPool_lockOrBurn} from
  "../HybridLockReleaseUSDCTokenPool/HybridLockReleaseUSDCTokenPool.lockOrBurn.t.sol";

contract USDCBridgeMigrator_BurnLockedUSDC is HybridLockReleaseUSDCTokenPool_lockOrBurn {
  function test_lockOrBurn_then_BurnInCCTPMigration() public {
    bytes32 receiver = bytes32(uint256(uint160(STRANGER)));
    address CIRCLE = makeAddr("CIRCLE CCTP Migrator");

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

    // Ensure that the tokens are properly locked
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), amount, "Incorrect token amount in the tokenPool");

    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      amount,
      "Internal locked token accounting is incorrect"
    );

    vm.startPrank(OWNER);

    vm.expectEmit();
    emit USDCBridgeMigrator.CircleMigratorAddressSet(CIRCLE);

    s_usdcTokenPool.setCircleMigratorAddress(CIRCLE);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationProposed(DEST_CHAIN_SELECTOR);

    // Propose the migration to CCTP
    s_usdcTokenPool.proposeCCTPMigration(DEST_CHAIN_SELECTOR);

    assertEq(
      s_usdcTokenPool.getCurrentProposedCCTPChainMigration(),
      DEST_CHAIN_SELECTOR,
      "Current proposed chain migration does not match expected for DEST_CHAIN_SELECTOR"
    );

    // Impersonate the set circle address and execute the proposal
    vm.startPrank(CIRCLE);

    vm.expectEmit();
    emit USDCBridgeMigrator.CCTPMigrationExecuted(DEST_CHAIN_SELECTOR, amount);

    // Ensure the call to the burn function is properly
    vm.expectCall(address(s_token), abi.encodeWithSelector(bytes4(keccak256("burn(uint256)")), amount));

    s_usdcTokenPool.burnLockedUSDC();

    // Assert that the tokens were actually burned
    assertEq(s_token.balanceOf(address(s_usdcTokenPool)), 0, "Tokens were not burned out of the tokenPool");

    // Ensure the proposal slot was cleared and there's no tokens locked for the destination chain anymore
    assertEq(s_usdcTokenPool.getCurrentProposedCCTPChainMigration(), 0, "Proposal Slot should be empty");
    assertEq(
      s_usdcTokenPool.getLockedTokensForChain(DEST_CHAIN_SELECTOR),
      0,
      "No tokens should be locked for DEST_CHAIN_SELECTOR after CCTP-approved burn"
    );

    assertFalse(
      s_usdcTokenPool.shouldUseLockRelease(DEST_CHAIN_SELECTOR), "Lock/Release mech should be disabled after a burn"
    );

    test_PrimaryMechanism();
  }

  function test_RevertWhen_invalidPermissions() public {
    address CIRCLE = makeAddr("CIRCLE");

    vm.startPrank(OWNER);

    // Set the circle migrator address for later, but don't start pranking as it yet
    s_usdcTokenPool.setCircleMigratorAddress(CIRCLE);

    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.onlyCircle.selector));

    // Should fail because only Circle can call this function
    s_usdcTokenPool.burnLockedUSDC();

    vm.startPrank(CIRCLE);

    vm.expectRevert(abi.encodeWithSelector(USDCBridgeMigrator.NoMigrationProposalPending.selector));
    s_usdcTokenPool.burnLockedUSDC();
  }
}
