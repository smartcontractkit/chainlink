// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseRewardManagerTest} from "./BaseRewardManager.t.sol";
import {Common} from "../../../libraries/Common.sol";
import {RewardManager} from "../../RewardManager.sol";
import {ERC20Mock} from "../../../vendor/openzeppelin-solidity/v4.8.0/contracts/mocks/ERC20Mock.sol";
import {IRewardManager} from "../../interfaces/IRewardManager.sol";

/**
 * @title BaseRewardManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the core functionality of the RewardManager contract
 */
contract RewardManagerSetupTest is BaseRewardManagerTest {
  uint256 internal constant POOL_DEPOSIT_AMOUNT = 10e18;

  function setUp() public override {
    //setup contracts
    super.setUp();
  }

  function test_rejectsZeroLinkAddressOnConstruction() public {
    //should revert if the contract is a zero address
    vm.expectRevert(INVALID_ADDRESS_ERROR_SELECTOR);

    //create a rewardManager with a zero link address
    new RewardManager(address(0));
  }

  function test_eventEmittedUponFeeManagerUpdate() public {
    //expect the event to be emitted
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit FeeManagerUpdated(FEE_MANAGER);

    //set the verifier proxy
    setFeeManager(FEE_MANAGER, ADMIN);
  }

  function test_eventEmittedUponFeePaid() public {
    //create pool and add funds
    createPrimaryPool();

    //change to the feeManager who is the one who will be paying the fees
    changePrank(FEE_MANAGER);

    //approve the amount being paid into the pool
    ERC20Mock(getAsset(POOL_DEPOSIT_AMOUNT).assetAddress).approve(address(rewardManager), POOL_DEPOSIT_AMOUNT);

    IRewardManager.FeePayment[] memory payments = new IRewardManager.FeePayment[](1);
    payments[0] = IRewardManager.FeePayment(PRIMARY_POOL_ID, uint192(POOL_DEPOSIT_AMOUNT));

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
}
