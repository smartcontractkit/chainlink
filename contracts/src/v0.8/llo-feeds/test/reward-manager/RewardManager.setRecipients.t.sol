// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseRewardManagerTest} from "./BaseRewardManager.t.sol";
import {Common} from "../../../libraries/Common.sol";
import {RewardManager} from "../../RewardManager.sol";

/**
 * @title BaseRewardManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the setRecipient functionality of the RewardManager contract
 */
contract RewardManagerSetRecipientsTest is BaseRewardManagerTest {
  function setUp() public override {
    //setup contracts
    super.setUp();
  }

  function test_setRewardRecipients() public {
    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), ADMIN);
  }

  function test_setRewardRecipientsIsEmpty() public {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

    //should revert if the recipients array is empty
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_setRewardRecipientWithZeroAddress() public {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = getPrimaryRecipients();

    //override the first recipient with a zero address
    recipients[0].addr = address(0);

    //should revert if the recipients array is empty
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_setRewardRecipientWeights() public {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](4);

    //init each recipient with even weights
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, 25);
    recipients[1] = Common.AddressAndWeight(DEFAULT_RECIPIENT_2, 25);
    recipients[2] = Common.AddressAndWeight(DEFAULT_RECIPIENT_3, 25);
    recipients[3] = Common.AddressAndWeight(DEFAULT_RECIPIENT_4, 25);

    //should revert if the recipients array is empty
    vm.expectRevert(INVALID_WEIGHT_ERROR_SELECTOR);

    //set the recipients with a recipient with a weight of 100%
    setRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_setSingleRewardRecipient() public {
    //array of recipients
    Common.AddressAndWeight[] memory recipients = new Common.AddressAndWeight[](1);

    //init each recipient with even weights
    recipients[0] = Common.AddressAndWeight(DEFAULT_RECIPIENT_1, POOL_SCALAR);

    //set the recipients with a recipient with a weight of 100%
    setRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }

  function test_setRewardRecipientTwice() public {
    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), ADMIN);

    //should revert if recipients for this pool have already been set
    vm.expectRevert(INVALID_POOL_ID_ERROR_SELECTOR);

    //set the recipients again
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), ADMIN);
  }

  function test_setRewardRecipientFromNonOwnerOrFeeManagerAddress() public {
    //should revert if the sender is not the owner or proxy
    vm.expectRevert(UNAUTHORIZED_ERROR_SELECTOR);

    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), USER);
  }

  function test_setRewardRecipientFromManagerAddress() public {
    //update the proxy address
    setFeeManager(FEE_MANAGER, ADMIN);

    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), FEE_MANAGER);
  }

  function test_eventIsEmittedUponSetRecipients() public {
    //expect an emit
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit RewardRecipientsUpdated(PRIMARY_POOL_ID, getPrimaryRecipients());

    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, getPrimaryRecipients(), ADMIN);
  }

  function test_setRecipientContainsDuplicateRecipients() public {
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

    //set the recipients
    setRewardRecipients(PRIMARY_POOL_ID, recipients, ADMIN);
  }
}
