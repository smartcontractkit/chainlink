// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseRewardManagerTest} from "./BaseRewardManager.t.sol";
import {Common} from "../../../libraries/internal/Common.sol";

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
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
  }

  function test_claimAllRecipients() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }
  }

  function test_claimSingleRecipient() public {
    //get the recipient we're claiming for
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
    //adding 1 to the pool should leave 1 wei worth of dust, which we don't handle due to it being economically infeasible
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(1));

    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
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
    //get the recipient we're claiming for
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //claim the individual rewards for this recipient
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //check the recipients balance is still 0 as there's no pool to receive fees from
    assertEq(getAssetBalance(recipient.addr), 0);

    //check the rewardManager has the full amount
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_singleRecipientClaimMultipleDeposits() public {
    //get the recipient we're claiming for
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
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

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
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //the reward manager balance should be 0 as all of the funds have been claimed
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
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
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
    addFundsToPool(SECONDARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
  }

  function test_claimAllRecipientsSinglePool() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //check the pool balance is still equal to DEPOSIT_AMOUNT as we only claimed for one of the pools
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
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //claim funds for each recipient within the pool
    for (uint256 i = 1; i < getSecondaryRecipients().length; i++) {
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getSecondaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //special case to handle the first recipient of each pool as they're the same address
    Common.AddressAndWeight memory recipient = getPrimaryRecipients()[0];

    //claim the individual rewards for each pool
    claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);
    claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

    //check the balance matches the ratio the recipient should have received, which is 1/4 of each deposit for each pool
    assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount * 2);
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
    //adding an unevent amount of dust to each pool, this should round down to the nearest whole number with 4 remaining in the contract
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(3));
    addFundsToPool(SECONDARY_POOL_ID, USER, getAsset(1));

    //the recipient should have received 1/4 of the deposit amount for each pool
    uint256 recipientExpectedAmount = POOL_DEPOSIT_AMOUNT / 4;

    //claim funds for each recipient within the pool
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
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
      //get the recipient we're claiming for
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
    //get the recipient we're claiming for
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
    addFundsToPool(SECONDARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

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
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getSecondaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(SECONDARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //the reward manager balance should contain only the funds of the secondary pool
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);

    //add funds to the pool to be split among the recipients
    addFundsToPool(SECONDARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //special case to handle the first recipient of each pool as they're the same address
    claimRewards(SECONDARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);

    //check the balance matches the ratio the recipient should have received
    assertEq(getAssetBalance(getSecondaryRecipients()[0].addr), expectedRecipientAmount * 2);

    //claim funds for each recipient within the pool except the first
    for (uint256 i = 1; i < getSecondaryRecipients().length; i++) {
      //get the recipient we're claiming for
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
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //claim all the rewards again for the first recipient as that address is a member of both pools
    claimRewards(PRIMARY_POOL_ARRAY, getSecondaryRecipients()[0].addr);

    //the balance should still be 1/4 of the deposit amount
    assertEq(getAssetBalance(getSecondaryRecipients()[0].addr), expectedRecipientAmount);
  }

  function test_getRewardsAvailableToRecipientInBothPools() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(getPrimaryRecipients()[0].addr);

    //check the recipient is in both pools
    assertEq(poolIds[0], PRIMARY_POOL_ID);
    assertEq(poolIds[1], SECONDARY_POOL_ID);
  }

  function test_getRewardsAvailableToRecipientInSinglePool() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(getPrimaryRecipients()[1].addr);

    //check the recipient is in both pools
    assertEq(poolIds[0], PRIMARY_POOL_ID);
    assertEq(poolIds[1], ZERO_POOL_ID);
  }

  function test_getRewardsAvailableToRecipientInNoPools() public {
    //get index 0 as this recipient is in both default pools
    bytes32[] memory poolIds = rewardManager.getAvailableRewardPoolIds(USER);

    //check the recipient is in neither pool
    assertEq(poolIds[0], ZERO_POOL_ID);
    assertEq(poolIds[1], ZERO_POOL_ID);
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
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
  }

  function getPrimaryRecipients() public virtual override returns (Common.AddressAndWeight[] memory) {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

    //init each recipient with unevent weights. 1000 = 25% of pool
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 1000);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 8000);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, 600);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, 400);

    return recipients;
  }

  function test_allRecipientsClaimingReceiveExpectedAmount() public {
    //loop all the recipients and claim their expected amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
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

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
  }

  function getPrimaryRecipients() public virtual override returns (Common.AddressAndWeight[] memory) {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](2);

    //init each recipient with even weights. 2500 = 25% of pool
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 3333);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 6667);

    return recipients;
  }

  function test_allRecipientsClaimingReceiveExpectedAmount() public {
    //loop all the recipients and claim their expected amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //get the recipient we're claiming for
      Common.AddressAndWeight memory recipient = getPrimaryRecipients()[i];

      //claim the individual rewards for each recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipient.addr);

      //the recipient should have received a share proportional to their weight
      uint256 expectedRecipientAmount = (POOL_DEPOSIT_AMOUNT * recipient.weight) / POOL_SCALAR;

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipient.addr), expectedRecipientAmount);
    }

    //their should be 2 wei left over from rounding errors
    assertEq(getAssetBalance(address(rewardManager)), 2);
  }
}
