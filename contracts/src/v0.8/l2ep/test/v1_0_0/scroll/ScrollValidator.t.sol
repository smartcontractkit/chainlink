// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ISequencerUptimeFeed} from "../../../interfaces/ISequencerUptimeFeed.sol";

import {MockScrollL1CrossDomainMessenger} from "../../mocks/scroll/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/scroll/MockScrollL2CrossDomainMessenger.sol";
import {MockScrollL1MessageQueue} from "../../mocks/scroll/MockScrollL1MessageQueue.sol";
import {ScrollSequencerUptimeFeed} from "../../../scroll/ScrollSequencerUptimeFeed.sol";
import {ScrollValidator} from "../../../scroll/ScrollValidator.sol";
import {BaseValidator} from "../../../base/BaseValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ScrollValidator_Setup is L2EPTest {
  /// Helper constants
  address internal immutable L2_SEQ_STATUS_RECORDER_ADDRESS = makeAddr("L2_SEQ_STATUS_RECORDER_ADDRESS");
  uint32 internal constant INIT_GAS_LIMIT = 1900000;

  /// L2EP contracts
  MockScrollL1CrossDomainMessenger internal s_mockScrollL1CrossDomainMessenger;
  MockScrollL2CrossDomainMessenger internal s_mockScrollL2CrossDomainMessenger;
  MockScrollL1MessageQueue internal s_mockScrollL1MessageQueue;
  ScrollSequencerUptimeFeed internal s_scrollSequencerUptimeFeed;
  ScrollValidator internal s_scrollValidator;

  /// https://github.com/scroll-tech/scroll/blob/03089eaeee1193ff44c532c7038611ae123e7ef3/contracts/src/libraries/IScrollMessenger.sol#L22
  event SentMessage(
    address indexed sender,
    address indexed target,
    uint256 value,
    uint256 messageNonce,
    uint256 gasLimit,
    bytes message
  );

  /// Setup
  function setUp() public {
    s_mockScrollL1CrossDomainMessenger = new MockScrollL1CrossDomainMessenger();
    s_mockScrollL2CrossDomainMessenger = new MockScrollL2CrossDomainMessenger();
    s_mockScrollL1MessageQueue = new MockScrollL1MessageQueue();

    s_scrollSequencerUptimeFeed = new ScrollSequencerUptimeFeed(
      address(s_mockScrollL1CrossDomainMessenger),
      address(s_mockScrollL2CrossDomainMessenger),
      true
    );

    s_scrollValidator = new ScrollValidator(
      address(s_mockScrollL1CrossDomainMessenger),
      address(s_scrollSequencerUptimeFeed),
      address(s_mockScrollL1MessageQueue),
      INIT_GAS_LIMIT
    );
  }
}

contract ScrollValidator_Constructor is ScrollValidator_Setup {
  /// @notice Reverts when L1 message queue address is invalid
  function test_Constructor_RevertWhen_InvalidL1MessageQueueAddress() public {
    vm.startPrank(s_l1OwnerAddr);

    vm.expectRevert("Invalid L1 message queue address");
    new ScrollValidator(
      address(s_mockScrollL1CrossDomainMessenger),
      address(s_scrollSequencerUptimeFeed),
      address(0),
      INIT_GAS_LIMIT
    );
  }
}

contract ScrollValidator_Validate is ScrollValidator_Setup {
  /// @notice Reverts if called by an account with no access
  function test_Validate_RevertWhen_CalledByAccountWithNoAccess() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    s_scrollValidator.validate(0, 0, 1, 1);
  }

  /// @notice Posts sequencer status when there is no status change
  function test_Validate_PostSequencerStatus_NoStatusChange() public {
    // Gives access to the s_eoaValidator
    s_scrollValidator.addAccess(s_eoaValidator);

    uint256 previousRoundId = 0;
    int256 previousAnswer = 0;
    uint256 currentRoundId = 0;
    int256 currentAnswer = 0;
    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 5000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      address(s_scrollValidator), // sender
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      0, // value
      0, // nonce
      INIT_GAS_LIMIT, // gas limit
      abi.encodeWithSelector(ISequencerUptimeFeed.updateStatus.selector, false, futureTimestampInSeconds) // message
    );
    vm.expectEmit(address(s_scrollValidator));
    emit BaseValidator.ValidatedStatus(previousRoundId, previousAnswer, currentRoundId, currentAnswer);
    // Runs the function (which produces the event to test)
    s_scrollValidator.validate(previousRoundId, previousAnswer, currentRoundId, currentAnswer);
  }

  /// @notice Posts sequencer offline status
  function test_Validate_PostSequencerOffline() public {
    // Gives access to the s_eoaValidator
    s_scrollValidator.addAccess(s_eoaValidator);

    uint256 previousRoundId = 0;
    int256 previousAnswer = 0;
    uint256 currentRoundId = 1;
    int256 currentAnswer = 1;
    // Sets block.timestamp to a later date
    uint256 futureTimestampInSeconds = block.timestamp + 10000;
    vm.startPrank(s_eoaValidator);
    vm.warp(futureTimestampInSeconds);

    // Sets up the expected event data
    vm.expectEmit(false, false, false, true);
    emit SentMessage(
      address(s_scrollValidator), // sender
      L2_SEQ_STATUS_RECORDER_ADDRESS, // target
      0, // value
      0, // nonce
      INIT_GAS_LIMIT, // gas limit
      abi.encodeWithSelector(ISequencerUptimeFeed.updateStatus.selector, true, futureTimestampInSeconds) // message
    );

    vm.expectEmit(address(s_scrollValidator));
    emit BaseValidator.ValidatedStatus(previousRoundId, previousAnswer, currentRoundId, currentAnswer);
    // Runs the function (which produces the event to test)
    s_scrollValidator.validate(previousRoundId, previousAnswer, currentRoundId, currentAnswer);
  }
}
