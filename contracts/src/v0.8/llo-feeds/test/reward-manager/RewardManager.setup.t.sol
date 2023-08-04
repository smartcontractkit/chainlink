// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {BaseRewardManagerTest} from "./BaseRewardManager.t.sol";
import {Common} from "../../../libraries/Common.sol";
import {RewardManager} from "../../RewardManager.sol";

/**
 * @title BaseRewardManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the core functionality of the RewardManager contract
 */
contract RewardManagerSetupTest is BaseRewardManagerTest {
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

  function test_eventEmittedUponProxyUpdate() public {
    //expect the event to be emitted
    vm.expectEmit();

    //emit the event that is expected to be emitted
    emit FeeManagerUpdated(USER);

    //set the verifier proxy
    setFeeManager(USER, ADMIN);
  }
}
