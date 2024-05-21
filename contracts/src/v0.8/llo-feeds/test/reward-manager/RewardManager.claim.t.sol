// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseRewardManagerTest} from "./BaseRewardManager.t.sol";
import {Common} from "../../libraries/Common.sol";

/**
 * @title BaseRewardManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the claim functionality of the RewardManager contract.
 */
contract RewardManagerClaimTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function test_claimAllRecipients() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }
  }

  function test_claimRewardsWithDuplicatePoolIdsDoesNotPayoutTwice() public {
    //add funds to a different pool to ensure they're not claimed
    addFundsToPool(SECONDARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create an array containing duplicate poolIds
    bytes32[] memory poolIds = new bytes32[](2);
    poolIds[0] = PRIMARY_POOL_ID;
    poolIds[1] = PRIMARY_POOL_ID;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(poolIds, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //the pool should still have the remaining
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_claimSingleRecipient() public {
    //get the recipient that is claiming
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //claim the individual rewards for this recipient
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT - expectedRecipientAmount);
  }

  function test_claimMultipleRecipients() public {
    //claim the individual rewards for each recipient
    claimRewards(PRIMARY_POOL_ARRAY, getPrimaryRecipients()[0].addr);
    claimRewards(PRIMARY_POOL_ARRAY, getPrimaryRecipients()[1].addr);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(getPrimaryRecipients()[0].addr), expectedRecipientAmount);
    assertEq(getAssetBalance(getPrimaryRecipients()[1].addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT - (expectedRecipientAmount * 2));
  }

  function test_claimUnregisteredRecipient() public {
    //claim the rewards for a recipient who isn't in this pool
    claimRewards(PRIMARY_POOL_ARRAY, getSecondaryRecipients()[1].addr);

    //check the recipients didn't receive any fees from this pool
    assertEq(getAssetBalance(getSecondaryRecipients()[1].addr), 0);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_claimUnevenAmountRoundsDown() public {
    //adding 1 to the pool should leave 1 wei worth of dust, which the contract doesn't handle due to it being economically infeasible
    addFundsToPool(PRIMARY_POOL_ID, getAsset(1), FEE_MANAGER);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //check the rewardManager has the remaining quantity equals 1 wei
    assertEq(getAssetBalance(address(rewardManager)), 1);
  }

  function test_claimUnregisteredPoolId() public {
    //get the recipient that is claiming
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //claim the individual rewards for this recipient
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //check the recipients balance is still 0 as there's no pool to receive fees from
    assertEq(getAssetBalance(recipient.addr), 0);

    //check the rewardManager has the full amount
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_singleRecipientClaimMultipleDeposits() public {
    //get the recipient that is claiming
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //claim the individual rewards for this recipient
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity, which is 3/4 of the initial deposit
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT - expectedRecipientAmount);

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //claim the individual rewards for this recipient
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

    //check the recipients balance matches the ratio the recipient should have received, which is 1/4 of each deposit
    assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount * 2);

    //check the rewardManager has the remaining quantity, which is now 3/4 of both deposits
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 2 - (expectedRecipientAmount * 2));
  }

  function test_recipientsClaimMultipleDeposits() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //the reward manager balance should be 0 as all of the funds have been claimed
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //expected recipient amount is 1/4 of the pool deposit
      expectedRecipientAmount = (POOL_DEPOSIT_AMOUNT / 4) * 2;

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //the reward manager balance should again be 0 as all of the funds have been claimed
    assertEq(getAssetBalance(address(rewardManager)), 0);
  }

  function test_eventIsEmittedUponClaim() public {
    //get the recipient that is claiming
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //expect an emit
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit RewardsClaimed(PRIMARY_POOL_ID, recipient.addr, uint192(POOL_DEPOSIT_AMOUNT / 4));

    //claim the individual rewards for each recipient
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);
  }

  function test_eventIsNotEmittedUponUnsuccessfulClaim() public {
    //record logs to check no events were emitted
    vm.recordLogs();

    //get the recipient that is claiming
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //claim the individual rewards for each recipient
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //no logs should have been emitted
    assertEq(vm.getRecordedLogs().length, 0);
  }
}

