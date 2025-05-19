// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Vm} from "forge-std/Test.sol";
import {AddressAliasHelper} from "../../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import {BaseSequencerUptimeFeed} from "../../../base/BaseSequencerUptimeFeed.sol";
import {MockBaseSequencerUptimeFeed} from "../../../test/mocks/MockBaseSequencerUptimeFeed.sol";
import {FeedConsumer} from "../../FeedConsumer.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract BaseSequencerUptimeFeed_Setup is L2EPTest {
  /// Helper Variables
  address internal s_aliasedL1OwnerAddress = AddressAliasHelper.applyL1ToL2Alias(s_l1OwnerAddr);

  /// L2EP contracts
  BaseSequencerUptimeFeed internal s_sequencerUptimeFeed;

  /// Events
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);
  event RoundUpdated(int256 status, uint64 updatedAt);
  event L1SenderTransferred(address indexed from, address indexed to);

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_sequencerUptimeFeed = new MockBaseSequencerUptimeFeed(s_l1OwnerAddr, false, true);
  }
}

contract BaseSequencerUptimeFeed_Constructor is BaseSequencerUptimeFeed_Setup {
  /// @notice Tests initial state of the contract
  function test_Constructor_InitialState() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Checks L1 sender
    address actualL1Addr = s_sequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_sequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract BaseSequencerUptimeFeed_transferL1Sender is BaseSequencerUptimeFeed_Setup {
  /// @notice Tests transferring L1 sender
  function test_transferL1Sender_CorrectlyTransfersL1Sender() public {
    address initialSender = address(0);
    address newSender = makeAddr("newSender");

    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    MockBaseSequencerUptimeFeed sequencerUptimeFeed = new MockBaseSequencerUptimeFeed(initialSender, false, true);

    assertEq(sequencerUptimeFeed.l1Sender(), initialSender);

    // Transfers the L1 sender
    vm.expectEmit();
    emit L1SenderTransferred(initialSender, newSender);
    sequencerUptimeFeed.transferL1Sender(newSender);
    assertEq(sequencerUptimeFeed.l1Sender(), newSender);

    vm.recordLogs();
    // Transfers to the same L1 sender should not emit an event
    sequencerUptimeFeed.transferL1Sender(newSender);
    assertEq(vm.getRecordedLogs().length, 0);
  }

  /// @notice Reverts if called by an unauthorized account
  function test_transferL1Sender_RevertWhen_CalledByUnauthorizedAccount() public {
    address newSender = makeAddr("newSender");

    // Sets msg.sender and tx.origin to an unauthorized address
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    vm.expectRevert("Only callable by owner");
    s_sequencerUptimeFeed.transferL1Sender(newSender);
  }
}

contract BaseSequencerUptimeFeed_UpdateStatus is BaseSequencerUptimeFeed_Setup {
  /// @notice Reverts if called by an unauthorized account
  function test_updateStatus_RevertWhen_NotL2CrossDomainMessengerAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    BaseSequencerUptimeFeed s_sequencerUptimeFeedFailSenderCheck = new MockBaseSequencerUptimeFeed(
      s_l1OwnerAddr,
      false,
      false
    );

    // Tries to update the status from an unauthorized account
    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_sequencerUptimeFeedFailSenderCheck.updateStatus(true, uint64(1));
  }

  /// @notice Updates status when status has not changed and incoming timestamp is the same as latest
  function test_updateStatus_UpdateWhen_NoStatusChangeSameTimestamp() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Fetches the latest timestamp
    uint256 timestamp = s_sequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_sequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_sequencerUptimeFeed.latestRound(), 2);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Stores the current round data before updating it
    (
      uint80 roundIdBeforeUpdate,
      int256 answerBeforeUpdate,
      uint256 startedAtBeforeUpdate,
      ,
      uint80 answeredInRoundBeforeUpdate
    ) = s_sequencerUptimeFeed.latestRoundData();

    // Submit another status update with the same status
    vm.expectEmit();
    emit RoundUpdated(1, uint64(block.timestamp));
    s_sequencerUptimeFeed.updateStatus(true, uint64(timestamp + 200));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_sequencerUptimeFeed.latestRound(), 2);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Stores the current round data after updating it
    (
      uint80 roundIdAfterUpdate,
      int256 answerAfterUpdate,
      uint256 startedAtAfterUpdate,
      uint256 updatedAtAfterUpdate,
      uint80 answeredInRoundAfterUpdate
    ) = s_sequencerUptimeFeed.latestRoundData();

    // Verifies the latest round data has been properly updated
    assertEq(roundIdAfterUpdate, roundIdBeforeUpdate);
    assertEq(answerAfterUpdate, answerBeforeUpdate);
    assertEq(startedAtAfterUpdate, startedAtBeforeUpdate);
    assertEq(answeredInRoundAfterUpdate, answeredInRoundBeforeUpdate);
    assertEq(updatedAtAfterUpdate, block.timestamp);
  }

  /// @notice Updates status when status has changed and incoming timestamp is newer than the latest
  function test_updateStatus_UpdateWhen_StatusChangeAndTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Submits a status update
    uint256 timestamp = s_sequencerUptimeFeed.latestTimestamp();
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_sequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_sequencerUptimeFeed.latestRound(), 2);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, newer timestamp should update
    timestamp = timestamp + 200;
    vm.expectEmit();
    emit AnswerUpdated(0, 3, timestamp);
    s_sequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_sequencerUptimeFeed.latestRound(), 3);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));
  }

  /// @notice Updates status when status has changed and incoming timestamp is the same as latest
  function test_updateStatus_UpdateWhen_StatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Fetches the latest timestamp
    uint256 timestamp = s_sequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_sequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_sequencerUptimeFeed.latestRound(), 2);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, same timestamp should update
    vm.expectEmit();
    emit AnswerUpdated(0, 3, timestamp);
    s_sequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_sequencerUptimeFeed.latestRound(), 3);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));
  }

  /// @notice Ignores out-of-order updates
  function test_updateStatus_IgnoreOutOfOrderUpdates() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Submits a status update
    uint256 timestamp = s_sequencerUptimeFeed.latestTimestamp() + 10000;
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_sequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_sequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_sequencerUptimeFeed.latestRound(), 2);
    assertEq(s_sequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Update with different status, but stale timestamp, should be ignored
    timestamp = timestamp - 1000;
    vm.expectEmit(false, false, false, false);
    emit UpdateIgnored(true, 0, true, 0); // arguments are dummy values

    vm.recordLogs();

    // Tries to update with stale timestamp
    s_sequencerUptimeFeed.updateStatus(false, uint64(timestamp));

    Vm.Log[] memory entries = vm.getRecordedLogs();

    assertEq(entries.length, 1);
    assertEq(entries[0].topics[0], keccak256("UpdateIgnored(bool,uint64,bool,uint64)"));
  }
}

