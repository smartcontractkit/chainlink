// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ISequencerUptimeFeed} from "../../../interfaces/ISequencerUptimeFeed.sol";

import {MockOptimismL1CrossDomainMessenger} from "../../mocks/MockOptimismL1CrossDomainMessenger.sol";
import {MockOptimismL2CrossDomainMessenger} from "../../mocks/MockOptimismL2CrossDomainMessenger.sol";
import {OptimismSequencerUptimeFeed} from "../../../optimism/OptimismSequencerUptimeFeed.sol";
import {OptimismValidator} from "../../../optimism/OptimismValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract OptimismValidator_Setup is L2EPTest {
  /// Helper constants
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = 0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b;
  uint32 internal constant INIT_GAS_LIMIT = 1900000;

  /// L2EP contracts
  MockOptimismL1CrossDomainMessenger internal s_mockOptimismL1CrossDomainMessenger;
  MockOptimismL2CrossDomainMessenger internal s_mockOptimismL2CrossDomainMessenger;
  OptimismSequencerUptimeFeed internal s_optimismSequencerUptimeFeed;
  OptimismValidator internal s_optimismValidator;

  /// Events
  event SentMessage(address indexed target, address sender, bytes message, uint256 messageNonce, uint256 gasLimit);

  /// Setup
  function setUp() public {
    s_mockOptimismL1CrossDomainMessenger = new MockOptimismL1CrossDomainMessenger();
    s_mockOptimismL2CrossDomainMessenger = new MockOptimismL2CrossDomainMessenger();

    s_optimismSequencerUptimeFeed = new OptimismSequencerUptimeFeed(
      address(s_mockOptimismL1CrossDomainMessenger),
      address(s_mockOptimismL2CrossDomainMessenger),
      true
    );

    s_optimismValidator = new OptimismValidator(
      address(s_mockOptimismL1CrossDomainMessenger),
      address(s_optimismSequencerUptimeFeed),
      INIT_GAS_LIMIT
    );
  }
}

contract OptimismValidator_Validate is OptimismValidator_Setup {
  /// @notice it reverts if called by account with no access
  function test_Validate_RevertWhen_CalledByAccountWithNoAccess() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    s_optimismValidator.validate(0, 0, 1, 1);
  }

  /// @notice it posts sequencer status when there is no status change
  function test_Validate_PostSequencerStatus_NoStatusChange() public {
    // Gives access to the s_eoaValidator
    s_optimismValidator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      address(s_optimismValidator), // sender
      abi.encodeWithSelector(ISequencerUptimeFeed.updateStatus.selector, false, futureTimestampInSeconds), // message
      0, // nonce
      INIT_GAS_LIMIT // gas limit
    );

    // Runs the function (which produces the event to test)
    s_optimismValidator.validate(0, 0, 0, 0);
  }

  /// @notice it posts sequencer offline
  function test_Validate_PostSequencerOffline() public {
    // Gives access to the s_eoaValidator
    s_optimismValidator.addAccess(s_eoaValidator);

    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 10000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      address(s_optimismValidator), // sender
      abi.encodeWithSelector(ISequencerUptimeFeed.updateStatus.selector, true, futureTimestampInSeconds), // message
      0, // nonce
      INIT_GAS_LIMIT // gas limit
    );

    // Runs the function (which produces the event to test)
    s_optimismValidator.validate(0, 0, 1, 1);
  }
}
