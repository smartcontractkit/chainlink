// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {AddressAliasHelper} from "../../../../vendor/arb-bridge-eth/v0.8.0-custom/contracts/libraries/AddressAliasHelper.sol";
import {ZKSyncSequencerUptimeFeed} from "../../../dev/zksync/ZKSyncSequencerUptimeFeed.sol";
import {BaseSequencerUptimeFeed} from "../../../dev/shared/BaseSequencerUptimeFeed.sol";
import {FeedConsumer} from "../../../../tests/FeedConsumer.sol";
import {L2EPTest} from "../L2EPTest.t.sol";

contract ZKSyncSequencerUptimeFeedTest is L2EPTest {
  /// Helper Variables
  address internal s_aliasedL1OwnerAddress = AddressAliasHelper.applyL1ToL2Alias(s_l1OwnerAddr);

  /// L2EP contracts
  ZKSyncSequencerUptimeFeed internal s_zksyncSequencerUptimeFeed;

  /// Events
  event UpdateIgnored(bool latestStatus, uint64 latestTimestamp, bool incomingStatus, uint64 incomingTimestamp);
  event AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt);
  event RoundUpdated(int256 status, uint64 updatedAt);

  /// Setup
  function setUp() public {
    // Deploys contracts
    s_zksyncSequencerUptimeFeed = new ZKSyncSequencerUptimeFeed(s_l1OwnerAddr, false);
  }
}

contract ZKSyncSequencerUptimeFeed_Constructor is ZKSyncSequencerUptimeFeedTest {
  /// @notice it should have been deployed with the correct initial state
  function test_InitialState() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Checks L1 sender
    address actualL1Addr = s_zksyncSequencerUptimeFeed.l1Sender();
    assertEq(actualL1Addr, s_l1OwnerAddr);

    // Checks latest round data
    (uint80 roundId, int256 answer, , , ) = s_zksyncSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}

contract ZKSyncSequencerUptimeFeed_UpdateStatus is ZKSyncSequencerUptimeFeedTest {
  /// @notice it should revert if called by an unauthorized account
  function test_RevertIfNotL2CrossDomainMessengerAddr() public {
    // Sets msg.sender and tx.origin to an unauthorized address
    vm.startPrank(s_strangerAddr, s_strangerAddr);

    // Tries to update the status from an unauthorized account
    vm.expectRevert(BaseSequencerUptimeFeed.InvalidSender.selector);
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(1));
  }

  /// @notice it should update status when status has not changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenNoChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Fetches the latest timestamp
    uint256 timestamp = s_zksyncSequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Stores the current round data before updating it
    (
      uint80 roundIdBeforeUpdate,
      int256 answerBeforeUpdate,
      uint256 startedAtBeforeUpdate,
      ,
      uint80 answeredInRoundBeforeUpdate
    ) = s_zksyncSequencerUptimeFeed.latestRoundData();

    // Submit another status update with the same status
    vm.expectEmit();
    emit RoundUpdated(1, uint64(block.timestamp));
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(timestamp + 200));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Stores the current round data after updating it
    (
      uint80 roundIdAfterUpdate,
      int256 answerAfterUpdate,
      uint256 startedAtAfterUpdate,
      uint256 updatedAtAfterUpdate,
      uint80 answeredInRoundAfterUpdate
    ) = s_zksyncSequencerUptimeFeed.latestRoundData();

    // Verifies the latest round data has been properly updated
    assertEq(roundIdAfterUpdate, roundIdBeforeUpdate);
    assertEq(answerAfterUpdate, answerBeforeUpdate);
    assertEq(startedAtAfterUpdate, startedAtBeforeUpdate);
    assertEq(answeredInRoundAfterUpdate, answeredInRoundBeforeUpdate);
    assertEq(updatedAtAfterUpdate, block.timestamp);
  }

  /// @notice it should update status when status has changed and incoming timestamp is newer than the latest
  function test_UpdateStatusWhenStatusChangeAndTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Submits a status update
    uint256 timestamp = s_zksyncSequencerUptimeFeed.latestTimestamp();
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, newer timestamp should update
    timestamp = timestamp + 200;
    vm.expectEmit();
    emit AnswerUpdated(0, 3, timestamp);
    s_zksyncSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));
  }

  /// @notice it should update status when status has changed and incoming timestamp is the same as latest
  function test_UpdateStatusWhenStatusChangeAndNoTimeChange() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Fetches the latest timestamp
    uint256 timestamp = s_zksyncSequencerUptimeFeed.latestTimestamp();

    // Submits a status update
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Submit another status update, different status, same timestamp should update
    vm.expectEmit();
    emit AnswerUpdated(0, 3, timestamp);
    s_zksyncSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 0);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));
  }

  /// @notice it should ignore out-of-order updates
  function test_IgnoreOutOfOrderUpdates() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Submits a status update
    uint256 timestamp = s_zksyncSequencerUptimeFeed.latestTimestamp() + 10000;
    vm.expectEmit();
    emit AnswerUpdated(1, 2, timestamp);
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    assertEq(s_zksyncSequencerUptimeFeed.latestAnswer(), 1);
    assertEq(s_zksyncSequencerUptimeFeed.latestTimestamp(), uint64(timestamp));

    // Update with different status, but stale timestamp, should be ignored
    timestamp = timestamp - 1000;
    vm.expectEmit(false, false, false, false);
    emit UpdateIgnored(true, 0, true, 0); // arguments are dummy values
    // TODO: how can we check that an AnswerUpdated event was NOT emitted
    s_zksyncSequencerUptimeFeed.updateStatus(false, uint64(timestamp));
  }
}