contract RewardManagerRecipientClaimMultiplePoolsTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a two pools
    createPrimaryPool();
    createSecondaryPool();

    //add funds to each of the pools to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
    addFundsToPool(SECONDARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function test_claimAllRecipientsSinglePool() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //check the pool balance is still equal to DEPOSIT_AMOUNT as the test only claims for one of the pools
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_claimMultipleRecipientsSinglePool() public {
    //claim the individual rewards for each recipient
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[1].addr);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(getSecondaryRecipients()[0].addr), expectedRecipientAmount);
    assertEq(getAssetBalance(getSecondaryRecipients()[1].addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 2 - (expectedRecipientAmount * 2));
  }

  function test_claimMultipleRecipientsMultiplePools() public {
    //claim the individual rewards for each recipient
    claimRewards(PRIMARY_POOL_ARRAY, getPrimaryRecipients()[0].addr);
    claimRewards(PRIMARY_POOL_ARRAY, getPrimaryRecipients()[1].addr);
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[1].addr);

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received. The first recipient is shared across both pools so should receive 1/4 of each pool
    assertEq(getAssetBalance(getPrimaryRecipients()[0].addr), expectedRecipientAmount * 2);
    assertEq(getAssetBalance(getPrimaryRecipients()[1].addr), expectedRecipientAmount);
    assertEq(getAssetBalance(getSecondaryRecipients()[1].addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_claimAllRecipientsMultiplePools() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i = 1; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //claim funds for each recipient within the pool
    for (uint256 i = 1; i < getSecondaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory secondaryRecipient = getSecondaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(SECONDARY_POOL_ARRAY, secondaryRecipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(secondaryRecipient.addr), expectedRecipientAmount);
    }

    //special case to handle the first recipient of each pool as they're the same address
    Common.AddressAndWeight memory commonRecipient = getPrimaryRecipients()[0];

    //claim the individual rewards for each pool
    claimRewards(PRIMARY_POOL_ARRAY, commonRecipient.addr);
    claimRewards(SECONDARY_POOL_ARRAY, commonRecipient.addr);

    //check the balance matches the ratio the recipient should have received, which is 1/4 of each deposit for each pool
    assertEq(getAssetBalance(commonRecipient.addr), expectedRecipientAmount * 2);
  }

  function test_claimSingleUniqueRecipient() public {
    //the first recipient of the secondary pool is in both pools, so take the second recipient which is unique
    Common.AddressAndWeight memory recipient = getSecondaryRecipients()[1];

    //claim the individual rewards for this recipient
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //the recipient should have received 1/4 of the deposit amount
    uint256 recipientExpectedAmount = POOL_DEPOSIT_AMOUNT / 4;

    //the recipient should have received 1/4 of the deposit amount
    assertEq(getAssetBalance(recipient.addr), recipientExpectedAmount);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 2 - recipientExpectedAmount);
  }

  function test_claimSingleRecipientMultiplePools() public {
    //the first recipient of the secondary pool is in both pools
    Common.AddressAndWeight memory recipient = getSecondaryRecipients()[0];

    //claim the individual rewards for this recipient
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //the recipient should have received 1/4 of the deposit amount for each pool
    uint256 recipientExpectedAmount = (POOL_DEPOSIT_AMOUNT / 4) * 2;

    //this recipient belongs in both pools so should have received 1/4 of each
    assertEq(getAssetBalance(recipient.addr), recipientExpectedAmount);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 2 - recipientExpectedAmount);
  }

  function test_claimUnregisteredRecipient() public {
    //claim the individual rewards for this recipient
    claimRewards(PRIMARY_POOL_ARRAY, getSecondaryRecipients()[1].addr);
    claimRewards(SECONDARY_POOL_ARRAY, getPrimaryRecipients()[1].addr);

    //check the recipients didn't receive any fees from this pool
    assertEq(getAssetBalance(getSecondaryRecipients()[1].addr), 0);
    assertEq(getAssetBalance(getPrimaryRecipients()[1].addr), 0);

    //check the rewardManager has the remaining quantity
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 2);
  }

  function test_claimUnevenAmountRoundsDown() public {
    //adding an uneven amount of dust to each pool, this should round down to the nearest whole number with 4 remaining in the contract
    addFundsToPool(PRIMARY_POOL_ID, getAsset(3), FEE_MANAGER);
    addFundsToPool(SECONDARY_POOL_ID, getAsset(1), FEE_MANAGER);

    //the recipient should have received 1/4 of the deposit amount for each pool
    uint256 recipientExpectedAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), recipientExpectedAmount);
    }

    //special case to handle the first recipient of each pool as they're the same address
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);

    //check the balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(getSecondaryRecipients()[0].addr), recipientExpectedAmount * 2);

    //claim funds for each recipient of the secondary pool except the first
    for (uint256 i = 1; i < getSecondaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getSecondaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), recipientExpectedAmount);
    }

    //contract should have 4 remaining
    assertEq(getAssetBalance(address(rewardManager)), 4);
  }

  function test_singleRecipientClaimMultipleDeposits() public {
    //get the recipient that is claiming
    Common.AddressAndWeight memory recipient = getSecondaryRecipients()[0];

    //claim the individual rewards for this recipient
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //the recipient should have received 1/4 of the deposit amount
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity, which is 3/4 of the initial deposit plus the deposit from the second pool
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 2 - expectedRecipientAmount);

    //add funds to the pool to be split among the recipients
    addFundsToPool(SECONDARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //claim the individual rewards for this recipient
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //the recipient should have received 1/4 of the next deposit amount
    expectedRecipientAmount += POOL_DEPOSIT_AMOUNT / 4;

    //check the recipients balance matches the ratio the recipient should have received, which is 1/4 of each deposit
    assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);

    //check the rewardManager has the remaining quantity, which is now 3/4 of both deposits
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT * 3 - expectedRecipientAmount);
  }

  function test_recipientsClaimMultipleDeposits() public {
    //the recipient should have received 1/4 of the deposit amount
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getSecondaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getSecondaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //the reward manager balance should contain only the funds of the secondary pool
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);

    //add funds to the pool to be split among the recipients
    addFundsToPool(SECONDARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //special case to handle the first recipient of each pool as they're the same address
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);

    //check the balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(getSecondaryRecipients()[0].addr), expectedRecipientAmount * 2);

    //claim funds for each recipient within the pool except the first
    for (uint256 i = 1; i < getSecondaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getSecondaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount * 2);
    }

    //the reward manager balance should again be the balance of the secondary pool as the primary pool has been emptied twice
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_claimEmptyPoolWhenSecondPoolContainsFunds() public {
    //the recipient should have received 1/4 of the deposit amount
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim all rewards for each recipient in the primary pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //claim all the rewards again for the first recipient as that address is a member of both pools
    claimRewards(PRIMARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);

    //check the balance
    assertEq(getAssetBalance(getSecondaryRecipients()[0].addr), expectedRecipientAmount);
  }

  function test_getRewardsAvailableToRecipientInBothPools() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(
      getPrimaryRecipients()[0].addr,
      0,
      type(uint256).max
    );

    //check the recipient is in both pools
    assertEq(poolIds[0], PRIMARY_POOL_ID);
    assertEq(poolIds[1], SECONDARY_POOL_ID);
  }

  function test_getRewardsAvailableToRecipientInSinglePool() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(
      getPrimaryRecipients()[1].addr,
      0,
      type(uint256).max
    );

    //check the recipient is in both pools
    assertEq(poolIds[0], PRIMARY_POOL_ID);
    assertEq(poolIds[1], ZERO_POOL_ID);
  }

  function test_getRewardsAvailableToRecipientInNoPools() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(FEE_MANAGER, 0, type(uint256).max);

    //check the recipient is in neither pool
    assertEq(poolIds[0], ZERO_POOL_ID);
    assertEq(poolIds[1], ZERO_POOL_ID);
  }

  function test_getRewardsAvailableToRecipientInBothPoolsWhereAlreadyClaimed() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(
      getPrimaryRecipients()[0].addr,
      0,
      type(uint256).max
    );

    //check the recipient is in both pools
    assertEq(poolIds[0], PRIMARY_POOL_ID);
    assertEq(poolIds[1], SECONDARY_POOL_ID);

    //claim the rewards for each pool
    claimRewards(PRIMARY_POOL_ARRAY, getPrimaryRecipients()[0].addr);
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);

    //get the available pools again
    poolIds = rewardManager.getAvailableRewardPoolIds(getPrimaryRecipients()[0].addr, 0, type(uint256).max);

    //user should not be in any pool
    assertEq(poolIds[0], ZERO_POOL_ID);
    assertEq(poolIds[1], ZERO_POOL_ID);
  }

  function test_getAvailableRewardsCursorCannotBeGreaterThanTotalPools() public {
    vm.expectRevert(INVALID_POOL_LENGTH_SELECTOR);

    rewardManager.getAvailableRewardPoolIds(FEE_MANAGER, type(uint256).max, 0);
  }

  function test_getAvailableRewardsCursorAndTotalPoolsEqual() public {
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(getPrimaryRecipients()[0].addr, 2, 2);

    assertEq(poolIds.length, 0);
  }

  function test_getAvailableRewardsCursorSingleResult() public {
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(getPrimaryRecipients()[0].addr, 0, 1);

    assertEq(poolIds[0], PRIMARY_POOL_ID);
  }
}

