// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {GasLimitValidator} from "../../dev/GasLimitValidator.sol";
import {L2EPTest} from "./L2EPTest.t.sol";

abstract contract GasLimitValidatorTest is L2EPTest {
  /// Helper constants
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = 0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b;
  uint32 internal constant INIT_GAS_LIMIT = 1900000;

  /// L2EP contract(s)
  GasLimitValidator internal s_validator;

  /// @notice returns a new GasLimitValidator instance
  function newValidator() internal virtual returns (GasLimitValidator validator);

  /// @notice emits the event that should be expected after validate is called
  function emitExpectedSentMessageEvent(
    address validatorAddress,
    bool status,
    uint256 futureTimestampInSeconds
  ) internal virtual;

  /// Setup
  function setUp() public {
    s_validator = newValidator();
  }

  /// @notice it correctly updates the gas limit
  function test_CorrectlyUpdatesTheGasLimit() public {
    uint32 newGasLimit = 2000000;
    assertEq(s_validator.getGasLimit(), INIT_GAS_LIMIT);
    s_validator.setGasLimit(newGasLimit);
    assertEq(s_validator.getGasLimit(), newGasLimit);
  }

  /// @notice it reverts if called by account with no access
  function test_RevertsIfCalledByAnAccountWithNoAccess() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    s_validator.validate(0, 0, 1, 1);
  }

  /// @notice it posts sequencer status when there is not status change
  function test_PostSequencerStatusWhenThereIsNotStatusChange() public {
    // Gives access to the s_eoaValidator
    s_validator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emitExpectedSentMessageEvent(address(s_validator), false, futureTimestampInSeconds);

    // Runs the function (which produces the event to test)
    s_validator.validate(0, 0, 0, 0);
  }

  /// @notice it post sequencer offline
  function test_PostSequencerOffline() public {
    // Gives access to the s_eoaValidator
    s_validator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 10000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emitExpectedSentMessageEvent(address(s_validator), true, futureTimestampInSeconds);

    // Runs the function (which produces the event to test)
    s_validator.validate(0, 0, 1, 1);
  }

}

