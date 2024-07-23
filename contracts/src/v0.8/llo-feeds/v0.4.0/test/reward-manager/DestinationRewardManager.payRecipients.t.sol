// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseDestinationRewardManagerTest} from "./BaseDestinationRewardManager.t.sol";
import {IDestinationRewardManager} from "../../interfaces/IDestinationRewardManager.sol";

/**
 * @title DestinationRewardManagerPayRecipientsTest
 * @author Michael Fletcher
 * @notice This contract will test the payRecipients functionality of the RewardManager contract
 */
contract DestinationRewardManagerPayRecipientsTest is BaseDestinationRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();

    //create a single pool for these tests
    createPrimaryPool();

    //add funds to the pool to be split among the recipients
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function test_payAllRecipients() public {
    //pay all the recipients in the pool
    payRecipients(PRIMARY_POOL_ID, getPrimaryRecipientAddresses(), ADMIN);

    //each recipient should receive 1/4 of the pool
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check each recipient received the correct amount
    for (uint256 i = 0; i < getPrimaryRecipientAddresses().length; i++) {
      assertEq(getAssetBalance(getPrimaryRecipientAddresses()[i]), expectedRecipientAmount);
    }
  }

  function test_paySingleRecipient() public {
    //get the first individual recipient
    address recipient = getPrimaryRecipientAddresses()[0];

    //get a single recipient as an array
    address[] memory recipients = new address[](1);
    recipients[0] = recipient;

    //pay a single recipient
    payRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //the recipient should have received 1/4 of the deposit amount
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    assertEq(getAssetBalance(recipient), expectedRecipientAmount);
  }

  function test_payRecipientWithInvalidPool() public {
    //get the first individual recipient
    address recipient = getPrimaryRecipientAddresses()[0];

    //get a single recipient as an array
    address[] memory recipients = new address[](1);
    recipients[0] = recipient;

    //pay a single recipient
    payRecipients(SECONDARY_POOL_ID, recipients, ADMIN);

    //the recipient should have received nothing
    assertEq(getAssetBalance(recipient), 0);
  }

  function test_payRecipientsEmptyRecipientList() public {
    //get a single recipient
    address[] memory recipients = new address[](0);

    //pay a single recipient
    payRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //rewardManager should have the full balance
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_payAllRecipientsWithAdditionalUnregisteredRecipient() public {
    //load all the recipients and add an additional one who is not in the pool
    address[] memory recipients = new address[](getPrimaryRecipientAddresses().length + 1);
    for (uint256 i = 0; i < getPrimaryRecipientAddresses().length; i++) {
      recipients[i] = getPrimaryRecipientAddresses()[i];
    }
    recipients[recipients.length - 1] = DEFAULT_RECIPIENT_5;

    //pay the recipients
    payRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //each recipient should receive 1/4 of the pool except the last
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check each recipient received the correct amount
    for (uint256 i = 0; i < getPrimaryRecipientAddresses().length; i++) {
      assertEq(getAssetBalance(getPrimaryRecipientAddresses()[i]), expectedRecipientAmount);
    }

    //the unregistered recipient should receive nothing
    assertEq(getAssetBalance(DEFAULT_RECIPIENT_5), 0);
  }

  function test_payAllRecipientsWithAdditionalInvalidRecipient() public {
    //load all the recipients and add an additional one which is invalid, that should receive nothing
    address[] memory recipients = new address[](getPrimaryRecipientAddresses().length + 1);
    for (uint256 i = 0; i < getPrimaryRecipientAddresses().length; i++) {
      recipients[i] = getPrimaryRecipientAddresses()[i];
    }
    recipients[recipients.length - 1] = INVALID_ADDRESS;

    //pay the recipients
    payRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //each recipient should receive 1/4 of the pool except the last
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check each recipient received the correct amount
    for (uint256 i = 0; i < getPrimaryRecipientAddresses().length; i++) {
      assertEq(getAssetBalance(getPrimaryRecipientAddresses()[i]), expectedRecipientAmount);
    }
  }

  function test_paySubsetOfRecipientsInPool() public {
    //load a subset of the recipients into an array
    address[] memory recipients = new address[](getPrimaryRecipientAddresses().length - 1);
    for (uint256 i = 0; i < recipients.length; i++) {
      recipients[i] = getPrimaryRecipientAddresses()[i];
    }

    //pay the subset of recipients
    payRecipients(PRIMARY_POOL_ID, recipients, ADMIN);

    //each recipient should receive 1/4 of the pool except the last
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check each subset of recipients received the correct amount
    for (uint256 i = 0; i < recipients.length - 1; i++) {
      assertEq(getAssetBalance(recipients[i]), expectedRecipientAmount);
    }

    //check the pool has the remaining balance
    assertEq(
      getAssetBalance(address(rewardManager)),
      POOL_DEPOSIT_AMOUNT - expectedRecipientAmount * recipients.length
    );
  }

  function test_payAllRecipientsFromNonAdminUser() public {
    //should revert if the caller isn't an admin or recipient within the pool
    vm.expectRevert(UNAUTHORIZED_ERROR_SELECTOR);

    //pay all the recipients in the pool
    payRecipients(PRIMARY_POOL_ID, getPrimaryRecipientAddresses(), FEE_MANAGER);
  }

  function test_payAllRecipientsFromRecipientInPool() public {
    //pay all the recipients in the pool
    payRecipients(PRIMARY_POOL_ID, getPrimaryRecipientAddresses(), DEFAULT_RECIPIENT_1);

    //each recipient should receive 1/4 of the pool
    uint256 expectedRecipientAmount = POOL_DEPOSIT_AMOUNT / 4;

    //check each recipient received the correct amount
    for (uint256 i = 0; i < getPrimaryRecipientAddresses().length; i++) {
      assertEq(getAssetBalance(getPrimaryRecipientAddresses()[i]), expectedRecipientAmount);
    }
  }

  function test_payRecipientsWithInvalidPoolId() public {
    //pay all the recipients in the pool
    payRecipients(INVALID_POOL_ID, getPrimaryRecipientAddresses(), ADMIN);

    //pool should still contain the full balance
    assertEq(getAssetBalance(address(rewardManager)), POOL_DEPOSIT_AMOUNT);
  }

  function test_addFundsToPoolAsOwner() public {
    //add funds to the pool
    addFundsToPool(PRIMARY_POOL_ID, getAsset(POOL_DEPOSIT_AMOUNT), FEE_MANAGER);
  }

  function test_addFundsToPoolAsNonOwnerOrFeeManager() public {
    //should revert if the caller isn't an admin or recipient within the pool
    vm.expectRevert(UNAUTHORIZED_ERROR_SELECTOR);

    IDestinationRewardManager.FeePayment[] memory payments = new IDestinationRewardManager.FeePayment[](1);
    payments[0] = IDestinationRewardManager.FeePayment(PRIMARY_POOL_ID, uint192(POOL_DEPOSIT_AMOUNT));

    //add funds to the pool
    rewardManager.onFeePaid(payments, USER);
  }
}