contract RewardManagerRecipientClaimDifferentWeightsTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function getPrimaryRecipients() public virtual override returns (Common.AddressAndWeight[] memory) {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

    //init each recipient with uneven weights
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, TEN_PERCENT);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, TEN_PERCENT * 8);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, ONE_PERCENT * 6);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, ONE_PERCENT * 4);

    return recipients;
  }

  function test_allRecipientsClaimingReceiveExpectedAmount() public {
    //loop all the recipients and claim their expected amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //the recipient should have received a share proportional to their weight
      uint256 expectedRecipientAmount = (POOL_DEPOSIT_AMOUNT * recipient.weight) / POOL_SCALAR;

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }
  }
}

contract RewardManagerRecipientClaimUnevenWeightTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();
  }

  function getPrimaryRecipients() public virtual override returns (Common.AddressAndWeight[] memory) {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](2);

    uint64 oneThird = POOL_SCALAR / 3;

    //init each recipient with even weights.
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, oneThird);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 2 * oneThird + 1);

    return recipients;
  }

  function test_allRecipientsClaimingReceiveExpectedAmountWithSmallDeposit() public {
    //add a smaller amount of funds to the pool
    uint256 smallDeposit = 1e8;

    //add a smaller amount of funds to the pool
    addFundsToPool(PRIMARY_POOL_ID, getAsset(smallDeposit), FEE_MANAGER);

    //loop all the recipients and claim their expected amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //the recipient should have received a share proportional to their weight
      uint256 expectedRecipientAmount = (smallDeposit * recipient.weight) / POOL_SCALAR;

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //smaller deposits will consequently have less precision and will not be able to be split as evenly, the remaining 1 will be lost due to 333...|... being paid out instead of 333...4|
    assertEq(getAssetBalance(address(rewardManager)), 1);
  }

  function test_allRecipientsClaimingReceiveExpectedAmount() public {
    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //loop all the recipients and claim their expected amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //the recipient should have received a share proportional to their weight
      uint256 expectedRecipientAmount = (POOL_DEPOSIT_AMOUNT * recipient.weight) / POOL_SCALAR;

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //their should be 0 wei left over indicating a successful split
    assertEq(getAssetBalance(address(rewardManager)), 0);
  }
}

contract RewardManagerNoRecipientSet is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //add funds to the pool to be split among the recipients once registered
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function test_claimAllRecipientsAfterRecipientsSet() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //try and claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //there should be no rewards claimed as the recipient is not registered
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the recipient received nothing
      assertEq(getAssetBalance(recipient.addr), 0);
    }

    //Set the recipients after the rewards have been paid into the pool
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), ADMIN);

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient that is claiming
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //there should be no rewards claimed as the recipient is registered
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }
  }
}