contract BaseSequencerUptimeFeed_AggregatorV3Interface is BaseSequencerUptimeFeed_Setup {
  /// @notice Returns valid answer from getRoundData and latestRoundData
  function test_AggregatorV3Interface_ReturnsValidData() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Defines helper variables
    uint80 roundId;
    int256 answer;
    uint256 startedAt;
    uint256 updatedAt;
    uint80 answeredInRound;

    // Checks initial state
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_sequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Submits status update with different status and newer timestamp, should update
    uint256 timestamp = startedAt + 1000;
    s_sequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_sequencerUptimeFeed.getRoundData(2);
    assertEq(roundId, 2);
    assertEq(answer, 1);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, timestamp);
    assertLe(updatedAt, startedAt);

    // Saves round 2 data
    uint80 roundId2 = roundId;
    int256 answer2 = answer;
    uint256 startedAt2 = startedAt;
    uint256 updatedAt2 = updatedAt;
    uint80 answeredInRound2 = answeredInRound;

    // Checks that last round is still returning the correct data
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_sequencerUptimeFeed.getRoundData(1);
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Assert latestRoundData corresponds to latest round id
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_sequencerUptimeFeed.latestRoundData();
    assertEq(roundId2, roundId);
    assertEq(answer2, answer);
    assertEq(startedAt2, startedAt);
    assertEq(updatedAt2, updatedAt);
    assertEq(answeredInRound2, answeredInRound);
  }

  /// @notice Reverts when getRoundData is called for a round that does not exist yet
  function test_getRoundData_RevertWhen_RoundDoesNotExist() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(BaseSequencerUptimeFeed.NoDataPresent.selector);
    s_sequencerUptimeFeed.getRoundData(2);
  }

  /// @notice Returns the getAnswer for the latest round
  function test_getAnswer_ReturnsValidAnswer() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    uint256 startedAt;
    (, , startedAt, , ) = s_sequencerUptimeFeed.latestRoundData();

    s_sequencerUptimeFeed.updateStatus(true, uint64(startedAt + 1000));

    assertEq(0, s_sequencerUptimeFeed.getAnswer(1));
  }

  /// @notice Reverts when getAnswer is called for a round that does not exist yet
  function test_getAnswer_RevertWhen_RoundDoesNotExist() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(BaseSequencerUptimeFeed.NoDataPresent.selector);
    s_sequencerUptimeFeed.getAnswer(2);
  }

  /// @notice Returns the getTimestamp for the latest round
  function test_getTimestamp_ReturnsValidTimestamp() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    uint256 startedAt;
    (, , startedAt, , ) = s_sequencerUptimeFeed.latestRoundData();

    s_sequencerUptimeFeed.updateStatus(true, uint64(startedAt + 1000));

    assertEq(startedAt, s_sequencerUptimeFeed.getTimestamp(1));
  }

  /// @notice Reverts when getTimestamp is called for a round that does not exist yet
  function test_getTimestamp_RevertWhen_RoundDoesNotExist() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(BaseSequencerUptimeFeed.NoDataPresent.selector);
    s_sequencerUptimeFeed.getTimestamp(2);
  }
}

contract BaseSequencerUptimeFeed_ProtectReadsOnAggregatorV2V3InterfaceFunctions is BaseSequencerUptimeFeed_Setup {
  /// @notice Disallows reads on AggregatorV2V3Interface functions when consuming contract is not whitelisted
  function test_ProtectReads_DisallowWhen_NotWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_sequencerUptimeFeed));

    // Sanity - consumer is not whitelisted
    assertEq(s_sequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_sequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), false);

    // Asserts reads are not possible from consuming contract
    vm.expectRevert("No access");
    feedConsumer.latestAnswer();
    vm.expectRevert("No access");
    feedConsumer.latestRoundData();
  }

  /// @notice Allows reads on AggregatorV2V3Interface functions when consuming contract is whitelisted
  function test_ProtectReads_AllowWhen_Whitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_sequencerUptimeFeed));

    // Whitelist consumer
    s_sequencerUptimeFeed.addAccess(address(feedConsumer));

    // Sanity - consumer is whitelisted
    assertEq(s_sequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_sequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), true);

    // Asserts reads are possible from consuming contract
    (uint80 roundId, int256 answer, , , ) = feedConsumer.latestRoundData();
    assertEq(feedConsumer.latestAnswer(), 0);
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}
