// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {Test} from "forge-std/Test.sol";
import {FeeManager} from "../../FeeManager.sol";
import {Common} from "../../../libraries/internal/Common.sol";
import "./BaseFeeManager.t.sol";

/**
 * @title BaseFeeManagerTest
 * @author Michael Fletcher
 * @notice This contract will test the functionality of the fee managers processfee
 */
contract FeeManagerProcessFeeTest is BaseFeeManagerTest {
  function test_noFeeIsAppliedIfNativeIsZero() public {}
}
