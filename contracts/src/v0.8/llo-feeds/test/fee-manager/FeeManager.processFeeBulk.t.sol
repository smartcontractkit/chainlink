// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../dev/FeeManager.sol";
import {Common} from "../../../libraries/Common.sol";
import "./BaseFeeManager.t.sol";
import {IRewardManager} from "../../dev/interfaces/IRewardManager.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the feeManager processFee
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  function setUp() public override {
    super.setUp();
  }

  function test_processMultipleLinkReports() public {
    //change
  }

  function test_processMultipleWrappedNativeReports() public {

  }

  function test_processMultipleUnwrappedNativeReports() public {

  }

  function test_processMultipleLinkAndNativeWrappedReports() public {

  }

  function test_processMultipleLinkAndNativeUnwrappedReports() public {
    //change
  }

  function test_processV1V2V3Reports() public {
    //v2 link v3 native
  }

  function test_processMultipleV1Reports() public {

  }

  function test_processWrappedAndUnwrappedReportsDefaultsToUnwrapped() public {

  }

  function test_eventIsEmittedIfNotEnoughLink() public {
    //link and native
  }

  function test_processLinkReportWithMultipleV1Reports() public {

  }
}
