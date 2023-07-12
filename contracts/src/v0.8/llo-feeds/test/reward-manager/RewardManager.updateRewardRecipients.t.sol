// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseRewardManagerTest} from "./BaseRewardManagerTest.t.sol";
import {Common} from "../../../libraries/internal/Common.sol";
import {RewardManager} from "../../RewardManager.sol";
import "forge-std/console.sol";

/**
 * @title BaseRewardManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the updateRecipient functionality of the RewardManager contract
 */
contract RewardManagerUpdateRewardRecipientsTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
  }

  function test_onlyAdminCanUpdateRecipients() public {
    //should revert if the caller is not the admin
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //updating a recipient should force the funds to be paid out
    updateRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), USER);
  }

  function test_updateAllRecipientsWithSameAddressAndWeight() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //updating a recipient should force the funds to be paid out
    updateRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), ADMIN);

    //check each recipient received the correct amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(getPrimaryRecipients()[i].addr), expectedRecipientAmount);
    }
  }

  function test_updatePartialRecipientsWithSameAddressAndWeight() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //get a subset of the recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](2);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 2500);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 2500);

    //updating a recipient should force the funds to be paid out
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each recipient received the correct amount
    for (uint256 i; i < recipients.length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedRecipientAmount);
    }

    //the reward manager should still have half remaining funds
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT / 2);
  }

  function test_updateRecipientWithNewZeroAddress() public {
    //create a new array to hold the existing recipients plus a new zero address
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 1);

    //add all the existing recipients
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      recipients[i] = getPrimaryRecipients()[i];
    }
    //add a new address to the primary recipients
    recipients[recipients.length - 1] = Common.AddressAndWeight(address(0), 0);

    //should revert if the recipient is a zero address
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //update the recipients with invalid address
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsContainsDuplicateRecipients() public {
    //create a new array to hold the existing recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length * 2);

    //add all the existing recipients
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      recipients[i] = getPrimaryRecipients()[i];
    }
    //add all the existing recipients again
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      recipients[i + getPrimaryRecipients().length] = getPrimaryRecipients()[i];
    }

    //should revert as the list contains a duplicate
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //update the recipients with the duplicate addresses
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsToDifferentSet() public {
    //expected recipient amount is 1/4 of the pool deposit
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 4);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 2500);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 2500);
    recipients[6] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, 2500);
    recipients[7] = Common.AddressAndWeight(DEFAULT_RECIPIENT_8, 2500);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each primary recipient received the correct amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(getPrimaryRecipients()[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //loop through the new recipients and claim the rewards
    for (uint256 i = 4; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedRecipientAmount);
    }
  }

  function test_updateRecipientsToDifferentPartialSet() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 2);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 5000);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 5000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each primary recipient received the correct amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(getPrimaryRecipients()[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //expected recipient amount is 1/2 of the pool deposit for new recipients
    uint256 expectedNewRecipientAmount = POOL_DEPOSIT_AMOUNT / 2;

    //loop through the new recipients and claim the rewards
    for (uint256 i = 4; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedNewRecipientAmount);
    }
  }

  function test_updateRecipientsToDifferentLargerSet() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 5);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 2000);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 2000);
    recipients[6] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, 2000);
    recipients[7] = Common.AddressAndWeight(DEFAULT_RECIPIENT_8, 2000);
    recipients[8] = Common.AddressAndWeight(DEFAULT_RECIPIENT_9, 2000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each primary recipient received the correct amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(getPrimaryRecipients()[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //expected recipient amount is 1/2 of the pool deposit for new recipients
    uint256 expectedNewRecipientAmount = POOL_DEPOSIT_AMOUNT / 5;

    //loop through the new recipients and claim the rewards
    for (uint256 i = 4; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedNewRecipientAmount);
    }
  }

  function test_updateRecipientsUpdateAndRemoveExistingForLargerSet() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](9);

    //update the existing recipients
    recipients[0] = Common.AddressAndWeight(getPrimaryRecipients()[0].addr, 0);
    recipients[1] = Common.AddressAndWeight(getPrimaryRecipients()[1].addr, 0);
    recipients[2] = Common.AddressAndWeight(getPrimaryRecipients()[2].addr, 3000);
    recipients[3] = Common.AddressAndWeight(getPrimaryRecipients()[3].addr, 2000);

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 1000);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 1000);
    recipients[6] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, 1000);
    recipients[7] = Common.AddressAndWeight(DEFAULT_RECIPIENT_8, 1000);
    recipients[8] = Common.AddressAndWeight(DEFAULT_RECIPIENT_9, 1000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each primary recipient received the correct amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(getPrimaryRecipients()[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //expected recipient amount for the new recipients is 10% of the pool
    uint256 expectedNewRecipientAmount = POOL_DEPOSIT_AMOUNT / 10;

    //claim the rewards for the updated recipients manually
    claimRewards(PRIMARY_POOL_ARRAY, recipients[2].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[3].addr);

    //check the balance matches the ratio the recipient who were updated should have received
    assertEq(getAssetBalance(recipients[2].addr), (POOL_DEPOSIT_AMOUNT * 3000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(recipients[3].addr), (POOL_DEPOSIT_AMOUNT * 2000) / 10000 + expectedRecipientAmount);

    //loop through the new recipients and claim the rewards
    for (uint256 i = 4; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);

      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedNewRecipientAmount);
    }
  }

  function test_updateRecipientsUpdateAndRemoveExistingForSmallerSet() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](5);

    //update the existing recipients
    recipients[0] = Common.AddressAndWeight(getPrimaryRecipients()[0].addr, 0);
    recipients[1] = Common.AddressAndWeight(getPrimaryRecipients()[1].addr, 0);
    recipients[2] = Common.AddressAndWeight(getPrimaryRecipients()[2].addr, 3000);
    recipients[3] = Common.AddressAndWeight(getPrimaryRecipients()[3].addr, 2000);

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 5000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each primary recipient received the correct amount
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(getPrimaryRecipients()[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //expected recipient amount for the new recipient is 50% of the pool
    uint256 expectedNewRecipientAmount = POOL_DEPOSIT_AMOUNT / 2;

    //claim the rewards for the updated recipients manually
    claimRewards(PRIMARY_POOL_ARRAY, recipients[2].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[3].addr);

    //check the balance matches the ratio the recipient who were updated should have received
    assertEq(getAssetBalance(recipients[2].addr), (POOL_DEPOSIT_AMOUNT * 3000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(recipients[3].addr), (POOL_DEPOSIT_AMOUNT * 2000) / 10000 + expectedRecipientAmount);

    //the remaining recipient should receive 50% upon claiming
    claimRewards(PRIMARY_POOL_ARRAY, recipients[4].addr);
    assertEq(getAssetBalance(recipients[4].addr), expectedNewRecipientAmount);
  }

  function test_updateRecipientsToDifferentSetWithInvalidWeights() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 2);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 5000);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 1000);

    //should revert as the weights do not add up to 10000
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updatePartialRecipientsToSubset() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 0);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 0);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, 5000);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, 5000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each recipient received the correct amount
    for (uint256 i; i < recipients.length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //claim the rewards for the updated recipients manually
    claimRewards(PRIMARY_POOL_ARRAY, recipients[2].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[3].addr);

    //check the balance matches the ratio the recipient who were updated should have received
    assertEq(getAssetBalance(recipients[2].addr), (POOL_DEPOSIT_AMOUNT * 5000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(recipients[3].addr), (POOL_DEPOSIT_AMOUNT * 5000) / 10000 + expectedRecipientAmount);
  }

  function test_updatePartialRecipientsToDifferentSetUnderWeight() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 0);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 0);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 2500);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 2499);

    //should revert as the new weights exceed the previous weights being replaced
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updatePartialRecipientsToDifferentExcessiveWeight() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 0);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 0);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, 2500);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, 2501);

    //should revert as the new weights exceed the previous weights being replaced
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientWeights() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set with their new weights
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 1000);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 1000);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, 3000);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, 5000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each recipient received the correct amount
    for (uint256 i; i < recipients.length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have no funds remaining
    assertEq(getAssetBalance(address(rewardManager)), 0);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //loop each user and claim the rewards
    for (uint256 i; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);
    }

    //manually check the balance of each recipient which should be their original amount of 1/4 plus their new weighted amount
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_1), (POOL_DEPOSIT_AMOUNT * 1000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_2), (POOL_DEPOSIT_AMOUNT * 1000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_3), (POOL_DEPOSIT_AMOUNT * 3000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_4), (POOL_DEPOSIT_AMOUNT * 5000) / 10000 + expectedRecipientAmount);
  }

  function test_partialUpdateRecipientWeights() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set with their new weights, which should total 5000
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](2);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 1000);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 4000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each recipient received the correct amount
    for (uint256 i; i < recipients.length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedRecipientAmount);
    }

    //the reward manager should have half the funds remaining
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT / 2);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //loop each user and claim the rewards
    for (uint256 i; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);
    }

    //manually check the balance of each recipient which should be their original amount of 1/4 plus their new weighted amount
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_1), (POOL_DEPOSIT_AMOUNT * 1000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_2), (POOL_DEPOSIT_AMOUNT * 4000) / 10000 + expectedRecipientAmount);

    //the reward manager should have half the funds remaining
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }
}

