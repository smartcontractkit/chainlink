// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.19;

import {BaseDestinationRewardManagerTest} from "./BaseDestinationRewardManager.t.sol";
import {DestinationRewardManager} from "../../../v0.4.0/DestinationRewardManager.sol";
import {ERC20Mock} from "../../../../vendor/openzeppelin-solidity/v4.8.3/contracts/mocks/ERC20Mock.sol";
import {IDestinationRewardManager} from "../../interfaces/IDestinationRewardManager.sol";

/**
 * @title DestinationRewardManagerSetupTest
 * @author Michael Fletcher
 * @notice This contract will test the core functionality of the DestinationRewardManager contract
 */
contract DestinationRewardManagerSetupTest is BaseDestinationRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();
  }

  function test_rejectsZeroLinkAddressOnConstruction() public {
    //should revert if the contract is a zero address
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //create a rewardManager with a zero link address
    new DestinationRewardManager(address(0));
  }

  function test_eventEmittedUponFeeManagerUpdate() public {
    //expect the event to be emitted
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit FeeManagerUpdated(FEE_MANAGER_2);

    //set the verifier proxy
    setFeeManager(FEE_MANAGER_2, ADMIN);
  }

  function test_eventEmittedUponFeePaid() public {
    //create pool and add funds
    createPrimaryPool();

    //change to the feeManager who is the one who will be paying the fees
    changePrank(FEE_MANAGER);

    //approve the amount being paid into the pool
    ERC20Mock(getAsset(POOL_DEPOSIT_AMOUNT).assetAddress).approve(address(rewardManager), POOL_DEPOSIT_AMOUNT);

    IDestinationRewardManager.FeePayment[] memory payments = new IDestinationRewardManager.FeePayment[](1);
    payments[0] = IDestinationRewardManager.FeePayment(PRIMARY_POOL_ID, uint192(POOL_DEPOSIT_AMOUNT));

    //event is emitted when funds are added
    vm.expectEmit();
    emit FeePaid(payments, FEE_MANAGER);

    //this represents the verifier adding some funds to the pool
    rewardManager.onFeePaid(payments, FEE_MANAGER);
  }

  function test_setFeeManagerZeroAddress() public {
    //should revert if the contract is a zero address
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //set the verifier proxy
    setFeeManager(address(0), ADMIN);
  }

  function test_addFeeManagerZeroAddress() public {
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);
    rewardManager.addFeeManager(address(0));
  }

  function test_addFeeManagerExistingAddress() public {
    address dummyAddress = address(998);
    rewardManager.addFeeManager(dummyAddress);
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);
    rewardManager.addFeeManager(dummyAddress);
  }

  function test_removeFeeManagerNonExistentAddress() public {
    address dummyAddress = address(991);
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);
    rewardManager.removeFeeManager(dummyAddress);
  }

  function test_addRemoveFeeManager() public {
    address dummyAddress1 = address(1);
    address dummyAddress2 = address(2);
    rewardManager.addFeeManager(dummyAddress1);
    rewardManager.addFeeManager(dummyAddress2);
    assertEq(rewardManager.s_feeManagerAddressList(dummyAddress1), dummyAddress1);
    assertEq(rewardManager.s_feeManagerAddressList(dummyAddress2), dummyAddress2);
    rewardManager.removeFeeManager(dummyAddress1);
    assertEq(rewardManager.s_feeManagerAddressList(dummyAddress1), address(0));
    assertEq(rewardManager.s_feeManagerAddressList(dummyAddress2), dummyAddress2);
  }
}
