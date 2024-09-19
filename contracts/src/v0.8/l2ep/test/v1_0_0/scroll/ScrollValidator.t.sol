// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {ISequencerUptimeFeed} from "../../../dev/interfaces/ISequencerUptimeFeed.sol";

import {MockScrollL1CrossDomainMessenger} from "../../mocks/scroll/MockScrollL1CrossDomainMessenger.sol";
import {MockScrollL2CrossDomainMessenger} from "../../mocks/scroll/MockScrollL2CrossDomainMessenger.sol";
import {MockScrollL1MessageQueue} from "../../mocks/scroll/MockScrollL1MessageQueue.sol";
import {ScrollSequencerUptimeFeed} from "../../../dev/scroll/ScrollSequencerUptimeFeed.sol";
import {ScrollValidator} from "../../../dev/scroll/ScrollValidator.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ScrollValidatorTest is L2EPTest {
  /// Helper constants
  address internal constant L2_SEQ_STATUS_RECORDER_ADDRESS = 0x491B1dDA0A8fa069bbC1125133A975BF4e85a91b;
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

contract ScrollValidator_SetGasLimit is ScrollValidatorTest {
  /// @notice it correctly updates the gas limit
  function test_CorrectlyUpdatesTheGasLimit() public {
    uint32 newGasLimit = 2000000;
    assertEq(s_scrollValidator.getGasLimit(), INIT_GAS_LIMIT);
    s_scrollValidator.setGasLimit(newGasLimit);
    assertEq(s_scrollValidator.getGasLimit(), newGasLimit);
  }
}

contract ScrollValidator_Validate is ScrollValidatorTest {
  /// @notice it reverts if called by account with no access
  function test_RevertsIfCalledByAnAccountWithNoAccess() public {
    vm.startPrank(s_strangerAddr);
    vm.expectRevert("No access");
    s_scrollValidator.validate(0, 0, 1, 1);
  }

  /// @notice it posts sequencer status when there is not status change
  function test_PostSequencerStatusWhenThereIsNotStatusChange() public {
    // Gives access to the s_eoaValidator
    s_scrollValidator.addAccess(s_eoaValidator);

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

    // Runs the function (which produces the event to test)
    s_scrollValidator.validate(0, 0, 0, 0);
  }

  /// @notice it post sequencer offline
  function test_PostSequencerOffline() public {
    // Gives access to the s_eoaValidator
    s_scrollValidator.addAccess(s_eoaValidator);

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

    // Runs the function (which produces the event to test)
    s_scrollValidator.validate(0, 0, 1, 1);
  }
}