contract ZKSyncSequencerUptimeFeed_AggregatorV3Interface is ZKSyncSequencerUptimeFeedTest {
  /// @notice it should return valid answer from getRoundData and latestRoundData
  function test_AggregatorV3Interface() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_aliasedL1OwnerAddress, s_aliasedL1OwnerAddress);

    // Defines helper variables
    uint80 roundId;
    int256 answer;
    uint256 startedAt;
    uint256 updatedAt;
    uint80 answeredInRound;

    // Checks initial state
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_zksyncSequencerUptimeFeed.latestRoundData();
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Submits status update with different status and newer timestamp, should update
    uint256 timestamp = startedAt + 1000;
    s_zksyncSequencerUptimeFeed.updateStatus(true, uint64(timestamp));
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_zksyncSequencerUptimeFeed.getRoundData(2);
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
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_zksyncSequencerUptimeFeed.getRoundData(1);
    assertEq(roundId, 1);
    assertEq(answer, 0);
    assertEq(answeredInRound, roundId);
    assertEq(startedAt, updatedAt);

    // Assert latestRoundData corresponds to latest round id
    (roundId, answer, startedAt, updatedAt, answeredInRound) = s_zksyncSequencerUptimeFeed.latestRoundData();
    assertEq(roundId2, roundId);
    assertEq(answer2, answer);
    assertEq(startedAt2, startedAt);
    assertEq(updatedAt2, updatedAt);
    assertEq(answeredInRound2, answeredInRound);
  }

  /// @notice it should revert from #getRoundData when round does not yet exist (future roundId)
  function test_RevertGetRoundDataWhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(BaseSequencerUptimeFeed.NoDataPresent.selector);
    s_zksyncSequencerUptimeFeed.getRoundData(2);
  }

  /// @notice it should revert from #getAnswer when round does not yet exist (future roundId)
  function test_RevertGetAnswerWhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(BaseSequencerUptimeFeed.NoDataPresent.selector);
    s_zksyncSequencerUptimeFeed.getAnswer(2);
  }

  /// @notice it should revert from #getTimestamp when round does not yet exist (future roundId)
  function test_RevertGetTimestampWhenRoundDoesNotExistYet() public {
    // Sets msg.sender and tx.origin to a valid address
    vm.startPrank(s_l1OwnerAddr, s_l1OwnerAddr);

    // Gets data from a round that has not happened yet
    vm.expectRevert(BaseSequencerUptimeFeed.NoDataPresent.selector);
    s_zksyncSequencerUptimeFeed.getTimestamp(2);
  }
}

contract ZKSyncSequencerUptimeFeed_ProtectReadsOnAggregatorV2V3InterfaceFunctions is ZKSyncSequencerUptimeFeedTest {
  /// @notice it should disallow reads on AggregatorV2V3Interface functions when consuming contract is not whitelisted
  function test_AggregatorV2V3InterfaceDisallowReadsIfConsumingContractIsNotWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_zksyncSequencerUptimeFeed));

    // Sanity - consumer is not whitelisted
    assertEq(s_zksyncSequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_zksyncSequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), false);

    // Asserts reads are not possible from consuming contract
    vm.expectRevert("No access");
    feedConsumer.latestAnswer();
    vm.expectRevert("No access");
    feedConsumer.latestRoundData();
  }

  /// @notice it should allow reads on AggregatorV2V3Interface functions when consuming contract is whitelisted
  function test_AggregatorV2V3InterfaceAllowReadsIfConsumingContractIsWhitelisted() public {
    // Deploys a FeedConsumer contract
    FeedConsumer feedConsumer = new FeedConsumer(address(s_zksyncSequencerUptimeFeed));

    // Whitelist consumer
    s_zksyncSequencerUptimeFeed.addAccess(address(feedConsumer));

    // Sanity - consumer is whitelisted
    assertEq(s_zksyncSequencerUptimeFeed.checkEnabled(), true);
    assertEq(s_zksyncSequencerUptimeFeed.hasAccess(address(feedConsumer), abi.encode("")), true);

    // Asserts reads are possible from consuming contract
    (uint80 roundId, int256 answer, , , ) = feedConsumer.latestRoundData();
    assertEq(feedConsumer.latestAnswer(), 0);
    assertEq(roundId, 1);
    assertEq(answer, 0);
  }
}
