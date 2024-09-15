// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseDestinationRewardManagerTest} from "./BaseDestinationRewardManager.t.sol";
import {Common} from "../../../libraries/Common.sol";

/**
 * @title DestinationRewardManagerUpdateRewardRecipientsTest
 * @author Michael Fletcher
 * @notice This contract will test the updateRecipient functionality of the RewardManager contract
 */
contract DestinationRewardManagerUpdateRewardRecipientsTest is BaseDestinationRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function test_onlyAdminCanUpdateRecipients() public {
    //should revert if the caller is not the admin
    vm.expectRevert(ONLY_CALLABLE_BY_OWNER_ERROR);

    //updating a recipient should force the funds to be paid out
    updateRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), FEE_MANAGER);
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
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, ONE_PERCENT * 25);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, ONE_PERCENT * 25);

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
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 4);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, ONE_PERCENT * 25);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, ONE_PERCENT * 25);
    recipients[6] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, ONE_PERCENT * 25);
    recipients[7] = Common.AddressAndWeight(DEFAULT_RECIPIENT_8, ONE_PERCENT * 25);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsToDifferentPartialSet() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 2);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, FIFTY_PERCENT);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, FIFTY_PERCENT);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsToDifferentLargerSet() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 5);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, TEN_PERCENT * 2);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, TEN_PERCENT * 2);
    recipients[6] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, TEN_PERCENT * 2);
    recipients[7] = Common.AddressAndWeight(DEFAULT_RECIPIENT_8, TEN_PERCENT * 2);
    recipients[8] = Common.AddressAndWeight(DEFAULT_RECIPIENT_9, TEN_PERCENT * 2);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsUpdateAndRemoveExistingForLargerSet() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](9);

    //update the existing recipients
    recipients[0] = Common.AddressAndWeight(getPrimaryRecipients()[0].addr, 0);
    recipients[1] = Common.AddressAndWeight(getPrimaryRecipients()[1].addr, 0);
    recipients[2] = Common.AddressAndWeight(getPrimaryRecipients()[2].addr, TEN_PERCENT * 3);
    recipients[3] = Common.AddressAndWeight(getPrimaryRecipients()[3].addr, TEN_PERCENT * 3);

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, TEN_PERCENT);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, TEN_PERCENT);
    recipients[6] = Common.AddressAndWeight(DEFAULT_RECIPIENT_7, TEN_PERCENT);
    recipients[7] = Common.AddressAndWeight(DEFAULT_RECIPIENT_8, TEN_PERCENT);
    recipients[8] = Common.AddressAndWeight(DEFAULT_RECIPIENT_9, TEN_PERCENT);

    //should revert as the weight does not equal 100%
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsUpdateAndRemoveExistingForSmallerSet() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](5);

    //update the existing recipients
    recipients[0] = Common.AddressAndWeight(getPrimaryRecipients()[0].addr, 0);
    recipients[1] = Common.AddressAndWeight(getPrimaryRecipients()[1].addr, 0);
    recipients[2] = Common.AddressAndWeight(getPrimaryRecipients()[2].addr, TEN_PERCENT * 3);
    recipients[3] = Common.AddressAndWeight(getPrimaryRecipients()[3].addr, TEN_PERCENT * 2);

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, TEN_PERCENT * 5);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updateRecipientsToDifferentSetWithInvalidWeights() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](getPrimaryRecipients().length + 2);
    for (uint256 i; i < getPrimaryRecipients().length; i++) {
      //copy the recipient and set the weight to 0 which implies the recipient is being replaced
      recipients[i] = Common.AddressAndWeight(getPrimaryRecipients()[i].addr, 0);
    }

    //add the new recipients individually
    recipients[4] = Common.AddressAndWeight(DEFAULT_RECIPIENT_5, TEN_PERCENT * 5);
    recipients[5] = Common.AddressAndWeight(DEFAULT_RECIPIENT_6, TEN_PERCENT);

    //should revert as the weight will not equal 100%
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updatePartialRecipientsToSubset() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 0);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 0);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, TEN_PERCENT * 5);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, TEN_PERCENT * 5);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updatePartialRecipientsWithUnderWeightSet() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, TEN_PERCENT);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, TEN_PERCENT);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, TEN_PERCENT);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, TEN_PERCENT);

    //should revert as the new weights exceed the previous weights being replaced
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //updating a recipient should force the funds to be paid out for the primary recipients
    updateRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_updatePartialRecipientsWithExcessiveWeight() public {
    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, TEN_PERCENT);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, TEN_PERCENT);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, TEN_PERCENT);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, POOL_SCALAR);

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
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, TEN_PERCENT);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, TEN_PERCENT);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, TEN_PERCENT * 3);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, TEN_PERCENT * 5);

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
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //loop each user and claim the rewards
    for (uint256 i; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);
    }

    //manually check the balance of each recipient
    assertEq(
      getAssetBalance(DEFAULT_RECIPIENT_1),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(DEFAULT_RECIPIENT_2),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(DEFAULT_RECIPIENT_3),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT * 3) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(DEFAULT_RECIPIENT_4),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT * 5) / POOL_SCALAR + expectedRecipientAmount
    );
  }

  function test_partialUpdateRecipientWeights() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set with their new weights
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](2);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, TEN_PERCENT);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, TEN_PERCENT * 4);

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
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //loop each user and claim the rewards
    for (uint256 i; i < recipients.length; i++) {
      //claim the rewards for this recipient
      claimRewards(PRIMARY_POOL_ARRAY, recipients[i].addr);
    }

    //manually check the balance of each recipient
    assertEq(
      getAssetBalance(DEFAULT_RECIPIENT_1),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(DEFAULT_RECIPIENT_2),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT * 4) / POOL_SCALAR + expectedRecipientAmount
    );

    //the reward manager should have half the funds remaining
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_eventIsEmittedUponUpdateRecipients() public {
    //expect an emit
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit RewardRecipientsUpdated(PRIMARY_POOL_ID, getPrimaryRecipients());

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
}

contract DestinationRewardManagerUpdateRewardRecipientsMultiplePoolsTest is BaseDestinationRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();
    createSecondaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
    addFundsToPool(SECONDARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function getSecondaryRecipients() public override returns (Common.AddressAndWeight[] memory) {
    //for testing purposes, the primary and secondary pool to contain the same recipients
    return getPrimaryRecipients();
  }

  function test_updatePrimaryRecipientWeights() public {
    //expected recipient amount is 1/4 of the pool deposit for original recipients
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //create a list of containing recipients from the primary configured set, and new recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, TEN_PERCENT * 4);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, TEN_PERCENT * 4);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, TEN_PERCENT);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, TEN_PERCENT);

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
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);

    //claim the rewards for the updated recipients manually
    claimRewards(PRIMARY_POOL_ARRAY, recipients[0].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[1].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[2].addr);
    claimRewards(PRIMARY_POOL_ARRAY, recipients[3].addr);

    //check the balance matches the ratio the recipient who were updated should have received
    assertEq(
      getAssetBalance(recipients[0].addr),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT * 4) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(recipients[1].addr),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT * 4) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(recipients[2].addr),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT) / POOL_SCALAR + expectedRecipientAmount
    );
    assertEq(
      getAssetBalance(recipients[3].addr),
      (POOL_DEPOSIT_AMOUNT * TEN_PERCENT) / POOL_SCALAR + expectedRecipientAmount
    );
  }
}