contract RewardManagerUpdateRewardRecipientsMultiplePoolsTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();
    createSecondaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
    addFundsToPool(SECONDARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));
  }

  function getSecondaryRecipients() public override returns (Common.AddressAndWeight[] memory) {
    //for testing purposes, we want the primary and secondary pool to contain the same recipients
    return getPrimaryRecipients();
  }

  function test_updatePrimaryRecipientWeights() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 0);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 0);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, 5000);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, 5000);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //check each recipient received the correct amount
    for (uint256 i; i < recipients.length; i++) {
      //check the balance matches the ratio the recipient should have received
      assertEq(getAssetBalance(recipients[i].addr), expectedRecipientAmount);
    }

    //the reward manager should still have the funds for the secondary pool
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);

    //add more funds to the pool to check new distribution
    addFundsToPool(PRIMARY_POOL_ID, USER, getAsset(POOL_DEPOSIT_AMOUNT));

    //expected recipient amount for the new recipient is 50% of the pool
    uint256 expectedNewRecipientAmount = POOL_DEPOSIT_AMOUNT / 2;

    //claim the rewards for the updated recipients manually
    claimRewards(PRIMARY_POOL_ARRAY, recipients[2].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[3].addr);

    //check the balance matches the ratio the recipient who were updated should have received
    assertEq(getAssetBalance(recipients[2].addr), (POOL_DEPOSIT_AMOUNT * 5000) / 10000 + expectedRecipientAmount);
    assertEq(getAssetBalance(recipients[3].addr), (POOL_DEPOSIT_AMOUNT * 5000) / 10000 + expectedRecipientAmount);
  }
}
